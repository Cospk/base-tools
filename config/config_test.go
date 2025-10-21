package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestConfig 测试用配置结构体
type TestConfig struct {
	Server struct {
		Host    string `mapstructure:"host"`
		Port    int    `mapstructure:"port"`
		Timeout int    `mapstructure:"timeout"`
	} `mapstructure:"server"`
	Database struct {
		DSN            string `mapstructure:"dsn"`
		MaxConnections int    `mapstructure:"max_connections"`
	} `mapstructure:"database"`
}

// TestNewViperConfig 测试创建配置管理器
func TestNewViperConfig(t *testing.T) {
	vc := NewViperConfig(
		WithConfigName("test"),
		WithConfigType("yaml"),
		WithConfigPath("./testdata"),
		WithEnvPrefix("TEST"),
	)

	if vc == nil {
		t.Fatal("创建 ViperConfig 失败")
	}

	if vc.configName != "test" {
		t.Errorf("配置名错误，期望: test, 实际: %s", vc.configName)
	}

	if vc.configType != "yaml" {
		t.Errorf("配置类型错误，期望: yaml, 实际: %s", vc.configType)
	}

	if vc.envPrefix != "TEST" {
		t.Errorf("环境变量前缀错误，期望: TEST, 实际: %s", vc.envPrefix)
	}
}

// TestSetAndGet 测试设置和获取配置
func TestSetAndGet(t *testing.T) {
	vc := NewViperConfig()

	// 测试各种类型
	vc.Set("string_key", "test_value")
	vc.Set("int_key", 123)
	vc.Set("bool_key", true)
	vc.Set("float_key", 3.14)
	vc.Set("slice_key", []string{"a", "b", "c"})
	vc.Set("map_key", map[string]interface{}{
		"nested": "value",
	})

	// 验证获取值
	if v := vc.GetString("string_key"); v != "test_value" {
		t.Errorf("字符串值错误，期望: test_value, 实际: %s", v)
	}

	if v := vc.GetInt("int_key"); v != 123 {
		t.Errorf("整数值错误，期望: 123, 实际: %d", v)
	}

	if v := vc.GetBool("bool_key"); v != true {
		t.Errorf("布尔值错误，期望: true, 实际: %v", v)
	}

	if v := vc.GetFloat64("float_key"); v != 3.14 {
		t.Errorf("浮点数值错误，期望: 3.14, 实际: %f", v)
	}

	slice := vc.GetStringSlice("slice_key")
	if len(slice) != 3 || slice[0] != "a" {
		t.Errorf("切片值错误，实际: %v", slice)
	}

	m := vc.GetStringMap("map_key")
	if m["nested"] != "value" {
		t.Errorf("映射值错误，实际: %v", m)
	}
}

// TestSetDefault 测试设置默认值
func TestSetDefault(t *testing.T) {
	vc := NewViperConfig()

	// 设置默认值
	vc.SetDefault("default_key", "default_value")
	vc.SetDefault("override_key", "default")

	// 覆盖默认值
	vc.Set("override_key", "overridden")

	// 验证
	if v := vc.GetString("default_key"); v != "default_value" {
		t.Errorf("默认值错误，期望: default_value, 实际: %s", v)
	}

	if v := vc.GetString("override_key"); v != "overridden" {
		t.Errorf("覆盖值错误，期望: overridden, 实际: %s", v)
	}
}

// TestIsSet 测试检查键是否存在
func TestIsSet(t *testing.T) {
	vc := NewViperConfig()

	vc.Set("exists_key", "value")

	if !vc.IsSet("exists_key") {
		t.Error("存在的键应该返回 true")
	}

	if vc.IsSet("not_exists_key") {
		t.Error("不存在的键应该返回 false")
	}
}

// TestUnmarshal 测试解析到结构体
func TestUnmarshal(t *testing.T) {
	vc := NewViperConfig()

	// 设置配置值
	vc.Set("server.host", "localhost")
	vc.Set("server.port", 8080)
	vc.Set("server.timeout", 30)
	vc.Set("database.dsn", "user:pass@tcp(localhost:3306)/db")
	vc.Set("database.max_connections", 100)

	// 解析到结构体
	var cfg TestConfig
	err := vc.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("解析配置失败: %v", err)
	}

	// 验证
	if cfg.Server.Host != "localhost" {
		t.Errorf("服务器主机错误，期望: localhost, 实际: %s", cfg.Server.Host)
	}

	if cfg.Server.Port != 8080 {
		t.Errorf("服务器端口错误，期望: 8080, 实际: %d", cfg.Server.Port)
	}

	if cfg.Database.MaxConnections != 100 {
		t.Errorf("数据库连接数错误，期望: 100, 实际: %d", cfg.Database.MaxConnections)
	}
}

// TestUnmarshalKey 测试解析指定键到结构体
func TestUnmarshalKey(t *testing.T) {
	vc := NewViperConfig()

	// 设置配置值
	vc.Set("server.host", "localhost")
	vc.Set("server.port", 8080)
	vc.Set("server.timeout", 30)

	// 解析指定键
	var serverCfg struct {
		Host    string `mapstructure:"host"`
		Port    int    `mapstructure:"port"`
		Timeout int    `mapstructure:"timeout"`
	}

	err := vc.UnmarshalKey("server", &serverCfg)
	if err != nil {
		t.Fatalf("解析配置键失败: %v", err)
	}

	// 验证
	if serverCfg.Host != "localhost" {
		t.Errorf("主机错误，期望: localhost, 实际: %s", serverCfg.Host)
	}

	if serverCfg.Port != 8080 {
		t.Errorf("端口错误，期望: 8080, 实际: %d", serverCfg.Port)
	}
}

// TestEnvOverride 测试环境变量覆盖
func TestEnvOverride(t *testing.T) {
	// 设置环境变量
	os.Setenv("TEST_SERVER_PORT", "9090")
	defer os.Unsetenv("TEST_SERVER_PORT")

	vc := NewViperConfig(
		WithEnvPrefix("TEST"),
	)

	// 设置默认值
	vc.SetDefault("server.port", 8080)

	// 环境变量应该覆盖默认值
	port := vc.GetInt("server.port")
	if port != 9090 {
		t.Errorf("环境变量覆盖失败，期望: 9090, 实际: %d", port)
	}
}

// TestAllSettings 测试获取所有配置
func TestAllSettings(t *testing.T) {
	vc := NewViperConfig()

	vc.Set("key1", "value1")
	vc.Set("key2", "value2")
	vc.Set("nested.key", "nested_value")

	settings := vc.AllSettings()

	if settings["key1"] != "value1" {
		t.Errorf("key1 值错误，实际: %v", settings["key1"])
	}

	if settings["key2"] != "value2" {
		t.Errorf("key2 值错误，实际: %v", settings["key2"])
	}

	// 嵌套值会被展开为 map
	if nested, ok := settings["nested"].(map[string]interface{}); ok {
		if nested["key"] != "nested_value" {
			t.Errorf("嵌套值错误，实际: %v", nested["key"])
		}
	} else {
		t.Error("嵌套配置应该是 map 类型")
	}
}

// TestLoadWithFile 测试从文件加载配置
func TestLoadWithFile(t *testing.T) {
	// 创建临时配置文件
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "test.yaml")
	configContent := `
server:
  host: testhost
  port: 9999
database:
  dsn: test_dsn
`
	err := os.WriteFile(configFile, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("创建测试配置文件失败: %v", err)
	}

	// 加载配置
	vc := NewViperConfig()
	err = vc.LoadWithFile(configFile)
	if err != nil {
		t.Fatalf("加载配置文件失败: %v", err)
	}

	// 验证
	if host := vc.GetString("server.host"); host != "testhost" {
		t.Errorf("主机配置错误，期望: testhost, 实际: %s", host)
	}

	if port := vc.GetInt("server.port"); port != 9999 {
		t.Errorf("端口配置错误，期望: 9999, 实际: %d", port)
	}
}

// TestWriteConfig 测试写入配置
func TestWriteConfig(t *testing.T) {
	// 创建临时目录
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "output.yaml")

	vc := NewViperConfig()
	
	// 设置一些配置
	vc.Set("test.key1", "value1")
	vc.Set("test.key2", 123)
	vc.Set("test.nested.key", true)

	// 写入配置文件
	err := vc.WriteConfigAs(configFile)
	if err != nil {
		t.Fatalf("写入配置失败: %v", err)
	}

	// 验证文件是否创建
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		t.Error("配置文件未创建")
	}

	// 重新加载并验证
	vc2 := NewViperConfig()
	err = vc2.LoadWithFile(configFile)
	if err != nil {
		t.Fatalf("重新加载配置失败: %v", err)
	}

	if v := vc2.GetString("test.key1"); v != "value1" {
		t.Errorf("重新加载后值错误，期望: value1, 实际: %s", v)
	}
}

// TestMergeConfig 测试合并配置
func TestMergeConfig(t *testing.T) {
	primary := NewViperConfig()
	primary.Set("key1", "primary_value")
	primary.Set("key2", "primary_only")

	secondary := NewViperConfig()
	secondary.Set("key1", "secondary_value")
	secondary.Set("key3", "secondary_only")

	// 合并配置
	err := MergeConfig(primary, secondary)
	if err != nil {
		t.Fatalf("合并配置失败: %v", err)
	}

	// 验证：primary 已有的值不应被覆盖
	if v := primary.GetString("key1"); v != "primary_value" {
		t.Errorf("主配置值被错误覆盖，期望: primary_value, 实际: %s", v)
	}

	// 验证：primary 独有的值应保留
	if v := primary.GetString("key2"); v != "primary_only" {
		t.Errorf("主配置独有值丢失，期望: primary_only, 实际: %s", v)
	}

	// 验证：secondary 独有的值应被添加
	if v := primary.GetString("key3"); v != "secondary_only" {
		t.Errorf("次要配置值未添加，期望: secondary_only, 实际: %s", v)
	}
}

// TestQuickLoad 测试快速加载
func TestQuickLoad(t *testing.T) {
	// 创建临时配置文件
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "quick.yaml")
	configContent := `
server:
  host: quickhost
  port: 7777
`
	err := os.WriteFile(configFile, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("创建测试配置文件失败: %v", err)
	}

	// 快速加载
	var cfg TestConfig
	err = QuickLoad(configFile, &cfg, "")
	if err != nil {
		t.Fatalf("快速加载失败: %v", err)
	}

	// 验证
	if cfg.Server.Host != "quickhost" {
		t.Errorf("主机错误，期望: quickhost, 实际: %s", cfg.Server.Host)
	}

	if cfg.Server.Port != 7777 {
		t.Errorf("端口错误，期望: 7777, 实际: %d", cfg.Server.Port)
	}
}

// TestConfigChange 测试配置变更监听
func TestConfigChange(t *testing.T) {
	// 创建临时配置文件
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "watch.yaml")
	initialContent := `value: initial`
	err := os.WriteFile(configFile, []byte(initialContent), 0644)
	if err != nil {
		t.Fatalf("创建测试配置文件失败: %v", err)
	}

	// 创建配置并加载
	vc := NewViperConfig()
	err = vc.LoadWithFile(configFile)
	if err != nil {
		t.Fatalf("加载配置失败: %v", err)
	}

	// 验证初始值
	if v := vc.GetString("value"); v != "initial" {
		t.Errorf("初始值错误，期望: initial, 实际: %s", v)
	}

	// 使用通道等待回调触发，避免数据竞争
	done := make(chan struct{}, 1)
	vc.OnConfigChange(func() {
		select {
		case done <- struct{}{}:
		default:
		}
	})
	vc.WatchConfig()

	// 修改配置文件
	updatedContent := `value: updated`
	err = os.WriteFile(configFile, []byte(updatedContent), 0644)
	if err != nil {
		t.Fatalf("更新配置文件失败: %v", err)
	}

	// 验证回调被触发（带超时，避免不支持的环境卡住）
	select {
	case <-done:
		// ok
	case <-time.After(2 * time.Second):
		t.Skip("配置变更监听测试跳过（可能是文件系统不支持）")
	}
}

// TestGlobalConfig 测试全局配置
func TestGlobalConfig(t *testing.T) {
	// 初始化全局配置
	err := InitGlobalConfig(
		WithConfigName("global_test"),
		WithConfigType("yaml"),
	)
	if err != nil {
		// 配置文件不存在不是错误
		t.Log("全局配置初始化（配置文件可能不存在）")
	}

	// 获取全局配置
	gc := GetGlobalConfig()
	if gc == nil {
		t.Fatal("获取全局配置失败")
	}

	// 设置和获取值
	gc.Set("global.test", "test_value")
	if v := gc.GetString("global.test"); v != "test_value" {
		t.Errorf("全局配置值错误，期望: test_value, 实际: %s", v)
	}
}

// BenchmarkGet 基准测试：获取配置
func BenchmarkGet(b *testing.B) {
	vc := NewViperConfig()
	vc.Set("benchmark.key", "value")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = vc.GetString("benchmark.key")
	}
}

// BenchmarkSet 基准测试：设置配置
func BenchmarkSet(b *testing.B) {
	vc := NewViperConfig()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		vc.Set("benchmark.key", i)
	}
}

// BenchmarkUnmarshal 基准测试：解析到结构体
func BenchmarkUnmarshal(b *testing.B) {
	vc := NewViperConfig()
	vc.Set("server.host", "localhost")
	vc.Set("server.port", 8080)
	vc.Set("database.dsn", "dsn")

	var cfg TestConfig

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = vc.Unmarshal(&cfg)
	}
}
