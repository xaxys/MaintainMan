package service

import (
	"errors"

	"github.com/xaxys/maintainman/core/dao"
	"github.com/xaxys/maintainman/core/model"
	"github.com/xaxys/maintainman/core/util"

	"gorm.io/gorm"
)

func GetDivision(id uint, auth *model.AuthInfo) *model.ApiJson {
	division, err := dao.GetDivisionByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrorNotFound(err)
		}
		return model.ErrorQueryDatabase(err)
	}
	return model.Success(DivisionToJson(division), "获取成功")
}

func GetDivisionsByParentID(id uint, auth *model.AuthInfo) *model.ApiJson {
	divisions, err := dao.GetDivisionsByParentID(id)
	if err != nil {
		return model.ErrorQueryDatabase(err)
	}
	return model.Success(util.TransSlice(divisions, DivisionToJson), "获取成功")
}

func CreateDivision(aul *model.CreateDivisionRequest, auth *model.AuthInfo) *model.ApiJson {
	if err := util.Validator.Struct(aul); err != nil {
		return model.ErrorValidation(err)
	}
	division, err := dao.CreateDivision(aul)
	if err != nil {
		return model.ErrorInsertDatabase(err)
	}
	return model.SuccessCreate(DivisionToJson(division), "创建成功")
}

func UpdateDivision(id uint, aul *model.UpdateDivisionRequest, auth *model.AuthInfo) *model.ApiJson {
	if err := util.Validator.Struct(aul); err != nil {
		return model.ErrorValidation(err)
	}
	division, err := dao.UpdateDivision(id, aul)
	if err != nil {
		return model.ErrorUpdateDatabase(err)
	}
	return model.SuccessUpdate(DivisionToJson(division), "更新成功")
}

func DeleteDivision(id uint, auth *model.AuthInfo) *model.ApiJson {
	if err := dao.DeleteDivision(id); err != nil {
		return model.ErrorDeleteDatabase(err)
	}
	return model.SuccessUpdate(nil, "删除成功")
}

func DivisionToJson(division *model.Division) *model.DivisionJson {
	if division == nil {
		return nil
	} else {
		return &model.DivisionJson{
			ID:       division.ID,
			Name:     division.Name,
			ParentID: uint(division.ParentID.Int64),
			Children: util.TransSlice(division.Children, DivisionToJson),
		}
	}
}
