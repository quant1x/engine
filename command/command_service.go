package command

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"
	"syscall"

	"gitee.com/quant1x/engine/cache"
	"gitee.com/quant1x/engine/services"
	"gitee.com/quant1x/gox/daemon"
	"gitee.com/quant1x/gox/logger"
	nix "github.com/sevlyar/go-daemon"
	cmder "github.com/spf13/cobra"
)

const (
	serviceCommand = "service"
)

var (
	// CmdService 守护进程
	CmdService *cmder.Command
)

// FirstUpper 字符串首字母大写
func FirstUpper(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// FirstLower 字符串首字母小写
func FirstLower(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToLower(s[:1]) + s[1:]
}

func initService() {
	CmdService = &cmder.Command{
		Use:     serviceCommand,
		Example: Application + " " + serviceCommand + " install | uninstall | remove | start | stop | list | status",
		Short:   "守护进程/服务",
		//Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cmder.Command, args []string) {
			if len(args) > 0 {
				serviceSubCommand = args[0]
			}
			//logger.Infof("serviceCommand:%s", serviceCommand)
			daemonKind := daemon.SystemDaemon
			applicationName, _, _ := strings.Cut(Application, ".")
			serviceName := applicationName
			switch runtime.GOOS {
			case "darwin":
				daemonKind = daemon.UserAgent
			case "windows":
				//serviceName = "Quant1X-Stock"
				serviceName = "Quant1X-" + FirstUpper(applicationName)
				serviceDescription += " V" + MinVersion
			}
			srv, err := daemon.New(serviceName, serviceDescription, daemonKind)
			if err != nil {
				logger.Errorf("Error: %+v", err)
				fmt.Println("Error: ", err)
				os.Exit(1)
			}
			service := &Service{srv}
			switch runtime.GOOS {
			case "linux": // linux使用守护进程方式, 不需要安装
				cntxt := &nix.Context{
					PidFileName: cache.GetVariablePath() + "/" + Application + ".pid",
					PidFilePerm: 0644,
					LogFileName: cache.GetVariablePath() + "/" + Application + ".log",
					LogFilePerm: 0640,
					WorkDir:     cache.GetVariablePath(),
					Umask:       027,
					Args:        []string{os.Args[0], serviceCommand},
				}
				logger.Warnf("stock args:%+v", os.Args)
				if len(serviceSubCommand) > 1 {
					switch serviceSubCommand {
					case "install":
						//return service.daemon.Install(serviceProgramArguments)
						fmt.Println("success")
					case "remove", "uninstall":
						//return service.daemon.Remove()
						fmt.Println("success")
					case "start":
						logger.Warnf("start service")
						d, err := cntxt.Reborn()
						if err != nil {
							fmt.Println("Unable to run: ", err)
							os.Exit(1)
						}
						if d != nil {
							//fmt.Println("Unable to run")
							//os.Exit(1)
							return
						}

						fmt.Println("- - - - - - - - - - - - - - -")
						fmt.Println("daemon started")
						//service.daemon.Run(service)
					case "stop":
						// No need to explicitly stop cron since job will be killed
						//return service.daemon.Stop()
						//coroutine.Shutdown()
						d, err := cntxt.Search()
						if err != nil {
							logger.Fatalf("Unable send signal to the daemon: %s", err.Error())
						}
						if d == nil {
							logger.Fatalf("Not found %s", applicationName)
						}
						err = d.Signal(syscall.SIGQUIT)
						if err != nil {
							logger.Fatalf("send signal failed: %s", err.Error())
						}
					case "list":
						services.PrintJobList()
						//return "", nil
					case "status":
						//return service.daemon.Status()
						fmt.Println("success")
					default:
						//status, err := service.daemon.Run(service)
						//fmt.Println(status, err)
					}
				} else {
					logger.Warnf("run service")
					d, err := cntxt.Reborn()
					if err != nil {
						fmt.Println("Unable to run: ", err)
						os.Exit(1)
					}
					if d != nil {
						//fmt.Println("Unable to run")
						//os.Exit(1)
						return
					}
					defer cntxt.Release()
					services.DaemonService()
				}

			default:
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
			}
		},
	}
}

var (
	serviceDescription      = "Quant1X量化系统数据服务"
	serviceSubCommand       = "" // 守护进程维护指令
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
	services.DaemonService()
}

// Manage by daemon commands or run the daemon
func (service *Service) Manage() (string, error) {
	if len(serviceSubCommand) > 1 {
		switch serviceSubCommand {
		case "install":
			return service.daemon.Install(serviceProgramArguments)
		case "remove", "uninstall":
			return service.daemon.Remove()
		case "start":
			return service.daemon.Start()
		case "stop":
			// No need to explicitly stop cron since job will be killed
			return service.daemon.Stop()
		case "list":
			services.PrintJobList()
			return "", nil
		case "status":
			return service.daemon.Status()
		default:
			_ = CmdService.Usage()
			return "", errors.New("unknown service flags=" + serviceSubCommand)
		}
	} else {
		// serviceCommand = service
		return service.daemon.Run(service)
	}
}
