package tester

import (
	"encoding/json"

	"github.com/HotCodeGroup/warscript-tester/atod"

	"github.com/HotCodeGroup/warscript-tester/games"
	"github.com/HotCodeGroup/warscript-tester/pong"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

// StatusQueue сообщение полученное для очереди задач
type StatusQueue struct {
	Type string          `json:"type"`
	Body json.RawMessage `json:"body"`
}

// StatusUpdate обновление статуса
type StatusUpdate struct {
	NewStatus string `json:"new_status"`
}

// StatusError сообщение об ошибке при проверке
type StatusError struct {
	Error string `json:"error"`
}

// StatusResult результат проверки
type StatusResult struct {
	Info   games.Info    `json:"info"`
	States []games.State `json:"states"`
	Winner int           `json:"result"`
}

// Lang по сути ENUM с доступными языками
type Lang string

// TestTask представление задачи на проверку
type TestTask struct {
	Code1    string `json:"code1"`
	Code2    string `json:"code2"`
	GameSlug string `json:"game_slug"`
	Language Lang   `json:"lang"`
}

const (
	pongSlug = "pong"
	atodSlug = "atod"
)

var (
	receivedMessage = &StatusUpdate{
		NewStatus: "Received. Starting containers",
	}
)

func sendReplyTo(ch *amqp.Channel, to, correlationID, t string, message interface{}) error {
	body, err := json.Marshal(message)
	if err != nil {
		return errors.Wrap(err, "can not marshal receivedMessage")
	}

	newMessage, err := json.Marshal(&StatusQueue{
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
			CorrelationId: correlationID,
			Body:          newMessage,
		},
	)
}

// ReceiveVerifyRPC обработка запросов полученных из очереди
func (t *Tester) ReceiveVerifyRPC(d amqp.Delivery) error {
	err := sendReplyTo(t.ch, d.ReplyTo, d.CorrelationId, "status", receivedMessage)
	if err != nil {
		return errors.Wrap(err, "can not send receive confirmation")
	}

	task := &TestTask{}
	err = json.Unmarshal(d.Body, task)
	if err != nil {
		return errors.Wrap(err, "can not unmarshal delivery body")
	}

	var game games.Game
	if task.GameSlug == pongSlug {
		game = &pong.Pong{}
	} else if task.GameSlug == atodSlug {
		game = &atod.Atod{}
	}

	info, states, result, err := t.Test(task.Code1, task.Code2, game)
	if err != nil {
		firstErr := err
		if errors.Cause(firstErr) == ErrTimeount {
			err = d.Reject(true)
			if err != nil {
				return errors.Wrap(err, "can not reject delivery")
			}

			return errors.Wrap(firstErr, "timeout error")
		}

		err = sendReplyTo(t.ch, d.ReplyTo, d.CorrelationId, "error", &StatusError{firstErr.Error()})
		if err != nil {
			return errors.Wrap(err, "can not send internal error")
		}

		err = d.Ack(true)
		if err != nil {
			return errors.Wrap(err, "can not ack delivery")
		}

		return errors.Wrap(firstErr, "tester error")
	}

	err = sendReplyTo(t.ch, d.ReplyTo, d.CorrelationId, "result",
		&StatusResult{
			Info:   info,
			States: states,
			Winner: result.GetWinner(),
		})
	if err != nil {
		return errors.Wrap(err, "can not send result state")
	}

	err = d.Ack(false)
	if err != nil {
		return errors.Wrap(err, "can not ack delivery")
	}
	return nil
}
