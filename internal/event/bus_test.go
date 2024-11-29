package event

import (
	"fmt"
	"net"
	"testing"

	natsserver "github.com/nats-io/nats-server/v2/test"
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

	err = Connect(fmt.Sprintf("nats://localhost:%d", port))
	assert.NoError(t, err)

	topic := "test"
	msg := Message{Text: "Hello World"}

	ch, err := Subscribe[Message](topic)
	assert.NoError(t, err)
	assert.NotNil(t, ch)

	err = Publish(msg, topic)
	assert.NoError(t, err)

	received := <-ch
	assert.Equal(t, msg.Text, received.Text)

	Close()
}
