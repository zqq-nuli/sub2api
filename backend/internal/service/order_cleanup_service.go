package service

import (
	"context"
	"log"
	"sync"
	"time"
)

// OrderCleanupService 订单清理服务（定时将过期订单标记为expired）
type OrderCleanupService struct {
	orderRepo OrderRepository
	interval  time.Duration
	stopCh    chan struct{}
	wg        sync.WaitGroup
}

// NewOrderCleanupService 创建订单清理服务
func NewOrderCleanupService(orderRepo OrderRepository) *OrderCleanupService {
	return &OrderCleanupService{
		orderRepo: orderRepo,
		interval:  1 * time.Minute, // 每分钟检查一次
		stopCh:    make(chan struct{}),
	}
}

// Start 启动订单清理服务
func (s *OrderCleanupService) Start() {
	s.wg.Add(1)
	go s.run()
	log.Printf("[OrderCleanup] Service started, interval: %v", s.interval)
}

// Stop 停止订单清理服务
func (s *OrderCleanupService) Stop() {
	close(s.stopCh)
	s.wg.Wait()
	log.Printf("[OrderCleanup] Service stopped")
}

func (s *OrderCleanupService) run() {
	defer s.wg.Done()

	// 启动时立即执行一次
	s.cleanup()

	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.cleanup()
		case <-s.stopCh:
			return
		}
	}
}

func (s *OrderCleanupService) cleanup() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	affected, err := s.orderRepo.MarkExpiredOrders(ctx)
	if err != nil {
		log.Printf("[OrderCleanup] Failed to mark expired orders: %v", err)
		return
	}

	if affected > 0 {
		log.Printf("[OrderCleanup] Marked %d orders as expired", affected)
	}
}
