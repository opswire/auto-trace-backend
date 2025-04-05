package v1

import "github.com/gin-gonic/gin"

const (
	DefaultInternalServerErrorMessage = "Something went wrong."
)

type response struct {
	Error string `json:"error" example:"message"`
}

func errorResponse(c *gin.Context, code int, msg string) {
	c.AbortWithStatusJSON(code, response{msg})
}
