package guardsman

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/floodedrealms/borderland-keep/internal/services"
	"github.com/floodedrealms/borderland-keep/internal/util"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

const sessionCookie = "session_token"

type PasswordForm struct {
	NameError     string
	PasswordError string
}

func (g Guardsman) DisplayLoginPage(w http.ResponseWriter, r *http.Request) {
	output, err := g.renderer.RenderPage("login.html", PasswordForm{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write([]byte(output))
}

func (g Guardsman) HandleLogin(w http.ResponseWriter, r *http.Request) {
	errors := PasswordForm{}
	isFormValid, username, providedPassword := errors.validateLoginForm(r)
	if !isFormValid {
		output, err := g.renderer.RenderPage("loginForm.html", errors)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write([]byte(output))
		return
	}

	id, hashedPassword, salt, err := g.userService.RetrieveWebUserInformation(username)
	if err != nil {
		_, isNotFoundError := err.(services.UserNotFoundError)
		if isNotFoundError {
			errors.NameError = "That name was not found in our database."
			output, err := g.renderer.RenderPage("loginForm.html", errors)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Write([]byte(output))
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	user := WebUser{
		Id:            id,
		Friendly_name: username,
		Password:      APIKey(*NewPasswordFromDatabase(providedPassword, hashedPassword, salt)),
	}

	_, err = user.Validate()
	if err != nil {
		_, isBadPassword := err.(BadPasswordError)
		if isBadPassword {
			errors.PasswordError = "That password is incorrect"
			output, err := g.renderer.RenderPage("loginForm.html", errors)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Write([]byte(output))
			return

		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	g.StoreSession(w, user)
	w.Header().Set("HX-Redirect", "/")
	http.Redirect(w, r, "/", http.StatusOK)
}

func (p *PasswordForm) validateLoginForm(r *http.Request) (bool, string, string) {
	r.ParseForm()
	valid := true
	username := r.Form["username"]
	password := r.Form["password"]
	if len(username) < 1 || username[0] == "" {
		p.NameError = "Username must be provided"
		valid = false
	}
	if len(password) < 1 || password[0] == "" {
		p.PasswordError = "Password cannot be blank"
		valid = false
	}
	return valid, username[0], password[0]
}

func (g Guardsman) StoreSession(w http.ResponseWriter, u WebUser) error {
	sessionToken := uuid.NewString()
	expiresAt := time.Now().Add(2 * time.Hour)

	// Set the token in the session map, along with the session information
	err := g.userService.StoreSession(sessionToken, u.Id, u.Friendly_name, expiresAt)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	// Finally, we set the client cookie for "session_token" as the session token we just generated
	// we also set an expiry time of 120 seconds
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookie,
		Value:    sessionToken,
		Expires:  expiresAt,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	return nil
}

func (g Guardsman) RefreshSession(w http.ResponseWriter, r *http.Request) {
	u := g.SimpleLoginCheck(r)
	if !u.LoggedIn {
		return
	}
	sessionToken := uuid.NewString()
	expiresAt := time.Now().Add(120 * time.Second)

	err := g.userService.StoreSession(sessionToken, u.Id, u.Friendly_name, expiresAt)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookie,
		Value:    sessionToken,
		Expires:  expiresAt,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

}

func (g Guardsman) Logout(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie(sessionCookie)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	sessionToken := c.Value
	g.userService.DeleteSession(sessionToken)
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookie,
		Value:    "",
		Expires:  time.Now(),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	w.Header().Set("HX-Redirect", "/")
	http.Redirect(w, r, "/", http.StatusOK)
}

type SesssionExiredError struct {
}

func (u SesssionExiredError) Error() string {
	return "Provided Session was expired"
}

// Used for actions don't require authentication, but have additional functionality when logged in.
// Eg: Navigating to a campaign page from the main list will check for login to determine whether or not to display the edit page.
// Nagigating from the /my-campaigns page won't need this, as it will go through the login middle-ware
// Also called when using the refresh token middleware
func (g Guardsman) SimpleLoginCheck(r *http.Request) WebUser {
	user := WebUser{
		LoggedIn: false,
	}
	c, err := r.Cookie(sessionCookie)
	if err != nil {
		return user
	}
	sessionToken := c.Value
	userId, userName, expiry, _, err := g.userService.RetrieveSession(sessionToken)
	expired := expiry.Before(time.Now())

	if err != nil || expired {
		return user
	}
	user.Id = userId
	user.Friendly_name = userName
	user.LoggedIn = true
	return user
}

func (g Guardsman) UserMustBeLoggedIn(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		c, err := r.Cookie(sessionCookie)
		if err != nil {
			if err == http.ErrNoCookie {
				// If the cookie is not set, return an unauthorized status
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			// For any other type of error, return a bad request status
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		sessionToken := c.Value
		_, _, expiry, exists, err := g.userService.RetrieveSession(sessionToken)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if !exists {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if expiry.Before(time.Now()) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (g Guardsman) CheckLoggedIn(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Header.Set(LoggedInHeader, "false")

		c, err := r.Cookie(sessionCookie)
		if err != nil {
			if err == http.ErrNoCookie {
				next.ServeHTTP(w, r)
				return
			}
			// For any other type of error, return a bad request status
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		sessionToken := c.Value
		_, _, expiry, exists, err := g.userService.RetrieveSession(sessionToken)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if !exists {
			next.ServeHTTP(w, r)
			return
		}
		if expiry.Before(time.Now()) {
			next.ServeHTTP(w, r)
			return
		}
		r.Header.Set(LoggedInHeader, "true")
		next.ServeHTTP(w, r)
	})
}

// TODO: This could likely be modified to make only a single database request by adjusting the resources view to have
// expiry information as well
func (g Guardsman) CheckUserhasEditAccessToCampaign(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Header.Set(EditAccessHeader, "false")

		c, err := r.Cookie(sessionCookie)
		if err != nil {
			if err == http.ErrNoCookie {
				next.ServeHTTP(w, r)
				return
			}
			// For any other type of error, return a bad request status
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		sessionToken := c.Value
		campaignId, err := util.ExtractCampaignId(r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		canEdit, err := g.userService.UserhasEditAccessToCampaign(sessionToken, campaignId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if !canEdit {
			next.ServeHTTP(w, r)
			return
		}
		r.Header.Set(EditAccessHeader, "true")
		next.ServeHTTP(w, r)
	})
}

func (g Guardsman) CheckUserhasEditAccessToAdventure(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Header.Set(EditAccessHeader, "false")

		c, err := r.Cookie(sessionCookie)
		if err != nil {
			if err == http.ErrNoCookie {
				next.ServeHTTP(w, r)
				return
			}
			// For any other type of error, return a bad request status
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		sessionToken := c.Value
		adventureId, err := util.ExtractAdventureId(r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		canEdit, err := g.userService.UserhasEditAccessToAdventure(sessionToken, adventureId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if !canEdit {
			next.ServeHTTP(w, r)
			return
		}
		r.Header.Set(EditAccessHeader, "true")
		next.ServeHTTP(w, r)
	})
}

func (g Guardsman) UserLoggedInAndHasEditAccessToCampaign(next http.HandlerFunc) http.HandlerFunc {
	return g.CheckLoggedIn(g.CheckUserhasEditAccessToCampaign(next))
}

func (g Guardsman) UserLoggedInAndHasEditAccessToAdventure(next http.HandlerFunc) http.HandlerFunc {
	return g.CheckLoggedIn(g.CheckUserhasEditAccessToAdventure(next))
}

// TODO: remove JWT stuff. Going to use sessions for now
func GenerateJWT(username string) (string, error) {
	secretKey := os.Getenv("BORDERLAND_KEEP_WATCHWORD")
	if secretKey == "" {
		return "", fmt.Errorf("could not get secret key (watch word)")
	}
	token := jwt.New(jwt.SigningMethodEdDSA)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(10 * time.Minute)
	claims["auhtorized"] = true
	claims["user"] = username
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// TODO: This is a mess and Logrocket is garbage.
// I need to return to this on a few hours of sleep at least and then try to untangle the actual logic.
func VerifyJWT(endpointHandler func(writer http.ResponseWriter, request *http.Request)) http.HandlerFunc {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Header["Token"] != nil {
			token, err := jwt.Parse(request.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
				_, ok := token.Method.(*jwt.SigningMethodECDSA)
				if !ok {
					writer.WriteHeader(http.StatusUnauthorized)
					_, err := writer.Write([]byte("You're Unauthorized!"))
					if err != nil {
						return nil, err

					}
				}
				return "", nil

			})
			if err != nil {
				writer.WriteHeader(http.StatusUnauthorized)
				_, err2 := writer.Write([]byte("You're Unauthorized due to error parsing the JWT"))
				if err2 != nil {
					return
				}
			}
			if token.Valid {
				endpointHandler(writer, request)
			} else {
				writer.WriteHeader(http.StatusUnauthorized)
				_, err := writer.Write([]byte("You're Unauthorized due to invalid token"))
				if err != nil {
					return
				}
			}
		} else {
			writer.WriteHeader(http.StatusUnauthorized)
			_, err := writer.Write([]byte("You're Unauthorized due to No token in the header"))
			if err != nil {
				return
			}
		}
	})
}
