package kafkamock

import (
	"context"
	"testing"

	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
)

func TestMockKafka(t *testing.T) {
	mockKafka := NewMockKafka(10)

	// Test writing messages
	msgsToWrite := []kafka.Message{
		{Value: []byte("message1")},
		{Value: []byte("message2")},
	}

	err := mockKafka.WriteMessages(context.Background(), msgsToWrite...)
	assert.NoError(t, err, "expected no error when writing messages")

	// Test reading messages
	for _, expectedMsg := range msgsToWrite {
		msg, err := mockKafka.ReadMessage(context.Background())
		assert.NoError(t, err, "expected no error when reading messages")
		assert.Equal(t, expectedMsg.Value, msg.Value, "expected message values to be equal")
	}

	// Test reading from empty queue
	_, err = mockKafka.ReadMessage(context.Background())
	assert.Error(t, err, "expected error when reading from empty queue")

	// Test closing MockKafka
	err = mockKafka.Close()
	assert.NoError(t, err, "expected no error when closing MockKafka")

	// Test writing after closing
	err = mockKafka.WriteMessages(context.Background(), kafka.Message{Value: []byte("message3")})
	assert.Error(t, err, "expected error when writing after closing")

	// Test reading after closing
	_, err = mockKafka.ReadMessage(context.Background())
	assert.Error(t, err, "expected error when reading after closing")
}
