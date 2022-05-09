BINARY = s3get
TARGETDIR = ./bin
TARGETOS = linux windows darwin

configure:
	go mod tidy

test:
	go test ../

build:
	mkdir -p $(TARGETDIR)
	$(foreach OS, $(TARGETOS), echo "Building $(BINARY)-$(OS)"; mkdir $(TARGETDIR)/$(OS); GOOS=$(OS) GOARCH=amd64 go build -v -o $(TARGETDIR)/$(OS)/$(BINARY) .;)

clean:
	$(foreach OS, $(TARGETOS), echo "Rmoving $(OS)"; rm -rv ./$(TARGETDIR)/$(OS) ;)