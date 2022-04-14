package order

import (
	"errors"

	"github.com/xaxys/maintainman/core/dao"
	"github.com/xaxys/maintainman/core/logger"
	"github.com/xaxys/maintainman/core/model"

	"gorm.io/gorm"
)

func dbGetCommentByID(id uint) (*Comment, error) {
	return txGetCommentByID(mctx.Database, id)
}

func txGetCommentByID(tx *gorm.DB, id uint) (*Comment, error) {
	comment := &Comment{}
	if err := tx.First(comment, id).Error; err != nil {
		logger.Logger.Debugf("GetCommentByIDErr: %v\n", err)
		return nil, err
	}
	return comment, nil
}

func dbGetCommentsByOrder(id uint, param *model.PageParam) (comments []*Comment, err error) {
	return txGetCommentsByOrder(mctx.Database, id, param)
}

func txGetCommentsByOrder(tx *gorm.DB, oid uint, param *model.PageParam) (comments []*Comment, err error) {
	comment := &Comment{OrderID: oid}
	if err = dao.TxPageFilter(tx, param).Where(comment).Find(&comments).Error; err != nil {
		logger.Logger.Debugf("GetCommentsByOrderErr: %v\n", err)
		return
	}
	return
}

func dbCreateComment(oid, uid uint, name string, aul *CreateCommentRequest) (comment *Comment, err error) {
	mctx.Database.Transaction(func(tx *gorm.DB) error {
		if comment, err = txCreateComment(tx, oid, uid, name, aul); err != nil {
			logger.Logger.Debugf("CreateCommentErr: %v\n", err)
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
		logger.Logger.Debugf("DeleteCommentErr: %v\n", err)
		return err
	}
	return nil
}
