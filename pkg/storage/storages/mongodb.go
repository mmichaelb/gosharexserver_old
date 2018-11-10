package storages

import (
	"bytes"
	cryptRand "crypto/rand"
	"fmt"
	"github.com/google/uuid"
	"github.com/kataras/iris/core/errors"
	"github.com/mmichaelb/gosharexserver/pkg/storage"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io"
	"log"
	"math/big"
	mathRand "math/rand"
)

const (
	// Data to generate new call/delete references
	referenceChars        = "abcdefghijklmnopqrstuvxyzABCDEFGHIJKLMNOPQRSTUVXYZ1234567890"
	callReferenceLength   = 6
	deleteReferenceLength = 16
	// MongoDB index names
	referenceIndexName = "reference_index"
	// MongoDB key names
	iDField              = "_id"
	metadataField        = "metadata"
	callReferenceField   = "call_reference"
	deleteReferenceField = "delete_reference"
	authorField          = "author"
	metadataFieldScheme  = "%s.%s"
)

// MongoStorage is the FileStorage implementation using MongoDB GridFS.
type MongoStorage struct {
	// MongoDB Database
	Database *mgo.Database
	// GridFS prefix
	GridFSPrefix string
	// GridFS prefix name
	GridFSChunkSize int
	// internal values
	gridFS *mgo.GridFS
}

// Initialize is the implementation of the FileStorage.Initialize method.
func (mongoStorage *MongoStorage) Initialize() (err error) {
	mongoStorage.gridFS = mongoStorage.Database.GridFS(mongoStorage.GridFSPrefix)
	// check whether db/collection exists
	collectionNames, err := mongoStorage.Database.CollectionNames()
	if err != nil {
		return
	}
	fileCollection := mongoStorage.gridFS.Files
collectionNameCheck:
	for _, collectionName := range collectionNames {
		if collectionName == fileCollection.Name {
			indexes, err := fileCollection.Indexes()
			if err != nil {
				return err
			}
			for _, index := range indexes {
				if index.Name == referenceIndexName {
					break collectionNameCheck
				}
				if err = fileCollection.EnsureIndex(mgo.Index{
					Name: referenceIndexName,
					Key:  []string{callReferenceField, deleteReferenceField},
				}); err != nil {
					return err
				}
			}
		}
	}
	return
}

// Store is the implementation of the FileStorage.Store method.
func (mongoStorage *MongoStorage) Store(entry *storage.Entry) (writer io.WriteCloser, err error) {
	// insert the file details into the collection
	gridFile, err := mongoStorage.gridFS.Create(entry.Filename)
	if err != nil {
		return nil, err
	}
	// set values
	entry.ID = gridFile.Id()
	entry.CallReference = mongoStorage.newCallReference()
	entry.DeleteReference = mongoStorage.newDeleteReference()
	gridFile.SetChunkSize(mongoStorage.GridFSChunkSize)
	gridFile.SetContentType(entry.ContentType)
	gridFile.SetUploadDate(entry.UploadDate)
	gridFile.SetMeta(bson.M{
		authorField:          entry.Author,
		callReferenceField:   entry.CallReference,
		deleteReferenceField: entry.DeleteReference,
	})
	return gridFile, nil
}

// method which randomly creates a new call reference
func (mongoStorage *MongoStorage) newCallReference() (callReference string) {
	buf := bytes.NewBuffer([]byte{})
	randomMaximum := big.NewInt(int64(len(referenceChars)))
	for i := 0; i < callReferenceLength; i++ {
		var randomIndex int
		randomIntIndex, err := cryptRand.Int(cryptRand.Reader, randomMaximum)
		if err != nil {
			log.Printf("Could not get create random call reference with crypto/rand. "+
				"Falling back to (insecure) math/rand package. %T: %v\n", err, err)
			randomIndex = mathRand.Intn(len(referenceChars))
		} else {
			randomIndex = int(randomIntIndex.Int64())
		}
		buf.WriteString(referenceChars[randomIndex : randomIndex+1])
	}
	callReference = buf.String()
	if duplicate, err := mongoStorage.checkForDuplicate(fmt.Sprintf(metadataFieldScheme, metadataField, callReferenceField), callReference); err != nil {
		fmt.Printf("Could not check for duplicates: %v\n", err.Error())
	} else if duplicate {
		return mongoStorage.newCallReference()
	}
	return buf.String()
}

//method which randomly creates a new delete reference
func (mongoStorage *MongoStorage) newDeleteReference() (deleteReference string) {
	buf := bytes.NewBuffer([]byte{})
	randomMaximum := big.NewInt(int64(len(referenceIndexName)))
	for i := 0; i < deleteReferenceLength; i++ {
		var randomIndex int
		randomIntIndex, err := cryptRand.Int(cryptRand.Reader, randomMaximum)
		if err != nil {
			log.Printf("Could not get create random call reference with crypto/rand. "+
				"Falling back to (insecure) math/rand package. %T: %v\n", err, err)
			randomIndex = mathRand.Intn(len(referenceIndexName))
		} else {
			randomIndex = int(randomIntIndex.Int64())
		}
		buf.WriteString(referenceIndexName[randomIndex : randomIndex+1])
	}
	deleteReference = buf.String()
	if duplicate, err := mongoStorage.checkForDuplicate(fmt.Sprintf(metadataFieldScheme, metadataField, deleteReferenceField), deleteReference); err != nil {
		fmt.Printf("Could not check for duplicates: %v\n", err.Error())
	} else if duplicate {
		return mongoStorage.newDeleteReference()
	}
	return deleteReference
}

// checkForDuplicate returns whether the value is already present in the remote database.
func (mongoStorage *MongoStorage) checkForDuplicate(path, value string) (bool, error) {
	count, err := mongoStorage.gridFS.Files.Find(bson.M{path: value}).Count()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Request is the implementation of the Storage.Request method
func (mongoStorage *MongoStorage) Request(callReference string) (*storage.Entry, error) {
	// read result to a simple bson map
	result := &bson.M{}
	// find the entry by its call reference
	if err := mongoStorage.gridFS.Find(bson.M{fmt.Sprintf(metadataFieldScheme, metadataField, callReferenceField): callReference}).One(result); err == mgo.ErrNotFound {
		// return error that entry was not found
		return nil, storage.ErrEntryNotFound
	} else if err != nil {
		// return unwrapped error because something gone horrifically wrong
		return nil, err
	}
	id := storage.ID((*result)[iDField])
	gridFile, err := mongoStorage.gridFS.OpenId(id)
	if err != nil {
		// an error occurred while opening the GridFile
		return nil, err
	}
	// set all entry values except for the reader
	metadata, ok := (*result)[metadataField].(bson.M)
	if !ok {
		return nil, errors.New("could not parse metadata field from GridFS document")
	}
	author, err := uuid.FromBytes(metadata[authorField].([]byte))
	if err != nil {
		return nil, errors.New("could not parse author uuid from database file metadata")
	}
	entry := &storage.Entry{
		ID:            id,
		CallReference: metadata[callReferenceField].(string),
		Author:        author,
		Filename:      gridFile.Name(),
		ContentType:   gridFile.ContentType(),
		UploadDate:    gridFile.UploadDate(),
		Reader:        gridFile,
	}
	return entry, nil
}

// Delete is the implementation of the Storage.Delete method
func (mongoStorage *MongoStorage) Delete(deleteReference string) (err error) {
	// initiate result instance
	result := &bson.M{}
	// find id and return not found if the entry could not be found
	if err = mongoStorage.gridFS.Files.Find(bson.M{fmt.Sprintf(metadataFieldScheme, metadataField, deleteReferenceField): deleteReference}).Select(bson.M{iDField: 1}).One(result); err == mgo.ErrNotFound {
		// return error that entry was not found
		return storage.ErrEntryNotFound
	} else if err != nil {
		// return unwrapped error because something gone horrifically wrong
		return
	}
	// delete entry by its deleteReference
	if err = mongoStorage.gridFS.RemoveId((*result)[iDField]); err == mgo.ErrNotFound {
		// return error that entry was not found
		return storage.ErrEntryNotFound
	} else if err != nil {
		// return unwrapped error because something gone horrifically wrong
		return
	}
	return nil
}

// Close is the implementation of the Storage.Close method
func (mongoStorage *MongoStorage) Close() error {
	// logout from Mongo database and revoke sent credentials
	mongoStorage.Database.Logout()
	return nil
}
