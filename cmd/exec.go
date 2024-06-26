package cmd

import (
	"go-screenshot/utils"
	// "io/ioutil"
	"log"
	"runtime"
	"strconv"
	"time"
	// "os"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	// "github.com/go-rod/rod/lib/proto"
	"github.com/spf13/cobra"
	"fmt"
		//  "github.com/go-rod/rod/lib/proto"

)

// execCmd represents the exec command
var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "给我一个URL,我截图成功发给钉钉机器人🤖",
	Long:  `命令行截图小工具,可使用此工具订阅一些你关心的网页截图，然后添加到定时任务中。`,
	Run: func(cmd *cobra.Command, args []string) {
		url := cmd.Flags().Lookup("url").Value.String()
		element := cmd.Flags().Lookup("element").Value.String()
		webhook := cmd.Flags().Lookup("webhook").Value.String()
		width, _ := strconv.Atoi(cmd.Flags().Lookup("kuan").Value.String())
		height, _ := strconv.Atoi(cmd.Flags().Lookup("gao").Value.String())
		// webhookURL := cmd.Flags().Lookup("webhook").Value.String()
		launch := launcher.New().Headless(true)
		if runtime.GOOS == "darwin" {
			launch = launch.Bin(`/Applications/Google Chrome.app/Contents/MacOS/Google Chrome`)
		}
		if runtime.GOOS == "linux" {
			launch = launch.Set("--no-sandbox")
		}
		page := rod.New().ControlURL(launch.MustLaunch()).MustConnect().MustPage()

		defer page.Close()
		withTimeout := page.Timeout(2*time.Minute)
		withTimeout.
			MustSetViewport(width, height, 1, false).
			MustNavigate(url).
			MustWaitLoad().
			MustWindowMaximize()

		// 自动登录阿里云
		// log.Println("开始自动登录阿里云")
		// page.MustElement("#fm-login-id").MustInput("hzmpm")
		// page.MustElement("#fm-login-password").MustInput("mvwnTGPN71BUOx/2uTiP")
		// page.MustElement("#login-form > div.fm-btn > button").MustClick()
		// page.MustWaitLoad()
		// log.Println("登录成功，等待页面加载完成")

		// 添加日志输出以确认页面导航成功

		// 添加日志输出以确认页面加载完成
		log.Println("Waiting for page to be idle and ready for screenshot...")

		// 等待页面请求空闲
		page.WaitRequestIdle(time.Duration(time.Second*10), []string{}, []string{},nil)()
		log.Println("确认元素是否存在")
		log.Println("Page is idle and ready for screenshot")
		elementExists := page.Timeout(30 * time.Second).MustElement(element).MustWaitVisible()
		

		if elementExists == nil {
			log.Fatal("元素未找到或不可见:", element)
		}
		el := elementExists.MustScreenshot()
		// 上传截图到阿里云OSS并返回URL
		url, err := utils.UploadToOSS(el)
		if err != nil {
			log.Fatalf("Failed to upload screenshot: %v", err)
		}		
		fmt.Println(url)
		utils.SendImageToDingtalk(url,webhook)

	},
}

func init() {
	fset := execCmd.Flags()
	fset.StringP("url", "u", "https://baidu.com", "给我一个你想要截图的URL")
	fset.StringP("element", "e", "#s_lg_img_new", "给我你关心的页面元素")
	fset.StringP("kuan", "k", "1200", "页面宽度")
	fset.StringP("gao", "g", "800", "页面高度")
	fset.StringP("webhook", "w", "", "钉钉机器人Webhook地址")
	execCmd.MarkFlagRequired("url")
	execCmd.MarkFlagRequired("element")
	execCmd.MarkFlagRequired("webhook")
	rootCmd.AddCommand(execCmd)
}