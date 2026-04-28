package service

import (
	"My_AI_Assistant/internal/model"
	"My_AI_Assistant/internal/repository"
	"encoding/json"
	"fmt"
	"strings"
)

type RAGService struct {
	embeddingService *EmbeddingService                //文本向量化
	chromaService    *ChromaService                   //向量数据库服务
	aiService        *AIService                       //大模型调用服务
	ragHistoryRepo   *repository.RAGHistoryRepository //历史记录仓储
	collectionID     string                           // 与 DocumentService 共用同一个集合的 UUID
}

func NewRAGService( //实例化注入
	embeddingService *EmbeddingService,
	chromaService *ChromaService,
	aiService *AIService,
	ragHistoryRepo *repository.RAGHistoryRepository,
	collectionID string,
) *RAGService {
	return &RAGService{
		embeddingService: embeddingService,
		chromaService:    chromaService,
		aiService:        aiService,
		ragHistoryRepo:   ragHistoryRepo,
		collectionID:     collectionID,
	}
}

// GenerateWithRAG 基于文档的 RAG 生成（优化版）
func (s *RAGService) GenerateWithRAG(userID uint, question string, topK int) (*model.RAGResponse, error) {
	if topK <= 0 { //用户输入需要检索的文档，如果小于0，默认5
		topK = 5
	}

	fmt.Printf(">>> RAG 开始处理问题: %s\n", question)

	// 1. 将问题向量化
	queryEmbedding, err := s.embeddingService.EmbedText(question) //将用户输入的问题向量化
	if err != nil {
		return nil, fmt.Errorf("问题向量化失败: %v", err)
	}
	fmt.Printf(">>> 问题向量化成功，向量维度: %d\n", len(queryEmbedding))

	// 2. 从 Chroma 检索
	fmt.Printf(">>> 正在从集合 %s 检索 top %d 条...\n", s.collectionID, topK)
	relevantDocs, err := s.chromaService.SearchSimilar(s.collectionID, queryEmbedding, topK)
	if err != nil {
		return nil, fmt.Errorf("检索文档失败: %v", err)
	}

	fmt.Printf(">>> 实际检索到 %d 条相关文档\n", len(relevantDocs))

	// 3. 构建上下文
	var contextBuilder strings.Builder //高效拼接字符串
	var references []string            //保存原始参考资料文本

	if len(relevantDocs) == 0 {
		contextBuilder.WriteString("（未检索到相关参考资料）\n\n")
	} else {
		for i, doc := range relevantDocs {
			//遍历结果，提取里面的text字段
			if text, ok := doc["text"].(string); ok && text != "" { //断言
				//将格式化后的字符串写入 builder
				contextBuilder.WriteString(fmt.Sprintf("【参考资料 %d】\n%s\n\n", i+1, text))
				references = append(references, text) //将text内容保存
			}
		}
	}
	//所有的拼接到一起，转换成一个
	context := contextBuilder.String()
	fmt.Printf(">>> 构建的上下文长度: %d 字符\n", len(context))

	// 4. 构建优化后的 Prompt
	//传入用户问题和检索的上下文，构成一个完整的提示词
	prompt := s.buildPrompt(question, context)

	// 5. 调用 DeepSeek 生成答案
	answer, err := s.aiService.GenerateContent(prompt)
	if err != nil {
		return nil, fmt.Errorf("生成答案失败: %v", err)
	}

	// 6. 保存历史记录（异步，不阻塞返回）
	go s.saveHistory(userID, question, answer, references)

	return &model.RAGResponse{
		Answer:     answer,
		References: references,
	}, nil
}

// buildPrompt 构建高质量 Prompt
func (s *RAGService) buildPrompt(question, context string) string {
	return fmt.Sprintf(`【角色定义】
你是一位专业、严谨的知识助手。你必须严格基于提供的参考资料进行回答，不得虚构参考资料中不存在的内容。

【参考资料】
%s

【用户问题】
%s

【回答要求】
1. 先判断参考资料是否能回答问题：
   - 如果能：基于参考资料，用自己的语言组织答案，逻辑清晰，层次分明。
   - 如果不能：明确告知用户"根据现有参考资料无法回答该问题"，不要编造内容。
2. 回答结构：先给出核心结论，再分点详细说明。
3. 如果参考资料部分相关，只使用相关的部分，其余部分明确标注"参考资料中未涉及"。
4. 回答语言：全程使用中文，表述准确、专业，避免口语化。
5. 引用方式：如需引用参考资料中的原文，请用【X】标注来源序号。

请开始回答：`, context, question)
}

// saveHistory 异步保存 RAG 生成历史
func (s *RAGService) saveHistory(userID uint, question, answer string, references []string) {
	refsJSON, _ := json.Marshal(references) //序列化参考资料
	history := &model.RAGHistory{
		UserID:     userID,
		Question:   question,
		Answer:     answer,
		References: string(refsJSON), // // 将 []byte 转为 string
	}
	if err := s.ragHistoryRepo.Create(history, 20); err != nil {
		fmt.Printf(">>> 保存 RAG 历史失败: %v\n", err)
	}
}

// GetHistory 获取用户 RAG 历史记录
func (s *RAGService) GetHistory(userID uint) ([]model.RAGHistory, error) {
	histories, err := s.ragHistoryRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	// 转换 References JSON 字段为数组
	for i := range histories {
		if histories[i].References != "" {
			//转换回[]string
			json.Unmarshal([]byte(histories[i].References), &histories[i].ReferencesArr)
		}
	}

	return histories, nil
}

// GetHistoryDetail 获取单条 RAG 历史详情
func (s *RAGService) GetHistoryDetail(historyID uint, userID uint) (*model.RAGHistory, error) {
	history, err := s.ragHistoryRepo.FindByID(historyID, userID)
	if err != nil {
		return nil, err
	}
	if history.References != "" {
		json.Unmarshal([]byte(history.References), &history.ReferencesArr)
	}
	return history, nil
}

// DeleteHistory 删除单条 RAG 历史
func (s *RAGService) DeleteHistory(historyID uint, userID uint) error {
	return s.ragHistoryRepo.Delete(historyID, userID)
}

// ClearHistory 清空用户所有 RAG 历史
func (s *RAGService) ClearHistory(userID uint) error {
	return s.ragHistoryRepo.ClearAll(userID)
}
