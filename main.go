package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/HotCodeGroup/warscript-tester/tester"
	"github.com/jcftang/logentriesrus"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"

	log "github.com/sirupsen/logrus"
)

const (
	amqpPattern     = "amqp://%s:%s@%s:%s/"
	testerQueueName = "tester_rpc_queue"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)

	// собираем логи в хранилище
	le, err := logentriesrus.NewLogentriesrusHook(os.Getenv("LOGENTRIESRUS_TOKEN"))
	if err != nil {
		os.Exit(-1)
	}
	log.AddHook(le)
}

func failOnError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	var err error
	routineCount, err := strconv.ParseInt(os.Getenv("ROUTINE_COUNT"), 10, 16)
	failOnError(errors.Wrap(err, "ROUTINE_COUNT int parse error"))

	_, err = strconv.ParseInt(os.Getenv("QUEUE_PORT"), 10, 16)
	failOnError(errors.Wrap(err, "QUEUE_PORT int parse error"))

	conn, err := amqp.Dial(fmt.Sprintf(amqpPattern, os.Getenv("QUEUE_USER"),
		os.Getenv("QUEUE_PASS"), os.Getenv("QUEUE_HOST"), os.Getenv("QUEUE_PORT")))
	failOnError(errors.Wrap(err, "failed to connect to RabbitMQ"))
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(errors.Wrap(err, "failed to open a channel"))
	defer ch.Close()

	q, err := ch.QueueDeclare(
		testerQueueName, // name
		false,           // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	failOnError(errors.Wrap(err, "failed to declare a queue"))

	err = ch.Qos(
		int(routineCount), // prefetch count
		0,                 // prefetch size
		false,             // global
	)
	failOnError(errors.Wrap(err, "failed to set QoS"))

	msgs, err := ch.Consume(
		q.Name,          // queue
		os.Getenv("ID"), // consumer
		false,           // auto-ack
		false,           // exclusive
		false,           // no-local
		false,           // no-wait
		nil,             // args
	)
	failOnError(errors.Wrap(err, "failed to register a consumer"))

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			go func(ch *amqp.Channel, d amqp.Delivery) {
				err := tester.ReceiveVerifyRPC(ch, d)
				if err != nil {
					log.Errorf("[%s] %s", d.CorrelationId, err)
				}
			}(ch, d)
		}
	}()

	log.Printf("[%s] Awaiting RPC requests", os.Getenv("ID"))
	<-forever // блокируемся
}
