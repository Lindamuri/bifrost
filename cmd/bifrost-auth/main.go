package main

// DONE: 1.权限管理
// DONE: 2.nginx配置定期备份机制
// DONE: 3.日志规范化输出
// DONE: 4.优化守护进程输出，调整标准输出到守护进程日志中

import (
	"errors"
	"fmt"
	"os"

	"github.com/tremendouscan/bifrost/internal/pkg/auth/daemon"
)

func main() {
	defer daemon.Logf.Close()
	defer daemon.Stdoutf.Close()

	err := errors.New("unknown signal")
	switch *daemon.Signal {
	case "":
		err = daemon.Start()
		if err == nil {
			os.Exit(0)
		}
	case "stop":
		err = daemon.Stop()
		if err == nil {
			if os.Getppid() != 1 {
				fmt.Println("bifrost is stopping...")
			}
			os.Exit(0)
		}
	case "restart":
		err = daemon.Restart()
		if err == nil {
			os.Exit(0)
		}
	case "status":
		pid, statErr := daemon.Status()
		if statErr != nil {
			fmt.Printf("bifrost-auth is abnormal with error: %s\n", statErr.Error())
			os.Exit(1)
		} else {
			fmt.Printf("bifrost-auth <PID %d> is running\n", pid)
			os.Exit(0)
		}
	}
	fmt.Println(err.Error())
	os.Exit(1)
}
