FROM golang:1.26.0
WORKDIR /backend
COPY . .
RUN go build -o /twitter-clone main.go 
CMD ["/twitter-clone", "serve", "--seed"] 