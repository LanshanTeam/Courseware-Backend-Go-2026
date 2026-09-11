# 蓝山工作室——Golang第〇节课

### 前言
欢迎大家来到蓝山工作室的Golang编程课程！在接下来的课程中，我们将一起探索Golang这门强大且高效的编程语言。
这节课仅作为帮助大家了解Go语言最基础的知识，不作为正式课程
***

### 关于Golang语言
Go 是一个开源的编程语言，它能让构造简单、可靠且高效的软件变得容易。

Go是从2007年末由Robert Griesemer, Rob Pike, Ken Thompson主持开发，后来还加入了Ian Lance Taylor, Russ Cox等人，并最终于2009年11月开源，在2012年早些时候发布了Go 1稳定版本。现在Go的开发已经是完全开放的，并且拥有一个活跃的社区。

***

### 认识一下Go程序
学习一门语言都是从Hello world开始的，在Go中，我们也不例外。下面是一个简单的Go程序，它会在控制台输出“Hello, World!”。
```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
```

我们来逐行解释一下这段代码：
- `package main`：每个Go程序都是由包（package）组成的。`main`包是一个特殊的包，它定义了一个可执行程序的入口点。
- `import "fmt"`：这行代码导入了Go的标准库中的`fmt`包，它提供了格式化输入输出的功能。（下面的输出就用到了它）
- `func main()`：这是程序的入口函数。当你运行这个程序时，Go会从这里开始执行代码。
- `fmt.Println("Hello, World!")`：这行代码调用了`fmt`包中的`Println`函数，用于在控制台打印输出“Hello, World!”。需要打印的内容放在括号内，并用双引号括起来表示这是一个字符串。

如果大家对于上面的程序还是无法理解，比如为什么需要`package main`，为什么要导入`fmt`包，为什么要定义`main`函数等等，不用担心，接下来的课程中我们会逐步讲解这些概念。现在只需要认识一下一个完整的Go程序是由哪些部分组成的就可以了。

***

### Go语法基础

#### 数据类型
1. 布尔型 `bool`

布尔型的值只可以是常量 true 或者 false

2. 数字类型

| 类型       | 描述                                                           |
| ---------- | -------------------------------------------------------------- |
| uint8      | 无符号 8 位整型 (0 到 255)                                     |
| uint16     | 无符号 16 位整型 (0 到 65535)                                  |
| uint32     | 无符号 32 位整型 (0 到 4294967295)                             |
| uint64     | 无符号 64 位整型 (0 到 18446744073709551615)                   |
| int8       | 有符号 8 位整型 (-128 到 127)                                  |
| int16      | 有符号 16 位整型 (-32768 到 32767)                             |
| int32      | 有符号 32 位整型 (-2147483648 到 2147483647)                   |
| int64      | 有符号 64 位整型 (-9223372036854775808 到 9223372036854775807) |
| float32    | IEEE-754 32位浮点型数                                          |
| float64    | IEEE-754 64位浮点型数                                          |
| complex64  | 32 位实数和虚数                                                |
| complex128 | 64 位实数和虚数                                                |
| int        | 根据你的操作系统架构，可以是32位或64位的整数                   |
| byte       | uint8的别名                                                    |
| rune       | int32的别名                                                    |
3. 字符串类型

   字符串的底层实现是基于一个不可变的字节数组，Go 的字符串在内存中使用 UTF-8 编码表示，每个字符串由两个部分组成：一个指向字节数组的指针和字符串的长度。
4. 派生类型

   Go 的派生类型是通过结构、组合或引用等方式，从基本类型派生出来的
   - (a) 指针类型（Pointer）
   - (b) 数组类型
   - (c) 结构化类型(struct)
   - (d) Channel 类型
   - (e) 函数类型
   - (f) 切片类型
   - (g) 接口类型（interface）
   - (h) Map 类型

### 变量声明
Go语言中的变量需要声明后才能使用，同一作用域内不支持重复声明。并且Go语言的变量声明后必须使用。
```go
// 第一种方式
// var 变量名 [类型] = 表达式
// example:
// 使用类型推导
var a = "first"
var b,c = 1,2
// 指定类型
var d int8 = 1
// 严格意义上讲，上述语句是两个过程：声明+赋值
var e float64 //声明
e = 1.0 //赋值

// 第二种方式
f := 114514
```
### 常量
相对于变量，常量是恒定不变的值，多用于定义程序运行期间不会改变的那些值。 常量的声明和变量声明非常类似，只是把var换成了const，常量在定义的时候必须赋值。
```go
const s string = "constant"
const h = 500000000
const i = 3e20 / h
fmt.Println(s, h, i, math.Sin(h), math.Sin(i))
```

### 运算符

Go 的运算符和其他主流语言类似，常用的有以下几类：

1. **算术运算符**：`+`（加）、`-`（减）、`*`（乘）、`/`（除）、`%`（取余）
2. **关系运算符**：`==`（等于）、`!=`（不等于）、`<`、`>`、`<=`、`>=`
3. **逻辑运算符**：`&&`（与）、`||`（或）、`!`（非）
4. **赋值运算符**：`=`、`+=`、`-=`、`*=`、`/=`、`%=`

```go
a, b := 10, 3
fmt.Println(a + b) // 13
fmt.Println(a - b) // 7
fmt.Println(a * b) // 30
fmt.Println(a / b) // 3，整数相除结果仍是整数，小数部分被丢弃
fmt.Println(a % b) // 1，取余

// 关系运算符的结果是 bool 类型
fmt.Println(a > b)  // true
fmt.Println(a == b) // false

// 逻辑运算符
fmt.Println(a > b && a != b) // true
fmt.Println(!(a > b))        // false

// 赋值运算符
c := 5
c += 3 // 等价于 c = c + 3
fmt.Println(c) // 8
```

> 注意：整数相除 `/` 会直接丢弃小数部分，如需保留小数，至少要有一个操作数是浮点数，例如 `10.0 / 3`。

### 类型转换

Go 是强类型语言，不同类型之间不能直接运算或赋值，需要显式转换，语法为：`类型(值)`。

```go
var a int = 10
var b float64 = float64(a) // int 转 float64
fmt.Println(b) // 10

var c float64 = 3.7
var d int = int(c) // float64 转 int，直接丢弃小数部分
fmt.Println(d) // 3

// 字符串与数字互转需要借助 strconv 包（需要 import "strconv"）
s := strconv.Itoa(123) // 数字转字符串
fmt.Println(s) // 123

n, _ := strconv.Atoi("456") // 字符串转数字，_ 忽略可能出现的错误
fmt.Println(n) // 456
```

- 整数和浮点数之间可以直接用 `类型(值)` 转换，浮点转整型会**丢失小数部分**
- 字符串与数字之间需要借助标准库 `strconv` 的 `Itoa` / `Atoi` 函数
- `Atoi` 会返回两个值：转换结果和错误（error），暂时用 `_` 忽略错误即可

### 字符串基础

字符串由字符组成，底层是不可变的字节序列，支持很多基础操作。

```go
s := "hello"
fmt.Println(len(s)) // 5，获取字符串长度（字节数）

// 拼接
s2 := s + " world"
fmt.Println(s2) // hello world

// 用下标访问单个字节（取出来的是数字，即该字符的编码）
fmt.Println(s[0]) // 104

// 截取子串，[start:end] 左闭右开
fmt.Println(s[1:3]) // el
fmt.Println(s[:2])  // he，省略 start 表示从 0 开始
fmt.Println(s[2:])  // llo，省略 end 表示取到末尾

// 常用判断函数在 strings 包里（需要 import "strings"）
fmt.Println(strings.Contains(s2, "world")) // true
fmt.Println(strings.HasPrefix(s, "he"))    // true
```

> 注意：`s[0]` 取出的是字节而不是字符，中文字符占多个字节，直接用下标访问中文会得到乱码。需要逐个取出字符时，推荐使用后面循环章节的 `for range`。

### 循环
和C++不同，Golang 只有一个循环关键字`for`，接下来讲`for`的几种使用模式
#### 三段式
```go
// for 初始化;条件;循环表达式
for i:=0;i<1;i++{
   fmt.Println(i)
}
```
第一个位置是单次表达式，循环开始时会执行一次这里，一般用于初始化变量。

第二个位置是条件表达式，即循环条件，只要满足循环条件就会执行循环体。

第三个位置是末尾循环体，每次执行完一遍循环体之后会执行一次此位置中中的表达式。

执行末尾循环体后将再次进行条件判断，若条件还成立，则继续重复上述循环，当条件不成立时则跳出当下for循环。

并且，你可以选择性的留空，就是让三段式的任何位置为空。例如：
```go
for ;i<5;{
  // 循环体
}
```
#### 一段式
一段式那就是只写条件
例如：

```go
for i<5{
  // 循环体
}
```
#### 无限循环
```go
for{
  // 循环体,记得要跳出循环
}
```
#### for range

`for range` 是 Go 中非常常用的循环写法，可以遍历字符串、数组、切片、map、通道（channel）等。我们先以字符串为例：

```go
// 遍历字符串：i 是下标，c 是每个字符
for i, c := range "hello" {
	fmt.Printf("i=%d, c=%c\n", i, c)
}
// i=0, c=h
// i=1, c=e
// i=2, c=l
// i=3, c=l
// i=4, c=o

// 只需要字符、不关心下标时，用 _ 忽略
for _, c := range "go" {
	fmt.Printf("%c ", c)
}
// g o
```

- `for range` 每轮会返回两个值：**下标**和**元素值**，不需要的那个用 `_` 忽略
- Go 语言中**声明了的变量必须使用**，所以用不到的变量统一写成 `_`
- 遍历中文时，`for range` 会按字符（rune）逐个取出，不会像按下标访问那样出现乱码
- `for range` 也常用于遍历数组、切片、map 等集合类型，这些会在后续课程中详细讲解

##### break 关键词

break 放在循环体中，只要执行到 break，则会立马跳出所在最里循环（注意，是所在最里循环，若嵌套，则无法跳出更外层循环）

例如：

```go
for i:=1;i<4;i++{
  j := i
  for j<4{
    if j == 2{
      break
    }
  }
  fmt.Println("hello lanshan")
}
```

上面这个函数不需要去体会其中的意思，我只是举个例子。若执行到了break，则只会跳出条件是 j < 4 这个循环，依然会执行println打印

这是总的示例，可以回顾一下：

```go
package main

import "fmt"

func main() {
	i := 1
	for {
		fmt.Println("loop")
		break // 跳出循环
	}
	
	// 打印7、8
	for j := 7; j < 9; j++ {
		fmt.Println(j)
	}

	for n := 0; n < 5; n++ {
		if n%2 == 0 {
			continue
			// 当n模2为0时不打印，进到下一次的循环
		}
		fmt.Println(n)
	}
	// 直到i>3
	for i <= 3 {
		fmt.Println(i)
		i = i + 1
	}
  // for 循环嵌套
  for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			fmt.Printf("i = %d, j = %d\n", i, j)
		}
	}
}
```
### if

```go
if 条件表达式 {
	//当条件表达式结果为true时，执行此处代码   
}

if 条件表达式 {
    //当条件表达式结果为true时，执行此处代码  
} else {
    //当条件表达式结果为false时，执行此处代码  
}
```

```go
package main

import "fmt"

func main() {
	// 条件表达式为false，打印出"7 是奇数"
	if 7%2 == 0 {
		fmt.Println("7 是偶数")
	} else {
		fmt.Println("7 是奇数")
	}

	// 条件表达式为ture，打印出"8 被 4 整除"
	if 8%4 == 0 {
		fmt.Println("8 被 4 整除")
	}

	// 这是一个短声明，效果等效于
	//num := 9
	//if num < 0{
	//	...
	//}
	if num := 9; num < 0 {
		fmt.Println(num, "is negative")
	} else if num < 10 {
		fmt.Println(num, "has 1 digit")
	} else {
		fmt.Println(num, "has multiple digits")
	}
}
```

### switch

当分支过多的时候，使用if-else语句会降低代码的可阅读性，这个时候，我们就可以考虑使用switch语句

- switch 语句用于基于不同条件执行不同动作，每一个 case 分支都是唯一的，从上至下逐一测试，直到匹配为止。
- switch 语句在默认情况下 case 相当于自带 break 语句，匹配一种情况成功之后就不会执行其它的case，这一点和 c/c++ 不同
- 如果我们希望在匹配一条 case 之后，继续执行后面的 case ，可以使用 fallthrough

```go
package main

import (
	"fmt"
	"time"
)

func main() {

	a := 2
	switch a {
	case 1:
		fmt.Println("one")
	case 2:
		// 在此打印"two"并跳出
		fmt.Println("two")
        //fallthrough，这里如果使用了fallthrough，所以会继续执行下一个case
	case 3:
		fmt.Println("three")
	case 4, 5:
		fmt.Println("four or five")
	default:
		fmt.Println("other")
	}

	t := time.Now()
	switch {
	// 根据现在的时间判断是上午还是下午
	case t.Hour() < 12:
		fmt.Println("It's before noon")
	default:
		fmt.Println("It's after noon")
	}
}
```

### func

函数是指一段可以直接被另一段程序或代码引用的程序或代码，一个较大的程序一般应分为若干个程序块，每一个模块用来实现一个特定的功能。

1. **函数的声明和定义**：

   在Go语言中，函数的定义以 `func` 关键字开始，然后是函数名、参数列表、返回类型和函数体。以下是一个函数的典型定义：

   ```go
   func add(x int, y int) int {
       return x + y
   }
   ```

   这个函数名为 `add`，接受两个整数参数 `x` 和 `y`，并返回一个整数。

2. **函数的参数**：

   函数可以接受零个或多个参数，参数在参数列表中定义，并且需要指定参数的类型。例如：

   ```go
   func greet(name string) {
       fmt.Println("Hello, " + name)
   }
   ```

   这个函数接受一个字符串参数 `name`。

3. **函数的返回值**：

   函数可以返回一个或多个值，返回值的类型也需要在函数定义中指定。如果函数没有返回值，可以将返回类型留空。例如：

   ```go
   func addAndMultiply(x, y int) (int, int) {
       sum := x + y
       product := x * y
       return sum, product
   }
   ```

   这个函数返回两个整数值。

4. **函数的调用**：

   要调用函数，只需使用函数名并传递参数。例如：

   ```go
   result := add(3, 5)
   fmt.Println(result)
   ```

   这里我们调用了 `add` 函数，将参数 `3` 和 `5` 传递给它，并将返回值赋给 `result` 变量。
   ```go
   package main

    import "fmt"

    func add(x int, y int) int {
        return x + y
    }
    
    func main() {
    	result := add(3, 5)
    	fmt.Println(result)
    }
   ```
5. **函数是一等公民**  
   函数是一等公民是指函数在语言中具有与其他数据类型（如数字、字符串等）相同的地位。

   这意味着函数可以被赋值给变量、作为参数传递给其他函数、作为返回值返回，甚至可以嵌套在其他函数中。 函数是一等公民的语言具有更高的表达能力，因为它可以用更简单的方式来编写代码。例如，在函数是一等公民的语言中，可以使用匿名函数（lambda）来创建临时函数，而无需为其分配名称。匿名函数可以用作回调函数，或者在其他函数中作为参数传递。

   还意味着函数可以作为数据结构的元素。例如，可以创建一个数组，其中每个元素都是一个函数。
```go
   
// 定义一个函数类型
type mathOperation func(int, int) int

// 一个普通的加法函数
func add(a, b int) int {
    return a + b
}

// 一个函数，接受另一个函数作为参数
func calculate(op func(int, int) int, a, b int) int {
    return op(a, b)
}

// 一个函数返回另一个函数
func getMultiplier() mathOperation {
    return func(a, b int) int {
        return a * b
    }
}

func main() {
    // 将函数赋值给变量
    var operation mathOperation
    operation = add

    // 调用函数
    result := operation(3, 4)
    fmt.Println("3 + 4 =", result)

    // 将函数作为参数传递给另一个函数
    result = calculate(add, 5, 6)
    fmt.Println("5 + 6 =", result)

    // 将函数作为返回值
    multiplier := getMultiplier()
    result = multiplier(3, 7)
    fmt.Println("3 * 7 =", result)
}
```

### fmt

fmt 库函数

```go
package main

import "fmt"

type point struct {
	x, y int
}

func main() {
	s := "hello"
	n := 123
	p := point{1, 2}
	fmt.Println(s, n) // hello 123
	fmt.Println(p)    // {1 2}

	// 使用 Printf 进行格式化输出
	// %v 以默认格式打印值，打印任意类型的变量的值
	fmt.Printf("s=%v\n", s)  // s=hello
	fmt.Printf("n=%v\n", n)  // n=123
	fmt.Printf("p=%v\n", p)  // p={1 2}
	fmt.Printf("p=%+v\n", p) // p={x:1 y:2}，%+v:结构体时会显示字段名和值
	fmt.Printf("p=%#v\n", p) // p=main.point{x:1, y:2}, %#v:值的Go语法表示

	f := 3.141592653
	fmt.Println(f)          // 3.141592653
	fmt.Printf("%.2f\n", f) // 3.14,以 浮点数格式（f）输出，并保留两位小数（.2）
}
```

#### 输入

除了输出，fmt 包还提供了读取用户输入的函数，最常用的是 `fmt.Scanln` 和 `fmt.Scan`：

```go
package main

import "fmt"

func main() {
	var name string
	var age int

	// Scanln 按空格/换行读取输入，依次存入后面的变量
	// 注意：要传变量的地址（用 &），这样读到的值才能写进变量里
	fmt.Print("请输入姓名和年龄：")
	fmt.Scanln(&name, &age)
	fmt.Printf("你好 %s，你今年 %d 岁\n", name, age)
}
```

- `fmt.Scanln(&变量1, &变量2, ...)` 读取用户输入，按空格/换行分隔后依次存入变量
- `fmt.Scan` 用法类似，但会按任意空白（空格、换行）连续读取
- 读取单个数字的常用写法：`fmt.Scanln(&num)`
- 变量前要加 `&` 传地址，这个符号在后面的指针课程中会深入讲解

### Go 语言特色语法

下面这些写法是 Go 比较有特色的地方，在其他主流语言里很少见，写代码时经常会遇到。

#### 空白标识符 `_`

`_`（下划线）是 Go 的**空白标识符（blank identifier）**，专门用来"丢弃"暂时用不到的值。它不能被读取，只能拿来占位忽略，而且可以反复使用。

```go
// 1. 忽略函数返回的某个值
// 字符串转数字会返回两个值：结果和错误，这里用 _ 忽略错误
n, _ := strconv.Atoi("123")
fmt.Println(n) // 123

// 2. 忽略 for range 返回的下标，只取元素值
for _, c := range "go" {
	fmt.Printf("%c ", c)
}
// g o
```

另外，数字字面量中也可以插入 `_` 来分隔位数、提高可读性，Go 会直接忽略它：

```go
price := 1_000_000 // 等价于 1000000
fmt.Println(price) // 1000000
```

#### 零值（zero value）

Go 的变量声明后即使不赋值，也会自动拥有一个默认值，叫**零值**，绝不会像 C 语言那样是一块随机的内存垃圾：

```go
var i int     // 0
var f float64 // 0
var s string  // ""（空字符串）
var b bool    // false
fmt.Println(i, f, b) // 0 0 false
fmt.Println(s == "") // true
```

> 所以 `var x int` 之后直接使用 `x` 也是安全的，不会报错或得到随机值。

#### iota：常量计数器

`iota` 是 Go 独有的常量计数器，只能用在 `const` 块中，从 0 开始，每换一行声明会自动加 1：

```go
const (
	A = iota // 0
	B        // 1，省略了 = iota
	C        // 2
)

const (
	_  = iota             // iota = 0，用 _ 跳过
	KB = 1 << (10 * iota) // iota = 1，KB = 1 << 10 = 1024
	MB                    // iota = 2，MB = 1 << 20 = 1048576
	GB                    // iota = 3，GB = 1 << 30 = 1073741824
)

fmt.Println(A, B, C)    // 0 1 2
fmt.Println(KB, MB, GB) // 1024 1048576 1073741824
```

`1 << n` 表示把 1 的二进制左移 n 位，也就是 2 的 n 次方。用 `iota` 可以很方便地定义一组递进的常量。

#### 首字母大小写：Go 的"公有/私有"

Go 没有 `public` / `private` 这样的关键字，**标识符首字母是否大写**就决定了它的可见性：

- 首字母**大写**：对外可见（导出），可以被其他包使用，比如标准库的 `fmt.Println`
- 首字母**小写**：仅当前包内可用，比如我们自己代码里的 `add` 函数

```go
// 包内私有：首字母小写，只能在当前包使用
func add(a, b int) int {
	return a + b
}

// 导出（公有）：首字母大写，其他包可以通过 包名.Add 调用
func Add(a, b int) int {
	return a + b
}
```

> 现在回头看，`fmt.Println`、`strings.Contains` 这些函数名首字母都大写，就是这个原因。

***

可能我在本节课中有很多没有讲到的细节，大家课下可以参考一下[Go语言教程](https://www.runoob.com/go/go-tutorial.html)以及[李文周的博客](https://www.liwenzhou.com/posts/Go/golang-menu/)结合进行学习

***

# 年轻人的第一个Go Project
***
## 课后作业
***
### LV0
编写一个Go程序，输出你的Hello➕你的名字（也可以是网名等）


### LV1
定义一个常量 pi = 3.14，再定义变量 r = 5，计算圆的面积（area = pi * r * r）并打印结果。

### LV2
使用for循环计算1到1000的和，并打印结果。

### LV3
编写一个函数，接受一个整数参数n，返回n的阶乘（factorial）。在main函数中调用该函数并打印结果。

### LVX
编写一个 Go 程序，实现以下功能：

1. 不断让用户输入整数（输入 `0` 表示结束输入）；
2. 使用 `for` 循环统计输入数字的总和与个数；
3. 定义一个函数 `average(sum int, count int) float64` 用于计算平均值；
4. 根据平均值输出结果：
   - 若平均值 ≥ 60，输出“平均成绩为 xx.xx，成绩合格”；
   - 否则输出“平均成绩为 xx.xx，成绩不合格”。

示例输入：
```
请输入一个整数(输入0结束): 80
请输入一个整数(输入0结束): 90
请输入一个整数(输入0结束): 40
请输入一个整数(输入0结束): 0
```
示例输出：
```
平均成绩为 70.00，成绩合格
```

思路提示
1. **定义变量** 保存总和 (`sum`) 与计数 (`count`)；
2. 使用 `for` 循环不断读取用户输入；
3. 当输入为 `0` 时，用 `break` 结束循环；
4. 调用自定义函数 `average(sum, count)` 计算平均值；
5. 使用 `if` 判断平均值是否合格并打印结果。  
6. 输入函数可以使用 `fmt.Scanln()` 来实现。

***

## 作业提交
在GitHub上新建一个仓库，命名为 `Lanshan-Go-2026-Homework`，将你的作业代码上传到该仓库中，
并把仓库地址发送到`wanghaihang@lanshan.email`  
建议大家自己研究一下git的使用，使用git把作业提交到仓库是更好的选择，关于git的教程在lesson0中有介绍，如果感觉命令行太复杂不会敲，可以下一个github desktop能够在图形界面上点点点，后面的课程我们也会详细讲解git的使用  
有问题在飞书群里提问或者飞书私聊我均可






