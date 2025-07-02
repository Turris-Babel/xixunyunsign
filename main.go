package main

import "C"
import (
	"github.com/spf13/cobra"
	"log"
	"os"
	"strconv"
	"xixunyunsign/cmd"
	"xixunyunsign/internal/wire" // Import the wire package
)

// simple aplit args
func splitArgs(s string) []string {
	var result []string
	var current string
	inQuote := false
	for _, r := range s {
		if r == '"' {
			inQuote = !inQuote
			continue
		}
		if r == ' ' && !inQuote {
			if len(current) > 0 {
				result = append(result, current)
			}
			current = ""
		} else {
			current += string(r)
		}
	}
	if len(current) > 0 {
		result = append(result, current)
	}
	return result
}

//export RunCommand
func RunCommand(args *C.char) {
	// Convert C string to Go string
	goArgs := splitArgs(C.GoString(args))
	log.Printf("Received args: %#v\n", goArgs)

	// Set the arguments for the root command
	rootCmd.SetArgs(goArgs)

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
	}
}

var rootCmd = &cobra.Command{Use: "xixun"}

var ServeCmd = &cobra.Command{
	Use:   "serve",
	Short: "启动Web服务器",
	Run: func(cmd *cobra.Command, args []string) {
		// 设置默认端口为8080，如果没有通过命令行传递端口
		var port int
		port, _ = cmd.Flags().GetInt("port")
		if port == 0 {
			port = 8080
		}
		// Initialize server using wire injector
		server, err := wire.InitializeServer()
		if err != nil {
			log.Fatalf("无法初始化服务器: %v", err)
		}
		// 启动服务器，传递可变端口
		server.Run(":" + strconv.Itoa(port)) // 可以通过参数指定端口
	},
}

func init() {
	rootCmd.SetOut(os.Stdout)
	rootCmd.SetErr(os.Stderr)

	// 添加端口参数
	ServeCmd.Flags().Int("port", 8080, "指定Web服务器监听的端口")

	// 添加其他命令
	rootCmd.AddCommand(cmd.LoginCmd, cmd.QueryCmd, cmd.SignCmd, cmd.SchoolSearchIDCmd, cmd.ExperimentalCmd)
	rootCmd.AddCommand(ServeCmd)
}

func main() {
	// The main function is only called when running as a standalone executable
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
