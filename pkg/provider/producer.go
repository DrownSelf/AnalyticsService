package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	Writer *kafka.Conn
}

func NewProducer(topic string) (*Producer, error) {
	conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9093", topic, 0)
	if err != nil {
		return nil, err
	}

	err = conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		return nil, err
	}

	return &Producer{
		Writer: conn,
	}, nil
}

func (p *Producer) Produce(key []byte, value []byte) {
	_, err := p.Writer.WriteMessages(kafka.Message{
		Offset: 0,
		Key:    key,
		Value:  value,
	})

	if err != nil {
		fmt.Printf("delivery failed %s \n", err.Error())
	} else {
		fmt.Printf("message delivered | key: %s\n", string(key))
	}
}
