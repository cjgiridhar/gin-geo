package main

import (
	"net/http"

	geo "github.com/cjgiridhar/gin-geo"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(geo.Default("github.com/cjgiridhar/gin-geo/db/GeoLite2-City.mmdb"))
	r.GET("/geo", func(c *gin.Context) {
		geoResponse, ok := c.Get("GeoResponse")
		if ok {
			c.JSON(http.StatusOK, gin.H{
				"geo": geoResponse,
			})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"detail": "Could not get geographical information",
			})
		}
	})
	r.Run()
}
