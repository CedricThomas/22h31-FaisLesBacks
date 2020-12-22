package store

import (
	"github.com/CedricThomas/22h31-FaisLesBacks/internal/store/model/memo"
	"github.com/CedricThomas/22h31-FaisLesBacks/internal/store/model/subscription"
)

type Store interface {
	// Memo
	NewMemo(title, content, userId string) (*memo.Memo, error)
	GetMemo(memoId string) (*memo.Memo, error)
	ListMemo(userId string) ([]memo.Memo, error)
	UpdateMemo(memoId string, memo *memo.Fields) (*memo.Memo, error)
	DeleteMemo(memoId string) error

	// Subscription
	NewSubscription(registrationId, userId string) (*subscription.Subscription, error)
}
