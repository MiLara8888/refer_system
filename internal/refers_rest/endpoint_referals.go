package refersrest

import (
	"encoding/json"
	"io"
	"net/http"
	er "refers_rest/pkg/errors"
	"strings"

	"github.com/gin-gonic/gin"
)

// создать реферала на рефер код
func (m *Rest) CreateReferals(c *gin.Context) {

	data := struct {
		Code  string `json:"code"`
		Email string `json:"email"`
	}{}

	jsonData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(jsonData, &data)
	if err != nil || data.Code == "" || data.Email == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, er.DataNotValid)
		return
	}

	data.Email = strings.ToLower(data.Email)

	res, err := m.DB.CreateReferals(c.Request.Context(), data.Code, data.Email)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, res)

}

// список рефералсов по id рефера
func (m *Rest) ReferalsList(c *gin.Context) {
	data := struct {
		Id int `json:"user_id"`
	}{}

	jsonData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(jsonData, &data)
	if err != nil || data.Id == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, er.DataNotValid)
		return
	}
	res, err := m.DB.ReferalsList(c.Request.Context(), data.Id)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, res)
}
