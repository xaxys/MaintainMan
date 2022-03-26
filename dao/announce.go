package dao

import (
	"maintainman/database"
	"maintainman/logger"
	"maintainman/model"
	"time"

	"gorm.io/gorm"
)

func GetAnnounceByID(id uint) (*model.Announce, error) {
	announce := &model.Announce{}
	if err := database.DB.First(announce, id).Error; err != nil {
		logger.Logger.Debugf("GetAnnounceByIDErr: %v\n", err)
		return nil, err
	}
	return announce, nil
}

func GetAnnounceByTitle(title string) (*model.Announce, error) {
	announce := &model.Announce{Title: title}
	if err := database.DB.Where(announce).First(announce).Error; err != nil {
		logger.Logger.Debugf("GetAnnounceByNameErr: %v\n", err)
		return nil, err
	}
	return announce, nil
}

func GetAllAnnouncesWithParam(aul *model.AllAnnounceRequest) (announces []*model.Announce, err error) {
	db := Filter(aul.OrderBy, aul.Offset, aul.Limit)
	if aul.Title != "" {
		db = db.Where("title like ?", aul.Title)
	}
	if aul.StartTime != -1 {
		time := time.Unix(aul.StartTime, 0)
		if aul.Inclusive {
			db = db.Where("start_time >= ?", time)
		} else {
			db = db.Where("start_time <= ?", time)
		}
	}
	if aul.EndTime != -1 {
		time := time.Unix(aul.EndTime, 0)
		if aul.Inclusive {
			db = db.Where("end_time <= ?", time)
		} else {
			db = db.Where("end_time >= ?", time)
		}
	}
	if err = db.Find(&announces).Error; err != nil {
		logger.Logger.Debugf("GetAllAnnounceErr: %v\n", err)
	}
	return
}

func CreateAnnounce(json *model.ModifyAnnounceRequest, operator uint) (*model.Announce, error) {
	announce := JsonToAnnounce(json)
	if announce.StartTime == nil {
		now := time.Now()
		announce.StartTime = &now
	}
	if announce.EndTime == nil {
		now := time.Unix(253370764799, 0)
		announce.EndTime = &now
	}
	announce.CreatedBy = operator

	if err := database.DB.Create(announce).Error; err != nil {
		logger.Logger.Debugf("CreateAnnounceErr: %v\n", err)
		return nil, err
	}

	return announce, nil
}

func UpdateAnnounce(id uint, json *model.ModifyAnnounceRequest, operator uint) (*model.Announce, error) {
	announce := JsonToAnnounce(json)
	announce.ID = id
	announce.UpdatedBy = operator

	if err := database.DB.Model(announce).Updates(announce).Error; err != nil {
		logger.Logger.Debugf("UpdateAnnounceErr: %v\n", err)
		return nil, err
	}

	return announce, nil
}

func DeleteAnnounce(id uint) error {
	if err := database.DB.Delete(&model.Announce{}, id).Error; err != nil {
		logger.Logger.Debugf("DeleteAnnounceErr: %v\n", err)
		return err
	}
	return nil
}

func HitAnnounce(id uint) error {
	announce := &model.Announce{}
	announce.ID = id
	if err := database.DB.Model(announce).Update("hits", gorm.Expr("hits + ?", 1)).Error; err != nil {
		logger.Logger.Debugf("HitAnnounceErr: %v\n", err)
		return err
	}
	return nil
}

func JsonToAnnounce(json *model.ModifyAnnounceRequest) (ret *model.Announce) {
	ret = &model.Announce{
		Title:   json.Title,
		Content: json.Content,
	}
	if json.StartTime != -1 {
		time := time.Unix(json.StartTime, 0)
		ret.StartTime = &time
	}
	if json.EndTime != -1 {
		time := time.Unix(json.EndTime, 0)
		ret.EndTime = &time
	}
	return
}
