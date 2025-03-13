package isfga

import (
	"crypto/sha256"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/sfborg/sflib/pkg/arch"
	"github.com/sfborg/sflib/pkg/sfga"
	_ "modernc.org/sqlite"
)

type schema struct {
	repo sfga.GitRepo
}

func NewSchema(repo sfga.GitRepo) arch.SchemaSFGA {
	res := schema{repo: repo}
	return &res
}

func (s *schema) Fetch() ([]byte, error) {
	tempDir, err := os.MkdirTemp("", "git-schema-")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tempDir)

	err = s.cloneRepo(tempDir)
	if err != nil {
		return nil, err
	}
	schemaPath := filepath.Join(tempDir, "schema.sql")

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
	return res, err
}

func (s *schema) cloneRepo(tmpDir string) error {
	cmd := exec.Command(
		"git", "clone", "--depth=1", "--branch", s.repo.Tag, s.repo.URL, tmpDir,
	)
	err := cmd.Run()
	if err != nil {
		err = fmt.Errorf("cannot clone GitHub Repo %s: %w", s.repo.URL, err)
		return &arch.ErrRepoClone{URL: s.repo.URL, Err: err}
	}

	return nil
}
