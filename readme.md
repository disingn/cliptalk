# ClipTalk
[![forthebadge made-with-go](http://ForTheBadge.com/images/badges/made-with-go.svg)](https://go.dev/)

[![Build Status](https://img.shields.io/badge/build-passing-brightgreen)](https://github.com/disingn/cliptalk/actions)
[![License](https://img.shields.io/badge/license-MIT-blue)](LICENSE)
[![GitHub stars](https://img.shields.io/github/stars/disingn/cliptalk.svg?style=social&label=Star)](https://GitHub.com/disingn/cliptalk/stargazers/)


ClipTalk 是一个用于去除抖音视频水印和将视频转为文本的工具。

## 目录

- [安装](#安装)
- [使用](#使用)
- [Docker 部署](#docker-部署)
- [其他](#其他)
- [联系我们](#联系我们)

## 安装

### 克隆代码

```shell
git clone https://github.com/disingn/cliptalk.git
```

### 构建程序

```shell
cd cliptalk
export GOOS=linux
export GOARCH=amd64
go build -o cliptalk
```

### 配置文件

复制示例配置文件并修改：

```shell
cp config.yaml.example config.yaml
```

编辑 `config.yaml` 文件，填入必要的配置信息：

```yaml
App:
  #Gemini 的 apikey
  GeminiKey:
    - key1
    - key2
  #浏览器的 UserAgent 用来解析抖音链接
  UserAgents:
    - Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.6 Safari/605.2.15
    - Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.104 Safari/537.66
  #注意：这里的 OpenaiUrl 最后面不带 / 你也可以修改成你自己反代的地址或者兼容 OpenAI 的地址
  OpenaiUrl: https://api.openai.com
  OpenaiKey:
    - key1
    - key2
#服务器配置
Sever:
  Port: 3100
  Host: localhost

# #代理配置 用代理( http|https|socks5://ip:port ) 
# Proxy:
#     Protocol: socks5://192.168.1.10:3200

#代理配置 不用代理 
Proxy:
  Protocol: 
```

如果你觉得配置过程繁琐，可以直接使用实例的配置文件。

### 启动程序

```shell
./cliptalk
```

### 配置 Nginx 反代

请参考 Nginx 官方文档进行配置，或使用宝塔、1panel 等工具。

## 使用

### 接口

#### 抖音去水印接口

请求方式：POST
请求地址：`/remove`

示例：

```shell
curl --location --request POST 'localhost:3100/remove' \
--header 'Content-Type: application/json' \
--data-raw '{
    "url":"https://v.douyin.com/iLYNG8vA/"
}'
```

返回的 JSON 参数：

```json
{
  "finalUrl": "去除水印的视频链接",
  "message": "success",
  "title": "视频标题"
}
```

#### 抖音视频转文本接口

请求方式：POST
请求地址：`/video`

示例：

```shell
curl --location --request POST 'localhost:3100/video' \
--header 'Content-Type: application/json' \
--data-raw '{
    "url":"https://v.douyin.com/iLYnjXbA/",
    "model":"openai" //这里的 model 可以是 openai 或者 gemini
}'
```

返回的 JSON 参数：

```json
{
  "finalUrl": "去除水印的视频链接",
  "message": "success",
  "title": "视频标题",
  "desc": "视频文本"
}
```

## Docker 部署

### 准备工作

确保已安装 Docker 和 Docker Compose。

### 部署

```shell
cd cliptalk
docker-compose up -d
```

## 其他

如果在使用过程中遇到问题，请加入我们的 QQ 群进行讨论。

QQ 群: 814702872

## 联系我们

如有任何疑问或需要支持，请通过以下方式联系我们：

[![联系我们 !](https://img.shields.io/badge/Ask%20me-anything-1abc9c.svg)]([https://GitHub.com/Naereen/ama](https://github.com/disingn/cliptalk/issues))
