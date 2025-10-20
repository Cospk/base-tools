package jsonutil

type Json interface {
    // Interface 返回底层数据
    Interface() any
    // Encode 返回序列化后的 `[]byte`
    Encode() ([]byte, error)
    // EncodePretty 返回带缩进的序列化 `[]byte`
    EncodePretty() ([]byte, error)
    // 实现 json.Marshaler 接口
    MarshalJSON() ([]byte, error)
    // Set 通过 `key` 和 `value` 修改 `Json` 的 map
    // 便于在 `Json` 对象中快速修改单个键值
    Set(key string, val any)
    // SetPath 递归检查/创建路径上的 map 键，最终写入对应的值
    SetPath(branch []string, val any)
    // Del 如果存在则删除 `key`
    Del(key string)
    // Get 返回 `key` 对应的新 `Json` 对象指针
    //
    // 便于链式调用（遍历嵌套 JSON）：
    //    js.Get("top_level").Get("dict").Get("value").Int()
    Get(key string) Json
    // GetPath 按提供的路径查找，不必多次调用 Get()
    //
    //   js.GetPath("top_level", "dict")
    GetPath(branch ...string) Json
    // CheckGet 返回新的 `Json` 对象指针和是否成功的布尔值
    //
    // 当需要判断成功与否时的链式调用：
    //    if data, ok := js.Get("top_level").CheckGet("inner"); ok {
    //        log.Println(data)
    //    }
    CheckGet(key string) (Json, bool)
    // Map 断言为 `map`
    Map() (map[string]any, error)
    // Array 断言为 `array`
    Array() ([]any, error)
    // Bool 断言为 `bool`
    Bool() (bool, error)
    // String 断言为 `string`
    String() (string, error)
    // Bytes 断言为 `[]byte`
    Bytes() ([]byte, error)
    // StringArray 断言为 `[]string`
    StringArray() ([]string, error)

    // 实现 json.Unmarshaler 接口
    UnmarshalJSON(p []byte) error
    // Float64 转为 float64
    Float64() (float64, error)
    // Int 转为 int
    Int() (int, error)
    // Int64 转为 int64
    Int64() (int64, error)
    // Uint64 转为 uint64
    Uint64() (uint64, error)
}
