package service

import (
	"errors"
	"maintainman/dao"
	"maintainman/model"
	"maintainman/util"

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
			ParentID: division.ParentID,
			Children: util.TransSlice(division.Children, DivisionToJson),
		}
	}
}
