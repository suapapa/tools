package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/transform"
)

const (
	maxZipMemberSize = 512 << 20 // 512 MiB
	maxZipTotalSize  = 2 << 30   // 2 GiB
	dirPerm          = 0o755
	filePerm         = 0o644
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	dryRun := flag.Bool("n", false, "dry-run")
	flag.Parse()

	if flag.NArg() == 0 {
		fmt.Printf("Usage: %s [-n] zipfile ...\n", os.Args[0])
		return nil
	}

	euckrDec := korean.EUCKR.NewDecoder()

	for _, zipFilename := range flag.Args() {
		fmt.Printf("Unzipping %s ...\n", zipFilename)
		if err := convUnzip(zipFilename, euckrDec, *dryRun); err != nil {
			log.Printf("failed to unzip %s: %v", zipFilename, err)
		}
	}
	return nil
}

func convUnzip(fileName string, dec transform.Transformer, dryRun bool) error {
	r, err := zip.OpenReader(fileName)
	if err != nil {
		return fmt.Errorf("failed to open zip: %w", err)
	}
	defer r.Close()

	destDir, err := filepath.Abs(".")
	if err != nil {
		return fmt.Errorf("failed to resolve destination: %w", err)
	}

	var total uint64
	for _, f := range r.File {
		fixed, _, err := transform.String(dec, f.Name)
		if err != nil {
			log.Printf("failed to transform filename %s: %v", f.Name, err)
			fixed = f.Name
		}
		fixed = strings.ReplaceAll(fixed, "\\", "/")

		fmt.Println(fixed)

		if dryRun {
			continue
		}

		target, err := safeExtractPath(destDir, fixed)
		if err != nil {
			return err
		}

		info := f.FileInfo()
		if info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("refusing symlink zip entry %q", fixed)
		}
		if !info.IsDir() && !info.Mode().IsRegular() && info.Mode()&os.ModeType != 0 {
			return fmt.Errorf("refusing non-regular zip entry %q", fixed)
		}

		if info.IsDir() || strings.HasSuffix(fixed, "/") {
			if err := os.MkdirAll(target, dirPerm); err != nil {
				return fmt.Errorf("failed to create directory: %w", err)
			}
			if err := rejectSymlinkPath(destDir, target); err != nil {
				return err
			}
			continue
		}

		if f.UncompressedSize64 > maxZipMemberSize {
			return fmt.Errorf("zip entry %q exceeds size limit", fixed)
		}
		if total+f.UncompressedSize64 > maxZipTotalSize {
			return fmt.Errorf("zip total uncompressed size exceeds limit")
		}

		if err := os.MkdirAll(filepath.Dir(target), dirPerm); err != nil {
			return fmt.Errorf("failed to create parent directory: %w", err)
		}
		if err := rejectSymlinkPath(destDir, target); err != nil {
			return err
		}

		n, err := extractFile(f, target)
		if err != nil {
			return err
		}
		total += uint64(n)
	}
	return nil
}

// safeExtractPath joins name under destDir and rejects Zip Slip / absolute paths.
func safeExtractPath(destDir, name string) (string, error) {
	name = strings.ReplaceAll(name, "\\", "/")
	name = filepath.Clean(filepath.FromSlash(name))
	if name == "." {
		return destDir, nil
	}
	if !filepath.IsLocal(name) {
		return "", fmt.Errorf("refusing unsafe zip path %q", name)
	}
	target := filepath.Join(destDir, name)
	rel, err := filepath.Rel(destDir, target)
	if err != nil || !filepath.IsLocal(rel) {
		return "", fmt.Errorf("refusing unsafe zip path %q", name)
	}
	return target, nil
}

func rejectSymlinkPath(root, target string) error {
	rel, err := filepath.Rel(root, target)
	if err != nil {
		return fmt.Errorf("refusing path outside destination: %w", err)
	}
	cur := root
	for _, part := range strings.Split(rel, string(os.PathSeparator)) {
		if part == "." || part == "" {
			continue
		}
		cur = filepath.Join(cur, part)
		fi, err := os.Lstat(cur)
		if os.IsNotExist(err) {
			return nil
		}
		if err != nil {
			return err
		}
		if fi.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("refusing to write through symlink %q", cur)
		}
	}
	return nil
}

func extractFile(f *zip.File, dest string) (int64, error) {
	rc, err := f.Open()
	if err != nil {
		return 0, fmt.Errorf("failed to open zip file member: %w", err)
	}
	defer rc.Close()

	out, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, filePerm)
	if err != nil {
		return 0, fmt.Errorf("failed to create output file: %w", err)
	}
	defer out.Close()

	limit := int64(f.UncompressedSize64)
	if limit <= 0 || limit > maxZipMemberSize {
		limit = maxZipMemberSize
	}
	n, err := io.Copy(out, io.LimitReader(rc, limit+1))
	if err != nil {
		return n, fmt.Errorf("failed to copy content: %w", err)
	}
	if n > limit {
		return n, fmt.Errorf("zip entry %q exceeds declared size", f.Name)
	}
	return n, nil
}
