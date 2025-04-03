package lsp

type NotificationMessageTextDocumentDidOpen struct {
	NotificationMessage
	Params NotificationMessageTextDocumentDidOpenParams `json:"params"`
}

type NotificationMessageTextDocumentDidOpenParams struct {
	// The document that was opened.
	TextDocument TextDocumentItem `json:"textDocument"`
}

func (l *LSP) handleTextDocumentDidOpen(rawMessage []byte) (any, error) {
	var notification NotificationMessageTextDocumentDidOpen
	if err := decode(rawMessage, &notification); err != nil {
		return nil, err
	}

	if err := l.docstore.open(notification.Params); err != nil {
		return nil, err
	}

	l.logger.Info("text document did open", "uri", notification.Params.TextDocument.URI)

	return nil, nil
}
