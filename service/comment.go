package service

import (
	"fmt"
	"maintainman/dao"
	"maintainman/model"
	"maintainman/util"
)

func GetCommentsByOrder(id uint, param *model.PageParam, auth *model.AuthInfo) *model.ApiJson {
	order, err := dao.GetOrderWithLastStatus(id)
	if err != nil {
		return model.ErrorNotFound(err)
	}
	if order.UserID != auth.User && util.LastElem(order.StatusList).RepairerID != auth.User {
		return model.ErrorNoPermissions(fmt.Errorf("您不是订单的创建者或指派人，不能查看评论"))
	}
	return ForceGetCommentsByOrder(id, param, auth)
}

func ForceGetCommentsByOrder(id uint, param *model.PageParam, auth *model.AuthInfo) *model.ApiJson {
	param.OrderBy = util.NotEmpty(param.OrderBy, "id desc")
	comments, err := dao.GetCommentsByOrder(id, param)
	if err != nil {
		return model.ErrorQueryDatabase(err)
	}
	cs := util.TransSlice(comments, CommentToJson)
	return model.Success(cs, "获取成功")
}

func CreateComment(id uint, aul *model.CreateCommentRequest, auth *model.AuthInfo) *model.ApiJson {
	order, err := dao.GetOrderWithLastStatus(id)
	if err != nil {
		return model.ErrorNotFound(err)
	}
	if order.UserID != auth.User && util.LastElem(order.StatusList).RepairerID != auth.User {
		return model.ErrorNoPermissions(fmt.Errorf("您不是订单的创建者或指派人，不能创建评论"))
	}
	if order.AllowComment == model.CommentDisallow {
		return model.ErrorNoPermissions(fmt.Errorf("该订单不允许评论"))
	}
	return ForceCreateComment(id, aul, auth)
}

func ForceCreateComment(id uint, aul *model.CreateCommentRequest, auth *model.AuthInfo) *model.ApiJson {
	if err := util.Validator.Struct(aul); err != nil {
		return model.ErrorValidation(err)
	}
	comment, err := dao.CreateComment(id, auth.User, auth.Name, aul)
	if err != nil {
		return model.ErrorInsertDatabase(err)
	}
	return model.SuccessCreate(CommentToJson(comment), "创建成功")
}

func DeleteComment(id uint, auth *model.AuthInfo) *model.ApiJson {
	comment, err := dao.GetCommentByID(id)
	if err != nil {
		return model.ErrorNotFound(err)
	}
	if comment.UserID != auth.User {
		return model.ErrorNoPermissions(fmt.Errorf("操作人不是评论创建者"))
	}
	return ForceDeleteComment(id, auth)
}

func ForceDeleteComment(id uint, auth *model.AuthInfo) *model.ApiJson {
	err := dao.DeleteComment(id)
	if err != nil {
		return model.ErrorDeleteDatabase(err)
	}
	return model.SuccessUpdate(nil, "删除成功")
}

func CommentToJson(comment *model.Comment) *model.CommentJson {
	if comment == nil {
		return nil
	} else {
		return &model.CommentJson{
			ID:          comment.ID,
			OrderID:     comment.OrderID,
			UserID:      comment.UserID,
			UserName:    comment.UserName,
			SequenceNum: comment.SequenceNum,
			Content:     comment.Content,
			CreatedAt:   comment.CreatedAt.Unix(),
		}
	}
}
