package refersrest

import (
	"encoding/json"
	"io"
	"net/http"
	er "refers_rest/pkg/errors"
	storage "refers_rest/pkg/storage/refersdb"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// регистрация пользователя
func (m *Rest) Register(c *gin.Context) {

	jsonData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	data := &storage.UserSerializer{}

	err = json.Unmarshal(jsonData, &data)
	if err != nil || (!data.Valid()) {
		c.AbortWithStatusJSON(http.StatusBadRequest, er.DataNotValid)
		return
	}

	//проверка, есть ли такой пользователь уже
	UserIsThere, err := m.DB.UserIs(c.Request.Context(), data.Email)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if UserIsThere {
		c.AbortWithStatusJSON(http.StatusBadRequest, er.UserIsThere)
		return
	}

	//если пользователя ещё нет и данные валидны, то его можно сохранить
	UserId, err := m.DB.UserSave(c.Request.Context(), data)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	res := struct {
		Id int `json:"user_id"`
	}{
		Id: UserId,
	}
	c.AbortWithStatusJSON(http.StatusCreated, res)

}

// аутентификация пользователя и получение токена
func (m *Rest) Auth(c *gin.Context) {
	jsonData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	data := &storage.UserSerializer{}

	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, er.DataNotValid)
		return
	}

	data.Email = strings.ToLower(data.Email)

	//проверка, есть ли такой пользователь уже
	UserIsThere, err := m.DB.UserIs(c.Request.Context(), data.Email)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if !UserIsThere {
		c.AbortWithStatusJSON(http.StatusBadRequest, er.DataNotValid)
		return
	}

	//получение данных для проверки авторизации
	password, err := m.DB.GetPasswordByLogin(c.Request.Context(), data.Email)
	if err != nil {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	// проверка на совпадение пароля
	if err := storage.PasswordMatched(data.Password, password.Password, password.Salt); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, er.ErrUserPassword)
		return
	}

	// генерация токена
	token, err := GenerateTokenJwt(data.Email, m.SecretKey, time.Duration(m.ExpTokenDay*86400*int(time.Second)))
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	res := struct {
		Token  string `json:"token"`
		UserID int    `json:"user_id"`
	}{
		Token:  token,
		UserID: password.ID,
	}

	c.AbortWithStatusJSON(http.StatusOK, res)

}
