package utils

import (
	"bytes"
	"fmt"
	"roby-backend-golang/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func UploadImageS3(conf *config.AppConfig, file []byte, content, fileName string, sess *session.Session) (int, error) {
	uploader := s3manager.NewUploader(sess)

	buff := bytes.NewReader(file)
	// Upload the file to S3.
	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(conf.AwsS3.Bucket),
		ACL:         aws.String("public-read"),
		Key:         aws.String(fileName),
		Body:        buff,
		ContentType: aws.String(content),
	})
	if err != nil {
		fmt.Println("Error uploading file to S3: ", err)
		return 0, err
	}

	return buff.Len(), nil
}

func InitAwss3(conf *config.AppConfig) *session.Session {
	defaultResolver := endpoints.DefaultResolver()
	s3CustResolverFn := func(service, region string, optFns ...func(*endpoints.Options)) (endpoints.ResolvedEndpoint, error) {
		if service == "s3" {
			return endpoints.ResolvedEndpoint{
				URL: conf.AwsS3.URL,
			}, nil
		}

		return defaultResolver.EndpointFor(service, region, optFns...)
	}

	// sess, err := session.NewSessionWithOptions(session.Options{
	// 	Config: aws.Config{
	// 		EndpointResolver: endpoints.ResolverFunc(s3CustResolverFn),
	// 		S3ForcePathStyle: aws.Bool(true),
	// 		Credentials:      credentials.NewStaticCredentials(conf.AwsS3.Access, conf.AwsS3.Secret, ""),
	// 		DisableSSL:       aws.Bool(false),
	// 	},
	// })

	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String(conf.AwsS3.Zone),
		EndpointResolver: endpoints.ResolverFunc(s3CustResolverFn),
		S3ForcePathStyle: aws.Bool(true),
		Credentials:      credentials.NewStaticCredentials(conf.AwsS3.Access, conf.AwsS3.Secret, ""),
		DisableSSL:       aws.Bool(false),
	})

	if err != nil {
		panic(fmt.Sprintf("Error creating session: %v", err))
	}

	return sess
}
