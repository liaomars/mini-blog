package mini_blog

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/liaomars/mini-blog/internal/pkg/log"
	mw "github.com/liaomars/mini-blog/internal/pkg/middleware"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// 配置文件
var cfgFile string

func NewMiniBlogCommand() *cobra.Command {
	cmd := &cobra.Command{
		// 指定命令的名字，该名字会出现在帮助信息中
		Use: "mini-blog",
		//命令的简短描述
		Short: "一个好的Go语言练手项目",
		// 命令的详细描述
		Long: `A good Go practical project, used to create user with basic information.

Find more miniblog information at:
        https://github.com/marmotedu/miniblog#readme`,
		// 命令出错时，不打印帮助信息。不需要打印帮助信息，设置为 true 可以保持命令出错时一眼就能看到错误信息
		SilenceUsage: true,
		// 指定调用 cmd.Execute() 时，执行的 Run 函数，函数执行失败会返回错误信息
		RunE: func(cmd *cobra.Command, args []string) error {
			// 初始化日志
			log.Init(logOptions())
			// Sync 将缓存中的日志刷新到磁盘文件中
			defer log.Sync()
			return run()
		},
		Args: func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {
				if len(arg) > 0 {
					return fmt.Errorf("%q does not take any arguments,got %q", cmd.CommandPath(), args)
				}
			}
			return nil
		},
	}

	// 以下设置，使得 initConfig 函数在每个命令运行时都会被调用以读取配置
	cobra.OnInitialize(initConfig)

	// Cobra 支持持久性标志(PersistentFlag)，该标志可用于它所分配的命令以及该命令下的每个子命令
	cmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "项目的配置路径")

	cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	return cmd
}

func run() error {
	//fmt.Println("Hello MiniBlog!!!")

	// 初始化数据库链接
	if err := initStore(); err != nil {
		return err
	}

	setting, _ := json.Marshal(viper.AllSettings())
	log.Infow(string(setting))

	log.Infow(viper.GetString("db.database"))

	// 设置 Gin 模式
	gin.SetMode(viper.GetString("runmode"))
	// 初始化gin框架对象
	g := gin.New()

	// gin.Recovery() 中间件，用来捕获任何 panic，并恢复
	mws := []gin.HandlerFunc{gin.Recovery(), mw.Cors, mw.NoCache, mw.Secure, mw.RequestID()}

	g.Use(mws...)

	// 创建http server服务实例
	httpSrv := &http.Server{Addr: viper.GetString("addr"), Handler: g}

	log.Infow("Start to listening the incoming requests on http address", "addr", viper.GetString("addr"))
	go func() {
		if err := httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalw(err.Error())
		}
	}()

	// 等待中断信号优雅关闭服务器（10秒超时）
	quit := make(chan os.Signal, 1)
	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号，我们常用的 CTRL + C 就是触发系统 SIGINT 信号
	// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) //此处不阻塞
	<-quit                                               // 此处开始阻塞
	log.Infow("Shutting down server...")

	// 创建 ctx 用于通知服务器 goroutine, 它有 10 秒时间完成当前正在处理的请求
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 10 秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过 10 秒就超时退出
	if err := httpSrv.Shutdown(ctx); err != nil {
		log.Errorw("Insecure Server forced to shutdown", "err", err)
		return err
	}
	log.Infow("Server exiting")

	return nil
}
