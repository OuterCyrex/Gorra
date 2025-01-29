package GorraAPI

import (
	"github.com/OuterCyrex/Gorra/GorraAPI/middlewares"
	"github.com/gin-gonic/gin"
	"net/http"
)

type RouterGroup func(*gin.RouterGroup)

func KeepAliveRouters(version string, groups ...RouterGroup) *gin.Engine {
	R := gin.Default()
	R.Use(middlewares.Cors())

	R.GET("health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "alive",
		})
	})

	ApiGroup := R.Group(version)

	if len(groups) > 0 {
		for _, group := range groups {
			group(ApiGroup)
		}
	}

	return R
}
