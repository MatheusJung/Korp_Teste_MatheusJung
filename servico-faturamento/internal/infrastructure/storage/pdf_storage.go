package storage

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

type LocalPDFStorage struct {
	basePath string
}

func NewLocalPDFStorage(basePath string) *LocalPDFStorage {
	return &LocalPDFStorage{basePath: basePath}
}

func (s *LocalPDFStorage) Salvar(ctx context.Context, notaID uuid.UUID, numero int64, data []byte) error {
	if err := os.MkdirAll(s.basePath, 0o755); err != nil {
		return err
	}

	path := filepath.Join(s.basePath, notaID.String()+".pdf")
	return os.WriteFile(path, data, 0o644)
}

func (s *LocalPDFStorage) Obter(ctx context.Context, notaID uuid.UUID) ([]byte, string, error) {
	path := filepath.Join(s.basePath, notaID.String()+".pdf")

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, "", fmt.Errorf("pdf da nota não encontrado")
		}
		return nil, "", err
	}

	return data, notaID.String() + ".pdf", nil
}