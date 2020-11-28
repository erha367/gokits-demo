package main

import (
	"flag"
	"fmt"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"gokits/services"
	"gokits/utils"
	"golang.org/x/time/rate"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

/**
 * @description：TODO
 * @author     ：yangsen
 * @date       ：2020/11/20 4:18 下午
 * @company    ：eeo.cn
 */

func main() {
	//参数解析
	flag.IntVar(&services.Port, "port", 88, "使用端口")
	flag.Parse()
	//加载服务
	user := services.UserService{}
	//限流
	limits := rate.NewLimiter(1, 3)
	end := utils.RateLimit(limits)(services.GetUserEndpoint(user))
	//自定义错误
	options := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(utils.MyErrorEncoder),
	}
	serverHandler := kithttp.NewServer(end, services.DecodeUserRequest, services.EncodeUserResponse, options...)
	//路由
	r := mux.NewRouter()
	r.Methods("GET").Path(`/user/{uid:\d+}`).Handler(serverHandler)
	r.Methods("GET").Path(`/health`).HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set(`Content-type`, `application/json`)
		writer.Write([]byte(`{"status":"OK"}`))
	})
	//错误处理
	errChan := make(chan error)
	go func() {
		utils.RegService(services.Port)
		err := http.ListenAndServe(`:`+strconv.Itoa(services.Port), r)
		if err != nil {
			errChan <- err
		}
	}()
	go func() {
		sigc := make(chan os.Signal)
		signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-sigc)
	}()
	e := <-errChan
	utils.UnRegService()
	log.Println(e)
}
