# use an official Golang runtime as the base image
FROM golang:latest

ENV PORT "8080"
ENV WORK_DIR "/workDir"

WORKDIR /app

# copy all files from the current directory to /app
COPY . ./

# build application
RUN go build -o server /app/fileServer

# expose the port that the application will listen on
EXPOSE $PORT

# run the application
CMD ["./server"]
