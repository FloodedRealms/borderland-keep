package repository

import "fmt"

type RepositoryError struct {
	AttemptedOperation string
	Resource           string
	InternalErr        error
}

func (e RepositoryError) Error() string {
	if e.InternalErr != nil {
		return fmt.Sprintf("The %s on %s failed, the database threw the following error:\n\t", e.AttemptedOperation, e.Resource, e.InternalErr.Error())
	}
	return fmt.Sprintf("The %s on %s failed.", e.AttemptedOperation, e.Resource)
}

type ParameterNotProvidedError struct {
	ExpectedMin   int
	Provided      int
	ParameterName string
}

func (e ParameterNotProvidedError) Error() string {
	return fmt.Sprintf("Not enough parameters provide for %s. expected at least %d, received %d", e.ParameterName, e.ExpectedMin, e.Provided)
}

type ParameterValueMismatch struct {
	ParamaterizedStatementName string
	Operation                  string
	ParametersProvided         int
	ValuesProvided             int
}

func (e ParameterValueMismatch) Error() string {
	return fmt.Sprintf("Mismatch on number of parameters and values for %s for attempted %s. There were %d parameters provided and %d values for those parameters", e.ParamaterizedStatementName, e.Operation, e.ParametersProvided, e.ValuesProvided)
}
