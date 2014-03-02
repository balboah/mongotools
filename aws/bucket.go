package aws

import (
	"io"
	"net/http"
)

type BucketObject struct{}

func (o BucketObject) Write(b []byte) (int, error) {

}

func CreateObject(bucketName, objectName string, client *http.Request) (io.WriteCloser, error) {

}
