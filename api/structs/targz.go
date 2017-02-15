package structs

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// TarGzipWriter type provides methods to wrap data in a Gzip tar archive
type TarGzipWriter struct {
	gw *gzip.Writer
	tw *tar.Writer
}

// TarGzipReader type provides methods to read data in a Gzip tar archive
type TarGzipReader struct {
	gr *gzip.Reader
	tr *tar.Reader
}

// NewTarGzipWriter creates a new TarGzipWriter
func NewTarGzipWriter(w io.Writer) *TarGzipWriter {
	tgw := &TarGzipWriter{}
	tgw.gw = gzip.NewWriter(w)
	tgw.tw = tar.NewWriter(tgw.gw)

	return tgw
}

// NewTarGzipReader creates a new NewTarGzipReader
func NewTarGzipReader(r io.Reader) (*TarGzipReader, error) {
	gz, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &TarGzipReader{
		gr: gz,
		tr: tar.NewReader(gz),
	}, nil
}

// Write writes the data to a Gzip tar archive
func (tgw *TarGzipWriter) Write(b []byte) (int, error) {
	return tgw.tw.Write(b)
}

// WriteHeader writes the header information for a tar archive
func (tgw *TarGzipWriter) WriteHeader(header *tar.Header) error {
	return tgw.tw.WriteHeader(header)
}

// WriteData is a helper method that writes the data and it's header
func (tgw *TarGzipWriter) WriteData(data []byte, header *tar.Header) error {
	if err := tgw.WriteHeader(header); err != nil {
		return err
	}

	if _, err := tgw.Write(data); err != nil {
		return err
	}

	return nil
}

// WriteDirectory walks the given directory and writes all files (keeping the path) to the archive
func (tgw *TarGzipWriter) WriteDirectory(dir string) error {
	files := []string{}

	if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info != nil && !info.IsDir() {
			files = append(files, path)
		}
		return nil
	}); err != nil {
		return err
	}

	for _, f := range files {
		file, err := os.Open(f)
		if err != nil {
			return nil

		}
		defer file.Close()

		stat, err := os.Stat(f)
		if err != nil {
			return err
		}

		header := &tar.Header{
			Typeflag: tar.TypeReg,
			Name:     strings.TrimPrefix(f, path.Clean(dir)), // strip the dir from path as we don't want the path above dir
			Mode:     int64(stat.Mode()),
			Size:     stat.Size(),
		}

		if err := tgw.WriteHeader(header); err != nil {
			return err
		}

		if _, err := io.Copy(tgw, file); err != nil {
			return err
		}

	}
	return nil
}

// Close closes the underlying Gzip and tar archive types
func (tgw *TarGzipWriter) Close() error {
	if err := tgw.tw.Close(); err != nil {
		return err
	}

	if err := tgw.gw.Close(); err != nil {
		return err
	}

	return nil
}

// Next advances to the next entry in the tar archive
func (tgr *TarGzipReader) Next() (*tar.Header, error) {
	return tgr.tr.Next()

}

// Read reads from the current entry in the tar archive
func (tgr *TarGzipReader) Read(b []byte) (n int, err error) {
	return tgr.tr.Read(b)

}

// ExtractArchive extracts the archive into the given destination directory
func (tgr *TarGzipReader) ExtractArchive(dst string) error {

	for {
		header, err := tgr.tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if header.Typeflag != tar.TypeReg {
			continue
		}

		df := path.Join(dst, header.Name)

		if err := os.MkdirAll(path.Dir(df), 0755); err != nil {
			return err
		}

		file, err := os.Create(df)
		if err != nil {
			return err
		}
		defer file.Close()

		if _, err := io.Copy(file, tgr.tr); err != nil {
			return err
		}

	}

	return nil
}
