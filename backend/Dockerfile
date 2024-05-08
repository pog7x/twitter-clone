FROM golang:1.22.2
WORKDIR /backend
COPY . .
RUN go build -o /twitter-clone main.go 
CMD ["/twitter-clone"]