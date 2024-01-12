package api

import (
	"douyinshibie/cfg"
	"douyinshibie/model"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// randomUserAgent 从userAgents列表中随机选择一个User-Agent字符串并返回
func randomUserAgent() string {
	rand.Seed(time.Now().UnixNano())
	return cfg.Config.App.UserAgents[rand.Intn(len(cfg.Config.App.UserAgents))]
}

func ProcessUserInput(input string) string {
	linkRegex := regexp.MustCompile(`v\.douyin\.com\/[a-zA-Z0-9]+`)
	idRegex := regexp.MustCompile(`\d{19}`)

	if linkRegex.MatchString(input) {
		return linkRegex.FindString(input)
	} else if idRegex.MatchString(input) {
		return idRegex.FindString(input)
	}
	return ""
}

func ExtractVideoId(link string) string {
	// 确保链接包含协议方案
	if !strings.HasPrefix(link, "http://") && !strings.HasPrefix(link, "https://") {
		link = "https://" + link
	}

	// 发送请求并获取重定向后的URL
	resp, err := http.Get(link)
	if err != nil {
		fmt.Println("Error fetching URL:", err)
		return ""
	}
	defer resp.Body.Close()

	// 使用最终请求的URL，可能包含重定向
	finalURL := resp.Request.URL.String()
	log.Print("Final URL: " + finalURL)
	finalURL = resp.Request.URL.String()
	idRegex := regexp.MustCompile(`/video/(\d+)`)
	matches := idRegex.FindStringSubmatch(finalURL)
	if len(matches) > 1 {
		log.Println("Video ID: " + matches[1])
		return matches[1]
	}
	return ""
}

func IsNumeric(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

func GetDouYinVideoInfo(videoId string) (string, string, error) {
	url := fmt.Sprintf("https://www.iesdouyin.com/web/api/v2/aweme/iteminfo/?item_ids=%s&a_bogus=64745b2b5bdc4e75b720a9a85b19867a", videoId)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return "", "", err
	}
	req.Header.Add("User-Agent", randomUserAgent())
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", "", err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return "", "", err
	}
	if res.StatusCode != 200 {
		fmt.Println(string(body))
		return "", "", fmt.Errorf("response status code is not 200")
	}
	//fmt.Println(string(body))

	var videoInfo model.DouYinVideoInfo
	err = json.Unmarshal(body, &videoInfo)
	if err != nil {
		log.Println("Error unmarshalling JSON:", err)
		return "", "", err
	}

	if len(videoInfo.ItemList) > 0 && videoInfo.ItemList[0].Video.PlayAddr.Uri != "" {
		uri := videoInfo.ItemList[0].Video.PlayAddr.Uri
		desc := videoInfo.ItemList[0].Desc
		finalUrl := fmt.Sprintf("https://www.iesdouyin.com/aweme/v1/play/?video_id=%s&ratio=1080p&line=0", uri)
		return finalUrl, desc, nil
	}
	return "", "", fmt.Errorf("no video found")
}

// GetDouYinInfo 获取抖音视频信息
func GetDouYinInfo(link string) (string, string, error) {
	videoIdOrLink := ProcessUserInput(link)
	var videoId string
	if videoIdOrLink != "" {
		if IsNumeric(videoIdOrLink) {
			videoId = videoIdOrLink
		} else {
			videoId = ExtractVideoId(videoIdOrLink)
		}
	}
	if len(videoId) == 0 {
		return "", "", fmt.Errorf("videoId is not found")
	}
	finalUrl, title, err := GetDouYinVideoInfo(videoId)
	if err != nil {
		return "", "", err
	}
	return finalUrl, title, nil
}
