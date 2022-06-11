package router

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"runtime"
)

// 初始化方法，注册全局处理器，类似于 java 中的 DispatcherServlet 的作用
func init() {
	http.HandleFunc("/", routerHandle)
}

// 全局路由处理器，用于匹配请求方法、请求路径和处理函数
func routerHandle(write http.ResponseWriter, request *http.Request) {
	// 全局异常处理
	defer func() {
		err := recover()
		if err == nil {
			return
		}
		switch err.(type) {
		case runtime.Error:
			log.Println(err)
			ResponseBadRequest(write, err.(error).Error())
		default:
			log.Println(err)
		}
	}()

	// 寻找请求方法对应的map
	m, exist := routerMap[request.Method]
	if !exist {
		responseNotFound(write)
		return
	}

	// 寻找请求路径对于的 controller 函数
	u, exist := m[request.URL.Path]
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
	data := invokeFunc(u, request)
	responseOk(write, data)
}

func invokeFunc(u Controller, request *http.Request) interface{} {
	// map[string][]string
	query := request.URL.Query()
	// 转换，只取第一个参数
	queryMap := make(map[string]string)
	for k, v := range query {
		queryMap[k] = v[0]
	}

	// 读取请求体
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Panic(err)
	}

	// 将 query 参数转成 json，后面通过反序列号的方式将值设置到接受对象中
	queryJson, err := json.Marshal(queryMap)
	if err != nil {
		log.Panic(err)
	}

	// 反射获取处理函数的参数类型，并生产对应的值，调用函数
	v := reflect.TypeOf(u)
	var parameters []reflect.Value
	for i := 0; i < v.NumIn(); i++ {
		in := v.In(i)

		if in.String() == "*http.Request" {
			parameters = append(parameters, reflect.ValueOf(request))
		} else {
			if in.Kind() == reflect.Pointer {
				// 如果参数是指针类型，需要调用 Elem 方法获取真实类型
				in = in.Elem()
			}
			value := reflect.New(in)
			err = json.Unmarshal(body, value.Interface())
			if err != nil {
				log.Panic(err)
			}
			err := json.Unmarshal(queryJson, value.Interface())
			if err != nil {
				log.Panic(err)
			}
			parameters = append(parameters, value)
		}
	}

	of := reflect.ValueOf(u)
	res := of.Call(parameters)
	return res[0].Interface()

}

// 通用 controller 注册函数，仅限内部使用
func register(method string, path string, handle Controller) {
	of := reflect.TypeOf(handle)
	if of.Kind() != reflect.Func {
		log.Panicf("controller注册的必须是函数")
	}
	if _, ok := routerMap[method]; !ok {
		routerMap[method] = make(map[string]Controller)
	}
	m := routerMap[method]
	if _, exist := m[path]; exist {
		log.Panicf("method %s, path %s has exist", method, path)
	}
	m[path] = handle
}
