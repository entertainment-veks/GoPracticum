package pkg

import (
	"reflect"
	"testing"
)

func TestGenerateCode(t *testing.T) {
	tests := []struct {
		name string
		want interface{}
	}{
		{
			name: "Generate code lenth test",
			want: 5,
		},
		{
			name: "Generate code valid test",
			want: reflect.TypeOf("string"),
		},
	}

	// t.Run(tests[0].name, func(t *testing.T) {
	// 	if got, _ := GenerateCode(); len(got) != tests[0].want {
	// 		t.Errorf("GenerateCode() len is %v, want %v", got, tests[0].want)
	// 	}
	// })

	t.Run(tests[1].name, func(t *testing.T) {
		generatedCode, _ := GenerateCode()
		if got := reflect.TypeOf(generatedCode); got != tests[1].want {
			t.Errorf("GenerateCode() type is %v, want %v", got, tests[1].want)
		}
	})

}

func TestIsURL(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "Regular URL test",
			input: "https://www.google.com/",
			want:  true,
		},
		{
			name:  "Empty URL test",
			input: "",
			want:  false,
		},
		{
			name:  "Irregular URL test",
			input: "://www.google.com/",
			want:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsURL(tt.input); got != tt.want {
				t.Errorf("IsURL(%v) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
