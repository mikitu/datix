package queue

import (
	"github.com/mikitu/datix/util"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"os"
)

type Queue struct {
	RpcQueue  	  amqp.Queue
	Conn          *amqp.Connection
	Channel       *amqp.Channel
	RpcQueueName string
}

func New() *Queue {
	return &Queue{RpcQueueName: os.Getenv("QUEUE_NAME")}
}

func (q *Queue) Open() {
	conn, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
	util.FailOnError(err, "Failed to connect to RabbitMQ")
	q.Conn = conn
	log.Debug("Connection opened")
}

func (q *Queue) OpenChannel() {
	ch, err := q.Conn.Channel()
	util.FailOnError(err, "Failed to open a channel")
	q.Channel = ch
	log.Debug("Channel opened")
}

func (q *Queue) CloseChannel() {
	err := q.Channel.Close()
	util.FailOnError(err, "Failed to close channel")
	log.Debug("Channel closed")

}

func (q *Queue) CloseConnection() {
	err := q.Conn.Close()
	util.FailOnError(err, "Failed to close connection")
	log.Debug("Connection closed")
}

func (q Queue) waitForMessages(autAck bool, deliveryCh chan amqp.Delivery) {
	msgs, err := q.Channel.Consume(
		q.RpcQueue.Name, // queue
		"",      // consumer
		autAck,    // auto-ack
		false,   // exclusive
		true,     // no-local
		false,     // no-wait
		nil,         // args
	)
	util.FailOnError(err, "Failed to register a consumer")
	log.Infof("received: %+v messages", len(msgs))

	go func() {
		for d := range msgs {
			deliveryCh <- d
		}
	}()
}

func (q Queue) WaitForResponse(deliveryCh chan amqp.Delivery) {
	q.waitForMessages(true, deliveryCh)
}

func (q Queue) WaitForRequest(deliveryCh chan amqp.Delivery) {
	q.waitForMessages(false, deliveryCh)
}

func (q Queue) SendRequest(corrId string, reqBodyBytes []byte) {
	err := q.Channel.Publish(
		"", 			// exchange
		q.RpcQueueName,     	// routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: 	"text/plain",
			CorrelationId: 	corrId,
			ReplyTo: 		q.RpcQueue.Name,
			Body:    		reqBodyBytes,
		})
	util.FailOnError(err, "Failed to close connection")
}

func (q Queue) SendResponse(replyTo, reqId string, reqBodyBytes []byte) {
	err := q.Channel.Publish(
		"", 			// exchange
		replyTo,     			// routing key
		false,       	// mandatory
		false,		// immediate
		amqp.Publishing{
			ContentType: "text/plain",
			CorrelationId: reqId,
			Body:        reqBodyBytes,
		})
	util.FailOnError(err, "Failed to close connection")
}

func (q *Queue) SetUpRequest() {
	rpcQueue, err := q.Channel.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // noWait
		nil,   // arguments
	)
	util.FailOnError(err, "Failed to declare a queue")
	q.RpcQueue = rpcQueue
}

func (q *Queue) SetUpResponse() {
	rpcQueue, err := q.Channel.QueueDeclare(
		q.RpcQueueName, // name
		false,       // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	util.FailOnError(err, "Failed to declare a queue")
	q.RpcQueue = rpcQueue

	err = q.Channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	util.FailOnError(err, "Failed to set QoS")
}

