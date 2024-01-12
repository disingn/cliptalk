package model

type Config struct {
	App struct {
		GeminiKey  []string `yaml:"GeminiKey"`  // GeminiKey
		GeminiUrl  string   `yaml:"GeminiUrl"`  // GeminiUrl
		UserAgents []string `yaml:"UserAgents"` // UserAgents
		OpenaiUrl  string   `yaml:"OpenaiUrl"`
		OpenaiKey  []string `yaml:"OpenaiKey"`
	} `yaml:"App"`
	Sever struct {
		Host        string `yaml:"Host"`
		Port        string `yaml:"Port"`
		MaxFileSize int    `yaml:"MaxFileSize"`
	}
	Proxy struct {
		Protocol string `yaml:"Protocol"` //Protocol 留空 表示不用代理
	}
}
