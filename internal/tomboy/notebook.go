package tomboy

import (
	"path/filepath"
	"sort"
)

// Notebook is list of Notes which sorted by LastChangeDate
type Notebook []*Note

func (b Notebook) Len() int {
	return len(b)
}

func (b Notebook) Less(i, j int) bool {
	return b[i].LastChangeDate.After(b[j].LastChangeDate)
}

func (b Notebook) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

// MakeNotebookFromFileSystemSync makes Notebook from file system sync directory
func MakeNotebookFromFileSystemSync(root string) (Notebook, error) {
	sync, err := ReadSyncFile(filepath.Join(root, "manifest.xml"))
	if err != nil {
		return nil, err
	}

	var notes Notebook
	for _, e := range sync.Items {
		n, err := ReadNoteFile(filepath.Join(root, e.Path()))
		if err != nil {
			return nil, err
		}
		n.ID = e.ID
		n.Rev = e.Revision
		notes = append(notes, n)
	}

	sort.Sort(notes)
	return notes, nil
}
