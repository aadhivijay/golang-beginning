package main

import (
	"fmt"
	"golang/students"
	"golang/websocket"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Hello World!")

	startServer()
}

func startServer() {
	ginServer := gin.Default()

	ginServer.GET("/ping", ping)

	students.Init(ginServer)
	websocket.Init(ginServer)

	if err := ginServer.Run(":3000"); err != nil {
		log.Fatal(err)
	}
}

func ping(con *gin.Context) {
	con.JSON(http.StatusOK, gin.H{
		"msg": "PONG",
	})
}
