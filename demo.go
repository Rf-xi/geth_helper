package main

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	gomail "gopkg.in/gomail.v2"
)

var target_time string = "09/03/2019, 00:00:00"

func email(mess string) {
	fmt.Printf("send email %v", mess)
	msg := gomail.NewMessage()
	msg.SetHeader("From", "18030129824@163.com")
	msg.SetHeader("To", "2997080200@qq.com")
	msg.SetHeader("Subject", "geth 同步进展")
	msg.SetBody("text/html", fmt.Sprintf("<b>Notify: %s</b>", mess))
	// msg.Attach("/home/User/cat.jpg")

	n := gomail.NewDialer("smtp.163.com", 465, "18030129824@163.com", "ICNSMFTJSWFMXRFQ")

	// Send the email
	if err := n.DialAndSend(msg); err != nil {
		panic(err)
	}

	// // 配置发送邮件的参数
	// from := "aws geth downloader demo(xrf)"
	// to := []string{"2997080200@qq.com"}
	// subject := "geth 同步监听程序"
	// body := mess

	// // 构建邮件内容
	// message := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", to[0], subject, body))

	// // 配置 SMTP 服务器信息
	// auth := smtp.PlainAuth("", "18030129824@163.com", "123456aA...", "smtp.163.com")

	// // 发送邮件
	// err := smtp.SendMail("smtp.163.com:25", auth, from, to, message)
	// if err != nil {
	// 	fmt.Println("Failed to send email:", err)
	// 	return
	// }

	fmt.Println("Email sent successfully!")
}

func lastest() (string, bool) {
	cmd := exec.Command("./download-helper.sh", "query")

	// 获取命令输出
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error executing command:", err)
		return "Err", false
	}
	fmt.Printf("result %v", string(output))
	// 输出结果
	return string(output), true
}

// 从字符串中提取包含时间信息的子串
func extractTimeString(s string) string {
	layout := "01/02/2006, 15:04:05"
	for i := 0; i < len(s); i++ {
		for j := i + len(layout); j <= len(s); j++ {
			if _, err := time.Parse(layout, s[i:j]); err == nil {
				return s[i:j]
			}
		}
	}
	return ""
}
func prasetime(mess string) (time.Time, bool) {
	fmt.Printf("prasetime %v\n", mess)
	layout := "01/02/2006, 15:04:05"
	dateTime, err := time.Parse(layout, extractTimeString(mess))
	if err != nil {
		fmt.Println("Error parsing datetime:", err)
		return dateTime, false
	}
	return dateTime, true
}

func shutdown() (string, bool) {
	cmd := exec.Command("./download-helper.sh", "stop")

	// 获取命令输出
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error executing command:", err)
		return "Err", false
	}

	time.Sleep(time.Second)

	if strings.Contains(string(output), "Stopped the Tmux program.") {
		// 输出结果
		return string(output), true
	}

	return string(output), false
}

func server() {
	// 创建一个 10 分钟间隔的定时器
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	times := 0
	// 在一个无限循环中处理定时器事件

	target, _ := prasetime(target_time)
	fmt.Printf("target %v", target)

	for range ticker.C {
		// 每 10 分钟触发一次
		mess, ok := lastest()
		if !ok {
			email("lastest err: 获取最新块失败")
		}
		cur_time, ok := prasetime(mess)
		if !ok {
			email("prasetime err: 解析时间失败")
		}
		if cur_time.After(target) {
			info, ok := shutdown()
			if ok {
				email(fmt.Sprintf("到达目标位置： %v\n %v", cur_time, info))
			} else {
				email(fmt.Sprintf("geth 停止失败！！！到达目标位置： %v\n %v", cur_time, info))
			}
			email("********************demo exit**********************")
			return
		}
		times++

		if times%6 == 0 { // 每一个小时发送一次同步情况
			email(fmt.Sprintf("geth 同步情况, lastest %v", cur_time))
		}
	}
}

func main() {
	server()
}
