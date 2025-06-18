package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	// 参数解析
	startTimeStr := flag.String("start", "", "开始时间 (格式: 2006-01-02T15:04:05)")
	commandFile := flag.String("file", "", "命令文件路径")
	totalOps := flag.Int("ops", 1000, "总操作数")
	maxConc := flag.Int("concurrency", 100, "最大并发数")
	flag.Parse()

	if *startTimeStr == "" || *commandFile == "" || *totalOps <= 0 {
		fmt.Println("必须提供 -start, -file, -ops 参数，且 ops > 0")
		flag.Usage()
		return
	}

	// 时间解析
	startTime, err := time.Parse("2006-01-02T15:04:05", *startTimeStr)
	if err != nil {
		fmt.Printf("开始时间格式错误: %v\n", err)
		return
	}

	// 读取命令文件
	lines, err := readLines(*commandFile)
	if err != nil {
		fmt.Printf("读取命令文件失败: %v\n", err)
		return
	}
	if len(lines) == 0 {
		fmt.Println("命令文件为空")
		return
	}

	fmt.Printf("等待开始时间: %s...\n", startTime.Format(time.RFC3339))
	for time.Now().Before(startTime) {
		time.Sleep(200 * time.Millisecond)
	}

	fmt.Println("开始执行固定操作数压力测试")
	var count int64 = 0
	var wg sync.WaitGroup
	sem := make(chan struct{}, *maxConc) // 并发控制

	startReal := time.Now()

	for i := 0; i < *totalOps; i++ {
		cmdStr := lines[i%len(lines)] // 循环使用命令
		sem <- struct{}{}
		wg.Add(1)
		go func(c string) {
			defer wg.Done()
			defer func() { <-sem }()
			_ = exec.Command("bash", "-c", strings.TrimSpace(c)).Run()
			atomic.AddInt64(&count, 1)
		}(cmdStr)
	}

	wg.Wait()
	duration := time.Since(startReal).Seconds()
	fmt.Println("测试结束")
	fmt.Printf("总操作数: %d\n", count)
	fmt.Printf("总耗时: %.2f 秒\n", duration)
	fmt.Printf("平均吞吐量: %.2f ops/sec\n", float64(count)/duration)
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if line := strings.TrimSpace(scanner.Text()); line != "" {
			lines = append(lines, line)
		}
	}
	return lines, scanner.Err()
}
