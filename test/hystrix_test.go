package method

import (
	"errors"
	"github.com/afex/hystrix-go/hystrix"
	"math/rand"
	"sync"
	"testing"
	"time"
)

/**
 * @description：TODO
 * @author     ：yangsen
 * @date       ：2020/11/28 下午2:34
 * @company    ：eeo.cn
 */

func TestHysYB(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	hysConfig := hystrix.CommandConfig{
		Timeout:               2000, //2s
		MaxConcurrentRequests: 5,    //最大并发数
	}
	hystrix.ConfigureCommand(`cmd`, hysConfig)
	resChan := make(chan product, 1)
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		go func() {
			wg.Add(1)
			defer wg.Done()
			//Do是同步方法 Go是异步方法
			errList := hystrix.Go(`cmd`, func() error {
				p, _ := getProduct()
				resChan <- p
				return errors.New(`超时了`)
			}, func(err error) error {
				t.Log(`jiang ji`, err)
				p, e := cachedProduct()
				resChan <- p
				return e
			})
			select {
			case p := <-resChan:
				t.Log(p)
			case err := <-errList:
				t.Log(err)
			}
		}()
	}
	wg.Wait()
	t.Log(`done`)
}

func TestHysTB(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	hysConfig := hystrix.CommandConfig{
		Timeout:                2000,  //2s
		MaxConcurrentRequests:  2,     //最大并发
		RequestVolumeThreshold: 3,     //请求阈值，默认20
		SleepWindow:            60000, //熔断器打开后，多久尝试恢复系统 60秒
		ErrorPercentThreshold:  10,    //超过10次请求后，触发熔断
	}
	hystrix.ConfigureCommand(`cmd`, hysConfig)
	//获取熔断器状态
	c, _, _ := hystrix.GetCircuit(`cmd`)
	for {
		t.Log(`--------------------------`)
		//Do是同步方法
		err := hystrix.Do(`cmd`, func() error {
			p, _ := getProduct()
			t.Log(p)
			return nil
		}, func(err error) error {
			t.Log(cachedProduct())
			t.Log(time.Now().Format("2006-01-02 15:04:05"))
			return errors.New(`降级是不能取远端服务的`)
		})
		if err != nil {
			t.Log(`hystrix`, err)
		}
		t.Log(c.IsOpen())
		time.Sleep(time.Second)
	}
}

type product struct {
	Pid  int
	Name string
	Desc string
}

func getProduct() (product, error) {
	r := rand.Intn(10)
	if r < 6 {
		time.Sleep(time.Second * 3)
	}
	return product{
		Pid:  100,
		Name: `wsk13-live.eeo.im`,
		Desc: `test-13`,
	}, nil
}

func cachedProduct() (product, error) {
	return product{
		Pid:  10,
		Name: `降级推荐商品`,
		Desc: `prod`,
	}, nil
}
