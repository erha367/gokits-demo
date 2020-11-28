package utils

import (
	"context"
	"errors"
	"github.com/go-kit/kit/endpoint"
	"golang.org/x/time/rate"
	"net/http"
)

/**
 * @description：TODO
 * @author     ：yangsen
 * @date       ：2020/11/28 11:21 上午
 * @company    ：eeo.cn
 */

func RateLimit(limit *rate.Limiter) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if !limit.Allow() {
				return nil, errors.New(`{"err":1,"msg":"too many requests."}`)
			}
			return next(ctx, request)
		}
	}
}

func MyErrorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	contentType, body := `text/plain; charset=utf8`, []byte(err.Error())
	w.Header().Set(`content-type`, contentType)
	w.WriteHeader(502)
	w.Write(body)
}
