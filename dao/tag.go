package dao

import (
	"fmt"
	"maintainman/database"
	"maintainman/logger"
	"maintainman/model"

	"gorm.io/gorm/clause"
)

func GetTagByID(id uint) (*model.Tag, error) {
	tag := &model.Tag{}

	if err := database.DB.First(tag, id).Error; err != nil {
		logger.Logger.Debugf("GetTagByIDErr: %v\n", err)
		return nil, err
	}

	return tag, nil
}

func GetTagsByIDs(ids []uint) ([]*model.Tag, error) {
	tags := []*model.Tag{}
	errs := []error{}
	for _, id := range ids {
		tag, err := GetTagByID(id)
		if err != nil {
			errs = append(errs, err)
		} else {
			tags = append(tags, tag)
		}
	}
	if len(errs) > 0 {
		return tags, fmt.Errorf("%v", errs)
	}
	return tags, nil
}

func GetAllTagSorts() (sorts []string, err error) {
	if err = database.DB.Model(&model.Tag{}).Distinct().Pluck("sort", &sorts).Error; err != nil {
		logger.Logger.Debugf("GetAllTagSortsErr: %v\n", err)
	}
	return
}

func GetAllTagsBySort(sort string) (tags []*model.Tag, err error) {
	tag := &model.Tag{
		Sort: sort,
	}
	if err = database.DB.Where(tag).Find(&tags).Error; err != nil {
		logger.Logger.Debugf("GetAllTagsBySortErr: %v\n", err)
	}
	return
}

func CreateTag(aul *model.CreateTagJson) (tag *model.Tag, err error) {
	tag = JsonToTag(aul)
	if err = database.DB.Create(tag).Error; err != nil {
		logger.Logger.Debugf("CreateTagErr: %v\n", err)
	}
	return
}

func UpdateTag(id uint, aul *model.CreateTagJson) (tag *model.Tag, err error) {
	tag = JsonToTag(aul)
	tag.ID = id
	if err = database.DB.Model(tag).Updates(tag).Error; err != nil {
		logger.Logger.Debugf("UpdateTagErr: %v\n", err)
	}
	return
}

func DeleteTag(id uint) (err error) {
	if err = database.DB.Select(clause.Associations).Delete(&model.Tag{}).Error; err != nil {
		logger.Logger.Debugf("DeleteTagErr: %v\n", err)
	}
	return
}

func JsonToTag(aul *model.CreateTagJson) *model.Tag {
	return &model.Tag{
		Name:  aul.Name,
		Sort:  aul.Sort,
		Level: aul.Level,
	}
}