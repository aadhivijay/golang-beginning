package main

import (
	"fmt"
	"golang/config"
	"golang/students"
	"golang/websocket"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

var env config.Config = *config.GetConfig()

func main() {
	fmt.Println("Hello World!")

	startServer()
}

func startServer() {
	ginServer := gin.Default()

	ginServer.GET("/ping", ping)

	students.Init(ginServer)
	websocket.Init(ginServer)

	if err := ginServer.Run(":" + env.Port); err != nil {
		log.Fatal(err)
	}
}

func ping(con *gin.Context) {
	con.JSON(http.StatusOK, gin.H{
		"msg": "PONG",
	})
}
