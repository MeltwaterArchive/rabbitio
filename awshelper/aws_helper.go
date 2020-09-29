// Copyright © 2020 Meltwater
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package awshelper

import (
	"bytes"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var (
	s3Bucket, awsKey, awsSecret string
)

func NewAWSHelper(awsKeyID string, awsSecretAcessKey string, awsS3Bucket string) {
	s3Bucket = awsS3Bucket
	awsKey = awsKeyID
	awsSecret = awsSecretAcessKey
}

func UploadToS3(path string) (string, error) {
	// If AWS credentials were not passed as args, skip uploading to S3
	if len(strings.TrimSpace(awsKey)) == 0 || len(strings.TrimSpace(awsSecret)) == 0 || len(strings.TrimSpace(s3Bucket)) == 0 {
		return "", errors.New("Empty or Invalid AWS Credentials, Skipping S3 Uploads..")
	}

	// All clients require a Session. The Session provides the client with
	// shared configuration such as region, endpoint, and credentials. A
	// Session should be shared where possible to take advantage of
	// configuration and credential caching. See the session package for
	// more information.
	sess := session.Must(session.NewSession(
		&aws.Config{
			Region: aws.String("us-west-2"),
			Credentials: credentials.NewStaticCredentials(
				awsKey,
				awsSecret,
				"", // a token will be created when the session is used.
			),
		}),
	)

	// Open the file for use
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Get file size and read the file content into a buffer
	fileInfo, _ := file.Stat()
	var size int64 = fileInfo.Size()
	buffer := make([]byte, size)
	file.Read(buffer)

	// Create a new instance of the service's client with a Session.
	// Optional aws.Config values can also be provided as variadic arguments
	// to the New function. This option allows you to provide service
	// specific configuration.
	svc := s3.New(sess)

	// Uploads the object to S3. Config settings: this is where you choose the bucket,
	// filename, content-type and storage class of the file you're uploading
	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(s3Bucket),
		Key:                  aws.String(path),
		Body:                 bytes.NewReader(buffer),
		ACL:                  aws.String("private"),
		ContentType:          aws.String(http.DetectContentType(buffer)),
		ContentLength:        aws.Int64(int64(size)),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
		StorageClass:         aws.String("INTELLIGENT_TIERING"),
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == request.CanceledErrorCode {
			// If the SDK can determine the request or retry delay was canceled
			// by a context the CanceledErrorCode error code will be returned.
			log.Fatalf("Upload canceled due to timeout, %s\n", err)
		} else {
			log.Fatalf("Failed to upload object, %s\n", err)
		}
		return "", err
	}

	log.Printf("successfully uploaded file to %s/%s\n", s3Bucket, path)
	return path, nil
}
