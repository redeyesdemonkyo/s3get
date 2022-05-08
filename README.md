# s3get simple cli object storage downloader

## WHy should you use it
s3get provides a simple cli go binary that allows downloading of S3 compatible
storage such as AWS S3, Google Cloud Storage (Using HMAC credentials), or Ceph
(Using Ceph Object Gateway)

The idea is to follow the `linux/Unix Principle` that a cli utility should be
small and perform one task very well.  Think of `wget` as an example.

Unline many other already available tools, this one is mainly interface to use
the `access key` and `secret key` for authentication and additionaly does not
require external usage of libraries like boto

s3get is written by leveraging the `aws-sdk-go` golang package

## requirements

* You need to have an `access key` and `secret key` with access to the bucket and object


## Usage
* TODO

```sh
```

## development

* golang >= v1.18 (see installing [Go](https://go.dev/doc/manage-install))
* initialize `go.mod` ie: `go mod init <package-name>`

```sh
go mod init s3get
```

* Once you have your `go.mod` initialized you can use `go get <package>` to
intall require packages.  This will also add the package to your `go.mod`
require packages.

```sh
go get github.com/aws/aws-sdk-go/
```

* add module requirements and sums

```sh
go mod tidy
```

## building

* view list of build supported platforms (operating systems)

```sh
go tool dist list
```

* building

```sh
GOOS=linux GOARCH=amd64 go build -v -o ./bin/s3get-linux 
```


* building using `Makefile` which will build for all 3 OS: linux, windows, darwin (OSx)


```sh
make build
```


### Testing

* test run with out compiling

```sh
go run s3get.go -b mybucket -o subDirectory/my.object
```

* test downloading public file (uses anonymous authentication) 


```sh
go run s3get.go -e objects-us-east-1.dream.io -p -b imgun -o capeta_v1.jpg
```

## TODO

* add flag for handling sha256 checks
* add a percentage or progress bar (try using async)
* add flag for splitting using `https://s3provider.tld/mybucket/dir/to/object.file` link rather then having to pass individual flags for provider, bucket & object

## References

* aws documentation: https://docs.aws.amazon.com/sdk-for-go/api/aws/session
* helpful google golang documentation example using HMAC credentials: https://cloud.google.com/storage/docs/samples/storage-s3-sdk-list-objects
* google HMAC docs: https://cloud.google.com/storage/docs/authentication/hmackeys
* installing Golang: https://go.dev/doc/install or https://go.dev/doc/manage-install or https://golangdocs.com/install-go-linux

## License
* TODO!!
This is free software under the terms of MIT the license
