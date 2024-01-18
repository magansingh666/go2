package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/magansingh666/go2/db"
	"github.com/magansingh666/go2/models"
)

var clientList sync.Map
var inQ chan models.IMsg
var outQ chan models.OMsg

func main() {
	db.InitDB()
	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Static("/", "./views")
	app.Get("/index", func(c *fiber.Ctx) error {

		return c.Render("index", fiber.Map{})
	})
	go sendTimeToEachClient()
	inQ = make(chan models.IMsg, 10000)
	outQ = make(chan models.OMsg, 10000)

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
		fmt.Println("\n Client Id and User Id: ", clientId, userId)
		clientList.Store(clientId, c)
		log.Println("added client id in client list : ", clientId)

		for {

			mt, mb, err := c.ReadMessage()
			if err != nil {
				fmt.Println("Calling close handler; message type is  ", mt)
				clientList.Delete(clientId)
				break

			}
			im := models.IMsg{}
			json.Unmarshal(mb, &im)
			im.ClientId = clientId
			im.UserId = userId
			inQ <- im

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

				wc, ok := value.(*websocket.Conn)
				if !ok {
					i++
					return true

				}
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
	fmt.Println("\n Staring InQueue handler")
	for {
		m := <-inQ
		fmt.Println("Processing message from InQueue:  ", m)
		if len(m.C) < 2 {
			oM := models.OMsg{}
			oM.ClientId = m.ClientId
			oM.UserId = m.UserId
			oM.E = "true"
			oM.D = "json not formatted properly or wrong c parameter"
			outQ <- oM
			return
		}
		switch m.C[0] {
		case "product":
			oM := models.OMsg{}
			p, e := db.CreateProuduct(m.D)
			if e != nil {
				oM.E = "true"
				oM.D = e.Error()
			}
			oM.ClientId = m.ClientId
			oM.UserId = m.UserId
			oM.D = fmt.Sprint(p)
			outQ <- oM
		default:
			oM := models.OMsg{}
			oM.ClientId = m.ClientId
			oM.UserId = m.UserId
			oM.E = "true"
			oM.D = "json not formatted properly or wrong c parameter"
			outQ <- oM

		}

	}

}

func outQHandler() {
	fmt.Println("\n Staring out Q handler")
	for {
		m := <-outQ
		fmt.Println("sending processed message ", m)
		cIf, _ := clientList.Load(m.ClientId)

		conn, ok := cIf.(*websocket.Conn)
		if !ok {
			return
		}

		conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprint(m)))

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
