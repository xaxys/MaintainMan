package announce

import (
	"time"

	"github.com/xaxys/maintainman/core/dao"

	"gorm.io/gorm"
)

func dbGetAnnounceCount() (count uint, err error) {
	return txGetAnnounceCount(mctx.Database)
}

func txGetAnnounceCount(tx *gorm.DB) (uint, error) {
	count := int64(0)
	if err := tx.Model(&Announce{}).Count(&count).Error; err != nil {
		mctx.Logger.Warnf("GetAnnounceCountErr: %v\n", err)
		return 0, err
	}
	return uint(count), nil
}

func dbGetAnnounceByID(id uint) (announce *Announce, err error) {
	return txGetAnnounceByID(mctx.Database, id)
}

func txGetAnnounceByID(tx *gorm.DB, id uint) (*Announce, error) {
	announce := &Announce{}
	if err := tx.First(announce, id).Error; err != nil {
		mctx.Logger.Warnf("GetAnnounceByIDErr: %v\n", err)
		return nil, err
	}
	return announce, nil
}

func dbGetAnnounceByTitle(title string) (announce *Announce, err error) {
	return txGetAnnounceByTitle(mctx.Database, title)
}

func txGetAnnounceByTitle(tx *gorm.DB, title string) (*Announce, error) {
	announce := &Announce{Title: title}
	if err := tx.Where(announce).First(announce).Error; err != nil {
		mctx.Logger.Warnf("GetAnnounceByTitleErr: %v\n", err)
		return nil, err
	}
	return announce, nil
}

func dbGetAllAnnouncesWithParam(aul *AllAnnounceRequest) (announces []*Announce, count uint, err error) {
	mctx.Database.Transaction(func(tx *gorm.DB) error {
		if announces, count, err = txGetAllAnnouncesWithParam(tx, aul); err != nil {
			mctx.Logger.Warnf("GetAllAnnouncesWithParamErr: %v\n", err)
		}
		return err
	})
	return
}

func txGetAllAnnouncesWithParam(tx *gorm.DB, aul *AllAnnounceRequest) (announces []*Announce, count uint, err error) {
	tx = dao.TxFilter(tx, aul.OrderBy, aul.Offset, aul.Limit)
	if aul.Title != "" {
		tx = tx.Where("title like ?", aul.Title)
	}

	if aul.StartTime != -1 {
		unix := time.Unix(aul.StartTime, 0)
		if aul.Inclusive {
			tx = tx.Where("start_time >= ?", unix)
		} else {
			tx = tx.Where("start_time <= ?", unix)
		}
	}

	if aul.EndTime != -1 {
		unix := time.Unix(aul.EndTime, 0)
		if aul.Inclusive {
			tx = tx.Where("end_time <= ?", unix)
		} else {
			tx = tx.Where("end_time >= ?", unix)
		}
	}

	if err = tx.Find(&announces).Error; err != nil {
		return
	}
	cnt := int64(0)
	if err = tx.Count(&cnt).Error; err != nil {
		return
	}
	count = uint(cnt)
	return
}

func dbCreateAnnounce(json *ModifyAnnounceRequest, operator uint) (*Announce, error) {
	return txCreateAnnounce(mctx.Database, json, operator)
}

func txCreateAnnounce(tx *gorm.DB, json *ModifyAnnounceRequest, operator uint) (*Announce, error) {
	announce := jsonToAnnounce(json)
	if announce.StartTime == nil {
		now := time.Now()
		announce.StartTime = &now
	}
	if announce.EndTime == nil {
		now := time.Unix(253370764799, 0)
		announce.EndTime = &now
	}
	announce.CreatedBy = operator

	if err := tx.Create(announce).Error; err != nil {
		mctx.Logger.Warnf("CreateAnnounceErr: %v\n", err)
		return nil, err
	}
	return announce, nil
}

func dbUpdateAnnounce(id uint, json *ModifyAnnounceRequest, operator uint) (*Announce, error) {
	return txUpdateAnnounce(mctx.Database, id, json, operator)
}

func txUpdateAnnounce(tx *gorm.DB, id uint, json *ModifyAnnounceRequest, operator uint) (*Announce, error) {
	announce := jsonToAnnounce(json)
	announce.ID = id
	announce.UpdatedBy = operator

	if err := tx.Model(announce).Updates(announce).Error; err != nil {
		mctx.Logger.Warnf("UpdateAnnounceErr: %v\n", err)
		return nil, err
	}
	return announce, nil
}

func dbDeleteAnnounce(id uint) error {
	return txDeleteAnnounce(mctx.Database, id)
}

func txDeleteAnnounce(tx *gorm.DB, id uint) error {
	if err := tx.Delete(&Announce{}, id).Error; err != nil {
		mctx.Logger.Warnf("DeleteAnnounceErr: %v\n", err)
		return err
	}
	return nil
}

func dbHitAnnounce(id uint) error {
	return txHitAnnounce(mctx.Database, id)
}

func txHitAnnounce(tx *gorm.DB, id uint) error {
	announce := &Announce{}
	announce.ID = id
	if err := tx.Model(announce).Update("hits", gorm.Expr("hits + ?", 1)).Error; err != nil {
		mctx.Logger.Warnf("HitAnnounceErr: %v\n", err)
		return err
	}
	return nil
}

func jsonToAnnounce(json *ModifyAnnounceRequest) (ret *Announce) {
	ret = &Announce{
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
