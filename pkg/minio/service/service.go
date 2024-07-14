package minioservice

import (
	"bytes"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/ardihikaru/example-golang-google-chat-sender/pkg/config"
	"github.com/ardihikaru/example-golang-google-chat-sender/pkg/logger"
	"github.com/ardihikaru/example-golang-google-chat-sender/pkg/minio/utility"
	e "github.com/ardihikaru/example-golang-google-chat-sender/pkg/utils/error"
)

// utility provides the interface for the functionality of MinIO
type utility interface {
	CreateBucketIfNotExist() error

	StoreFileToMinio(objectName, fullFilePath, contentType string) error
	StoreToMinioFromBuffer(fullPath string, byteBuffer []byte) error

	GetObjectByObjectName(filename string) (*bytes.Buffer, error)
}

// Service prepares the interfaces related with this minio service
type Service struct {
	log        *logger.Logger
	cli        *minio.Client
	bucketName string

	// utility extends utility interface{} to implement minioutility.Utility
	// this variable should to be closed from modification: Open/Closed Principle (OCP)
	utility utility
}

// NewService creates a new user service
func NewService(log *logger.Logger, mioCfg config.Minio) *Service {
	minioClient := Connect(mioCfg.Endpoint, mioCfg.AccessKeyId, mioCfg.SecretAccessKey, mioCfg.UseSSL)

	return &Service{
		log:        log,
		cli:        minioClient,
		bucketName: mioCfg.BucketName,
		utility: &minioutility.Utility{
			Log:        log,
			Cli:        minioClient,
			BucketName: mioCfg.BucketName,
		},
	}
}

// Connect initializes minio client object
func Connect(endpoint, accessKeyId, secretAccessKey string, useSSL bool) *minio.Client {
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyId, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		e.FatalOnError(err, "failed to initialize Minio client")
	}

	return minioClient
}

// CreateBucketIfNotExist creates bucket if not exist
func (svc *Service) CreateBucketIfNotExist() error {
	return svc.utility.CreateBucketIfNotExist()
}

// StoreFileToMinio stores file tp minio
func (svc *Service) StoreFileToMinio(objectName, fullFilePath, contentType string) error {
	return svc.utility.StoreFileToMinio(objectName, fullFilePath, contentType)
}

// StoreToMinioFromBuffer stores buffer data to minio
func (svc *Service) StoreToMinioFromBuffer(fullPath string, byteBuffer []byte) error {
	return svc.utility.StoreToMinioFromBuffer(fullPath, byteBuffer)
}

// GetObjectByObjectName gets object data by filename as a byte buffer
func (svc *Service) GetObjectByObjectName(filename string) (*bytes.Buffer, error) {
	return svc.utility.GetObjectByObjectName(filename)
}
