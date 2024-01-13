# ClipTalk
[![forthebadge made-with-go](http://ForTheBadge.com/images/badges/made-with-go.svg)](https://go.dev/)

[![Build Status](https://img.shields.io/badge/build-passing-brightgreen)](https://github.com/disingn/cliptalk/actions)
[![License](https://img.shields.io/badge/license-MIT-blue)](LICENSE)
[![GitHub stars](https://img.shields.io/github/stars/disingn/cliptalk.svg?style=social&label=Star)](https://GitHub.com/disingn/cliptalk/stargazers/)

[简体中文版](./readme_cn.md) 

ClipTalk is a tool designed for removing watermarks from TikTok videos and converting video content into text. It is now compatible with TikTok.

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
- [Docker Deployment](#docker-deployment)
- [Local Development](#local-development)
- [Miscellaneous](#miscellaneous)
- [Contact Us](#contact-us)

## Installation <a name="installation"></a>

### Clone the Repository

```shell
git clone https://github.com/disingn/cliptalk.git
```

### Build the Application
`Note: I assume that you have already installed ffmpeg and the go environment locally or on your server. If these are not installed, please install them first!!! Otherwise, it won't run.`

```shell
cd cliptalk
export GOOS=linux
export GOARCH=amd64
go build -o cliptalk
```

### Configuration File

Copy the example configuration file and modify it:

```shell
cp config.yaml.example config.yaml
```

Edit the `config.yaml` file and fill in the necessary configuration information:

```yaml
App:
  # Gemini's apikey
  GeminiKey:
    - key1
    - key2
  # Custom Gemini URL, you can use https://zhile.io/2023/12/24/gemini-pro-proxy.html#more-587 as a proxy
  # PS: Do not include a trailing slash in the proxy address
  # If you configure a GeminiUrl, you do not need to configure a Proxy
  GeminiUrl: https://gemini.baipiao.io
  # Browser UserAgent for parsing TikTok links
  UserAgents:
    - Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.6 Safari/605.2.15
    - Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.104 Safari/537.66
  # Note: Do not include a trailing slash in the OpenaiUrl. You can also change it to your own reverse proxy address or a compatible OpenAI address
  OpenaiUrl: https://api.openai.com
  OpenaiKey:
    - key1
    - key2
# Server Configuration
Sever:
  Port: 3100
  Host: localhost
  # Maximum file size for uploads in MB, default is 10MB, do not write 0
  MaxFileSize: 10
# # Proxy Configuration, use a proxy (http|https|socks5://ip:port)
# Proxy:
#     Protocol: socks5://192.168.1.10:3200

# Proxy Configuration, no proxy 
Proxy:
  Protocol: 
```

If you find the configuration process cumbersome, you can directly use the example configuration file.

### Start the Application

```shell
./cliptalk
```

### Configure Nginx Reverse Proxy

Please refer to the official Nginx documentation for configuration or use tools like Baota or 1panel.

## Usage <a name="usage"></a>

### API Endpoints

#### TikTok Watermark Removal API

Method: POST
Endpoint: `/remove`

Example:

```shell
curl --location --request POST 'localhost:3100/remove' \
--header 'Content-Type: application/json' \
--data-raw '{
    "url":"https://v.douyin.com/iLYNG8vA/"
}'
```

Returned JSON parameters:

```json
{
  "finalUrl": "Watermark-free video link",
  "message": "success",
  "title": "Video title"
}
```

#### TikTok Video to Text API

Method: POST
Endpoint: `/video`

Example:

```shell
curl --location --request POST 'localhost:3100/video' \
--header 'Content-Type: application/json' \
--data-raw '{
    "url":"https://v.douyin.com/iLYnjXbA/",
    "model":"openai" // The 'model' here can be 'openai' or 'gemini'
}'
```

Returned JSON parameters:

```json
{
  "finalUrl": "Watermark-free video link",
  "message": "success",
  "title": "Video title",
  "content": "Video text"
}
```
### Local Video to Text API
Method: POST
Endpoint: `/video-file`

Example:

```shell
curl --location --request POST 'localhost:3100/video-file' \
--form 'file=@"/test.mp4"' \
--form 'model="openai"'
```
Returned JSON parameters:
```json
{
  "content": "Video text"
}
```

## Docker Deployment <a name="docker-deployment"></a>

### Prerequisites

Make sure Docker and Docker Compose are installed.

### Deployment

```shell
cd cliptalk
docker-compose up -d
```
## Local Development <a name="local-development"></a>
`Some experience with writing Go code is required`
### Required Environment (assuming you already have it)
- Install Go
- Install ffmpeg

### Development
```shell
cd cliptalk
go mod tidy
go run main.go
```

`The code directory is also written in a simple and clear manner, no further explanation is needed.`
## Miscellaneous <a name="miscellaneous"></a>

For further assistance or if you have any questions, feel free to join our Telegram group.

[![cliptalk](https://img.shields.io/badge/Telegram-2CA5E0?style=for-the-badge&logo=telegram&logoColor=white)](https://t.me/cliptalk)

## Contact Us <a name="contact-us"></a>

If you have any questions or need support, please contact us through the following means:

[![Contact Us!](https://img.shields.io/badge/Ask%20me-anything-1abc9c.svg)](https://github.com/disingn/cliptalk/issues)
