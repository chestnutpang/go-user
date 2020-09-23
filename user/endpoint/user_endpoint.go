package endpoint
/*
构建 RegisterEndpoint 与 LoginEndpoint
将请求转化成 UserService 接口可处理的参数
并将处理的结果封装成对应 response 结构体返回 transport 包
*/

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/longjoy/micro-go-course/service"
	"log"
)

type UserEndpoints struct{
	RegisterEndpoint endpoint.Endpoint
	LoginEndpoint endpoint.Endpoint
}


// 登录处理相关函数
type LoginRequest struct{
	Email string
	Password string
}

type LoginResponse struct{
	UserInfo *service.UserInfoDTO `json:"user_info"`
}

func MakeLoginEndpoint(userService service.UserService) endpoint.Endpoint{
	// 解析 LoginRequest 中的参数传递给 UserService.Login 方法处理并将处理结果封装为 LoginResponse 返回
	return func(ctx context.Context, request interface{})(response interface{}, err error){
		req := request.(*LoginRequest)
		userInfo, err := userService.Login(ctx, req.Email, req.Password)
		return &LoginResponse{UserInfo:userInfo}, err
	}
}


// 注册处理相关函数
type RegisterRequest struct{
	Username string
	Email string
	Password string
}

type RegisterResponse struct{
	UserInfo *service.UserInfoDTO `json:"user_info"`
}

func MakeRegisterEndpoint(userService service.UserService) endpoint.Endpoint{
	// 解析 RegisterRequest 中的参数传递给 UserService.Register 方法处理并将处理的结果封装
	return func(ctx context.Context, request interface{})(response interface{}, err error){
		req := request.(*RegisterRequest)
		log.Println(req.Username, req.Password,req.Email)
		userInfo, err := userService.Register(ctx, &service.RegisterUserVO{
			Username: req.Username,
			Password: req.Password,
			Email: req.Email,
		})
		return &RegisterResponse{UserInfo: userInfo}, err
	}
}
