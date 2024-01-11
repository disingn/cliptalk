package api

import (
	"bytes"
	"douyinshibie/cfg"
	"douyinshibie/model"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
	"time"
)

func randKey() string {
	rand.Seed(time.Now().UnixNano())
	return cfg.Config.App.GeminiKey[rand.Intn(len(cfg.Config.App.GeminiKey))]
}

func NewClient() *http.Client {

	if cfg.Config.Proxy.Protocol != "" {
		// 设置代理地址
		proxyURL, err := url.Parse(cfg.Config.Proxy.Protocol)
		if err != nil {
			log.Println("设置代理出错:", err)
			log.Println("用默认直连")
			client := &http.Client{}
			return client
		}

		log.Println("gemimni 用代理\n代理地址:", proxyURL)

		// 创建一个自定义的 Transport
		transport := &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}

		// 使用自定义的 Transport 创建一个 http.Client
		client := &http.Client{
			Transport: transport,
		}
		return client

	} else {
		log.Println("直连所有")
		client := &http.Client{}
		return client
	}
}

var FrameDescriptions []string
var mu sync.Mutex // 用于保护frameDescriptions切片的互斥锁

func SetGeminiV(data model.FrameInfo) error {
	// 确保Base64数据非空
	if data.Base64Data == "" {
		return fmt.Errorf("base64 data is empty")
	}

	_url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-pro-vision:generateContent?key=%s", randKey())
	method := "POST"

	payload := model.GeminiData{
		Contents: []model.Contents{
			{
				Parts: []model.Parts{
					{
						Text: fmt.Sprintf("这个图片是一段视频中第%d的第%d帧，他的详细内容内容是什么？比如有什么人物，他们在做什么动作，说什么话。这个时候你就是一个视频脚本分析大师，你应该剖析他们原本的剧情或者画面呈现的东西，你应该直接输出告诉我视频这一帧呈现的内容，现在请开始你分析：", data.SegmentIndex, data.FrameIndex),
					},
					{
						InlineData: &model.InlineData{
							MimeType: "image/jpeg",
							Data:     data.Base64Data,
						},
					},
				},
			},
		},
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshaling payload: %v", err)
	}

	//client := &http.Client{}
	client := NewClient()

	req, err := http.NewRequest(method, _url, bytes.NewReader(payloadBytes))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	// 添加必要的请求头
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %v", err)
	}

	// 确保响应状态码为200
	if res.StatusCode != 200 {
		return fmt.Errorf("status code error: %d %s", res.StatusCode, string(body))
	}

	// 解析响应体
	var geminiResponse model.GeminiResponse
	err = json.Unmarshal(body, &geminiResponse)
	if err != nil {
		return fmt.Errorf("error unmarshalling JSON: %v", err)
	}

	// 检查Candidates切片是否非空
	if len(geminiResponse.Candidates) == 0 {
		return fmt.Errorf("candidates slice is empty")
	}

	// 检查Parts切片是否非空
	if len(geminiResponse.Candidates[0].Content.Parts) == 0 {
		return fmt.Errorf("parts slice is empty")
	}

	frameDescription := fmt.Sprintf("片段 %d中的第 %d帧的内容是%s",
		data.SegmentIndex,
		data.FrameIndex,
		geminiResponse.Candidates[0].Content.Parts[0].Text)
	// 使用互斥锁来保护对共享切片的写入
	mu.Lock()
	FrameDescriptions = append(FrameDescriptions, frameDescription)
	mu.Unlock()
	// 输出响应内容
	//log.Printf("片段 %d中的第 %d帧的内容是%s",
	//	data.SegmentIndex,
	//	data.FrameIndex,
	//	geminiResponse.Candidates[0].Content.Parts[0].Text)

	return nil
}

func SetGemini(content string) (error, string) {
	_url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-pro:generateContent?key=%s", randKey())
	method := "POST"

	payload := model.GeminiPro{
		Contents: []model.GeminiProContent{
			{
				Role: "USER",
				Parts: []model.GeminiProPart{
					{
						Text: fmt.Sprintf("你现在是一个视频脚本整合大师。你的任务是将一系列乱序的视频片段整合成一个完整的故事。每个片段都包含了一系列帧的详细内容描述。由于这些片段是乱序的，你需要先找到第 0 片段的第 0 帧，这是视频的开头。从那里开始，确定每个片段及其帧的正确顺序，然后按照这个顺序来分析整个视频的内容。你的最终目标是输出一个连贯的视频内容脚本，该脚本详细地叙述了视频的全部故事线，包括所有关键的对话、场景和情感转变。请注意，你不需要输出处理的过程，只需要提供视频的完整内容概要。现在开始，请查看以下视频片段及其内容描述，并根据这些信息，从第 0 片段的第 0 帧开始，创建一个完整的视频内容脚本：'%s'", content),
					},
				},
			},
		},
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshaling payload: %v", err), ""
	}

	// client := &http.Client{}
	client := NewClient()

	req, err := http.NewRequest(method, _url, bytes.NewReader(payloadBytes))

	if err != nil {
		return fmt.Errorf("error creating request: %v", err), ""
	}
	req.Header.Add("User-Agent", "Apifox/1.0.0 (https://apifox.com)")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Host", "generativelanguage.googleapis.com")
	req.Header.Add("Connection", "keep-alive")

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %v", err), ""
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %v", err), ""
	}
	if res.StatusCode != 200 {
		return fmt.Errorf("status code error: %d %s", res.StatusCode, string(body)), ""
	}
	var geminiResponse model.GeminiResponse
	if err = json.Unmarshal(body, &geminiResponse); err != nil {
		return fmt.Errorf("error unmarshalling JSON: %v", err), ""
	}
	if len(geminiResponse.Candidates) == 0 {
		return fmt.Errorf("candidates slice is empty"), ""
	}
	if len(geminiResponse.Candidates[0].Content.Parts) == 0 {
		return fmt.Errorf("parts slice is empty"), ""
	}

	return nil, geminiResponse.Candidates[0].Content.Parts[0].Text
}
