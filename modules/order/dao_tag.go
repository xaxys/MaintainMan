package order

import (
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func dbGetTagByID(id uint) (*Tag, error) {
	return txGetTagByID(mctx.Database, id)
}

func txGetTagByID(tx *gorm.DB, id uint) (*Tag, error) {
	tag := &Tag{}
	if err := tx.First(tag, id).Error; err != nil {
		mctx.Logger.Warnf("GetTagByIDErr: %v\n", err)
		return nil, err
	}
	return tag, nil
}

func dbGetTagsByIDs(ids []uint) (tags []*Tag, err error) {
	return txGetTagsByIDs(mctx.Database, ids)
}

func txGetTagsByIDs(tx *gorm.DB, ids []uint) (tags []*Tag, err error) {
	for _, id := range ids {
		tag, err := txGetTagByID(tx, id)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

func dbGetAllTagSorts() ([]string, error) {
	return txGetAllTagSorts(mctx.Database)
}

func txGetAllTagSorts(tx *gorm.DB) (sorts []string, err error) {
	if err = tx.Model(&Tag{}).Distinct().Pluck("Sort", &sorts).Error; err != nil {
		mctx.Logger.Warnf("GetAllTagSortsErr: %v\n", err)
	}
	return
}

func dbGetAllTagsBySort(sort string) ([]*Tag, error) {
	return txGetAllTagsBySort(mctx.Database, sort)
}

func txGetAllTagsBySort(tx *gorm.DB, sort string) (tags []*Tag, err error) {
	tag := &Tag{
		Sort: sort,
	}
	if err = tx.Where(tag).Find(&tags).Error; err != nil {
		mctx.Logger.Warnf("GetAllTagsBySortErr: %v\n", err)
	}
	return
}

func dbCreateTag(aul *CreateTagRequest, operator uint) (*Tag, error) {
	return txCreateTag(mctx.Database, aul, operator)
}

func txCreateTag(tx *gorm.DB, aul *CreateTagRequest, operator uint) (tag *Tag, err error) {
	tag = jsonToTag(aul)
	tag.CreatedBy = operator
	cond := &Tag{
		Sort: tag.Sort,
		Name: tag.Name,
	}
	if err = tx.Where(cond).Attrs(tag).FirstOrCreate(tag).Error; err != nil {
		mctx.Logger.Warnf("CreateTagErr: %v\n", err)
	}
	return
}

func dbUpdateTag(id uint, aul *CreateTagRequest, operator uint) (*Tag, error) {
	return txUpdateTag(mctx.Database, id, aul, operator)
}

func txUpdateTag(tx *gorm.DB, id uint, aul *CreateTagRequest, operator uint) (tag *Tag, err error) {
	tag = jsonToTag(aul)
	tag.ID = id
	tag.UpdatedBy = operator
	if err = tx.Model(tag).Updates(tag).Error; err != nil {
		mctx.Logger.Warnf("UpdateTagErr: %v\n", err)
	}
	return
}

func dbDeleteTag(id uint) error {
	return txDeleteTag(mctx.Database, id)
}

func txDeleteTag(tx *gorm.DB, id uint) (err error) {
	if err = tx.Select(clause.Associations).Delete(&Tag{}, id).Error; err != nil {
		mctx.Logger.Warnf("DeleteTagErr: %v\n", err)
	}
	return
}

func dbCheckTagsCongener(tags []*Tag) error {
	count := map[string]uint{}
	min := map[string]uint{}
	for _, t := range tags {
		count[t.Sort]++
		min[t.Sort] = t.Congener
		if min[t.Sort] != 0 && count[t.Sort] > min[t.Sort] {
			return fmt.Errorf("[%s] 标签超过最大数量", t.Sort)
		}
	}
	return nil
}

func jsonToTag(aul *CreateTagRequest) *Tag {
	return &Tag{
		Name:     aul.Name,
		Sort:     aul.Sort,
		Level:    aul.Level,
		Congener: aul.Congener,
	}
}
