package main

import (
	"bufio"
	"flag"
	"math"
	"os"
	"runtime"
	"strings"
	"time"

	hook "github.com/robotn/gohook"
)

// FEAT:
// 1. ctrl+c+启动键 快速问答
// 2. 预测用户是否选中文本
func main() {

	// 读取启动参数
	key := flag.String("key", "", "the key to be pressed")
	ctrlKey := "ctrl"

	if runtime.GOOS == "darwin" {
		ctrlKey = "cmd"
	}
	flag.Parse()
	// 记录快速两次 ctrl+c 按键
	preStamp := int64(0)
	stamp := int64(0)
	defer hook.End()

	if *key == "C" {
		hook.Register(hook.KeyDown, []string{ctrlKey, "c"}, func(e hook.Event) {
			stamp = e.When.UnixNano()
			// 间隔小于 0.4s 认为是快速两次按键
			if stamp-preStamp < 400000000 {
				os.Stdout.WriteString("quickly-ans")
			}
			preStamp = stamp
		})
	} else {
		hook.Register(hook.KeyDown, []string{ctrlKey, "c"}, func(e hook.Event) {
			preStamp = e.When.UnixNano()
		})
		hook.Register(hook.KeyDown, []string{ctrlKey, strings.ToLower(*key)}, func(e hook.Event) {
			stamp = e.When.UnixNano()
			// 间隔小于 0.8s 认为是快速两次按键
			if stamp-preStamp < 800000000 {
				os.Stdout.WriteString("quickly-ans")
			}
		})
	}

	// 预测用户是否选中文本
	isDragged := false
	mousePosition := [2]int{-1, -1}
	hook.Register(hook.MouseHold, []string{}, func(e hook.Event) {
		isDragged = false
		if e.Button == hook.MouseMap["left"] {
			mousePosition[0], mousePosition[1] = int(e.X), int(e.Y)
		}
	})
	hook.Register(hook.MouseDrag, []string{}, func(e hook.Event) {
		x, y := int(e.X), int(e.Y)
		if mousePosition[0] == -1 || mousePosition[1] == -1 {
			return
		}
		if math.Abs(float64(x)-float64(mousePosition[0])) <= 14 && math.Abs(float64(y)-float64(mousePosition[1])) <= 14 {
			isDragged = false
			return
		} else {
			isDragged = true
			time.AfterFunc(3*time.Second, func() {
				isDragged = false
			})
		}
	})
	preMouseDown := [2]int{-1, -1}
	hook.Register(hook.MouseDown, []string{}, func(e hook.Event) {
		if math.Abs(float64(e.X)-float64(preMouseDown[0])) <= 5 && math.Abs(float64(e.Y)-float64(preMouseDown[1])) <= 5 {
			isDragged = true
			time.AfterFunc(3*time.Second, func() {
				isDragged = false
			})
		}
		if e.Button == hook.MouseMap["left"] {
			preMouseDown[0], preMouseDown[1] = int(e.X), int(e.Y)
			mousePosition[0], mousePosition[1] = -1, -1
		}
	})

	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			text, _ := reader.ReadString('\n')
			text = strings.TrimSpace(text) // 去掉换行符和空白字符
			switch text {
			case "isDragged":
				if isDragged {
					isDragged = false
					os.Stdout.WriteString("true\n")
				} else {
					os.Stdout.WriteString("false\n")
				}
			}
		}
	}()

	s := hook.Start()
	<-hook.Process(s)
}
