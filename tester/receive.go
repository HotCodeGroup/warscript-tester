package tester

import (
	"encoding/json"

	"github.com/HotCodeGroup/warscript-tester/games"
	"github.com/HotCodeGroup/warscript-tester/pong"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

type TesterStatusQueue struct {
	Type string          `json:"type"`
	Body json.RawMessage `json:"body"`
}

type TesterStatusUpdate struct {
	NewStatus string `json:"new_status"`
}

type TesterStatusError struct {
	Error string `json:"error"`
}

type TesterStatusResult struct {
	Winner int           `json:"result"`
	States []games.State `json:"states"`
}

type Lang string
type TestTask struct {
	Code1    string `json:"code1"`
	Code2    string `json:"code2"`
	GameSlug string `json:"game_slug"`
	Language Lang   `json:"lang"`
}

const (
	pongSlug = "pong"
)

var (
	receivedMessage = &TesterStatusUpdate{
		NewStatus: "Received. Starting containers",
	}
)

func sendReplyTo(ch *amqp.Channel, to, correlationId, t string, message interface{}) error {
	body, err := json.Marshal(message)
	if err != nil {
		return errors.Wrap(err, "can not marshal receivedMessage")
	}

	newMessage, err := json.Marshal(&TesterStatusQueue{
		Type: t,
		Body: body,
	})
	if err != nil {
		return errors.Wrap(err, "can not marshal TesterStatusQueue")
	}

	return ch.Publish(
		"",    // exchange
		to,    // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType:   "application/json",
			CorrelationId: correlationId,
			Body:          newMessage,
		},
	)
}

func ReceiveVerifyRPC(ch *amqp.Channel, d amqp.Delivery) error {
	err := sendReplyTo(ch, d.ReplyTo, d.CorrelationId, "status", receivedMessage)
	if err != nil {
		return errors.Wrap(err, "can not send receive confirmation")
	}

	task := &TestTask{}
	err = json.Unmarshal(d.Body, task)
	if err != nil {
		return errors.Wrap(err, "can not unmarshal delivery body")
	}

	var game games.Game
	switch task.GameSlug {
	case pongSlug:
		game = &pong.Pong{}
	}

	states, result, err := Test(task.Code1, task.Code2, game)
	if err != nil {
		err := sendReplyTo(ch, d.ReplyTo, d.CorrelationId, "error", &TesterStatusError{err.Error()})
		if err != nil {
			return errors.Wrap(err, "can not send internal error")
		}
	} else {
		err = sendReplyTo(ch, d.ReplyTo, d.CorrelationId, "result",
			&TesterStatusResult{
				Winner: result.GetWinner(),
				States: states,
			})
		if err != nil {
			return errors.Wrap(err, "can not send result state")
		}
	}

	err = d.Ack(false)
	if err != nil {
		return errors.Wrap(err, "can not ack delivery")
	}

	return nil
}
