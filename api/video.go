package api

import (
	"bytes"
	"douyinshibie/model"
	"encoding/base64"
	"encoding/json"
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

func VideoSlice(videoURL, m string) (error, *model.VideoReply) {
	segments := 6      // 分成6段
	intercept := "1/5" // 每5帧截取一帧
	duration, err := getVideoDuration(videoURL)
	//log.Println("Video duration:", duration)
	if err != nil {
		return fmt.Errorf("Error getting video duration: %s\n", err), nil
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
	var FrameDescriptions []string
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
			var mu sync.Mutex
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
						var frameDescription string
						if m == "gemini" {
							err, frameDescription = SetGeminiV(frameInfo)
						} else if m == "openai" {
							err, frameDescription = SetGptV(frameInfo)
						}
						if err != nil {
							log.Printf("Error creating stdout pipe for segment %d: %s", segmentIndex, err)
							continue
						}
						if len(frameDescription) != 0 {
							mu.Lock()
							FrameDescriptions = append(FrameDescriptions, frameDescription)
							mu.Unlock()
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
	d := ""
	if m == "gemini" {
		err, d = SetGemini(fullDescription)
		if err != nil {
			return err, nil
		}
	} else if m == "openai" {
		err, d = SetGpt(fullDescription)
		if err != nil {
			return err, nil
		}
	}
	return nil, &model.VideoReply{
		Content:  d,
		Duration: duration,
	}
}

// getVideoDurationFromStream 从视频流中获取视频的持续时间
func getVideoDurationFromStream(videoStream io.Reader) (float64, error) {
	cmd := exec.Command("ffprobe",
		"-i", "pipe:0", // 从标准输入读取视频数据
		"-show_entries", "format=duration",
		"-print_format", "json",
		"-v", "quiet",
	)
	cmd.Stdin = videoStream

	output, err := cmd.Output()
	if err != nil {
		log.Println("Error running ffprobe:", err)
		return 0, err
	}

	// 解析输出以获取视频持续时间
	var info model.VideoInfo
	if err := json.Unmarshal(output, &info); err != nil {
		log.Println("Error running ffprobe Unmarshal:", err)
		return 0, err
	}

	// 将持续时间字符串转换为 float64
	duration, err := strconv.ParseFloat(info.Format.Duration, 64)
	if err != nil {
		log.Println("Error parsing duration:", err)
		return 0, err
	}

	return duration, nil
}

func VideoFileSlice(videoStream io.Reader, m string) (error, string) {

	cmd := exec.Command(
		"ffmpeg",
		"-i", "pipe:0", // 使用管道作为输入
		"-vf", "fps=1/5,scale=iw/5:-1", // 降低帧率和分辨率
		"-f", "image2pipe",
		"-vcodec", "mjpeg",
		"pipe:1",
	)
	cmd.Stdin = videoStream

	// 创建管道
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		log.Printf("Error creating stdout pipe: %s\n", err)
		return err, ""
	}

	// 启动命令
	if err := cmd.Start(); err != nil {
		log.Printf("Error starting ffmpeg: %s\n", err)
		return err, ""
	}
	// 读取数据并处理
	buffer := make([]byte, 4096) // 用于存储从管道读取的数据
	imageBuffer := new(bytes.Buffer)
	var frameInfo model.FrameInfo // 存储图像数据和索引
	var FrameDescriptions []string
	var wg sync.WaitGroup
	frameStringChan := make(chan string, 10)
	frameIndex := 0
	for {
		n, err := stdoutPipe.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break // 管道关闭，没有更多的数据
			}
			log.Printf("Error reading from stdout pipe: %s\n", err)
			return err, ""
		}
		if n > 0 {
			imageBuffer.Write(buffer[:n])
			// 检查imageBuffer中是否存在JPEG结束标记
			if idx := bytes.Index(imageBuffer.Bytes(), []byte("\xff\xd9")); idx != -1 {
				// 截取到结束标记的部分作为一张完整的JPEG图像
				jpegData := imageBuffer.Bytes()[:idx+2] // 包含结束标记
				imageBuffer.Next(idx + 2)               // 移除已处理的JPEG图像数据

				// 将 JPEG 图像编码为 Base64
				base64Data := base64.StdEncoding.EncodeToString(jpegData)
				frameInfo = model.FrameInfo{
					SegmentIndex: 1,
					FrameIndex:   frameIndex,
					Base64Data:   base64Data,
				}
				frameIndex++
				wg.Add(1)
				go func(mode string, f model.FrameInfo) {
					defer wg.Done()
					var desc string
					var err error
					if m == "gemini" {
						err, desc = SetGeminiV(f)
					} else if m == "openai" {
						err, desc = SetGptV(f)
					}
					if err != nil {
						log.Printf("Error processing frame info: %s\n", err)
						desc = "Error: " + err.Error()
					}
					frameStringChan <- desc
				}(m, frameInfo)

			}
		}
	}
	go func() {
		wg.Wait()
		close(frameStringChan) // 确保所有goroutine完成后关闭通道
	}()
	for desc := range frameStringChan {
		if len(desc) != 0 {
			FrameDescriptions = append(FrameDescriptions, desc)
		}
	}
	if err := cmd.Wait(); err != nil {
		log.Printf("Error waiting for ffmpeg command to finish: %s\n", err)
		return err, ""
	}
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
