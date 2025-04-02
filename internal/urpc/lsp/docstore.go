package lsp

import "fmt"

type docstoreItem struct {
	rawText string
}

type docstore struct {
	docs map[string]docstoreItem
}

func newDocstore() *docstore {
	return &docstore{
		docs: make(map[string]docstoreItem),
	}
}

func (d *docstore) open(params NotificationMessageTextDocumentDidOpenParams) error {
	if params.TextDocument.LanguageID != "urpc" {
		return nil
	}

	if _, ok := d.docs[params.TextDocument.URI]; ok {
		return fmt.Errorf("document already exists")
	}

	d.docs[params.TextDocument.URI] = docstoreItem{
		rawText: params.TextDocument.Text,
	}
	return nil
}

func (d *docstore) change(params NotificationMessageTextDocumentDidChangeParams) error {
	if _, ok := d.docs[params.TextDocument.URI]; !ok {
		return fmt.Errorf("document not found")
	}

	if len(params.ContentChanges) == 0 {
		return fmt.Errorf("no content changes")
	}

	lastChange := params.ContentChanges[len(params.ContentChanges)-1]
	newDoc := d.docs[params.TextDocument.URI]
	newDoc.rawText = lastChange.Text
	d.docs[params.TextDocument.URI] = newDoc
	return nil
}

func (d *docstore) close(params NotificationMessageTextDocumentDidCloseParams) error {
	if _, ok := d.docs[params.TextDocument.URI]; !ok {
		return fmt.Errorf("document not found")
	}
	delete(d.docs, params.TextDocument.URI)
	return nil
}

func (d *docstore) get(uri string) (docstoreItem, error) {
	if _, ok := d.docs[uri]; !ok {
		return docstoreItem{}, fmt.Errorf("document not found")
	}
	return d.docs[uri], nil
}
