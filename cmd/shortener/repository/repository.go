package repository

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"os"
	"sync"
)

type Service struct {
	Repository *Repository
	BaseURL    string
}

type Repository struct {
	fileName string
	mu       *sync.Mutex // <- нужно для потокобезопасной записи в мапку, об этом в курсе рассказывается далее, пока просто добавь ее)
}

type entity struct {
	key   string
	value string
}

func (r *Repository) Get(key string) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	file, err := os.OpenFile(r.fileName, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return "", err
	}
	scanner := bufio.NewScanner(file)
	decoded := entity{}

	for currentKey := ""; currentKey != key; {
		if !scanner.Scan() {
			return "", scanner.Err()
		}

		if err := gob.NewDecoder(bytes.NewReader(scanner.Bytes())).Decode(&decoded); err != nil {
			return "", err
		}

		currentKey = decoded.key
	}

	return decoded.value, nil
}

func (r *Repository) Set(key string, value string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	file, err := os.OpenFile(r.fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return err
	}

	var encoded bytes.Buffer
	if err := gob.NewEncoder(&encoded).Encode(entity{key, value}); err != nil {
		return err
	}

	data := append(encoded.Bytes(), '\n')

	_, err = file.Write(data)
	return err
}

func NewRepository() *Repository {
	return &Repository{
		fileName: os.Getenv("FILE_STORAGE_PATH"),
		mu:       &sync.Mutex{},
	}
}
