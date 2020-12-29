package airtable

import (
	"time"

	"github.com/CedricThomas/22h31-FaisLesBacks/internal/store/model/reminder"
)

func (at *Storer) NewReminder(memoId, title, content string, reminderDate time.Time) (*reminder.Reminder, error) {
	entity := reminder.Reminder{
		Fields: reminder.Fields{
			MemoId:       memoId,
			Title:        title,
			Content:      content,
			ReminderDate: reminderDate,
			Triggered:    false,
		},
	}
	table := at.client.Table(entity.TableName())
	if err := table.Create(&entity); err != nil {
		return nil, err
	}
	return &entity, nil
}
