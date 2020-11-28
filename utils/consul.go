package utils

/**
 * @description：TODO
 * @author     ：yangsen
 * @date       ：2020/11/23 11:15 上午
 * @company    ：eeo.cn
 */

import (
	consulApi "github.com/hashicorp/consul/api"
	"log"
	"strconv"
)

var ConsulClient *consulApi.Client
var err error

func init() {
	//配置
	config := consulApi.DefaultConfig()
	config.Address = `127.0.0.1:8500`
	ConsulClient, err = consulApi.NewClient(config)
	if err != nil {
		log.Fatal(`consul connect err`, err)
	}
}

func RegService(port int) {
	//注册
	reg := consulApi.AgentServiceRegistration{
		ID:      `user-` + strconv.Itoa(port),
		Name:    `user_service`,
		Address: `127.0.0.1`,
		Port:    port,
		Tags:    []string{`primary`},
	}
	//健康检查
	check := consulApi.AgentServiceCheck{
		Interval: `5s`,
		HTTP:     `http://127.0.0.1:` + strconv.Itoa(port) + `/health`,
	}
	reg.Check = &check
	err = ConsulClient.Agent().ServiceRegister(&reg)
	if err != nil {
		log.Fatal(`consul reg err`, err)
	}
	log.Println(`consul reg success`)
}

func UnRegService() {
	ConsulClient.Agent().ServiceDeregister(`user01`)
}
