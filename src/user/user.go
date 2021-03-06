package user

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"errors"
)

// common interface for user
type user interface {
	auth(input string) bool
}

// plain password
type plainUser struct {
	token string
}

func (u *plainUser) auth(input string) bool {
	return u.token == input
}

func newPlainUser(pass string) *plainUser {
	return &plainUser{pass}
}

// base64 password
type base64User struct {
	token string
}

func (u *base64User) auth(input string) bool {
	inputToken := base64.StdEncoding.EncodeToString([]byte(input))
	return u.token == inputToken
}

func newBase64User(encPass string) *base64User {
	return &base64User{encPass}
}

// md5 hashed password
type md5User struct {
	token [md5.Size]byte
}

func (u *md5User) auth(input string) bool {
	inputToken := md5.Sum([]byte(input))
	return u.token == inputToken
}

func newMd5User(encPass string) (*md5User, error) {
	tokenSlice, err := hex.DecodeString(encPass)
	if err != nil {
		return nil, err
	}
	if len(tokenSlice) != md5.Size {
		return nil, errors.New("unrecognized hash")
	}
	token := [md5.Size]byte{}
	copy(token[:], tokenSlice)
	return &md5User{token}, nil
}

// sha1 hashed password
type sha1User struct {
	token [sha1.Size]byte
}

func (u *sha1User) auth(input string) bool {
	inputToken := sha1.Sum([]byte(input))
	return u.token == inputToken
}

func newSha1User(encPass string) (*sha1User, error) {
	tokenSlice, err := hex.DecodeString(encPass)
	if err != nil {
		return nil, err
	}
	if len(tokenSlice) != sha1.Size {
		return nil, errors.New("unrecognized hash")
	}
	token := [sha1.Size]byte{}
	copy(token[:], tokenSlice)
	return &sha1User{token}, nil
}

// sha256 hashed password
type sha256User struct {
	token [sha256.Size]byte
}

func (u *sha256User) auth(input string) bool {
	inputToken := sha256.Sum256([]byte(input))
	return u.token == inputToken
}

func newSha256User(encPass string) (*sha256User, error) {
	tokenSlice, err := hex.DecodeString(encPass)
	if err != nil {
		return nil, err
	}
	if len(tokenSlice) != sha256.Size {
		return nil, errors.New("unrecognized hash")
	}
	token := [sha256.Size]byte{}
	copy(token[:], tokenSlice)
	return &sha256User{token}, nil
}

// sha512 hashed password
type sha512User struct {
	token [sha512.Size]byte
}

func (u *sha512User) auth(input string) bool {
	inputToken := sha512.Sum512([]byte(input))
	return u.token == inputToken
}

func newSha512User(encPass string) (*sha512User, error) {
	tokenSlice, err := hex.DecodeString(encPass)
	if err != nil {
		return nil, err
	}
	if len(tokenSlice) != sha512.Size {
		return nil, errors.New("unrecognized hash")
	}
	token := [sha512.Size]byte{}
	copy(token[:], tokenSlice)
	return &sha512User{token}, nil
}
