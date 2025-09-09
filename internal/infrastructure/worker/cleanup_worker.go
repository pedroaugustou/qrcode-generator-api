package worker

import (
	"context"
	"log"
	"time"

	"github.com/pedroaugustou/qrcode-generator-api/internal/domain/entity"
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

	expiredQRCodes, err := w.getExpiredQRCodes(ctx)
	if err != nil {
		log.Printf("failed to get expired QR codes: %v", err)
		return
	}

	if len(expiredQRCodes) == 0 {
		log.Println("no expired QR codes found")
		return
	}

	log.Printf("found %d expired QR codes to cleanup", len(expiredQRCodes))

	for _, qrCode := range expiredQRCodes {
		if err := w.deleteQRCode(ctx, qrCode); err != nil {
			log.Printf("failed to delete QR code %s: %v", qrCode.ID, err)
			continue
		}
		log.Printf("successfully deleted QR code: %s", qrCode.ID)
	}

	log.Println("cleanup process completed successfully")
}

func (w *CleanupWorker) getExpiredQRCodes(ctx context.Context) ([]entity.QRCode, error) {
	allQRCodes, err := w.qrCodeRepository.GetAllQRCodes(ctx)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC().Truncate(time.Hour)
	var expiredQRCodes []entity.QRCode

	for _, qrCode := range allQRCodes {
		if qrCode.ExpiresAt.Before(now) || qrCode.ExpiresAt.Equal(now) {
			expiredQRCodes = append(expiredQRCodes, qrCode)
		}
	}

	return expiredQRCodes, nil
}

func (w *CleanupWorker) deleteQRCode(ctx context.Context, qrCode entity.QRCode) error {
	if err := w.storageService.DeleteQRCode(ctx, qrCode.URL); err != nil {
		return err
	}

	if err := w.qrCodeRepository.DeleteQRCode(ctx, qrCode.ID); err != nil {
		return err
	}

	return nil
}
