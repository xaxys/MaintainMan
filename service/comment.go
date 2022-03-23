package service

import (
	"fmt"
	"maintainman/dao"
	"maintainman/model"
	"maintainman/util"
)

func GetCommentsByOrder(id, offset, operator uint) *model.ApiJson {
	order, err := dao.GetOrderWithLastStatus(id)
	if err != nil {
		return model.ErrorNotFound(err)
	}
	if order.UserID != operator && order.StatusList[len(order.StatusList)-1].RepairerID != operator {
		return model.ErrorNoPermissions(fmt.Errorf("您不是订单的创建者或指派人，不能查看评论"))
	}
	return GetCommentsByOrderID(id, offset)
}

func GetCommentsByOrderID(id, offset uint) *model.ApiJson {
	comments, err := dao.GetCommentsByOrder(id, offset)
	if err != nil {
		return model.ErrorQueryDatabase(err)
	}
	cs := util.TransSlice(comments, CommentToJson)
	return model.Success(cs, "获取成功")
}

func CreateComment(oid uint, aul *model.CreateCommentJson) *model.ApiJson {
	order, err := dao.GetOrderWithLastStatus(oid)
	if err != nil {
		return model.ErrorNotFound(err)
	}
	if order.UserID != aul.OperatorID && order.StatusList[len(order.StatusList)-1].RepairerID != aul.OperatorID {
		return model.ErrorNoPermissions(fmt.Errorf("您不是订单的创建者或指派人，不能查看评论"))
	}
	return CreateCommentOverride(oid, aul)
}

func CreateCommentOverride(oid uint, aul *model.CreateCommentJson) *model.ApiJson {
	if err := util.Validator.Struct(aul); err != nil {
		return model.ErrorValidation(err)
	}
	comment, err := dao.CreateComment(oid, aul.OperatorID, aul.Content)
	if err != nil {
		return model.ErrorInsertDatabase(err)
	}
	return model.SuccessCreate(CommentToJson(comment), "创建成功")
}

func DeleteComment(id, operator uint) *model.ApiJson {
	comment, err := dao.GetCommentByID(id)
	if err != nil {
		return model.ErrorNotFound(err)
	}
	if comment.UserID != operator {
		return model.ErrorUpdateDatabase(fmt.Errorf("操作人不是评论创建者"))
	}
	return DeleteCommentByID(id)
}

func DeleteCommentByID(id uint) *model.ApiJson {
	err := dao.DeleteComment(id)
	if err != nil {
		return model.ErrorQueryDatabase(err)
	}
	return model.SuccessUpdate(nil, "删除成功")
}

func CommentToJson(comment *model.Comment) *model.CommentJson {
	return &model.CommentJson{
		ID:        comment.ID,
		OrderID:   comment.OrderID,
		UserID:    comment.UserID,
		UserName:  comment.UserName,
		Content:   comment.Content,
		CreatedAt: comment.CreatedAt.Unix(),
	}
}
