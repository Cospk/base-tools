package config

import (
	"fmt"
	"log"
)

// ExampleConfig 示例配置结构体
type ExampleConfig struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Log      LogConfig      `mapstructure:"log"`
}

type ServerConfig struct {
	Host    string `mapstructure:"host"`
	Port    int    `mapstructure:"port"`
	Timeout int    `mapstructure:"timeout"`
}

type DatabaseConfig struct {
	Driver          string `mapstructure:"driver"`
	DSN             string `mapstructure:"dsn"`
	MaxConnections  int    `mapstructure:"max_connections"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"`
}

type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

type LogConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
	Output string `mapstructure:"output"`
}

// Example1_BasicUsage 基本使用示例
func Example1_BasicUsage() {
	// 创建配置管理器
	vc := NewViperConfig(
		WithConfigName("app"),
		WithConfigType("yaml"),
		WithConfigPath("./config", "."),
		WithEnvPrefix("APP"),
	)

	// 加载配置
	if err := vc.Load(); err != nil {
		log.Printf("警告: 无法加载配置文件: %v", err)
	}

	// 解析到结构体
	var cfg ExampleConfig
	if err := vc.Unmarshal(&cfg); err != nil {
		log.Fatalf("解析配置失败: %v", err)
	}

	fmt.Printf("服务器配置: %+v\n", cfg.Server)
}

// Example2_QuickLoad 快速加载示例
func Example2_QuickLoad() {
	var cfg ExampleConfig
	
	// 方式1: 指定配置文件路径
	err := QuickLoad("./config/app.yaml", &cfg, "APP")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 方式2: 自动搜索配置文件
	err = QuickLoad("", &cfg, "APP")
	if err != nil {
		log.Printf("使用默认配置")
	}
}

// Example3_EnvOverride 环境变量覆盖示例
func Example3_EnvOverride() {
	vc := NewViperConfig(
		WithConfigName("app"),
		WithConfigType("yaml"),
		WithConfigPath("./config"),
		WithEnvPrefix("APP"), // 环境变量前缀
	)

	// 设置默认值
	vc.SetDefault("server.port", 8080)
	vc.SetDefault("database.max_connections", 100)

	// 加载配置文件
	_ = vc.Load()

	// 环境变量会覆盖配置文件的值
	// 例如: APP_SERVER_PORT=9090 会覆盖 server.port
	port := vc.GetInt("server.port")
	fmt.Printf("服务器端口: %d\n", port)
}

// Example4_WatchConfig 监听配置变化示例
func Example4_WatchConfig() {
	vc := NewViperConfig(
		WithConfigName("app"),
		WithConfigType("yaml"),
		WithConfigPath("./config"),
	)

	if err := vc.Load(); err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	var cfg ExampleConfig
	_ = vc.Unmarshal(&cfg)

	// 注册配置变更回调
	vc.OnConfigChange(func() {
		fmt.Println("配置已更新，重新加载...")
		if err := vc.Unmarshal(&cfg); err != nil {
			log.Printf("重新加载配置失败: %v", err)
		} else {
			fmt.Printf("新配置: %+v\n", cfg)
		}
	})

	// 开始监听配置文件变化
	vc.WatchConfig()
}

// Example5_GlobalConfig 全局配置示例
func Example5_GlobalConfig() {
	// 初始化全局配置
	err := InitGlobalConfig(
		WithConfigName("app"),
		WithConfigType("yaml"),
		WithConfigPath("./config"),
		WithEnvPrefix("APP"),
	)
	if err != nil {
		log.Printf("初始化配置失败: %v", err)
	}

	// 在任何地方使用全局配置
	gc := GetGlobalConfig()
	port := gc.GetInt("server.port")
	dbDSN := gc.GetString("database.dsn")
	
	fmt.Printf("端口: %d, 数据库: %s\n", port, dbDSN)
}

// Example6_MultipleConfigs 多配置文件示例
func Example6_MultipleConfigs() {
	// 主配置
	mainConfig := NewViperConfig(
		WithConfigName("app"),
		WithConfigType("yaml"),
		WithConfigPath("./config"),
	)
	_ = mainConfig.Load()

	// 数据库专用配置
	dbConfig := NewViperConfig(
		WithConfigName("database"),
		WithConfigType("yaml"),
		WithConfigPath("./config"),
	)
	_ = dbConfig.Load()

	// 合并配置（数据库配置优先级更高）
	_ = MergeConfig(mainConfig, dbConfig)

	// 使用合并后的配置
	var cfg ExampleConfig
	_ = mainConfig.Unmarshal(&cfg)
}

// Example7_DynamicConfig 动态配置示例
func Example7_DynamicConfig() {
	vc := NewViperConfig(
		WithConfigName("app"),
		WithConfigType("yaml"),
		WithConfigPath("./config"),
	)
	_ = vc.Load()

	// 动态设置配置
	vc.Set("feature.enabled", true)
	vc.Set("api.rate_limit", 1000)

	// 保存配置到文件
	if err := vc.WriteConfig(); err != nil {
		// 如果文件不存在，使用 WriteConfigAs
		if err := vc.WriteConfigAs("./config/app_new.yaml"); err != nil {
			log.Printf("保存配置失败: %v", err)
		}
	}
}

// Example8_ConfigValidation 配置验证示例
func Example8_ConfigValidation() {
	vc := NewViperConfig(
		WithConfigName("app"),
		WithConfigType("yaml"),
		WithConfigPath("./config"),
	)
	_ = vc.Load()

	// 验证必需的配置项
	requiredKeys := []string{
		"server.host",
		"server.port",
		"database.dsn",
	}

	for _, key := range requiredKeys {
		if !vc.IsSet(key) {
			log.Fatalf("缺少必需的配置项: %s", key)
		}
	}

	// 验证配置值的有效性
	port := vc.GetInt("server.port")
	if port < 1 || port > 65535 {
		log.Fatalf("无效的端口号: %d", port)
	}
}

// Example9_ConfigDebug 配置调试示例
func Example9_ConfigDebug() {
	vc := NewViperConfig(
		WithConfigName("app"),
		WithConfigType("yaml"),
		WithConfigPath("./config"),
		WithEnvPrefix("APP"),
	)

	// 设置一些默认值
	vc.SetDefault("debug.enabled", false)
	vc.SetDefault("debug.verbose", false)

	_ = vc.Load()

	// 打印调试信息
	vc.Debug()

	// 获取特定配置部分
	serverConfig := vc.GetStringMap("server")
	fmt.Printf("服务器配置: %+v\n", serverConfig)

	// 获取所有键
	allKeys := vc.AllSettings()
	for key := range allKeys {
		fmt.Printf("配置键: %s\n", key)
	}
}
