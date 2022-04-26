package user

// GetUserByID returns the user with the given ID.
func GetUserByID(id uint) (*User, error) {
	return dbGetUserByID(id)
}

// UserToJson converts a user to a json string.
func UserToJson(user *User) *UserJson {
	return userToJson(user)
}
