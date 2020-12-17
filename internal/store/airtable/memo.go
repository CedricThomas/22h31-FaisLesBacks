package airtable

import (
	"fmt"

	"github.com/brianloveswords/airtable"

	"github.com/CedricThomas/22h31-FaisLesBacks/internal/store/model/memo"
)

func (at *Storer) NewMemo(title, content, userId string) (*memo.Memo, error) {
	entity := memo.Memo{
		Fields: memo.Fields{
			Title:   title,
			Content: content,
			UserId:  userId,
		},
	}
	table := at.client.Table(entity.TableName())
	if err := table.Create(&entity); err != nil {
		return nil, err
	}
	return &entity, nil
}

func (at *Storer) GetMemo(memoId string) (*memo.Memo, error) {
	panic("implement me")
}

func (at *Storer) ListMemo(userId string) ([]memo.Memo, error) {
	table := at.client.Table(memo.Memo{}.TableName())
	var entities []memo.Memo
	if err := table.List(&entities, &airtable.Options{
		Filter: fmt.Sprintf("{user_id} = \"%s\"", userId),
	}); err != nil {
		return nil, err
	}
	return entities, nil
}

func (at *Storer) UpdateMemo(memoId string, memo *memo.Memo) (*memo.Memo, error) {
	panic("implement me")
}

func (at *Storer) DeleteMemo(memoId string) (*memo.Memo, error) {
	panic("implement me")
}
