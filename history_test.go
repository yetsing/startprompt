package startprompt

import (
	"os"
	"testing"
)

func TestMemHistory(t *testing.T) {
	mem := NewMemHistory()
	if len(mem.GetAll()) != 0 {
		t.Fatalf("want=0, but got=%d", len(mem.GetAll()))
	}

	tests := []string{
		"hello world",
		"hello world2",
		"hello\nworld",
		"hello\nworld2",
	}
	for i, want := range tests {
		mem.Append(want)
		gots := mem.GetAll()
		if len(gots) != i+1 {
			t.Fatalf("want=%d, but got=%d", i+1, len(gots))
		}
		if gots[i] != want {
			t.Fatalf("want=%q, but got=%q", want, gots[i])
		}
		for j := 0; j < i; j++ {
			if gots[j] != tests[j] {
				t.Fatalf("want=%q, but got=%q", tests[j], gots[j])
			}
		}
	}
}

func TestFileHistory(t *testing.T) {
	tempFile, err := os.CreateTemp("", "example")
	if err != nil {
		t.Fatalf("CreateTemp error: %v", err)
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			t.Errorf("Error remove file: %v", err)
		}
	}(tempFile.Name())
	history := NewFileHistory(tempFile.Name())
	//history := NewFileHistory("example")
	if len(history.GetAll()) != 0 {
		t.Fatalf("want=0, but got=%d", len(history.GetAll()))
	}

	tests := []string{
		"hello world",
		"hello world2",
		"hello\nworld",
		"hello\nworld2",
		"hello\nworld\n",
		"hello\nworld\n3",
	}
	for i, want := range tests {
		history.Append(want)
		gots := history.GetAll()
		if len(gots) != i+1 {
			t.Fatalf("want=%d, but got=%d", i+1, len(gots))
		}
		if gots[i] != want {
			t.Fatalf("want=%q, but got=%q", want, gots[i])
		}
		for j := 0; j < i; j++ {
			if gots[j] != tests[j] {
				t.Fatalf("want=%q, but got=%q", tests[j], gots[j])
			}
		}
	}
}
