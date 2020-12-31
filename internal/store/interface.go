package store

import (
	"time"

	"github.com/CedricThomas/22h31-FaisLesBacks/internal/store/model/memo"
	"github.com/CedricThomas/22h31-FaisLesBacks/internal/store/model/reminder"
	"github.com/CedricThomas/22h31-FaisLesBacks/internal/store/model/subscription"
)

type Store interface {
	// Memo
	NewMemo(title, content, location, userId string) (*memo.Memo, error)
	GetMemo(memoId string) (*memo.Memo, error)
	ListMemo(userId string) ([]memo.Memo, error)
	UpdateMemo(memo *memo.Memo) (*memo.Memo, error)
	DeleteMemo(memoId string) error

	// Subscription
	NewSubscription(registrationId, userId string) (*subscription.Subscription, error)
	GetSubscription(registrationId string) (*subscription.Subscription, error)
	ListSubscription(userId string) ([]subscription.Subscription, error)
	DeleteSubscription(registrationId string) error

	// Reminder
	NewReminder(memoId, userId, title, content string, reminderDate time.Time) (*reminder.Reminder, error)
	GetReminder(reminderId string) (*reminder.Reminder, error)
	ListReminder(memoId string) ([]reminder.Reminder, error)
	ListReminderToTrigger() ([]reminder.Reminder, error)
	UpdateReminder(reminderId string, reminder *reminder.Fields) (*reminder.Reminder, error)
	DeleteAllReminder(memoId string) error
	DeleteReminder(reminderId string) error
}
