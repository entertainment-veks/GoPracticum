package repository

import "sync"

type Repository struct {
	data map[string]string
	mu   *sync.Mutex // <- нужно для потокобезопасной записи в мапку, об этом в курсе рассказывается далее, пока просто добавь ее)
}

func (r *Repository) Get(key string) string {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.data[key]
}

func (r *Repository) Set(key, value string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.data[key] = value
}

func NewRepository() *Repository {
	return &Repository{
		data: make(map[string]string),
	}
}
