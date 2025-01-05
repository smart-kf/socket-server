package main

import "github.com/gin-gonic/gin"

func main() {
	g := gin.Default()
	g.Static(".", ".")
	g.Run(":9091")
}
