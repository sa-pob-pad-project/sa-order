package apperr

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

type Code int

const (
	CodeBadRequest Code = iota + 1
	CodeUnauthorized
	CodeForbidden
	CodeNotFound
	CodeConflict
	CodeInternal
)

type Error struct {
	Code   Code
	Msg    string
	Err    error
	Fields map[string]any
}

func (e *Error) Error() string {
	if e.Err != nil {
		return e.Msg + ": " + e.Err.Error()
	}
	return e.Msg
}
func (e *Error) Unwrap() error { return e.Err }

func New(code Code, msg string, err error) *Error {
	return &Error{Code: code, Msg: msg, Err: err}
}

func IsCode(err error, code Code) bool {
	var ae *Error
	if errors.As(err, &ae) {
		return ae.Code == code
	}
	return false
}

func WriteError(c *fiber.Ctx, err error) error {
	var ae *Error
	status := fiber.StatusInternalServerError
	msg := "internal error"

	if errors.As(err, &ae) {
		msg = ae.Msg
		switch ae.Code {
		case CodeBadRequest:
			status = fiber.StatusBadRequest
		case CodeUnauthorized:
			status = fiber.StatusUnauthorized
		case CodeForbidden:
			status = fiber.StatusForbidden
		case CodeNotFound:
			status = fiber.StatusNotFound
		case CodeConflict:
			status = fiber.StatusConflict
		default:
			status = fiber.StatusInternalServerError
		}
	}
	return c.Status(status).JSON(fiber.Map{"error": msg})
}
