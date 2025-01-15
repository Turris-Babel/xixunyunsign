package main

import (
	"github.com/spf13/cobra"
	"log"
	"xixunyunsign/cmd"
	"xixunyunsign/utils"
	"xixunyunsign/web"
)

func main() {
	// 初始化数据库
	err := utils.InitDB()
	if err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}
	//
	//// 初始化定时任务调度器
	//_, err = utils.InitScheduler(utils.PerformSign)
	//if err != nil {
	//	log.Printf("初始化定时任务调度器失败: %v\n", err)
	//}

	// 设置根命令
	var rootCmd = &cobra.Command{Use: "xixun"}
	var ServeCmd = &cobra.Command{
		Use:   "serve",
		Short: "启动Web服务器",
		Run: func(cmd *cobra.Command, args []string) {
			server := web.NewServer()
			server.Run(":8080") // 可以通过参数或配置文件指定端口
		},
	}
	rootCmd.AddCommand(cmd.LoginCmd, cmd.QueryCmd, cmd.SignCmd, cmd.SchoolSearchIDCmd, cmd.ExperimentalCmd)
	rootCmd.AddCommand(ServeCmd)
	//rootCmd.CompletionOptions.DisableDefaultCmd = false
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
