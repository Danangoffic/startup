package helper

import "github.com/go-playground/validator/v10"

// Response helper json
type Response struct {
	// `json:"meta"` untuk renaming dari "Meta" menjadi "meta" pada response json
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data"`
}

type Meta struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Status  string `json:"status"`
}

func APIResponse(message string, code int, status string, data interface{}) Response {
	meta := Meta{
		Message: message,
		Code:    code,
		Status:  status,
	}

	response := Response{
		Meta: meta,
		Data: data,
	}

	return response
}

//to format listed validation error
func FormatValidationError(err error) []string {
	var errors []string
	for _, e := range err.(validator.ValidationErrors) {
		errors = append(errors, e.Error())
	}
	return errors
}
