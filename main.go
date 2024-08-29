package main

import (
	"dict_tagging/dict"
	"dict_tagging/funcs"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	Server struct {
		Port int `toml:"port"`
	} `toml:"server"`
	App struct {
		RootDir       string
		DictDir       string `toml:"dict_dir"`
		StaticDir     string `toml:"static_dir"`
		AccessLogPath string `toml:"access_log_path"`
		ErrorLogPath  string `toml:"error_log_path"`
		AppLogPath    string `toml:"app_log_path"`
	} `toml:"app"`
}

type ApiResult struct {
	Code   int         `json:"code"`
	Msg    string      `json:"msg"`
	Result interface{} `json:"result"`
	Micros int         `json:"micros"`
}

var (
	// 全局变量
	root          *dict.TrieNode
	infos         []dict.DictInfo
	config        Config
	dictWriteLock sync.Mutex
)

func init() {
	// 设置Windows控制台为UTF-8编码
	// if os.Getenv("OS") == "Windows_NT" {
	// 	handle := windows.Handle(os.Stdout.Fd())
	// 	var mode uint32
	// 	windows.GetConsoleMode(handle, &mode)
	// 	mode |= windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING
	// 	windows.SetConsoleMode(handle, mode)
	// }

	// 读取配置文件
	baseDir := funcs.GetExecutionPath()
	file, err := os.Open(filepath.Join(baseDir, "config.toml"))
	if err != nil {
		fmt.Printf("Failed to open config file: %v", err)
	}
	defer file.Close()
	decoder := toml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		fmt.Printf("Failed to decode config file: %v", err)
	}

	config.App.RootDir = baseDir

	// dict dir
	if !filepath.IsAbs(config.App.DictDir) {
		config.App.DictDir = filepath.Join(baseDir, config.App.DictDir)
	}
	funcs.TouchDir(config.App.DictDir)

	// static dir
	if !filepath.IsAbs(config.App.StaticDir) {
		config.App.StaticDir = filepath.Join(baseDir, config.App.StaticDir)
	}
	funcs.TouchDir(config.App.StaticDir)

	// access log path
	if !filepath.IsAbs(config.App.AccessLogPath) {
		config.App.AccessLogPath = filepath.Join(baseDir, config.App.AccessLogPath)
	}
	funcs.TouchDir(filepath.Dir(config.App.AccessLogPath))
	accessLogFile, err := os.OpenFile(config.App.AccessLogPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer file.Close()
	gin.DefaultWriter = io.MultiWriter(accessLogFile, os.Stdout)

	// error log path
	if !filepath.IsAbs(config.App.ErrorLogPath) {
		config.App.ErrorLogPath = filepath.Join(baseDir, config.App.ErrorLogPath)
	}
	funcs.TouchDir(filepath.Dir(config.App.ErrorLogPath))
	errorLogFile, err := os.OpenFile(config.App.ErrorLogPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer file.Close()
	gin.DefaultErrorWriter = io.MultiWriter(errorLogFile, os.Stderr)

	// app log path
	if !filepath.IsAbs(config.App.AppLogPath) {
		config.App.AppLogPath = filepath.Join(baseDir, config.App.AppLogPath)
	}
	funcs.InitializeLogFile(config.App.AppLogPath, true)

	// 加载字典
	root, infos = dict.LoadData()
}

func main() {
	// 设置 Gin 运行模式为 release
	gin.SetMode(gin.ReleaseMode)

	// 创建Gin引擎
	engine := gin.Default()

	engine.Static("/static", config.App.StaticDir)
	engine.GET("/", func(ctx *gin.Context) {
		ctx.File(filepath.Join(config.App.StaticDir, "index.html"))
	})
	engine.GET("/index.html", func(ctx *gin.Context) {
		ctx.File(filepath.Join(config.App.StaticDir, "index.html"))
	})
	engine.GET("/favicon.ico", func(ctx *gin.Context) {
		ctx.File(filepath.Join(config.App.StaticDir, "favicon.ico"))
	})

	handleTag(engine)
	handlePut(engine)
	handleReload(engine)
	handleInfo(engine)

	engine.Run(fmt.Sprintf(":%d", config.Server.Port))
}
