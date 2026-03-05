FROM golang:1.24.0
WORKDIR /backend
COPY . .
RUN go build -o /twitter-clone main.go 
CMD ["/twitter-clone", "serve", "--seed"] 