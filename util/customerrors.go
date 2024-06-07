package util

import "errors"

func NotYetImplmented() error {
	return errors.New("not yet implemented")
}

func RequiredParameterNotProvided() error {
	return errors.New("A required parameter was not provided.")
}
