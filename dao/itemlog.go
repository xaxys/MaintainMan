package dao

import "maintainman/model"

func ItemLogAdd(aul *model.AddItemJson) *model.ItemLog {
	itemlog := &model.ItemLog{
		ItemID:      aul.ItemID,
		ChangeNum:   int(aul.Num),
		ChangePrice: aul.Price,
		BaseModel: model.BaseModel{
			CreatedBy: aul.OperatorID,
			UpdatedBy: aul.OperatorID,
		},
	}
	return itemlog
}

func ItemLogConsume(aul *model.ConsumeItemJson) *model.ItemLog {
	itemlog := &model.ItemLog{
		ItemID:      aul.ItemID,
		ChangeNum:   -int(aul.Num),
		ChangePrice: -aul.Price,
		BaseModel: model.BaseModel{
			CreatedBy: aul.OperatorID,
			UpdatedBy: aul.OperatorID,
		},
	}
	return itemlog
}
