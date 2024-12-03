package event

import (
	"encoding/json"
	"fmt"
	"net"
	"testing"

	"github.com/hashicorp/go-hclog"
	natsserver "github.com/nats-io/nats-server/v2/test"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
)

type Message struct {
	Text string `json:"text"`
}

func TestBus(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	// Get a random free port
	conn, err := net.Listen("tcp", ":0")
	port := conn.Addr().(*net.TCPAddr).Port
	if err == nil {
		conn.Close()
	}

	options := natsserver.DefaultTestOptions
	options.Port = port

	s := natsserver.RunServer(&options)
	defer s.Shutdown()

	nb := NewNatsBus(hclog.Default())

	err = nb.Connect(fmt.Sprintf("nats://localhost:%d", port))
	assert.NoError(t, err)

	topic := "test"
	msg := Message{Text: "Hello World"}

	ch := make(chan Message)

	_, err = nb.conn.Subscribe(topic, func(m *nats.Msg) {
		var msg Message
		json.Unmarshal(m.Data, &msg)
		ch <- msg
	})

	err = Publish(nb, msg, topic)
	assert.NoError(t, err)

	received := <-ch
	assert.Equal(t, msg.Text, received.Text)

	nb.Close()
}
