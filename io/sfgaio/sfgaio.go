package sfgaio

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/sfborg/sflib/ent/sfga"
)

type sfgaio struct {
	repo sfga.GitRepo
	path string
}

func New(repo sfga.GitRepo, path string) (sfga.SFGA, error) {
	var err error
	res := &sfgaio{repo: repo, path: path}
	err = res.cleanPath()
	return res, err
}

func (s *sfgaio) FetchSchema() ([]byte, error) {
	err := s.cloneRepo()
	if err != nil {
		return nil, err
	}

	schemaPath := filepath.Join(s.path, "schema.sql")
	res, err := os.ReadFile(schemaPath)
	if err != nil {
		err = fmt.Errorf("cannot read %s: %w", schemaPath, err)
		return nil, err
	}

	err = s.cleanPath()
	return res, err
}

func (s *sfgaio) cleanPath() error {
	err := os.RemoveAll(s.path)
	if err != nil {
		err = fmt.Errorf("cannot remove path %s: %w", s.path, err)
		return err
	}
	return nil
}

func (s *sfgaio) cloneRepo() error {
	var err error
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
