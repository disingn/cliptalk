package api

import (
	"bytes"
	"douyinshibie/model"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

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

func VideoSlice(videoURL, m string) (error, string) {
	segments := 6      // 分成6段
	intercept := "1/5" // 每5帧截取一帧
	duration, err := getVideoDuration(videoURL)
	//log.Println("Video duration:", duration)
	if err != nil {
		return fmt.Errorf("Error getting video duration: %s\n", err), ""
	}
	if duration < 60 && duration > 10 {
		segments = 1
		intercept = "1/5"
	} else if duration <= 10 {
		segments = 5
		intercept = "1" // 每秒截取一帧
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
				"-vf", fmt.Sprintf("fps=%s,scale=iw/5:-1", intercept), // 降低帧率和分辨率
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
						if m == "gemini" {
							if err := SetGeminiV(frameInfo); err != nil {
								log.Printf("Error processing frame info: %s\n", err)
								continue
							}
						} else if m == "openai" {
							if err := SetGptV(frameInfo); err != nil {
								log.Printf("Error processing frame info: %s\n", err)
								continue
							}
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
	d := ""
	if m == "gemini" {
		err, d = SetGemini(fullDescription)
		if err != nil {
			return err, ""
		}
	} else if m == "openai" {
		err, d = SetGpt(fullDescription)
		if err != nil {
			return err, ""
		}
	}
	return nil, d
}
