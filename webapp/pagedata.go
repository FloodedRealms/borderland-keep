package webapp

type PageData interface {
	IsAuthenticated() bool
	HasEditAccessToResource(resourceType string, resourceId int) bool
}

type FormData interface {
	Validate() (bool, []interface{})
}

type UserDetails struct {
	UserName        string
	PrefferredStyle string
}

func (u UserDetails) IsAuthenticated() bool {
	return true
}

func (u UserDetails) HasEditAccessToResource(r string, id int) bool {
	return true
}
