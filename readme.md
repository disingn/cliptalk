### 一、安装
#### 1.1 克隆代码
```shell
git clone https://github.com/disingn/cliptalk.git
```
#### 1.2 构建程序
```shell
cd cliptalk
export GOOS=linux                                                             
export GOARCH=amd64
go build -o cliptalk
```
#### 1.3 配置文件
```shell
cp config.yaml.example config.yaml
```
修改配置文件
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
#服务器配置
Sever:
  Port: 3100
  Host: localhost
```
嫌麻烦可以直接用实例的配置文件

#### 1.4 启动程序
```shell
./cliptalk
```
#### 1.5 配置Nginx反代
这个自己使用宝塔或者 1panel 都可以，这里就不多说了

### 二、使用
#### 2.1 接口
抖音去水印接口请求方式 POST 请求 请求地址：/remove 示例如下：
```shell
curl --location --request POST 'localhost:3100/remove' \
--header 'Content-Type: application/json' \
--data-raw '{
    "url":"https://v.douyin.com/iLYNG8vA/"
}'
```
返回的 json 参数：
```json
{
    "finalUrl": "去除水印的视频链接",
    "message": "success",
    "title": "视频标题"
}
```
---
抖音视频转文本接口请求方式 POST 请求 请求地址：/video 示例如下：
```shell
curl --location --request POST 'localhost:3100/video' \
--header 'Content-Type: application/json' \
--data-raw '{
    "url":"https://v.douyin.com/iLYNG8vA/"
}'
```
返回的 json 参数：
```json
{
    "finalUrl": "去除水印的视频链接",
    "message": "success",
    "title": "视频标题",
    "desc": "视频文本"
}
```
---
docker 部署：
#### 1.1准备工作
安装 docker 和 docker-compose （这里建议使用 docker-compose）
#### 1.2 部署
```shell
cd cliptalk
docker-compose up -d
```

### 三、其他
其他的就都和上面的一样了，如果有问题可以加群讨论

QQ: 814702872




