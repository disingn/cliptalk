# 使用 Node.js 镜像构建前端项目
FROM node:18 as frontend

WORKDIR /app

COPY web/cliptalk/package.json web/cliptalk/yarn.lock ./

RUN yarn install

COPY web/cliptalk/ ./

RUN yarn build

# 使用官方Go镜像作为构建环境
FROM golang:1.20 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ClipTalk .

# 使用alpine作为最终镜像
FROM alpine

RUN apk update \
    && apk upgrade \
    && apk add --no-cache ca-certificates tzdata \
    && update-ca-certificates 2>/dev/null || true

RUN apk add --no-cache ffmpeg

# 从构建器镜像中复制执行文件
COPY --from=builder /app/ClipTalk /ClipTalk

# 从前端构建阶段复制构建结果到最终镜像
COPY --from=frontend /app/dist /web/cliptalk/dist

EXPOSE 3100

CMD ["/ClipTalk"]