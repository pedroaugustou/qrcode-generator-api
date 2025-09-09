package worker

import (
	"context"
	"log"
	"time"

	"github.com/pedroaugustou/qrcode-generator-api/internal/domain/service"
	"github.com/pedroaugustou/qrcode-generator-api/internal/infrastructure/repository"
)

type CleanupWorker struct {
	storageService   service.StorageService
	qrCodeRepository repository.QRCodeRepository
	interval         time.Duration
	stopCh           chan struct{}
	doneCh           chan struct{}
}

func NewCleanupWorker(s service.StorageService,
	q repository.QRCodeRepository,
	i time.Duration) *CleanupWorker {
	return &CleanupWorker{
		storageService:   s,
		qrCodeRepository: q,
		interval:         i,
		stopCh:           make(chan struct{}),
		doneCh:           make(chan struct{}),
	}
}

func (w *CleanupWorker) Start(ctx context.Context) {
	go w.run(ctx)
}

func (w *CleanupWorker) Stop() {
	close(w.stopCh)
	<-w.doneCh
}

func (w *CleanupWorker) run(ctx context.Context) {
	defer close(w.doneCh)

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			w.runCleanup(ctx)
		case <-w.stopCh:
			log.Println("cleanup worker stopped")
			return
		case <-ctx.Done():
			log.Println("cleanup worker stopped due to context cancellation")
			return
		}
	}
}

func (w *CleanupWorker) runCleanup(parentCtx context.Context) {
	ctx, cancel := context.WithTimeout(parentCtx, 30*time.Second)
	defer cancel()

	log.Println("starting cleanup process...")

	if err := w.qrCodeRepository.DeleteExpiredQRCodes(ctx); err != nil {
		log.Printf("failed to delete expired records: %v", err)
		return
	}

	if err := w.storageService.CleanupExpiredFiles(ctx); err != nil {
		log.Printf("failed to cleanup expired files: %v", err)
		return
	}

	log.Println("cleanup process completed successfully")
}
