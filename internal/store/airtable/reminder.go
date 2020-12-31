package airtable

import (
	"fmt"
	"time"

	"github.com/brianloveswords/airtable"

	"github.com/CedricThomas/22h31-FaisLesBacks/internal/store/model"
	"github.com/CedricThomas/22h31-FaisLesBacks/internal/store/model/reminder"
)

func (at *Storer) NewReminder(memoId, userId, title, content string, reminderDate time.Time) (*reminder.Reminder, error) {
	entity := reminder.Reminder{
		Fields: reminder.Fields{
			MemoId:       memoId,
			UserId:       userId,
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

func (at *Storer) UpdateReminder(fields *reminder.Reminder) (*reminder.Reminder, error) {
	table := at.client.Table(fields.TableName())
	if err := table.Update(fields); isNotFoundErr(err) {
		return nil, model.NoSuchEntity
	} else if err != nil {
		return nil, err
	}
	return fields, nil
}

func (at *Storer) DeleteAllReminder(memoId string) error {
	reminders, err := at.ListReminder(memoId)
	if err != nil {
		return err
	}
	for _, rem := range reminders {
		if err := at.DeleteReminder(rem.ID); err != nil {
			return err
		}
	}
	return nil
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

func (at *Storer) ListReminderToTrigger() ([]reminder.Reminder, error) {
	table := at.client.Table(reminder.Reminder{}.TableName())
	var entities []reminder.Reminder
	if err := table.List(&entities, &airtable.Options{
		Filter: fmt.Sprintf("AND({reminder_date} <= NOW(), {triggered} = FALSE())"),
	}); err != nil {
		return nil, err
	}
	return entities, nil
}
