package model

import "github.com/xaxys/maintainman/core/util"

type Page struct {
	Entries interface{} `json:"entries"`
	Total   uint        `json:"total"`
}

func PagedData(entries interface{}, total uint) *Page {
	return &Page{Entries: entries, Total: total}
}

type ApiJson struct {
	Code   int         `json:"code"`
	Status bool        `json:"status"`
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data"`
}

func ApiResponse(code int, status bool, objects interface{}, msg string) *ApiJson {
	return &ApiJson{Code: code, Status: status, Data: objects, Msg: msg}
}

func combineError(errs ...error) (errMsg []string) {
	return util.TransSlice(errs, func(err error) string { return err.Error() })
}

// Success 成功
func Success(objects interface{}, msg string) *ApiJson {
	return ApiResponse(200, true, objects, msg)
}

// SuccessPaged 成功
func SuccessPaged(entries interface{}, total uint, msg string) *ApiJson {
	return ApiResponse(200, true, PagedData(entries, total), msg)
}

// SuccessUpdate 成功更新/删除资源
func SuccessUpdate(objects interface{}, msg string) *ApiJson {
	return ApiResponse(204, true, objects, msg)
}

// SuccessCreate 成功创建资源
func SuccessCreate(objects interface{}, msg string) *ApiJson {
	return ApiResponse(201, true, objects, msg)
}

// ErrorInsertDatabase 插入数据库失败
func ErrorInsertDatabase(errs ...error) *ApiJson {
	return ApiResponse(500, false, combineError(errs...), "插入数据库失败")
}

// ErrorQueryDatabase 查询数据库失败
func ErrorQueryDatabase(errs ...error) *ApiJson {
	return ApiResponse(500, false, combineError(errs...), "查询数据库失败")
}

// ErrorUpdateDatabase 更新数据库失败
func ErrorUpdateDatabase(errs ...error) *ApiJson {
	return ApiResponse(500, false, combineError(errs...), "更新数据库失败")
}

// ErrorDeleteDatabase 删除数据库失败
func ErrorDeleteDatabase(errs ...error) *ApiJson {
	return ApiResponse(500, false, combineError(errs...), "删除数据库失败")
}

// ErrorNotFound 未找到数据记录
func ErrorNotFound(errs ...error) *ApiJson {
	return ApiResponse(404, false, combineError(errs...), "未找到数据记录")
}

// ErrorInvalidData 数据解析失败
func ErrorInvalidData(errs ...error) *ApiJson {
	return ApiResponse(400, false, combineError(errs...), "数据解析失败")
}

// ErrorIncompleteData 数据不完整
func ErrorIncompleteData(errs ...error) *ApiJson {
	return ApiResponse(422, false, combineError(errs...), "数据不完整")
}

// ErrorValidation 数据检验失败
func ErrorValidation(errs ...error) *ApiJson {
	return ApiResponse(422, false, combineError(errs...), "数据检验失败")
}

// ErrorBuildJWT 生成凭证错误
func ErrorBuildJWT(errs ...error) *ApiJson {
	return ApiResponse(500, false, combineError(errs...), "生成凭证错误")
}

// ErrorUnauthorized 未认证登录
func ErrorUnauthorized(errs ...error) *ApiJson {
	return ApiResponse(401, false, combineError(errs...), "未认证登录")
}

// ErrorVerification 认证失败
func ErrorVerification(errs ...error) *ApiJson {
	return ApiResponse(403, false, combineError(errs...), "认证失败")
}

// ErrorNoPermissions 账号权限不足
func ErrorNoPermissions(errs ...error) *ApiJson {
	return ApiResponse(403, false, combineError(errs...), "账号权限不足")
}

// ErrorInternalServer 服务器内部错误
func ErrorInternalServer(errs ...error) *ApiJson {
	return ApiResponse(500, false, combineError(errs...), "服务器内部错误")
}
