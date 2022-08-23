package router

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strings"
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
		log.Println(err)
		switch err.(type) {
		case error:
			ResponseBadRequest(write, err.(error).Error())
		case string:
			ResponseBadRequest(write, err.(string))
		default:

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

	// 将 query 参数转成 json，后面通过反序列化的方式将值设置到接受对象中
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

			// 反序列化请求体
			if len(body) > 0 {
				err = json.Unmarshal(body, value.Interface())
				if err != nil {
					log.Panic(err)
				}
			}

			err := json.Unmarshal(queryJson, value.Interface())
			if err != nil {
				log.Panic(err)
			}
			parameters = append(parameters, value)
		}
	}

	valid(&parameters)

	of := reflect.ValueOf(u)
	res := of.Call(parameters)
	return res[0].Interface()

}

// 参数校验
func valid(parameters *[]reflect.Value) {
	for _, v := range *parameters {
		t := reflect.TypeOf(v.Interface())

		// 如果是 request 类型，不进行参数校验
		if t.String() == "*http.Request" {
			continue
		}
		if t.Kind() == reflect.Pointer {
			t = t.Elem()
			v = v.Elem()
		}

		// 如果不是结构体，则不进行参数校验
		if t.Kind() != reflect.Struct {
			return
		}

		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			value := v.Field(i)

			// 取结构体字段上标注 valid 的 tag
			valid := field.Tag.Get("valid")
			if valid == "" {
				continue
			}

			// 根据,分割，因为可能有多种类型限制，例如 not nil,len:3
			blocks := strings.Split(valid, ",")
			for _, block := range blocks {
				split := strings.Split(block, ":")
				if len(split) == 2 {
					if handler, ok := validMap[split[0]]; ok {
						handler(field, value, split[1])
					}
				} else if len(split) == 1 {
					if handler, ok := validMap[split[0]]; ok {
						handler(field, value, "")
					}
				}
			}
		}

	}
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
