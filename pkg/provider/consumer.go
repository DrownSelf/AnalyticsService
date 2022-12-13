package provider

import (
	"context"
	"encoding/json"
	"time"

	"github.com/segmentio/kafka-go"
)

type Consumer[T comparable] struct {
	reader *kafka.Reader
}

func (c *Consumer[T]) CreateConnection(connections []string, dialer *kafka.Dialer, topic string) error {
	c.reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:   connections,
		Topic:     topic,
		Partition: 0,
		MinBytes:  10e3,
		MaxBytes:  10e6,
		MaxWait:   time.Millisecond * 10,
		Dialer:    dialer,
	})
	if err := c.reader.SetOffset(0); err != nil {
		return err
	}
	return nil
}

func (c *Consumer[T]) Read(model T, callback func(T, error)) {
	for {
		ctx, _ := context.WithTimeout(context.Background(), time.Millisecond*80)
		message, err := c.reader.ReadMessage(ctx)

		if err != nil {
			callback(model, err)
			return
		}

		err = json.Unmarshal(message.Value, &model)
		if err != nil {
			callback(model, err)
			return
		}
		callback(model, nil)
	}
}
