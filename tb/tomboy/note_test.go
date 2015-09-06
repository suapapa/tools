package tomboy

import (
	"fmt"
	"testing"
)

func TestReadNoteFile(t *testing.T) {
	n, err := ReadNoteFile("_testdata/template.note")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(n.Tags)
}
