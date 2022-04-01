package dao

import (
	"maintainman/database"
	"maintainman/logger"
	"maintainman/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func GetTagByID(id uint) (*model.Tag, error) {
	return TxGetTagByID(database.DB, id)
}

func TxGetTagByID(tx *gorm.DB, id uint) (*model.Tag, error) {
	tag := &model.Tag{}
	if err := tx.First(tag, id).Error; err != nil {
		logger.Logger.Debugf("GetTagByIDErr: %v\n", err)
		return nil, err
	}
	return tag, nil
}

func GetTagsByIDs(ids []uint) (tags []*model.Tag, err error) {
	return TxGetTagsByIDs(database.DB, ids)
}

func TxGetTagsByIDs(tx *gorm.DB, ids []uint) (tags []*model.Tag, err error) {
	for _, id := range ids {
		tag, err := TxGetTagByID(tx, id)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

func GetAllTagSorts() ([]string, error) {
	return TxGetAllTagSorts(database.DB)
}

func TxGetAllTagSorts(tx *gorm.DB) (sorts []string, err error) {
	if err = tx.Model(&model.Tag{}).Distinct().Pluck("Sort", &sorts).Error; err != nil {
		logger.Logger.Debugf("GetAllTagSortsErr: %v\n", err)
	}
	return
}

func GetAllTagsBySort(sort string) ([]*model.Tag, error) {
	return TxGetAllTagsBySort(database.DB, sort)
}

func TxGetAllTagsBySort(tx *gorm.DB, sort string) (tags []*model.Tag, err error) {
	tag := &model.Tag{
		Sort: sort,
	}
	if err = tx.Where(tag).Find(&tags).Error; err != nil {
		logger.Logger.Debugf("GetAllTagsBySortErr: %v\n", err)
	}
	return
}

func CreateTag(aul *model.CreateTagRequest, operator uint) (*model.Tag, error) {
	return TxCreateTag(database.DB, aul, operator)
}

func TxCreateTag(tx *gorm.DB, aul *model.CreateTagRequest, operator uint) (tag *model.Tag, err error) {
	tag = JsonToTag(aul)
	tag.CreatedBy = operator
	cond := &model.Tag{
		Sort: tag.Sort,
		Name: tag.Name,
	}
	if err = tx.Where(cond).Attrs(tag).FirstOrCreate(tag).Error; err != nil {
		logger.Logger.Debugf("CreateTagErr: %v\n", err)
	}
	return
}

func UpdateTag(id uint, aul *model.CreateTagRequest, operator uint) (*model.Tag, error) {
	return TxUpdateTag(database.DB, id, aul, operator)
}

func TxUpdateTag(tx *gorm.DB, id uint, aul *model.CreateTagRequest, operator uint) (tag *model.Tag, err error) {
	tag = JsonToTag(aul)
	tag.ID = id
	tag.UpdatedBy = operator
	if err = tx.Where(tag).Updates(tag).Error; err != nil {
		logger.Logger.Debugf("UpdateTagErr: %v\n", err)
	}
	return
}

func DeleteTag(id uint) error {
	return TxDeleteTag(database.DB, id)
}

func TxDeleteTag(tx *gorm.DB, id uint) (err error) {
	if err = tx.Select(clause.Associations).Delete(&model.Tag{}, id).Error; err != nil {
		logger.Logger.Debugf("DeleteTagErr: %v\n", err)
	}
	return
}

func JsonToTag(aul *model.CreateTagRequest) *model.Tag {
	return &model.Tag{
		Name:  aul.Name,
		Sort:  aul.Sort,
		Level: aul.Level,
	}
}
