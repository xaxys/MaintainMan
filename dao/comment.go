package dao

import (
	"errors"
	"maintainman/database"
	"maintainman/logger"
	"maintainman/model"

	"gorm.io/gorm"
)

func GetCommentByID(id uint) (*model.Comment, error) {
	return TxGetCommentByID(database.DB, id)
}

func TxGetCommentByID(tx *gorm.DB, id uint) (*model.Comment, error) {
	comment := &model.Comment{}
	if err := tx.First(comment, id).Error; err != nil {
		logger.Logger.Debugf("GetCommentByIDErr: %v\n", err)
		return nil, err
	}
	return comment, nil
}

func GetCommentsByOrder(id uint, param *model.PageParam) (comments []*model.Comment, err error) {
	return TxGetCommentsByOrder(database.DB, id, param)
}

func TxGetCommentsByOrder(tx *gorm.DB, oid uint, param *model.PageParam) (comments []*model.Comment, err error) {
	comment := &model.Comment{OrderID: oid}
	if err = TxPageFilter(tx, param).Where(comment).Find(&comments).Error; err != nil {
		logger.Logger.Debugf("GetCommentsByOrderErr: %v\n", err)
		return
	}
	return
}

func CreateComment(oid, uid uint, aul *model.CreateCommentRequest) (comment *model.Comment, err error) {
	database.DB.Transaction(func(tx *gorm.DB) error {
		if comment, err = TxCreateComment(tx, oid, uid, aul); err != nil {
			logger.Logger.Debugf("CreateCommentErr: %v\n", err)
		}
		return err
	})
	return
}

func TxCreateComment(tx *gorm.DB, oid, uid uint, aul *model.CreateCommentRequest) (comment *model.Comment, err error) {
	seqNum := uint(0)
	cmt := &model.Comment{OrderID: oid}
	if err = tx.Where(cmt).Order("id desc").First(cmt).Error; err == nil {
		seqNum = cmt.SequenceNum
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return
	}
	user := &model.User{}
	if err = tx.First(user, uid).Error; err != nil {
		return
	}
	comment = &model.Comment{
		OrderID:     oid,
		UserID:      uid,
		UserName:    user.Name,
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

func DeleteComment(id uint) error {
	return TxDeleteComment(database.DB, id)
}

func TxDeleteComment(tx *gorm.DB, id uint) error {
	if err := tx.Delete(&model.Comment{}, id).Error; err != nil {
		logger.Logger.Debugf("DeleteCommentErr: %v\n", err)
		return err
	}
	return nil
}
