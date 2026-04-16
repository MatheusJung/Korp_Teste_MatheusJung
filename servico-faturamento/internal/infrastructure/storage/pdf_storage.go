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

func (s *LocalPDFStorage) Salvar(ctx context.Context, notaID uuid.UUID, numero int64, conteudo []byte) error {
	_ = ctx
	_ = notaID

	if err := os.MkdirAll(s.basePath, 0o755); err != nil {
		return err
	}

	nomeArquivo := fmt.Sprintf("NF-%06d.pdf", numero)
	caminho := filepath.Join(s.basePath, nomeArquivo)

	if err := os.WriteFile(caminho, conteudo, 0o644); err != nil {
		return err
	}

	return nil
}

func (s *LocalPDFStorage) Obter(ctx context.Context, notaID uuid.UUID, numero int64) ([]byte, string, error) {
	_ = ctx
	_ = notaID

	nomeArquivo := fmt.Sprintf("NF-%06d.pdf", numero)
	path := filepath.Join(s.basePath, nomeArquivo)

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, "", fmt.Errorf("pdf da nota não encontrado")
		}
		return nil, "", err
	}

	return data, nomeArquivo, nil
}