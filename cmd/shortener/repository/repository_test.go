package repository

import (
	"os"
	"testing"
)

func TestRepository(t *testing.T) {
	repository := NewRepository()
	fileName := "file"
	os.Setenv("FILE_STORAGE_PATH", fileName)

	repository.Set("aaaaa", "https://www.google.com/")
	repository.Set("bbbbb", "http://localhost:8080/")
	repository.Set("ccccc", "https://practicum.yandex.ru/")
	repository.Set("ddddd", "https://translate.google.com/")

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "Regular test #1",
			input: "aaaaa",
			want:  "https://www.google.com/",
		},
		{
			name:  "Regular test #2",
			input: "bbbbb",
			want:  "http://localhost:8080/",
		},
		{
			name:  "Regular test #3",
			input: "ccccc",
			want:  "https://practicum.yandex.ru/",
		},
		{
			name:  "Regular test #4",
			input: "ddddd",
			want:  "https://translate.google.com/",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := repository.Get(tt.input); got != tt.want {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}

	//os.Remove(fileName)
}
