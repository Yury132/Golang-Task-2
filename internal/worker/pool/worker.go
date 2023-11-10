package pool

import (
	"sync"
	"time"

	"github.com/rs/zerolog"
)

type Worker struct {
	log zerolog.Logger
	// Канал с пустой структурой (структура - потому что она самая легкая)
	quit chan struct{}
	// Выполняемая задача
	task func()
}

// Запускаем воркера выполнять задачу
func (w *Worker) Start(wg *sync.WaitGroup) {
	go func() {
		w.log.Info().Msg("Start")
		ticker := time.NewTicker(200 * time.Millisecond)
		defer ticker.Stop()

		wg.Add(1)
		defer wg.Done()

		// Бесконечный цикл
		for {
			select {
			// Каждый цикл тикер отправляет сигнал в канал
			case <-ticker.C:
				// И воркер может начать выполнять задачу
				w.task()

			// Если что-то получили в канал, выходим из бесконечного цикла
			case <-w.quit:
				w.log.Info().Msg("Stopped")
				return
			}
		}
	}()
}

// Устанавливаем конкретную задачу
func (w *Worker) SetTask(task func()) {
	w.task = task
}

// Останавливаем воркера
func (w *Worker) Stop() {
	// Записываем в канал пустую структуру
	w.quit <- struct{}{}
}

func NewWorker(log zerolog.Logger) *Worker {
	return &Worker{
		log:  log,
		quit: make(chan struct{}, 1),
	}
}
