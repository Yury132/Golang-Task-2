package media_service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"os"

	"github.com/Yury132/Golang-Task-2/internal/models"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/nfnt/resize"
	"github.com/rs/zerolog"
)

type MediaService interface {
	// Создание миниатюры
	CreateThumbnail(info *models.InfoForThumbnail) error
	// Получаем сообщение из Nats
	GetTaskForProcessing() (*models.InfoForThumbnail, error)
}

type Storage interface {
	// Загрузка данных в БД об изначальных изображениях
	SaveFileMeta(ctx context.Context, metaInfo *models.ImageMeta) error
	// Загрузка данных в БД о миниатюрах
	SaveFileMiniMeta(ctx context.Context, metaInfo *models.ImageMeta) error
}

type ObjectStorage interface {
	// Сохранение изображения в хранилище
	Save(data []byte, name string) error
}

type mediaService struct {
	log           zerolog.Logger
	storage       Storage
	objectStorage ObjectStorage
	jsConsumer    jetstream.Consumer
}

// Получаем сообщение из Nats
func (m *mediaService) GetTaskForProcessing() (*models.InfoForThumbnail, error) {
	// Получаем сообщение из Nats
	// Next - блокируется пока нет входящих сообщений или не прошло время таймаута
	msg, err := m.jsConsumer.Next()
	if err != nil {
		return nil, err
	}

	// Подтверждение сообщения
	if err = msg.Ack(); err != nil {
		return nil, err
	}

	// Получаем данные из сообщения
	if msg.Data() == nil {
		return nil, nil
	}

	// Декодируем
	var info = new(models.InfoForThumbnail)
	if err = json.Unmarshal(msg.Data(), info); err != nil {
		return nil, err
	}

	return info, nil
}

// Через resize
// Создание миниатюры
// Тут же сохраняем данные в БД
func (m *mediaService) CreateThumbnail(info *models.InfoForThumbnail) error {

	// Открываем ранее сохраненную картинку
	file, err := os.Open(info.Path)
	if err != nil {
		m.log.Error().Err(err).Msg("failed to open file...")
		return err
	}
	// Получаем image.Image
	imageData, _, err := image.Decode(file)
	if err != nil {
		m.log.Error().Err(err).Msg("failed to decode...")
		return err
	}
	// Закрываем файл
	err = file.Close()
	if err != nil {
		m.log.Error().Err(err).Msg("failed to close file...")
		return err
	}
	// Создаем миниатюру
	newImage := resize.Thumbnail(100, 100, imageData, resize.Lanczos3)

	// Преобразуем в байты
	buf := new(bytes.Buffer)
	err = png.Encode(buf, newImage)
	if err != nil {
		m.log.Error().Err(err).Msg("failed to encode...")
		return err
	}
	imgBytes := buf.Bytes()

	// Создаем уникальное имя
	pName := uuid.New().String()

	// Сохраняем миниатюру в память
	if err = m.objectStorage.Save(imgBytes, fmt.Sprintf("%s.png", pName)); err != nil {
		m.log.Error().Err(err).Msg("objectStorage.Save err")
		return err
	}

	// Подготавливаем данные
	dataMini := &models.ImageMeta{Name: fmt.Sprintf("%s.png", pName), Type: "png", Width: newImage.Bounds().Max.X, Height: newImage.Bounds().Max.Y}

	// Сохраняем данные о миниатюре в БД
	if err = m.storage.SaveFileMiniMeta(context.Background(), dataMini); err != nil {
		m.log.Error().Err(err).Msg("failed to save data about mini to DB")
		return err
	}

	return nil
}

// То же самое, но только через Vips
// func (m *mediaService) CreateThumbnail(info *models.InfoForThumbnail) error {
// 	img, err := vips.NewImageFromFile(info.Path)

// 	if err != nil {
// 		m.log.Error().Err(err).Msg("NewImageFromFile err")
// 		return err
// 	}

// 	err = img.Thumbnail(info.Size, info.Size, vips.InterestingNone)
// 	if err != nil {
// 		m.log.Error().Err(err).Msg("ResizeWithVScale err")
// 		return err
// 	}

// 	pName := uuid.New().String()

// 	imgBytes, _, err := img.ExportNative()
// 	if err = m.objectStorage.Save(imgBytes, fmt.Sprintf("%s.png", pName)); err != nil {
// 		m.log.Error().Err(err).Msg("objectStorage.Save err")
// 		return err
// 	}

// 	return nil
// }

func New(log zerolog.Logger, storage Storage, objectStorage ObjectStorage, jsConsumer jetstream.Consumer) MediaService {
	return &mediaService{
		log:           log,
		storage:       storage,
		objectStorage: objectStorage,
		jsConsumer:    jsConsumer,
	}
}
