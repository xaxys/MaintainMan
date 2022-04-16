package user

import (
	"errors"

	"github.com/xaxys/maintainman/core/model"
	"github.com/xaxys/maintainman/core/util"

	"gorm.io/gorm"
)

func getDivisionService(id uint, auth *model.AuthInfo) *model.ApiJson {
	division, err := dbGetDivisionByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrorNotFound(err)
		}
		return model.ErrorQueryDatabase(err)
	}
	return model.Success(divisionToJson(division), "获取成功")
}

func getDivisionsByParentIDService(id uint, auth *model.AuthInfo) *model.ApiJson {
	divisions, err := dbGetDivisionsByParentID(id)
	if err != nil {
		return model.ErrorQueryDatabase(err)
	}
	return model.Success(util.TransSlice(divisions, divisionToJson), "获取成功")
}

func createDivisionService(aul *CreateDivisionRequest, auth *model.AuthInfo) *model.ApiJson {
	if err := util.Validator.Struct(aul); err != nil {
		return model.ErrorValidation(err)
	}
	division, err := dbCreateDivision(aul)
	if err != nil {
		return model.ErrorInsertDatabase(err)
	}
	return model.SuccessCreate(divisionToJson(division), "创建成功")
}

func updateDivisionService(id uint, aul *UpdateDivisionRequest, auth *model.AuthInfo) *model.ApiJson {
	if err := util.Validator.Struct(aul); err != nil {
		return model.ErrorValidation(err)
	}
	division, err := dbUpdateDivision(id, aul)
	if err != nil {
		return model.ErrorUpdateDatabase(err)
	}
	return model.SuccessUpdate(divisionToJson(division), "更新成功")
}

func deleteDivisionService(id uint, auth *model.AuthInfo) *model.ApiJson {
	if err := dbDeleteDivision(id); err != nil {
		return model.ErrorDeleteDatabase(err)
	}
	return model.SuccessUpdate(nil, "删除成功")
}

func divisionToJson(division *Division) *DivisionJson {
	if division == nil {
		return nil
	} else {
		return &DivisionJson{
			ID:       division.ID,
			Name:     division.Name,
			ParentID: uint(division.ParentID.Int64),
			Children: util.TransSlice(division.Children, divisionToJson),
		}
	}
}
