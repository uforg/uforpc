package lsp

type NotificationMessageTextDocumentDidChange struct {
	NotificationMessage
	Params NotificationMessageTextDocumentDidChangeParams `json:"params"`
}

type NotificationMessageTextDocumentDidChangeParams struct {
	// The document that did change.
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	// The content of the document.
	ContentChanges []TextDocumentContentChangeEvent `json:"contentChanges"`
}

func (l *LSP) handleTextDocumentDidChange(rawMessage []byte) error {
	var notification NotificationMessageTextDocumentDidChange
	if err := decode(rawMessage, &notification); err != nil {
		return err
	}

	if err := l.docstore.change(notification.Params); err != nil {
		return err
	}

	l.logger.Info("text document did change", "uri", notification.Params.TextDocument.URI)

	return nil
}
