package storage

import (
	"errors"
	"io"
)

// ErrEntryNotFound is returned by the FileStorage.Request method if the entry could not be found.
var ErrEntryNotFound = errors.New("entry not found")

// FileStorage is an interface which is the scheme to store and request file entries. The implementations can vary.
type FileStorage interface {
	// Initialize is called at the start of the application to e.g. connect to a database or create data folders. It
	// returns an error if something goes wrong.
	Initialize() error
	// Store saves the provided entry and adjusts its the ID and CallReference field values. It returns a writer to
	// write the file data or an error if something goes wrong.
	Store(entry *Entry) (io.WriteCloser, error)
	// Request searches for an entry by the provided callReference which is the substring which is used in the uri.
	// It returns an entry or a specific error (see above) or an unwrapped one if something goes wrong.
	Request(callReference string) (*Entry, error)
	// Delete deletes an entry by the given deleteReference.
	Delete(deleteReference string) error
	// Close shutdowns/closes the FileStorage and allows the storage to exit gracefully. It returns an error if
	// something goes wrong.
	Close() error
}
