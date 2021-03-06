package system

import (
	"context"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"go-admin/global"
	"go-admin/model/system"
	"go-admin/model/system/request"
	"go-admin/model/system/response"
	"go-admin/utils"
	upload2 "go-admin/utils/upload"
	"time"
)

type UserService struct {
}

var ctx = context.Background()

// Register 用户注册
func (u *UserService) Register(user *request.Register) (err error) {
	//1.校验是否注册
	if _, num := system.CheckUserExist(user.Username); num > 0 {
		return response.ErrorUserExit
	}
	//2.生成UUID
	UUID := uuid.NewV4()
	password := utils.MD5V(user.Password)

	//用户多角色
	var authorities []system.SysAuthority
	for _, v := range user.AuthorityIds {
		authorities = append(authorities, system.SysAuthority{
			AuthorityId: v,
		})
	}
	//3.用户数据
	userData := system.SysUser{
		UUID:        UUID,
		Username:    user.Username,
		Password:    password,
		NickName:    user.NickName,
		Sex:         user.Sex,
		Email:       user.Email,
		Phone:       user.Phone,
		HeaderImg:   user.HeaderImg,
		AuthorityId: user.AuthorityId,
		Authorities: authorities,
	}
	return system.InsertUser(&userData)
}

// Login 用户登录
func (u *UserService) Login(l *request.Login) (user *system.SysUser, err error) {
	//1.查询用户
	//未注册用户
	user, num := system.CheckUserExist(l.Username)
	if num == 0 {
		return nil, response.ErrorUserNotExit
	}
	//2.校验用户密码
	oPasssword := l.Password
	password := utils.MD5V(oPasssword)
	if password != user.Password {
		return nil, response.ErrorPasswordWrong
	}
	return user, nil
}

// GetUserList 用户列表
func (u *UserService) GetUserList(p *request.GetUserList) (list interface{}, total int64, err error) {
	list, total, err = system.UserList(int(p.Page), int(p.Limit), p.Status, p.Username, p.Phone, p.AuthorityId)
	return
}

// GetUserList 删除用户
func (u *UserService) DelUser(id string) (err error) {
	ids := utils.StringToSlice(id)
	return system.DelUser(ids)
}

// UpdateUser 更新用户
func (u *UserService) UpdateUser(user *request.EditUserParam) (err error) {
	if len(user.AuthorityIds) > 0 {
		err = system.SetUserAuthorities(uint(user.ID), user.AuthorityIds)
		if err != nil {
			return
		}
	}
	u2 := system.SysUser{
		GA_MODEL:  global.GA_MODEL{ID: uint(user.ID)},
		Username:  user.Username,
		NickName:  user.NickName,
		HeaderImg: user.HeaderImg,
		Sex:       user.Sex,
		Email:     user.Email,
		Phone:     user.Phone,
		Status:    user.Status,
	}
	return system.UpdateUser(&u2)
}

// EditPassword 编辑密码
func (*UserService) EditPassword(p *request.ChangePasswordStruct) (err error) {
	oldPassword := p.Password
	if p.IsOldPwd == "0" {
		oldPassword = utils.MD5V(p.Password)
	}
	NewPassword := utils.MD5V(p.NewPassword)
	u2 := &system.SysUser{
		GA_MODEL: global.GA_MODEL{
			ID: uint(p.UserID),
		},
		Password: oldPassword,
	}
	return system.EditPassword(u2, NewPassword)
}

// EditUserStatus 更新状态
func (u *UserService) EditUserStatus(userid int, status string) (err error) {
	u2 := system.SysUser{
		GA_MODEL: global.GA_MODEL{ID: uint(userid)},
		Status:   status,
	}
	return system.UpdateUser(&u2)
}

// GetUserInfo 获取用户信息
func (u *UserService) GetUserInfo(uuid uuid.UUID) (user *system.SysUser, err error) {
	return system.GetUserInfo(uuid)
}

// DelUserAvater 删除用户头像
func (u *UserService) DelUserAvater(p *request.DelUserAvaterParam) (err error) {
	userid := p.UserId
	u2 := system.SysUser{
		GA_MODEL:  global.GA_MODEL{ID: uint(userid)},
		HeaderImg: "",
	}
	fmt.Println(u2)
	//更新用户信息
	err = system.UpdateUser(&u2)
	if err != nil {
		return
	}
	//删除头像图片
	upload := upload2.NewFileStore()
	return upload.DeleteFile(p.HeaderImg)
}

// SetUserTokenBlackList 设置用户token黑名单
func (u *UserService) SetUserTokenBlackList(token string) (err error) {
	// 检查是否存在token
	res, err := global.GA_REDIS.SMembers(ctx, token).Result()
	if len(res) == 0 {
		return global.GA_REDIS.SAdd(ctx, "user:blacklist", token).Err()
	}
	return nil
}

// GetUserTokenBlackList 获取用户token是否在黑名单中
func (u *UserService) GetUserTokenBlackList(token string) (val bool, err error) {
	return global.GA_REDIS.SIsMember(ctx, "user:blacklist", token).Result()
}

// GetUserToken 获取用户Token
func (u *UserService) GetUserToken(userName string) (token string, err error) {
	key := "user:token:" + userName
	return global.GA_REDIS.Get(context.Background(), key).Result()
}

// SetUserToken 设置用户Token
func (u *UserService) SetUserToken(token, userName string) (err error) {
	key := "user:token:" + userName
	timer := time.Duration(global.GA_CONFIG.JwtConfig.ExpiresTime) * time.Second
	return global.GA_REDIS.Set(ctx, key, token, timer).Err()
}

// DelUserToken 删除用户token
func (u *UserService) DelUserToken(userName string) (err error) {
	return global.GA_REDIS.Del(ctx, userName).Err()
}
