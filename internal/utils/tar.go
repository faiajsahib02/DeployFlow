package utils

import (
	"archive/tar"
	"bytes"
	"io"
	"time"
)

// CreateTarArchive creates an in-memory tarball from a map of filenames to content
func CreateTarArchive(files map[string]string) (io.Reader, error) {
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)

	for name, body := range files {
		hdr := &tar.Header{
			Name:    name,
			Mode:    0600,
			Size:    int64(len(body)),
			ModTime: time.Now(),
		}
		if err := tw.WriteHeader(hdr); err != nil {
			return nil, err
		}
		if _, err := tw.Write([]byte(body)); err != nil {
			return nil, err
		}
	}
	if err := tw.Close(); err != nil {
		return nil, err
	}
	return buf, nil
}