package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Cospk/base-tools/errs"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// ViperConfig 基于 Viper 的配置管理器
type ViperConfig struct {
	viper      *viper.Viper
	configName string
	configType string
	configPath []string
	envPrefix  string
	// 配置变更回调函数
	onChangeCallbacks []func()
}

// ViperOption 配置选项
type ViperOption func(*ViperConfig)

// WithConfigName 设置配置文件名（不含扩展名）
func WithConfigName(name string) ViperOption {
	return func(c *ViperConfig) {
		c.configName = name
	}
}

// WithConfigType 设置配置文件类型（json, yaml, toml, ini等）
func WithConfigType(configType string) ViperOption {
	return func(c *ViperConfig) {
		c.configType = configType
	}
}

// WithConfigPath 添加配置文件搜索路径
func WithConfigPath(paths ...string) ViperOption {
	return func(c *ViperConfig) {
		c.configPath = append(c.configPath, paths...)
	}
}

// WithEnvPrefix 设置环境变量前缀
func WithEnvPrefix(prefix string) ViperOption {
	return func(c *ViperConfig) {
		c.envPrefix = prefix
	}
}

// NewViperConfig 创建新的 Viper 配置管理器
func NewViperConfig(opts ...ViperOption) *ViperConfig {
	vc := &ViperConfig{
		viper:      viper.New(),
		configName: "config",
		configType: "yaml",
		configPath: []string{"."},
		envPrefix:  "",
	}

	// 应用选项
	for _, opt := range opts {
		opt(vc)
	}

	// 设置配置名和类型
	vc.viper.SetConfigName(vc.configName)
	vc.viper.SetConfigType(vc.configType)

	// 添加配置搜索路径
	for _, path := range vc.configPath {
		vc.viper.AddConfigPath(path)
	}

	// 设置环境变量
	if vc.envPrefix != "" {
		vc.viper.SetEnvPrefix(vc.envPrefix)
		vc.viper.AutomaticEnv()
		vc.viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	}

	return vc
}

// Load 加载配置文件
func (vc *ViperConfig) Load() error {
	// 尝试读取配置文件
	if err := vc.viper.ReadInConfig(); err != nil {
		// 如果配置文件不存在，不一定是错误（可能只使用环境变量）
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return errs.WrapMsg(err, "读取配置文件失败")
		}
	}

	return nil
}

// LoadWithFile 从指定文件加载配置
func (vc *ViperConfig) LoadWithFile(configFile string) error {
	vc.viper.SetConfigFile(configFile)
	if err := vc.viper.ReadInConfig(); err != nil {
		return errs.WrapMsg(err, "读取配置文件失败", "file", configFile)
	}

	return nil
}

// Unmarshal 将配置解析到结构体
func (vc *ViperConfig) Unmarshal(config interface{}) error {
	if err := vc.viper.Unmarshal(config); err != nil {
		return errs.WrapMsg(err, "解析配置到结构体失败")
	}
	return nil
}

// UnmarshalKey 将指定键的配置解析到结构体
func (vc *ViperConfig) UnmarshalKey(key string, rawVal interface{}) error {
	if err := vc.viper.UnmarshalKey(key, rawVal); err != nil {
		return errs.WrapMsg(err, "解析配置键失败", "key", key)
	}
	return nil
}

// Get 获取配置值
func (vc *ViperConfig) Get(key string) interface{} {
	return vc.viper.Get(key)
}

// GetString 获取字符串配置
func (vc *ViperConfig) GetString(key string) string {
	return vc.viper.GetString(key)
}

// GetInt 获取整型配置
func (vc *ViperConfig) GetInt(key string) int {
	return vc.viper.GetInt(key)
}

// GetBool 获取布尔配置
func (vc *ViperConfig) GetBool(key string) bool {
	return vc.viper.GetBool(key)
}

// GetFloat64 获取浮点数配置
func (vc *ViperConfig) GetFloat64(key string) float64 {
	return vc.viper.GetFloat64(key)
}

// GetStringSlice 获取字符串切片配置
func (vc *ViperConfig) GetStringSlice(key string) []string {
	return vc.viper.GetStringSlice(key)
}

// GetStringMap 获取字符串映射配置
func (vc *ViperConfig) GetStringMap(key string) map[string]interface{} {
	return vc.viper.GetStringMap(key)
}

// IsSet 检查键是否设置
func (vc *ViperConfig) IsSet(key string) bool {
	return vc.viper.IsSet(key)
}

// Set 设置配置值
func (vc *ViperConfig) Set(key string, value interface{}) {
	vc.viper.Set(key, value)
}

// SetDefault 设置默认值
func (vc *ViperConfig) SetDefault(key string, value interface{}) {
	vc.viper.SetDefault(key, value)
}

// BindEnv 绑定环境变量到配置键
func (vc *ViperConfig) BindEnv(key string, envVars ...string) error {
	if err := vc.viper.BindEnv(append([]string{key}, envVars...)...); err != nil {
		return errs.WrapMsg(err, "绑定环境变量失败", "key", key)
	}
	return nil
}

// WatchConfig 监听配置文件变化
func (vc *ViperConfig) WatchConfig() {
	vc.viper.WatchConfig()
}

// OnConfigChange 注册配置变更回调
func (vc *ViperConfig) OnConfigChange(fn func()) {
	vc.onChangeCallbacks = append(vc.onChangeCallbacks, fn)
	vc.viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Printf("配置文件已更新: %s\n", e.Name)
		for _, callback := range vc.onChangeCallbacks {
			callback()
		}
	})
}

// WriteConfig 将当前配置写入文件
func (vc *ViperConfig) WriteConfig() error {
	if err := vc.viper.WriteConfig(); err != nil {
		return errs.WrapMsg(err, "写入配置文件失败")
	}
	return nil
}

// WriteConfigAs 将配置写入指定文件
func (vc *ViperConfig) WriteConfigAs(filename string) error {
	if err := vc.viper.WriteConfigAs(filename); err != nil {
		return errs.WrapMsg(err, "写入配置文件失败", "file", filename)
	}
	return nil
}

// GetConfigFile 获取使用的配置文件路径
func (vc *ViperConfig) GetConfigFile() string {
	return vc.viper.ConfigFileUsed()
}

// AllSettings 获取所有配置
func (vc *ViperConfig) AllSettings() map[string]interface{} {
	return vc.viper.AllSettings()
}

// Debug 打印配置调试信息
func (vc *ViperConfig) Debug() {
	fmt.Printf("配置文件: %s\n", vc.GetConfigFile())
	fmt.Printf("所有配置:\n")
	for k, v := range vc.AllSettings() {
		fmt.Printf("  %s: %v\n", k, v)
	}
}

// GlobalConfig 全局配置实例
var globalConfig *ViperConfig

// InitGlobalConfig 初始化全局配置
func InitGlobalConfig(opts ...ViperOption) error {
	globalConfig = NewViperConfig(opts...)
	return globalConfig.Load()
}

// GetGlobalConfig 获取全局配置实例
func GetGlobalConfig() *ViperConfig {
	if globalConfig == nil {
		panic("全局配置未初始化，请先调用 InitGlobalConfig")
	}
	return globalConfig
}

// QuickLoad 快速加载配置的辅助函数
// configFile: 配置文件路径（可选）
// config: 目标配置结构体指针
// envPrefix: 环境变量前缀（可选）
func QuickLoad(configFile string, config interface{}, envPrefix string) error {
	opts := []ViperOption{}
	
	if configFile != "" {
		// 从文件路径提取配置信息
		dir := filepath.Dir(configFile)
		base := filepath.Base(configFile)
		ext := filepath.Ext(base)
		name := strings.TrimSuffix(base, ext)
		
		if ext != "" {
			ext = strings.TrimPrefix(ext, ".")
			opts = append(opts, WithConfigType(ext))
		}
		opts = append(opts, 
			WithConfigName(name),
			WithConfigPath(dir),
		)
	} else {
		// 使用默认配置路径
		executable, _ := os.Executable()
		execDir := filepath.Dir(executable)
		opts = append(opts, WithConfigPath(
			".",
			"./config",
			execDir,
			filepath.Join(execDir, "config"),
		))
	}
	
	if envPrefix != "" {
		opts = append(opts, WithEnvPrefix(envPrefix))
	}
	
	vc := NewViperConfig(opts...)
	
	if configFile != "" {
		if err := vc.LoadWithFile(configFile); err != nil {
			return err
		}
	} else {
		if err := vc.Load(); err != nil {
			return err
		}
	}
	
	return vc.Unmarshal(config)
}

// MergeConfig 合并多个配置源
func MergeConfig(primary, secondary *ViperConfig) error {
	// 获取次要配置的所有设置
	for key, value := range secondary.AllSettings() {
		// 如果主配置中没有设置该键，则使用次要配置的值
		if !primary.IsSet(key) {
			primary.Set(key, value)
		}
	}
	return nil
}
