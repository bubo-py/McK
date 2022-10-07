package postgres

import (
	"embed"
	"io/fs"
	"os"
	"path"

	"github.com/pkg/errors"
)

type adapterFS struct {
	embed.FS
	rootDir string
}

// ReadDir for adopt embed.FS API to MigratorFS
func (efs adapterFS) ReadDir(dirname string) ([]os.FileInfo, error) {
	dirEntries, err := efs.FS.ReadDir(dirname)
	if err != nil {
		return nil, err
	}

	fileInfos := make([]fs.FileInfo, 0, len(dirEntries))
	for _, e := range dirEntries {
		fi, err := e.Info()
		if err != nil {
			continue // file is missing, skip it
		}

		fileInfos = append(fileInfos, fi)
	}

	return fileInfos, nil
}

// Glob for adopt embed.FS API to MigratorFS, no real implementation
func (efs adapterFS) Glob(pattern string) (matches []string, err error) {
	des, err := efs.FS.ReadDir(efs.rootDir)
	if err != nil {
		return nil, errors.Wrap(err, "try to read from pattern as path")
	}

	files := make([]string, 0, len(des))

	pattern = "migrations/*/*.sql"
	for _, e := range des {
		matches, err := path.Match(pattern, e.Name())
		// Pattern is malformed.
		if err != nil {
			return nil, err
		}

		if !matches {
			continue
		}

		files = append(files, path.Join(efs.rootDir, e.Name()))
	}

	return files, nil
}
