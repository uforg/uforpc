package lsp

type NotificationMessageTextDocumentDidClose struct {
	NotificationMessage
	Params NotificationMessageTextDocumentDidCloseParams `json:"params"`
}

type NotificationMessageTextDocumentDidCloseParams struct {
	// The document that did close.
	TextDocument TextDocumentIdentifier `json:"textDocument"`
}

func (l *LSP) handleTextDocumentDidClose(rawMessage []byte) (any, error) {
	var notification NotificationMessageTextDocumentDidClose
	if err := decode(rawMessage, &notification); err != nil {
		return nil, err
	}

	if err := l.docstore.close(notification.Params); err != nil {
		return nil, err
	}

	l.logger.Info("text document did close", "uri", notification.Params.TextDocument.URI)

	return nil, nil
}
