package archivist

import "github.com/floodedrealms/borderland-keep/guardsman"

type Page struct {
	User guardsman.WebUser
}

type UserDetails struct {
	LoggedIn        bool
	UserId          string
	UserName        string
	PrefferredStyle string
}

func (u UserDetails) IsAuthenticated() bool {
	return u.LoggedIn
}

func (u UserDetails) HasEditAccessToResource(r string, id int) bool {
	return true
}
