package airtable

import (
	"strings"

	"github.com/brianloveswords/airtable"
)

type Storer struct {
	client airtable.Client
}

func New(apiKey, baseId string) *Storer {
	return &Storer{client: airtable.Client{
		APIKey: apiKey,
		BaseID: baseId,
	}}
}

func isNotFoundErr(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "NOT_FOUND")
}
