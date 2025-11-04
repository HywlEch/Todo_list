package apperrors

import (
	"fmt"
)

//定义一个包含HTTP状态码和错误信息的结构体
type AppError struct {
	Code int 
	Message string
	Err error
}

func (e *AppError)Error() string {
	if e.Err != nil {
		return fmt.Sprintf("Code: %d, Message: %s, InternalError: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("Code: %d, Message: %s", e.Code, e.Message)
}

//Unwarp用于errors.Is和errors.As
func (e *AppError)Unwarp() error {
	return e.Err
}

//构造函数
func NewAppError(code int, message string, err error) *AppError{
	return &AppError{
		Code: code,
		Message: message,
		Err: err,
	}
}

func NewNotFoundError(message string, err error) *AppError{
	if message == "" {
		message = "Not Found"
	}
	return NewAppError(404, message, err)
}

func NewInternalServerError(message string, err error) *AppError{
	if message == "" {
		message = "Internal Server Error"
	}
	return NewAppError(500, message, err)
}

func NewBadRequestError(message string, err error) *AppError{
	if message == "" {
		message = "Bad Request"
	}
	return NewAppError(400, message, err)
}

func NewUnauthorizedError(message string, err error) *AppError{
	if message == "" {
		message = "Unauthorized"
	}
	return NewAppError(401, message, err)
}
func NewConfilictError(message string, err error) *AppError{
	if message == "" {
		message = "Confilict"
	}
	return NewAppError(409, message, err)
}