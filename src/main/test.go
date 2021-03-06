package main

import (
	"fmt"
	"reflect"
	"sync"
	"time"
)

// switch
func trySwitch() {
	b := "0"
	switch {
	case b == "0":
		fmt.Println("猪")
	case b == "1":
		fmt.Println("牛")
	case b == "2":
		fmt.Println("马")
	default:
		fmt.Println("羊")
	}
}

// range遍历
func tryRange() {
	props := [...]int{1, 2, 3, 4, 5, 6}
	for key, val := range props {
		fmt.Println("key是", key)
		fmt.Println(val)
	}
}

// 数组分割
func trySpliceArray() {
	props := [...]int{1, 2, 3, 4, 5}
	props2 := props[1:]
	props3 := props[3:5]
	fmt.Println("props的数组元素有多少", len(props2))
	fmt.Println("props的数组元素有多少", len(props3))
}

// map
func trymMap() {
	a := map[string]string{}
	b := map[string]string{}
	c := make(map[string]string, 10)

	a["banana"] = "香蕉"
	b["banana"] = "香蕉"
	c["banana"] = "香蕉"
	fmt.Println("a的长度", len(a))
	fmt.Println("b的长度", len(b))
	fmt.Println("c的长度", len(c))
}

//函数式编程以及计算函数执行时间
func tryExeSpentFucTime() {
	test := spentFucTime(trySpentFucTime)
	test(999)
}

func spentFucTime(fn func(val int) int) func(val int) int {
	return func(val int) int {
		start := time.Now()
		var ret = fn(val)
		fmt.Println("该函数用时为:", time.Since(start).Seconds(), "秒")
		return ret
	}
}
func trySpentFucTime(val int) int {
	time.Sleep(time.Second * 2)
	return val
}

//可变参数
func canTransferOptFuc(opt ...int) {
	ret := 0
	for _, s := range opt {
		ret += s
	}
	fmt.Println("求和结果是", ret)
}

//结构和接口
type Animals struct {
	Name string
	Food string
}

func (e *Animals) Hungry() {
	fmt.Println(e.Name, "饿了要吃", e.Food)
}
func tryTypeAndInterface() {
	tiger := &Animals{Name: "老虎", Food: "肉"}
	tiger.Hungry()
}

//延迟函数
func deferFunc() {
	defer func() {
		fmt.Println("我是延迟函数")
	}()
	fmt.Println("我在延迟函数后面")
	// panic("err")
	time.Sleep(time.Second * 1)
	fmt.Println("睡了1秒")

}

//多线程
func mutilThread() {
	var mut sync.Mutex
	var wat sync.WaitGroup
	count := 0
	for i := 0; i < 5000; i++ {
		wat.Add(1)
		go func() {
			defer func() {
				mut.Unlock()
			}()
			mut.Lock()
			count++
			wat.Done()
		}()
		wat.Wait()
	}
	fmt.Println(count)
}

//chan
func tryprint(n int) string {
	time.Sleep(1000000000 * 2)
	return "我睡了" + fmt.Sprint(n) + "s醒了"
}
func tryChan(n int) int {
	ret := make(chan string)
	go func() {
		ret <- tryprint(n)
	}()
	fmt.Println("我开始睡了")
	fmt.Println(<-ret)
	return n
}

//关闭chan的实例
func isCancel(n chan struct{}) bool {
	select {
	case <-n:
		return true
	default:
		return false
	}
}
func closeChan(n chan struct{}) {
	n <- struct{}{}
}
func tryCloseChan() {
	var mut sync.WaitGroup
	// 开启通道
	var ret = make(chan struct{}, 1)
	// 开始循环创建线程
	func() {
		mut.Add(1)
		for i := 0; i < 5; i++ {
			go func(n chan struct{}, i int) {
				for {
					if isCancel(n) {
						break
					}
					time.Sleep(time.Second * 2)
				}
				fmt.Println(i, "已经被结束了")
				mut.Done()
			}(ret, i)
		}
		closeChan(ret)
		mut.Wait()
	}()

}

//执行一次异步任务
func tryOnlyAsyncTask() {
	var once sync.Once
	for i := 0; i < 50; i++ {
		go func(num int) {
			once.Do(func() {
				fmt.Println("我是", num)
			})
		}(i)
	}
}

// 第一个异步执行成功就返回
func firstExeSusRetern() int {
	var ch = make(chan int, 1)
	for i := 0; i < 50; i++ {
		go func(num int) {
			time.Sleep(time.Second * 1)
			ch <- num
		}(i)
	}
	return <-ch
}

//reflect实验
func tryReflect(n interface{}) {
	c := reflect.TypeOf(n)
	// println(c.Kind())
	println(c)
}

func main() {
	a := 1
	b := &a
	print(*b)
}
