FROM golang
COPY . /go/src/github.com/mrturkmen06/scheduler
WORKDIR /go/src/github.com/mrturkmen06/scheduler
RUN  go mod download
ENTRYPOINT ["go","run","/go/src/github.com/mrturkmen06/scheduler/main.go"]
# this approach used to parse flags which will be provided by client
CMD ["-c"]