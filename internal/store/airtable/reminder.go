package airtable

import (
	"fmt"
	"time"

	"github.com/brianloveswords/airtable"

	"github.com/CedricThomas/22h31-FaisLesBacks/internal/store/model"
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

func (at *Storer) GetReminder(reminderId string) (*reminder.Reminder, error) {
	var entity reminder.Reminder
	table := at.client.Table(entity.TableName())
	if err := table.Get(reminderId, &entity); isNotFoundErr(err) {
		return nil, model.NoSuchEntity
	} else if err != nil {
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

func (at *Storer) UpdateReminder(reminderId string, reminder *reminder.Reminder) (*reminder.Reminder, error) {
	return nil, nil
}

func (at *Storer) DeleteReminder(reminderId string) error {
	table := at.client.Table(reminder.Reminder{}.TableName())
	if err := table.Delete(&reminder.Reminder{
		Record: airtable.Record{
			ID: reminderId,
		},
	}); isNotFoundErr(err) {
		return model.NoSuchEntity
	} else if err != nil {
		return err
	}
	return nil
}
