package service

import (
	"context"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

type AIService struct { //ai服务结构体
	//变量名：客户端  类型： openai 库提供的 AI 客户端
	client *openai.Client //连接AI的客户端
	//拿着它，才能调用AI接口，发请求，拿回答
}

func NewAIService(apiKey string) *AIService { //创建并返回一个可用的ai服务
	config := openai.DefaultConfig(apiKey)      //传入 API key，
	config.BaseURL = "https://api.deepseek.com" //DeepSeek 官方地址这样才可以使用模型

	return &AIService{
		//拿着配置去创建一个 AI 客户端，
		client: openai.NewClientWithConfig(config),
	}
}

// GenerateContent 根据主题生成文章
func (s *AIService) GenerateContent(topic string) (string, error) {
	resp, err := s.client.CreateChatCompletion( //通过AI 客户端，发送请求给 DeepSeek，让 AI 开始写东西。
		context.Background(), //空的上下文
		//控制请求的生命周期，ai请求现在开始，必须传入
		openai.ChatCompletionRequest{ //
			//用来包装 “发给 AI 的所有请求内容”
			/*结构体内容{
				使用模型，
				发什么消息，
				创意程度

			}
			*/
			Model: "deepseek-v4-flash", //使用的模型DeepSeek 官方的聊天模型。
			Messages: []openai.ChatCompletionMessage{ //发给ai的话，一组话
				{
					Role: openai.ChatMessageRoleSystem, //给ai设定角色
					Content: "你是一位精通中文表达的资深技术文档与内容创作者，" +
						"擅长撰写逻辑严谨、结构清晰、语言精炼且专业度高的文章。" +
						"写作风格要求：\n" +
						"1.  用词精准、表述客观，避免口语化与模糊表达；\n" +
						"2.  结构层次分明，善用标题、分点与逻辑连接词；\n" +
						"3.  内容基于事实与行业标准，兼具实用性与权威性；\n" +
						"4.  优先使用专业术语，必要时提供通俗解释；\n" +
						"5.  输出内容完整、自洽，不冗余、不跑题，直接回应需求。", //设定角色的具体内容
				},
				{
					Role: openai.ChatMessageRoleUser, //用户提的需求
					Content: fmt.Sprintf("请根据以下主题，生成一篇高质量的中文文章。\n"+
						"主题：《%s》\n"+
						"要求：\n"+
						"1.  字数控制在800字左右，正负误差不超过100字；\n"+
						"2.  结构清晰，建议包含引言、核心内容、总结三部分，可使用小标题分层；\n"+
						"3.  语言正式专业，避免口语化表达、网络热词和冗余描述；\n"+
						"4.  内容需逻辑严谨，观点明确，基于事实或行业通用认知；\n"+
						"5.  输出格式为纯文本，使用标准中文标点，无需Markdown格式。", topic),
				},
			},
			Temperature: 0.7, //创意程度
			/*0 = 很老实、很死板
			1 = 很会编、很放飞
			0.7 = 写文章最完美的数值*/
		},
	)

	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
	/*
			resp:  AI 返回给你的整个响应结果（一大包数据）\
			resp.Choices:  AI 给你的回答列表
			[0]下标为0，就是第一个回答
			.Message : 这个回答里的消息内容
			.Content : 只取Content后面的内容

		AI返回的第一个回答的消息
		Message: {
		    Role: "assistant"
		    Content: "这是一篇关于Go语言的文章..."
		}
	*/
}
