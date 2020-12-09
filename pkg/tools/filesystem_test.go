package tools

import (
	"os"
	"testing"
)

const fileName = "c:/temp/testfile.txt"

func TestFileExists(t *testing.T) {
	if FileExists(fileName) {
		if err := os.Remove(fileName); err != nil {
			t.Errorf("cann't delete file  %v", fileName)
		}
	}

	if FileExists(fileName) {
		t.Errorf("non existing file %v dedected", fileName)
	}

	CreateFile(fileName)
	if !FileExists(fileName) {
		t.Errorf("file %v doesn't exists", fileName)
	}

	if err := os.Remove(fileName); err != nil {
		t.Errorf("cann't delete file  %v", fileName)
	}
}
