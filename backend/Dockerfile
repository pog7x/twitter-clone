FROM golang:1.22.2

RUN go install github.com/cespare/reflex@latest

WORKDIR /backend

COPY . .

CMD reflex -c reflex.conf

# CMD [ "go", "run", "main.go" ]