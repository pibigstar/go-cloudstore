package service

import (
	"context"
	"github.com/pibigstar/go-cloudstore/constant"
	"github.com/pibigstar/go-cloudstore/db"
	"github.com/pibigstar/go-cloudstore/service/account/proto"
	"github.com/pibigstar/go-cloudstore/utils"
)

type User struct{}

// 提供User注册服务，这个是为其他rpc client提供的
func (*User) UserSignup(context context.Context, request *proto.ReqSignup, response *proto.RespSignup) error {

	username := request.Username
	password := request.Password

	enc_pwd := utils.Sha1([]byte(password + cont.PASSWORD_SALT))
	b := db.UserSignup(username, enc_pwd)
	if b {
		response.Code = cont.OperatorSuccess
		response.Message = "注册成功！"
	} else {
		response.Code = cont.UserSignupFailed
		response.Message = "注册失败！"
	}
	return nil
}
