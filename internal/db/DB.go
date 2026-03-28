package db

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/sirupsen/logrus"
)

const (
	File      = "./db/db.json"
	DefaultID = 1
)

type DB struct {
	mu      sync.RWMutex
	resumes map[int]Resume
	logger  *logrus.Logger
}

func NewDB(logger *logrus.Logger) (*DB, error) {
	db := &DB{
		resumes: make(map[int]Resume),
		logger:  logger,
	}
	err := db.load()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (db *DB) load() error {
	db.mu.Lock()
	defer db.mu.Unlock()

	dir := filepath.Dir(File)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	if _, err := os.Stat(File); os.IsNotExist(err) {
		// File doesn't exist, create it with an empty JSON object "{}"
		if err := db.save(); err != nil {
			return err
		}
	} else {
		fileData, err := os.ReadFile(File)
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}

		if len(fileData) == 0 {
			return nil
		}

		return json.Unmarshal(fileData, &db.resumes)
	}

	return nil
}

func (db *DB) save() error {
	db.mu.Lock()
	defer db.mu.Unlock()
	return db.saveLocked()
}

func (db *DB) saveLocked() error {
	data, err := json.MarshalIndent(db.resumes, "", "")
	if err != nil {
		return err
	}
	return os.WriteFile(File, data, os.ModePerm)
}

func (db *DB) UpsertResume(resume Resume) (int, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	if resume.ID <= 0 {
		resume.ID = len(db.resumes) + 1
	}

	db.resumes[resume.ID] = resume

	err := db.saveLocked()
	if err != nil {
		return -1, err
	}

	return resume.ID, nil
}

func (db *DB) GetResume(id int) (Resume, bool) {
	db.mu.Lock()
	defer db.mu.Unlock()

	resume, ok := db.resumes[id]
	return resume, ok
}

func (db *DB) GetDefaultResume() (Resume, bool) {
	return db.GetResume(DefaultID)
}

func (db *DB) DeleteResume(id int) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	delete(db.resumes, id)

	return db.saveLocked()
}
