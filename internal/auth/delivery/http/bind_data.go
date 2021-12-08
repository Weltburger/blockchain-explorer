package http

import (
	"fmt"
	"log"

	"explorer/internal/apperrors"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type bindError struct {
	Error     *apperrors.Error  `json:"error"`
	Arguments []invalidArgument `json:"invalidArguments"`
}

// used to help extract validation errors
type invalidArgument struct {
	Field string `json:"field"`
	Value string `json:"value"`
	Tag   string `json:"tag"`
	Param string `json:"param"`
}

// bindData is helper function, returns false if data is not bound
func bindData(c echo.Context, req interface{}) bool {
	// extract request content type
	h := c.Request().Header
	if h.Get("Content-Type") != "application/json" {
		msg := fmt.Sprintf("%s only accepts Content-Type application/json", c.Path())

		err := apperrors.NewUnsupportedMediaType(msg)

		c.JSON(err.Status(), err)
		return false
	}
	// Bind incoming json to struct and check for validation errors
	if err := c.Bind(req); err != nil {
		log.Printf("Error binding data: %+v\n", err)

		if errs, ok := err.(validator.ValidationErrors); ok {
			// could probably extract this, it is also in middleware_auth_user
			var invalidArgs []invalidArgument

			for _, err := range errs {
				invalidArgs = append(invalidArgs, invalidArgument{
					err.Field(),
					err.Value().(string),
					err.Tag(),
					err.Param(),
				})
			}

			err := apperrors.NewBadRequest("Invalid request parameters.")

			c.JSON(err.Status(), bindError{
				Error:     err,
				Arguments: invalidArgs,
			})
			return false
		}

		//return an internal server error
		fallBack := apperrors.NewInternal()

		c.JSON(fallBack.Status(), fallBack)
		return false
	}

	return true
}

// validateError
type validateError struct {
	Error     *apperrors.Error  `json:"error"`
	Arguments []invalidArgument `json:"invalidArguments"`
}

// validateData is helper function, returns false if data is not valid to requirements
func validData(ctx echo.Context, req interface{}) bool {
	// define receive data validator
	validate := validator.New()

	if err := validate.Struct(req); err != nil {

		// could probably extract this, it is also in middleware_auth_user
		var invalidArgs []invalidArgument

		for _, e := range err.(validator.ValidationErrors) {
			invalidArgs = append(invalidArgs, invalidArgument{
				e.Field(),
				e.Value().(string),
				e.Tag(),
				e.Param(),
			})
		}

		err := apperrors.NewBadRequest("Parameters validation error.")

		ctx.JSON(err.Status(), validateError{
			Error:     err,
			Arguments: invalidArgs,
		})
		return false
	}

	return true
}
