BINARY = s3get
TARGETDIR = ./bin
TARGETOS = linux windows darwin

configure:
	go mod tidy

test:
	go test ../

build:
	mkdir -p $(TARGETDIR)
	$(foreach OS, $(TARGETOS), echo "Building $(BINARY)-$(OS)"; GOOS=$(OS) GOARCH=amd64 go build -v -o $(TARGETDIR)/$(BINARY)-$(OS) .;)