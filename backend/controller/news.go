package controller

import (
	"crypto/tls"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const NEWS_API_URL = "https://newsapi.org/v2/top-headlines?country=us&apiKey=804f6bdac084463ba3fadb53f9efce90"

// News 负责从 NewsAPI 获取头条新闻并返回给前端
func News(c *gin.Context) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Get(NEWS_API_URL)
	if err != nil {
		log.Println("请求 NewsAPI 失败:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取新闻失败"})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("读取 NewsAPI 响应失败:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "读取响应失败"})
		return
	}

	c.Data(http.StatusOK, "application/json", body)
}
