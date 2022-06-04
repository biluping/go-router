package router

import (
	"log"
	"net/http"
	"strconv"
)

// 初始化方法，注册全局处理器，类似于 java 中的 DispatcherServlet 的作用
func init() {
	http.HandleFunc("/", routerHandle)
}

// Get 添加 Get 方法处理器
func Get(path string, handle Controller) {
	register(http.MethodGet, path, handle)
}

func Post(path string, handle Controller) {
	register(http.MethodPost, path, handle)
}

func Put(path string, handle Controller) {
	register(http.MethodPut, path, handle)
}

func Delete(path string, handle Controller) {
	register(http.MethodDelete, path, handle)
}

// AddFilter 添加请求过滤器
func AddFilter(filter Filter) {
	filterList = append(filterList, filter)
}

func Start(port int) {
	log.Printf("http server start successful, port %d\n", port)
	err := http.ListenAndServe(":"+strconv.Itoa(port), nil)
	if err != nil {
		log.Println(err)
		return
	}
}
