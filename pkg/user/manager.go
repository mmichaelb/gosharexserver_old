package user

import (
	"github.com/google/uuid"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	collectionName        = "users"
	searchIndexName       = "search_index"
	uuidField             = "uuid"
	usernameField         = "username"
	hashingAlgorithmField = "hashing_algorithm"
	hashField             = "hash"
)

// Manager manages the users by connecting to and using a MongoDB server.
type Manager struct {
	// MongoDB Database
	Database *mgo.Database
	// internal field for MongoDB management
	collection *mgo.Collection
}

// InitializeCollection initializes the collection.
func (manager *Manager) InitializeCollection() (err error) {
	manager.collection = manager.Database.C(collectionName)
	var indexes []mgo.Index
	if indexes, err = manager.collection.Indexes(); err != nil || len(indexes) > 0 {
		return
	}
	if err = manager.collection.Create(&mgo.CollectionInfo{
		Validator: bson.M{
			uuidField:             bson.M{"$exists": true},
			usernameField:         bson.M{"$exists": true},
			hashingAlgorithmField: bson.M{"$exists": true},
			hashField:             bson.M{"$exists": true},
		},
		ValidationLevel:  "strict",
		ValidationAction: "error",
	}); err != nil {
		return
	}
	if err = manager.collection.EnsureIndex(mgo.Index{
		Name:       searchIndexName,
		Key:        []string{uuidField, usernameField},
		Background: false,
		Unique:     true,
	}); err != nil {
		return
	}
	return
}

// GetNewUserInstance returns a new user with only the name and internal collection field set. This will not effect the
// Mongo database.
func (manager *Manager) GetNewUserInstance(username string) (user *User) {
	return &User{
		Username:   username,
		collection: manager.collection,
	}
}

// LoadUser loads the user by using the given uuid from the Mongo database and returns the parsed object.
func (manager *Manager) LoadUser(uuid uuid.UUID) (user *User, err error) {
	user = &User{}
	if err = manager.collection.Find(bson.M{uuidField: uuid}).One(user); err != nil {
		return
	}
	return
}
