package router

import "net/http"

type Controller func(request *http.Request) (interface{}, error)

// {请求方法: {请求路径: 执行函数}}
var routerMap = make(map[string]map[string]Controller)

// 过滤器链
var filterList []Filter
