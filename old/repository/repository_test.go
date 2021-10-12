package repository

import (
	"os"
	"testing"
)

func TestRepository(t *testing.T) {

	fileName := "file"

	file, _ := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	repository := NewRepository(file)

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

	repository.Set(tests[0].input, tests[0].want)
	repository.Set(tests[1].input, tests[1].want)
	repository.Set(tests[2].input, tests[2].want)
	repository.Set(tests[3].input, tests[3].want)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := repository.Get(tt.input); got != tt.want {
				if err != nil {
					t.Errorf("Error during repository.Get: %v", err)
					return
				}
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}

	os.Remove(fileName)
}
