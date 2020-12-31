package memo

import (
	"github.com/brianloveswords/airtable"

	"github.com/CedricThomas/22h31-FaisLesBacks/api/model"
)

type (
	Fields struct {
		Title    string `json:"title"`
		Content  string `json:"content"`
		UserId   string `json:"user_id"`
		Location string `json:"location"`
	}
	Memo struct {
		airtable.Record
		Fields Fields
	}
)

func (Memo) TableName() string {
	return "memo"
}

func (m *Memo) ToModel() *model.Memo {
	return &model.Memo{
		Id:        m.ID,
		CreatedAt: m.CreatedTime,
		Title:     m.Fields.Title,
		Content:   m.Fields.Content,
		Location:  m.Fields.Location,
	}
}
