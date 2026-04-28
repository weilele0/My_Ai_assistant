package service

import (
	"My_AI_Assistant/internal/model"
	"My_AI_Assistant/internal/repository"
	"fmt"
)

// DocumentService 文档业务逻辑层
type DocumentService struct {
	documentRepo     *repository.DocumentRepository
	embeddingService *EmbeddingService
	chromaService    *ChromaService
	collectionID     string // 存储集合 UUID
} //把repository注入进来，不自己创建db，调用repository中的db进行操作

// NewDocumentService 创建 DocumentService
func NewDocumentService(
	documentRepo *repository.DocumentRepository,
	embeddingService *EmbeddingService,
	chromaService *ChromaService,
) *DocumentService {
	//调用这个方法 获取或创建名为 "documents_v2" 的集合
	collectionID, err := chromaService.GetOrCreateCollection("my_knowledge_base")
	if err != nil {
		fmt.Printf("初始化 Chroma 集合失败: %v\n", err)
		// 如果失败，使用空字符串，后续 AddDocument 会报错
		collectionID = ""
	} else {
		fmt.Printf("Chroma 集合 UUID: %s\n", collectionID)
	}

	return &DocumentService{
		documentRepo:     documentRepo,
		embeddingService: embeddingService,
		chromaService:    chromaService,
		collectionID:     collectionID, // ✅ 保存 UUID
	}
} //创建 service 实例 传入 repository

// UploadDocument 上传文档（业务逻辑）
func (s *DocumentService) UploadDocument(userID uint, title, filePath, content string) (*model.Document, error) {
	doc := &model.Document{
		UserID:   userID,
		Title:    title,
		FilePath: filePath,
		Content:  content,
	} //组装 Document 结构体
	err := s.documentRepo.Create(doc) //调用数据库创建文档
	if err != nil {
		return nil, err
	}
	go s.vectorizeDocument(doc.ID, content)
	return doc, nil
}

// vectorizeDocument 异步向量化文档并存入 Chroma
func (s *DocumentService) vectorizeDocument(documentID uint, content string) {
	// 1. 文本分块（简单按 500 字符分块）
	chunks := chunkText(content, 500) //AI 模型有输入长度限制（如一次只能处理 512 个字符）
	//遍历每一块
	if s.collectionID == "" {
		fmt.Printf("文档 %d 存入 Chroma 失败: 集合未初始化\n", documentID)
		return
	}
	for i, chunk := range chunks {
		// 2. 调用 Embedding API 生成向量， 调用 AI 模型生成向量
		embedding, err := s.embeddingService.EmbedText(chunk)
		if err != nil {
			fmt.Printf("文档 %d 向量化失败: %v\n", documentID, err)
			continue
		}

		// 3. 存入 Chroma                   文章的id——块索引，生成唯一的分块iD
		docID := fmt.Sprintf("%d_%d", documentID, i) //// 不同文档的块不会冲突
		metadata := map[string]interface{}{
			//元数据
			"document_id": documentID, //原始文档ID
			"chunk_index": i,          //块在文档中的位置
			"text":        chunk,      //块的原始文本内容
		}
		//将向量化的文档块存入 Chroma 数据库            Chroma 集合名称     文档块唯一ID  文本的向量表示  附加信息（元数据
		err = s.chromaService.AddDocument(s.collectionID, docID, embedding, metadata)
		if err != nil {
			fmt.Printf("文档 %d 存入 Chroma 失败: %v\n", documentID, err)
		}
	}

	fmt.Printf(">>> 文档 %d 向量化完成\n", documentID)
}

// chunkText 简单文本分块
//
//	原始文本       每块的大小
func chunkText(text string, chunkSize int) []string {
	var chunks []string
	runes := []rune(text) //因为rune能正确处理中文

	for i := 0; i < len(runes); i += chunkSize {
		end := i + chunkSize
		if end > len(runes) {
			end = len(runes)
		} //调整结束位置长度，如果end为1200 文章长度为1000，将end改为1000，不会浪费内存
		chunks = append(chunks, string(runes[i:end]))
	}
	return chunks
}

// GetUserDocuments 获取当前用户的所有文档
func (s *DocumentService) GetUserDocuments(userID uint) ([]model.Document, error) {
	return s.documentRepo.FindByUserID(userID)
}

// GetDocumentDetail 获取文档详情（带权限校验）                       //返回指针能表示没查到，nil 如果返回值类型 是一个空结构体
func (s *DocumentService) GetDocumentDetail(id uint, userID uint) (*model.Document, error) {
	return s.documentRepo.FindByID(id, userID)
}

// DeleteDocument 删除文档（带权限校验）
func (s *DocumentService) DeleteDocument(id uint, userID uint) error {
	return s.documentRepo.Delete(id, userID)
}
