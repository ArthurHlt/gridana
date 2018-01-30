package amqp

import (
	"encoding/json"
	"github.com/ArthurHlt/gridana/model"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

const (
	QUEUE_NAME = "gridana_alerts"
)

type AmqpRemote struct {
	ch      *amqp.Channel
	msgChan <-chan amqp.Delivery
}

func NewAmqpRemote(conn *amqp.Connection) (*AmqpRemote, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	_, err = ch.QueueDeclare(
		QUEUE_NAME, // name
		false,      // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		return nil, err
	}
	msgs, err := ch.Consume(
		QUEUE_NAME, // queue
		"",         // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		return nil, err
	}
	return &AmqpRemote{
		ch:      ch,
		msgChan: msgs,
	}, nil
}
func (r AmqpRemote) Emit(alert model.FormattedAlert) {
	entry := log.WithField("alert_id", alert.ID).
		WithField("remote_emitter", "amqp").
		WithField("order", "emit")
	b, err := json.Marshal(alert)
	if err != nil {
		entry.Error(err)
	}
	err = r.ch.Publish(
		"",         // exchange
		QUEUE_NAME, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        b,
		})
	if err != nil {
		entry.Error(err)
	}
}
func (r AmqpRemote) Receive(emit func(model.FormattedAlert)) {
	entry := log.WithField("remote_emitter", "amqp").WithField("order", "receive")
	go func() {
		for d := range r.msgChan {
			var alert model.FormattedAlert
			err := json.Unmarshal(d.Body, &alert)
			if err != nil {
				entry.Error(err)
			}
			emit(alert)
		}
	}()
}
