package sfgaio

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

type sfgaio struct {
	repo sfga.GitRepo
	path string
}

func New(repo sfga.GitRepo, path string) sfga.SFGA {
	res := &sfgaio{repo: repo, path: path}
	return res
}

func (s *sfgaio) FetchSchema() ([]byte, error) {
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
func (s *sfgaio) Clean() error {
	err := os.RemoveAll(s.path)
	if err != nil {
		err = fmt.Errorf("cannot remove path %s: %w", s.path, err)
		return err
	}
	return nil
}

// GitRepo returns GitRepo of its instance.
func (s *sfgaio) GitRepo() sfga.GitRepo {
	return s.repo
}

// Path returns temporary path where SFGA schema is downloaded.
func (s *sfgaio) Path() string {
	return s.path
}

func (s *sfgaio) loadSchema() ([]byte, error) {
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

func (s *sfgaio) cloneRepo() error {
	var err error
	err = s.Clean()
	if err != nil {
		return err
	}

	var currentDir string
	cmd := exec.Command("git", "clone", s.repo.URL, s.path)
	err = cmd.Run()
	if err != nil {
		err = fmt.Errorf("cannot clone GitHub Repo %s: %w", s.repo.URL, err)
		return err
	}

	currentDir, err = os.Getwd()
	if err != nil {
		return err
	}
	defer os.Chdir(currentDir)

	err = os.Chdir(s.path)
	if err != nil {
		err = fmt.Errorf("cannot chdir %s: %w", s.path, err)
		return err
	}

	if s.repo.Tag == "" {
		return nil
	}

	cmd = exec.Command("git", "checkout", s.repo.Tag)
	err = cmd.Run()
	if err != nil {
		err = fmt.Errorf("cannot checkout tag %s: %w", s.repo.Tag, err)
		return err
	}

	return nil
}
