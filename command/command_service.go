package command

import (
	"fmt"
	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/storages"
	"os"
	"runtime"
	"strings"

	"gitee.com/quant1x/gox/daemon"
	"gitee.com/quant1x/gox/logger"
	"github.com/spf13/cobra"
)

var (
	serviceDescription      = "Quant1X量化系统数据服务"
	serviceCommand          = "" // 守护进程维护指令
	serviceProgramArguments = "service"
)

// Service is the daemon service struct
type Service struct {
	daemon daemon.Daemon
}

func (service *Service) Start() {
	// TODO: 启动服务需要做的事情
}

func (service *Service) Stop() {
	// TODO: 关闭服务需要做的事情
}

func (service *Service) Run() {
	// 运行服务
	storages.DaemonService()
}

// Manage by daemon commands or run the daemon
func (service *Service) Manage() (string, error) {
	if len(serviceCommand) > 1 {
		switch serviceCommand {
		case "install":
			return service.daemon.Install(serviceProgramArguments)
		case "remove", "uninstall":
			return service.daemon.Remove()
		case "start":
			return service.daemon.Start()
		case "stop":
			// No need to explicitly stop cron since job will be killed
			return service.daemon.Stop()
		case "status":
			return service.daemon.Status()
		default:
			//return usage, nil
		}
	}
	// serviceCommand = service
	return service.daemon.Run(service)
}

// CmdService 守护进程
var CmdService = &cobra.Command{
	Use:     "service",
	Example: Application + " service install | uninstall | remove | start | stop | status",
	Short:   "守护进程/服务",
	//Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			serviceCommand = args[0]
		}
		//logger.Infof("serviceCommand:%s", serviceCommand)
		daemonKind := daemon.SystemDaemon
		serviceName := Application
		switch runtime.GOOS {
		case "darwin":
			daemonKind = daemon.UserAgent
		case "windows":
			serviceName = "Quant1X-Stock"
		}
		srv, err := daemon.New(serviceName, serviceDescription, daemonKind)
		if err != nil {
			logger.Errorf("Error: %+v", err)
			fmt.Println("Error: ", err)
			os.Exit(1)
		}
		service := &Service{srv}
		replacer := strings.NewReplacer("${ROOT_PATH}", cache.GetRootPath(), "${LOG_PATH}", cache.GetLoggerPath())
		properties := replacer.Replace(propertyList)
		_ = service.daemon.SetTemplate(properties)
		status, err := service.Manage()
		if err != nil {
			logger.Errorf("Error: %+v", err)
			fmt.Println(status, "\nError: ", err)
			os.Exit(1)
		}
		fmt.Println(status)
	},
}
