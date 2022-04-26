package wxnotify

import (
	"fmt"

	"github.com/xaxys/maintainman/core/module"
	"github.com/xaxys/maintainman/core/util"
	"github.com/xaxys/maintainman/modules/order"
	"github.com/xaxys/maintainman/modules/user"
)

var Module = module.Module{
	ModuleName:    "wxnotify",
	ModuleVersion: "1.0.0",
	ModuleEnv:     map[string]any{},
	ModuleExport:  map[string]any{},
	ModulePerm:    map[string]string{},
	EntryPoint:    entry,
}

var mctx *module.ModuleContext

func entry(ctx *module.ModuleContext) {
	mctx = ctx
	initAccessToken()
	go listener()
}

const sendMessageURL = "https://api.weixin.qq.com/cgi-bin/message/template/send"

type wxSendMessageResponse struct {
	ErrCode int64  `json:"errcode"` // 错误码
	ErrMsg  string `json:"errmsg"`  // 错误信息
}

func listener() {
	defer func() {
		if err := recover(); err != nil {
			mctx.Logger.Errorf("wxnotify listener panic: %s", err)
		}
	}()
	if getAccessToken() == "" {
		mctx.Logger.Infof("access token is empty, wechat notification service will be unavailable")
		return
	}

	orderModule := mctx.Registry.Get("order")
	if orderModule == nil {
		mctx.Logger.Errorf("order module not found")
		return
	}

	statusTmplID := getExportString(orderModule, "wechat.status.tmpl")
	if statusTmplID == "" {
		mctx.Logger.Errorf("status template id not found, wechat status notification service will be unavailable")
	}
	keyStatusOrder := getExportString(orderModule, "wechat.status.order")
	keyStatusTitle := getExportString(orderModule, "wechat.status.title")
	keyStatusStatus := getExportString(orderModule, "wechat.status.status")
	keyStatusTime := getExportString(orderModule, "wechat.status.time")
	keyStatusOther := getExportString(orderModule, "wechat.status.other")

	commentTmplID := getExportString(orderModule, "wechat.comment.tmpl")
	if commentTmplID == "" {
		mctx.Logger.Errorf("comment template id not found, wechat comment notification service will be unavailable")
	}
	keyCommentTitle := getExportString(orderModule, "wechat.comment.title")
	keyCommentName := getExportString(orderModule, "wechat.comment.name")
	keyCommentMessage := getExportString(orderModule, "wechat.comment.message")
	keyCommentTime := getExportString(orderModule, "wechat.comment.time")

	for {
		select {
		// order status changed notification
		case ch := <-mctx.EventBus.On("order:update:status:*"):
			if statusTmplID == "" {
				continue
			}
			orderID, _ := ch.Args[0].(uint)
			status, _ := ch.Args[1].(int)
			var odr *order.Order
			var err error
			if status == order.StatusAssigned {
				odr, err = order.GetOrderWithLastStatus(orderID)
			} else {
				odr, err = order.GetOrderByID(orderID)
			}
			if err != nil {
				mctx.Logger.Warnf("get order failed: %s", err)
				continue
			}
			usr, err := user.GetUserByID(odr.UserID)
			if err != nil {
				mctx.Logger.Warnf("get user failed: %s", err)
				continue
			}
			if usr.OpenID == "" {
				mctx.Logger.Infof("user %d has no openid, skipped", usr.ID)
				continue
			}

			// get template data
			data := map[string]string{}
			if keyStatusOrder != "" {
				data[keyStatusOrder] = fmt.Sprintf("%d", odr.ID)
			}
			if keyStatusTitle != "" {
				data[keyStatusTitle] = odr.Title
			}
			if keyStatusStatus != "" {
				data[keyStatusStatus] = order.StatusName(status)
			}
			if keyStatusTime != "" {
				data[keyStatusTime] = odr.UpdatedAt.Local().Format("2006-01-02 15:04:05")
			}
			if keyStatusOther != "" && status == order.StatusAssigned && odr.Status == uint(status) {
				// add repairer info if status is assigned
				repairerID, _ := ch.Args[2].(uint)
				repairer, err := user.GetUserByID(repairerID)
				if err != nil {
					mctx.Logger.Warnf("get repairer failed: %s", err)
					continue
				}
				data[keyStatusOther] = fmt.Sprintf("维修师傅 %s 将尽快为您维修", repairer.Name)
			}

			// send notification
			param := map[string]string{
				"access_token": getAccessToken(),
			}

			payload := map[string]any{
				"touser":      usr.OpenID,
				"template_id": statusTmplID,
				"data":        data,
			}

			wxResp, err := util.HTTPRequest[wxSendMessageResponse](sendMessageURL, "POST", param, payload)
			if err != nil {
				mctx.Logger.Warnf("send wechat message failed: %s", err)
				continue
			}
			if wxResp.ErrCode != 0 {
				mctx.Logger.Warnf("send wechat message failed: %s", wxResp.ErrMsg)
				continue
			}
		// order comment notification
		case ch := <-mctx.EventBus.On("order:update:comment"):
			if commentTmplID == "" {
				continue
			}
			orderID, _ := ch.Args[0].(uint)
			commentID, _ := ch.Args[1].(uint)
			comment, err := order.GetCommentByID(commentID)
			if err != nil {
				mctx.Logger.Warnf("get comment failed: %s", err)
				continue
			}
			odr, err := order.GetOrderWithLastStatus(orderID)
			if err != nil {
				mctx.Logger.Warnf("get order failed: %s", err)
				continue
			}

			// get template data
			data := map[string]string{}
			if keyCommentTitle != "" {
				data[keyCommentTitle] = odr.Title
			}
			if keyCommentName != "" {
				data[keyCommentName] = comment.UserName
			}
			if keyCommentMessage != "" {
				data[keyCommentMessage] = comment.Content
			}
			if keyCommentTime != "" {
				data[keyCommentTime] = comment.CreatedAt.Local().Format("2006-01-02 15:04:05")
			}

			openIDs := []string{}
			// send notification to user
			if odr.UserID != comment.UserID {
				usr, err := user.GetUserByID(odr.UserID)
				if err != nil {
					mctx.Logger.Warnf("get user failed: %s", err)
					continue
				}
				if usr.OpenID != "" {
					openIDs = append(openIDs, usr.OpenID)
				}
			}
			// send notification to current repairer
			if odr.Status == order.StatusAssigned {
				repairerID := util.LastElem(odr.StatusList).RepairerID
				if repairerID.Int64 == 0 || repairerID.Valid == false {
					mctx.Logger.Warnf("repairer id not found")
					continue
				}
				if uint(repairerID.Int64) == comment.UserID {
					continue
				}
				repairer, err := user.GetUserByID(uint(repairerID.Int64))
				if err != nil {
					mctx.Logger.Warnf("get repairer failed: %s", err)
					continue
				}
				if repairer.OpenID != "" {
					openIDs = append(openIDs, repairer.OpenID)
				}
			}

			// send notification
			for _, openID := range openIDs {
				param := map[string]string{
					"access_token": getAccessToken(),
				}
				payload := map[string]any{
					"touser":      openID,
					"template_id": statusTmplID,
					"data":        data,
				}

				wxResp, err := util.HTTPRequest[wxSendMessageResponse](sendMessageURL, "POST", param, payload)
				if err != nil {
					mctx.Logger.Warnf("send wechat message failed: %s", err)
					continue
				}
				if wxResp.ErrCode != 0 {
					mctx.Logger.Warnf("send wechat message failed: %s", wxResp.ErrMsg)
					continue
				}
			}
		}
	}
}

func getExportString(mod module.IModule, key string) string {
	exp, ok := mod.Export(key)
	if !ok {
		mctx.Logger.Warnf(fmt.Sprintf("export key %s not found", key))
		return ""
	}
	s, _ := exp.(string)
	return s
}
