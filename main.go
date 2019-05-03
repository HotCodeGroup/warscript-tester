package main

import (
	"log"
	"os"
	"strconv"

	"github.com/HotCodeGroup/warscript-tester/tester"
	"github.com/HotCodeGroup/warscript-utils/logging"
	"github.com/HotCodeGroup/warscript-utils/rabbitmq"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"

	"github.com/sirupsen/logrus"

	consulapi "github.com/hashicorp/consul/api"
	vaultapi "github.com/hashicorp/vault/api"
)

var logger *logrus.Logger

const (
	testerQueueName = "tester_rpc_queue"
)

func failOnError(err error) {
	if err != nil {
		logger.Fatal(err)
	}
}

func main() {
	var err error
	logger, err = logging.NewLogger(os.Stdout, os.Getenv("LOGENTRIESRUS_TOKEN"))
	if err != nil {
		log.Printf("can not create logger: %s", err)
		return
	}

	consulConfig := consulapi.DefaultConfig()
	consulConfig.Address = os.Getenv("CONSUL_ADDR")
	consul, err := consulapi.NewClient(consulConfig)
	if err != nil {
		logger.Errorf("can not connect consul service: %s", err)
		return
	}

	vaultConfig := vaultapi.DefaultConfig()
	vaultConfig.Address = os.Getenv("VAULT_ADDR")
	vault, err := vaultapi.NewClient(vaultConfig)
	if err != nil {
		logger.Errorf("can not connect vault service: %s", err)
		return
	}
	rabbitConf, err := vault.Logical().Read("warscript-bots/rabbitmq")
	if err != nil || rabbitConf == nil || len(rabbitConf.Warnings) != 0 {
		logger.Errorf("can read warscript-bots/rabbitmq key: %+v; %+v", err, rabbitConf)
		return
	}
	kv, _, err := consul.KV().Get("warscript-tester/routine-count", nil)
	if err != nil || kv == nil {
		logger.Errorf("can read warscript-tester/routine-count key: %+v;", err)
		return
	}
	routineCount, err := strconv.ParseInt(string(kv.Value), 10, 16)
	failOnError(errors.Wrap(err, "ROUTINE_COUNT int parse error"))

	endpoint := "unix:///var/run/docker.sock"
	client, err := docker.NewClient(endpoint)
	failOnError(errors.Wrap(err, "can not create docker client"))

	rabbitConn, err := rabbitmq.Connect(rabbitConf.Data["user"].(string), rabbitConf.Data["pass"].(string),
		rabbitConf.Data["host"].(string), rabbitConf.Data["port"].(string))
	if err != nil {
		logger.Errorf("can not connect to rabbitmq: %s", err.Error())
		return
	}
	defer rabbitConn.Close()

	rabbitChannel, err := rabbitConn.Channel()
	failOnError(errors.Wrap(err, "failed to open a channel"))
	defer rabbitChannel.Close()

	q, err := rabbitChannel.QueueDeclare(
		testerQueueName, // name
		false,           // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	failOnError(errors.Wrap(err, "failed to declare a queue"))

	err = rabbitChannel.Qos(
		int(routineCount), // prefetch count
		0,                 // prefetch size
		false,             // global
	)
	failOnError(errors.Wrap(err, "failed to set QoS"))

	msgs, err := rabbitChannel.Consume(
		q.Name,          // queue
		os.Getenv("ID"), // consumer
		false,           // auto-ack
		false,           // exclusive
		false,           // no-local
		false,           // no-wait
		nil,             // args
	)
	failOnError(errors.Wrap(err, "failed to register a consumer"))

	t := tester.NewTester(client, rabbitChannel)
	forever := make(chan bool)
	go func() {
		for d := range msgs {
			go func(ch *amqp.Channel, d amqp.Delivery) {
				logger.Infof("[NEW] Delivery %s", d.CorrelationId)
				err := t.ReceiveVerifyRPC(d)
				if err != nil {
					logger.Errorf("[ERROR] Delivery %s: %s", d.CorrelationId, err)
					return
				}
				log.Printf("[DONE] Delivery %s", d.CorrelationId)
			}(rabbitChannel, d)
		}
	}()

	log.Printf("[%s] Awaiting RPC requests", os.Getenv("ID"))
	<-forever // блокируемся
}
