package service

import (
	"My_AI_Assistant/pkg/httpclient"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// EmbeddingService 向量嵌入服务（使用通义千问）
type EmbeddingService struct {
	apiKey string
	client *http.Client
}

// NewEmbeddingService 创建 EmbeddingService
func NewEmbeddingService(apiKey string) *EmbeddingService {
	return &EmbeddingService{
		apiKey: apiKey,
		client: httpclient.GetExternalAPIClient(),
	}
}

// EmbedText 将文本转换为向量（1536维）
func (s *EmbeddingService) EmbedText(text string) ([]float32, error) {
	url := "https://dashscope.aliyuncs.com/api/v1/services/embeddings/text-embedding/text-embedding"
	//因为阿里云准换回来得就是这种样式，所以把map写成这样样式方便传入
	payload := map[string]interface{}{
		"model": "text-embedding-v2", //模型模板
		"input": map[string]interface{}{
			"texts": []string{text}, //要转换的文本
		},
		"parameters": map[string]interface{}{
			"dimension": 1536, //输出向量维度
		},
	}

	jsonData, _ := json.Marshal(payload) //转为json字符串
	//创建 HTTP POST 请求对象
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	//                     HTTP 标准认证头  认证方式  API 密钥
	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	req.Header.Set("Content-Type", "application/json") //告诉请求是json格式
	resp, err := s.client.Do(req)                      //发送请求
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() //延迟关闭响应体

	body, err := io.ReadAll(resp.Body) // 从响应体中读取所有字节数据
	if err != nil {
		return nil, err
	}
	// API 返回的原始 JSON
	var result struct { //精确匹配 API 返回格式
		Output struct { //只提取需要的字段
			Embeddings []struct {
				Embedding []float32 `json:"embedding"`
			} `json:"embeddings"`
		} `json:"output"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析 Embedding 返回失败: %v", err)
	}

	if len(result.Output.Embeddings) == 0 {
		return nil, fmt.Errorf("Embedding 返回为空")
	}

	return result.Output.Embeddings[0].Embedding, nil
}
