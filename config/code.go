package config

type ResCode int64

//定义返回状态
const (
	CodeSuccess ResCode = 1000 + iota
	CodeInvalidParam
	CodeUserExist
	CodeUserNotExist
	CodeInvalidPassword
	CodeServerBusy
	CodeNeedLogin
	CodeInvalidToken
	CodeMenuExist
	CodeMenuChildExist
	CodeAuthExist
	CodeAuthChildExit
	CodeAuthApiExit
	CodeNoPermission
	CodeFileError
)

//定义返回信息
var codeMsgMap = map[ResCode]string{
	CodeSuccess:         "success",
	CodeInvalidParam:    "请求参数有误",
	CodeUserExist:       "用户存在",
	CodeUserNotExist:    "用户不存在",
	CodeInvalidPassword: "用户名或密码错误",
	CodeServerBusy:      "服务繁忙",
	CodeInvalidToken:    "无效Token",
	CodeNeedLogin:       "需要登录",
	CodeMenuExist:       "菜单已存在",
	CodeMenuChildExist:  "菜单存在下级",
	CodeAuthExist:       "角色已存在",
	CodeAuthChildExit:   "角色存在下级",
	CodeAuthApiExit:     "api已存在",
	CodeNoPermission:    "权限不足",
	CodeFileError:       "文件接受失败",
}

func (c ResCode) Msg() string {
	msg, ok := codeMsgMap[c]
	if !ok {
		return codeMsgMap[CodeServerBusy]
	}
	return msg
}
