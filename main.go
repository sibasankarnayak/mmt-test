package main

import (
	"mmt/controllers"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	v1 := r.Group("/v1")
	{
		flight := new(controllers.FilghtController)
		v1.POST("/find/fastestFlight", flight.GetFastestFlight)
	}
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
