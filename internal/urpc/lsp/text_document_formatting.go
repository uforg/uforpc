package lsp

import (
	"fmt"
	"strings"

	"github.com/uforg/uforpc/internal/urpc/formatter"
)

type RequestMessageTextDocumentFormatting struct {
	RequestMessage
	Params RequestMessageTextDocumentFormattingParams `json:"params"`
}

type RequestMessageTextDocumentFormattingParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	// Options are not used because the formatting rules are not configurable
}

type ResponseMessageTextDocumentFormatting struct {
	ResponseMessage
	Result *[]TextDocumentTextEdit `json:"result,omitempty"`
}

func (l *LSP) handleTextDocumentFormatting(rawMessage []byte) error {
	var request RequestMessageTextDocumentFormatting
	if err := decode(rawMessage, &request); err != nil {
		return fmt.Errorf("failed to decode text document formatting request: %w", err)
	}

	doc, err := l.docstore.get(request.Params.TextDocument.URI)
	if err != nil {
		return fmt.Errorf("failed to get text document: %w", err)
	}

	response := ResponseMessageTextDocumentFormatting{
		ResponseMessage: ResponseMessage{
			Message: DefaultMessage,
			ID:      request.ID,
		},
		Result: &[]TextDocumentTextEdit{},
	}

	formattedText, err := formatter.Format(doc.rawText)
	if err != nil {
		return l.sendMessage(response)
	}

	lines := strings.Split(doc.rawText, "\n")
	lastLine := max(len(lines)-1, 0)
	lastLineChar := max(len(lines[lastLine])-1, 0)

	response.Result = &[]TextDocumentTextEdit{
		{
			Range: TextDocumentRange{
				Start: TextDocumentPosition{Line: 0, Character: 0},
				End:   TextDocumentPosition{Line: lastLine, Character: lastLineChar},
			},
			NewText: formattedText,
		},
	}

	return l.sendMessage(response)
}
