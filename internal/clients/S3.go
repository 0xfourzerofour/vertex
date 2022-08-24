package clients

import (
	"sync"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-xray-sdk-go/xray"
)

var (
	onceS3 sync.Once
	s3conn *s3.S3
)

func initialiseS3() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	s3conn = s3.New(sess)
}

func S3() *s3.S3 {
	onceS3.Do(initialiseS3)
	xray.AWS(s3conn.Client)
	return s3conn
}
