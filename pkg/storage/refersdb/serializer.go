package storage

import (
	eml "net/mail"
	"regexp"
	"strings"
	"time"
)

type RefSubsirbeSerializer struct {
	ID     int64             `json:"id"`
	Email  string            `json:"email"`
	RefKey string             `json:"code"`
	UserId RefCodeSerializer `json:"ref_key"`
}

type RefCodeSerializer struct {
	UserID  int64     `json:"user_id"`
	Code    string    `json:"code"`
	ExpDate time.Time `json:"exp_date"`
}

type UserToken struct {
	Token  string `json:"token"`
	UserID int64  `json:"user_id"`
}

type UserSerializer struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Проверка авторизации
type UserPasswordSerializer struct {
	ID       int    `json:"id"`
	Password string `json:"password"`
	Salt     string `json:"salt"`
}

// валидация пароля и почты пользователя
// пароль должен быть из букв и цифр, email проверим стандартно
func (u *UserSerializer) Valid() bool {
	//если какое-то из полей будет пустое отдам ошибку
	if u.Email == "" || u.Password == "" {
		return false
	}
	if _, err := eml.ParseAddress(u.Email); err != nil {
		return false
	}
	//будет соответствовать строкам, которые состоят только из арабских цифр и букв латинского алфавита верхнего и нижнего регистров
	matched, err := regexp.MatchString(`^[0-9A-Za-z]+$`, u.Password)
	if (!matched) || (err != nil) {
		return false
	}

	u.Email = strings.ToLower(u.Email)

	return true
}
