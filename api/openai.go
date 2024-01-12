package api

import (
	"bytes"
	"douyinshibie/cfg"
	"douyinshibie/model"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func SetGptV(data model.FrameInfo) error {
	// 确保Base64数据非空
	if data.Base64Data == "" {
		return fmt.Errorf("base64 data is empty")
	}
	//fmt.Printf("data:image/jpeg;base64,%s\r\n", data.Base64Data)
	url := cfg.Config.App.OpenaiUrl + "/v1/chat/completions"
	method := "POST"

	payload := model.RequestPayload{
		Stream: false,
		Model:  "gpt-4-vision-preview",
		Messages: []model.Message{
			{
				Role: "user",
				Content: []model.Content{
					{
						Type: "text",
						Text: fmt.Sprintf("这个图片是一段视频中第%d片段的第%d帧，他的详细内容内容是什么？比如有什么人物，他们在做什么动作，说什么话。这个时候你就是一个视频脚本分析大师，你应该剖析他们原本的剧情或者画面呈现的东西，你应该直接输出告诉我视频这一帧呈现的内容，现在请开始你分析：", data.SegmentIndex, data.FrameIndex),
					},
					{
						Type: "image_url",
						ImageURL: &model.ImageURL{
							URL: "data:image/jpeg;base64," + data.Base64Data,
						},
					},
				},
			},
		},
		MaxTokens: 300,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshaling payload: %v", err)
	}
	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewReader(payloadBytes))

	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Add("Authorization", "Bearer "+RandKey(cfg.Config.App.OpenaiKey))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Connection", "keep-alive")

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %v", err)
	}
	if res.StatusCode != 200 {
		return fmt.Errorf("status code error: %d %s", res.StatusCode, string(body))
	}
	var d model.GptResponse
	if err = json.Unmarshal(body, &d); err != nil {
		return fmt.Errorf("error unmarshal response body: %v", err)
	}
	//fmt.Println(d.Choices[0].Message.Content)
	if len(d.Choices[0].Message.Content) == 0 {
		return fmt.Errorf("content slice is empty")
	}

	frameDescription := fmt.Sprintf("片段 %d中的第 %d帧的内容是%s",
		data.SegmentIndex,
		data.FrameIndex,
		d.Choices[0].Message.Content,
	)
	// 使用互斥锁来保护对共享切片的写入
	Mu.Lock()
	FrameDescriptions = append(FrameDescriptions, frameDescription)
	Mu.Unlock()

	return nil
}

func SetGpt(content string) (error, string) {
	url := cfg.Config.App.OpenaiUrl + "/v1/chat/completions"
	method := "POST"
	//contentStr := strings.TrimSuffix(content, "\n```")
	contentStr := strings.ReplaceAll(content, "\n", "")
	contentStr = strings.ReplaceAll(contentStr, "\\", "")
	contentStr = strings.ReplaceAll(contentStr, `"`, `\"`)

	payload := model.GPTRequest{
		Stream: false,
		Model:  "gpt-4-0613",
		Messages: []model.GPTMessage{
			{
				Role:    "system",
				Content: "你现在是一个视频脚本整合大师。你的任务是将一系列乱序的视频片段整合成一个完整的故事。每个片段都包含了一系列帧的详细内容描述。由于这些片段是乱序的，你需要先找到第 0 片段的第 0 帧，这是视频的开头。从那里开始，确定每个片段及其帧的正确顺序，然后按照这个顺序来分析整个视频的内容。你的最终目标是输出一个连贯的视频内容脚本，该脚本详细地叙述了视频的全部故事线，包括所有关键的对话、场景和情感转变。请注意，你不需要输出处理的过程，只需要提供视频的完整内容概要。现在开始，请查看以下视频片段及其内容描述，并根据这些信息，从第 0 片段的第 0 帧开始，创建一个完整的视频内容脚本，最后请输出一个完成或一个大致的视频内容信息",
			},
			{
				Role:    "user",
				Content: "视频片段解析内容：" + contentStr,
			},
		},
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshaling payload: %v", err), ""
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewReader(payloadBytes))

	if err != nil {
		return fmt.Errorf("error creating request: %v", err), ""
	}
	req.Header.Add("Authorization", "Bearer "+RandKey(cfg.Config.App.OpenaiKey))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "*/*")
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
	var d model.GptResponse
	if err = json.Unmarshal(body, &d); err != nil {
		return fmt.Errorf("error unmarshal response body: %v", err), ""
	}
	if len(d.Choices[0].Message.Content) == 0 {
		return fmt.Errorf("content slice is empty"), ""
	}
	return nil, d.Choices[0].Message.Content
}
