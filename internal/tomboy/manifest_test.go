package tomboy

import "testing"

func TestReadManifestFile(t *testing.T) {
	m, err := ReadManifestFile("_testdata/manifest.xml")
	if err != nil {
		t.Fatal(err)
	}
	if len(m.Items) != 432 {
		t.Errorf("expect 432 got %d", len(m.Items))
	}
	if m.ServerID.String() != "d3d83ee6-78ed-473c-ac7f-678a4d777b41" {
		t.Errorf("expect d3d83ee6-78ed-473c-ac7f-678a4d777b41 got %s", m.ServerID)
	}
}
