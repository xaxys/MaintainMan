package service

import (
	"errors"
	"maintainman/dao"
	"maintainman/model"
	"maintainman/util"

	"gorm.io/gorm"
)

func GetTagByID(id uint) *model.ApiJson {
	tag, err := dao.GetTagByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrorNotFound(err)
		} else {
			return model.ErrorQueryDatabase(err)
		}
	}
	return model.Success(TagToJson(tag), "获取成功")
}

func GetAllTagSorts() *model.ApiJson {
	tags, err := dao.GetAllTagSorts()
	if err != nil {
		return model.ErrorQueryDatabase(err)
	}
	return model.Success(tags, "获取成功")
}

func GetAllTagsBySort(sort string) *model.ApiJson {
	tags, err := dao.GetAllTagsBySort(sort)
	if err != nil {
		return model.ErrorQueryDatabase(err)
	}
	ts := util.TransSlice(tags, TagToJson)
	return model.Success(ts, "获取成功")
}

func CreateTag(aul *model.CreateTagJson) *model.ApiJson {
	tag, err := dao.CreateTag(aul)
	if err != nil {
		return model.ErrorInsertDatabase(err)
	}
	return model.SuccessCreate(TagToJson(tag), "创建成功")
}

// TODO: Add func UpdateTag ?

func DeleteTag(id uint) *model.ApiJson {
	err := dao.DeleteTag(id)
	if err != nil {
		return model.ErrorDeleteDatabase(err)
	}
	return model.SuccessUpdate(nil, "删除成功")
}

func TagToJson(tag *model.Tag) *model.TagJson {
	return util.NotNil(tag, &model.TagJson{
		ID:    tag.ID,
		Sort:  tag.Sort,
		Name:  tag.Name,
		Level: tag.Level,
	})
}
