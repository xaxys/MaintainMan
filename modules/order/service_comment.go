package order

import (
	"fmt"

	"github.com/xaxys/maintainman/core/model"
	"github.com/xaxys/maintainman/core/util"
)

func getCommentsByOrderService(id uint, param *model.PageParam, auth *model.AuthInfo) *model.ApiJson {
	order, err := dbGetOrderWithLastStatus(id)
	if err != nil {
		return model.ErrorNotFound(err)
	}
	if order.UserID != auth.User && uint(util.LastElem(order.StatusList).RepairerID.Int64) != auth.User {
		return model.ErrorNoPermissions(fmt.Errorf("您不是订单的创建者或指派人，不能查看评论"))
	}
	return forceGetCommentsByOrderService(id, param, auth)
}

func forceGetCommentsByOrderService(id uint, param *model.PageParam, auth *model.AuthInfo) *model.ApiJson {
	param.OrderBy = util.NotEmpty(param.OrderBy, "id desc")
	comments, err := dbGetCommentsByOrder(id, param)
	if err != nil {
		return model.ErrorQueryDatabase(err)
	}
	cs := util.TransSlice(comments, commentToJson)
	return model.Success(cs, "获取成功")
}

func createCommentService(id uint, aul *CreateCommentRequest, auth *model.AuthInfo) *model.ApiJson {
	order, err := dbGetOrderWithLastStatus(id)
	if err != nil {
		return model.ErrorNotFound(err)
	}
	if order.UserID != auth.User && uint(util.LastElem(order.StatusList).RepairerID.Int64) != auth.User {
		return model.ErrorNoPermissions(fmt.Errorf("您不是订单的创建者或指派人，不能创建评论"))
	}
	if order.AllowComment == CommentDisallow {
		return model.ErrorNoPermissions(fmt.Errorf("该订单不允许评论"))
	}
	return forceCreateCommentService(id, aul, auth)
}

func forceCreateCommentService(id uint, aul *CreateCommentRequest, auth *model.AuthInfo) *model.ApiJson {
	if err := util.Validator.Struct(aul); err != nil {
		return model.ErrorValidation(err)
	}
	comment, err := dbCreateComment(id, auth.User, auth.Name, aul)
	if err != nil {
		return model.ErrorInsertDatabase(err)
	}
	return model.SuccessCreate(commentToJson(comment), "创建成功")
}

func DeleteCommentService(id uint, auth *model.AuthInfo) *model.ApiJson {
	comment, err := dbGetCommentByID(id)
	if err != nil {
		return model.ErrorNotFound(err)
	}
	if comment.UserID != auth.User {
		return model.ErrorNoPermissions(fmt.Errorf("操作人不是评论创建者"))
	}
	return forceDeleteCommentService(id, auth)
}

func forceDeleteCommentService(id uint, auth *model.AuthInfo) *model.ApiJson {
	err := dbDeleteComment(id)
	if err != nil {
		return model.ErrorDeleteDatabase(err)
	}
	return model.SuccessUpdate(nil, "删除成功")
}

func commentToJson(comment *Comment) *CommentJson {
	if comment == nil {
		return nil
	} else {
		return &CommentJson{
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
