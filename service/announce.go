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

func GetAnnounceByID(id uint, auth *model.AuthInfo) *model.ApiJson {
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

func GetAnnounceByTitle(title string, auth *model.AuthInfo) *model.ApiJson {
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

func GetAllAnnounces(aul *model.AllAnnounceRequest, auth *model.AuthInfo) *model.ApiJson {
	if err := util.Validator.Struct(aul); err != nil {
		return model.ErrorValidation(err)
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

func GetLatestAnnounces(param *model.PageParam, auth *model.AuthInfo) *model.ApiJson {
	now := time.Now().Unix()
	aul := &model.AllAnnounceRequest{
		StartTime: now,
		EndTime:   now,
		Inclusive: false,
		PageParam: model.PageParam{
			OrderBy: "id desc",
			Offset:  param.Offset,
			Limit:   param.Limit,
		},
	}
	return GetAllAnnounces(aul, auth)
}

func CreateAnnounce(aul *model.CreateAnnounceRequest, auth *model.AuthInfo) *model.ApiJson {
	// TODO: Localize error info: https://blog.xizhibei.me/2019/06/16/an-introduction-to-golang-validator/
	if err := util.Validator.Struct(aul); err != nil {
		return model.ErrorValidation(err)
	}
	req := model.ModifyAnnounceRequest(*aul)
	announce, err := dao.CreateAnnounce(&req, auth.User)
	if err != nil {
		return model.ErrorInsertDatabase(err)
	}
	return model.SuccessCreate(AnnounceToJson(announce), "创建成功")
}

func UpdateAnnounce(id uint, aul *model.UpdateAnnounceRequest, auth *model.AuthInfo) *model.ApiJson {
	if err := util.Validator.Struct(aul); err != nil {
		return model.ErrorValidation(err)
	}
	req := model.ModifyAnnounceRequest(*aul)
	announce, err := dao.UpdateAnnounce(id, &req, auth.User)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrorNotFound(err)
		} else {
			return model.ErrorUpdateDatabase(err)
		}
	}
	return model.SuccessUpdate(AnnounceToJson(announce), "更新成功")
}

func DeleteAnnounce(id uint, auth *model.AuthInfo) *model.ApiJson {
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

func HitAnnounce(id uint, auth *model.AuthInfo) *model.ApiJson {
	key := fmt.Sprintf("%d:%d", id, auth.User)
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
	return util.NilOrValue(announce, &model.AnnounceJson{
		ID:        announce.ID,
		Title:     announce.Title,
		Content:   announce.Content,
		StartTime: announce.StartTime.Unix(),
		EndTime:   announce.EndTime.Unix(),
		Hits:      announce.Hits,
		CreatedAt: announce.CreatedAt.Unix(),
		UpdatedAt: announce.UpdatedAt.Unix(),
	})
}
