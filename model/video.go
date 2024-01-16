package model

// FrameInfo 用于保存帧的信息和Base64编码
type FrameInfo struct {
	SegmentIndex int    `json:"segment_index"`
	FrameIndex   int    `json:"frame_index"`
	Base64Data   string `json:"base64_data"`
}

// VideoInfo 包含视频的基本信息
type VideoInfo struct {
	Format struct {
		Duration string `json:"duration"` // 视频持续时间，以秒为单位
	} `json:"format"`
}

type DouYinVideoInfo struct {
	ItemList []struct {
		Video struct {
			PlayAddr struct {
				Uri string `json:"uri"`
			} `json:"play_addr"`
		} `json:"video"`
		Desc string `json:"desc"`
	} `json:"item_list"`
}

type TikTokVideoData struct {
	StatusCode int `json:"status_code"`
	ItemList   []struct {
		Video struct {
			PlayAddr struct {
				URLList []string `json:"url_list"`
			} `json:"play_addr"`
		} `json:"video"`
		Desc string `json:"desc"`
	} `json:"aweme_list"`
}

type VideoReply struct {
	Content  string  `json:"content"`
	Duration float64 `json:"duration"`
}
