package lsp

import (
	"bufio"
	"fmt"
	"io"
)

type LSP struct {
	reader io.Reader
	writer io.Writer
	logger *LSPLogger
}

// New creates a new LSP instance. It uses the given reader and writer to read and write
// messages to the LSP server.
func New(reader io.Reader, writer io.Writer) *LSP {
	return &LSP{
		reader: reader,
		writer: writer,
		logger: NewLSPLogger(),
	}
}

// Run starts the LSP server. It will read messages from the reader and write responses
// to the writer.
func (l *LSP) Run() error {
	if l.reader == nil || l.writer == nil {
		return fmt.Errorf("reader and writer are required")
	}

	scanner := bufio.NewScanner(l.reader)
	scanner.Split(scannerSplitFunc)

	for scanner.Scan() {
		messageBytes := scanner.Bytes()
		var message Message
		if err := decode(messageBytes, &message); err != nil {
			return fmt.Errorf("failed to decode message or notification: %w", err)
		}
		if err := l.handleMessage(messageBytes, message); err != nil {
			return fmt.Errorf("failed to handle message with method %s and id %s: %w", message.Method, message.ID, err)
		}
	}

	return nil
}

func (l *LSP) handleMessage(messageBytes []byte, message Message) error {
	if message.ID != "" {
		l.logger.Info("message received", "id", message.ID, "method", message.Method)
	} else {
		l.logger.Info("notification received", "method", message.Method)
	}

	switch message.Method {
	case "initialize":
	case "initialized":
	}

	return nil
}
