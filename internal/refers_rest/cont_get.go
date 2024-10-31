package refersrest


import (
	"github.com/gin-gonic/gin"
)

func ContextGet[T comparable](n string, c *gin.Context, d T) T {

	val, ok := c.Get(n)
	if !ok {
		return d
	}
	t, ok := val.(T)

	if !ok {
		return d
	}
	return t
}
