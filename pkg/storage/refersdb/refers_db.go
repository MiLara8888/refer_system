package storage

import (
	"context"
)

// источник данных
type RefersDB interface {
	// закрытие базы
	Close(ctx context.Context) error

	//проверка есть ли пользователь в базе
	UserIs(ctx context.Context, email string) (bool, error)

	//сохранение пользователя
	UserSave(ctx context.Context, data *UserSerializer) (int, error)

	//получение пароля
	GetPasswordByLogin(ctx context.Context, email string) (*UserPasswordSerializer, error)

	//проверка токена
	IsAuch(ctx context.Context, token string) (*UserSerializer, error)

	//изменить/ создать код
	RefCodeUpdate(ctx context.Context, user *UserSerializer, exp_day int) (*RefCodeSerializer, error)

	//удалить дейструющий код
	RefCodeDelete(ctx context.Context, user *UserSerializer) error

	//получить код по email
	GetCode(ctx context.Context, email string) (*RefCodeSerializer, error)

	//создание подписчика на код
	CreateReferals(ctx context.Context, code, email string) (*RefSubsirbeSerializer, error)

	//список по id рефера
	ReferalsList(ctx context.Context, id int) ([]*RefSubsirbeSerializer, error)
}
