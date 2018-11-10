package user

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// DefaultAuthor holds the default uuid value for e.g. old entries with no author set
var DefaultAuthor, _ = uuid.FromBytes(make([]byte, 16))

// User contains the meta data of a user.
type User struct {
	// UUID is a unique identifier different for each user.
	UUID uuid.UUID `bson:"uuid"`
	// Username is the username of the user (changeable) but still unique.
	Username string `bson:"username,omitempty"`
	// AuthorizationToken is the authorization token which is used by the user to upload files.
	AuthorizationToken AuthorizationToken `bson:"authorization_token"`
	// HashingAlgorithm indicates the hashing algorithm by providing an identical and unique number to match the algorithm.
	HashingAlgorithm HashingAlgorithm `bson:"hashing_algorithm,omitempty"`
	// Hash contains the hashed password with all its needed fields.
	Hash PasswordHash `bson:"hash,omitempty"`
	// internal MongoDB specific values
	collection *mgo.Collection
}

// SetBSON implements the mgo.Setter interface.
func (user *User) SetBSON(raw bson.Raw) (err error) {
	rawWrapped := bson.M{}
	if err = raw.Unmarshal(rawWrapped); err != nil {
		return
	}
	if value, ok := rawWrapped[uuidField]; ok {
		if user.UUID, err = uuid.FromBytes(value.([]byte)); err != nil {
			return
		}
	}
	if value, ok := rawWrapped[usernameField]; ok {
		user.Username = value.(string)
	}
	var unwrappedHashingAlgorithmId interface{}
	var ok bool
	if unwrappedHashingAlgorithmId, ok = rawWrapped[hashingAlgorithmField]; !ok {
		return
	}
	user.HashingAlgorithm = HashingAlgorithm(unwrappedHashingAlgorithmId.(int))
	user.Hash = GetPasswordHashByAlgorithmId(user.HashingAlgorithm)
	if user.Hash == nil {
		return fmt.Errorf("could not find a matching algorithm for id: %d", user.HashingAlgorithm)
	}
	var hashData []byte
	if hashData, err = bson.Marshal(rawWrapped[hashField]); err != nil {
		return
	}
	if err = bson.Unmarshal(hashData, user.Hash); err != nil {
		return
	}
	return
}

var UsernameTakenErr = errors.New("the provided username is already taken")

// CreateNewEntry creates a new user entry in the database and returns the new generated uuid.
func (user *User) CreateNewEntry(rawPassword []byte) (uid uuid.UUID, err error) {
	// validate if collection field is set
	if err = user.validateCollectionField(); err != nil {
		return
	}
	// check if username already exists
	var exists bool
	if exists, err = entryExists(user.collection, bson.M{usernameField: user.Username}); err != nil {
		return
	} else if exists {
		err = UsernameTakenErr
		return
	}
	// find unused uuid
findUUID:
	user.UUID = uuid.New()
	if exists, err = entryExists(user.collection, bson.M{uuidField: user.UUID}); err != nil {
		return
	} else if exists {
		goto findUUID
	}
	// set hash field of user object and initialize with given password
	user.Hash, user.HashingAlgorithm = GetDefaultPasswordHash()
	if err = user.Hash.New(rawPassword); err != nil {
		return
	}
	// insert to MongoDB database
	if err = user.collection.Insert(user); err != nil {
		return
	}
	return
}

// RegenerateAuthorizationToken (re-)generates the authorization token and updates the database entry. To acquire this,
// the uuid field needs to be set. If there is no such user entry an mgo.ErrNotFound is returned.
func (user *User) RegenerateAuthorizationToken() (token AuthorizationToken, err error) {
	// validate if collection field is set
	if err = user.validateCollectionField(); err != nil {
		return
	}
	// generate token and check if it already exists
tokenGeneration:
	token = AuthorizationToken(make([]byte, viper.GetInt("webserver.authorization_token_length")))
	var exists bool
	if exists, err = entryExists(user.collection, bson.M{authorizationTokenField: token}); err != nil {
		return
	} else if exists {
		goto tokenGeneration
	}
	if err = user.collection.Update(bson.M{uuidField: user.UUID}, bson.M{"$set": bson.M{authorizationTokenField: token}}); err != nil {
		return
	}
	user.AuthorizationToken = token
	return
}

// UpdateUsername updates the user's name. To acquire this, the uuid field needs to be set. If there is no such user
// entry an mgo.ErrNotFound is returned.
func (user *User) UpdateUsername(username string) (err error) {
	// validate if collection field is set
	if err = user.validateCollectionField(); err != nil {
		return
	}
	if err = user.collection.Update(bson.M{uuidField: user.UUID}, bson.M{"$set": bson.M{usernameField: username}}); err != nil {
		return
	}
	user.Username = username
	return
}

// UpdatePassword updates the user's password.  To acquire this, the uuid field needs to be set. If there is no entry to
// update an mgo.ErrNotFound is returned.
func (user *User) UpdatePassword(rawPassword []byte) (err error) {
	// validate if collection field is set
	if err = user.validateCollectionField(); err != nil {
		return
	}
	user.Hash, user.HashingAlgorithm = GetDefaultPasswordHash()
	if err = user.Hash.New(rawPassword); err != nil {
		return
	}
	if err = user.collection.Update(bson.M{uuidField: user.UUID}, bson.M{"$set": bson.M{
		hashingAlgorithmField: user.HashingAlgorithm,
		hashField:             user.Hash,
	}}); err != nil {
		return
	}
	return
}

// Delete deletes the user entry. To acquire this, the uuid field needs to be set. If there is no such user entry an
// mgo.ErrNotFound is returned.
func (user *User) Delete() (err error) {
	// validate if collection field is set
	if err = user.validateCollectionField(); err != nil {
		return
	}
	if err = user.collection.Remove(bson.M{uuidField: user.UUID}); err != nil {
		return
	}
	return
}

// internal utility functions

func (user *User) validateCollectionField() (err error) {
	if user.collection != nil {
		err = errors.New("internal 'collection' field cannot be unset")
	}
	return
}

func entryExists(collection *mgo.Collection, keyValue bson.M) (bool, error) {
	if count, err := collection.Find(keyValue).Count(); err != nil {
		return false, err
	} else if count > 0 {
		return true, err
	}
	return false, nil
}
