package worker

import (
	"sync"

	"github.com/Yury132/Golang-Task-2/internal/models"
	"github.com/rs/zerolog"

	"github.com/Yury132/Golang-Task-2/internal/worker/pool"
)

// Эти функции будет вызывать воркер пул
// Связь с "media_service"
type MediaService interface {
	// Создание миниатюры
	CreateThumbnail(info *models.InfoForThumbnail) error
	// Получаем сообщение из Nats
	GetTaskForProcessing() (*models.InfoForThumbnail, error)
}

type MediaHandler struct {
	log          zerolog.Logger
	mediaService MediaService
	pool         *pool.Pool
}

// Запуск
func (mh *MediaHandler) Start() {
	// Передаем функцию, которую будет выполнять воркер пул
	mh.pool.RunBackground(mh.createThumbnail)
}

// Остановка
func (mh *MediaHandler) Shutdown() {
	mh.pool.Stop()
}

// Функция, которую будет выполнять воркер пул
func (mh *MediaHandler) createThumbnail() {
	info, err := mh.mediaService.GetTaskForProcessing()
	if err != nil {
		mh.log.Error().Err(err).Send()
		return
	}

	var wg = new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err = mh.mediaService.CreateThumbnail(info); err != nil {
			mh.log.Error().Err(err).Send()
			return
		}
	}()
	wg.Wait()
}

func New(log zerolog.Logger, mediaService MediaService, workersNum int) *MediaHandler {
	return &MediaHandler{
		log:          log,
		mediaService: mediaService,
		pool:         pool.New(log, workersNum),
	}
}
