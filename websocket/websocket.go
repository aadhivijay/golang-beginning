package websocket

import (
	"fmt"
	"golang/middleware"
	"math/rand"
	"net"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gobwas/ws"
)

type Client struct {
	Id         string      `json:"id"`
	Conn       net.Conn    `json:"conn"`
	MsgChannel chan string `json:"msgChannel"`
}

var clientsList = map[string]Client{}
var mutex = &sync.RWMutex{}

type Event struct {
	Event string   `json:"event"`
	Id    []string `json:"id"`
	Data  string   `json:"data"`
}

func Init(server *gin.Engine) {
	wsRouter := server.Group("/ws")
	wsRouter.Use(middleware.Authorize)
	{
		wsRouter.GET("", handleSocketConnection)
		wsRouter.POST("/send", postMessage)
	}
}

func RandStringBytes(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func handleSocketConnection(c *gin.Context) {
	conn, _, _, err := ws.UpgradeHTTP(c.Request, c.Writer)
	if err != nil {
		fmt.Println("Error in Upgrading: ", err)
		return
	}

	id, ok := c.GetQuery("id")
	if !ok {
		id = RandStringBytes(32)
	}

	newClient := Client{
		Id:         id,
		Conn:       conn,
		MsgChannel: make(chan string),
	}
	mutex.Lock()
	clientsList[id] = newClient
	mutex.Unlock()
	fmt.Printf("New client added! %+v\n", newClient)

	go initReadMessage(conn, newClient)
	go initSendMessage(conn, newClient)
}

func initReadMessage(con net.Conn, client Client) {
	defer func() {
		fmt.Printf("Client disconnected! %+v\n", client)
		close(client.MsgChannel)
		con.Close()
		mutex.Lock()
		delete(clientsList, client.Id)
		mutex.Unlock()
	}()
	for {
		header, errHeader := ws.ReadHeader(con)
		if errHeader != nil {
			fmt.Printf("[ERROR] readMessage: ReadHeader: %v\n", errHeader)
			return
		}

		if header.Length > 0 {
			payload := make([]byte, header.Length)
			_, errRead := con.Read(payload)
			if errRead != nil {
				fmt.Printf("[ERROR] readMessage: Read: %v\n", errRead)
				return
			}

			if header.Masked {
				ws.Cipher(payload, header.Mask, 0)
			}

			fmt.Printf("Msg: `%v` from client: %v\n", string(payload), client.Id)

			// sending back to client
			client.MsgChannel <- "From Server: " + string(payload)
		}

		if header.OpCode == ws.OpClose {
			return
		}
	}
}

func initSendMessage(con net.Conn, client Client) {
	for {
		msg, ok := <-client.MsgChannel
		if !ok {
			fmt.Printf("Channel already closed! %v\n", client)
			return
		}

		fmt.Printf("Sending: `%v` to client: %v\n", msg, client.Id)

		header := ws.Header{
			Fin:    true,
			Rsv:    0,
			OpCode: ws.OpText,
			Masked: false,
			Length: int64(len(msg)),
		}
		if err := ws.WriteHeader(con, header); err != nil {
			fmt.Printf("[ERROR] sendMessage: WriteHeader: %v\n", err)
			return
		}

		if _, err := con.Write([]byte(msg)); err != nil {
			fmt.Printf("[ERROR] sendMessage: Write: %v\n", err)
			return
		}
	}
}

func postMessage(c *gin.Context) {
	var event Event
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	fmt.Printf("POST MSG: `%v` to client: %v\n", event.Data, event.Id)

	if len(event.Id) > 0 {
		mutex.RLock()
		for _, v := range event.Id {
			if client, ok := clientsList[v]; ok {
				client.MsgChannel <- event.Data
				continue
			}
			fmt.Printf("Client: %v, already disconnected!\n", v)
		}
		mutex.RUnlock()
		return
	}

	// send to all
	mutex.RLock()
	for _, client := range clientsList {
		client.MsgChannel <- event.Data
	}
	mutex.RUnlock()

	c.JSON(http.StatusOK, event)
}
