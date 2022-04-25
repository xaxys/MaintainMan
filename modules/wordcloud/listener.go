package wordcloud

import "github.com/xaxys/maintainman/modules/order"

func listener() {
	defer func() {
		if err := recover(); err != nil {
			mctx.Logger.Errorf("wordcloud listener panic: %s", err)
		}
	}()

	orderModule := mctx.Registry.Get("order")
	if orderModule == nil {
		mctx.Logger.Errorf("order module not found")
		return
	}

	for {
		select {
		// order created
		case ch := <-mctx.EventBus.On("order:create"):
			orderID, _ := ch.Args[0].(uint)
			odr, err := order.GetOrderByID(orderID)
			if err != nil {
				mctx.Logger.Warnf("get order failed: %s", err)
				continue
			}
			uploadWordsService(odr.ID, odr.Title)
			uploadWordsService(odr.ID, odr.Content)
		// order title changed
		case ch := <-mctx.EventBus.On("order:update:title"):
			orderID, _ := ch.Args[0].(uint)
			odr, err := order.GetOrderByID(orderID)
			if err != nil {
				mctx.Logger.Warnf("get order failed: %s", err)
				continue
			}
			uploadWordsService(odr.ID, odr.Title)
		// order content changed
		case ch := <-mctx.EventBus.On("order:update:content"):
			orderID, _ := ch.Args[0].(uint)
			odr, err := order.GetOrderByID(orderID)
			if err != nil {
				mctx.Logger.Warnf("get order failed: %s", err)
				continue
			}
			uploadWordsService(odr.ID, odr.Content)
		// order comment
		case ch := <-mctx.EventBus.On("order:update:comment"):
			commentID, _ := ch.Args[1].(uint)
			comment, err := order.GetCommentByID(commentID)
			if err != nil {
				mctx.Logger.Warnf("get comment failed: %s", err)
				continue
			}
			uploadWordsService(comment.OrderID, comment.Content)
		}
	}
}
