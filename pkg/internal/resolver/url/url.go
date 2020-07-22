package url

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

// Resolver resolves operator package from URLs pointing to package tarballs.
type Resolver struct {
	URL string
}

// NewResolver creates a new Resolver for a URL.
func NewResolver(url string) Resolver {
	return Resolver{
		URL: url,
	}
}

// Resolve downloads an operator package tarball and extracts it into a file system.
func (r Resolver) Resolve(ctx context.Context) (fs afero.Fs, rem func() error, err error) {
	req, err := http.NewRequestWithContext(ctx, "GET", r.URL, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create HTTP request for %q: %v", r.URL, err)
	}

	log.WithField("url", r.URL).
		Info("Downloading operator tarball")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get HTTP response for %q: %v", r.URL, err)
	}

	tarball, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read HTTP response for %q: %v", r.URL, err)
	}

	fs = afero.NewMemMapFs()

	// Using 'MemMapFs' with the default base path causes all kinds of trouble.
	// To avoid potential issues, all files are created in directories and
	// 'BasePathFs' is used to point to a different base path.
	operatorDir := filepath.Join(string(filepath.Separator), "operator")

	if err := fs.Mkdir(operatorDir, 0755); err != nil {
		return nil, nil, fmt.Errorf("failed to create operator directory: %v", err)
	}

	gzr, err := gzip.NewReader(bytes.NewBuffer(tarball))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to unzip tarball: %v", err)
	}

	defer func() {
		if cerr := gzr.Close(); cerr != nil {
			err = fmt.Errorf("failed to close tarball: %v", cerr)
		}
	}()

	tr := tar.NewReader(gzr)

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, nil, fmt.Errorf("failed to read tarball entry: %v", err)
		}

		switch hdr.Typeflag {
		case tar.TypeReg:
			buf, err := ioutil.ReadAll(tr)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to extract tarball entry: %v", err)
			}

			filename := filepath.Join(operatorDir, hdr.Name)

			// 'WriteFile' won't create directories, let's do this here instead.
			// 'MkdirAll' won't fail if directories already exists which makes
			// it safe to call all the time.
			dir := filepath.Dir(filename)
			if err := fs.MkdirAll(dir, 0755); err != nil {
				return nil, nil, fmt.Errorf("failed to create operator directory %q: %v", dir, err)
			}

			if err := afero.WriteFile(fs, filename, buf, hdr.FileInfo().Mode()); err != nil {
				return nil, nil, fmt.Errorf("failed to write operator file %q: %v", filename, err)
			}
		default:
			continue
		}
	}

	rem = func() error {
		// Nothing to clean up, because we're providing the file system in memory.
		return nil
	}

	return afero.NewBasePathFs(fs, operatorDir), rem, nil
}
