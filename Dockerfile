FROM golang:latest

ENV PORT "8080"
ENV WORK_DIR "/workDir"

WORKDIR /app

COPY . ./

RUN go build -o server /app/fileServer

EXPOSE $PORT

CMD ["./server"]
