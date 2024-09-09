package renderer

type Form interface {
	FormName() string
	GetInputFields() []InputField
	GetErrors() []FormInputError
	Validate() (bool, Form, error)
}

type InputField struct {
	InputType  string
	FieldName  string
	FieldLabel string
	Options    []SelectOptions
}

type SelectOptions struct {
	Options []string
	Values  []string
}

func NewTextInput(fieldname, fieldLabel string) *InputField {
	return &InputField{
		InputType:  "text",
		FieldLabel: fieldLabel,
		FieldName:  fieldname,
	}
}

func NewDropDown(fieldName, fieldLabel string, opts []SelectOptions) *InputField {
	return &InputField{
		InputType:  "select",
		FieldName:  fieldName,
		FieldLabel: fieldLabel,
		Options:    opts,
	}
}

type FormInputError struct {
	Fieldname    string
	ErrorMessage string
}

func (f FormInputError) Error() string {
	return f.ErrorMessage
}
