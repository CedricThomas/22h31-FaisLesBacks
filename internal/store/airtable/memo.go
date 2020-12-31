package airtable

import (
	"fmt"

	"github.com/brianloveswords/airtable"

	"github.com/CedricThomas/22h31-FaisLesBacks/internal/store/model"
	"github.com/CedricThomas/22h31-FaisLesBacks/internal/store/model/memo"
)

func (at *Storer) NewMemo(title, content, location, userId string) (*memo.Memo, error) {
	entity := memo.Memo{
		Fields: memo.Fields{
			Title:    title,
			Content:  content,
			UserId:   userId,
			Location: location,
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
	if err := table.Get(memoId, &entity); isNotFoundErr(err) {
		return nil, model.NoSuchEntity
	} else if err != nil {
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

func (at *Storer) UpdateMemo(toUpdate *memo.Memo) (*memo.Memo, error) {
	table := at.client.Table(memo.Memo{}.TableName())
	if err := table.Update(toUpdate); isNotFoundErr(err) {
		return nil, model.NoSuchEntity
	} else if err != nil {
		return nil, err
	}
	return toUpdate, nil
}

func (at *Storer) DeleteMemo(memoId string) error {
	table := at.client.Table(memo.Memo{}.TableName())
	if err := table.Delete(&memo.Memo{
		Record: airtable.Record{
			ID: memoId,
		},
	}); isNotFoundErr(err) {
		return model.NoSuchEntity
	} else if err != nil {
		return err
	}
	return nil
}
