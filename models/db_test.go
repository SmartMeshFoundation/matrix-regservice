package models

import (
	"os"
	"path"
	"testing"
)

var dbPath string

func setupDb(t *testing.T) {
	dbPath = path.Join(os.TempDir(), "testxxxx.db")
	os.Remove(dbPath)
	os.Remove(dbPath + ".lock")
	err := SetupDB(dbPath)
	if err != nil {
		t.Error(err)
		return
	}
	return
}

func TestIsUserAlreadyExists(t *testing.T) {
	user := &User{
		LocalPart:   "aaa",
		DisplayName: "aaa",
		Password:    "123",
	}
	setupDb(t)
	if IsUserAlreadyExists(user.LocalPart) {
		t.Error("must not exist")
		return
	}
	err := NewUser(user.LocalPart, user.DisplayName, user.Password)
	if err != nil {
		t.Error(err)
		return
	}
	if !IsUserAlreadyExists(user.LocalPart) {
		t.Error("must exist")
		return
	}
}
