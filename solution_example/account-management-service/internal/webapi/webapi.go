package webapi

import "context"

type GDrive interface {
	UploadCSVFile(ctx context.Context, name string, data []byte) (string, error)
	DeleteFile(ctx context.Context, name string) error
	GetAllFilenames(ctx context.Context) ([]string, error)
	IsAvailable() bool
}
