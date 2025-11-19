package main

import (
	"fmt"
	"my_go_learning/handlers"
	"net/http"
)

func main() {

	dsn := "root:123456@tcp(127.0.0.1:3306)/short_url?charset=utf8mb4&parseTime=True&loc=Local"

	// 创建处理器实例
	handler, err := handlers.NewShortURLHandler(dsn)
	if err != nil {
		fmt.Printf("❌ 初始化失败: %v\n", err)
		fmt.Println("请检查：")
		fmt.Println("1. MySQL服务是否启动")
		fmt.Println("2. 数据库连接信息是否正确")
		fmt.Println("3. 数据库和表是否创建")
		return
	}
	http.HandleFunc("/create", handler.CreateShortURL)
	http.HandleFunc("/stats", handler.Stats)
	http.HandleFunc("/", handler.Redirect)

	fmt.Println("短链接服务启动成功!")
	fmt.Println("访问 http://localhost:8080/stats 查看统计")
	fmt.Println("示例:curl -X POST -d \"url=https://www.example.com\"http://localhost:8080/create")

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("服务器启动失败:%v\n", err)
	}

}
