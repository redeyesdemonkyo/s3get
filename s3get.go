package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"hash"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/cheggaaa/pb/v3"
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
var Dest string
var Anonyous bool
var checksum string
var Help bool

func init() {
	flag.StringVar(&EndPoint, "e", "https://storage.googleapis.com", "URL endpoint for where to get your object.  Using `url`") // by using `url` it overrides flag type from string to url
	flag.StringVar(&srcBucket, "b", "", "Bucket name")
	flag.StringVar(&Dest, "d", "", "Destination path ie for linux/Mac: /path/2/save/ or for Windows: C:\\temp\\ ")
	flag.StringVar(&srcObject, "o", "", "Object to download.  If the object is under a directory include the whole path: subdir/myobject.file")
	flag.StringVar(&secretKey, "s", os.Getenv("AWS_SECRET_KEY"), "Secret key.  Defaults to using environment variable: AWS_SECRET_KEY")
	flag.StringVar(&accessKey, "a", os.Getenv("AWS_ACCESS_KEY"), "Access key.  Defaults to using environment variable: AWS_ACCESS_KEY")
	flag.BoolVar(&Anonyous, "p", false, "For public objects.  Will skip authentication")
	flag.StringVar(&checksum, "checksum", "", "the algo:hash to verify the oject checksum.  Algos supported are: sha256, sha1 & md5")
	flag.StringVar(&checksum, "c", "", "the algo:hash to verify the oject checksum.  Algos supported are: sha256, sha1 & md5")
	flag.BoolVar(&Help, "h", false, "Print usage info")
}

func usage() {
	fmt.Printf(usageMessage, os.Args[0])
	flag.PrintDefaults()
	os.Exit(1)
}

// method for the progressbar
type ProgressWriter struct {
	w  io.WriterAt
	pb *pb.ProgressBar
}

// define method function of type ProgressWriter
func (pw ProgressWriter) WriteAt(p []byte, off int64) (int, error) {
	pw.pb.Add(len(p))
	return pw.w.WriteAt(p, off)
}

func main() {
	var err error

	// define custom usage message
	flag.Usage = usage
	flag.Parse()

	// define Creds varriable as aws.Config struct type
	var Creds = &aws.Config{
		Endpoint:         aws.String(EndPoint),
		Region:           aws.String("us-west-2"), // only here cause its required by the library but does not matter for the actual download procees.
		S3ForcePathStyle: aws.Bool(true),
	}

	// print usage or verify we got at least bucket
	if Help || len(srcBucket) < 1 {
		usage()
	}

	if Anonyous {
		// append struct field to the struct variable
		Creds.Credentials = credentials.AnonymousCredentials
	} else {
		Creds.Credentials = credentials.NewStaticCredentials(accessKey, secretKey, "")

		// verify we have require keys to auth
		if len(accessKey) < 1 || len(secretKey) < 1 {
			fmt.Print("!!You did NOT specified access & secret key.  See usage for more info!!\n\n")
			usage()
		}
	}

	// aws: https://docs.aws.amazon.com/sdk-for-go/api/aws/session
	// helpful google golang example using HMAC credentials: https://cloud.google.com/storage/docs/samples/storage-s3-sdk-list-objects
	mySession, err := session.NewSession(Creds)
	if err != nil {
		msg := fmt.Errorf("failed to initilize session %v", err)
		catch(msg)
	}

	// create a client from the session and pass additional configuration
	s3Client := s3.New(mySession)
	s3Downloader := s3manager.NewDownloader(mySession)

	BaseObject := filepath.Base(srcObject)
	var Dir, Base, Ext, DestFile string
	if len(Dest) > 1 {
		Dir = filepath.Dir(Dest)
		Base = filepath.Base(Dest)
		Ext = filepath.Ext(Dest)

		if Ext != "" {
			DestFile = Dir + "/" + Base
		} else {
			DestFile = Dir + "/" + BaseObject
		}
	} else {
		DestFile = BaseObject
	}

	objectSize, err := getFileSize(s3Client, srcBucket, srcObject)
	if err != nil {
		catch(fmt.Errorf("There was an error getting file size %v", err))
	}

	// template for bar.  Uses unicode characters for cycle
	tmpl := `{{ counters . | red }} {{ bar . "[" "#" (cycle . "↖" "↗" "↘" "↙" | yellow ) "." "]" | green }} {{ percent . }} {{ rtime . | yellow }} {{ speed . | green }}`

	// initiate bar
	//bar := pb.Full.Start64(objectSize) // for using with out template
	bar := pb.ProgressBarTemplate(tmpl).Start64(objectSize)
	bar.Set(pb.Bytes, true) // format as bytes output ie B, KiB, MiB
	temp, err := ioutil.TempFile(Dir, "s3get-tmp-")
	writer := &ProgressWriter{temp, bar}

	fmt.Printf("Creating file object: %s with total size of %d\n", DestFile, objectSize)

	// download it and write to file
	fmt.Printf("Downloading oject: %s from bucket: %s\n", srcObject, srcBucket)
	n, err := s3Downloader.Download(writer, &s3.GetObjectInput{
		Bucket: aws.String(srcBucket),
		Key:    aws.String(srcObject),
	})
	if err != nil {
		catch(fmt.Errorf("failed to download file, %v", err))
	}

	// checksum verification
	if len(checksum) > 1 {
		if err := verifyCheckSum(checksum, temp); err != nil {
			temp.Close()
			os.Remove(temp.Name())
			catch(fmt.Errorf("checksum verification failed: %v", err))
		}
	}

	if err := temp.Close(); err != nil {
		catch(fmt.Errorf("failed to close temp file %v", err))
	}

	// move and rename temp file base on flag argument
	if err := os.Rename(temp.Name(), DestFile); err != nil {
		catch(fmt.Errorf("failed to rename temp file error: %v", err))
	}

	fmt.Printf("file downloaded, %d bytes\n", n)
}

func verifyCheckSum(c string, f *os.File) error {
	// which hash algo from flags
	sp := strings.Split(c, ":")
	fmt.Printf("algo: %s hash: %s\n", sp[0], sp[1])

	var h hash.Hash
	if sp[0] == "sha256" {
		h = sha256.New()
	} else if sp[0] == "sha1" {
		h = sha1.New()
	} else if sp[0] == "md5" {
		h = md5.New()
	} else {
		return fmt.Errorf("Error no supported algo defined")
	}

	if _, err := io.Copy(h, f); err != nil {
		return fmt.Errorf("Unable to checksum data")
	}

	// calculate hash
	CalcSum := h.Sum(nil)
	strSum := hex.EncodeToString(CalcSum) // Sum hash requires to be hex decoded
	fmt.Printf("calculated sum: %s\n", strSum)
	if strSum != sp[1] {
		return fmt.Errorf("Error hash did not match")
	} else {
		return nil
	}
}

func getFileSize(c *s3.S3, bucket string, key string) (size int64, error error) {
	params := &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	s, err := c.HeadObject(params)
	if err != nil {
		return 0, err
	}

	return *s.ContentLength, nil
}

func catch(err error) {
	if err != nil {
		fmt.Printf("[Error] We encountered the following error:\n\n\t%s\n\n", err)
		os.Exit(1)
	}
}
