package wordcloud

import "github.com/xaxys/maintainman/modules/order"

func listener() {
	defer func() {
		if err := recover(); err != nil {
			mctx.Logger.Errorf("wordcloud listener panic: %s", err)
		}
	}()

	for {
		select {
		// order created
		case ch := <-mctx.EventBus.On("order:create"):
			orderID, _ := ch.Args[0].(uint)
			odr, err := order.GetOrderByID(orderID)
			if err != nil {
				mctx.Logger.Warnf("Get order failed: %s", err)
				continue
			}
			if res := uploadWordsService(odr.ID, odr.Title); !res.Status {
				mctx.Logger.Warnf("Upload words failed: [order: %d, content: %s] errors: %v", odr.ID, odr.Title, res.Data)
			} else {
				mctx.Logger.Infof("Upload words success: [order: %d, content: %s]", odr.ID, odr.Title)
			}
			if res := uploadWordsService(odr.ID, odr.Content); !res.Status {
				mctx.Logger.Warnf("Upload words failed: [order: %d, content: %s] errors: %v", odr.ID, odr.Content, res.Data)
			} else {
				mctx.Logger.Infof("Upload words success: [order: %d, content: %s]", odr.ID, odr.Content)
			}
		// order title changed
		case ch := <-mctx.EventBus.On("order:update:title"):
			orderID, _ := ch.Args[0].(uint)
			odr, err := order.GetOrderByID(orderID)
			if err != nil {
				mctx.Logger.Warnf("Get order failed: %s", err)
				continue
			}
			if res := uploadWordsService(odr.ID, odr.Title); !res.Status {
				mctx.Logger.Warnf("Upload words failed: [order: %d, content: %s] errors: %v", odr.ID, odr.Title, res.Data)
			} else {
				mctx.Logger.Infof("Upload words success: [order: %d, content: %s]", odr.ID, odr.Title)
			}
		// order content changed
		case ch := <-mctx.EventBus.On("order:update:content"):
			orderID, _ := ch.Args[0].(uint)
			odr, err := order.GetOrderByID(orderID)
			if err != nil {
				mctx.Logger.Warnf("Get order failed: %s", err)
				continue
			}
			if res := uploadWordsService(odr.ID, odr.Content); !res.Status {
				mctx.Logger.Warnf("Upload words failed: [order: %d, content: %s] errors: %v", odr.ID, odr.Content, res.Data)
			} else {
				mctx.Logger.Infof("Upload words success: [order: %d, content: %s]", odr.ID, odr.Content)
			}
		// order comment
		case ch := <-mctx.EventBus.On("order:update:comment"):
			commentID, _ := ch.Args[1].(uint)
			comment, err := order.GetCommentByID(commentID)
			if err != nil {
				mctx.Logger.Warnf("Get comment failed: %s", err)
				continue
			}
			if res := uploadWordsService(comment.OrderID, comment.Content); !res.Status {
				mctx.Logger.Warnf("Upload words failed: [order: %d, content: %s] errors: %v", comment.OrderID, comment.Content, res.Data)
			} else {
				mctx.Logger.Infof("Upload words success: [order: %d, content: %s]", comment.OrderID, comment.Content)
			}
		}
	}
}
