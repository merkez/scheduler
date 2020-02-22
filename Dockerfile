FROM golang
COPY . /go/src/github.com/mrturkmen06/scheduler
WORKDIR /go/src/github.com/mrturkmen06/scheduler
RUN  go mod download
CMD ["go","run","/go/src/github.com/mrturkmen06/scheduler/main.go"]
