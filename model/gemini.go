package model

type GeminiData struct {
	Contents []Contents `json:"contents"`
}
type InlineData struct {
	MimeType string `json:"mime_type"`
	Data     string `json:"data"`
}
type Parts struct {
	Text       string      `json:"text,omitempty"`
	InlineData *InlineData `json:"inline_data,omitempty"`
}
type Contents struct {
	Parts []Parts `json:"parts"`
}

type GeminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
			Role string `json:"role"`
		} `json:"content"`
		FinishReason  string `json:"finishReason"`
		Index         int    `json:"index"`
		SafetyRatings []struct {
			Category    string `json:"category"`
			Probability string `json:"probability"`
		} `json:"safetyRatings"`
	} `json:"candidates"`
	PromptFeedback struct {
		SafetyRatings []struct {
			Category    string `json:"category"`
			Probability string `json:"probability"`
		} `json:"safetyRatings"`
	} `json:"promptFeedback"`
}

type GeminiProPart struct {
	Text string `json:"text"`
}

type GeminiProContent struct {
	Role  string          `json:"role"`
	Parts []GeminiProPart `json:"parts"`
}

type GeminiPro struct {
	Contents []GeminiProContent `json:"contents"`
}
