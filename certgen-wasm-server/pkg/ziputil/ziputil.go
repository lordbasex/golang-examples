package ziputil

import (
	"archive/zip"
	"bytes"
	"fmt"
	"sort"
	"strings"
)

// ZipFiles creates a zip archive from a map of filename -> bytes.
// Filenames should use forward slashes (example: "certs/server.crt").
func ZipFiles(files map[string][]byte) ([]byte, error) {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)

	// Stable output ordering helps with debugging / reproducible builds.
	names := make([]string, 0, len(files))
	for name := range files {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		data := files[name]
		if name == "" {
			_ = zw.Close()
			return nil, fmt.Errorf("zip entry name cannot be empty")
		}

		header := &zip.FileHeader{
			Name:   name,
			Method: zip.Deflate,
		}
		// Mark shell scripts as executable inside the zip.
		if strings.HasSuffix(name, ".sh") {
			header.SetMode(0o755)
		} else {
			header.SetMode(0o644)
		}

		w, err := zw.CreateHeader(header)
		if err != nil {
			_ = zw.Close()
			return nil, err
		}
		if _, err := w.Write(data); err != nil {
			_ = zw.Close()
			return nil, err
		}
	}

	if err := zw.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

