package storage

import (
	"github.com/google/uuid"
	"io"
	"time"
)

// ID can vary and is therefore mutable.
type ID interface{}

// ReadCloseSeeker combines both, io.ReaderCloser and io.Seeker into one interface.
type ReadCloseSeeker interface {
	// Allow access via the built in interface and implement the Read, Close and Seek methods.
	io.ReadCloser
	io.Seeker
}

// Entry represents an uploaded file and its metadata in the storage system.
type Entry struct {
	// ID is an identical token which identifies the entry.
	ID ID
	// CallReference is also an identical token but this one is used in the request uri.
	CallReference string
	//DeleteReference is an token used for deleting the entry.
	DeleteReference string
	// Author determines the uploader by providing a unique uuid.
	Author uuid.UUID
	// Filename is the name of the file (contains the application name and date) which is sent with by the ShareX client.
	Filename string
	// ContentType is the MIME-Type of the uploaded file.
	ContentType string
	// UploadDate is the unix timestamp when the file was uploaded.
	UploadDate time.Time
	// ReadCloseSeeker allows to read the image data while controlling the reading start process.
	Reader ReadCloseSeeker
}
