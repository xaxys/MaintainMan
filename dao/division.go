package dao

import (
	"database/sql"
	"maintainman/database"
	"maintainman/logger"
	"maintainman/model"

	"gorm.io/gorm"
)

func GetDivisionByID(id uint) (*model.Division, error) {
	return TxGetDivisionByID(database.DB, id)
}

func TxGetDivisionByID(tx *gorm.DB, id uint) (*model.Division, error) {
	division := &model.Division{}
	if err := tx.Preload("Children").First(division, id).Error; err != nil {
		logger.Logger.Debugf("GetDivisionByIDErr: %v\n", err)
		return nil, err
	}
	return division, nil
}

func GetDivisionsByParentID(id uint) ([]*model.Division, error) {
	return TxGetDivisionsByParentID(database.DB, id)
}

func TxGetDivisionsByParentID(tx *gorm.DB, id uint) (divisions []*model.Division, err error) {
	if id != 0 {
		tx = tx.Where("parent_id = (?)", id)
	} else {
		tx = tx.Where("parent_id is null")
	}
	if err = tx.Find(&divisions).Error; err != nil {
		logger.Logger.Debugf("GetDivisionsByParentIDErr: %v\n", err)
	}
	return
}

func CreateDivision(aul *model.CreateDivisionRequest) (*model.Division, error) {
	return TxCreateDivision(database.DB, aul)
}

func TxCreateDivision(tx *gorm.DB, aul *model.CreateDivisionRequest) (*model.Division, error) {
	division := &model.Division{
		Name:     aul.Name,
		ParentID: sql.NullInt64{Int64: int64(aul.ParentID), Valid: aul.ParentID != 0},
	}
	if err := tx.Create(division).Error; err != nil {
		logger.Logger.Debugf("CreateDivisionErr: %v\n", err)
		return nil, err
	}
	return division, nil
}

func UpdateDivision(id uint, aul *model.UpdateDivisionRequest) (*model.Division, error) {
	return TxUpdateDivision(database.DB, id, aul)
}

// FIXME: Can't update division correctly
func TxUpdateDivision(tx *gorm.DB, id uint, aul *model.UpdateDivisionRequest) (*model.Division, error) {
	division := &model.Division{
		Name:     aul.Name,
		ParentID: sql.NullInt64{Int64: 0, Valid: false},
	}
	division.ID = id
	tx = tx.Model(division).Updates(division)
	if aul.ParentID != 0 {
		tx = tx.Update("parent_id", sql.NullInt64{Int64: int64(aul.ParentID), Valid: aul.ParentID != -1})
	}
	if err := tx.Error; err != nil {
		logger.Logger.Debugf("UpdateDivisionErr: %v\n", err)
		return nil, err
	}
	return division, nil
}

func DeleteDivision(id uint) error {
	return TxDeleteDivision(database.DB, id)
}

func TxDeleteDivision(tx *gorm.DB, id uint) (err error) {
	if err = tx.Delete(&model.Division{}, id).Error; err != nil {
		logger.Logger.Debugf("DeleteDivisionErr: %v\n", err)
	}
	return
}
