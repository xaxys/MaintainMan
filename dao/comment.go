package dao

import (
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

func GetCommentsByOrder(id, offset uint) (comments []*model.Comment, err error) {
	comment := &model.Comment{
		OrderID: id,
	}
	if err = Filter("id desc", offset, 0).Where(comment).Find(&comments).Error; err != nil {
		logger.Logger.Debugf("GetCommentsByOrderErr: %v\n", err)
		return
	}
	return
}

func CreateComment(oid, uid uint, content string) (comment *model.Comment, err error) {
	comment = &model.Comment{
		OrderID: oid,
		UserID:  uid,
		Content: content,
		BaseModel: model.BaseModel{
			CreatedBy: uid,
			UpdatedBy: uid,
		},
	}
	if err = database.DB.Transaction(func(tx *gorm.DB) error {
		cmt := &model.Comment{OrderID: oid}
		if err := tx.Where(cmt).Order("id desc").First(cmt).Error; err != nil {
			return err
		}
		comment.SequenceNum = cmt.SequenceNum + 1
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
