package main

import (
	"bytes"
	"encoding/json"
	"github.com/joho/godotenv"
	"github.com/mikitu/datix/db"
	"github.com/mikitu/datix/model"
	"github.com/mikitu/datix/queue"
	"github.com/mikitu/datix/util"
	"github.com/streadway/amqp"
	"net/http"
	"strconv"
)

var Db *db.Db
func init() {
	Db = db.New()
	Db.Import()
}

func main() {
	godotenv.Load("../../.env.dist")
	godotenv.Overload("../../.env")
	q := queue.New()
	q.Open()
	defer q.CloseConnection()
	q.OpenChannel()
	defer q.CloseChannel()
	q.SetUpResponse()

	deliveryCh := make(chan amqp.Delivery, 10)
	responseCh := make(chan model.Response, 10)

	q.WaitForRequest(deliveryCh)

	for {
		select {
		case msg := <-deliveryCh:
			req := model.Request{}
			err := json.Unmarshal(msg.Body, &req)
			util.FailOnError(err, "Cannot decode request")
			resp := prepareResponse(req, msg.ReplyTo)
			responseCh <- resp
			msg.Ack(false)
		case rmsg := <-responseCh:
			reqBodyBytes := new(bytes.Buffer)
			json.NewEncoder(reqBodyBytes).Encode(rmsg)
			q.SendResponse(rmsg.ReplyTo, rmsg.Request.RequestId, reqBodyBytes.Bytes())
		}
	}
}

func prepareResponse(req model.Request, replyTo string) model.Response {
	var res interface{}
	var err error
	id, err := strconv.ParseUint(req.Id, 10, 64)
	if id == 0 {
		res, err = Db.FindAll()
	} else {
		res, err = Db.FindById(id)
	}
	status := http.StatusOK
	errStr := ""

	if err != nil {
		status = http.StatusNotFound
		errStr = err.Error()
		res = nil
	}

	return model.Response{
		ReplyTo: replyTo,
		Request: req,
		HttpResponse: model.HttpResponse {
			Data:    res,
			Err:     errStr,
			Status: status,
		},
	}
}
