package order

// GetSimpleOrderByID returns the order with the given ID.
func GetSimpleOrderByID(id uint) (*Order, error) {
	return dbGetSimpleOrderByID(id)
}

// GetOrderByID returns the order with the given ID and Tags and Comments.
func GetOrderByID(id uint) (*Order, error) {
	return dbGetOrderByID(id)
}

// GetOrderWithLastStatus returns the order with the given ID and the last status.
func GetOrderWithLastStatus(id uint) (*Order, error) {
	return dbGetOrderWithLastStatus(id)
}

// GetCommentByID returns the comment with the given ID.
func GetCommentByID(id uint) (*Comment, error) {
	return dbGetCommentByID(id)
}
