package lsp

import (
	"fmt"

	"github.com/uforg/uforpc/internal/version"
)

type RequestMessageInitialize struct {
	RequestMessage
	Params RequestMessageInitializeParams `json:"params"`
}

type RequestMessageInitializeParams struct {
	ClientInfo struct {
		Name    string `json:"name"`
		Version string `json:"version,omitzero,omitempty"`
	} `json:"clientInfo,omitzero"`
}

type ResponseMessageInitialize struct {
	ResponseMessage
	Result ResponseMessageInitializeResult `json:"result"`
}

type ResponseMessageInitializeResult struct {
	ServerInfo   ResponseMessageInitializeResultServerInfo   `json:"serverInfo"`
	Capabilities ResponseMessageInitializeResultCapabilities `json:"capabilities"`
}

type ResponseMessageInitializeResultServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type ResponseMessageInitializeResultCapabilities struct {
	DocumentFormattingProvider bool `json:"documentFormattingProvider"`
	TextDocumentSync           int  `json:"textDocumentSync"`
	// Advertise diagnostic capabilities
	DiagnosticProvider bool `json:"diagnosticProvider,omitempty"`
}

func (l *LSP) handleInitialize(rawMessage []byte) (any, error) {
	var message RequestMessageInitialize
	if err := decode(rawMessage, &message); err != nil {
		return nil, fmt.Errorf("failed to decode initialize message: %w", err)
	}

	l.logger.Info(
		"initialize message received",
		"id", message.ID,
		"method", message.Method,
		"clientName", message.Params.ClientInfo.Name,
		"clientVersion", message.Params.ClientInfo.Version,
	)

	response := ResponseMessageInitialize{
		ResponseMessage: ResponseMessage{
			Message: DefaultMessage,
			ID:      message.ID,
		},
		Result: ResponseMessageInitializeResult{
			ServerInfo: ResponseMessageInitializeResultServerInfo{
				Name:    "UFO RPC Language Server",
				Version: version.VersionWithPrefix,
			},
			Capabilities: ResponseMessageInitializeResultCapabilities{
				// Documents are synced by always sending the full content of the document.
				TextDocumentSync: 1,
				// Document formatting is supported.
				DocumentFormattingProvider: true,
				// Diagnostics are supported if analyzer is available
				DiagnosticProvider: l.analyzer != nil,
			},
		},
	}

	return response, nil
}
