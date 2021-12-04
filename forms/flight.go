package forms

import (
	"encoding/json"
	"fmt"

	"github.com/go-playground/validator/v10"
)

//FlightForm ...
type FlightForm struct{}

//GetFlightForm ...
type GetFlightForm struct {
	ToCode   string `form:"to_code" json:"to_code" binding:"required"`
	FromCode string `form:"from_code" json:"from_code" binding:"required"`
}

//Code ...
func (f FlightForm) Code(tag, field string, errMsg ...string) (message string) {
	switch tag {
	case "required":
		if len(errMsg) == 0 {
			return "Please enter the " + field + " title"
		}
		return errMsg[0]
	default:
		return "Something went wrong, please try again later"
	}
}

//Create ...
func (f FlightForm) Create(err error) string {
	switch err.(type) {
	case validator.ValidationErrors:

		if _, ok := err.(*json.UnmarshalTypeError); ok {
			return "Something went wrong, please try again later"
		}

		for _, err := range err.(validator.ValidationErrors) {
			if err.Field() == "ToCode" {
				fmt.Println(err.Namespace())
				return f.Code(err.Tag(), err.Field())
			}
			if err.Field() == "FromCode" {
				return f.Code(err.Tag(), err.Field())
			}
		}

	default:
		return "Invalid request"
	}

	return "Something went wrong, please try again later"
}
