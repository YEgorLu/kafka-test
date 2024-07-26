package util

import (
	"encoding/json"

	"github.com/segmentio/kafka-go"
)

func FormatMessagesJSON(topic string, objects ...json.Marshaler) ([]kafka.Message, error) {
	messages := make([]kafka.Message, 0, len(objects))
	for _, v := range objects {
		data, err := v.MarshalJSON()
		if err != nil {
			return []kafka.Message{}, err
		}
		msg := kafka.Message{
			Topic: topic,
			Value: data,
		}
		messages = append(messages, msg)
	}
	return messages, nil
}
