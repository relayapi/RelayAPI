package services

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/fatih/color"
)

type Stats struct {
	TotalRequests      uint64
	SuccessfulRequests uint64
	FailedRequests     uint64
	BytesReceived      uint64
	BytesSent          uint64
	StartTime          time.Time
}

func NewStats() *Stats {
	return &Stats{
		StartTime: time.Now(),
	}
}

// GetUptime 返回服务器运行时间
func (s *Stats) GetUptime() time.Duration {
	return time.Since(s.StartTime)
}

func (s *Stats) IncrementTotal() {
	atomic.AddUint64(&s.TotalRequests, 1)
}

func (s *Stats) IncrementSuccess() {
	atomic.AddUint64(&s.SuccessfulRequests, 1)
}

func (s *Stats) IncrementFailed() {
	atomic.AddUint64(&s.FailedRequests, 1)
}

func (s *Stats) AddBytesReceived(n uint64) {
	atomic.AddUint64(&s.BytesReceived, n)
}

func (s *Stats) AddBytesSent(n uint64) {
	atomic.AddUint64(&s.BytesSent, n)
}

// 格式化字节大小
func formatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := uint64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// StartConsoleDisplay 开始在控制台显示实时统计信息
func (s *Stats) StartConsoleDisplay(stopChan chan struct{}) {
	// 创建颜色输出
	titleColor := color.New(color.FgHiCyan, color.Bold)
	labelColor := color.New(color.FgHiYellow)
	valueColor := color.New(color.FgHiGreen)
	errorColor := color.New(color.FgHiRed)
	successColor := color.New(color.FgHiGreen)

	// 创建进度条字符
	progressChars := []string{"⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷"}
	progressIdx := 0

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	// 清除控制台并隐藏光标
	fmt.Print("\033[2J\033[?25l")
	defer fmt.Print("\033[?25h") // 恢复光标

	for {
		select {
		case <-stopChan:
			return
		case <-ticker.C:
			// 获取当前统计数据
			uptime := s.GetUptime()
			totalReqs := atomic.LoadUint64(&s.TotalRequests)
			successReqs := atomic.LoadUint64(&s.SuccessfulRequests)
			failedReqs := atomic.LoadUint64(&s.FailedRequests)
			bytesRecv := atomic.LoadUint64(&s.BytesReceived)
			bytesSent := atomic.LoadUint64(&s.BytesSent)
			tps := float64(totalReqs) / uptime.Seconds()

			// 移动光标到顶部
			fmt.Print("\033[H")

			// 显示标题
			titleColor.Printf("\n  %s RelayAPI Server Statistics %s\n\n", progressChars[progressIdx], progressChars[progressIdx])
			progressIdx = (progressIdx + 1) % len(progressChars)

			// 显示运行时间
			labelColor.Print("  ⏱️  Uptime: ")
			valueColor.Printf("%s\n", uptime.Round(time.Second))

			// 显示请求统计
			labelColor.Print("  🔄 Total Requests: ")
			valueColor.Printf("%d\n", totalReqs)

			// 显示成功/失败请求
			labelColor.Print("  ✅ Successful: ")
			successColor.Printf("%d", successReqs)
			labelColor.Print("  ❌ Failed: ")
			errorColor.Printf("%d\n", failedReqs)

			// 显示 TPS
			labelColor.Print("  ⚡ TPS: ")
			valueColor.Printf("%.2f\n", tps)

			// 显示流量统计
			labelColor.Print("  📥 Bytes Received: ")
			valueColor.Printf("%s", formatBytes(bytesRecv))
			labelColor.Print("  📤 Bytes Sent: ")
			valueColor.Printf("%s\n", formatBytes(bytesSent))

			// 显示成功率
			successRate := float64(0)
			if totalReqs > 0 {
				successRate = float64(successReqs) / float64(totalReqs) * 100
			}
			labelColor.Print("  📊 Success Rate: ")
			if successRate >= 90 {
				successColor.Printf("%.2f%%\n", successRate)
			} else if successRate >= 70 {
				valueColor.Printf("%.2f%%\n", successRate)
			} else {
				errorColor.Printf("%.2f%%\n", successRate)
			}

			// 添加分隔线
			fmt.Println("\n  " + color.HiBlackString(string(repeat('─', 50))))
		}
	}
}

func repeat(char rune, count int) []rune {
	result := make([]rune, count)
	for i := range result {
		result[i] = char
	}
	return result
}
