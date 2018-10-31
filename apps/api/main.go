package main

import (
	"bytes"
	"encoding/json"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	echo_middleware "github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"github.com/mikitu/datix/middleware"
	"github.com/mikitu/datix/model"
	"github.com/mikitu/datix/queue"
	"github.com/streadway/amqp"
	"net/http"
)

func main() {
	godotenv.Load("../../.env.dist")
	godotenv.Overload("../../.env")
	q := queue.New()
	q.Open()
	defer q.CloseConnection()

	e := echo.New()
	e.Use(middleware.RequestID)
	e.Use(echo_middleware.CORSWithConfig(echo_middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.GET("/users/:id", func(c echo.Context) error {
		return getUser(c, q)
	})
	e.GET("/users", func(c echo.Context) error {
		return getUsers(c, q)
	})
	e.Logger.Fatal(e.Start(":8080"))
}

func getUser(c echo.Context, q *queue.Queue) error {
	// User ID from path `users/:id`
	id := c.Param("id")
	log.Infof("Request user with Id: %v", id)
	// Request ID from context
	reqId := c.Get("RequestID").(string)
	log.Infof("Request Id: %v", reqId)
	return find(c, q, reqId, id)
}
func getUsers(c echo.Context, q *queue.Queue) error {
	// Request ID from context
	reqId := c.Get("RequestID").(string)
	return find(c, q, reqId, "0")
}

func find(c echo.Context, q *queue.Queue, reqId, id string) error {
	q.OpenChannel()
	defer q.CloseChannel()
	q.SetUpRequest()

	reqBody := model.Request{Id: id, RequestId: reqId}
	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(reqBody)

	q.SendRequest(reqId, reqBodyBytes.Bytes())

	deliveryCh := make(chan amqp.Delivery)
	q.WaitForResponse(deliveryCh)

	for {
		select {
		case rmsg := <-deliveryCh:
			res := &model.Response{}
			err := json.Unmarshal(rmsg.Body, res)
			if err != nil {
				c.Error(err)
			}
			body, err := json.Marshal(res.HttpResponse)
			if rmsg.CorrelationId == reqId {
				return c.String(res.HttpResponse.Status, string(body))
			}
		}
	}
	return nil
}
