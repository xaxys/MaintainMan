package announce

import (
	"errors"
	"fmt"
	"time"

	"github.com/xaxys/maintainman/core/model"
	"github.com/xaxys/maintainman/core/util"

	"gorm.io/gorm"
)

func getAnnounceService(id uint, auth *model.AuthInfo) *model.ApiJson {
	announce, err := dbGetAnnounceByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrorNotFound(err)
		}
		return model.ErrorQueryDatabase(err)
	}
	return model.Success(announceToJson(announce), "获取成功")
}

func getAnounceByTitleService(title string, auth *model.AuthInfo) *model.ApiJson {
	announce, err := dbGetAnnounceByTitle(title)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrorNotFound(err)
		}
		return model.ErrorQueryDatabase(err)
	}
	return model.Success(announceToJson(announce), "获取成功")
}

func getAllAnnouncesService(aul *AllAnnounceRequest, auth *model.AuthInfo) *model.ApiJson {
	if err := util.Validator.Struct(aul); err != nil {
		return model.ErrorValidation(err)
	}
	announces, count, err := dbGetAllAnnouncesWithParam(aul)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrorNotFound(err)
		}
		return model.ErrorQueryDatabase(err)
	}
	as := util.TransSlice(announces, announceToJson)
	return model.SuccessPaged(as, count, "获取成功")
}

func getLatestAnnouncesService(param *model.PageParam, auth *model.AuthInfo) *model.ApiJson {
	now := time.Now().Unix()
	aul := &AllAnnounceRequest{
		StartTime: now,
		EndTime:   now,
		Inclusive: false,
		PageParam: model.PageParam{
			OrderBy: "id desc",
			Offset:  param.Offset,
			Limit:   param.Limit,
		},
	}
	return getAllAnnouncesService(aul, auth)
}

func createAnnounceService(aul *CreateAnnounceRequest, auth *model.AuthInfo) *model.ApiJson {
	// TODO: Localize error info: https://blog.xizhibei.me/2019/06/16/an-introduction-to-golang-validator/
	if err := util.Validator.Struct(aul); err != nil {
		return model.ErrorValidation(err)
	}
	req := ModifyAnnounceRequest(*aul)
	announce, err := dbCreateAnnounce(&req, auth.User)
	if err != nil {
		return model.ErrorInsertDatabase(err)
	}
	return model.SuccessCreate(announceToJson(announce), "创建成功")
}

func updateAnnounceService(id uint, aul *UpdateAnnounceRequest, auth *model.AuthInfo) *model.ApiJson {
	if err := util.Validator.Struct(aul); err != nil {
		return model.ErrorValidation(err)
	}
	req := ModifyAnnounceRequest(*aul)
	announce, err := dbUpdateAnnounce(id, &req, auth.User)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrorNotFound(err)
		}
		return model.ErrorUpdateDatabase(err)
	}
	return model.SuccessUpdate(announceToJson(announce), "更新成功")
}

func deleteAnnounceService(id uint, auth *model.AuthInfo) *model.ApiJson {
	err := dbDeleteAnnounce(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrorNotFound(err)
		}
		return model.ErrorDeleteDatabase(err)
	}
	return model.SuccessUpdate(nil, "删除成功")
}

func hitAnnounceService(id uint, auth *model.AuthInfo) *model.ApiJson {
	key := fmt.Sprintf("%d:%d", id, auth.User)
	if _, ok := mctx.Cache.Get(key); ok {
		return model.Success(nil, "浏览过了")
	}
	expire, err := time.ParseDuration(announceConfig.GetString("hit_expire"))
	if err != nil {
		return model.ErrorInternalServer(err)
	}
	announce, err := dbGetAnnounceByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrorNotFound(err)
		}
		return model.ErrorQueryDatabase(err)
	}
	if time.Now().Before(*announce.StartTime) || time.Now().After(*announce.EndTime) {
		return model.ErrorNotFound(errors.New("不在公告期间"))
	}
	mctx.Cache.Set(key, nil, expire)
	if err := dbHitAnnounce(id); err != nil {
		return model.ErrorUpdateDatabase(err)
	}
	return model.SuccessUpdate(nil, "浏览成功")
}

func announceToJson(announce *Announce) *AnnounceJson {
	if announce == nil {
		return nil
	} else {
		return &AnnounceJson{
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
}
