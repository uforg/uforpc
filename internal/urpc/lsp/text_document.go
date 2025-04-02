package lsp

// TextDocumentIdentifier: Text documents are identified using a URI.
type TextDocumentIdentifier struct {
	// The text document's URI.
	URI string `json:"uri"`
}

// TextDocumentItem: An item to transfer a text document from the client to the server.
type TextDocumentItem struct {
	// The text document's URI.
	URI string `json:"uri"`
	// The text document's language identifier.
	LanguageID string `json:"languageId"`
	// The version number of this document (it will increase after each change, including undo/redo).
	Version int `json:"version"`
	// The content of the opened text document.
	Text string `json:"text"`
}

// TextDocumentContentChangeEvent: An event describing a change to a text document. If only
// a text is provided it is considered to be the full content of the document.
type TextDocumentContentChangeEvent struct {
	// The new text of the whole document.
	Text string `json:"text"`
}
