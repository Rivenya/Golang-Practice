package micro

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"
)

var Wating = 0
var Running = 1
var WrongStateError = error{msg: "agent state error"}

/*
 * 错误处理模块
 */
type error struct {
	msg            string
	CollectorError []string
}

func (e *error) Error() string {
	return e.msg
}

/*
 * Micro主程序模块
 */

/* 定义事件内容，例如有名字和内容 */
type Event struct {
	name    string
	content string
}

/* 定义代理类，内部拥有很多collector */
type Agent struct {
	collectors map[string]Collector
	evtBuf     chan Event
	cancel     context.CancelFunc
	ctx        context.Context
	state      int
}

/* 定义事件接受者，每个插件含有event事件 */
type EventReceiver interface {
	OnEvent(evt Event)
}

/*
	定义每个事件的四个方法
*/
type Collector interface {
	Init(evtReceiver EventReceiver) error
	Start(ctx context.Context) error
	Stop() error
	Destroy() error
}

/* 给代理类实现事件接收方法 */
func (agt *Agent) OnEvent(evt Event) {
	agt.evtBuf <- evt
}

/* 事件运行机 */
func (agt *Agent) EventProcessGroutine() {
	var evtArg [10]Event
	for i := 0; i < 10; i++ {
		select {
		case evtArg[i] = <-agt.evtBuf:
			fmt.Println("我是Groutine，我取到值了")
		case <-agt.ctx.Done():
			fmt.Println("我是Groutine，我接收到取消了")
			return
		}
	}
	fmt.Println(evtArg)
}

/* 快捷返回agent的指针 */
func NewAgent(sizeEvtBuf int) *Agent {
	agt := Agent{
		collectors: map[string]Collector{},
		evtBuf:     make(chan Event, sizeEvtBuf),
		state:      Wating,
	}
	return &agt
}

/*
 * name: RegisterCollector
 * description: 完成注册
 * return: 注册的init调用
 */
func (agt *Agent) RegisterCollector(name string, collector Collector) error {
	if agt.state != Wating {
		return WrongStateError
	}
	agt.collectors[name] = collector
	fmt.Println(name, "注册完成")
	return collector.Init(agt)
}

func (agt *Agent) Start() error {
	if agt.state != Wating {
		return WrongStateError
	}
	agt.state = Running
	agt.ctx, agt.cancel = context.WithCancel(context.Background())
	go agt.EventProcessGroutine()
	return agt.startCollectors()
}

func (agt *Agent) startCollectors() error {
	var err error
	var errs error
	var mut sync.Mutex
	for name, collector := range agt.collectors {
		go func(name string, collector Collector, ctx context.Context) {
			defer func() {
				mut.Unlock()
			}()
			err = collector.Start(ctx)
			mut.Lock()
			if err.msg != "" {
				errs.CollectorError = append(errs.CollectorError, errors.New(name+":"+err.Error()).Error())
			}
		}(name, collector, agt.ctx)
	}
	return errs
}

func (agt *Agent) Stop() error {
	if agt.state != Running {
		return WrongStateError
	}
	agt.state = Wating
	agt.cancel()
	fmt.Println("我进行context的取消了")
	return agt.stopCollectors()
}

func (agt *Agent) stopCollectors() error {
	var err error
	var errs error
	for name, collector := range agt.collectors {
		if err = collector.Stop(); err.msg != "" {
			errs.CollectorError = append(errs.CollectorError, errors.New(name+":"+err.Error()).Error())
		}
	}
	return errs
}

func (agt *Agent) destoryCollectors() error {
	var err error
	var errs error
	for name, collector := range agt.collectors {
		if err = collector.Destroy(); err.msg != "" {
			errs.CollectorError = append(errs.CollectorError, errors.New(name+":"+err.Error()).Error())
		}
	}
	return errs
}

func (agt *Agent) Desroty() error {
	if agt.state != Wating {
		return WrongStateError
	}
	return agt.destoryCollectors()
}

type DemoCollector struct {
	evtReceiver EventReceiver
	stopChan    chan struct{}
	name        string
	content     string
}

func NewCollect(name string, content string) *DemoCollector {
	return &DemoCollector{
		stopChan: make(chan struct{}),
		name:     name,
		content:  content,
	}
}

func (d *DemoCollector) Init(evt EventReceiver) error {
	fmt.Println(d.name, ":初始化")
	d.evtReceiver = evt
	return error{msg: ""}
}

func (d *DemoCollector) Start(ctx context.Context) error {
	fmt.Println(d.name, ":开始")
	for {
		select {
		case <-ctx.Done():
			fmt.Println("我接收到cancel了")
			d.stopChan <- struct{}{}
		default:
			time.Sleep(time.Millisecond * 50)
			fmt.Println("我在输出", d.name, ",", d.content)
			d.evtReceiver.OnEvent(Event{d.name, d.content})
		}
	}
}

func (d *DemoCollector) Stop() error {
	fmt.Println(d.name, ":停止")
	select {
	case <-d.stopChan:
		fmt.Println("取消成功")
		return error{msg: ""}
	case <-time.After(time.Second * 2):
		fmt.Println("取消超时")
		return error{msg: errors.New("fail to stop for timeout").Error()}
	}
}

func (d *DemoCollector) Destroy() error {
	fmt.Println(d.name, ":终止")
	return error{msg: ""}
}

func Test_micro(t *testing.T) {
	agt := NewAgent(100)

	lion := NewCollect("师子", "喜欢吃肉")
	rabbit := NewCollect("兔子", "喜欢吃胡萝卜")

	agt.RegisterCollector("师子", lion)
	agt.RegisterCollector("兔子", rabbit)

	agt.Start()

	time.Sleep(time.Second * 1)

	agt.Stop()

	agt.Desroty()
}
