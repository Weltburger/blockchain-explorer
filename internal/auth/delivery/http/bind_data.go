package http

import (
	"fmt"
	"log"

	"explorer/internal/apperrors"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type invalidArgument struct {
	Field string `json:"field"`
	Value string `json:"value"`
	Tag   string `json:"tag"`
	Param string `json:"param"`
}

// bindData is helper function, returns false if data is not bound
func bindData(c echo.Context, req interface{}) (bool, error) {
	// extract request content type
	h := c.Request().Header
	if h.Get("Content-Type") != "application/json" {
		msg := fmt.Sprintf("%s only accepts Content-Type application/json", c.Path())
		err := apperrors.NewUnsupportedMediaType(msg)

		return false, c.JSON(err.Status(), err)
	}
	// Bind incoming json to struct and check for validation errors
	if err := c.Bind(req); err != nil {
		log.Printf("Error binding data: %+v\n", err)

		//return an internal server error
		bindError := apperrors.NewInternal()
		return false, c.JSON(bindError.Status(), bindError)
	}

	return true, nil
}

// validateData is helper function, returns false if data is not valid to requirements
func validData(ctx echo.Context, req interface{}) (bool, error) {
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

		return false, ctx.JSON(err.Status(), map[string]interface{}{
			"error":            err,
			"invalidArguments": invalidArgs})
	}

	return true, nil
}
