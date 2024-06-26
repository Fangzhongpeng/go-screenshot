package cmd

import (
	// "care-screenshot/public"
	// "io/ioutil"
	"log"
	"runtime"
	"strconv"
	"time"
	"os"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	// "github.com/go-rod/rod/lib/proto"
	"github.com/spf13/cobra"
		//  "github.com/go-rod/rod/lib/proto"

)

// execCmd represents the exec command
var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "给我一个URL,我截图成功发给钉钉机器人🤖",
	Long:  `命令行工具，可使用此工具订阅一些你关心的网页服务状态，然后添加到定时任务中。`,
	Run: func(cmd *cobra.Command, args []string) {
		url := cmd.Flags().Lookup("url").Value.String()
		element := cmd.Flags().Lookup("element").Value.String()
		width, _ := strconv.Atoi(cmd.Flags().Lookup("kuan").Value.String())
		height, _ := strconv.Atoi(cmd.Flags().Lookup("gao").Value.String())
		// webhookURL := cmd.Flags().Lookup("webhook").Value.String()
		log.Println("Successfully navigated to:", "1")
		launch := launcher.New().Headless(true)
		if runtime.GOOS == "darwin" {
			launch = launch.Bin(`/Applications/Google Chrome.app/Contents/MacOS/Google Chrome`)
			log.Println("Successfully navigated to:", "2")
		}
		if runtime.GOOS == "linux" {
			launch = launch.Set("--no-sandbox")
		}
		page := rod.New().ControlURL(launch.MustLaunch()).MustConnect().MustPage()

		log.Println("Successfully navigated to:", "3")
		defer page.Close()
		log.Println("Successfully navigated to:", "4")
		withTimeout := page.Timeout(2*time.Minute)
		log.Println("Successfully navigated to:", "5")
		withTimeout.
			MustSetViewport(width, height, 1, false).
			MustNavigate(url).
			MustWaitLoad().
			MustWindowMaximize()

		// 添加日志输出以确认页面导航成功
		log.Println("Successfully navigated to:", url)

		// 添加日志输出以确认页面加载完成
		log.Println("Waiting for page to be idle and ready for screenshot...")

		// 等待页面请求空闲
		page.WaitRequestIdle(time.Duration(time.Second*10), []string{}, []string{},nil)()
		log.Println("确认元素是否存在")
		log.Println("Page is idle and ready for screenshot")
		log.Println("Successfully navigated to:", "6")
		elementExists := page.Timeout(30 * time.Second).MustElement(element).MustWaitVisible()
		

		if elementExists == nil {
			log.Fatal("元素未找到或不可见:", element)
		}
		el := elementExists.MustScreenshot()


		
		
		// log.Println("Successfully navigated to:", "7")
		if err := os.WriteFile("tmp.png", el, 0644); err != nil {
			log.Fatal("Failed to write screenshot file:", err)
		}
		log.Println("Successfully navigated to:", "7")
		log.Println("Screenshot saved to tmp.png")

		// public.SendImage("tmp.png", webhookURL)
	},
}

func init() {
	fset := execCmd.Flags()
	fset.StringP("url", "u", "https://baidu.com", "给我一个你想要截图的URL")
	fset.StringP("element", "e", "#s_lg_img", "给我你关心的页面元素")
	fset.StringP("kuan", "k", "1200", "页面宽度")
	fset.StringP("gao", "g", "800", "页面高度")
	fset.StringP("webhook", "w", "", "钉钉机器人Webhook地址")
	execCmd.MarkFlagRequired("url")
	execCmd.MarkFlagRequired("element")
	execCmd.MarkFlagRequired("webhook")
	rootCmd.AddCommand(execCmd)
}
