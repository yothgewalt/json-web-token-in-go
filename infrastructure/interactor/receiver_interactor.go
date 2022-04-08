package interactor

import (
	"WebTokenAuthorization/infrastructure/presenter"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Receiver(c *gin.Context) {
	c.JSON(http.StatusOK, presenter.ResponseOnlyMessage{Message: "hello, you can reach the first endpoint of service for testing response"})
}
