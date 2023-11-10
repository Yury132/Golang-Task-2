package pool

import (
	"sync"

	"github.com/rs/zerolog"
)

type Pool struct {
	workers []*Worker

	log        zerolog.Logger
	workersNum int
	wg         sync.WaitGroup
}

// Выполнение задачи
func (p *Pool) RunBackground(f func()) {
	// Проходимся по всем воркерам
	for i := 0; i < p.workersNum; i++ {
		// Создаем воркера
		worker := NewWorker(p.log)
		// Добавляем в массив
		p.workers = append(p.workers, worker)
		// Устанавливаем конкретную задачу
		worker.SetTask(f)
		// Запускаем воркера выполнять эту задачу
		worker.Start(&p.wg)
	}
}

// Остановка
func (p *Pool) Stop() {
	for _, worker := range p.workers {
		// Останавливаем каждого воркера
		worker.Stop()
	}
	// Ждем когда остановятся все воркеры
	p.wg.Wait()
}

func New(log zerolog.Logger, workersNum int) *Pool {
	return &Pool{
		log:        log,
		workersNum: workersNum,
	}
}
