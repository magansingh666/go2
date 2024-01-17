package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/magansingh666/go2/db"
)

var clientList sync.Map

func main() {
	db.InitDB()
	app := fiber.New()
	//sendTimeToEachClient()

	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws/:clientId", websocket.New(func(c *websocket.Conn) {
		clientId := c.Params("clientId")
		fmt.Println("This is client Id :-", clientId)
		clientList.Store(clientId, c)
		log.Println("added client id in client list : ", clientId)
		var msg []byte

		for {
			if mt, msg, err := c.ReadMessage(); err != nil {
				log.Println("Read Error := ", err)
				clientList.Delete(clientId)
				log.Println("\n delete cliend id from client list map : ", clientId, mt, msg)

				break
			}

			go func() {
				fmt.Println("starting go routine to process request")
				e := db.CreateProuduct(msg)
				if e != nil {
					c.WriteJSON(e)
				}
				c.WriteJSON(map[string]string{"result": "success"})

			}()

		}

	}))

	log.Fatal(app.Listen(":3000"))
	// Access the websocket server: ws://localhost:3000/ws/123?v=1.0
	// https://www.websocket.org/echo.html
}

func sendTimeToEachClient() {
	ticker := time.NewTicker(5 * time.Second)
	i := 0
	go func() {
		for v := range ticker.C {
			clientList.Range(func(key, value interface{}) bool {
				wc := value.(*websocket.Conn)
				e := wc.WriteMessage(websocket.TextMessage, []byte(fmt.Sprint(v)))
				if e != nil {
					fmt.Println(e)
				}
				i++
				return true
			})

		}
	}()

}
