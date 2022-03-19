package service

import (
	"errors"
	"maintainman/config"
	"maintainman/dao"
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

func GetAllAnnounces(aul *model.AllAnnounceJson) *model.ApiJson {
	if err := util.Validator.Struct(aul); err != nil {
		return model.ErrorVerification(err)
	}
	if aul.Limit == 0 {
		aul.Limit = config.AppConfig.GetInt("app.page_limit_default")
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
		Limit:     config.AppConfig.GetInt("app.page_limit_default"),
		Offset:    int(offset),
		OrderBy:   "id",
	}
	return GetAllAnnounces(aul)
}

func AnnounceToJson(announce *model.Announce) *model.AnnounceJson {
	return &model.AnnounceJson{
		ID:        announce.ID,
		Title:     announce.Title,
		Content:   announce.Content,
		StartTime: announce.StartTime.Unix(),
		EndTime:   announce.EndTime.Unix(),
		CreatedAt: announce.CreatedAt.Unix(),
		UpdatedAt: announce.UpdatedAt.Unix(),
		CreatedBy: announce.CreatedBy,
		UpdatedBy: announce.UpdatedBy,
	}
}
