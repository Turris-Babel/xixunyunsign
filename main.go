package main

import (
	"github.com/spf13/cobra"
	"log"
	"strconv"
	"xixunyunsign/cmd"
	"xixunyunsign/web"
)

func main() {
	var port int

	var rootCmd = &cobra.Command{Use: "xixun"}

	var ServeCmd = &cobra.Command{
		Use:   "serve",
		Short: "启动Web服务器",
		Run: func(cmd *cobra.Command, args []string) {
			// 设置默认端口为8080，如果没有通过命令行传递端口
			if port == 0 {
				port = 8080
			}
			server := web.NewServer()
			// 启动服务器，传递可变端口
			server.Run(":" + strconv.Itoa(port)) // 可以通过参数指定端口
		},
	}

	// 添加端口参数
	ServeCmd.Flags().IntVar(&port, "port", 8080, "指定Web服务器监听的端口")

	// 添加其他命令
	rootCmd.AddCommand(cmd.LoginCmd, cmd.QueryCmd, cmd.SignCmd, cmd.SchoolSearchIDCmd, cmd.ExperimentalCmd)
	rootCmd.AddCommand(ServeCmd) //rootCmd.CompletionOptions.DisableDefaultCmd = false
	//rootCmd.AddCommand(cmd.ScheduleCmd) // 添加 schedule 命令
	//// 设置优雅关闭信号
	//stopChan := make(chan os.Signal, 1)
	//signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)
	//
	//// 在单独的 goroutine 中执行命令
	//go func() {
	//	if err := rootCmd.Execute(); err != nil {
	//		log.Println(err)
	//		os.Exit(1)
	//	}
	//}()
	//
	//// 等待中断信号
	//<-stopChan
	//log.Println("收到停止信号，正在关闭调度器...")
	//
	//// 停止调度器
	//ctx := scheduler.StopScheduler()
	//<-ctx.Done()
	//log.Println("调度器已停止，程序退出")
	// 执行根命令
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}

}
