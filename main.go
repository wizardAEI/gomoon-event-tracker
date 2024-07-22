package main

import (
	"bufio"
	"math"
	"os"
	"strings"

	hook "github.com/robotn/gohook"
)

func main() {

	// 记录快速两次 ctrl+c 按键
	preStamp := int64(0)
	stamp := int64(0)
	defer hook.End()

	// 查看Access

	// win
	hook.Register(hook.KeyDown, []string{"ctrl", "c"}, func(e hook.Event) {
		stamp = e.When.UnixNano()
		// 间隔小于 0.4s 认为是快速两次按键
		if stamp-preStamp < 400000000 {
			os.Stdout.WriteString("multi-copy")
		}
		preStamp = stamp
	})
	// mac
	hook.Register(hook.KeyDown, []string{"cmd", "c"}, func(e hook.Event) {
		stamp = e.When.UnixNano()
		// 间隔小于 0.4s 认为是快速两次按键
		if stamp-preStamp < 400000000 {
			os.Stdout.WriteString("multi-copy")
		}
		preStamp = stamp
	})

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
		}
	})
	preMouseDown := [2]int{-1, -1}
	hook.Register(hook.MouseDown, []string{}, func(e hook.Event) {
		if math.Abs(float64(e.X)-float64(preMouseDown[0])) <= 5 && math.Abs(float64(e.Y)-float64(preMouseDown[1])) <= 5 {
			isDragged = true
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
