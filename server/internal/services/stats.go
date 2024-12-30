package services

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"golang.org/x/term"

	"relayapi/server/internal/middleware/logger"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type Stats struct {
	TotalRequests      uint64
	SuccessfulRequests uint64
	FailedRequests     uint64
	BytesReceived      uint64
	BytesSent          uint64
	StartTime          time.Time
	errorStats         sync.Map // 用于存储每个错误状态码的计数
	Version            string   // 版本号
	ServerAddr         string   // 服务器地址
}

func NewStats(version, serverAddr string) *Stats {
	return &Stats{
		StartTime:  time.Now(),
		Version:    version,
		ServerAddr: serverAddr,
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

// IncrementErrorStatus 增加特定错误状态码的计数
func (s *Stats) IncrementErrorStatus(statusCode int) {
	if value, ok := s.errorStats.Load(statusCode); ok {
		atomic.AddUint64(value.(*uint64), 1)
	} else {
		var counter uint64 = 1
		s.errorStats.Store(statusCode, &counter)
	}
}

// GetErrorStats 获取错误状态码统计
func (s *Stats) GetErrorStats() map[int]uint64 {
	stats := make(map[int]uint64)
	s.errorStats.Range(func(key, value interface{}) bool {
		stats[key.(int)] = atomic.LoadUint64(value.(*uint64))
		return true
	})
	return stats
}

func (s *Stats) AddBytesReceived(n uint64) {
	atomic.AddUint64(&s.BytesReceived, n)
}

func (s *Stats) AddBytesSent(n uint64) {
	atomic.AddUint64(&s.BytesSent, n)
}

// 获取错误状态码的描述
func getStatusCodeDesc(code int) string {
	switch code {
	case 400:
		return "Bad Request"
	case 401:
		return "Unauthorized"
	case 403:
		return "Forbidden"
	case 404:
		return "Not Found"
	case 405:
		return "Method Not Allowed"
	case 408:
		return "Request Timeout"
	case 429:
		return "Too Many Requests"
	case 500:
		return "Internal Server Error"
	case 502:
		return "Bad Gateway"
	case 503:
		return "Service Unavailable"
	case 504:
		return "Gateway Timeout"
	default:
		return "Unknown Error"
	}
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
	var uiActive bool = true
	var uiQuit bool = false

	// 保存原始终端设置
	oldState, err := term.MakeRaw(int(syscall.Stdin))
	if err != nil {
		log.Printf("无法设置终端为原始模式: %v", err)
		return
	}
	defer term.Restore(int(syscall.Stdin), oldState)

	// 创建一个函数来启动 UI
	startUI := func() error {
		if err := ui.Init(); err != nil {
			log.Printf("failed to initialize termui: %v", err)
			return err
		}
		uiActive = true
		return nil
	}

	// 初始启动 UI
	if err := startUI(); err != nil {
		return
	}
	defer ui.Close()

	// 创建标题
	title := widgets.NewParagraph()
	title.Title = "RelayAPI Server"
	title.Text = fmt.Sprintf("Version: %s   |   Server: %s", s.Version, s.ServerAddr)
	title.TextStyle.Fg = ui.ColorCyan
	title.BorderStyle.Fg = ui.ColorCyan
	title.TitleStyle.Fg = ui.ColorCyan

	// 创建基本统计信息区域
	basicStats := widgets.NewParagraph()
	basicStats.Title = "Basic Statistics"
	basicStats.BorderStyle.Fg = ui.ColorYellow

	// 创建请求统计图表
	requestsPlot := widgets.NewPlot()
	requestsPlot.Title = "Requests Per Second"
	requestsPlot.Data = make([][]float64, 1)
	requestsPlot.Data[0] = []float64{0, 0} // 初始化为两个零点
	requestsPlot.LineColors = []ui.Color{ui.ColorYellow}
	requestsPlot.BorderStyle.Fg = ui.ColorYellow
	requestsPlot.AxesColor = ui.ColorWhite
	requestsPlot.DrawDirection = widgets.DrawLeft
	requestsPlot.MaxVal = 100

	// 创建错误统计区域
	errorStats := widgets.NewParagraph()
	errorStats.Title = "Error Statistics"
	errorStats.BorderStyle.Fg = ui.ColorRed

	// 创建日志区域
	logView := widgets.NewParagraph()
	logView.Title = "Recent Logs"
	logView.BorderStyle.Fg = ui.ColorBlue

	// 初始化计数器和数据切片
	lastTotal := atomic.LoadUint64(&s.TotalRequests)
	tpsData := []float64{0, 0} // 初始化为两个零点

	// 创建事件处理通道
	uiEvents := ui.PollEvents()
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	// 设置布局函数
	updateUI := func() {
		if !uiActive {
			return
		}
		// 获取终端大小
		width, height := ui.TerminalDimensions()

		// 设置各个组件的位置和大小
		title.SetRect(0, 0, width, 3)
		basicStats.SetRect(0, 3, width/2, height/2)
		requestsPlot.SetRect(width/2, 3, width, height/2)
		errorStats.SetRect(0, height/2, width/2, height-3)
		logView.SetRect(width/2, height/2, width, height-3)

		// 更新统计数据
		uptime := s.GetUptime()
		totalReqs := atomic.LoadUint64(&s.TotalRequests)
		successReqs := atomic.LoadUint64(&s.SuccessfulRequests)
		failedReqs := atomic.LoadUint64(&s.FailedRequests)
		bytesRecv := atomic.LoadUint64(&s.BytesReceived)
		bytesSent := atomic.LoadUint64(&s.BytesSent)

		// 计算 TPS
		currentTPS := float64(totalReqs-lastTotal) / 1.0
		lastTotal = totalReqs

		// 更新图表数据，确保至少有两个点
		if len(tpsData) < 2 {
			tpsData = []float64{0, currentTPS}
		} else {
			tpsData = append(tpsData, currentTPS)
			if len(tpsData) > 60 {
				tpsData = tpsData[1:]
			}
		}

		// 动态调整最大值
		maxTPS := currentTPS
		for _, v := range tpsData {
			if v > maxTPS {
				maxTPS = v
			}
		}
		requestsPlot.MaxVal = maxTPS * 1.2 // 设置为最大值的 1.2 倍，留出一些空间
		if requestsPlot.MaxVal < 10 {      // 设置最小值，避免图表太扁
			requestsPlot.MaxVal = 10
		}

		requestsPlot.Data[0] = tpsData
		requestsPlot.Title = fmt.Sprintf("Requests Per Second (Current: %.2f)", currentTPS)

		// 计算成功率
		successRate := float64(0)
		if totalReqs > 0 {
			successRate = float64(successReqs) / float64(totalReqs) * 100
		}

		// 更新基本统计信息
		basicStats.Text = fmt.Sprintf(
			"⏱️  Uptime: %s\n"+
				"🔄 Total Requests: %d\n"+
				"✅ Successful: %d\n"+
				"❌ Failed: %d\n"+
				"📥 Bytes Received: %s\n"+
				"📤 Bytes Sent: %s\n"+
				"📊 Success Rate: %.2f%%",
			uptime.Round(time.Second),
			totalReqs,
			successReqs,
			failedReqs,
			formatBytes(bytesRecv),
			formatBytes(bytesSent),
			successRate,
		)

		// 更新错误统计信息
		if failedReqs > 0 {
			var errorText strings.Builder
			errStats := s.GetErrorStats()
			var codes []int
			for code := range errStats {
				codes = append(codes, code)
			}
			sort.Ints(codes)

			for _, code := range codes {
				count := errStats[code]
				percentage := float64(count) / float64(failedReqs) * 100
				errorText.WriteString(fmt.Sprintf("%d %s\nCount: %d (%.2f%%)\n\n",
					code,
					getStatusCodeDesc(code),
					count,
					percentage,
				))
			}
			errorStats.Text = errorText.String()
		} else {
			errorStats.Text = "No errors reported"
		}

		// 更新日志视图
		logView.Text = logger.GetRecentLogs()

		// 渲染所有组件
		ui.Render(title, basicStats, requestsPlot, errorStats, logView)
	}

	// 主事件循环
	logUpdateChan := logger.GetLogUpdateChan()

	// 创建按键读取缓冲区
	buf := make([]byte, 1)

	for !uiQuit {
		if uiActive {
			select {
			case e := <-uiEvents:
				switch e.ID {
				case "q", "<C-q>":
					// 切换到普通模式
					ui.Close()
					uiActive = false
					fmt.Println("\n“Press ‘q’ or ‘Ctrl+q’ to switch to UI mode, or ‘Ctrl+C’ to exit.”")
				case "<C-c>":
					uiQuit = true
				case "<Resize>":
					updateUI()
				}
			case <-ticker.C:
				updateUI()
			case <-logUpdateChan:
				updateUI()
			case <-stopChan:
				return
			}
		} else {
			// 普通模式下的事件处理
			select {
			case <-ticker.C:
				// 在普通模式下打印基本统计信息
				uptime := s.GetUptime()
				totalReqs := atomic.LoadUint64(&s.TotalRequests)
				successReqs := atomic.LoadUint64(&s.SuccessfulRequests)
				failedReqs := atomic.LoadUint64(&s.FailedRequests)
				fmt.Printf("\rUptime: %s | Requests: %d | Success: %d | Failed: %d | TPS: %.2f",
					uptime.Round(time.Second),
					totalReqs,
					successReqs,
					failedReqs,
					float64(totalReqs)/uptime.Seconds())
			case <-stopChan:
				return
			default:
				// 检查键盘输入
				if n, err := os.Stdin.Read(buf); err == nil && n == 1 {
					switch buf[0] {
					case 'q':
						// 切换回 UI 模式
						fmt.Print("\n") // 在切换回 UI 模式前换行，保持输出整洁
						if err := startUI(); err == nil {
							updateUI()
						}
					case 3: // Ctrl+C
						uiQuit = true
					}
				}
			}
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
