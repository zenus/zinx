package znet

import (
	"github.com/zenus/zinx/ziface"
	"sync"
)

// BaseRouter 实现router时，先嵌入这个基类，然后根据需要对这个基类的方法进行重写
type BaseRouter struct{}

//这里之所以BaseRouter的方法都为空，
// 是因为有的Router不希望有PreHandle或PostHandle
// 所以Router全部继承BaseRouter的好处是，不需要实现PreHandle和PostHandle也可以实例化

// PreHandle -
func (br *BaseRouter) PreHandle(req ziface.IRequest) {}

// Handle -
func (br *BaseRouter) Handle(req ziface.IRequest) {}

// PostHandle -
func (br *BaseRouter) PostHandle(req ziface.IRequest) {}

//
//
//
//

// 新切片集合式路由
// 新版本路由基本逻辑,用户可以传入不等数量的路由路由处理器
// 路由本体会讲这些路由处理器函数全部保存,在请求来的时候找到，并交由IRequest去执行
// 路由可以设置全局的共用组件通过Use方法
// 路由可以分组,通过Group,分组也有自己对应Use方法设置组共有组件

type RouterSlices struct {
	Apis     map[string][]ziface.RouterHandler
	Handlers []ziface.RouterHandler
	sync.RWMutex
}

func NewRouterSlices() *RouterSlices {
	return &RouterSlices{
		Apis:     make(map[string][]ziface.RouterHandler, 10),
		Handlers: make([]ziface.RouterHandler, 0, 6),
	}
}

func (r *RouterSlices) Use(handles ...ziface.RouterHandler) {
	r.Handlers = append(r.Handlers, handles...)
}

func (r *RouterSlices) AddHandler(cmd string, Handlers ...ziface.RouterHandler) {
	//1 判断当前msg绑定的API处理方法是否已经存在
	if _, ok := r.Apis[cmd]; ok {
		panic("repeated api , cmd = " + cmd)
	}

	finalSize := len(r.Handlers) + len(Handlers)
	mergedHandlers := make([]ziface.RouterHandler, finalSize)
	copy(mergedHandlers, r.Handlers)
	copy(mergedHandlers[len(r.Handlers):], Handlers)
	r.Apis[cmd] = append(r.Apis[cmd], mergedHandlers...)
}

func (r *RouterSlices) GetHandlers(cmd string) ([]ziface.RouterHandler, bool) {
	r.RLock()
	defer r.RUnlock()
	handlers, ok := r.Apis[cmd]
	return handlers, ok
}

//func (r *RouterSlices) Group(start, end uint32, Handlers ...ziface.RouterHandler) ziface.IGroupRouterSlices {
//	return NewGroup(start, end, r, Handlers...)
//}

//type GroupRouter struct {
//	start    uint32
//	end      uint32
//	Handlers []ziface.RouterHandler
//	router   ziface.IRouterSlices
//}
//
//func NewGroup(start, end uint32, router *RouterSlices, Handlers ...ziface.RouterHandler) *GroupRouter {
//	g := &GroupRouter{
//		start:    start,
//		end:      end,
//		Handlers: make([]ziface.RouterHandler, 0, len(Handlers)),
//		router:   router,
//	}
//	g.Handlers = append(g.Handlers, Handlers...)
//	return g
//}
//
//func (g *GroupRouter) Use(Handlers ...ziface.RouterHandler) {
//	g.Handlers = append(g.Handlers, Handlers...)
//}
//
//func (g *GroupRouter) AddHandler(cmd string, Handlers ...ziface.RouterHandler) {
//	finalSize := len(g.Handlers) + len(Handlers)
//	mergedHandlers := make([]ziface.RouterHandler, finalSize)
//	copy(mergedHandlers, g.Handlers)
//	copy(mergedHandlers[len(g.Handlers):], Handlers)
//	//回调实际路由的添加组件
//	g.router.AddHandler(cmd, mergedHandlers...)
//}
