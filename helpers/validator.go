package helpers

import (
	"fmt"
	"strings"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
)

var once sync.Once
var v *validator.Validate

type ErrResponse struct {
	Errors []string `json:"errors"`
}

// GetValidator is responsible for returning a single instance of the validator.
func GetValidator() *validator.Validate {
	once.Do(func() {
		log.Info().Msg("Validator initialized.")
		v = validator.New()
	})

	return v
}

func ToErrorResponse(err error) *ErrResponse {
	if fielderrors, ok := err.(validator.ValidationErrors); ok {
		resp := ErrResponse{
			Errors: make([]string, len(fielderrors)),
		}

		for i, err := range fielderrors {
			switch err.Tag() {
			case "required":
				resp.Errors[i] = fmt.Sprintf("%s is a required field", err.Field())
			case "eqfield":
				resp.Errors[i] = fmt.Sprintf("%s must be equal to %s", err.Field(), err.Param())
			case "email":
				resp.Errors[i] = fmt.Sprintf("%s must be a valid email", err.Field())
			case "oneof":
				resp.Errors[i] = fmt.Sprintf("%s must be one of %s", err.Field(), err.Param())
			case "gt":
				resp.Errors[i] = fmt.Sprintf("%s must be greater than %s", err.Field(), err.Param())
			case "required_if":
				// For required_if, the Param() returns the condition in format "FieldName FieldValue"
				params := strings.SplitN(err.Param(), " ", 2)
				if len(params) == 2 {
					if params[1] == "." {
						// Handle the case where any non-empty value triggers the requirement
						resp.Errors[i] = fmt.Sprintf("%s is required when %s is not empty", err.Field(), params[0])
					} else {
						// Handle the case where a specific value triggers the requirement
						resp.Errors[i] = fmt.Sprintf("%s is required when %s is %s", err.Field(), params[0], params[1])
					}
				} else {
					// Fallback for unexpected param format
					resp.Errors[i] = fmt.Sprintf("%s is required based on other field values", err.Field())
				}
			default:
				resp.Errors[i] = fmt.Sprintf("something wrong on %s; %s", err.Field(), err.Tag())

			}
		}
		return &resp
	}
	return nil
}
