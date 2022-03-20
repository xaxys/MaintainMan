package service

import (
	"errors"
	"fmt"
	"maintainman/config"
	"maintainman/dao"
	"maintainman/database"
	"maintainman/model"
	"maintainman/util"
	"time"

	"gorm.io/gorm"
)

func GetAnnounceByID(id uint) *model.ApiJson {
	announce, err := dao.GetAnnounceByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrorNotFound(err)
		} else {
			return model.ErrorQueryDatabase(err)
		}
	}
	return model.Success(AnnounceToJson(announce), "获取成功")
}

func GetAnnounceByTitle(title string) *model.ApiJson {
	announce, err := dao.GetAnnounceByTitle(title)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrorNotFound(err)
		} else {
			return model.ErrorQueryDatabase(err)
		}
	}
	return model.Success(AnnounceToJson(announce), "获取成功")
}

func GetAllAnnounces(aul *model.AllAnnounceJson) *model.ApiJson {
	if err := util.Validator.Struct(aul); err != nil {
		return model.ErrorVerification(err)
	}
	announces, err := dao.GetAllAnnouncesWithParam(aul)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrorNotFound(err)
		} else {
			return model.ErrorQueryDatabase(err)
		}
	}
	as := util.TransSlice(announces, AnnounceToJson)
	return model.Success(as, "获取成功")
}

func GetLatestAnnounces(offset uint) *model.ApiJson {
	now := time.Now().Unix()
	aul := &model.AllAnnounceJson{
		StartTime: now,
		EndTime:   now,
		Inclusive: false,
		Offset:    offset,
		OrderBy:   "id",
	}
	return GetAllAnnounces(aul)
}

func CreateAnnounce(aul *model.ModifyAnnounceJson) *model.ApiJson {
	// TODO: Localize error info: https://blog.xizhibei.me/2019/06/16/an-introduction-to-golang-validator/
	if err := util.Validator.Struct(aul); err != nil {
		return model.ErrorVerification(err)
	}
	announce, err := dao.CreateAnnounce(aul)
	if err != nil {
		return model.ErrorInsertDatabase(err)
	}
	return model.SuccessCreate(AnnounceToJson(announce), "创建成功")
}

func UpdateAnnounce(id uint, aul *model.ModifyAnnounceJson) *model.ApiJson {
	if err := util.Validator.Struct(aul); err != nil {
		return model.ErrorVerification(err)
	}
	announce, err := dao.UpdateAnnounce(id, aul)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrorNotFound(err)
		} else {
			return model.ErrorUpdateDatabase(err)
		}
	}
	return model.SuccessUpdate(AnnounceToJson(announce), "更新成功")
}

func DeleteAnnounce(id uint) *model.ApiJson {
	err := dao.DeleteAnnounce(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrorNotFound(err)
		} else {
			return model.ErrorDeleteDatabase(err)
		}
	}
	return model.SuccessUpdate(nil, "删除成功")
}

func HitAnnounce(id, uid uint) *model.ApiJson {
	key := fmt.Sprintf("%d:%d", id, uid)
	if _, ok := database.Cache.Get(key); ok {
		return model.Success(nil, "浏览过了")
	}
	expire := time.Duration(config.AppConfig.GetInt("app.hit_expire_time")) * time.Second
	database.Cache.Set(key, nil, expire)
	if err := dao.HitAnnounce(id); err != nil {
		return model.ErrorUpdateDatabase(err)
	}
	return model.SuccessUpdate(nil, "浏览成功")
}

func AnnounceToJson(announce *model.Announce) *model.AnnounceJson {
	return &model.AnnounceJson{
		ID:        announce.ID,
		Title:     announce.Title,
		Content:   announce.Content,
		StartTime: announce.StartTime.Unix(),
		EndTime:   announce.EndTime.Unix(),
		Hits:      announce.Hits,
		CreatedAt: announce.CreatedAt.Unix(),
		UpdatedAt: announce.UpdatedAt.Unix(),
	}
}
