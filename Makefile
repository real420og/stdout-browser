build_linux:
	CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags "-s" -o stdout-browser .
build_darwin:
	CGO_ENABLED=0 GOOS=darwin go build -trimpath -ldflags "-s" -o stdout-browser .
