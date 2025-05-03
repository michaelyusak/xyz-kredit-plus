package repository

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/michaelyusak/xyz-kredit-plus/entity"
)

type mediaRepositoryLocal struct {
	storagePath string
}

func NewMediaRepositoryLocal(storagePath string) *mediaRepositoryLocal {
	return &mediaRepositoryLocal{
		storagePath: storagePath,
	}
}

type MediaOpt struct {
	Extension string
	Key       string
	Bytes     []byte
	File      *entity.Media
}

func (r *mediaRepositoryLocal) Store(ctx context.Context, media MediaOpt) error {
	if media.Key == "" || media.Extension == "" {
		return fmt.Errorf("[local_media_repository][Store] media key and extension is required")
	}

	fullPath := r.storagePath + media.Key + media.Extension

	switch {
	case len(media.Bytes) > 0:
		if err := os.WriteFile(fullPath, media.Bytes, 0644); err != nil {
			return fmt.Errorf("[local_media_repository][Store][os.WriteFile] error: %w", err)
		}
		return nil

	case media.File != nil && media.File.File != nil:
		dstFile, err := os.Create(fullPath)
		if err != nil {
			return fmt.Errorf("[local_media_repository][Store][os.Create] error: %w", err)
		}
		defer dstFile.Close()

		if _, err := io.Copy(dstFile, media.File.File); err != nil {
			return fmt.Errorf("[local_media_repository][Store][io.Copy] error: %w", err)
		}
		return nil

	default:
		return fmt.Errorf("[local_media_repository][Store] no file data provided")
	}
}
