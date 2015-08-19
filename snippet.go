package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Snippet struct {
	Group     string
	Name      string
	GroupPath string
	Path      string
}

func NewSnippet(snippetDir, snippetPath string) *Snippet {
	sl := strings.Split(snippetPath, "/")
	group := strings.Join(sl[:len(sl)-1], "/")
	groupPath := filepath.Join(append([]string{snippetDir}, sl[:len(sl)-1]...)...)
	name := sl[len(sl)-1]
	return &Snippet{
		Group:     group,
		Name:      name,
		GroupPath: groupPath,
		Path:      filepath.Join(groupPath, name),
	}
}

func (s *Snippet) Create() (*os.File, error) {
	if err := s.mkGroupDir(); err != nil {
		return nil, err
	}

	f, err := os.Create(s.Path)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (s *Snippet) Delete() error {
	err := os.Remove(s.Path)
	if err != nil {
		return err
	}
	return nil
}

func (s *Snippet) Write(contents []byte) error {
	f, err := s.Create()
	if err != nil {
		return err
	}

	defer f.Close()

	if _, err := f.Write(contents); err != nil {
		return err
	}

	return nil
}

func (s *Snippet) ReadContents() (string, error) {
	f, err := os.Open(s.Path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (s *Snippet) Exists() bool {
	st, err := os.Stat(s.Path)
	if err != nil {
		return false
	}

	if st.IsDir() {
		return false
	}

	return true
}

func (s *Snippet) mkGroupDir() error {
	if err := os.MkdirAll(s.GroupPath, 0777); err != nil {
		if !os.IsExist(err) {
			return err
		}
	}
	return nil
}
