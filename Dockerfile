FROM golang
RUN go get github.com/go-martini/martini
RUN go get github.com/lib/pq

COPY . /go/src/github.com/golang/IdeaEvolver

RUN go install /go/src/github.com/golang/IdeaEvolver
ENTRYPOINT /go/bin/IdeaEvolver
EXPOSE 8080
