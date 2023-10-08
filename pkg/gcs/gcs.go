package gcs

import (
	"context"
	"fmt"

	"cloud.google.com/go/storage"
)

func NewClient() *GCSClientWrapper {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		panic(fmt.Errorf("storage.NewClient: %w", err))
	}

	return &GCSClientWrapper{client}
}

// downloadFile downloads an object to a file.
func (c *GCSClientWrapper) DownloadFile(bucket, object string) (*storage.Reader, error) {

	ctx := context.Background()

	rc, err := c.Bucket(bucket).Object(object).NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("Object(%q).NewReader: %w", object, err)
	}
	defer rc.Close()

	return rc, nil

}
