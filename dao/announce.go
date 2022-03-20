package dao

import (
	"fmt"
	"maintainman/database"
	"maintainman/model"
	. "maintainman/model"
	"time"

	"gorm.io/gorm"
)

func GetAnnounceByID(id uint) (*Announce, error) {
	announce := &Announce{}

	if err := database.DB.First(announce, id).Error; err != nil {
		fmt.Printf("GetAnnounceByIDErr: %v\n", err)
		return nil, err
	}

	return announce, nil
}

func GetAnnounceByTitle(title string) (*Announce, error) {
	announce := &Announce{Title: title}

	if err := database.DB.Where(announce).First(announce).Error; err != nil {
		fmt.Printf("GetAnnounceByNameErr: %v\n", err)
		return nil, err
	}

	return announce, nil
}

func GetAllAnnouncesWithParam(aul *model.AllAnnounceJson) (announces []*Announce, err error) {
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
	dry := db.Session(&gorm.Session{DryRun: true}).Find(&announces).Statement
	fmt.Println(dry.SQL.String())
	fmt.Println(dry.Vars)
	if err = db.Find(&announces).Error; err != nil {
		fmt.Printf("GetAllAnnounceErr: %v\n", err)
	}
	return
}

func CreateAnnounce(json *ModifyAnnounceJson) (*Announce, error) {
	announce := JsonToAnnounce(json)
	if announce.StartTime == nil {
		now := time.Now()
		announce.StartTime = &now
	}
	if announce.EndTime == nil {
		now := time.Unix(253370764799, 0)
		announce.EndTime = &now
	}
	announce.CreatedBy = json.OperatorID

	if err := database.DB.Create(announce).Error; err != nil {
		fmt.Printf("CreateAnnounceErr: %v\n", err)
		return nil, err
	}

	return announce, nil
}

func UpdateAnnounce(id uint, json *ModifyAnnounceJson) (*Announce, error) {
	announce := JsonToAnnounce(json)
	announce.ID = id
	announce.UpdatedBy = json.OperatorID

	if err := database.DB.Model(announce).Updates(announce).Error; err != nil {
		fmt.Printf("UpdateAnnounceErr: %v\n", err)
		return nil, err
	}

	return announce, nil
}

func DeleteAnnounce(id uint) error {
	if err := database.DB.Delete(&Announce{}, id).Error; err != nil {
		fmt.Printf("DeleteAnnounceErr: %v\n", err)
		return err
	}
	return nil
}

func JsonToAnnounce(json *ModifyAnnounceJson) (ret *Announce) {
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
