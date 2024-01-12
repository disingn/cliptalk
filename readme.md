# ClipTalk
[![forthebadge made-with-go](http://ForTheBadge.com/images/badges/made-with-go.svg)](https://go.dev/)

[![Build Status](https://img.shields.io/badge/build-passing-brightgreen)](https://github.com/disingn/cliptalk/actions)
[![License](https://img.shields.io/badge/license-MIT-blue)](LICENSE)
[![GitHub stars](https://img.shields.io/github/stars/disingn/cliptalk.svg?style=social&label=Star)](https://GitHub.com/disingn/cliptalk/stargazers/)


ClipTalk 是一个用于去除抖音视频水印和将视频解析成文本的工具，
目前已经兼容 tiktok。

## 目录

- [安装](#安装)
- [使用](#使用)
- [Docker 部署](#docker-部署)
- [本地开发](#本地开发)
- [其他](#其他)
- [联系我们](#联系我们)

## 安装 <a name="安装"></a>

### 克隆代码

```shell
git clone https://github.com/disingn/cliptalk.git
```

### 构建程序
`注意：这里我默认你本地或者服务器已经安装了 ffmpeg 和 go 环境，如果没有安装这两个，请先安装一下！！！不然跑不起来`

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
  # 自定义的 Gemini 的 url 地址 可以使用https://zhile.io/2023/12/24/gemini-pro-proxy.html#more-587来做代理
  # ps: 代理地址不要带最后的/
  #配置了 GeminiUrl 就不需要配置 Proxy
  GeminiUrl: https://gemini.baipiao.io
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
  #可以上传的文件大小 单位MB 默认10MB 不要写 0
  MaxFileSize: 10
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

## 使用 <a name="使用"></a>

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
  "content": "视频文本"
}
```
### 本地视频转文本接口
请求方式：POST
请求地址：`/video-file`

示例：

```shell
curl --location --request POST 'localhost:3100/video-file' \
--form 'file=@"/test.mp4"' \
--form 'model="openai"'
```
返回的 json 参数：
```json
{
  "content": "视频文本"
}
```

## Docker 部署 <a name="Docker 部署"></a>

### 准备工作

确保已安装 Docker 和 Docker Compose。

### 部署

```shell
cd cliptalk
docker-compose up -d
```
## 本地开发 <a name="本地开发"></a>
` 需要有一点的 go 的代码编写的一点经验`
### 需要的环境 （默认你都具备了）
- 安装 go
- 安装 ffmpeg

### 开发
```shell
cd cliptalk
go mod tidy
go run main.go
```

`代码目录写的也比较简单明了了，不再赘述了`
## 其他 <a name="其他"></a>

如果在使用过程中遇到问题，请加入我们的 QQ 群进行讨论。

QQ 群: 814702872

## 联系我们 <a name="联系我们"></a>

如有任何疑问或需要支持，请通过以下方式联系我们：

[![联系我们 !](https://img.shields.io/badge/Ask%20me-anything-1abc9c.svg)]([https://GitHub.com/Naereen/ama](https://github.com/disingn/cliptalk/issues))
