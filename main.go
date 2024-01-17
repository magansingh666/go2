package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/magansingh666/go2/db"
	"github.com/magansingh666/go2/models"
)

var clientList sync.Map
var inQ chan models.Massage
var outQ chan models.Massage

func main() {
	db.InitDB()
	app := fiber.New()
	go sendTimeToEachClient()
	inQ = make(chan models.Massage, 10000)
	outQ = make(chan models.Massage, 10000)

	go inQHandler()
	go outQHandler()

	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws/:clientId/:userId", websocket.New(func(c *websocket.Conn) {
		clientId := c.Params("clientId")
		userId := c.Params("userId")
		fmt.Println("This is client Id :-", clientId, userId)
		clientList.Store(clientId, c)
		log.Println("added client id in client list : ", clientId)
		var msg []byte
		var mt int
		var err error

		for {
			if mt, msg, err = c.ReadMessage(); err != nil {
				log.Println("Read Error := ", err)
				clientList.Delete(clientId)
				log.Println("\n delete cliend id from client list map : ", clientId, mt, msg)

				break
			}
			m := models.Massage{ClientId: clientId, UserId: userId, Msg: msg}
			inQ <- m
			fmt.Println("this is in q ....", inQ, outQ)

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

func inQHandler() {
	fmt.Println("\n Staring in Q handler")
	for {
		m := <-inQ
		fmt.Println("message recived processin it ", m)
		e := db.CreateProuduct(m.Msg)
		fmt.Print(e)
		m.Msg = []byte(fmt.Sprint(map[string]string{"message": "success"}))
		outQ <- m

	}

}

func outQHandler() {
	fmt.Println("\n Staring out Q handler")
	for {
		m := <-outQ
		fmt.Println("sending processed message ", m)
		cIf, _ := clientList.Load(m.ClientId)
		conn := cIf.(*websocket.Conn)
		conn.WriteMessage(websocket.TextMessage, m.Msg)

	}

}

/*
go func() {
				fmt.Println("starting go routine to process request")
				e := db.CreateProuduct(msg)
				if e != nil {
					c.WriteJSON(e)
					return
				}
				c.WriteJSON(map[string]string{"result": "success"})

			}()



*/
