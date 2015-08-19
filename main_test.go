package main

import (
	"fmt"
	"os"
	"testing"
)

func TestLoadConfigExists(t *testing.T) {
	f, err := os.Create(RC_FILE)
	if err != nil {
		t.Errorf("unexpected error happend. %v", err)
	}
	defer f.Close()
	defer os.Remove(RC_FILE)

	expected := &Config{
		SnippetDirectory: "~/.chickenheadsnippetdir",
		Editor:           "emacs",
	}

	f.WriteString(fmt.Sprintf("SnippetDirectory = \"%s\"\n", expected.SnippetDirectory))
	f.WriteString(fmt.Sprintf("Editor = \"%s\"\n", expected.Editor))

	actual := loadConfig(".")

	if actual.String() != expected.String() {
		t.Errorf("unexpected config. expected:%s, actual:%s.", expected, actual)
	}
}

func TestLoadConfigNotExists(t *testing.T) {
	expected := &Config{
		SnippetDirectory: "/home/foo/.chickenhead",
		Editor:           "",
	}

	actual := loadConfig("/home/foo")

	if actual.String() != expected.String() {
		t.Errorf("unexpected config. expected:%s, actual:%s.", expected, actual)
	}
}
