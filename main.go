package main

import (
	"github.com/gin-gonic/gin"
	"poc-pushover/router"
)

func main() {
	app := gin.Default()
	app.LoadHTMLGlob("templates/*")
	router.Routes(app)
	app.Run()
}
