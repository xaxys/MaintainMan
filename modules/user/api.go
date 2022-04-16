package user

// GetUserByID returns the user with the given ID.
func GetUserByID(id uint) (*User, error) {
	return dbGetUserByID(id)
}
