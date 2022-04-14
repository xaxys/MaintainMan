package order

import (
	"errors"
	"fmt"

	"github.com/xaxys/maintainman/core/dao"
	"github.com/xaxys/maintainman/core/model"
	"github.com/xaxys/maintainman/core/util"

	"gorm.io/gorm"
)

func getTagByIDService(id uint, auth *model.AuthInfo) *model.ApiJson {
	tag, err := dbGetTagByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrorNotFound(err)
		}
		return model.ErrorQueryDatabase(err)
	}
	role := util.NilOrBaseValue(auth, func(v *model.AuthInfo) string { return v.Role }, "")
	if err := dao.CheckPermission(role, fmt.Sprintf("tag.view.%d", tag.Level)); err != nil {
		return model.ErrorNoPermissions(err)
	}
	return model.Success(tagToJson(tag), "获取成功")
}

func getAllTagSortsService(auth *model.AuthInfo) *model.ApiJson {
	tags, err := dbGetAllTagSorts()
	if err != nil {
		return model.ErrorQueryDatabase(err)
	}
	return model.Success(tags, "获取成功")
}

func getAllTagsBySortService(sort string, auth *model.AuthInfo) *model.ApiJson {
	tags, err := dbGetAllTagsBySort(sort)
	if err != nil {
		return model.ErrorQueryDatabase(err)
	}
	role := util.NilOrBaseValue(auth, func(v *model.AuthInfo) string { return v.Role }, "")
	ts := util.TransSlice(tags, func(t *Tag) *TagJson {
		if dao.HasPermission(role, fmt.Sprintf("tag.view.%d", t.Level)) {
			return tagToJson(t)
		}
		return nil
	})
	return model.Success(ts, "获取成功")
}

func createTagService(aul *CreateTagRequest, auth *model.AuthInfo) *model.ApiJson {
	tag, err := dbCreateTag(aul, auth.User)
	if err != nil {
		return model.ErrorInsertDatabase(err)
	}
	return model.SuccessCreate(tagToJson(tag), "创建成功")
}

func deleteTagService(id uint, auth *model.AuthInfo) *model.ApiJson {
	err := dbDeleteTag(id)
	if err != nil {
		return model.ErrorDeleteDatabase(err)
	}
	return model.SuccessUpdate(nil, "删除成功")
}

func checkTagsService(tagIDs []uint, perm, role string) *model.ApiJson {
	tags, err := dbGetTagsByIDs(tagIDs)
	if err != nil {
		return model.ErrorQueryDatabase(err)
	}
	for _, t := range tags {
		if err := dao.CheckPermission(role, fmt.Sprintf("%s.%d", perm, t.Level)); err != nil {
			return model.ErrorNoPermissions(err)
		}
	}
	return nil
}

func tagToJson(tag *Tag) *TagJson {
	if tag == nil {
		return nil
	} else {
		return &TagJson{
			ID:       tag.ID,
			Sort:     tag.Sort,
			Name:     tag.Name,
			Level:    tag.Level,
			Congener: tag.Congener,
		}
	}
}
