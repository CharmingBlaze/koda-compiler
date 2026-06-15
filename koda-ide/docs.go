package main

import "koda/api"

// DocPage is one markdown help page in the SDK.
type DocPage = api.DocPage

// ListDocPages returns all markdown documentation shipped with the SDK.
func (a *App) ListDocPages() ([]DocPage, error) {
	a.bootstrapSDK()
	return api.ListDocPages()
}

// ReadDocPage reads markdown content for rel (e.g. docs/faq.md).
func (a *App) ReadDocPage(rel string) (string, error) {
	a.bootstrapSDK()
	return api.ReadDocPage(rel)
}
