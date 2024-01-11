package model

type Config struct {
	App struct {
		GeminiKey  []string `yaml:"GeminiKey"`  // GeminiKey 是一个字符串切片
		UserAgents []string `yaml:"UserAgents"` // UserAgents 是一个字符串切片
	} `yaml:"App"`
	Sever struct {
		Host string `yaml:"Host"`
		Port string `yaml:"Port"`
	}
}
