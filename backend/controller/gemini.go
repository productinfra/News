package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 请求结构
type ChatRequest struct {
	Message string `json:"message"`
}

// Gemini API 请求体结构
type GeminiAPIRequest struct {
	Contents []struct {
		Parts []struct {
			Text string `json:"text"`
		} `json:"parts"`
	} `json:"contents"`
}

// Gemini API 响应结构（只关注需要的字段）
type GeminiAPIResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

// Gemini 处理函数
func Gemini(c *gin.Context) {
	var req ChatRequest

	// 绑定请求体
	if err := c.ShouldBindJSON(&req); err != nil || req.Message == "" {
		log.Printf("请求解析失败: %v", err) // 记录错误
		c.JSON(http.StatusBadRequest, gin.H{"error": "Message is required"})
		return
	}

	// 构造 Gemini 请求体
	requestBody := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]string{
					{"text": req.Message},
				},
			},
		},
	}

	// 将请求体序列化为 JSON 字符串
	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		log.Printf("请求体序列化失败: %v", err) // 记录错误
		c.JSON(http.StatusInternalServerError, gin.H{"error": "请求体序列化失败"})
		return
	}

	// 设置 Gemini API
	const GEMINI_API_KEY = "AIzaSyASFi2GlRG72jYFCqoXltB6z0qtSlaQnnA"
	apiURL := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent?key=%s", GEMINI_API_KEY)

	// 发起 POST 请求
	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(bodyBytes))
	if err != nil {
		log.Printf("Gemini API 网络错误: %v", err) // 记录错误
		c.JSON(http.StatusInternalServerError, gin.H{"error": "请求 Gemini 失败"})
		return
	}
	defer resp.Body.Close()

	// 读取响应体
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("读取响应体失败: %v", err) // 记录错误
		c.JSON(http.StatusInternalServerError, gin.H{"error": "读取响应体失败"})
		return
	}

	// 解析 Gemini 响应
	var geminiResp GeminiAPIResponse
	if err := json.Unmarshal(respBody, &geminiResp); err != nil {
		log.Printf("Gemini 响应解析失败: %v", err) // 记录错误
		c.JSON(http.StatusInternalServerError, gin.H{"error": "响应解析失败"})
		return
	}

	// 提取 Gemini 回复
	var reply string
	if len(geminiResp.Candidates) > 0 &&
		len(geminiResp.Candidates[0].Content.Parts) > 0 {
		reply = geminiResp.Candidates[0].Content.Parts[0].Text
	} else {
		reply = "Gemini 没有返回结果"
	}

	// 返回回复
	c.JSON(http.StatusOK, gin.H{"reply": reply})
}
