FROM golang:1.17.1-alpine3.14 as builder

### Install Tesseract and its dependencies
RUN apk update && apk add \
    g++ \
    musl-dev \
    tesseract-ocr-dev \
    tesseract-ocr-data-por

### Start build flow
WORKDIR /app

# Set Go build envs
ENV GO111MODULE=on \
    CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64

# Copy all files from source and compile code
COPY . .

RUN go build -ldflags '-s -w' -o ./bin/birus ./cmd/api

FROM alpine:3.14

# Copy required 
COPY --from=builder /app/bin /bin
COPY --from=builder /usr/lib /usr/lib
COPY --from=builder /usr/share/tessdata/ /usr/share/tessdata/

# Set Tesseract training data dir
ENV TESSDATA_PREFIX=/usr/share/tessdata/

EXPOSE 8000

CMD [ "birus" ]
