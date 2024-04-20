/**
 * Author: John Dupuy
 * File: gittextdb.go
 */

package lib

import (
	"errors"
	"os"
)

type GitTextDatabase struct {
	directory string
}

func CreateDb(directory string) (GitTextDatabase, error) {
	db := GitTextDatabase{directory: directory}
	if db.StructureExists() {
		return db, errors.New("directory already exists")
	}
	return db, nil
}

func (db *GitTextDatabase) StructureExists() bool {
	info, err := os.Stat(db.directory)
	if err != nil {
		return false
	}
	if !info.IsDir() {
		return false
	}
	return true
}
