package er

import (
	"fmt"
)

var (
	// rest api
	ErrorIsNotMatch = New(0, "")
	DataNotValid    = New(50, "Не валидные данные")
	UserIsThere     = New(60, "Такой пользователь уже зарегистрирован")
	ErrSaltDecode   = New(15, "ошибка преобразования соли")
	ErrGenPassword  = New(16, "ошибка в генерации пароля")
	ErrIsNotMatch   = New(17, "пароли не совпадают")
	ErrUserPassword = New(6, "пароль не верен")
)

type Error struct {
	ErrorCode        int    `json:"error_code"`
	ErrorDescription string `json:"error_description"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("err_type:%d , err_des:%s", e.Code, e.ErrorDescription)
}
func (e *Error) Code() int {
	return e.ErrorCode
}

func New(code int, desc string) *Error {
	return &Error{
		ErrorCode:        code,
		ErrorDescription: desc,
	}
}
