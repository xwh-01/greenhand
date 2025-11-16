package handlers

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"my_go_learning/models"
	"my_go_learning/storage"
)

type ShortURLHandler struct {
	storage *storage.MemoryStorage
}

func NewShortURLHandler() *ShortURLHandler {
	return &ShortURLHandler{
		storage: storage.NewMemoryStorage(),
	}
}

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

	shortCode := generateShortCode()
	shortURL := models.NewShortURL(shortCode, longURL)

	err := h.storage.Save(shortURL)
	if err != nil {
		http.Error(w, "保存失败", http.StatusInternalServerError)
		return
	}

	result := fmt.Sprintf("短链接创建成功!\n短码:%s\n长链接:%s\n访问地址:http://localhost:8080%s", shortCode, longURL, shortCode)

	w.Header().Set("Content-Type", "text/plain;charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(result))

}

func (h *ShortURLHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	shortCode := r.URL.Path[1:]

	if shortCode == "" {
		http.Error(w, "短码不能为空", http.StatusBadRequest)
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

	shortURL.IncrementClick()

	http.Redirect(w, r, shortURL.LongURL, http.StatusFound)

	fmt.Printf("短码 %s 被访问,跳转到: %s (总点击: %d)\n", shortCode, shortURL.LongURL, shortURL.ClickCount)

}

func (h *ShortURLHandler) Stats(w http.ResponseWriter, r *http.Request) {
	allURLs := h.storage.GetAll()

	w.Header().Set("Content-Type", "text/plain;charset=utf-8")

	if len(allURLs) == 0 {
		w.Write([]byte("还没有创建任何短链接\n"))
		return
	}

	result := "短链接统计信息:\n\n"
	for _, shortURL := range allURLs {
		result += fmt.Sprintf("短码:%s\n长链接:%s\n点击次数:%d\n创建时间:%s\n\n", shortURL.ShortCode, shortURL.LongURL, shortURL.ClickCount, shortURL.CreatedAt.Format("2006-01-02 15:04:05"))

	}
	w.Write([]byte(result))
}
