package service

import (
	"My_AI_Assistant/pkg/httpclient"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// 是一个 Go 语言客户端，用于连接和操作 Chroma 向量数据库，实现向量的存储和检索。
// ChromaService 用于与 Chroma 向量数据库交互
type ChromaService struct { //这个结构体代表一个 Chroma 服务客户端
	baseURL string //存储chroma服务器的访问地址
	client  *http.Client
}

// NewChromaService 创建 ChromaService
func NewChromaService(baseURL string) *ChromaService {
	return &ChromaService{
		baseURL: baseURL,
		client:  httpclient.GetExternalAPIClient(), // 使用连接池 Client
	}
}

// AddDocument 将文档向量添加到 Chroma
// collectionName：集合名称（类似于数据库中的表名）                                 文档的向量表示   文档的元数据
func (s *ChromaService) AddDocument(collectionName, documentID string, embedding []float32, metadata map[string]interface{}) error {
	//构建请求URL：http://localhost:8000/api/v1/collections/我的集合/add
	url := fmt.Sprintf("%s/api/v1/collections/%s/add", s.baseURL, collectionName)
	//构建请求体（Chroma API 要求的格式
	payload := map[string]interface{}{
		"ids":        []string{documentID},               //文档ID列表
		"embeddings": [][]float32{embedding},             //向量列表
		"metadatas":  []map[string]interface{}{metadata}, //元数据列表
	}
	// 将 payload 转换成 JSON 格式
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	// 创建 HTTP POST 请求                           bytes.NewBuffer() 将字节数组转换成 io.Reader 接口类型
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	/*http.NewRequest 创建http请求对象
	参数说明：
	method  string类型   http方法
	url     string类型   请求地址
	body	io.Reader 	请求体数据
	*/
	if err != nil {
		return err
	}
	//设置 HTTP 请求头中的 Content-Type 字段
	req.Header.Set("Content-Type", "application/json") //不设置这个，服务器就不知道你发的是什么格式的数据，可能无法正确解析。
	/*req	HTTP 请求对象
	Header	请求头字段的 map
	Set()	设置或覆盖请求头字段的方法
	"Content-Type"	要设置的请求头字段名*/
	resp, err := s.client.Do(req) //实际发送之前创建好的 HTTP 请求
	if err != nil {
		return err
	}
	defer resp.Body.Close() //延迟执行
	//resp.Body.Close() 关闭响应体连接
	/*为什么必须关闭？
	防止内存泄漏
	释放网络连接资源
	如果不关闭，连接可能一直被占用*/
	if resp.StatusCode < 200 || resp.StatusCode >= 300 { //服务器返回的 HTTP 状态码
		/*	200	成功	一切正常
			400	错误请求	客户端发送的数据格式不对
			401	未授权	需要认证
			403	禁止访问	没有权限
			404	未找到	资源不存在
			500	服务器错误	服务器内部出错*/
		body, _ := io.ReadAll(resp.Body) //读取服务器返回的完整数据
		return fmt.Errorf("Chroma 添加失败: %s", string(body))
	}

	return nil
}

// SearchSimilar 检索相似文档（用于 RAG）     集合名称                 查询问题的向量     返回最相似的 K 个结果
func (s *ChromaService) SearchSimilar(collectionName string, queryEmbedding []float32, topK int) ([]map[string]interface{}, error) {
	//构建请求 URL
	url := fmt.Sprintf("%s/api/v1/collections/%s/query", s.baseURL, collectionName)
	//构建请求体（Chroma API 要求的格式
	payload := map[string]interface{}{
		//Chroma API 要求的字段名       二维浮点数数组
		"query_embeddings": [][]float32{queryEmbedding}, //因为Chroma API支持查询多个向量，所以要用二维
		"n_results":        topK,                        //Chroma 返回多少个最相似的文档
	}
	//请求体转换为json字符串
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	// 创建 HTTP POST 请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	//设置 HTTP 请求头中的 Content-Type 字段
	req.Header.Set("Content-Type", "application/json")
	//发送请求
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() //延迟关闭响应体

	body, err := io.ReadAll(resp.Body) //读取响应体
	if err != nil {
		return nil, err
	}

	// ── 调试：打印 Chroma 返回的状态码和响应体 ──
	fmt.Printf(">>> Chroma query 状态码: %d\n", resp.StatusCode)
	fmt.Printf(">>> Chroma query 响应: %s\n", string(body))

	//检查 Chroma 是否返回了错误
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("Chroma 查询失败 [HTTP %d]: %s", resp.StatusCode, string(body))
	}
	//声明匿名结构体 //用来接收上面读取的响应体数据
	var result struct {
		Ids       [][]string                 `json:"ids"`
		Distances [][]float64                `json:"distances"`
		Metadatas [][]map[string]interface{} `json:"metadatas"`
	}
	//将 JSON 格式的字节数组解析成 Go 结构体 并放入结构体里
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	//提取第一个查询的元数据列表
	//检查是否有结果						返回第一个查询的结果
	if len(result.Metadatas) > 0 {
		return result.Metadatas[0], nil
	}
	//没有结果时返回空切片
	return []map[string]interface{}{}, nil
}

// GetOrCreateCollection 获取或创建集合，返回集合 ID (UUID)
func (s *ChromaService) GetOrCreateCollection(name string) (string, error) {
	// 先尝试获取集合
	//构建请求 URL
	url := fmt.Sprintf("%s/api/v1/collections/%s", s.baseURL, name)
	//创建 HTTP GET 请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	//发送请求
	resp, err := s.client.Do(req)
	if err != nil {
		return "", err
	}
	//延迟关闭响应体
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		//如果集合存在，定义临时结构体 获取id字段
		var result struct {
			ID string `json:"id"`
		}
		//读取响应体
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		if err := json.Unmarshal(body, &result); err != nil {
			return "", err
		}
		return result.ID, nil //找到集合，将id取出来放弃结构体中，返回出去
	}

	// 集合不存在，创建新集合
	//构建创建集合的 URL
	url = fmt.Sprintf("%s/api/v1/collections", s.baseURL)
	payload := map[string]interface{}{
		"name": name,
	}
	//将请求体转为json格式
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	//创建POST请求
	req, err = http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	//设置请求头
	req.Header.Set("Content-Type", "application/json")
	//发送请求
	resp, err = s.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close() //延迟关闭响应体

	body, err := io.ReadAll(resp.Body) //读取响应体
	if err != nil {
		return "", err
	}

	var result struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	return result.ID, nil
}
