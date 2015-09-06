package tomboy

import (
	"encoding/xml"
	"os"
	"time"
)

// Manifest is ...
type Manifest struct {
	XMLName  xml.Name       `xml:"manifest"`
	Date     time.Time      `xml:"last-sync-date"`
	Revision int            `xml:"last-sync-rev"`
	ServerID UUID           `xml:"server-id"`
	Items    []NoteRevision `xml:"note-revisions>note"`
}

// NoteRevision is ...
type NoteRevision struct {
	ID       UUID `xml:"guid,attr"`
	Revision int  `xml:"lastest-revision,attr"`
}

// ReadManifestFile read Manifest from file
func ReadManifestFile(name string) (*Manifest, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	dec := xml.NewDecoder(f)

	var m Manifest
	err = dec.Decode(&m)
	if err != nil {
		return nil, err
	}

	return &m, nil
}
