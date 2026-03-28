package tomboy

import "testing"

func TestReadSyncFile(t *testing.T) {
	m, err := ReadSyncFile("_testdata/sync.xml")
	if err != nil {
		t.Fatal(err)
	}
	if 398 != len(m.Items) {
		t.Errorf("got %d expected 398", len(m.Items))
	}
}
