package gcs

import "cloud.google.com/go/storage"

type GCSClientWrapper struct {
	*storage.Client
}
