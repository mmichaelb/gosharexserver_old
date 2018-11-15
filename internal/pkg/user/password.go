package user

import (
	"bytes"
	"crypto/rand"
	"errors"
	"golang.org/x/crypto/argon2"
)

// HashingAlgorithm is an enumerate containing all valid and implemented hashing methods for password hashing entries.
type HashingAlgorithm int

const (
	// Argon2ID hashing algorithm defaults
	Argon2IDSaltLen int    = 16
	Argon2IDTime    uint32 = 1
	Argon2IDMemory  uint32 = 1024 * 64
	Argon2IDThreads uint8  = 4
	Argon2IDKeyLen  uint32 = 32
)

const (
	// IDs of all hashing algorithms
	HashingArgon2ID HashingAlgorithm = iota
)

// PasswordHash contains basic methods to initiate and validate a password's entry.
type PasswordHash interface {
	// New creates the password entry and returns no error if the creation was successful.
	New(rawPassword []byte) error
	// Check validates the given raw password and returns whether the validation was successful.
	Check(rawPassword []byte) (bool, error)
}

// GetDefaultPasswordHash returns the current default PasswordHash algorithm and its ID.
func GetDefaultPasswordHash() (PasswordHash, HashingAlgorithm) {
	return GetPasswordHashByAlgorithmId(HashingArgon2ID), HashingArgon2ID
}

// GetPasswordHashByAlgorithmId returns a PasswordHash by using the provided HashingAlgorithm.
func GetPasswordHashByAlgorithmId(algorithmId HashingAlgorithm) PasswordHash {
	switch algorithmId {
	case HashingArgon2ID:
		return &Argon2IDPasswordHash{
			Time:    Argon2IDTime,
			Memory:  Argon2IDMemory,
			Threads: Argon2IDThreads,
			KeyLen:  Argon2IDKeyLen,
		}
	default:
		return nil
	}
}

// Argon2IDPasswordHash is the PasswordHash interface implementation by using the Argon2ID algorithm - official winner of the PHC.
type Argon2IDPasswordHash struct {
	// Hash contains the hashed version of the raw password.
	Hash []byte
	// Salt is a randomly generated salt to prevent the entry from getting cracked by e.g. simple rainbow tables.
	Salt []byte
	// Iterations declares the number of iterations the algorithm should run (argon2 synonym: 'time').
	Time uint32
	// Memory declares the number of bytes (part of the costs) used in each round.
	Memory uint32
	// Threads declares the number of threads used in each round.
	Threads uint8
	// KeyLen declares the length of the key length.
	KeyLen uint32
}

// New is the implementation of the PasswordHash#New function.
func (passwordHash *Argon2IDPasswordHash) New(rawPassword []byte) (err error) {
	// generate random salt
	passwordHash.Salt = make([]byte, Argon2IDSaltLen)
	_, err = rand.Read(passwordHash.Salt)
	if err != nil {
		return
	}
	// generate passwordHash hash by using Argon2ID
	passwordHash.Hash = argon2.IDKey(rawPassword, passwordHash.Salt, passwordHash.Time, passwordHash.Memory, passwordHash.Threads, passwordHash.KeyLen)
	return
}

// Check is the implementation of the PasswordHash#Check function.
func (passwordHash *Argon2IDPasswordHash) Check(rawPassword []byte) (ok bool, err error) {
	if err = passwordHash.validateFields(); err != nil {
		return
	}
	rawPasswordHash := argon2.IDKey(rawPassword, passwordHash.Salt, passwordHash.Time, passwordHash.Memory, passwordHash.Threads, passwordHash.KeyLen)
	ok = bytes.Equal(rawPasswordHash, passwordHash.Hash)
	return
}

func (passwordHash *Argon2IDPasswordHash) validateFields() (err error) {
	if passwordHash.Hash == nil {
		err = errors.New("field 'Hash' cannot be unset")
	}
	if passwordHash.Salt == nil {
		err = errors.New("field 'Salt' cannot be unset")
	}
	if passwordHash.Time <= 0 {
		err = errors.New("field 'Time' must be greater than zero")
	}
	if passwordHash.Memory <= 0 {
		err = errors.New("field 'Memory' must be greater than zero")
	}
	if passwordHash.Threads <= 0 {
		err = errors.New("field 'Threads' must be greater than zero")
	}
	if passwordHash.KeyLen <= 0 {
		err = errors.New("field 'KeyLen' must be greater than zero")
	}
	return nil
}
