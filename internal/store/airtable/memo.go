package airtable

import (
	"errors"
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
	table := at.client.Table(memo.Memo{}.TableName())
	var entity memo.Memo
	if err := table.Get(memoId, &entity); err != nil {
		if clientErr, ok := err.(airtable.ErrClientRequest); ok {
			_ = clientErr

		}
		return nil, err
	}
	return &entity, nil
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

func (at *Storer) UpdateMemo(memoId string, toUpdate *memo.Fields) (*memo.Memo, error) {
	table := at.client.Table(memo.Memo{}.TableName())
	if toUpdate == nil {
		return nil, errors.New("invalid update fields: nil received")
	}
	entity := memo.Memo{
		Record: airtable.Record{
			ID: memoId,
		},
		Fields: *toUpdate,
	}
	if err := table.Update(&entity); err != nil {
		return nil, err
	}
	return &entity, nil
}

func (at *Storer) DeleteMemo(memoId string) error {
	table := at.client.Table(memo.Memo{}.TableName())
	return table.Delete(&memo.Memo{
		Record: airtable.Record{
			ID: memoId,
		},
	})
}
