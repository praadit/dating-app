server:
	go run main.go

t:
	go test ./test/ -count=1

tv:
	go test ./test/ -count=1 -v