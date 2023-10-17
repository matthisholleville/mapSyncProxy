package gcs

import (
	"context"
	"fmt"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

func NewClient() *GCSClientWrapper {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		panic(fmt.Errorf("storage.NewClient: %w", err))
	}

	return &GCSClientWrapper{client}
}

func (c *GCSClientWrapper) ListFiles(bucket string) (*[]storage.ObjectAttrs, error) {
	ctx := context.Background()
	files := []storage.ObjectAttrs{}
	items := c.Client.Bucket(bucket).Objects(ctx, nil)
	for {
		attrs, err := items.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return &files, fmt.Errorf("Bucket(%q).Objects: %w", bucket, err)
		}
		files = append(files, *attrs)
	}
	return &files, nil
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
