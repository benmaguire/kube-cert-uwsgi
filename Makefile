build:
	go get . && CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags "-s -w" -o gokubewatcher .

run:
	go run main.go
