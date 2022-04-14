package order

func dbItemLogAdd(aul *AddItemRequest) *ItemLog {
	itemlog := &ItemLog{
		ItemID:      aul.ItemID,
		ChangeNum:   int(aul.Num),
		ChangePrice: aul.Price,
	}
	return itemlog
}

func dbItemLogConsume(aul *ConsumeItemRequest) *ItemLog {
	itemlog := &ItemLog{
		ItemID:      aul.ItemID,
		ChangeNum:   -int(aul.Num),
		ChangePrice: -aul.Price,
	}
	return itemlog
}
