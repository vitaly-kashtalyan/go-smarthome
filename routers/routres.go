package routers

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/location"
	"github.com/gin-gonic/gin"
	"time"
)

const (
	pathToApi = "/api/"
	v1        = "v1"
)

func Setup() *gin.Engine {
	fmt.Println("Starting routes with gin")

	r := gin.Default()
	r.Use(
		cors.Default(),
		location.Default(),
		gzip.Gzip(gzip.BestCompression),
	)
	InitializeRoutes(r)
	return r
}

func InitializeRoutes(r *gin.Engine) {
	v1 := r.Group(pathToApi + v1)
	{
		cars := v1.Group("/")
		{
			cars.GET("/ping", func(c *gin.Context) {
				c.String(200, "pong "+fmt.Sprint(time.Now().Unix()))
			})
		}
	}
}
