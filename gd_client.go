package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
)

type Executor[T any] func(ctx context.Context, opts *Options) (T, error)

type Options struct {
	Query     string
	Fields    string
	PageToken string
	PageSize  int64
	FolderID  string
	DriveID   string
}

type Option func(*Options)

func Query(query string) Option {
	return func(o *Options) {
		o.Query = query
	}
}

func Fields(fields string) Option {
	return func(o *Options) {
		o.Fields = fields
	}
}

func PageToken(pageToken string) Option {
	return func(o *Options) {
		o.PageToken = pageToken
	}
}

func PageSize(pageSize int64) Option {
	return func(o *Options) {
		o.PageSize = pageSize
	}
}

func FolderID(folderID string) Option {
	return func(o *Options) {
		o.FolderID = folderID
	}
}

func DriveID(driveID string) Option {
	return func(o *Options) {
		o.DriveID = driveID
	}
}

func NewOptions(opts ...Option) *Options {
	options := Options{}
	for _, o := range opts {
		o(&options)
	}
	return &options
}

type GoogleDriveClient struct {
	srv *drive.Service
}

func NewGoogleDriveClient(ctx context.Context) *GoogleDriveClient {
	srv, err := drive.NewService(ctx)
	if err != nil {
		panic(err)
	}
	return &GoogleDriveClient{
		srv: srv,
	}
}

func Retry[T any](executor Executor[T], retries int, delay time.Duration) Executor[T] {
	return func(ctx context.Context, opts *Options) (T, error) {
		var zero T
		for r := 0; ; r++ {
			res, err := executor(ctx, opts)
			if err == nil || r >= retries {
				return res, err
			}

			log.Printf("Attempt %d failed; retrying in %v", r+1, delay)

			select {
			case <-time.After(delay):
			case <-ctx.Done():
				return zero, ctx.Err()
			}
		}
	}
}

func (c GoogleDriveClient) GetFileName(ctx context.Context, opts *Options) (string, error) {
	var res *drive.File
	var err error
	if opts.DriveID == "" {
		res, err = c.srv.Files.Get(opts.FolderID).Fields("name").Do()
	} else {
		res, err = c.srv.Files.Get(opts.FolderID).SupportsAllDrives(true).Fields("name").Do()
	}
	if err != nil {
		return "", err
	}
	return res.Name, nil
}

func (c GoogleDriveClient) ListFiles(ctx context.Context, opts *Options) (*drive.FileList, error) {
	flc := c.srv.Files.List()
	if opts.DriveID != "" {
		flc.
			Corpora("drive").
			IncludeItemsFromAllDrives(true).
			SupportsAllDrives(true).
			DriveId(opts.DriveID)
	}
	res, err := flc.
		Q(opts.Query).
		Fields(googleapi.Field(opts.Fields)).
		PageToken(opts.PageToken).
		PageSize(opts.PageSize).
		Do()
	if err != nil {
		return nil, err
	}

	return res, nil
}
