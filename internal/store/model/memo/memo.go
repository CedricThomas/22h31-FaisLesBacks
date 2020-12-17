package memo

import "github.com/brianloveswords/airtable"

type (
	Fields struct {
		Title   string `json:"title"`
		Content string `json:"content"`
		UserId  string `json:"user_id"`
	}
	Memo struct {
		airtable.Record
		Fields Fields
	}
)

func (Memo) TableName() string {
	return "memo"
}
