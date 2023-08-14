package startprompt

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
	"time"
)

type History interface {
	GetAll() []string
	Append(s string)
}

type MemHistory struct {
	texts []string
}

func NewMemHistory() *MemHistory {
	return &MemHistory{}
}

func (m *MemHistory) GetAll() []string {
	return m.texts
}

func (m *MemHistory) Append(s string) {
	m.texts = append(m.texts, s)
}

type FileHistory struct {
	MemHistory
	filename string
}

func NewFileHistory(filename string) *FileHistory {
	fh := &FileHistory{
		MemHistory: MemHistory{},
		filename:   filename,
	}
	fh.load()
	return fh
}

func (fh *FileHistory) load() {
	if !fileExists(fh.filename) {
		return
	}
	file, err := os.Open(fh.filename)
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "+") {
			lines = append(lines, line[1:])
		} else {
			if len(lines) > 0 {
				fh.MemHistory.Append(strings.Join(lines, "\n"))
				lines = nil
			}
		}
	}
	if len(lines) > 0 {
		fh.MemHistory.Append(strings.Join(lines, "\n"))
		lines = nil
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}
}

func (fh *FileHistory) Append(s string) {
	fh.MemHistory.Append(s)

	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("\n # %s\n", time.Now().Format(time.RFC3339)))
	for _, line := range strings.Split(s, "\n") {
		buf.WriteString(fmt.Sprintf("+%s\n", line))
	}

	file, err := os.OpenFile(fh.filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)
	_, err = file.Write(buf.Bytes())
	if err != nil {
		panic(err)
	}
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}
