package order

import (
	"errors"

	"github.com/xaxys/maintainman/core/dao"
	"github.com/xaxys/maintainman/core/model"

	"gorm.io/gorm"
)

func dbGetCommentCountByOrder(id uint) (uint, error) {
	return txGetCommentCountByOrder(mctx.Database, id)
}

func txGetCommentCountByOrder(tx *gorm.DB, id uint) (uint, error) {
	count := int64(0)
	comment := &Comment{OrderID: id}
	if err := tx.Model(comment).Where(comment).Count(&count).Error; err != nil {
		mctx.Logger.Debugf("GetCommentCountByOrderErr: %v\n", err)
		return 0, err
	}
	return uint(count), nil
}

func dbGetCommentByID(id uint) (*Comment, error) {
	return txGetCommentByID(mctx.Database, id)
}

func txGetCommentByID(tx *gorm.DB, id uint) (*Comment, error) {
	comment := &Comment{}
	if err := tx.First(comment, id).Error; err != nil {
		mctx.Logger.Debugf("GetCommentByIDErr: %v\n", err)
		return nil, err
	}
	return comment, nil
}

func dbGetCommentsByOrder(id uint, param *model.PageParam) (comments []*Comment, count uint, err error) {
	mctx.Database.Transaction(func(tx *gorm.DB) error {
		comments, count, err = txGetCommentsByOrder(tx, id, param)
		mctx.Logger.Debugf("GetCommentsByOrder: %v\n", err)
		return err
	})
	return
}

func txGetCommentsByOrder(tx *gorm.DB, oid uint, param *model.PageParam) (comments []*Comment, count uint, err error) {
	comment := &Comment{OrderID: oid}
	tx = dao.TxPageFilter(tx, param).Where(comment)
	if err = tx.Find(&comments).Error; err != nil {
		return
	}
	cnt := int64(0)
	if err = tx.Count(&cnt).Error; err != nil {
		return
	}
	count = uint(cnt)
	return
}

func dbCreateComment(oid, uid uint, name string, aul *CreateCommentRequest) (comment *Comment, err error) {
	mctx.Database.Transaction(func(tx *gorm.DB) error {
		if comment, err = txCreateComment(tx, oid, uid, name, aul); err != nil {
			mctx.Logger.Debugf("CreateCommentErr: %v\n", err)
		}
		return err
	})
	return
}

func txCreateComment(tx *gorm.DB, oid, uid uint, name string, aul *CreateCommentRequest) (comment *Comment, err error) {
	seqNum := uint(0)
	cmt := &Comment{OrderID: oid}
	if err = tx.Where(cmt).Order("id desc").First(cmt).Error; err == nil {
		seqNum = cmt.SequenceNum
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return
	}
	comment = &Comment{
		OrderID:     oid,
		UserID:      uid,
		UserName:    name,
		Content:     aul.Content,
		SequenceNum: seqNum + 1,
		BaseModel: model.BaseModel{
			CreatedBy: uid,
			UpdatedBy: uid,
		},
	}
	if err = tx.Create(comment).Error; err != nil {
		return
	}
	return
}

func dbDeleteComment(id uint) error {
	return txDeleteComment(mctx.Database, id)
}

func txDeleteComment(tx *gorm.DB, id uint) error {
	if err := tx.Delete(&Comment{}, id).Error; err != nil {
		mctx.Logger.Debugf("DeleteCommentErr: %v\n", err)
		return err
	}
	return nil
}
