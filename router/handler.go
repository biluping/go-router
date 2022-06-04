package router

import (
	"log"
	"net/http"
)

// 全局路由处理器，用于匹配请求方法、请求路径和处理函数
func routerHandle(write http.ResponseWriter, request *http.Request) {
	// 全局异常处理
	defer func() {
		err := recover()
		if err != nil {
			log.Println(err)
			ResponseBadRequest(write, "")
		}
	}()

	// 寻找请求方法对应的map
	m, exist := routerMap[request.Method]
	if !exist {
		responseNotFound(write)
		return
	}

	// 寻找请求路径对于的 controller 函数
	u, exist := m[request.RequestURI]
	if !exist {
		responseNotFound(write)
		return
	}

	// 执行过滤器
	chain := FilterChain{index: -1}
	chain.DoFilter(write, request)
	if chain.index < len(filterList) {
		// 过滤器应该自己处理响应，这里不做处理
		return
	}

	// 执行 controller 方法
	data, err := u(request)
	if err != nil {
		ResponseBadRequest(write, err.Error())
		return
	}
	responseOk(write, data)
}

// 通用 controller 注册函数，仅限内部使用
func register(method string, path string, handle Controller) {
	if _, ok := routerMap[method]; !ok {
		routerMap[method] = make(map[string]Controller)
	}
	m := routerMap[method]
	if _, exist := m[path]; exist {
		log.Panicf("method %s, path %s has exist", method, path)
	}
	m[path] = handle
}
