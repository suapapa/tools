package tomboy

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"strconv"
)

// Sync is ...
type Sync struct {
	XMLName  xml.Name   `xml:"sync"`
	Revision int        `xml:"revision,attr"`
	ServerID UUID       `xml:"server-id,attr"`
	Items    []SyncItem `xml:"note"`
}

// SyncItem is ...
type SyncItem struct {
	ID       UUID `xml:"id,attr"`
	Revision int  `xml:"rev,attr"`
}

// Path returns file path for NoteItem
func (n SyncItem) Path() string {
	first := strconv.Itoa(n.Revision / 100)
	second := strconv.Itoa(n.Revision)
	base := n.ID.String() + ".note"

	return filepath.Join(first, second, base)
}

// ReadSyncFile read Sync from file
func ReadSyncFile(name string) (*Sync, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	dec := xml.NewDecoder(f)

	var m Sync
	err = dec.Decode(&m)
	if err != nil {
		return nil, err
	}

	return &m, nil
}
