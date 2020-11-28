package services

import "strconv"

/**
 * @description：TODO
 * @author     ：yangsen
 * @date       ：2020/11/20 3:32 下午
 * @company    ：eeo.cn
 */

//定义接口
type IUserService interface {
	GetName(userId int) string
}

//接口继承
type UserService struct{}
var Port int

//定义方法
func (u UserService) GetName(userId int) string {
	if userId > 100 {
		return `yangsen-` + strconv.Itoa(Port)
	}
	return `none`
}
