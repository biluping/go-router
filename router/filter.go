package router

import "net/http"

// Filter 过滤器，用于执行 controller 方法的前置处理，通过 router.AddFilter() 方法注册
type Filter func(write http.ResponseWriter, request *http.Request, chain *FilterChain)

// FilterChain 过滤器链，多个过滤器串成一条链
type FilterChain struct {
	index int
}

// DoFilter 每个过滤器中需要调用 FilterChain 的 DoFilter 方法才能执行下一个过滤器
// 如果中途某个过滤器没有执行 FilterChain 的 DoFilter 方法，那么请求就被拦截不再继续往下
func (f *FilterChain) DoFilter(write http.ResponseWriter, request *http.Request) {
	f.index += 1
	if len(filterList) == f.index {
		return
	}
	filterList[f.index](write, request, f)
}
