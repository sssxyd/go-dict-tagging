package main

import (
	"dict_tagging/dict"
	"dict_tagging/funcs"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/pelletier/go-toml/v2"
	"golang.org/x/sys/windows"
)

type Config struct {
	Server struct {
		Port int `toml:"port"`
	} `toml:"server"`
	App struct {
		RootDir   string
		DictDir   string `toml:"dict_dir"`
		LogPath   string `toml:"log_path"`
		StaticDir string `toml:"static_dir"`
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
	log.SetFlags(log.LstdFlags)

	// 设置Windows控制台为UTF-8编码
	if os.Getenv("OS") == "Windows_NT" {
		handle := windows.Handle(os.Stdout.Fd())
		var mode uint32
		windows.GetConsoleMode(handle, &mode)
		mode |= windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING
		windows.SetConsoleMode(handle, mode)
	}

	// 读取配置文件
	baseDir := funcs.GetExecutionPath()
	file, err := os.Open(filepath.Join(baseDir, "config.toml"))
	if err != nil {
		log.Fatalf("Failed to open config file: %v", err)
	}
	defer file.Close()
	decoder := toml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		log.Fatalf("Failed to decode config file: %v", err)
	}

	config.App.RootDir = baseDir

	// 初始化目录及日志
	if !filepath.IsAbs(config.App.DictDir) {
		config.App.DictDir = filepath.Join(baseDir, config.App.DictDir)
	}
	funcs.TouchDir(config.App.DictDir)

	if !filepath.IsAbs(config.App.StaticDir) {
		config.App.StaticDir = filepath.Join(baseDir, config.App.StaticDir)
	}
	funcs.TouchDir(config.App.StaticDir)

	if !filepath.IsAbs(config.App.LogPath) {
		config.App.LogPath = filepath.Join(baseDir, config.App.LogPath)
	}
	funcs.InitializeLogFile(config.App.LogPath, true)

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
