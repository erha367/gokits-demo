package method

import (
	"context"
	"golang.org/x/time/rate"
	"testing"
	"time"
)

/**
 * @description：TODO
 * @author     ：yangsen
 * @date       ：2020/11/25 5:56 下午
 * @company    ：eeo.cn
 */

func TestTong(t *testing.T){
	//每秒放1个，最大值存5个
	r := rate.NewLimiter(1,5)
	ctx := context.Background()
	for {
		//每次消耗2个，相当于放3个使用2个，后面每2秒执行1次 【阻塞等待】
		err := r.WaitN(ctx,2)
		if err != nil {
			t.Error(err)
		}
		time.Sleep(time.Second)
		t.Log(time.Now().Format(`2006-01-02 15:04:05`))
	}
}

func TestTong2(t *testing.T){
	//每秒放1个，最大值存5个
	r := rate.NewLimiter(1,5)
	for {
		//每次消耗2个，相当于放3个使用2个，后面每2秒执行1次 【阻塞等待】
		if r.AllowN(time.Now(),2) {
			//模拟业务代码
			t.Log(time.Now().Format(`2006-01-02 15:04:05`))
		}else{
			t.Log(`too quickly.`)
		}
		time.Sleep(time.Second)
	}
}