# go-screenshot

go语言截图小工具


go语言截图小工具

go run main.go exec -h
命令行截图小工具,可使用此工具订阅一些你关心的网页截图，然后添加到定时任务中。

Usage:
  go-screenshot exec [flags]

Flags:
  -e, --element string   给我你关心的页面元素 (default "#s_lg_img_new")
  -g, --gao string       页面高度 (default "800")
  -h, --help             help for exec
  -k, --kuan string      页面宽度 (default "1200")
  -u, --url string       给我一个你想要截图的URL (default "https://baidu.com")
  -w, --webhook string   钉钉机器人Webhook地址
