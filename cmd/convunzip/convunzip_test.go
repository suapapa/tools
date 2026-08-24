package main

import (
	"archive/zip"
	"os"
	"path/filepath"
	"testing"

	"golang.org/x/text/transform"
)

func TestSafeExtractPath(t *testing.T) {
	dest := t.TempDir()

	tests := []struct {
		name    string
		in      string
		wantErr bool
	}{
		{name: "relative file", in: "dir/file.txt", wantErr: false},
		{name: "dot parent", in: "../evil.txt", wantErr: true},
		{name: "nested parent", in: "foo/../../evil.txt", wantErr: true},
		{name: "absolute unix", in: "/tmp/evil.txt", wantErr: true},
		{name: "backslash parent", in: `..\evil.txt`, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := safeExtractPath(dest, tt.in)
			if tt.wantErr && err == nil {
				t.Fatalf("safeExtractPath(%q) = nil, want error", tt.in)
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("safeExtractPath(%q) = %v, want nil", tt.in, err)
			}
		})
	}
}

func TestConvUnzipRejectsZipSlip(t *testing.T) {
	dir := t.TempDir()
	zipPath := filepath.Join(dir, "slip.zip")

	f, err := os.Create(zipPath)
	if err != nil {
		t.Fatal(err)
	}
	zw := zip.NewWriter(f)
	w, err := zw.Create("../evil.txt")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := w.Write([]byte("pwned")); err != nil {
		t.Fatal(err)
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}

	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(wd) })
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	err = convUnzip(zipPath, transform.Nop, false)
	if err == nil {
		t.Fatal("convUnzip accepted zip-slip archive, want error")
	}

	if _, err := os.Stat(filepath.Join(filepath.Dir(dir), "evil.txt")); err == nil {
		t.Fatal("zip-slip wrote a file outside the destination directory")
	}
}
