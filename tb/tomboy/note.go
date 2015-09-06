package tomboy

import (
	"encoding/xml"
	"os"
	"time"
)

// Note is ...
type Note struct {
	XMLName xml.Name `xml:"note"`
	Version float64  `xml:"version,attr"`
	Title   string   `xml:"title"`
	Text    struct {
		XMLName xml.Name
		Content string `xml:",innerxml"`
	} `xml:"text>note-content"`
	LastChangeDate         time.Time `xml:"last-change-date"`
	LastMetadataChangeDate time.Time `xml:"last-metadata-change-date"`
	CreateDate             time.Time `xml:"create-date"`
	Tags                   []string  `xml:"tags>tag"` // system:template or system:notebook:golang

	// GUI fileds
	CursorPosition         int  `xml:"cursor-position"`
	SelectionBoundPosition int  `xml:"selection-bound-position"`
	Width                  int  `xml:"width"`
	Height                 int  `xml:"height"`
	X                      int  `xml:"x"`
	Y                      int  `xml:"y"`
	OpenOnStartup          bool `xml:"open-on-startup"`

	ID  UUID
	Rev int
}

// ReadNoteFile read a Note from file
func ReadNoteFile(name string) (*Note, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	dec := xml.NewDecoder(f)

	var n Note
	err = dec.Decode(&n)
	if err != nil {
		return nil, err
	}

	return &n, nil
}
