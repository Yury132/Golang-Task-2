package postgres

import (
	"context"
	"time"

	"github.com/Yury132/Golang-Task-2/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type Storage interface {
	// Загрузка данных в БД об изначальных изображениях
	SaveFileMeta(ctx context.Context, metaInfo *models.ImageMeta) error
	// Загрузка данных в БД о миниатюрах
	SaveFileMiniMeta(ctx context.Context, metaInfo *models.ImageMeta) error
	// Получаем информацию о картинках
	GetData(ctx context.Context) ([]models.AllImages, error)
	// Получаем информацию о картинках по id
	GetDataId(ctx context.Context, id int) ([]models.AllImages, error)
}

type storage struct {
	conn *pgxpool.Pool
}

// Загрузка данных в БД об изначальных изображениях
func (s *storage) SaveFileMeta(ctx context.Context, metaInfo *models.ImageMeta) error {
	query := "INSERT INTO public.uploads_info (name, type, width, height) VALUES ($1, $2, $3, $4)"

	// 10 секунд на выполнение операции с этим контекстом
	ctxDb, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	_, err := s.conn.Exec(ctxDb, query, metaInfo.Name, metaInfo.Type, metaInfo.Width, metaInfo.Height)
	if err != nil {
		return errors.Wrap(err, "failed to write file meta to db")
	}

	return nil
}

// Загрузка данных в БД о миниатюрах
func (s *storage) SaveFileMiniMeta(ctx context.Context, metaInfo *models.ImageMeta) error {
	query := "INSERT INTO public.mini_info (name, type, width, height) VALUES ($1, $2, $3, $4)"

	ctxDb, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	_, err := s.conn.Exec(ctxDb, query, metaInfo.Name, metaInfo.Type, metaInfo.Width, metaInfo.Height)
	if err != nil {
		return errors.Wrap(err, "failed to write fileMini meta to db")
	}

	return nil
}

// Получаем информацию о картинках
func (s *storage) GetData(ctx context.Context) ([]models.AllImages, error) {
	//query := "SELECT id, name, type, height, width FROM public.mini_info"

	query := "SELECT ui.id, ui.name, ui.type, ui.width, ui.height, mi.name, mi.width, mi.height FROM public.mini_info mi INNER JOIN public.uploads_info ui ON mi.id = ui.id"

	rows, err := s.conn.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var images = make([]models.AllImages, 0)

	for rows.Next() {
		var image models.AllImages
		if err = rows.Scan(&image.ID, &image.Name, &image.Type, &image.Width, &image.Height, &image.NameMini, &image.WidthMini, &image.HeightMini); err != nil {
			return nil, err
		}
		images = append(images, image)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return images, nil
}

// Получаем информацию о картинках по id
func (s *storage) GetDataId(ctx context.Context, id int) ([]models.AllImages, error) {
	//query := "SELECT id, name, type, height, width FROM public.mini_info"

	query := "SELECT ui.id, ui.name, ui.type, ui.width, ui.height, mi.name, mi.width, mi.height FROM public.mini_info mi INNER JOIN public.uploads_info ui ON mi.id = ui.id WHERE ui.id = $1"

	rows, err := s.conn.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var images = make([]models.AllImages, 0)

	for rows.Next() {
		var image models.AllImages
		if err = rows.Scan(&image.ID, &image.Name, &image.Type, &image.Width, &image.Height, &image.NameMini, &image.WidthMini, &image.HeightMini); err != nil {
			return nil, err
		}
		images = append(images, image)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return images, nil
}

func New(conn *pgxpool.Pool) Storage {
	return &storage{
		conn: conn,
	}
}
