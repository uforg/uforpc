package lsp

import (
	"bufio"
	"fmt"
	"io"
	"sync"
)

type LSP struct {
	reader    io.Reader
	writer    io.Writer
	handlerMu sync.Mutex
	logger    *LSPLogger
	docstore  *docstore
}

// New creates a new LSP instance. It uses the given reader and writer to read and write
// messages to the LSP server.
func New(reader io.Reader, writer io.Writer) *LSP {
	return &LSP{
		reader:    reader,
		writer:    writer,
		handlerMu: sync.Mutex{},
		logger:    NewLSPLogger(),
		docstore:  newDocstore(),
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

		shouldExit, err := l.handleMessage(messageBytes, message)
		if err != nil {
			return fmt.Errorf("failed to handle message with method %s and id %s: %w", message.Method, message.ID, err)
		}

		if shouldExit {
			return nil
		}
	}

	return nil
}

func (l *LSP) sendMessage(message any) error {
	messageBytes, err := encode(message)
	if err != nil {
		return fmt.Errorf("failed to encode message: %w", err)
	}

	_, err = l.writer.Write(messageBytes)
	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	if m, ok := message.(ResponseMessage); ok {
		l.logger.Info("message response sent", "id", m.ID, "method", m.Method)
	}
	if m, ok := message.(NotificationMessage); ok {
		l.logger.Info("notification sent", "method", m.Method)
	}

	return nil
}

func (l *LSP) handleMessage(messageBytes []byte, message Message) (bool, error) {
	l.handlerMu.Lock()
	defer l.handlerMu.Unlock()

	if message.ID != "" {
		l.logger.Info("message received", "id", message.ID, "method", message.Method)
	} else {
		l.logger.Info("notification received", "method", message.Method)
	}

	switch message.Method {
	// Lifecycle operations
	case "initialize":
		return false, l.handleInitialize(messageBytes)
	case "initialized":
		return false, l.handleInitialized(messageBytes)
	case "shutdown":
		return false, l.handleShutdown(messageBytes)
	case "exit":
		return true, nil

	// Text document operations
	case "textDocument/didOpen":
		return false, l.handleTextDocumentDidOpen(messageBytes)
	case "textDocument/didChange":
		return false, l.handleTextDocumentDidChange(messageBytes)
	case "textDocument/didClose":
		return false, l.handleTextDocumentDidClose(messageBytes)
	case "textDocument/formatting":
		return false, l.handleTextDocumentFormatting(messageBytes)
	}

	return false, nil
}
