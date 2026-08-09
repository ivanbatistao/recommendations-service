package workerpool

import (
	"context"
	"log/slog"
	"sync"

	"github.com/ivanbatistao/recommendations-service/internal/domain/event"
)

type EventProcessor interface {
	ProcessEvent(ctx context.Context, ev event.Event) error
}

type WorkerPool struct {
	numWorkers   int
	eventChan   chan event.Event
	processor   EventProcessor
	wg          sync.WaitGroup
	logger      *slog.Logger
	ctx         context.Context
	cancel      context.CancelFunc
}

func NewWorkerPool(
	numWorkers int,
	bufferSize int,
	processor EventProcessor,
	logger *slog.Logger,
) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())

	return &WorkerPool{
		numWorkers: numWorkers,
		eventChan:  make(chan event.Event, bufferSize),
		processor:  processor,
		logger:     logger,
		ctx:        ctx,
		cancel:     cancel,
	}
}

func (wp *WorkerPool) Start() {
	wp.logger.Info(
		"starting worker pool",
		slog.Int("workers", wp.numWorkers),
		slog.Int("buffer_size", cap(wp.eventChan)),
	)

	for i := 0; i < wp.numWorkers; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}
}

func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()

	wp.logger.Info(
		"worker started",
		slog.Int("worker_id", id),
	)

	for {
		select {
		case <-wp.ctx.Done():
			wp.logger.Info(
				"worker stopping",
				slog.Int("worker_id", id),
			)
			return

		case ev := <-wp.eventChan:
			wp.logger.Debug(
				"worker processing event",
				slog.Int("worker_id", id),
				slog.String("event_id", ev.EventID),
				slog.String("user_id", ev.UserID),
			)

			if err := wp.processor.ProcessEvent(wp.ctx, ev); err != nil {
				wp.logger.Error(
					"worker failed to process event",
					slog.Int("worker_id", id),
					slog.String("event_id", ev.EventID),
					slog.String("error", err.Error()),
				)
			} else {
				wp.logger.Debug(
					"worker processed event successfully",
					slog.Int("worker_id", id),
					slog.String("event_id", ev.EventID),
				)
			}
		}
	}
}

func (wp *WorkerPool) Submit(ev event.Event) {
	select {
	case wp.eventChan <- ev:
		// Event submitted successfully
	default:
		wp.logger.Warn(
			"event channel full, dropping event",
			slog.String("event_id", ev.EventID),
		)
	}
}

func (wp *WorkerPool) Stop() {
	wp.logger.Info("stopping worker pool")

	wp.cancel()
	close(wp.eventChan)
	wp.wg.Wait()

	wp.logger.Info("worker pool stopped")
}

func (wp *WorkerPool) Stats() map[string]interface{} {
	return map[string]interface{}{
		"num_workers":  wp.numWorkers,
		"buffer_size":  cap(wp.eventChan),
		"buffer_used":  len(wp.eventChan),
	}
}
