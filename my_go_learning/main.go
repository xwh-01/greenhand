package main

import (
	"fmt"
	"my_go_learning/handlers"
	"net/http"
)

func main() {
	handler := handlers.NewShortURLHandler()

	http.HandleFunc("/create", handler.CreateShortURL)
	http.HandleFunc("/stats", handler.Stats)
	http.HandleFunc("/", handler.Redirect)

	fmt.Println("短链接服务启动成功!")
	fmt.Println("访问 http://localhost:8080/stats 查看统计")
	fmt.Println("示例:curl -X POST -d \"url=https://www.example.com\"http://localhost:8080/create")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("服务器启动失败:%v\n", err)
	}
}
