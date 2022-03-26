package dao

import (
	"maintainman/model"
)

func ItemLogAdd(aul *model.AddItemRequest) *model.ItemLog {
	itemlog := &model.ItemLog{
		ItemID:      aul.ItemID,
		ChangeNum:   int(aul.Num),
		ChangePrice: aul.Price,
	}
	return itemlog
}

func ItemLogConsume(aul *model.ConsumeItemRequest) *model.ItemLog {
	itemlog := &model.ItemLog{
		ItemID:      aul.ItemID,
		ChangeNum:   -int(aul.Num),
		ChangePrice: -aul.Price,
	}
	return itemlog
}
