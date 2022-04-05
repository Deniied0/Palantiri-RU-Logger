env GOOS=windows GOARCH=amd64 go build -o palantiri-windows-amd64
env GOOS=linux GOARCH=amd64 go build -o palantiri-linux-amd64
env GOOS=dragonfly GOARCH=amd64 go build -o palantiri-dragonfly-amd64
env GOOS=freebsd GOARCH=amd64 go build -o palantiri-freebsd-amd64
env GOOS=darwin GOARCH=amd64 go build -o palantiri-darwin-amd64
env GOOS=netbsd GOARCH=amd64 go build -o palantiri-netbsd-amd64
env GOOS=plan9 GOARCH=amd64 go build -o palantiri-plan9-amd64
env GOOS=solaris GOARCH=amd64 go build -o palantiri-solaris-amd64
