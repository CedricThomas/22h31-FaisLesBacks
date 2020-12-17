package airtable

import "github.com/brianloveswords/airtable"

type Storer struct {
	client airtable.Client
}

func New(apiKey, baseId string) *Storer {
	return &Storer{client: airtable.Client{
		APIKey: apiKey,
		BaseID: baseId,
	}}
}
