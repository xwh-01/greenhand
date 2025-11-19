package handlers

import (
	"fmt"
	"math/rand"
	"my_go_learning/models"
	"my_go_learning/storage"
	"net/http"
	"time"
)

type ShortURLHandler struct {
	storage *storage.MySQLStorage
}

// NewShortURLHandler 创建使用MySQL存储的处理器
func NewShortURLHandler(dsn string) (*ShortURLHandler, error) {
	mysqlStorage, err := storage.NewMySQLStorage(dsn)
	if err != nil {
		return nil, err
	}
	return &ShortURLHandler{storage: mysqlStorage}, nil
}

// // NewShortURLHandlerMemory 创建使用内存存储的处理器
// func NewShortURLHandlerMemory() *ShortURLHandler {
// 	memoryStorage := storage.NewMemoryStorage()
// 	return &ShortURLHandler{storage: memoryStorage}
// }

func generateShortCode() string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const codeLength = 6

	rand.Seed(time.Now().UnixNano())
	code := make([]byte, codeLength)
	for i := range code {
		code[i] = chars[rand.Intn(len(chars))]
	}
	return string(code)
}

func (h *ShortURLHandler) CreateShortURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "只支持POST请求", http.StatusMethodNotAllowed)
		return
	}

	longURL := r.FormValue("url")
	if longURL == "" {
		http.Error(w, "URL不能为空", http.StatusBadRequest)
		return
	}

	// shortCode := generateShortCode()
	// shortURL := models.NewShortURL(shortCode, longURL)
	// 生成唯一短码
	var shortCode string
	for i := 0; i < 5; i++ {
		shortCode = generateShortCode()
		existing, _ := h.storage.FindByShortCode(shortCode)
		if existing == nil {
			break
		}
	}

	shortURL := models.NewShortURL(shortCode, longURL)

	err := h.storage.Save(shortURL)
	if err != nil {
		http.Error(w, "保存失败", http.StatusInternalServerError)
		return
	}

	result := fmt.Sprintf("短链接创建成功!\n短码:%s\n长链接:%s\n访问地址:http://localhost:8080/%s", shortCode, longURL, shortCode)

	w.Header().Set("Content-Type", "text/plain;charset=utf-8")
	//w.WriteHeader(http.StatusCreated)
	w.Write([]byte(result))

}

func (h *ShortURLHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	shortCode := r.URL.Path[1:]

	if shortCode == "" && r.URL.Path == "/" {
		welcome := `欢迎使用短链接服务 V2.0！创建短链接: POST http://localhost:8080/create查看统计: GET http://localhost:8080/stats访问短链接: GET http://localhost:8080/{短码}`

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte(welcome))
		//http.Error(w, "短码不能为空", http.StatusBadRequest)
		return
	}

	shortURL, err := h.storage.FindByShortCode(shortCode)
	if err != nil {
		http.Error(w, "查找失败", http.StatusInternalServerError)
		return
	}
	if shortURL == nil {
		http.Error(w, "短链接不存在", http.StatusNotFound)
		return
	}

	h.storage.IncrementClick(shortCode)
	http.Redirect(w, r, shortURL.LongURL, http.StatusFound)
	//fmt.Printf("短码 %s 被访问,跳转到: %s (总点击: %d)\n", shortCode, shortURL.LongURL, shortURL.ClickCount)

}

func (h *ShortURLHandler) Stats(w http.ResponseWriter, r *http.Request) {
	allURLs, err := h.storage.GetAll()
	if err != nil {
		http.Error(w, "获取数据失败", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/plain;charset=utf-8")

	if len(allURLs) == 0 {
		w.Write([]byte("还没有创建任何短链接\n"))
		return
	}

	result := "短链接统计信息:\n\n"
	totalClicks := 0
	for _, item := range allURLs {
		result += fmt.Sprintf("短码: %s\n长链接: %s\n点击次数: %d\n创建时间: %s\n\n",
			item.ShortCode, item.LongURL, item.ClickCount,
			item.CreatedAt.Format("2006-01-02 15:04:05"))
		totalClicks += item.ClickCount
	}
	result += fmt.Sprintf("总计: %d 个短链接, %d 次点击\n", len(allURLs), totalClicks)
	w.Write([]byte(result))
}
