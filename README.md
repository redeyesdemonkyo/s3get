# s3get simple cli object storage downloader

## Why use it
`s3get` provides a simple cli go binary that allows downloading of S3 compatible
storage objects such as AWS S3, Google Cloud Storage (Using HMAC credentials), or Ceph
(Using Ceph Object Gateway)

The idea is to follow the `linux/Unix Principle` that a cli utility should be
small and perform one task very well.  Think of `wget` as an example.

Unline many other already available tools, this one is mainly interface to use
the `access key` and `secret key` for authentication and additionaly does not
require external usage of libraries like boto.

A good use case for `s3get` is where you are using private buckets and objects as repos
for your packages or other private binaries.

`s3get` is written by leveraging the `aws-sdk-go` golang package.

## binaries

You can download any of these pre-build binaries for your operating system

* [linux s3get](https://objects-us-east-1.dream.io/pub-binaries/linux/s3get) sha256: 4771e6924befceaa3fb9fb235d8b7d85395fdc8e014bae90f941437bb02741bd
  
* [darwin s3get](https://objects-us-east-1.dream.io/pub-binaries/darwin/s3get) sha256: d83788f39b39827f7ac056e08c03e891e8e8ac26ec3767f522d817fc36405d27

* [windows s3get](https://objects-us-east-1.dream.io/pub-binaries/windows/s3get) sha256: 8d3792c6f3694369b089eed24a0d887fb5557dc3b85fcedbd7c65648d993ae06


## Installation

No installation required if you get the pre-build [binaries](#binaries) else checkout the [building](#building) section


## requirements

You need to have an `access key` and `secret key` with `READ` access to the bucket and object for non public objects


## Usage
* get full usage help menu

```sh
s3get -h
NOTE: You must specify the bucket and object file to download.
Remember you must also specify your access key and secret key as either environment variables
or pass them as flags to the command!!!


Usage: /tmp/go-build3389196233/b001/exe/s3get -b <bucket> -o <path/to/my.object>
  -a string
        Access key.  Defaults to using environment variable: AWS_ACCESS_KEY
  -b string
        Bucket name
  -d string
        Destination path ie for linux/Mac: /path/2/save/ or for Windows: C:\temp\ 
  -e url
        URL endpoint for where to get your object.  Using url (default "https://storage.googleapis.com")
  -h    Print usage info
  -o string
        Object to download.  If the object is under a directory include the whole path: subdir/myobject.file
  -p    For public objects.  Will skip authentication
  -s string
        Secret key.  Defaults to using environment variable: AWS_SECRET_KEY
```

* downloading a public object from a Ceph object storage

```sh
s3get -e objects-us-east-1.dream.io -b imgun -o Downloads/linuxmint-20.3-mate-64bit.iso -p -d /home/flynn/tmp/
```

* download object file from Google Storage bucket.  Using [HMAC keys](https://cloud.google.com/storage/docs/authentication/hmackeys) and exported as environment variables; `AWS_SECRET_KEY` and `AWS_ACCESS_KEY`  The file will be save to `~/tmp` dir

> NOTE: no need to specify endpoint since GCP is the default endpoint for s3get

```sh
s3get -b mybucket -o path/to/my.object -d ~/tmp/
```

* download from AWS S3 bucket.  TODO usage example!!

```sh
s3get -e <https://aws-endpoint.tld> -a Acc3sKey -s SeCR3TKey -b thebucketName -o myfile.object
```

## development

* golang >= v1.18 (see installing [Go](https://go.dev/doc/manage-install))
* initialize `go.mod` ie: `go mod init <package-name>`

```sh
go mod init s3get
```

* Once you have your `go.mod` initialized you can use `go get <package>` to
intall require packages or simply run `go mod tidy`.  This will also add the package to your `go.mod`
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
go run s3get.go -e objects-us-east-1.dream.io -p -b imgun -o Downloads/linuxmint-20.3-mate-64bit.iso -d ~/tmp/
```

* testing windows TODO!!

```ps1
s3get TODO
```

* testing MacOSX (darwin)

```sh
s3get TODO
```

## Feature Improvements

* add test cases
* add flag for handling sha256 checks
* add rsync like functionality
* add flag for splitting using `https://s3provider.tld/mybucket/dir/to/object.file` link rather then having to pass individual flags for provider, bucket & object

## References

* aws documentation: https://docs.aws.amazon.com/sdk-for-go/api/aws/session
* helpful google golang documentation example using HMAC credentials: https://cloud.google.com/storage/docs/samples/storage-s3-sdk-list-objects
* google HMAC docs: https://cloud.google.com/storage/docs/authentication/hmackeys
* installing Golang: https://go.dev/doc/install or https://go.dev/doc/manage-install or https://golangdocs.com/install-go-linux
* using bar: https://github.com/cheggaaa/pb

## License
  
This is free software under the terms of MIT the license, See [LICENSE](https://github.com/redeyesdemonkyo/s3get/blob/main/LICENSE) for more information.
