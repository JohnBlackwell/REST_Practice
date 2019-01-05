FROM golang
RUN go get github.com/go-martini/martini
RUN go get github.com/lib/pq

COPY . /go/src/github.com/golang/IdeaEvolver

EXPOSE 3000
WORKDIR /go/src/github.com/golang/IdeaEvolver
CMD ["go", "run", "server.go", "dbfuncs.go", "validator.go"]
