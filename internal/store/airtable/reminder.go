package airtable

import (
	"fmt"
	"time"

	"github.com/brianloveswords/airtable"

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

func (at *Storer) ListReminder(memoId string) ([]reminder.Reminder, error) {
	table := at.client.Table(reminder.Reminder{}.TableName())
	var entities []reminder.Reminder
	if err := table.List(&entities, &airtable.Options{
		Filter: fmt.Sprintf("{memo_id} = \"%s\"", memoId),
	}); err != nil {
		return nil, err
	}
	return entities, nil
}
