FROM golang:1.24.5-alpine

WORKDIR /app


COPY . .

RUN go mod init my-app
RUN go mod tidy 

COPY go.* ./

RUN go mod download

RUN go build -o main

CMD [ "./main" ]

# docker buildx build --platform linux/amd64,linux/arm64 \
#   -t tomtris02/paas_image:latest \
#   --push .