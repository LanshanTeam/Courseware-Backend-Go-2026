# 蓝山工作室——Golang第三节课：常用标准库

### 前言

前两节课讲完了语法、复合数据类型和面向对象。这节课讲标准库。

标准库的包比较多，按包逐个讲容易散。这节课按一条日志的处理流程来讲：一条日志从终端读入，经过解析和格式化，写入文件，再读回来统计。每个包对应流程中的一环。

代码以阅读和理解为主，建议边看边敲一遍。

标准库文档：

- 中文： https://studygolang.com/pkgdoc
- 英文： https://pkg.go.dev/std

***

## 接收日志：`fmt` 与 `bufio`

`fmt` 包实现了类似 C 语言 printf 和 scanf 的格式化 I/O(input/output)。

### 输出

**fmt.Print**：直接输出，不带格式，输出到标准输出

```go
func Print(a ...interface{}) (n int, err error)
```

```go
fmt.Print("Hello, Go!")
```

多个参数之间的空格规则：字符串和字符串之间不加空格，其他情况加空格。

```go
fmt.Print("a", "b", 1, 2, "\n") // 输出：ab1 2
```

**fmt.Println**：在输出内容的末尾加一个换行符

```go
func Println(a ...interface{}) (n int, err error)
```

```go
fmt.Println("Hello, Go!")
```

多个参数之间统一用空格分隔，是最常用的调试输出方式。

**fmt.Printf**：格式化输出，支持使用格式化占位符

```go
func Printf(format string, a ...interface{}) (n int, err error)
```

```go
fmt.Printf("我的名字是 %s，我今年 %d 岁。\n", "小明", 20)
// 输出：我的名字是小明，我今年 20 岁。
```

三个函数都返回 `(n int, err error)`，表示写出的字节数和错误。一般直接忽略返回值。

### 格式化占位符

`fmt.Printf` 函数支持多种格式化占位符，常用的包括：

| 占位符 | 含义 | 示例输出 |
| --- | --- | --- |
| `%v` | 值的默认格式 | `map[a:1]` |
| `%+v` | 默认格式，结构体带字段名 | `{Level:INFO Code:200}` |
| `%#v` | Go 语法表示 | `map[string]int{"a":1}` |
| `%T` | 值的类型 | `map[string]int` |
| `%t` | 布尔值 | `true` |
| `%d` | 十进制整数 | `200` |
| `%b` | 二进制整数 | `11001000` |
| `%o` | 八进制整数 | `310` |
| `%x` / `%X` | 十六进制整数 | `c8` |
| `%c` | 字符 | `A` |
| `%q` | 带双引号的字符串 | `"abc"` |
| `%s` | 字符串 | `abc` |
| `%f` | 浮点数 | `3.140000` |
| `%e` / `%g` | 科学计数法 / 自动选择 | |
| `%p` | 指针 | `0xc000010000` |
| `%%` | 一个百分号 | `%` |

%v 和 %T 的示例

```go
m := map[string]int{"a": 1, "b": 2}
fmt.Printf("map 的值: %v\n", m)
fmt.Printf("map 的类型: %T\n", m)
// 输出：
// map 的值: map[a:1 b:2]
// map 的类型: map[string]int
```

`%v` 输出内容，`%T` 输出类型。

占位符可以带宽度和精度，常用于对齐输出：

```go
fmt.Printf("%5d|%-5d|%05d\n", 42, 42, 42) // "   42|42   |00042"
fmt.Printf("%.2f\n", 3.14159)             // "3.14"
fmt.Printf("%8.3f\n", 3.14159)            // "   3.142"
```

**fmt.Sprintf()**:根据format参数生成格式化的字符串并返回该字符串

```go
func Sprintf(format string, a ...interface{}) string
```

```go
str := fmt.Sprintf("我的名字是 %s", "小明")
fmt.Println(str) // 输出：我的名字是小明
```

**fmt.Fprintf()**:将格式化的输出写入实现了 `io.Writer` 接口的对象

```go
func Fprintf(w io.Writer, format string, a ...interface{}) (n int, err error)
```

```go
file, err := os.Create("output.txt")
if err != nil {
    fmt.Println("文件创建失败:", err)
    return
}
defer file.Close()

name, age := "小明", 20
_, err = fmt.Fprintf(file, "我的名字是 %s，今年 %d 岁。\n", name, age)
if err != nil {
    fmt.Println("写入文件失败:", err)
    return
}
```

`*os.File` 实现了 `io.Writer` 接口，所以可以传给 `Fprintf`。终端、文件、网络连接、内存缓冲区都可以，只要实现了 `Write` 方法。这个接口在讲 `io` 时展开。

同理还有Sprint，Sprintln，Fprint，Fprintln。`S` 对应 String，`F` 对应某个 Writer。

### 输入

**fmt.Scan**：从标准输入中读取数据，以空格分隔，换行也视为空格（类似于 c++ 中的 cin）

```go
func Scan(a ...interface{}) (n int, err error)
```

```go
var name string
var age int

fmt.Print("请输入你的名字和年龄：")
if _, err := fmt.Scan(&name, &age); err != nil {
    fmt.Println("输入有误:", err)
}
```

参数必须是指针，否则读到的值写不回变量，且不会报错。

**fmt.Scanln**：读取一行输入，遇到换行符停止

```go
func Scanln(a ...interface{}) (n int, err error)
```

```go
var sentence string
fmt.Print("请输入一句话：")
if _, err := fmt.Scanln(&sentence); err != nil {
    fmt.Println("输入有误:", err)
}
```

注意：`Scanln` 仍然按空格分词。输入 `用户登录 成功` 时，`sentence` 只能得到 `用户登录`，剩下的内容留在输入缓冲区里。名字里的 `ln` 表示读完这一项后要求换行，不表示读取整行文本。

**fmt.Scanf**:根据指定格式读取输入

```go
func Scanf(format string, a ...interface{}) (n int, err error)
```

```go
var name string
var age int

fmt.Print("请输入你的名字和年龄（格式：名字 年龄）：")
if _, err := fmt.Scanf("%s %d", &name, &age); err != nil {
    fmt.Println("输入有误:", err)
}
```

格式串中的空白匹配任意数量的空白字符，`"%s %d"` 与 `"%s%d"` 效果相同。

**输入系列的数据来源**

| 系列 | 数据来源 |
| --- | --- |
| `Scan` / `Scanln` / `Scanf` | 标准输入 `os.Stdin` |
| `Sscan` / `Sscanln` / `Sscanf` | 一个字符串 |
| `Fscan` / `Fscanln` / `Fscanf` | 任意 `io.Reader` |

**fmt.Sscanf**:从字符串中读取格式化数据

```go
func Sscanf(str string, format string, a ...interface{}) (n int, err error)
```

```go
str := "Alice 30"

var name string
var age int

n, err := fmt.Sscanf(str, "%s %d", &name, &age)
if err != nil {
    fmt.Println("解析错误:", err)
    return
}
fmt.Println(n, name, age) // 2 Alice 30
```

**fmt.Fscanf**:从指定的 `io.Reader` 中根据指定格式读取输入

```go
func Fscanf(r io.Reader, format string, a ...interface{}) (n int, err error)
```

```go
file, err := os.Open("data.txt")
if err != nil {
    fmt.Println("无法打开文件:", err)
    return
}
defer file.Close()

var name string
var age int

for {
    _, err := fmt.Fscanf(file, "%s %d\n", &name, &age)
    if err != nil {
        break // 读取完毕或出现错误
    }
    fmt.Printf("姓名: %s, 年龄: %d\n", name, age)
}
```

上面的循环不区分文件结束和真正的错误。更严谨的写法是判断 `err == io.EOF`，读文件统计时再展开。

### fmt.Errorf 与错误链

**fmt.Errorf**:根据format参数生成格式化字符串并返回一个包含该字符串的错误

```go
func Errorf(format string, a ...interface{}) error
```

```go
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, fmt.Errorf("不能将 %f 除以 0", a)
    }
    return a / b, nil
}
```

跟 `errors.New()` 功能类似，区别是可以格式化。

`fmt.Errorf` 的 `%w` 动作用来包装原错误：

```go
if err != nil {
    return fmt.Errorf("读取输入失败: %w", err)
}
```

用 `%w` 时，新错误把原错误包在里面，调用方可以用 `errors.Is` 或 `errors.As` 取出原错误判断，同时错误信息保留上下文。日志程序每个环节都建议这样包一层，否则只能看到一句 `EOF`。

### bufio.Reader：读一整行

`fmt` 的输入函数按空格分词，读不到含空格的整行文本。读整行要用 `bufio.Reader`。

`bufio` 包在内存中加一层缓冲区，数据先读进缓冲区，程序再从缓冲区里取，这样能减少访问底层数据流的次数。`os.Stdin` 就是被包装的底层读取器（`os` 包和 `io.Reader` 接口后面再讲）。

**NewReader**：创建带缓冲的读取器

```go
func NewReader(rd io.Reader) *Reader
```

**ReadString**：读取直到遇到分隔符，返回值包含分隔符

```go
func (b *Reader) ReadString(delim byte) (string, error)
```

```go
reader := bufio.NewReader(os.Stdin)
line, err := reader.ReadString('\n')
```

读到的 `line` 结尾带 `\n`（Windows 上是 `\r\n`），需要自己清理。

`bufio.Reader` 的其余方法，以及 `bufio.Writer` 和缓冲的完整说明，在后面讲 `os` 写文件时展开。这里只需要会用 `NewReader` 和 `ReadString`。

### 动手：读入一条日志

日志格式约定为三段：

```
级别,内容,状态码
INFO,用户登录,200
```

日志内容可能包含空格，`Scanln` 只能读到第一个空格之前，因此要用 `bufio.Reader` 读到换行符。

```go
// 片段：从终端读入一行日志，去掉首尾空白
func readLogLine() (string, error) {
    fmt.Print("请输入日志（格式：级别,内容,状态码）：")

    reader := bufio.NewReader(os.Stdin)
    line, err := reader.ReadString('\n')
    if err != nil {
        return "", fmt.Errorf("读取输入失败: %w", err)
    }

    return strings.TrimSpace(line), nil
}
```

- `reader.ReadString('\n')` 的返回值包含换行符，交给 `strings.TrimSpace` 去掉（`strings` 后面细讲）。
- Windows 上回车得到 `\r\n`，不清理的话最后一段会带 `\r`，`Atoi` 会报错。

### 注意

1. `fmt.Scan` 系列不要和 `bufio.NewReader(os.Stdin)` 混用。`bufio` 会预先读走一大块数据，之后 `Scan` 读不到内容。
2. `Scanln` 读不了带空格的句子。
3. `Scan` 系列参数要传指针。
4. `ReadString` 的返回值包含分隔符，需要自己处理。
5. 用 `%s` 打印 `int` 会输出 `%!s(int=200)`，看到这种输出说明占位符用错了。

***

## 解析日志：`strings` 与 `strconv`

一行 `INFO,用户登录,200`，需要拆成三个字段并转换类型。

### strings

`strings` 包提供字符串操作函数。函数较多，下面列出常用的。

**Contains / HasPrefix / HasSuffix**：判断包含、前缀、后缀

```go
func Contains(s, substr string) bool
func HasPrefix(s, prefix string) bool
func HasSuffix(s, suffix string) bool
```

```go
fmt.Println(strings.Contains("user login", "login")) // true
fmt.Println(strings.HasPrefix("INFO", "IN"))         // true
fmt.Println(strings.HasSuffix("app.log", ".log"))    // true
```

**TrimSpace / Trim / TrimPrefix / TrimSuffix**：修剪字符串

```go
func TrimSpace(s string) string           // 去首尾空白（空格、\t、\r、\n）
func Trim(s, cutset string) string        // 去首尾出现的任意一个字符
func TrimPrefix(s, prefix string) string  // 去前缀
func TrimSuffix(s, suffix string) string  // 去后缀
```

`" INFO "` 和 `"INFO"` 在 `switch` 中是不同的值，解析前统一 `TrimSpace`。

**Split / SplitN**:将字符串按指定分隔符分割成切片，`SplitN` 限制分割的次数

```go
func Split(s, sep string) []string
func SplitN(s, sep string, n int) []string
```

```go
fmt.Println(strings.Split("a,b,c", ","))     // [a b c]
fmt.Println(strings.SplitN("a,b,c", ",", 2)) // [a b,c]
```

`n` 表示最多切几段。`n=2` 只切第一刀，剩余部分作为最后一段。

`Split(s, "")` 按 UTF-8 字符切分。

**Index / LastIndex**:查找子串在字符串中首次或最后一次出现的位置，找不到则返回 `-1`

```go
func Index(s, substr string) int
func LastIndex(s, substr string) int
```

```go
fmt.Println(strings.Index("app.log", "."))   // 3
fmt.Println(strings.LastIndex("a.b.c", ".")) // 3
```

类似的还有 `IndexByte`、`IndexAny`、`IndexRune`。只判断存在与否用 `Contains`，需要位置才用 `Index`。

**Replace**:返回将s中前n个不重叠old子串都替换为new的新字符串，如果n<0会替换所有old子串

```go
func Replace(s, old, new string, n int) string
func ReplaceAll(s, old, new string) string
```

```go
fmt.Println(strings.Replace("a,b,c", ",", "|", 1)) // a|b,c
fmt.Println(strings.ReplaceAll("a,b,c", ",", "|")) // a|b|c
```

`Replace(s, old, new, -1)` 与 `ReplaceAll` 等价。

**Fields**：按空白切分，并丢弃空字段

```go
fmt.Println(strings.Fields("  a  b   c ")) // [a b c]
```

`strings.Split` 同样的输入会得到多个空字符串。

**ToUpper / ToLower**:将字符串转换为大写或小写

```go
func ToLower(s string) string
func ToUpper(s string) string
```

```go
fmt.Println(strings.ToUpper("Hello")) // HELLO
fmt.Println(strings.ToLower("WORLD")) // world
```

日志级别统一转大写存储，`info`、`Info`、`INFO` 都能识别。

### strconv

`strconv` 包用于在字符串和基本数据类型之间进行转换。

#### 字符串转整数/浮点数

**Atoi**:将字符串转换为 int 类型(Atoi是ParseInt(s, 10, 0)的简写)

```go
func Atoi(s string) (int, error)
```

```go
n, err := strconv.Atoi("200")
if err != nil {
    fmt.Println("转换失败:", err) // "abc" 会走到这里
}
fmt.Println(n) // 200
```

Atoi 有两个返回值，第二个必须检查。

**ParseInt**：字符串转任意进制整数

**ParseUint**：字符串转无符号整数

```go
func ParseInt(s string, base int, bitSize int) (i int64, err error)
func ParseUint(s string, base int, bitSize int) (uint64, error)
```

```go
num, err := strconv.ParseInt("1010", 2, 64) // 解析二进制 "1010" 为 64 位整数
if err != nil {
    fmt.Println("Error:", err)
} else {
    fmt.Println("Binary to int:", num) // 输出：Binary to int: 10
}
```

参数含义：

- `base`：进制，2~36；传 0 表示根据前缀自动识别（`0x` 十六进制、`0b` 二进制、`0o` 八进制）。
- `bitSize`：位宽，可取 0、8、16、32、64；传 0 表示 `int`。
- 数值超出范围时返回 `strconv.ErrRange`，同时返回被截断后的最大或最小值。

**ParseFloat**：字符串转浮点数

```go
func ParseFloat(s string, bitSize int) (f float64, err error)
```

`bitSize` 传 32 时先转 `float32` 再转回 `float64`，精度会降低。

#### 整数/浮点数转字符串

**Itoa**：int 转字符串(Itoa是FormatInt(i, 10) 的简写)

```go
func Itoa(i int) string
```

**FormatInt**：任意进制整数转字符串

**FormatUint**：FormatInt的无符号整数版本

```go
func FormatInt(i int64, base int) string
func FormatUint(i uint64, base int) string
```

```go
fmt.Println(strconv.Itoa(200))           // "200"
fmt.Println(strconv.FormatInt(255, 16))  // "ff"
```

**FormatFloat**：浮点数转字符串

```go
func FormatFloat(f float64, fmt byte, prec, bitSize int) string
```

```go
fmt.Println(strconv.FormatFloat(3.14159, 'f', 2, 64)) // "3.14"
```

四个参数依次为值、格式（`'f'` 定点、`'e'` 科学计数、`'g'` 自动）、小数位数（`-1` 表示用最少位数完整表示）、位宽。

#### 字符串扩展

```go
func AppendUint(dst []byte, i uint64, base int) []byte
```

```go
b := []byte("number:")
b = strconv.AppendInt(b, 123, 10)
fmt.Println(string(b)) // 输出：number:123
```

如 `AppendInt`，它实际上等价于 `append(dst, FormatInt(I, base)...)`，省去一次临时字符串分配。

类似的函数还有 `AppendBool`、`AppendInt`、 `AppendUint`……

### 动手：把日志拆成三个字段

先把一条日志定义成结构体：

```go
type Entry struct {
    Level   string // 日志级别：INFO / WARN / ERROR
    Content string // 日志正文
    Code    int    // 状态码
}
```

```go
// 片段：解析 "INFO,用户登录,200"
func parseEntry(line string) (Entry, error) {
    parts := strings.Split(line, ",")
    if len(parts) != 3 {
        return Entry{}, fmt.Errorf("日志格式错误：期望 3 个字段，实际 %d 个", len(parts))
    }

    level := strings.ToUpper(strings.TrimSpace(parts[0]))
    switch level {
    case "INFO", "WARN", "ERROR":
    default:
        return Entry{}, fmt.Errorf("未知的日志级别: %q", parts[0])
    }

    content := strings.TrimSpace(parts[1])
    if content == "" {
        return Entry{}, errors.New("日志内容不能为空")
    }

    code, err := strconv.Atoi(strings.TrimSpace(parts[2]))
    if err != nil {
        return Entry{}, fmt.Errorf("状态码 %q 不是合法整数: %w", parts[2], err)
    }

    return Entry{Level: level, Content: content, Code: code}, nil
}
```

这段代码用到了 `Split` 切分、`TrimSpace` 去空白、`ToUpper` 统一大小写、`Atoi` 转数字、`errors.New` 构造错误、`fmt.Errorf` 包装错误。

`errors.New` 返回一个只有错误信息的 error，不支持格式化，和 `fmt.Errorf` 是一对。

先校验 `len(parts)` 的原因：用户只输入 `INFO,用户登录`（少一段）时，`parts[2]` 会越界 panic。先判断长度，再按下标取值。

### 注意

1. 切分之后必须验证字段数量，下标越界是 panic，不是 error。
2. `Atoi` 的错误必须检查。`n, _ := strconv.Atoi("abc")` 得到 `n = 0`，会静默记下一条状态码为 0 的日志。
3. 切分前对整行 `TrimSpace`，切分后对每段再 `TrimSpace`。输入 `INFO , 用户登录 , 200` 很常见。
4. `strings.Split("", "")` 返回 `[""]`，长度为 1。
5. `strings` 的函数返回新字符串，不修改原值，因为 Go 的字符串不可变。

***

## 格式化日志：`time` 与 `fmt`

### 时间

`time` 包中的时间使用结构体 `Time` 表示。内部字段不导出，只能通过方法访问：

```go
type Time struct {
    wall uint64    // 内部表示
    ext  int64     // 内部表示
    loc  *Location // 时区
}
```

**time.Now()** 获取当前时间

```go
func Now() Time
```

```go
now := time.Now()

fmt.Println(now.Year(), now.Month(), now.Day())
fmt.Println(now.Hour(), now.Minute(), now.Second())
fmt.Println(now.Weekday())
```

常用的取值方法：

```go
func (t Time) Date() (year int, month Month, day int) // 返回 t 对应的年、月、日
func (t Time) Clock() (hour, min, sec int)            // 返回 t 对应的时、分、秒

func (t Time) Year() int
func (t Time) Month() Month
func (t Time) Day() int
func (t Time) Hour() int
func (t Time) Minute() int
func (t Time) Second() int
func (t Time) Nanosecond() int
func (t Time) Weekday() Weekday
func (t Time) YearDay() int
```

`Month` 和 `Weekday` 是自定义类型，底层为 `int`，直接打印输出英文名，取数字要 `int(now.Month())`。

### 操作

#### 格式化方法

**`Format`** 方法将 `time.Time` 转换为格式化的字符串

```go
func (t Time) Format(layout string) string
```

```go
fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
// 2026-10-21 17:52:51
```

**time.Parse**:将字符串通过格式转化成time.Time

```go
func Parse(layout, value string) (Time, error)
```

```go
t, err := time.Parse("2006-01-02 15:04:05", "2026-10-21 17:52:51")
if err != nil {
    fmt.Println("解析失败:", err)
}
fmt.Println(t.Year()) // 2026
```

tips：Go语言中格式化时间模板不是常见的 `Y-m-d H:M:S` 而是使用Go语言的诞生时间 2006 年 1 月 2 号 15 点 04 分 05 秒，也就是 `2006-01-02 15:04:05`（记不住没关系，IDE会帮你的）

layout 中各元素的对应关系：

| 想要 | 写法 |
| --- | --- |
| 年 | `2006`（`06` 为两位年份） |
| 月 | `01`（`1` 不补零） |
| 日 | `02`（`2` 不补零） |
| 时 | `15`（24 小时制，`03` 为 12 小时制） |
| 分 | `04` |
| 秒 | `05` |
| 毫秒 | `.000` |
| 星期 | `Monday` / `Mon` |
| 月份英文 | `January` / `Jan` |
| 时区 | `MST` |
| 时区偏移 | `-0700` |

layout 中的数字是占位符，不是具体数值。`2006` 不能写成 `2026`。

`Parse` 得到的是 UTC 时间（layout 中没有时区信息）。按本地时区解析用 `ParseInLocation`：

```go
loc, _ := time.LoadLocation("Asia/Shanghai")
t, err := time.ParseInLocation("2006-01-02 15:04:05", "2026-10-21 17:52:51", loc)
```

日志解析出来差 8 小时，通常是这里没有处理时区。

#### 时间操作

**Add**：时间点t1+时间段d，返回新的时间点t2

```go
func (t Time) Add(d Duration) Time
```

**Sub**：时间点t1-时间点t2，返回一个时间段

```go
func (t Time) Sub(u Time) Duration
```

```go
now := time.Now()
fmt.Println(now.Add(2 * time.Hour))       // 两小时后
fmt.Println(now.AddDate(0, 1, 0))         // 一个月后
fmt.Println(now.Sub(now.Add(-time.Hour))) // 1h0m0s
```

`Sub` 返回 `Duration`，直接打印是 `1h0m0s`，要数字用 `d.Hours()`、`d.Minutes()`、`d.Milliseconds()`。

比较：

```go
func (t Time) Before(u Time) bool // t 是否早于 u
func (t Time) After(u Time) bool  // t 是否晚于 u
func (t Time) Equal(u Time) bool  // 是否为同一时刻
```

不要用 `==` 比较两个 `Time`。`Time` 带时区和单调时钟信息，同一时刻也可能返回 `false`，用 `Equal`。

其余常用函数：

```go
func Since(t Time) Duration            // t 到现在经过的时间
func Until(t Time) Duration            // 现在到 t 的时间
func (t Time) Truncate(d Duration) Time // 截断到某一精度
func (t Time) Round(d Duration) Time    // 四舍五入到某一精度
func (t Time) IsZero() bool             // 是否零值
```

`time.Since(start)` 是测耗时的标准写法，作业里会用到。

#### 时间段

`time` 包中，用来表示时间段（纳秒为单位）的是 `Duration`

```go
type Duration int64
```

```go
const (
    Nanosecond  Duration = 1
    Microsecond          = 1000 * Nanosecond
    Millisecond          = 1000 * Microsecond
    Second               = 1000 * Millisecond
    Minute               = 60 * Second
    Hour                 = 60 * Minute
)
```

`Duration` 最大大约能表示 290 年。

#### 时间戳

**时间戳**是一种用于标记特定时间点的数字表示，通常表示自某个基准时间点（如 1970 年 1 月 1 日 00:00:00 UTC，即 Unix 纪元）以来的秒数或毫秒数。时间戳广泛应用于计算机系统中，用来记录事件的发生时间，便于排序、同步和计算时间间隔

时间转换到Unix时间戳

```go
func (t Time) Unix() int64      // 秒
func (t Time) UnixMilli() int64 // 毫秒
func (t Time) UnixMicro() int64 // 微秒
func (t Time) UnixNano() int64  // 纳秒
```

```go
now := time.Now()
fmt.Println(now.Unix())      // 1792576371
fmt.Println(now.UnixMilli()) // 1792576371123
```

Unix时间戳转换到时间

```go
func Unix(sec int64, nsec int64) Time // 秒和纳秒
func UnixMilli(msec int64) Time       // 毫秒
func UnixMicro(usec int64) Time       // 微秒
```

```go
t := time.Unix(1792576371, 0)
fmt.Println(t.Format("2006-01-02 15:04:05"))
```

`time.Unix` 返回本地时区的时间。同一时间戳在不同时区格式化出的字符串不同，但表示同一时刻。

### 时区

在 Golang 的 `time` 包中，时区用结构体 `Location` 表示

UTC 标准时间是以 GMT（Greenwich Mean Time，格林尼治时间）这个时区为主，所以本地时间与 UTC 时间的时差就是本地时间与 GMT 时间的时差。

```
UTC + 时间差 = 本地时间
```

可以通过 `FixedZone` 和 `LoadLocation` 方法来获取变量

```go
func FixedZone(name string, offset int) *Location
```

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    // 创建一个固定偏移时区 (+8小时)
    loc := time.FixedZone("CST", 8*3600)

    // 获取当前时间并设置时区
    now := time.Now().In(loc)
    fmt.Println("自定义时区 (+8):", now)
}
```

```go
func LoadLocation(name string) (*Location, error)
```

```go
loc, err := time.LoadLocation("Asia/Shanghai")
if err != nil {
    fmt.Println("加载时区失败:", err)
    return
}

// 获取当前时间并转换为指定时区
now := time.Now().In(loc)
fmt.Println("上海时间:", now)
```

**`In`**: 返回采用 `loc` 指定的时区的时间

```go
func (t Time) In(loc *Location) Time
```

`LoadLocation` 依赖系统时区数据库。编进 `scratch` 镜像后常报 `unknown time zone Asia/Shanghai`，解决办法是加 `import _ "time/tzdata"`，把时区数据编进二进制（体积增加约 450KB）。

### 定时器

在 Go 语言中，`time` 包提供了 `Timer` 和 `Ticker` 两种类型，用于处理定时和周期性的时间事件

Timer

```go
func NewTimer(d Duration) *Timer
```

`Timer` 类型表示一个单次事件。在未来的某个时间点，`Timer` 会发送一个值到它的 `C` 通道，使用 `timer` 的 `Stop` 方法可以提前停止它

```go
timer := time.NewTimer(3 * time.Second)
defer timer.Stop()

fmt.Println("等待 3 秒...")

// 通过通道接收计时结束信号
<-timer.C
fmt.Println("计时结束！")
```

`Stop()` 返回 `true` 表示成功拦下了这次触发。

Ticker

```go
func NewTicker(d Duration) *Ticker
```

#### Ticker 使用方法

- 创建 `Ticker`：使用 `time.NewTicker(d)` 创建一个周期性触发的 `Ticker`，`d` 是时间间隔。
- 接收触发信号：每隔 `d` 时间，`Ticker` 的 `C` 通道会收到当前时间。
- 常用方法：
    - **`Stop()`**：停止 `Ticker`，不再发送信号。
    - **`Reset(d)`**：重置 `Ticker` 的间隔

```go
ticker := time.NewTicker(1 * time.Second)
defer ticker.Stop()

for t := range ticker.C {
    fmt.Println("当前时间:", t)
}
```

这个 `for range` 不会自己结束，需要在循环里 `break` 或另外控制退出。

`time.Tick(d)` 一行就能拿到通道，但无法停止。Go 1.23 之后可以被垃圾回收，仍然推荐 `NewTicker` 加 `defer Stop`。

`time.After(d)` 返回只发一次值的通道，常用于超时：

```go
select {
case msg := <-ch:
    fmt.Println("收到:", msg)
case <-time.After(2 * time.Second):
    fmt.Println("超时")
}
```

`select` 是并发那节课的内容。

### 动手：拼出标准日志行

```go
const timeLayout = "2006-01-02 15:04:05"

// 给 Entry 挂一个方法，实现 fmt.Stringer 接口
func (e Entry) String() string {
    return fmt.Sprintf("[%s] [%s] %s (code=%d)",
        time.Now().Format(timeLayout), e.Level, e.Content, e.Code)
}
```

```go
e := Entry{Level: "INFO", Content: "用户登录", Code: 200}
fmt.Println(e)
// [2026-10-21 17:52:51] [INFO] 用户登录 (code=200)
```

`fmt.Println(e)` 输出这个格式，是因为 `Entry` 实现了 `String() string` 方法。`fmt` 打印任何值前都会检查是否实现了 `fmt.Stringer` 接口。

需要同时输出 Unix 时间戳时：

```go
func (e Entry) StringWithUnix() string {
    now := time.Now()
    return fmt.Sprintf("%s (%d) [%s] %s (code=%d)",
        now.Format(timeLayout), now.Unix(), e.Level, e.Content, e.Code)
}
```

`now` 只取一次。同一个 `Sprintf` 里写两次 `time.Now()` 可能跨秒，前后半句对不上。

### 注意

1. layout 写错不会报错，只会输出错乱的结果。例如 `t.Format("2026-10-21")` 得到 `21216-100-2110`，因为其中每个数字都被当成占位符解析了。
2. `Format` 和 `Parse` 的 layout 要一致。
3. `Parse` 默认返回 UTC 时间。
4. 不要用 `==` 比较 `Time`，用 `Equal`。
5. `Before` 表示"早于"，`After` 表示"晚于"，注意方向。
6. `Ticker` 要 `Stop`。
7. `time.Now()` 调用多次不是同一时刻，一次取、多次用。

***

## 写入日志文件：`io` 与 `os`

### io

`io` 库是 Go 里面 IO 相关的最核心的库，因为许多用于处理 I/O 操作的接口和函数定义在这个库里。虽然这个库不涉及具体 I/O 的实现，但是描述的是最核心的 I/O 接口语义。Go 语言的 `io` 包充分利用接口的特性，通过为不同的 I/O 操作提供一致的接口，提高了可扩展性和灵活性.

#### io.Reader

io.Reader 表示一个读取器，它将数据从某个资源读取到传输缓冲区。在缓冲区中，数据可以被流式传输和使用。

对于要用作读取器的类型，它必须实现 io.Reader 接口的唯一一个方法 Read(p []byte)。
换句话说，只要实现了 Read(p []byte) ，那它就是一个读取器。

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}
```

Read() 方法有两个返回值，一个是读取到的字节数，一个是发生错误时的错误。

约定：

1. `n > 0` 时 `err` 可能为 `nil`，也可能为 `io.EOF`。
2. 数据读取完毕时返回 `n == 0, err == io.EOF`。
3. 除 `len(p) == 0` 外，实现方不应返回 `n == 0, err == nil`。

所以读取循环要这样写：

```go
buf := make([]byte, 1024)
for {
    n, err := r.Read(buf)
    if n > 0 {
        fmt.Print(string(buf[:n]))
    }
    if err == io.EOF {
        break
    }
    if err != nil {
        return err
    }
}
```

先处理 `n > 0` 的数据，再判断 `err`。顺序反了会丢掉最后一次读到的内容。

#### io.Writer

writer则与reader相对。

```go
type Writer interface {
    Write(p []byte) (n int, err error)
}
```

约定：`n < len(p)` 时必须返回非 `nil` 的 `err`。所以写操作要同时检查 `err`，不能只看 `n`。

两者的对应关系：

```
os.Stdin    ──┐
os.File     ──┤
net.Conn    ──┼──►  io.Reader  ◄── bufio.Reader / strings.Reader / bytes.Buffer
bytes.Buffer──┘

os.Stdout   ──┐
os.File     ──┤
net.Conn    ──┼──►  io.Writer  ◄── bufio.Writer / strings.Builder / bytes.Buffer
bytes.Buffer──┘
```

配套的工具函数：

```go
func Copy(dst Writer, src Reader) (written int64, err error) // 全量搬运
func MultiWriter(writers ...Writer) Writer                   // 一份数据写多处
func TeeReader(r Reader, w Writer) Reader                    // 读的同时写一份
func LimitReader(r Reader, n int64) Reader                   // 最多读 n 字节
```

`io.Copy` 用于把文件内容拷到另一个文件，内部自己管理缓冲区。

### os

负责与操作系统交互。文件操作、环境变量管理、进程控制等

#### 文件操作

说到这里，就不得不提及一个非常常用的数据类型：`os.File`

从字面上来看，`os.File`类型代表了操作系统中的文件。但实际上，它可以代表的远不止于此。或许你已经知道，对于类Unix的操作系统（包括Linux、macOS、FreeBSD等），其中的一切都可以被看做是文件。

除了文本文件、二进制文件、压缩文件、目录这些常见的形式之外，还有符号链接、各种物理设备（包括内置或外接的面向块或者字符的设备）、命名管道，以及套接字（也就是socket），等等。

因此，可以说，我们能够利用`os.File`类型操纵的东西太多了。不过为了方便理解，我们在这里主要把`os.File`类型应用于常规的文件

##### 文件打开

**os.Open**：以只读模式打开文件(想要写入需要用OpenFile)

```go
func Open(name string) (*File, error)
```

它等价于 `OpenFile(name, os.O_RDONLY, 0)`。

**`os.Create`**：创建一个新文件，如果文件已存在则将其清空

```go
func Create(name string) (*File, error)
```

它等价于 `OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)`。注意 `O_TRUNC`，这是清空重写，不是继续写。日志文件不能用它。

**`os.OpenFile`**：提供更多的选项，可以自定义打开文件的模式。

```go
func OpenFile(name string, flag int, perm FileMode) (*File, error)
```

openfile是一个更具一般性的打开文件方法

perm是一个关于文件权限的参数，如果文件已存在会忽略

`flag` 可以选下列的值:

- `os.O_RDONLY`：只读
- `os.O_WRONLY`：只写
- `os.O_RDWR`：读写
- `os.O_APPEND`：追加写
- `os.O_CREATE`：如果文件不存在则创建
- `os.O_TRUNC`：打开时清空文件
- `os.O_EXCL`: 和 `os.O_CREATE` 配合使用，文件必须不存在
- `os.O_SYNC`: 每次写完同步刷到磁盘（很慢，慎用）

例子

```go
file, err := os.OpenFile("example.txt", os.O_RDWR|os.O_CREATE, 0666)
if err != nil {
    // 处理错误
}
defer file.Close()
```

权限位受系统 `umask` 影响。写 `0666` 实际可能是 `0644`，需要精确控制用 `os.Chmod`。

##### 文件读取/写入

**`File.Read`**：`os.File`类型实现了`io.Reader`接口，因此可以调用`Read`方法来读取数据。Read方法从f中读取最多len(b)字节数据并写入b，文件终止标志是读取0个字节且返回值err为io.EOF。（要读完整的文件得加个for）

```go
func (f *File) Read(b []byte) (n int, err error)
```

```go
buf := make([]byte, 1024) // 创建缓冲区
n, err := file.Read(buf)
if err != nil && err != io.EOF {
    // 处理错误
}
fmt.Println("Read bytes:", n, "Content:", string(buf[:n]))
```

**`os.ReadFile`**：直接读取整个文件并返回内容，适用于小文件。

```go
data, err := os.ReadFile("example.txt")
if err != nil {
    // 处理错误
}
fmt.Println("File content:", string(data))
```

**`File.Write`**：`os.File`类型实现了`io.Writer`接口,向文件写入字节切片数据。

```go
func (f *File) Write(b []byte) (n int, err error)
```

```go
n, err := file.Write([]byte("Hello, Go!"))
if err != nil {
    // 处理错误
}
fmt.Println("Wrote bytes:", n)
```

**`File.WriteString`**：向文件写入字符串数据。

```go
func (f *File) WriteString(s string) (ret int, err error)
```

```go
_, err := file.WriteString("Hello, Go!")
if err != nil {
    // 处理错误
}
```

**`os.WriteFile`**：直接写入数据到文件，适用于简单写入。

```go
err := os.WriteFile("output.txt", []byte("Hello, World!"), 0644)
if err != nil {
    // 处理错误
}
```

它的内部是 `O_WRONLY|O_CREATE|O_TRUNC`，会覆盖原文件，日志场景不能用。

**`os.Stat`**：获取文件信息，返回`os.FileInfo`接口，包含文件的基本信息，如大小、模式、修改时间等。

```go
func Stat(name string) (FileInfo, error)
```

```go
info, err := os.Stat("example.txt")
if err != nil {
    // 处理错误
}
fmt.Println("Name:", info.Name())
fmt.Println("Size:", info.Size())
fmt.Println("Mode:", info.Mode())
fmt.Println("ModTime:", info.ModTime())
fmt.Println("IsDir:", info.IsDir())
```

判断文件不存在用 `os.IsNotExist(err)`。

##### 文件目录操作

**Remove**：删除name指定的文件或目录。

```go
func Remove(name string) error
```

**Mkdir**：使用指定的权限和名称创建一个目录

**MkdirAll**：使用指定的权限和名称创建一个创建多层级目录

```go
func Mkdir(name string, perm FileMode) error
func MkdirAll(path string, perm FileMode) error
```

`Mkdir` 只能创建一级，父目录必须已存在。

**Rename**：修改一个文件的名字，移动一个文件。可能会有一些个操作系统特定的限制。

```go
func Rename(oldpath, newpath string) error
```

**RemoveAll**：递归删除，慎用。

**os.ReadDir**：列出目录内容

```go
entries, err := os.ReadDir(".")
for _, e := range entries {
    fmt.Println(e.Name(), e.IsDir())
}
```

### bufio

`bufio` 包主要用于优化 I/O 操作，通过使用缓冲区在内存中临时存储数据，从而减少对文件或网络连接等底层数据流的频繁操作。这在处理大文件、网络连接或需要逐行读取的文本文件时非常有用。

直接调用 `file.Write` 每次都是一次系统调用。写 1 万行日志就切 1 万次内核态。

`bufio` 的做法是先在内存里攒着，攒够了一起交出去。

#### Writer

**NewWriter**：创建Writer

**NewWriterSize**：创建自定义大小的缓冲

```go
func NewWriter(w io.Writer) *Writer
```

```go
file, err := os.Create("output.txt")
if err != nil {
    panic(err)
}
defer file.Close()

writer := bufio.NewWriter(file)
```

默认缓冲区 4096 字节。

**Writer**：Write将p的内容写入缓冲。返回写入的字节数。如果返回值nn < len(p)，还会返回一个错误说明原因。

```go
func (b *Writer) Write(p []byte) (nn int, err error)
```

**WriteString**：写入一个字符串。返回写入的字节数。如果返回值nn < len(s)，还会返回一个错误说明原因。

```go
func (b *Writer) WriteString(s string) (int, error)
```

**WriteByte**：WriteByte写入单个字节。

```go
func (b *Writer) WriteByte(c byte) error
```

**WriteRune**：WriteRune写入一个unicode码值（的utf-8编码），返回写入的字节数和可能的错误。

```go
func (b *Writer) WriteRune(r rune) (size int, err error)
```

**Flush**：Flush方法将缓冲中的数据写入下层的io.Writer接口。

```go
func (b *Writer) Flush() error
```

`Flush` 是 `bufio.Writer` 最容易遗漏的一步。忘了它，缓冲区里剩下的数据不会写入文件。

#### Reader

前面读入日志时已经用过 `bufio.NewReader` 和 `ReadString` 读整行，这里补充它的内部结构和其余方法。

`bufio.Reader` 是一个带有缓冲区的 `io.Reader` 接口的实现，它会在内存中存储从底层 `io.Reader` 中读取到的数据，然后先从内存缓冲区中读取数据，这样可以减少访问底层 `io.Reader` 对象的次数以及减轻操作系统的压力。结构体定义：

```go
type Reader struct {
    buf          []byte
    rd           io.Reader // reader provided by the client
    r, w         int       // buf read and write positions
    err          error
    lastByte     int // last byte read for UnreadByte; -1 means invalid
    lastRuneSize int // size of last rune read for UnreadRune; -1 means invalid
}
```

其中`r`，`w` 是两个偏移量表示缓冲区中读写的位置。当从缓冲区中读取数据时，`r` 增加，当调用底层 `io.Reader` 的 `Read` 方法读取数据到缓冲区时，`w` 增加

**NewReader**：创建Reader

**NewReaderSize**：创建自定义大小的缓冲

```go
func NewReader(rd io.Reader) *Reader
```

```go
file, err := os.Open("example.txt")
if err != nil {
    panic(err)
}
defer file.Close()

reader := bufio.NewReader(file)
```

**主要方法**

**Read**：读取数据写入p。本方法返回写入p的字节数。本方法一次调用最多会调用下层Reader接口一次Read方法，因此返回值n可能小于len(p)。读取到达结尾时，返回值n将为0而err将为io.EOF。

```go
func (b *Reader) Read(p []byte) (n int, err error)
```

**ReadByte**：读取并返回一个字节。如果没有可用的数据，会返回错误。

```go
func (b *Reader) ReadByte() (c byte, err error)
```

**ReadRune**：读取一个utf-8编码的unicode码值，返回该码值、其编码长度和可能的错误。如果utf-8编码非法，读取位置只移动1字节，返回U+FFFD，返回值size为1而err为nil。如果没有可用的数据，会返回错误。

```go
func (b *Reader) ReadRune() (r rune, size int, err error)
```

**ReadString**：读取直到分隔符，返回值包含分隔符。

```go
func (b *Reader) ReadString(delim byte) (string, error)
```

**ReadLine**：返回一行，返回值包含 `isPrefix`，表示这一行没读完需要继续调用。

```go
func (b *Reader) ReadLine() (line []byte, isPrefix bool, err error)
```

### 动手：追加写 `app.log`

日志文件要求：不存在就创建，存在就追加。

```go
// 片段：把一行日志追加到文件
func appendLog(path, line string) error {
    file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
    if err != nil {
        return fmt.Errorf("打开日志文件 %s 失败: %w", path, err)
    }
    defer file.Close()

    writer := bufio.NewWriter(file)
    if _, err := writer.WriteString(line + "\n"); err != nil {
        return fmt.Errorf("写入缓冲失败: %w", err)
    }
    if err := writer.Flush(); err != nil {
        return fmt.Errorf("刷新缓冲失败: %w", err)
    }
    return nil
}
```

各项的说明：

- `O_CREATE|O_WRONLY|O_APPEND`：不存在就建、只写、每次追加到末尾。
- `0644`：`rw-r--r--`。
- `defer file.Close()` 写在 `OpenFile` 成功之后，避免后面加 `return` 时漏掉。
- `Flush` 的错误单独检查，`WriteString` 成功只表示写进了缓冲区。
- 错误用 `%w` 包装，调用方能区分是打开失败还是刷新失败。

关于 `defer` 的顺序：`defer writer.Flush()` 必须写在 `defer file.Close()` 之后。`defer` 是后进先出，后声明的先执行，先 `Flush` 再 `Close` 才是正确顺序。顺序写反会报 `file already closed`。

另外 `defer writer.Flush()` 拿不到返回值，出错无法感知，所以推荐上面显式 `Flush` 并检查错误的写法。

### 环境变量

环境变量是操作系统为运行中的程序提供的一组键值对配置信息。你可以把它想象成贴在程序工位上的便签纸条，上面写着"数据库地址是多少"、"日志放哪里"等提示信息。

本质特征：

- 全局可访问：同一系统中的所有进程都能读取
- 继承性：子进程会继承父进程的环境变量
- 字符串格式：键和值都必须是字符串（Key=Value）
- 生命周期：随进程创建而加载，进程结束而销毁

#### 为什么要使用环境变量

1. 分离代码与配置。配置硬编码在代码里，换环境就要改代码：

```go
// 坏实践
db.Connect("mysql://admin:123456@prod-server/main")

// 好实践
db.Connect(os.Getenv("DATABASE_URL"))
```

2. 保护敏感信息。数据库密码、API密钥不写进代码，尤其不能提交到 Git。

```shell
export DB_PASS="super_secret_123"
./myapp  # 读取 DB_PASS，代码仓库里看不到密码
```

3. 动态配置容器和云部署。Docker/Kubernetes 中环境变量是主要配置方式。

```yaml
services:
  web:
    environment:
      - DATABASE_URL=mysql://db:3306/myapp
      - LOG_LEVEL=info
```

4. 系统级功能定制。`PATH` 决定去哪里找可执行文件，`LANG` 设置语言和编码，`HOME` 是用户主目录。

5. 跨进程共享配置。父进程设置的环境变量，子进程都能读到。

6. 调试和诊断。通过环境变量控制日志和行为，无需改代码。

7. 兼容性和标准化。`HTTP_PROXY`/`HTTPS_PROXY`、`EDITOR`、`TZ` 都有约定俗成的名字。

常用函数：

**Environ**：获取所有环境变量，并返回格式为 `key=value` 的字符串切片

```go
func Environ() []string
```

**Getenv**：获取指定的环境变量

```go
func Getenv(key string) string
```

```go
path := os.Getenv("LOG_FILE")
if path == "" {
    path = "app.log"
}
```

`Getenv` 无法区分"变量不存在"和"变量为空串"，需要区分时用 `LookupEnv`。

**LookupEnv**：读取环境变量并返回是否存在

```go
func LookupEnv(key string) (string, bool)
```

**Setenv**：设置环境变量

```go
func Setenv(key, value string) error
```

**Unsetenv**：删除环境变量

```go
func Unsetenv(key string) error
```

`Setenv` 只影响当前进程及其子进程，程序退出后失效。

### 命令行参数

假设我们想创建一个程序，它接受一个名字作为命令行参数，并输出一条欢迎消息。

```go
package main

import (
    "fmt"
    "os"
)

func main() {
    var args []string = os.Args
    if len(args) < 2 {
        fmt.Println("你没有提供名字哦")
        return
    }
    name := args[1] // 获取第一个命令行参数
    fmt.Printf("你好, %s!\n", name)
}
```

如果你在Go文件所在的目录中运行这个程序，可以使用以下命令，并传入一个名字参数：

```shell
go run main.go Alice
```

这样程序会输出：你好，Alice

参数较多时用 `flag` 包：

```go
var file = flag.String("file", "app.log", "日志文件路径")
var stat = flag.Bool("stat", false, "只统计，不写入")

flag.Parse()

fmt.Println(*file, *stat)
```

```shell
go run . -file logs/app.log -stat
```

`flag.String` 返回指针，使用时解引用 `*file`。`-h` 会生成帮助信息。

`flag` 只支持 `-flag value` 和 `-flag=value` 两种形式，不支持 `-f` 短选项。

### 注意

1. `bufio.Writer` 忘了 `Flush`，数据留在缓冲区，文件是空的。
2. `os.Create` 会清空已有文件。
3. `defer file.Close()` 要写在 `err` 检查之后，`OpenFile` 失败时 `file` 是 `nil`。
4. `defer` 后进先出，`Flush` 要在 `Close` 之后声明。
5. `0644` 是八进制，写成 `644` 是十进制，权限不对。
6. `os.ReadFile` 会把整个文件读进内存，只适合小文件。
7. `O_APPEND` 在多进程下仍可能交叉，单次 `Write` 小于 `PIPE_BUF` 通常原子，但不要依赖它做并发控制。

***

## 读取日志并统计：`bufio.Scanner` 与 `strings`

### Scanner

`bufio.Scanner` 用于按行或按词读取。

```go
func NewScanner(r io.Reader) *Scanner
```

**Newscanner**：创建scanner

```go
scanner := bufio.NewScanner(file)
```

**Scan**：获取当前位置的token（该token可以通过Bytes或Text方法获得），并让Scanner的扫描位置移动到下一个token。当扫描因为抵达输入流结尾或者遇到错误而停止时，本方法会返回false。在Scan方法返回false后，Err方法将返回扫描时遇到的任何错误；除非是io.EOF，此时Err会返回nil。

```go
func (s *Scanner) Scan() bool
```

**Bytes**：返回最近一次Scan调用生成的token。底层数组指向的数据可能会被下一次Scan的调用重写。

```go
func (s *Scanner) Bytes() []byte
```

**Text**：返回最近一次Scan调用生成的token，会申请创建一个字符串保存token并返回该字符串。

```go
func (s *Scanner) Text() string
```

**Err** 返回Scanner遇到的第一个非EOF的错误

```go
func (s *Scanner) Err() error
```

典型用法：

```go
scanner := bufio.NewScanner(file)
for scanner.Scan() {
    fmt.Println(scanner.Text())
}
if err := scanner.Err(); err != nil {
    fmt.Println("扫描出错:", err)
}
```

`Scan()` 返回 `false` 有两种情况：正常读完，或者出错。必须再用 `Err()` 判断，返回 `nil` 表示正常结束。

扫描粒度可以切换，默认按行：

```go
scanner.Split(bufio.ScanLines) // 按行（默认）
scanner.Split(bufio.ScanWords) // 按空白分词
scanner.Split(bufio.ScanBytes) // 按字节
scanner.Split(bufio.ScanRunes) // 按 UTF-8 字符
```

### 动手：统计各级别条数

目标：统计日志文件中 `INFO` 和 `ERROR` 的数量。

```go
// 片段：统计日志文件中各级别的条数
func countLevels(path string) (map[string]int, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, fmt.Errorf("打开日志文件失败: %w", err)
    }
    defer file.Close()

    counts := map[string]int{}
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        for _, level := range []string{"INFO", "WARN", "ERROR"} {
            if strings.Contains(line, "["+level+"]") {
                counts[level]++
            }
        }
    }
    if err := scanner.Err(); err != nil {
        return nil, fmt.Errorf("扫描日志文件失败: %w", err)
    }

    return counts, nil
}
```

说明：

- 匹配 `"[INFO]"` 而不是 `"INFO"`，避免正文中出现这个单词时被误计数。
- `counts` 的零值就是可用的空 map，不需要 `make`。
- 返回 `map` 而不是直接打印，调用方决定怎么展示。

### 注意

1. `Scan()` 返回 `false` 不一定是读完，要看 `scanner.Err()`。
2. `Scanner` 单行默认上限 64KB，超长行会返回 `bufio.ErrTooLong`。需要放大：

```go
scanner := bufio.NewScanner(file)
scanner.Buffer(make([]byte, 64*1024), 1024*1024) // 最大 1MB
```

3. `Bytes()` 返回的切片会被下一次 `Scan` 覆盖，需要保留时拷贝：`append([]byte(nil), scanner.Bytes()...)`。
4. 同一个文件不要同时开两个 `Scanner`，会互相抢偏移量。
5. 用 `Contains` 统计前先确定匹配规则。

***

## 在内存里组装文本：`strings.Builder` 与 `bytes.Buffer`

`bytes` 和 `strings` 包分别提供了对字节切片 (`[]byte`) 和字符串 (`string`) 的高效操作函数。

### 字符串拼接的开销

```go
s := ""
for i := 0; i < 10000; i++ {
    s += "line\n"
}
```

这段代码是 O(n²)。Go 的字符串不可变，每次 `+=` 都要新分配内存、拷贝旧内容、拼上新内容。

### strings.Builder

`strings.Builder` 是 Go 语言中用于高效构建字符串的结构体，特别适合频繁的字符串拼接操作

```go
func (b *Builder) Write(p []byte) (int, error)     // 将字节切片追加到缓冲区
func (b *Builder) WriteString(s string) (int, error) // 将字符串追加到缓冲区
func (b *Builder) WriteByte(c byte) error          // 追加一个字节
func (b *Builder) WriteRune(r rune) (int, error)   // 追加一个字符
func (b *Builder) Len() int                        // 返回缓冲区中的字节数
func (b *Builder) Reset()                          // 将缓冲区长度重置为 0
func (b *Builder) String() string                  // 返回缓冲区的内容作为字符串
func (b *Builder) Grow(n int)                      // 预分配容量
```

示例

```go
var b strings.Builder
b.WriteString("[INFO] 用户登录\n")
b.WriteString("[ERROR] 数据库连接失败\n")
fmt.Print(b.String())
```

两点说明：

- `WriteString` 是追加，不是覆盖。
- `String()` 之后不要再写。它内部用 `unsafe` 把字节切片零拷贝转成字符串。

### Bytes

**Equal / EqualFold**：`Equal` 用于检查两个字节切片内容是否完全相等。`EqualFold` 用于不区分大小写的比较。（byte切片不能用==直接比较）

```go
func Equal(a, b []byte) bool
```

bytes包的很多方法跟strings非常像，这里不一一列举了

同时跟 `strings.Reader` 差不多 `bytes.Reader` 也实现了读相关的接口

**Buffer**：是一个实现了读写方法的可变大小的字节缓冲。本类型的零值是一个空的可用于读写的缓冲。

**NewBuffer**：使用buf作为初始内容创建并初始化一个Buffer

大多数情况下，new(Buffer)（或只是声明一个Buffer类型变量）就足以初始化一个Buffer了。

```go
func NewBuffer(buf []byte) *Buffer
```

常用的方法

```go
func (b *Buffer) Write(p []byte) (n int, err error)         // 将 p 写入缓冲区
func (b *Buffer) WriteString(s string) (n int, err error)   // 将 s 写入缓冲区
func (b *Buffer) Read(p []byte) (n int, err error)          // 从缓冲区读取数据到 p
```

常用的操作

```go
func (b *Buffer) Reset()        // 清空缓冲区的内容
func (b *Buffer) Len() int      // 返回缓冲区中未读取的字节数
func (b *Buffer) Cap() int      // 返回缓冲区的容量
func (b *Buffer) Bytes() []byte // 以字节切片返回所有未读取的数据
func (b *Buffer) String() string // 以字符串返回所有未读取的数据
```

示例

```go
var buffer bytes.Buffer

buffer.WriteString("[INFO] 用户登录\n")
buffer.WriteString("[ERROR] 数据库连接失败\n")

fmt.Print(buffer.String())
```

`*bytes.Buffer` 同时实现了 `io.Reader` 和 `io.Writer`，可以一边写一边读。

### 怎么选

| | `strings.Builder` | `bytes.Buffer` |
| --- | --- | --- |
| 用途 | 只拼字符串 | 既读又写 |
| 实现接口 | `io.Writer`、`io.ByteWriter`、`io.StringWriter` | `io.Reader`、`io.Writer`、`io.ByteScanner` |
| 能否读 | 不能 | 能 |
| 性能 | 略快 | 略慢 |

只往一个方向拼文本用 `strings.Builder`；需要读回来用 `bytes.Buffer`。

### 动手：拼统计报告

```go
// 片段：把统计结果拼成一段报告
func buildReport(counts map[string]int) string {
    var b strings.Builder
    b.WriteString("==== 日志统计 ====\n")

    // map 的遍历顺序是随机的，固定顺序输出，保证结果稳定
    for _, level := range []string{"INFO", "WARN", "ERROR"} {
        fmt.Fprintf(&b, "%-6s %d\n", level, counts[level])
    }

    return b.String()
}
```

`fmt.Fprintf(&b, ...)` 能用的原因是 `*strings.Builder` 实现了 `io.Writer`。`fmt` 只认 `io.Writer` 接口，不认具体类型。

`%-6s` 表示左对齐、宽度 6。

### 注意

1. 循环里 `s += x` 数据量大时性能明显下降。
2. `Builder` 和 `Buffer` 不要值拷贝，传 `*strings.Builder`、`*bytes.Buffer`。
3. `Buffer.Bytes()` 返回内部切片，修改它会影响缓冲区内容。
4. `Builder.String()` 之后继续 `Write` 属于未定义行为。
5. `Reset()` 不释放内存，只把长度归零，内存占用会停在峰值。

***

## 全流程回顾

前面各步的片段按顺序拼起来就是一个完整的日志工具。这里不给完整实现，各步骤对应的函数如下：

| 步骤 | 函数 |
| --- | --- |
| 读入一行 | `readLogLine` |
| 解析校验 | `parseEntry` |
| 拼日志行 | `Entry.String` |
| 追加写入 | `appendLog` |
| 逐行统计 | `countLevels` |
| 拼统计报告 | `buildReport` |

数据在各步骤之间的变化：

```
终端输入
   │  bufio.NewReader + ReadString('\n')
   ▼
"INFO,用户登录,200"
   │  strings.Split / TrimSpace / ToUpper
   │  strconv.Atoi
   ▼
Entry{INFO, 用户登录, 200}
   │  time.Now().Format + fmt.Sprintf
   ▼
"[2026-10-21 17:52:51] [INFO] 用户登录 (code=200)"
   │  os.OpenFile + bufio.Writer + Flush
   ▼
app.log
   │  os.Open + bufio.Scanner
   ▼
map[INFO:1 WARN:0 ERROR:1]
   │  strings.Builder + fmt.Fprintf
   ▼
"==== 日志统计 ===="
```

`main` 函数按顺序调用上面几个函数即可，参数用 `-file` 和 `-stat` 两个开关处理。

这个工具目前不支持日志级别过滤、不支持多行日志、并发写也不安全，作业里会补上。

***

## 实用第三方库

### Viper

#### 简介

Viper 是 Go 语言流行的配置解决方案，支持 JSON、YAML、TOML 等多种格式，提供配置文件、环境变量、命令行标志的统一管理。

#### 安装

```bash
go get github.com/spf13/viper
```

核心功能速览

**1. 设置默认值**

```go
package main

import (
    "fmt"
    "github.com/spf13/viper"
)

func main() {
    // 设置默认值
    viper.SetDefault("app.name", "myapp")
    viper.SetDefault("server.port", 8080)
    viper.SetDefault("debug", false)

    fmt.Println(viper.GetString("app.name")) // 输出: myapp
}
```

**2. 读取配置文件**

```go
func main() {
    // 设置配置文件名(无扩展名)
    viper.SetConfigName("config")
    // 设置配置文件路径
    viper.AddConfigPath(".")
    viper.AddConfigPath("./config")
    // 设置配置类型
    viper.SetConfigType("yaml")

    // 读取配置
    if err := viper.ReadInConfig(); err != nil {
        panic(fmt.Errorf("读取配置失败: %w", err))
    }

    fmt.Println(viper.ConfigFileUsed()) // 实际读到的是哪个文件
}
```

config.yaml 示例:

```yaml
app:
  name: "mywebapp"
  version: "1.0.0"

server:
  host: "0.0.0.0"
  port: 9090

database:
  host: "localhost"
  port: 3306
  username: "root"
  password: "secret"
```

**3. 自动读取环境变量**

```go
func main() {
    // 自动匹配环境变量
    viper.AutomaticEnv()

    // 设置环境变量前缀
    viper.SetEnvPrefix("MYAPP")

    // 显式绑定环境变量名
    viper.BindEnv("app.name")
    viper.BindEnv("server.port", "MYAPP_SERVER_PORT")

    // 环境变量会自动映射:
    // MYAPP_APP_NAME -> app.name
    // MYAPP_SERVER_PORT -> server.port
}
```

映射规则：key 中的 `.` 换成 `_`，全部大写，再加前缀。

`AutomaticEnv` 只在键没有其他来源时生效。`GetString` 能读到，`Unmarshal` 往往读不到。

**4. 绑定命令行标志**

```go
import "github.com/spf13/cobra"

func main() {
    // 结合 cobra 使用
    var port int
    rootCmd := &cobra.Command{
        Run: func(cmd *cobra.Command, args []string) {
            // 绑定命令行标志
            viper.BindPFlag("server.port", cmd.Flags().Lookup("port"))

            fmt.Printf("Server port: %d\n", viper.GetInt("server.port"))
        },
    }

    rootCmd.Flags().IntVarP(&port, "port", "p", 8080, "server port")
    rootCmd.Execute()
}
```

**5. 获取配置值**

```go
func main() {
    // 读取配置
    viper.ReadInConfig()

    // 按类型获取
    appName := viper.GetString("app.name")
    port := viper.GetInt("server.port")
    debug := viper.GetBool("debug")

    // 带单位获取
    timeout := viper.GetDuration("server.timeout") * time.Second

    // 获取子树
    if dbConfig := viper.Sub("database"); dbConfig != nil {
        dbHost := dbConfig.GetString("host")
    }

    // 获取所有配置
    all := viper.AllSettings()
}
```

`Sub` 在键不存在时返回 `nil`，直接调用其方法会 panic。

**6. 映射到结构体**

```go
type Config struct {
    App struct {
        Name    string `mapstructure:"name"`
        Version string `mapstructure:"version"`
    } `mapstructure:"app"`
    Server struct {
        Host string `mapstructure:"host"`
        Port int    `mapstructure:"port"`
    } `mapstructure:"server"`
}

func main() {
    viper.SetConfigName("config")
    viper.AddConfigPath(".")
    if err := viper.ReadInConfig(); err != nil {
        panic(err)
    }

    var cfg Config
    if err := viper.Unmarshal(&cfg); err != nil {
        panic(err)
    }

    fmt.Println(cfg.Server.Port)
}
```

tag 是 `mapstructure`，不是 `json`。

**7. 监听配置变化**

```go
viper.WatchConfig()

viper.OnConfigChange(func(e fsnotify.Event) {
    fmt.Printf("配置已更新: %s\n", e.Name)
    // 重新加载配置逻辑
})

select {} // 阻塞，保持程序运行
```

依赖 `github.com/fsnotify/fsnotify`。

#### 最佳实践

配置优先级: 默认值 < 配置文件 < 环境变量 < 命令行标志

- 敏感信息: 密码、密钥等务必使用环境变量
- 配置校验: 加载后验证关键字段
- 本地覆盖: 使用 config.local.yaml 忽略到 git
- 结构体映射: 复杂配置建议用结构体管理

#### 常见问题

| 问题 | 解决 |
| --- | --- |
| 无法找到配置文件 | 检查 `AddConfigPath` 路径是否正确，`SetConfigName` 不带扩展名 |
| 环境变量不生效 | 确认 `AutomaticEnv()` 已调用，变量名为大写加下划线 |
| `Unmarshal` 读不到环境变量 | `AutomaticEnv` 对 `Unmarshal` 不可靠，改用 `BindEnv` 逐个绑定 |
| 无法反序列化到结构体 | 确保 `mapstructure` tag 与配置 key 匹配 |
| `WatchConfig` 不触发 | 没有安装 `fsnotify`，或监听路径不对 |

### snowflake

Go语言中最主流的雪花算法实现是 bwmarrin/snowflake 库，这是一个轻量级、高性能的分布式唯一ID生成器。

```bash
go get github.com/bwmarrin/snowflake
```

雪花算法（Snowflake）是Twitter开发并开源的一种分布式全局唯一ID生成算法。其核心思想是将一个64位长整型（long）数字按特定结构分割，使不同节点能独立生成全局唯一的ID，且ID随时间大致有序递增。

#### 核心原理与ID结构

雪花算法生成的64位ID由四个部分组成：

```
0 (1位) | 时间戳(41位) | 工作机器ID(10位) | 序列号(12位)
```

各部分详细说明：

**符号位（1位）**

最高位，固定为0，确保生成的ID为正数

**时间戳（41位）**

存储当前时间与自定义纪元（epoch）的毫秒偏移量

可表示约69年的时间范围（2^41-1毫秒）

提供时间有序性和可排序性

**工作机器ID（10位）**

标识不同节点，通常拆分为5位数据中心ID + 5位机器ID

最多支持32个数据中心 × 32台机器 = 1024个节点

**序列号（12位）**

同一毫秒内生成的序号，从0开始递增

支持每节点每毫秒生成4096个ID

`bwmarrin/snowflake` 的默认 epoch 是 `1288834974657`（2010-11-04），可以通过 `snowflake.Epoch` 修改。

#### 工作流程

1. 获取当前毫秒时间戳
2. 与上次生成ID的时间戳比较：
   - 相同：序列号+1，若达到4095则等待下一毫秒
   - 更大：序列号重置为0
   - 更小（时钟回拨）：触发异常处理或等待

#### 示例

```go
func main() {
    // 1. 创建节点（参数为机器ID，范围0-1023）
    node, err := snowflake.NewNode(1)
    if err != nil {
        panic(err)
    }

    // 2. 生成ID
    id := node.Generate()

    // 3. 多种格式输出
    fmt.Printf("Int64  ID: %d\n", id)          // 整数形式
    fmt.Printf("String ID: %s\n", id)          // 字符串形式
    fmt.Printf("Base2  ID: %s\n", id.Base2())  // 二进制形式
    fmt.Printf("Base64 ID: %s\n", id.Base64()) // Base64编码

    // 4. 解析ID信息
    fmt.Printf("时间戳: %d\n", id.Time())  // 毫秒时间戳
    fmt.Printf("节点ID: %d\n", id.Node()) // 机器ID
    fmt.Printf("序列号: %d\n", id.Step())  // 序列号
}
```

#### 注意

- 机器ID必须全局唯一，两个节点用同一个ID会生成重复ID。生产环境一般从配置中心或 K8s StatefulSet 序号取。
- 时钟回拨要处理，服务器校时可能把时间往回拨。
- 不适合直接展示给用户，需要短ID时用 `id.Base36()` 之类的转换。

***

## 速查表

### 常用库

| 包 | 主要用途 | 代表函数 |
| --- | --- | --- |
| `fmt` | 输入、输出、格式化 | `Println`、`Scanln`、`Sprintf`、`Fprintf`、`Errorf` |
| `bufio` | 带缓冲的读写、逐行扫描 | `NewReader`、`NewWriter`、`NewScanner` |
| `strings` | 字符串处理 | `Split`、`TrimSpace`、`Contains`、`Index`、`Replace`、`Builder` |
| `strconv` | 字符串与基本类型转换 | `Atoi`、`Itoa`、`ParseInt`、`FormatFloat` |
| `time` | 时间获取、格式化、定时 | `Now`、`Format`、`Parse`、`Since`、`NewTicker` |
| `os` | 文件及操作系统操作 | `Open`、`OpenFile`、`ReadFile`、`Stat`、`Getenv` |
| `io` | 通用输入输出接口 | `Reader`、`Writer`、`Copy`、`MultiWriter` |
| `bytes` | 内存中的字节处理 | `Buffer`、`Equal`、`Contains` |
| `errors` | 创建简单错误 | `New`、`Is`、`As` |

### 按用途查

| 需求 | 用法 |
| --- | --- |
| 读一整行（含空格） | `bufio.NewReader(os.Stdin).ReadString('\n')` |
| 按分隔符切字符串 | `strings.Split` |
| 字符串转数字 | `strconv.Atoi` / `ParseInt` / `ParseFloat` |
| 数字转字符串 | `strconv.Itoa` / `FormatInt` |
| 时间转字符串 | `t.Format("2006-01-02 15:04:05")` |
| 字符串转时间 | `time.Parse(layout, s)` |
| 测耗时 | `start := time.Now()` … `time.Since(start)` |
| 追加写文件 | `os.OpenFile(p, os.O_CREATE\|os.O_WRONLY\|os.O_APPEND, 0644)` |
| 一次读小文件 | `os.ReadFile(p)` |
| 一次写小文件 | `os.WriteFile(p, data, 0644)` |
| 按行读大文件 | `bufio.NewScanner(f)` |
| 拼接字符串 | `strings.Builder` |
| 读写内存缓冲区 | `bytes.Buffer` |
| 同时写文件和终端 | `io.MultiWriter(file, os.Stdout)` |
| 读配置 | `viper` |
| 生成唯一 ID | `snowflake` |

### 处理流程

① 接收日志

`fmt`、`bufio`

② 解析与校验

`strings`、`strconv`

③ 格式化日志

`time`、`fmt.Sprintf`

④ 写入文件

`os`、`bufio.Writer`、`Flush`

⑤ 读取与统计

`os`、`bufio.Scanner`、`strings.Contains`

⑥ 内存中组装文本

`strings.Builder`、`bytes.Buffer`

***

## 作业

### LV0

将自己不熟悉的包敲一遍（无需提交）。

### LV1

借助这节课学习的 `time` 包和 `bufio` 包自己验证一下带缓冲 I/O 写是否比不带缓冲 I/O 写的效率更高（通常写入次数越多差异越明显）

要求：分别用 `file.WriteString` 和 `bufio.Writer.WriteString` 向文件写 10 万行，用 `time.Since` 记录耗时并对比。可以再测 1 千、1 万、10 万行。

TIPS:用 `time.Since`记录花费的时间；`bufio.Writer` 记得 `Flush`，否则测的是写入缓冲区的速度。

### LV2

通过这节课的学习，相信你们都对 Go 语言的 I/O 有了比较清晰的认识。那么自己动手来实现一个日志 I/O 接口吧

需要实现的功能：

1. 将日志信息保存在文件
2. 日志信息中要包含时间和Unix时间戳

参考下面的代码

```go
package main

import (
    "fmt"
    "io"
    "os"
    "time"
)

func main() {

    // 打开一个日志文件，如果文件不存在则创建，追加写入
    file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
    if err != nil {
        panic(err)
    }
    defer file.Close()

    // 创建一个带时间戳的写入器
    logWriter := newTimestampWriter(file)

    // 模拟用户操作并记录日志
    if _, err := fmt.Fprintln(logWriter, "用户登录"); err != nil {
        panic(err)
    }
    time.Sleep(2 * time.Second)
    if _, err := fmt.Fprintln(logWriter, "用户执行操作A"); err != nil {
        panic(err)
    }
    time.Sleep(1 * time.Second)
    if _, err := fmt.Fprintln(logWriter, "用户执行操作B"); err != nil {
        panic(err)
    }
}

// timestampWriter 是一个实现 io.Writer 接口的结构体，它在写入数据前添加时间和时间戳
type timestampWriter struct {
    logFile io.Writer
}

// 传入一个io.writer,file实现了io.writer接口
func newTimestampWriter(w io.Writer) *timestampWriter {
    return &timestampWriter{logFile: w}
}

func (tw *timestampWriter) Write(p []byte) (n int, err error) {
    // 添加时间戳和时间

    // 输出到文件
    return 0, nil
}
```

需要考虑：

- `Write` 的返回值是原始数据的字节数，不是加上时间戳之后的字节数。
- `p` 中有多行时，时间戳加在哪里。
- `fmt.Fprintln` 会在结尾加 `\n`。

### LV2 PLUS

将日志信息同时输出文件和控制台中

TIPS: `io.MultiWriter`

实现任何你想实现的功能：

1. 通过环境变量，命令行参数（可以自己了解一下 `flag` 包），配置文件等方式指定日志文件路径
2. 提供日志级别功能
3. 允许自定义日志格式
4. 支持按天切分文件
5. ……

### LV PRO MAX

实现一个文件同步工具，可以监视目录中的文件变化，并在文件修改时自动将变化同步到另一个目录。

TIPS:

- 可以用 `os.Stat()` 获取文件信息
- 可以用 `os.ReadDir` 获取目录下所有文件的信息
- `os.FileInfo.ModTime` 可以获取文件更新的时间，也可以通过别的方法判断文件是否被修改
- `time.Ticker` 实现周期事件，或者自己了解一下 `cron`

### LV ？（选做）

想提 PR 但是无从下手？来挑挑蓝山后端 Go 课件的毛病吧！可以是错别字、格式、代码错误……

欢迎各位提出意见

作业提交至3440131429@qq.com

***

## 初识开源

### 什么是开源

开源项目是指通过公开可用的源代码，允许任何人查看、使用、修改和分享的项目。一些比较出名的开源项目有：

- Linux 内核
- Apache HTTP Server
- Python
- Docker
- K8S
- ……

开源项目有很多种组织模式，比如个人项目、开源社区……这里主要介绍开源社区，也是许多开源项目常见的组织形式。

开源社区是一个由志愿者组成的社群，致力于共同开发和分享开源软件、硬件或其他项目的社会集体。在开源社区中，人们可以自由地查看、使用、修改和分享项目的源代码，使用分散和协作的方式进行开源软件的开发。

### 开源项目的运作

> 开源项目的一些重要角色和元素：

- 项目维护者Maintainer：维护者有一位或者多位，他们启动开源项目、管理项目、制定决策、集思广益并且和贡献者、用户以及营销平台密切接触。
- 项目贡献者Contributor：顾名思义对参与项目贡献的人，项目越大贡献者越多，具体的贡献方式下面会介绍到
- 用户User：开源项目的消费者、使用者
- 源代码和文档库Source Code & Documentation Library：维护者将项目源代码保存在一个中心化的源代码库（如 Github）。这样所有贡献者就可以访问代码并且参与贡献。开源项目的文档库包括安装指南、使用手册、贡献指南等，有助于新用户和开发者更好地理解和参与开源项目。
- 项目许可证License：每一个开源项目都必须指定一个分发许可证，主要作用是规定其他人使用、修改或分享你的开源项目的行为。一些常见的开源许可证包括MIT许可证、GNU通用公共许可证（GPL）、Apache许可证、BSD许可证等。
- 贡献指南：帮助贡献者理解合并请求（pull request）的流程、标准、范围等。

> 什么是开源项目贡献

开源贡献就是使用一切办法提升开源项目，很多以为开源贡献只能贡献源代码（开始我也这么觉得）但其实还有很多其他方式：

- 优化项目文档，帮助贡献者和用户更容易参与和使用
- 测试应用，发现问题，提 issue
- 编写单元测试、端到端测试，使得应用的质量更高
- 写文章或制作视频等，为项目布道
- ……

> 开源项目的开发流程

一个具有一定规模的开源项目一般有一套成熟的开发流程，主要包括以下步骤：

1. 新想法的提出，可以是一个新项目、新特性或是对现有项目的提升
2. 基于这个被提议的功能，社区里志愿的开发人员完成相关开发并创建一个对现有项目的合并请求（Pull Request）
3. 社区中更有经验的人（通常是 Maintainer 或 Reviewer）将审查被提交的代码更改并决定是否应将其合并到代码的主要版本中，也可能要求作出一定的修改
4. 将当前版本作为一个开发版本来发布，即使它可能包含已知和未知的 bugs。然后，该软件将由社区进行测试，并提供反馈，项目成员和维护人员会记录并考虑社区反馈并改进实施，然后将提供新的开发版本
5. 反馈的循环会发生多次，直到项目成员认为软件的执行足够稳定，最后作为稳定版本发布

### PR

一般来说，直接向主仓库（通常是上游仓库，即项目的原始仓库）推送修改是有限制的。这是因为直接推送可能会导致冲突、代码覆盖和其他问题。开源项目一般是通过 PR(Pull Request) 允许其他开发者为项目做出贡献同时确定代码质量的。

那么如何提一个 PR 呢？我们需要遵循以下步骤：

1. Fork 项目：访问你想要为其做出贡献的项目的 GitHub 页面，点击右上角的"Fork"按钮，将该项目复制到你的GitHub账户下，这会在你的账户中创建一个对应的 Fork 仓库
2. 克隆仓库：将上面说的 Fork 仓库克隆到本地
3. 创建分支：在你的本地仓库中，创建一个新的分支
4. 修改：在这个新的分支上进行修改和提交
5. 推送到 Fork 仓库：将这个本地分支推送到你账户上的 Fork 仓库中
6. 创建 PR：现在在你的仓库中出现了一个 `Compare & pull request` 按钮，单击创建一个 PR
7. 根据社区规范填写 PR 表单：提交 PR 时，你会看到一个表单，填写标题和描述，解释你的修改是什么，为什么需要这个修改等。
8. 等待审核：提交了 PR 之后，项目维护者会审查你的代码。可能会有讨论、修改的请求，直到最终合并或关闭。如果在审查中有修改，需要在本地更新分支，然后再次推送

### 怎么找到自己感兴趣的开源项目？

#### Github Explore

GitHub Explore 会向你展示你感兴趣的仓库，你可以给特定仓库设置提醒。

https://github.com/explore/

#### Github Trending

GitHub 提供了一个 "Trending" 页面，用于显示当前在 GitHub 上流行的开源项目。这个页面上列出了每日、每周和每月的趋势，涵盖了各种编程语言和主题。

https://github.com/trending

#### Github 直接搜索

- `in:name <关键字>` 仓库名称带关键字查询
- `in:description <关键字>` 仓库描述带关键字查询
- `in:readme <关键字>` README 文件带关键字查询
- `stars(fork): >(=) <数字> <关键字>` star 或 fork 数大于(或等于)指定数字的带关键字查询
- `stars(fork): 10..20 <关键词>` star 或 fork 数在 10 到 20 之间的带关键字查询
- `size:>=5000 <关键词>` 限定仓库大于等于 5000K 的带关键字查询
- `pushed(created):>2019-11-15 <关键字>` 更新 或 创建 日期在 2019 年 11 月 16 日之后的带关键字查询
- `license:apache-2.0 <关键字>` LICENSE 为 apache-2.0 的带关键字查询
- `language:java <关键词>` 仓库语言为 Java 的带关键字查询
- `user:<用户名>` 查询某个用户的项目
- `org:<组织名>` 查询某个组织的项目

这些可以混合使用，也可以先查找某一类的 awesome 仓库，然后从 awesome 库里找相关的资源，github 里有很多归纳仓库，可以先看看已有的收集，有时候会节省很多时间

可以根据下面几个标准判断项目热度：

1. Star 数量
2. Fork 数量
3. 项目代码更新日期

#### 技术论坛

- 掘金
- Stack Overflow
- Reddit
- ……

### 从开源项目中我们能学到什么

贡献开源项目对于开发者来说好处多多，但是主要的好处有：

- 抓住机会、提高技能水平
- 提高软件/应用的代码水平
- 接触想法相似的朋友、建立人脉和社区
- 理解应用开发和维护的流程
- 从 Pull Request 反馈中学习
- 学习如何管理自己的开源项目

## 拓展阅读

1. [浅析Go IO 的知识框架](https://zhuanlan.zhihu.com/p/394888082)
2. [Golang 中的 io 包详解](https://blog.csdn.net/luduoyuan/article/details/131057066)
3. [Go by Example: Reading Files](https://gobyexample.com/reading-files)
4. [Viper 官方文档](https://github.com/spf13/viper)
5. [bwmarrin/snowflake](https://github.com/bwmarrin/snowflake)
6. [github 参与开源项目开发或贡献代码流程](https://www.cnblogs.com/yuejunasia/p/15161915.html)
7. [初学者开源项目参与指南](https://www.freecodecamp.org/chinese/news/a-practical-guide-to-start-opensource-contributions/)
