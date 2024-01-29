package service

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestPasswordEncrypt(t *testing.T) {
	password := []byte("123123!123#passwd")
	hashCode, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	assert.NoError(t, err)
	println(string(hashCode))
	err = bcrypt.CompareHashAndPassword(hashCode, []byte("123123!123#passwd"))
	assert.NoError(t, err)
}
