package object_storage

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type ObjectStorage interface {
	// Сохранение изображения в хранилище
	Save(data []byte, name string) error
}

type objectStorage struct {
	log zerolog.Logger
}

// Сохранение изображения в хранилище
func (o *objectStorage) Save(data []byte, name string) error {
	path := fmt.Sprintf("uploads/%s", name)

	f, err := os.Create(path)
	if err != nil {
		return errors.Wrap(err, "failed to create file")
	}
	defer func() {
		if err = f.Close(); err != nil {
			o.log.Error().Err(err).Send()
		}
	}()

	_, err = f.Write(data)
	if err != nil {
		return errors.Wrap(err, "failed to write data to file")
	}

	return nil
}

func New(log zerolog.Logger) ObjectStorage {
	return &objectStorage{
		log: log,
	}
}
