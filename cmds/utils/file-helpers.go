package utils

import (
	"os"
)

var maxFiles = 256
var fileSemaphore = make(chan bool, maxFiles)

func init() {
	for i := 0; i < maxFiles; i++ {
		fileSemaphore <- true
	}
}

func CreateAndWriteFile(path string, bytes []byte) error {
	f, err := CreateFile(path)
	if err != nil {
		return err
	}

	_, err = f.Write(bytes)
	if err != nil {
		CloseFile(f)
		return err
	}
	return CloseFile(f)
}

func CreateFile(path string) (*os.File, error) {
	<-fileSemaphore
	f, err := os.Create(path)
	if err != nil {
		fileSemaphore <- true
		return nil, err
	}
	return f, nil
}

func CloseFile(file *os.File) error {
	err := file.Close()
	fileSemaphore <- true
	return err
}

type ConditionalFile struct {
	filename     string
	buffer       []byte
	writeEnabled bool
	file         *os.File
}

func NewConditionalFile(name string) *ConditionalFile {
	return &ConditionalFile{filename: name}
}

func (cf *ConditionalFile) Write(bs []byte, enableWrites bool) (int, error) {
	if enableWrites == true {
		err := cf.enableWrites()
		if err != nil {
			return 0, err
		}
	}

	if cf.writeEnabled {
		return cf.file.Write(bs)
	}

	cf.buffer = append(cf.buffer, bs...)
	return len(bs), nil
}

func (cf *ConditionalFile) WriteString(s string, enableWrites bool) (int, error) {
	if enableWrites == true {
		err := cf.enableWrites()
		if err != nil {
			return 0, err
		}
	}

	if cf.writeEnabled {
		return cf.file.WriteString(s)
	}

	cf.buffer = append(cf.buffer, []byte(s)...)
	return len(s), nil
}

func (cf *ConditionalFile) enableWrites() error {
	if cf.writeEnabled {
		return nil
	}

	cf.writeEnabled = true
	f, err := CreateFile(cf.filename)
	if err != nil {
		return err
	}

	cf.file = f
	_, err = cf.file.Write(cf.buffer)
	if err != nil {
		return err
	}

	cf.buffer = nil

	return nil
}

func (cf *ConditionalFile) Close() error {
	if cf.writeEnabled {
		return CloseFile(cf.file)
	}

	cf.buffer = nil
	return nil
}
