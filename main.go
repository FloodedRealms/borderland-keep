package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"

	"github.com/floodedrealms/borderland-keep/archivist"
	"github.com/floodedrealms/borderland-keep/guardsman"
	"github.com/floodedrealms/borderland-keep/internal/repository"
	"github.com/floodedrealms/borderland-keep/internal/services"
	"github.com/floodedrealms/borderland-keep/internal/util"
	"github.com/floodedrealms/borderland-keep/renderer"
)

const dbName = "keep.db"

func main() {
	flag.Parse()
	flags := flag.Args()

	if len(flags) == 0 {
		fmt.Print("usage: archivist [operation] <args>")
		return
	}

	operation := flag.Arg(0)
	switch operation {
	default:
		fmt.Printf("Unrecgonized Command: %s", operation)
	case "create-user":
		repo, _ := repository.NewSqliteRepo(util.NewLogger(true), dbName)
		userService := services.NewUserService(repo, *util.NewLogger(true))
		guardsman := guardsman.NewGuardsman(repo, *userService, nil, *util.NewLogger(true))
		if len(flags) == 1 {
			fmt.Println("usage: archivist create-user [friendly-name]")
			return
		}
		friendlyName := flag.Arg(1)
		if len(flags) == 3 {
			guardsman.CreateUser("api", friendlyName, "", "")
		} else if len(flags) == 4 {
			email := flag.Arg(2)
			password := flag.Arg(3)
			guardsman.CreateUser("web", friendlyName, email, password)
		} else {
			fmt.Println("make your code better 5head")
		}
	case "unlimit-user":
		repo, _ := repository.NewSqliteRepo(util.NewLogger(true), dbName)
		userService := services.NewUserService(repo, *util.NewLogger(true))
		guardsman := guardsman.NewGuardsman(repo, *userService, nil, *util.NewLogger(true))
		if len(flags) == 1 {
			fmt.Println("usage: archivist unlimit-user [user-id]")
			return
		}
		id := flag.Arg(1)
		err := guardsman.UnlimitUserCampaigns(id)
		if err != nil {
			fmt.Println("User not unlimted")
			fmt.Println(err.Error())
			return
		}
		fmt.Println("User unlimted")

	case "limit-user":
		repo, _ := repository.NewSqliteRepo(util.NewLogger(true), dbName)
		userService := services.NewUserService(repo, *util.NewLogger(true))
		guardsman := guardsman.NewGuardsman(repo, *userService, nil, *util.NewLogger(true))
		if len(flags) == 1 {

			return
		}
		id := flag.Arg(1)
		err := guardsman.LimitUserCampaigns(id)
		if err != nil {
			fmt.Println("User not limted")
			fmt.Println(err.Error())
			return
		}
		fmt.Println("User limted")
	case "server":
		debug := false
		if len(flags) == 2 {
			if flag.Arg(1) == "true" {
				debug = true
			}
		}

		logger := util.NewLogger(debug)

		//Turn on renderer for webpages (will panic if templates are wrong)
		renderer := renderer.NewRenderer()

		sqlRepo, err := repository.NewSqliteRepo(logger, dbName)
		util.CheckErr(err)

		//services
		campaignService := services.NewCampaignService(sqlRepo, logger, context.TODO())
		characterService := services.NewCharacterService(sqlRepo, logger, context.TODO())
		adventureRecordService := services.NewAdventureRecordService(sqlRepo, context.TODO())
		userService := services.NewUserService(sqlRepo, *logger)
		// campaignActionService := services.NewCampaignActionService(sqlRepo)

		//pages
		guardsman := guardsman.NewGuardsman(sqlRepo, *userService, renderer, *logger)
		homePages := archivist.NewHomePage(*renderer, *campaignService, *guardsman)
		campaignPages := archivist.NewCampaignPage(*campaignService, *characterService, *adventureRecordService, *renderer)
		adventurePages := archivist.NewAdventurePage(*adventureRecordService, *characterService, *renderer)
		calculatorPages := archivist.NewSimpleAdventurePage(*renderer)

		//router
		router := http.NewServeMux()

		// USER
		// router.HandleFunc(" GET /user/validate", userApi.ValidateClient)

		// static
		fs := http.FileServer(http.Dir("./static"))
		router.Handle("/static/", http.StripPrefix("/static/", fs))

		// archivist Pages
		router.HandleFunc("/", guardsman.CheckLoggedIn(homePages.Index))
		router.HandleFunc("/guild", guardsman.CheckLoggedIn(homePages.GuildLanding))
		router.HandleFunc("/tavern", guardsman.CheckLoggedIn(homePages.TavernLanding))
		router.HandleFunc("/campaign-list", homePages.LoadNextCampaignSet)
		router.HandleFunc("/about", homePages.About)
		router.HandleFunc("/error", homePages.ErrorPage)
		campaignPages.RegisterRoutes(router, *guardsman)
		adventurePages.RegisterRoutes(router, *guardsman)
		calculatorPages.RegisterRoutes(router, *guardsman)

		//User Pages

		router.HandleFunc("/user/{userId}/campaigns", guardsman.UserMustBeLoggedIn(homePages.MyCampaigns))

		router.HandleFunc("/user/{userId}/campaign/{campaignId}", guardsman.UserMustBeLoggedIn(campaignPages.CampaignPageForUser))
		router.HandleFunc("POST /user/{userId}/campaign", guardsman.UserMustBeLoggedIn(campaignPages.CampaignPageForUser))

		router.HandleFunc("DELETE /campaigns/{campaignId}", guardsman.UserLoggedInAndHasEditAccessToCampaign(campaignPages.CRUDRoutes))
		// guardsmen pages
		guardsman.RegisterRoutes(router)

		// Utility Pages
		router.HandleFunc("POST /lang", updateLangage)

		server := &http.Server{
			Addr:    ":9090",
			Handler: router,
		}
		logger.Print("Listening on 9090")
		for {
			server.ListenAndServe()
			logger.Print("Server crash... attempting restart")
		}

	}

}

func updateLangage(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	newLang := r.Form["lang"]
	http.SetCookie(w, &http.Cookie{
		Name:     "lang",
		Value:    newLang[0],
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	w.Header().Add("HX-Refresh", "true")
	w.WriteHeader(http.StatusOK)
}
