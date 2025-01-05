package adapter

import (
	"context"
	_ "encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"songs/internal/app/common"
)

type ginRequestReader struct {
	c *gin.Context
}

func (g *ginRequestReader) PathParam(name string) (string, error) {
	param := g.c.Param(name)
	if param == "" {
		return "", fmt.Errorf("parameter %s not found", name)
	}
	return param, nil
}

func (g *ginRequestReader) QueryParam(name string) string {
	return g.c.Query(name)
}

func (g *ginRequestReader) DefaultQueryParam(name, defaultValue string) string {
	return g.c.DefaultQuery(name, defaultValue)
}

func (g *ginRequestReader) DecodeBody(v interface{}) error {
	return g.c.ShouldBindJSON(v)
}

func (g *ginRequestReader) Context() context.Context {
	return g.c.Request.Context()
}

// ToGinHandler converts handler to gin.HandlerFunc
func ToGinHandler(handler func(common.RequestReader, http.ResponseWriter) error) gin.HandlerFunc {
	return func(c *gin.Context) {
		reader := &ginRequestReader{c: c}
		if err := handler(reader, c.Writer); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"ToGinHandler converts handler error": err.Error()})
		}
	}
}
