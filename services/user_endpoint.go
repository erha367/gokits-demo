package services

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

/**
 * @description：定义基本结构
 * @author     ：yangsen
 * @date       ：2020/11/20 3:35 下午
 * @company    ：eeo.cn
 */

type UserRequest struct {
	Uid int `json:"uid"`
}

type UserResponse struct {
	Res string `json:"result"`
}

//注意参数使用的是user接口，返回值固定为 endpoint.Endpoint
func GetUserEndpoint(userService IUserService) endpoint.Endpoint {
	//return 格式固定
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		//注意这里request的interface抽象为UserRequest
		r := request.(UserRequest)
		//方法调用
		res := userService.GetName(r.Uid)
		//返回
		return UserResponse{Res: res}, nil
	}
}
