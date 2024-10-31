package refersrest

import (
	"encoding/json"
	"io"
	"net/http"
	er "refers_rest/pkg/errors"
	storage "refers_rest/pkg/storage/refersdb"
	"strings"

	"github.com/gin-gonic/gin"
)

// создать реферальный код и задать ему срок годности
func (m *Rest) CreateRef(c *gin.Context) {

	body := struct {
		Day int `json:"day"`
	}{}

	jsonData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(jsonData, &body)
	if err != nil || body.Day == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, er.DataNotValid)
		return
	}

	user := ContextGet[*storage.UserSerializer]("user", c, &storage.UserSerializer{})

	data, err := m.DB.RefCodeUpdate(c.Request.Context(), user, body.Day)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, data)

}

// удалить действующий код
func (m *Rest) DeleteRef(c *gin.Context) {

	user := ContextGet[*storage.UserSerializer]("user", c, &storage.UserSerializer{})

	err := m.DB.RefCodeDelete(c.Request.Context(), user)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	c.AbortWithStatus(http.StatusOK)

}

// получить код по email пользователя
func (m *Rest) GetCode(c *gin.Context) {

	body := struct {
		Email string `json:"email"`
	}{}

	jsonData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(jsonData, &body)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, er.DataNotValid)
		return
	}

	email := strings.ToLower(body.Email)

	res, err := m.DB.GetCode(c.Request.Context(), email)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, res)

}
