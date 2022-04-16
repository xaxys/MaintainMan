package user

import (
	"database/sql"

	"gorm.io/gorm"
)

func dbGetDivisionByID(id uint) (*Division, error) {
	return txGetDivisionByID(mctx.Database, id)
}

func txGetDivisionByID(tx *gorm.DB, id uint) (*Division, error) {
	division := &Division{}
	if err := tx.Preload("Children").First(division, id).Error; err != nil {
		mctx.Logger.Debugf("GetDivisionByIDErr: %v\n", err)
		return nil, err
	}
	return division, nil
}

func dbGetDivisionsByParentID(id uint) ([]*Division, error) {
	return txGetDivisionsByParentID(mctx.Database, id)
}

func txGetDivisionsByParentID(tx *gorm.DB, id uint) (divisions []*Division, err error) {
	if id != 0 {
		tx = tx.Where("parent_id = (?)", id)
	} else {
		tx = tx.Where("parent_id is null")
	}
	if err = tx.Find(&divisions).Error; err != nil {
		mctx.Logger.Debugf("GetDivisionsByParentIDErr: %v\n", err)
	}
	return
}

func dbCreateDivision(aul *CreateDivisionRequest) (*Division, error) {
	return txCreateDivision(mctx.Database, aul)
}

func txCreateDivision(tx *gorm.DB, aul *CreateDivisionRequest) (*Division, error) {
	division := &Division{
		Name:     aul.Name,
		ParentID: sql.NullInt64{Int64: int64(aul.ParentID), Valid: aul.ParentID != 0},
	}
	if err := tx.Create(division).Error; err != nil {
		mctx.Logger.Debugf("CreateDivisionErr: %v\n", err)
		return nil, err
	}
	return division, nil
}

func dbUpdateDivision(id uint, aul *UpdateDivisionRequest) (*Division, error) {
	return txUpdateDivision(mctx.Database, id, aul)
}

func txUpdateDivision(tx *gorm.DB, id uint, aul *UpdateDivisionRequest) (*Division, error) {
	division := &Division{
		Name:     aul.Name,
		ParentID: sql.NullInt64{Int64: 0, Valid: false},
	}
	division.ID = id
	tx = tx.Model(division).Updates(division)
	if aul.ParentID != 0 {
		tx = tx.Update("parent_id", sql.NullInt64{Int64: int64(aul.ParentID), Valid: aul.ParentID != -1})
	}
	if err := tx.Error; err != nil {
		mctx.Logger.Debugf("UpdateDivisionErr: %v\n", err)
		return nil, err
	}
	return division, nil
}

func dbDeleteDivision(id uint) error {
	return txDeleteDivision(mctx.Database, id)
}

func txDeleteDivision(tx *gorm.DB, id uint) (err error) {
	if err = tx.Delete(&Division{}, id).Error; err != nil {
		mctx.Logger.Debugf("DeleteDivisionErr: %v\n", err)
	}
	return
}
