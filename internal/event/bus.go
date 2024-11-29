package event

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"sync"

	"github.com/nats-io/nats.go"
)

const NATS_RECONNECT_BUF_SIZE = 5*1024*1024

type chanHolder struct {
	Ch interface{}
}

var (
	conn  *nats.Conn
	subCh []chanHolder
	mu    sync.Mutex
)

func Connect(server string) error {
	mu.Lock()
	defer mu.Unlock()

	if conn != nil {
		return errors.New("already connected")
	}

	c, err := nats.Connect(server, nats.ReconnectBufSize(NATS_RECONNECT_BUF_SIZE))
	if err != nil {
		return err
	}
	conn = c
	subCh = make([]chanHolder, 0)

	return nil
}

func Subscribe[T any](topic string) (chan T, error) {
	ch := make(chan T)
	_, err := conn.Subscribe(topic, func(m *nats.Msg) {
		var msg T
		decoder := json.NewDecoder(bytes.NewReader(m.Data))
		decoder.DisallowUnknownFields()
		err := decoder.Decode(&msg)
		if err != nil {
			log.Printf("Error unmarshalling message: %v", err)
			return
		}
		ch <- msg
	})
	if err != nil {
		return nil, err
	}
	mu.Lock()
	subCh = append(subCh, chanHolder{Ch: ch})
	mu.Unlock()

	return ch, nil
}

func Publish[T any](msg T, topic string) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	log.Printf("Publishing message to %s: %s", topic, string(data))
	return conn.Publish(topic, data)
}

func Close() {
	conn.Close()
	for _, holder := range subCh {
		if ch, ok := holder.Ch.(chan interface{}); ok {
			close(ch)
		}
	}
}
