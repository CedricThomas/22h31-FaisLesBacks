package reminder

import (
	"time"

	"github.com/brianloveswords/airtable"

	"github.com/CedricThomas/22h31-FaisLesBacks/api/model"
)

type (
	Fields struct {
		MemoId       string    `json:"memo_id"`
		Title        string    `json:"title"`
		Content      string    `json:"content"`
		ReminderDate time.Time `json:"reminder_date"`
		Triggered    bool      `json:"triggered"`
	}
	Reminder struct {
		airtable.Record
		Fields Fields
	}
)

func (Reminder) TableName() string {
	return "reminder"
}

func (m *Reminder) ToModel() *model.Reminder {
	return &model.Reminder{
		Id:        m.ID,
		CreatedAt: m.CreatedTime,
		Title:     m.Fields.Title,
		Content:   m.Fields.Content,
		Date:      m.Fields.ReminderDate,
		Triggered: m.Fields.Triggered,
	}
}
