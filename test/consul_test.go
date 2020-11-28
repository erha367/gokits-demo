package method

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/consul"
	"github.com/go-kit/kit/sd/lb"
	htport `github.com/go-kit/kit/transport/http`
	consulApi "github.com/hashicorp/consul/api"
	"gokits/services"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"testing"
	"time"
)

/**
 * @description：TODO
 * @author     ：yangsen
 * @date       ：2020/11/23 3:12 下午
 * @company    ：eeo.cn
 */

func TestConsul(t *testing.T) {
	/*- 1.创建client -*/
	config := consulApi.DefaultConfig()
	config.Address = `127.0.0.1:8500`
	apiClient, _ := consulApi.NewClient(config)
	client := consul.NewClient(apiClient)
	//日志
	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stdout)
	/*- 2.创建实例-*/
	tags := []string{`primary`}
	instancer := consul.NewInstancer(client, logger, `user_service`, tags, true)
	/*- 3.查询服务状态 -*/
	endpointer := sd.NewEndpointer(instancer, factory, logger)
	//轮询
	//mylb := lb.NewRoundRobin(endpointer)
	//随机
	mylb := lb.NewRandom(endpointer, time.Now().UnixNano())
	t.Log(mylb.Endpoint())
	/*- 4.获取第一个 -*/
	//points, err := endpointer.Endpoints()
	//t.Log(points, err)
	//getUser := points[1]
	/*- 5.负载均衡算法=简单轮询，随机 -*/
	for i := 0; i < 4; i++ {
		time.Sleep(time.Second)
		if getUser, err := mylb.Endpoint();err == nil {
			ctx := context.Background()
			res, err := getUser(ctx, services.UserRequest{Uid: 888})
			t.Log(res, err)
		}else{
			t.Log(err)
		}
	}
}

func factory(serUrl string) (endpoint.Endpoint, io.Closer, error) {
	tart, _ := url.Parse(`http://` + serUrl)
	return htport.NewClient(`GET`, tart, enc, dec).Endpoint(), nil, nil
}

//请求部分
func enc(c context.Context, req *http.Request, r interface{}) error {
	userRequest := r.(services.UserRequest)
	req.URL.Path += `/user/` + strconv.Itoa(userRequest.Uid)
	return nil
}

//返回部分
func dec(c context.Context, res *http.Response) (response interface{}, err error) {
	if res.StatusCode != 200 {
		return nil, errors.New(`err response code`)
	}
	var userRes services.UserResponse
	err = json.NewDecoder(res.Body).Decode(&userRes)
	if err != nil {
		return nil, err
	}
	return userRes, nil
}
