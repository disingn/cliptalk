package api

import (
	"douyinshibie/model"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
)

func GetTikTokId(link string) string {
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

func GetTikTokVideoData(id string) (string, string, error) {
	apiURL := fmt.Sprintf("https://api16-normal-c-useast1a.tiktokv.com/aweme/v1/feed/?aweme_id=%s", id)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, apiURL, nil)
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
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", "", fmt.Errorf("error reading response body: %v", err)
	}
	if res.StatusCode != 200 {
		log.Print("Error:", res.StatusCode)
		return "", "", fmt.Errorf("error:%s", string(body))
	}
	var data model.TikTokVideoData
	if err = json.Unmarshal(body, &data); err != nil {
		return "", "", fmt.Errorf("error unmarshalling JSON: %v", err)
	}
	if len(data.ItemList) == 0 {
		return "", "", fmt.Errorf("aweme list is empty")
	}
	return data.ItemList[0].Video.PlayAddr.URLList[0], data.ItemList[0].Desc, nil
}

func GetTikTokInfo(link string) (string, string, error) {
	videoIdOrLink := GetTikTokId(link)
	if videoIdOrLink == "" {
		return "", "", fmt.Errorf("invalid link")
	}
	f, d, err := GetTikTokVideoData(videoIdOrLink)
	if err != nil {
		return "", "", err
	}
	return f, d, nil
}
