package method

import (
	"github.com/dgrijalva/jwt-go"
	"testing"
	"time"
)

/**
 * @description：TODO
 * @author     ：yangsen
 * @date       ：2020/11/28 下午4:11
 * @company    ：eeo.cn
 */

type MyCustomClaims struct {
	Foo string `json:"foo"`
	jwt.StandardClaims
}

func TestJwt(t *testing.T) {
	mySigningKey := []byte("AllYourBase")
	// Create the Claims
	claims := MyCustomClaims{
		"www.baidu.com",
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Second).Unix(),
			Issuer:    "test",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(mySigningKey)
	t.Log(ss, err)
	//方式1
	getTk, e := jwt.Parse(ss, func(token *jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	})
	if e != nil {
		t.Log(e)
	}
	if getTk.Valid {
		t.Log(getTk.Claims)
	}
	time.Sleep(time.Second * 2)
	//方式2
	getTk2, e := jwt.ParseWithClaims(ss, &claims, func(token *jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	})
	t.Log(e)
	if getTk2.Valid {
		t.Log(getTk2.Claims.(*MyCustomClaims).Foo)
	}
}
