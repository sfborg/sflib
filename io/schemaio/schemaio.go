package schemaio

import (
	"crypto/sha256"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gnames/gnsys"
	"github.com/sfborg/sflib/ent/sfga"
)

type schemaio struct {
	repo sfga.GitRepo
	path string
}

func New(repo sfga.GitRepo, path string) sfga.Schema {
	res := &schemaio{repo: repo, path: path}
	return res
}

func (s *schemaio) Fetch() ([]byte, error) {
	res, err := s.loadSchema()
	if err == nil {
		return res, nil
	}

	err = s.cloneRepo()
	if err != nil {
		return nil, err
	}

	res, err = s.loadSchema()
	return res, err
}

// Clean removes SFGA data directory.
func (s *schemaio) Clean() error {
	err := os.RemoveAll(s.path)
	if err != nil {
		return &sfga.ErrDirRemove{Dir: s.path, Err: err}
	}
	return nil
}

// GitRepo returns GitRepo of its instance.
func (s *schemaio) GitRepo() sfga.GitRepo {
	return s.repo
}

// Path returns temporary path where SFGA schema is downloaded.
func (s *schemaio) Path() string {
	return s.path
}

func (s *schemaio) loadSchema() ([]byte, error) {
	var err error
	var exists bool
	schemaPath := filepath.Join(s.path, "schema.sql")
	exists, err = gnsys.FileExists(schemaPath)

	if err != nil {
		err = fmt.Errorf("bad file %s: %w", schemaPath, err)
		return nil, err
	}

	if !exists {
		err = fmt.Errorf("file %s does not exist", schemaPath)
		return nil, err
	}

	res, err := os.ReadFile(schemaPath)
	if err != nil {
		err = fmt.Errorf("cannot read %s: %w", schemaPath, err)
		return nil, err
	}

	sum := sha256.Sum256(res)
	hash := fmt.Sprintf("%x", sum)
	if !strings.HasPrefix(hash, s.repo.ShaSchemaSQL) {
		err = fmt.Errorf("Schema does not match %s", s.repo.ShaSchemaSQL)
		return nil, err
	}

	return res, nil
}

func (s *schemaio) cloneRepo() error {
	var err error
	err = s.Clean()
	if err != nil {
		return &sfga.ErrRepoCacheClean{Dir: s.path, Err: err}
	}

	var currentDir string
	cmd := exec.Command("git", "clone", s.repo.URL, s.path)
	err = cmd.Run()
	if err != nil {
		err = fmt.Errorf("cannot clone GitHub Repo %s: %w", s.repo.URL, err)
		return &sfga.ErrRepoClean{URL: s.repo.URL, Err: err}
	}

	currentDir, err = os.Getwd()
	if err != nil {
		return err
	}
	defer os.Chdir(currentDir)

	err = os.Chdir(s.path)
	if err != nil {
		return &sfga.ErrDirChange{Src: currentDir, Dst: s.path, Err: err}
	}

	if s.repo.Tag == "" {
		return nil
	}

	cmd = exec.Command("git", "checkout", s.repo.Tag)
	err = cmd.Run()
	if err != nil {
		return &sfga.ErrRepoTagCheckout{Tag: s.repo.Tag, Err: err}
	}

	return nil
}
