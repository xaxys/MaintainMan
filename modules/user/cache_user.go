package user

import (
	"fmt"
	"strconv"
)

func cacheGetUserByID(id uint) (*User, error) {
	obj, ok := mctx.Cache.Get(strconv.FormatUint(uint64(id), 36))
	if !ok {
		return nil, fmt.Errorf("未找到用户: id: %d", id)
	}
	user, ok := obj.(User)
	if !ok {
		err := fmt.Errorf("缓存中的用户不是 User 类型: id: %d", id)
		mctx.Logger.Warn(err)
		cacheDeleteUser(id)
		return nil, err
	}
	return &user, nil
}

func cacheSaveUser(user *User) error {
	cacheDeleteUser(user.ID)
	ok := mctx.Cache.Set(strconv.FormatUint(uint64(user.ID), 36), *user, 0)
	if !ok {
		return fmt.Errorf("缓存用户失败: id: %d", user.ID)
	}
	return nil
}

func cacheDeleteUser(id uint) {
	mctx.Cache.Del(strconv.FormatUint(uint64(id), 36))
}
