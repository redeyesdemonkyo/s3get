package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// defining the string literal for usage
const usageMessage = `NOTE: You must specify the bucket and object file to download.
Remember you must also specify your access key and secret key as either environment variables
or pass them as flags to the command!!!


Usage: %s -b <bucket> -o <path/to/my.object>
`

// defining global vars
var EndPoint string
var srcBucket string
var srcObject string
var secretKey string
var accessKey string
var Help bool

func init() {
	flag.StringVar(&EndPoint, "e", "https://storage.googleapis.com", "URL endpoint for where to get your object.  Using `url`") // by using `url` it overrides flag type from string to url
	flag.StringVar(&srcBucket, "b", "", "Bucket name")
	flag.StringVar(&srcObject, "o", "", "Object to download.  If the object is under a directory include the whole path: 'subdir/my.object")
	flag.StringVar(&secretKey, "s", os.Getenv("AWS_SECRET_KEY"), "Secret key.  Defaults to using environment variable: AWS_SECRET_KEY")
	flag.StringVar(&accessKey, "a", os.Getenv("AWS_ACCESS_KEY"), "Access key.  Defaults to using environment variable: AWS_ACCESS_KEY")
	flag.BoolVar(&Help, "h", false ,"Print usage info")

}

func main() {
	var err error

	// define custom usage message
	flag.Usage = usage
	flag.Parse()

	// print usage or verify we got at least bucket
	if Help || len(srcBucket) < 1 {
		usage()
	}

	// verify we have require keys to auth
	if len(accessKey) < 1 || len(secretKey) < 1 {
		fmt.Print("!!You did NOT specified access & secret key.  See usage for more info!!\n\n")
		usage()
	}

	// aws: https://docs.aws.amazon.com/sdk-for-go/api/aws/session
	// helpful google golang example using HMAC credentials: https://cloud.google.com/storage/docs/samples/storage-s3-sdk-list-objects
	mySession, err := session.NewSession(&aws.Config{
		Endpoint:         aws.String(EndPoint),
		Region:           aws.String("us-west-2"), // only here cause its required by the library but does not matter for the actual download procees.
		S3ForcePathStyle: aws.Bool(true),
		Credentials:      credentials.NewStaticCredentials(accessKey, secretKey, ""),
	})
	if err != nil {
		msg := fmt.Errorf("failed to initilize session %v", err)
		catch(msg)
	}

	// create a client from the session and pass additional configuration
	//s3Client := s3.New(mySession)
	s3Client := s3manager.NewDownloader(mySession)

	// create file so we can write to it from NewDownloader
	BaseObject := filepath.Base(srcObject)
	fmt.Printf("Creating file object: %s\n", BaseObject)
	f, err := os.Create(BaseObject)
	if err != nil {
		//return fmt.Errorf("failed to create file %q, %v", srcObject, err)
		msg := fmt.Errorf("failed to create file %q, %v", srcObject, err)
		catch(msg)
	}

	// download it and write to file
	fmt.Printf("Downloading %s/%s\n", srcBucket, srcObject)
	n, err := s3Client.Download(f, &s3.GetObjectInput{
		Bucket: aws.String(srcBucket),
		Key:    aws.String(srcObject),
	})
	if err != nil {
		msg := fmt.Errorf("failed to download file, %v", err)
		catch(msg)
	}

	fmt.Printf("file downloaded, %d bytes\n", n)
}

func usage() {
	fmt.Printf(usageMessage, os.Args[0])
	flag.PrintDefaults()
	os.Exit(1)
}

func catch(err error) {
	if err != nil {
		fmt.Printf("[Error] We encountered an error:\n\n\t%s\n\n", err)
		os.Exit(1)
	}
}
