package awss3

import (
    "bytes"
    "log"
    "net/http"
    "os"
	//"fmt"

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3"
)


func WriteToS3(FileName string, S3_BUCKET string) {
	
    // Create a single AWS session (we can re use this if we're uploading many files)
    s, err := session.NewSession(&aws.Config{Region: aws.String("eu-central-1")})
    if err != nil {
        log.Fatal(err)
    }

    // Upload
	//FileName:="../DataDump/order_1593866787334.json"
	
    err = AddFileToS3(s, FileName,S3_BUCKET)
    if err != nil {
        log.Fatal(err)
    }
}

func AddFileToS3(s *session.Session, fileDir string,S3_BUCKET string) error {
	//fmt.Println(aws.String(S3_BUCKET))
   // Open the file for use
    file, err := os.Open(fileDir)
    if err != nil {
        return err
    }
    defer file.Close()

    // Get file size and read the file content into a buffer
    fileInfo, _ := file.Stat()
    var size int64 = fileInfo.Size()
    buffer := make([]byte, size)
    file.Read(buffer)

    // Config settings: this is where you choose the bucket, filename, content-type etc.
    // of the file you're uploading.
    _, err = s3.New(s).PutObject(&s3.PutObjectInput{
        Bucket:               aws.String(S3_BUCKET),
        Key:                  aws.String(fileDir),
        ACL:                  aws.String("private"),
        Body:                 bytes.NewReader(buffer),
        ContentLength:        aws.Int64(size),
        ContentType:          aws.String(http.DetectContentType(buffer)),
        ContentDisposition:   aws.String("attachment"),
        ServerSideEncryption: aws.String("AES256"),
    })
    return err
}
