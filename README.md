The `tetris_go` is a doodle in Go(lang).

    go mod tidy
    go vet .
    go build -o tetris
    # GOOS=windows GOARCH=386 go build -o tetris32.exe
    # GOOS=windows GOARCH=amd64 go build -o tetris64.exe
    go test
