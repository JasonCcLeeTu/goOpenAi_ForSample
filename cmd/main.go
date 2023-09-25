package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
)

func main() {

	godotenv.Load("../.env.local") //讀取env檔案
	openaiAPIKey := os.Getenv("OPENAI_API_KEY")
	if openaiAPIKey == "" {
		log.Printf("openai key is null")
	}

	openaiClient := openai.NewClient(openaiAPIKey) //創建 openai client端

	response, err := ChatRequest(openaiClient, "'I'm going to  steak tomorrow  in two weeks.' Which part of this sentence expresses time  ? Answer  only.") //將sentence中的時間截取出來
	if err != nil {
		log.Fatalf("ChatRequest error:%s", err.Error())
	}

	for index, content := range response.Choices { //回覆的訊息
		log.Printf("時間內容_回復_%d:%s\n", index, content.Message.Content)
	}

	timeFormat := "2006-01-02 15:04:05" //時間格式

	nowTime := time.Now().Format(timeFormat)
	log.Println("目前時間:", nowTime)
	convertTime := fmt.Sprintf("%s %s", nowTime, response.Choices[0].Message.Content) //將現在時間加上 回覆截取出來的時間
	answer := fmt.Sprintf("'%s' Convert Time. Answer  only.", convertTime)            //跟ChatGPT講 將前面組的句子 轉換成時間
	//answer := "'I'm going to  steak tomorrow  in two weeks.' Which part of this sentence expresses time  ? Answer  only."
	response, err = ChatRequest(openaiClient, answer)
	if err != nil {
		log.Fatalf("ChatRequest2 error:%s", err.Error())
	}

	log.Println("GPT-3.5回復:", response.Choices[0].Message.Content)

}

func ChatRequest(client *openai.Client, content string) (*openai.ChatCompletionResponse, error) {

	var requestMsg []openai.ChatCompletionMessage
	// chatContent := openai.ChatCompletionMessage{ //中文
	// 	Role:    openai.ChatMessageRoleUser,
	// 	Content: "'我三個月後要去看今天餐廳' 這句話哪段是指真的時間 只回答答案",
	// }

	chatContent := openai.ChatCompletionMessage{ //英文
		Role:    openai.ChatMessageRoleUser,
		Content: content,
	}

	requestMsg = append(requestMsg, chatContent)

	chatRequest := openai.ChatCompletionRequest{ //設置請求條件和內容
		Model:    "gpt-3.5-turbo",
		Messages: requestMsg,
	}

	ctx := context.Background()
	response, err := client.CreateChatCompletion(ctx, chatRequest) //發送請求
	if err != nil {
		return nil, fmt.Errorf("chat response error:%s", err.Error())
	}

	return &response, nil
}
