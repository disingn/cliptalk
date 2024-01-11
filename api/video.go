package api

import (
	"bytes"
	"douyinshibie/cfg"
	"douyinshibie/model"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"
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

func GetVideoInfo(videoId string) (string, string, error) {
	url := fmt.Sprintf("https://www.iesdouyin.com/web/api/v2/aweme/iteminfo/?item_ids=%s&a_bogus=64745b2b5bdc4e75b720a9a85b19867a", videoId)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return "", "", err
	}
	req.Header.Add("User-Agent", randomUserAgent())
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Host", "www.iesdouyin.com")
	req.Header.Add("Connection", "keep-alive")

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

func getVideoDuration(videoURL string) (float64, error) {
	// 使用ffmpeg获取视频时长
	cmd := exec.Command("ffmpeg", "-i", videoURL)
	output, err := cmd.CombinedOutput()
	if err == nil {
		return 0, fmt.Errorf("ffmpeg didn't return an error, but it should have")
	}

	// 解析ffmpeg输出，寻找时长
	outputStr := string(output)
	if strings.Contains(outputStr, "Duration") {
		start := strings.Index(outputStr, "Duration: ")
		if start != -1 {
			end := strings.Index(outputStr[start:], ",")
			if end != -1 {
				durationStr := outputStr[start+10 : start+end]
				t, err := time.Parse("15:04:05.00", durationStr)
				if err == nil {
					return float64(t.Hour()*3600+t.Minute()*60+t.Second()) + float64(t.Nanosecond())/1e9, nil
				}
			}
		}
	}
	return 0, fmt.Errorf("unable to parse duration from ffmpeg output")
}

func VideoSlice(videoURL string) (error, string) {
	segments := 6 // 分成6段
	duration, err := getVideoDuration(videoURL)
	if err != nil {
		return fmt.Errorf("Error getting video duration: %s\n", err), ""
	}
	if duration < 60 {
		segments = 1
	}
	// 计算每个段的时长
	segmentDuration := duration / float64(segments)

	var wg sync.WaitGroup
	for i := 0; i < segments; i++ {
		wg.Add(1)
		go func(segmentIndex int) {
			defer wg.Done()
			startTime := float64(segmentIndex) * segmentDuration

			// 构建ffmpeg命令，使用管道输出
			cmd := exec.Command(
				"ffmpeg",
				"-i", videoURL,
				"-ss", strconv.FormatFloat(startTime, 'f', -1, 64),
				"-t", strconv.FormatFloat(segmentDuration, 'f', -1, 64),
				"-vf", "fps=1/5,scale=iw/5:-1", // 降低帧率和分辨率
				"-f", "image2pipe",
				"-vcodec", "mjpeg",
				"pipe:1",
			)

			// 创建管道
			stdoutPipe, err := cmd.StdoutPipe()
			if err != nil {
				log.Printf("Error creating stdout pipe for segment %d: %s\n", segmentIndex, err)
				return
			}

			// 启动命令
			if err := cmd.Start(); err != nil {
				log.Printf("Error starting ffmpeg for segment %d: %s\n", segmentIndex, err)
				return
			}

			// 读取数据并处理
			buffer := make([]byte, 4096) // 用于存储从管道读取的数据
			imageBuffer := new(bytes.Buffer)
			frameIndex := 0
			for {
				n, err := stdoutPipe.Read(buffer)
				//log.Print(len(buffer))
				if err != nil {
					if err == io.EOF {
						break // 管道关闭，没有更多的数据
					}
					log.Printf("Error reading from stdout pipe: %s\n", err)
					return
				}
				if n > 0 {
					imageBuffer.Write(buffer[:n])
					// 检查imageBuffer中是否存在JPEG结束标记
					if idx := bytes.Index(imageBuffer.Bytes(), []byte("\xff\xd9")); idx != -1 {
						// 截取到结束标记的部分作为一张完整的JPEG图像
						jpegData := imageBuffer.Bytes()[:idx+2] // 包含结束标记
						base64Data := base64.StdEncoding.EncodeToString(jpegData)
						imageBuffer.Next(idx + 2) // 移除已处理的JPEG图像数据

						// 输出带有标记的Base64字符串
						frameInfo := model.FrameInfo{
							SegmentIndex: segmentIndex,
							FrameIndex:   frameIndex,
							Base64Data:   base64Data,
						}
						frameIndex++
						if err := SetGeminiV(frameInfo); err != nil {
							log.Printf("Error processing frame info: %s\n", err)
							continue
						}
					}
				}
			}

			// 等待命令完成
			if err := cmd.Wait(); err != nil {
				log.Printf("Error waiting for ffmpeg command to finish for segment %d: %s\n", segmentIndex, err)
				return
			}
		}(i)
	}
	wg.Wait() // 等待所有goroutine完成
	fullDescription := strings.Join(FrameDescriptions, " ")
	//fmt.Println("Full description:", fullDescription)
	err, d := SetGemini(fullDescription)
	if err != nil {
		return err, ""
	}
	return nil, d
}
