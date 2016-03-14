package fakedriver

import "os"

//go:generate counterfeiter -o ../volmanfakes/fake_file_system.go . FileSystem

// Interface on file system calls in order to facilitate testing
type FileSystem interface {
	MkdirAll(string, os.FileMode) error
	TempDir() string
}

type realFileSystem struct{}

func NewRealFileSystem() realFileSystem {
	return realFileSystem{}
}

func (f *realFileSystem) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (f *realFileSystem) TempDir() string {
	return os.TempDir()
}
