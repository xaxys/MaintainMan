package wxnotify

import "github.com/xaxys/maintainman/core/util"

const accessTokenURL = "https://api.weixin.qq.com/cgi-bin/token"

var accessToken util.AtomPtr[string]

func getAccessToken() string {
	return util.NilOrBaseValue(accessToken.Get(), func(v *string) string { return *v }, "")
}

func setAccessToken(token string) {
	accessToken.Set(&token)
}

type wxAccessTokenResponse struct {
	AccessToken string `json:"access_token"` // 获取到的凭证
	ExpiresIn   int64  `json:"expire_in"`    // 凭证有效时间，单位：秒。目前是7200秒之内的值。
	ErrCode     int64  `json:"errcode"`      // 错误码
	ErrMsg      string `json:"errmsg"`       // 错误信息
}

func initAccessToken() {
	userModule := mctx.Registry.Get("user")
	if userModule == nil {
		mctx.Logger.Errorf("user module not found")
		return
	}
	expAppid, ok := userModule.Export("appid")
	if !ok {
		mctx.Logger.Errorf("appid not found")
		return
	}
	expSecret, ok := userModule.Export("appsecret")
	if !ok {
		mctx.Logger.Errorf("appsecret not found")
		return
	}
	appid, _ := expAppid.(string)
	secret, _ := expSecret.(string)
	param := map[string]string{
		"grant_type": "client_credential",
		"appid":      appid,
		"secret":     secret,
	}

	var wxResp *wxAccessTokenResponse
	var err error
	queryAccessToken := func() {
		for {
			wxResp, err = util.HTTPRequest[wxAccessTokenResponse](accessTokenURL, "GET", param, nil)
			if err != nil {
				mctx.Logger.Errorf("get access token failed: %s", err)
				return
			}
			if wxResp.ErrCode != 0 {
				mctx.Logger.Errorf("get access token failed: %s", wxResp.ErrMsg)
				return
			}
			if wxResp.ErrCode == -1 {
				continue
			}
			setAccessToken(wxResp.AccessToken)
			break
		}
	}
	queryAccessToken()
	if err != nil {
		mctx.Logger.Error("get access token failed, access token service will be unavailable")
		return
	}
	mctx.Scheduler.Every(wxResp.ExpiresIn - 200).Seconds().SingletonMode().Do(queryAccessToken)
}
