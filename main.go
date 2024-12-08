package main

import (
	"github.com/spf13/cobra"
	"log"
	"xixuanyunsign/cmd"
	"xixuanyunsign/utils"
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
	rootCmd.AddCommand(cmd.LoginCmd)
	rootCmd.AddCommand(cmd.QueryCmd)
	rootCmd.AddCommand(cmd.SignCmd)
	rootCmd.AddCommand(cmd.SchoolSearchIDCmd)
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
	rootCmd.Execute()
}
