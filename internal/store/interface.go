package store

import "github.com/CedricThomas/22h31-FaisLesBacks/internal/store/model/memo"

type Store interface {
	NewMemo(title, content, userId string) (*memo.Memo, error)
	GetMemo(memoId string) (*memo.Memo, error)
	ListMemo(userId string) ([]memo.Memo, error)
	UpdateMemo(memoId string, memo *memo.Memo) (*memo.Memo, error)
	DeleteMemo(memoId string) (*memo.Memo, error)
}
