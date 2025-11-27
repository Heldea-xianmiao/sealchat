package service

import (
	"math"
	"sync"
	"time"
)

// ExportLimiterConfig 定义导出下载限速参数（单位：KB）。
type ExportLimiterConfig struct {
	BandwidthKBps int
	BurstKB       int
}

var (
	exportLimiterMu sync.RWMutex
	exportLimiter   *downloadLimiter
)

// InitExportLimiter 初始化全局导出带宽限制器。
func InitExportLimiter(cfg ExportLimiterConfig) {
	exportLimiterMu.Lock()
	defer exportLimiterMu.Unlock()

	rateBytes := int64(cfg.BandwidthKBps) * 1024
	burstBytes := int64(cfg.BurstKB) * 1024
	exportLimiter = newDownloadLimiter(rateBytes, burstBytes)
}

// WaitExportBandwidth 会在需要时阻塞，直到本次传输的字节数满足速率要求。
func WaitExportBandwidth(bytes int) {
	if bytes <= 0 {
		return
	}
	exportLimiterMu.RLock()
	limiter := exportLimiter
	exportLimiterMu.RUnlock()
	if limiter == nil {
		return
	}
	limiter.waitBytes(bytes)
}

type downloadLimiter struct {
	mu             sync.Mutex
	bytesPerSecond float64
	capacity       float64
	tokens         float64
	last           time.Time
}

func newDownloadLimiter(rateBytes, burstBytes int64) *downloadLimiter {
	if rateBytes <= 0 {
		return nil
	}
	if burstBytes <= 0 {
		burstBytes = rateBytes
	}
	return &downloadLimiter{
		bytesPerSecond: float64(rateBytes),
		capacity:       float64(burstBytes),
		tokens:         float64(burstBytes),
		last:           time.Now(),
	}
}

func (l *downloadLimiter) waitBytes(n int) {
	if l == nil || n <= 0 {
		return
	}
	required := float64(n)
	for {
		l.mu.Lock()
		now := time.Now()
		if l.last.IsZero() {
			l.last = now
		}
		elapsed := now.Sub(l.last).Seconds()
		if elapsed > 0 {
			l.tokens = math.Min(l.capacity, l.tokens+elapsed*l.bytesPerSecond)
			l.last = now
		}
		if l.tokens >= required {
			l.tokens -= required
			l.mu.Unlock()
			return
		}

		deficit := required - l.tokens
		waitSeconds := deficit / l.bytesPerSecond
		if waitSeconds < 0.001 {
			waitSeconds = 0.001
		}
		l.mu.Unlock()
		time.Sleep(time.Duration(waitSeconds * float64(time.Second)))
	}
}
