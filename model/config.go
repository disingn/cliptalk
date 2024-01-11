package model

type Config struct {
	App struct {
		GeminiKey  []string `yaml:"GeminiKey"`  // GeminiKey 是一个字符串切片
		UserAgents []string `yaml:"UserAgents"` // UserAgents 是一个字符串切片
		OpenaiUrl  string   `yaml:"OpenaiUrl"`
		OpenaiKey  []string `yaml:"OpenaiKey"`
	} `yaml:"App"`
	Sever struct {
		Host string `yaml:"Host"`
		Port string `yaml:"Port"`
	}
	Proxy struct {
		Protocol string `yaml:"Protocol"` //Protocol 留空 表示不用代理
	}
}
