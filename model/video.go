package model

// FrameInfo 用于保存帧的信息和Base64编码
type FrameInfo struct {
	SegmentIndex int    `json:"segment_index"`
	FrameIndex   int    `json:"frame_index"`
	Base64Data   string `json:"base64_data"`
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
