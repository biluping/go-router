package router

// Controller func(...) interface{}
type Controller interface{}

// {请求方法: {请求路径: 执行函数}}
var routerMap = make(map[string]map[string]Controller)

// 过滤器链
var filterList []Filter
