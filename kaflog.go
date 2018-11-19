package kaflog

import (
	"context"
	"fmt"
	"time"

	"github.com/eaciit/toolkit"

	kafka "github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/snappy"
)

func Hook(host string, topic string, brokers ...string) func(logType string, message string) {
	return func(logType string, message string) {
		m := toolkit.M{}.Set("Host", host).
			Set("TimeStamp", time.Now()).
			Set("LogType", logType).
			Set("Message", message)

		dialer := &kafka.Dialer{
			Timeout:  10 * time.Second,
			ClientID: "BEF",
		}

		config := kafka.WriterConfig{
			Brokers:          brokers,
			Topic:            topic,
			Balancer:         &kafka.LeastBytes{},
			Dialer:           dialer,
			WriteTimeout:     10 * time.Second,
			ReadTimeout:      10 * time.Second,
			CompressionCodec: snappy.NewCompressionCodec(),
		}

		w := kafka.NewWriter(config)
		err := w.WriteMessages(context.Background(), kafka.Message{Value: toolkit.Jsonify(m)})
		if err != nil {
			fmt.Printf("Kafka writing failed. %s.\n%s \n",
				err.Error(),
				m)
			return
		}
		w.Close()
	}
}
