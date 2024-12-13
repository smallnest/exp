package kafkamock

import (
	"context"
	"errors"
	"sync"

	"github.com/segmentio/kafka-go"
)

// MockReader simulates the kafka.Reader interface
type MockReader struct {
	messages chan *kafka.Message
	mu       sync.Mutex
	closed   bool
}

// NewMockReader creates a new MockReader
func NewMockReader(messages chan *kafka.Message) *MockReader {
	return &MockReader{
		messages: messages,
	}
}

// ReadMessage simulates reading a message
func (m *MockReader) ReadMessage(ctx context.Context) (kafka.Message, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.closed {
		return kafka.Message{}, errors.New("reader is closed")
	}

	if len(m.messages) == 0 {
		return kafka.Message{}, errors.New("no more messages")
	}

	msg := <-m.messages

	return *msg, nil
}

// Close closes the MockReader
func (m *MockReader) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.closed = true
	return nil
}

// MockWriter simulates the kafka.Writer interface
type MockWriter struct {
	messages chan *kafka.Message
	mu       sync.Mutex
	closed   bool
}

// NewMockWriter creates a new MockWriter
func NewMockWriter(messages chan *kafka.Message) *MockWriter {
	return &MockWriter{
		messages: messages,
	}
}

// WriteMessages simulates writing messages
func (m *MockWriter) WriteMessages(ctx context.Context, msgs ...kafka.Message) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.closed {
		return errors.New("writer is closed")
	}

	for i := range msgs {
		m.messages <- &msgs[i]
	}
	return nil
}

// Close closes the MockWriter
func (m *MockWriter) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.closed = true
	return nil
}

// MockKafka contains MockReader and MockWriter
type MockKafka struct {
	*MockReader
	*MockWriter
	messages chan *kafka.Message
}

// NewMockKafka creates a new MockKafka
func NewMockKafka(size int) *MockKafka {
	messages := make(chan *kafka.Message, size)
	return &MockKafka{
		MockReader: NewMockReader(messages),
		MockWriter: NewMockWriter(messages),
		messages:   messages,
	}
}

// Close closes the MockKafka
func (m *MockKafka) Close() error {
	m.MockReader.Close()
	m.MockWriter.Close()
	close(m.messages)
	return nil
}

// GetMessages gets the message channel of MockKafka
func (m *MockKafka) GetMessages() chan *kafka.Message {
	return m.messages
}

// MockKafkaInterface defines an interface to easily replace the real Kafka client
type MockKafkaInterface interface {
	ReadMessage(ctx context.Context) (kafka.Message, error)
	WriteMessages(ctx context.Context, msgs ...kafka.Message) error
	Close() error
}

// Ensure MockReader and MockWriter implement MockKafkaInterface
var (
	_ MockKafkaInterface = &MockKafka{}
)
