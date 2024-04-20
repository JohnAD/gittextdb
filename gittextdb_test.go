package lib

import (
	"os"
	"testing"
)

const directoryName = "_test"

func TestGitTextDatabase_SETUP(t *testing.T) {
	_ = os.RemoveAll(directoryName)
}

func TestGitTextDatabase_NewDirectoryCreated(t *testing.T) {
	// arrange
	// act
	db, err := CreateDb(directoryName)
	// assert
	if err != nil {
		t.Error(err)
	}
	if !db.StructureExists() {
		t.Errorf("directory %s did not auto-create", directoryName)
	}
}

func TestGitTextDatabase_TEARDOWN(t *testing.T) {
	_ = os.RemoveAll(directoryName)
}
