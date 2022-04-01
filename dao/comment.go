package dao

import (
	"errors"
	"maintainman/database"
	"maintainman/logger"
	"maintainman/model"

	"gorm.io/gorm"
)

func GetCommentByID(id uint) (*model.Comment, error) {
	comment := &model.Comment{}
	if err := database.DB.First(comment, id).Error; err != nil {
		logger.Logger.Debugf("GetCommentByIDErr: %v\n", err)
		return nil, err
	}
	return comment, nil
}

func GetCommentsByOrder(id uint, param *model.PageParam) (comments []*model.Comment, err error) {
	comment := &model.Comment{OrderID: id}
	if err = PageFilter(param).Where(comment).Find(&comments).Error; err != nil {
		logger.Logger.Debugf("GetCommentsByOrderErr: %v\n", err)
		return
	}
	return
}

func CreateComment(oid, uid uint, aul *model.CreateCommentRequest) (comment *model.Comment, err error) {
	seqNum := uint(0)
	cmt := &model.Comment{OrderID: oid}
	if err = database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where(cmt).Order("id desc").First(cmt).Error; err == nil {
			seqNum = cmt.SequenceNum
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		user := &model.User{}
		if err := tx.First(user, uid).Error; err != nil {
			return err
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
		if err := tx.Create(comment).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		logger.Logger.Debugf("CreateCommentErr: %v\n", err)
		return
	}
	return
}

func DeleteComment(id uint) error {
	if err := database.DB.Delete(&model.Comment{}, id).Error; err != nil {
		logger.Logger.Debugf("DeleteCommentErr: %v\n", err)
		return err
	}
	return nil
}
