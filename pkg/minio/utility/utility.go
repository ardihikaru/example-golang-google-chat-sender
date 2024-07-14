package minioutility

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"

	"github.com/ardihikaru/example-golang-google-chat-sender/pkg/logger"
)

// Utility is the executable functions in within minio service
type Utility struct {
	Log        *logger.Logger
	Cli        *minio.Client
	BucketName string
}

// CreateBucketIfNotExist creates bucket if not exist
func (mio *Utility) CreateBucketIfNotExist() error {
	var err error

	ctx := context.Background()
	err = mio.Cli.MakeBucket(ctx, mio.BucketName, minio.MakeBucketOptions{})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := mio.Cli.BucketExists(ctx, mio.BucketName)
		if errBucketExists == nil && exists {
			mio.Log.Warn(fmt.Sprintf("bucket %s exists. Do nothing", mio.BucketName))
		} else {
			return fmt.Errorf(fmt.Sprintf("failed to create a new bucket: %s", mio.BucketName))
		}
	}

	mio.Log.Debug(fmt.Sprintf("successfully created %s", mio.BucketName))

	return nil
}

// StoreFileToMinio stores file tp minio
// input param `objectName` expects a value such as: `<owner-name>/<filename>.<ext>` -> `userAbc/file.pdf`
// where it may contain the full path where we want to store it into the storage
// input param `contentType` example: "application/octet-stream", "application/pdf"
func (mio *Utility) StoreFileToMinio(objectName, fullFilePath, contentType string) error {
	_, err := mio.Cli.FPutObject(
		context.Background(), mio.BucketName, objectName, fullFilePath,
		minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return err
	}

	return nil
}

// StoreToMinioFromBuffer stores buffer data to minio
func (mio *Utility) StoreToMinioFromBuffer(fullPath string, byteBuffer []byte) error {
	r := bytes.NewReader(byteBuffer)
	fileSize := int64(len(byteBuffer))

	_, err := mio.Cli.PutObject(context.Background(), mio.BucketName, fullPath, r, fileSize, minio.PutObjectOptions{})
	if err != nil {
		mio.Log.Warn("failed to store data to minio instance")
		return err
	}

	return nil
}

// GetObjectByObjectName gets object data by filename as a byte buffer
func (mio *Utility) GetObjectByObjectName(filename string) (*bytes.Buffer, error) {
	mio.Log.Info(fmt.Sprintf("tries to load following file: [%s]", filename))

	object, err := mio.Cli.GetObject(context.Background(), mio.BucketName, filename, minio.GetObjectOptions{})
	if err != nil {
		mio.Log.Warn("failed to fetch object data from minio instance")
		return nil, err
	}

	buf := new(bytes.Buffer)

	if _, err = io.Copy(buf, object); err != nil {
		mio.Log.Warn("failed to copy fetched Minio Object to the designated buffer")
		return nil, err
	}

	return buf, nil
}
