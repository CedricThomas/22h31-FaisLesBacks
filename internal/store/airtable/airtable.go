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
	if clientErr, ok := err.(airtable.ErrClientRequest); ok {
		return strings.Contains(clientErr.Err.Error(), "MODEL_ID_NOT_FOUND")
	}
	return false
}
