package repository

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"go_practicum/internal/util"
	"os"
	"sync"
)

type Repository struct {
	file *util.FilerReader
	mu   *sync.Mutex // <- нужно для потокобезопасной записи в мапку, об этом в курсе рассказывается далее, пока просто добавь ее)
}

type Entity struct {
	Key   string
	Value string
}

func (r *Repository) Get(key string) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	scanner := bufio.NewScanner(r.file)
	decoded := Entity{}

	for currentKey := ""; currentKey != key; currentKey = decoded.Key {
		if !scanner.Scan() {
			return "", scanner.Err()
		}

		if err := gob.NewDecoder(bytes.NewReader(scanner.Bytes())).Decode(&decoded); err != nil {
			return "", err
		}
	}

	return decoded.Value, nil
}

func (r *Repository) Set(key string, value string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	var encoded bytes.Buffer
	if err := gob.NewEncoder(&encoded).Encode(Entity{key, value}); err != nil {
		return err
	}

	data := append(encoded.Bytes(), '\n')

	_, err := r.file.Write(data)
	return err
}

func NewRepository(file *os.File) *Repository {
	return &Repository{
		file: &util.FilerReader{
			file,
		},
		mu: &sync.Mutex{},
	}
}
