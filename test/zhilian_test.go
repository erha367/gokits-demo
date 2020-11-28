package method_test

import (
	"context"
	"encoding/json"
	"errors"
	htport `github.com/go-kit/kit/transport/http`
	"gokits/services"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"testing"
)

/**
 * @description：测试rpc服务
 * @author     ：yangsen
 * @date       ：2020/11/23 2:33 下午
 * @company    ：eeo.cn
 */

func TestZhiLian(t *testing.T) {
	tgt, _ := url.Parse(`http://127.0.0.1:89`)
	//创建client
	client := htport.NewClient(`GET`, tgt, enc, dec)
	//创建endpoint
	ept := client.Endpoint()
	//创建ctx
	ctx := context.Background()
	//传参
	res, err := ept(ctx, services.UserRequest{Uid: 50})
	t.Log(res, err)
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
		log.Println(res.StatusCode)
		return nil, errors.New(`err response code`)
	}
	var userRes services.UserResponse
	err = json.NewDecoder(res.Body).Decode(&userRes)
	if err != nil {
		return nil, err
	}
	return userRes, nil
}
