package util

import (
	"errors"
	"fmt"
)

func NotYetImplmented() error {
	return errors.New("not yet implemented")
}

func RequiredParameterNotProvided() error {
	return errors.New("a required parameter was not provided")
}

func NamedParameterNotProvided(paramName string) error {
	return fmt.Errorf("the required parameter %s was not provided", paramName)
}

func UnknownCharacterOperation(given string) error {
	return fmt.Errorf("an unknown operation was given. the valid operations are \"add\" and \"remove.\" the given operation was: %s ", given)
}

func UnableToUpdateAdventure(step, errorMessage string) error {
	return fmt.Errorf("unable to update the adventure on step %s, the following error was thrown: %s", step, errorMessage)
}

func UnableToFindResourceWithId(resource string, id int) error {
	return fmt.Errorf("unable to find %s with id %d", resource, id)
}
