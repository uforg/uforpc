package lsp

import "fmt"

func (l *LSP) handleShutdown(rawMessage []byte) error {
	var message RequestMessage
	if err := decode(rawMessage, &message); err != nil {
		return fmt.Errorf("failed to decode shutdown message: %w", err)
	}

	l.logger.Info("shutdown message received", "id", message.ID, "method", message.Method)

	response := ResponseMessage{
		Message: DefaultMessage,
		ID:      message.ID,
	}

	if err := l.sendMessage(response); err != nil {
		return fmt.Errorf("failed to send shutdown response: %w", err)
	}

	return nil
}
