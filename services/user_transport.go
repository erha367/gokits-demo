package services

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

/**
 * @description：TODO
 * @author     ：yangsen
 * @date       ：2020/11/20 4:06 下午
 * @company    ：eeo.cn
 */

//访问http://127.0.0.1:xxxx/user/100
func DecodeUserRequest(c context.Context, r *http.Request) (interface{}, error) {
	parms := mux.Vars(r)
	if uid, ok := parms["uid"]; ok {
		uidStr, _ := strconv.Atoi(uid)
		return UserRequest{Uid: uidStr}, nil
	}
	return UserRequest{}, errors.New(`参数错误`)
}

func EncodeUserResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set(`Content-type`, `application/json`)
	return json.NewEncoder(w).Encode(response)
}
