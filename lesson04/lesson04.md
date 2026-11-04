# 前言
前面我们学习了go语言的大部分语法和关键字、go的包管理、面向对象和go语言的常用官方库。
在lesson03中提到了我们可以通过`<-time.Ticker的ticker.C`进行定时任务，你可能会因为不知道`<-`的含义而迷茫，
那么这节课我们将会讲解一下go语言的一个特色：channel，以及相关的并发知识。  

# 前情提要
## 操作系统
我们的电脑在硬件层面由：中央处理器(cpu),储存器(磁盘和内存),输入输出设备以及其他附加设备组成。
但是除了涉及到硬件层面的处理的时候(比如单片机)，我们往往不需要直接处理硬件的数据，而操作系统就是硬件和软件之间的这个桥梁。  

操作系统将硬件层的具体的行为封装为系统调用的接口给软件层进行使用。虽然不同的操作系统之间提供的接口都差不多，但是还是有一定的区别。
所以有的语言，尤其是c、cpp、go这类编译为cpu二进制可执行文件的编译性语言来说，不同的操作系统之间不能共通，导致有的应用程序会有多个操作系统的版本，这也可以看到跨平台的能力对于一门语言来说也是一项加分项。  

虽然在后续的课程以及开发过程中，我们很少会直接通过系统调用，大多数情况都是使用封装好的包进行编程，但是了解底层知识可以有助于设计出更快、更优雅的~~屎山~~代码。  

## 进程和线程
你可能疑惑过，既然一个cpu的核心在同一时间之内只能执行一个任务，那为什么我在电脑上一次性开了几十个软件，它们却好像都在单独运行一样。
在回答问题之前我们先了解一下操作系统的最小任务：**进程**，以及cpu处理的最小任务：**线程**  
### 线程
首先我们需要知道cpu的计算速度是非常快的，计算单位可以以纳秒(10^-9秒)算。所以对于cpu来说，我们的世界是很慢的。  
想象一下，你如果有现在一百倍的速度，但是外界不变，相当于是原本100分钟的作业你可以1分钟做完。
那么是不是我们就可以同时写100份作业，期间我们可能不停的更换作业，但是我们的行为足够快，在外人眼里我们表现出来的就是我们可以同时写100份作业。  

那么我们说回cpu上，我们刚才就相当于是cpu的一个核心，而每一个线程就是一个作业。
这解释了，为什么即使单个核心在同一时间内只能执行一个任务，在我们使用者的视角来看是同时执行了很多的任务。  

当然线程的本身和切换也是有开销的，就好比你在写高数写到一半跑去做大物，那么你需要记住高数的思路等数据，还需要想一想大物做到哪了。
这些对于cpu来说就是线程本身的储存开销和切换开销，一个线程保底基本是8M的空间，所以线程的创建、销毁、调度和上下文切换本身是很重的。
此外线程的切换还涉及到lesson03提到过的用户态和内核态的切换，我们后面再说。  
### 进程
进程是一个正在运行的程序实例，操作系统会为每一个进程分配彼此隔离的用户空间，避免不同的程序访问对方的数据~~真要能这样我打游戏随便就可以一刀999~~。
一个程序可以有一个进程，也可以有多个，这取决于应用本身。  

同一个进程至少包含了一个线程，往上基本不封顶，但是如果进程像操作系统申请的内存太多了还是有被直接OOM的风险的。  

如果用刚才的例子做比如，那么进程就像是某个学科总的的作业，而线程就是该学科的具体的某一本作业，一个科目可以同时有很多本作业，但是不能没有作业。
以后面试的时候，如果面试官问进程和线程的区别，你就可以想想今天的例子。  
## 一个小结
上到这里，你可能还有点懵，不重要，这下东西可以以后慢慢消化，go本身的并发是很简单的。
go的runtime的**GMP模型**(课后可以去了解一下)已经将线程的问题处理好了，实际业务的时候不用考虑线程池等细节。  
# go并发编程
## goroutine
### 前~情提要
~~怎么这么多前情提要~~

goroutine是go语言处理并发的最基本的单位，这时候你可能有疑惑了，既然线程已经是cpu处理的最小单位了，那goroutine为什么会有比线程更轻量、开销更小的特点呢。
这时候就可以提到前面说的用户态和内核态了。  
前面说过，操作系统为每一个进程提供了相互隔离的运行环境，让进程不能直接访问硬件，所以像访问硬盘、网卡等操作需要由操作系统代为执行。
这里普通进程运行时就是处于用户态，当程序发起系统调用时，需要cpu先切换成权限更高的内核态，再执行相应的操作。  
一个典型的例子就是读取文件：
1. 进程调用了Read读取文件
2. 由于触发系统调用，线程进入等待
3. cpu切换为内核态
4. cpu访问硬盘将数据存在内核态的内存中
5. cpu切换回用户态
6. 程序继续执行

这个例子也说明了为什么直接读取文件的缓冲区只能读一次并且不能写入，因为用户态进程不能给操作内核态的内存数据。  
### 什么是goroutine
goroutine是用户级别的线程，相比传统操作系统级别的线程，goroutine的创建、销毁、调度和上下文切换是由go底层runtime在用户态进行管理，不用涉及到cpu频繁的切换状态，
并且单个goroutine保底只有8k(不同的系统可能有区别，但是都很少)，相比传统线程也会轻量很多。  
用前面的作业的例子来说的话，线程是一本作业，那goroutine就算是其中的一道题目(别管这个题目多大，就算占了10页的题目不也是一个题吗)，执行不同的goroutine只需要在同一个本子上找，就不用不断切换本子了。  

### 如何使用goroutine
前面讲了这么多，你可能暂时不是很理解，不重要，目前你们只需要知道相比传统线程，goroutine更轻量，可以创建更多。
在go语言中，创建一个goroutine非常简单，只需要在需要执行的函数前面加上`go`关键字即可：
```go
func Demo1(){
	fmt.Println(time.Now())
	go time.Sleep(5*time.Second)
	fmt.Println(time.Now())
}//这是直接调用已有函数

func Demo2(a *int){
	go func(){
		time.Sleep(5*time.Second)
        (*a)++
    }()
}//这是启动了一个闭包函数
```
总之你只需要在你想异步执行的函数前面加上`go`即可。  
当然，创建很简单，难点是在于如何管理goroutine以及处理并发问题~~真要这么点东西就不至于单独开一节课了~~

## 互斥锁和原子操作
现在我们学了如何启动goroutine并发执行任务，你看到你手上执行十万次加法的单线程操作，觉得太慢了，想给他升级为10个goroutine并发执行，这样不就快10倍了吗，说干就干：
```go
func Add1(){
	a:=0
	for range 100000{
		a++
    }
	fmt.Println(a)//输出 100000 很合理
}

func Add2(){
	a:=0
	for range 10{
		go func(){
			for range 10000{
				a++
            }
        }()//启动10个goroutine进行计算
    }
	time.Sleep(5*time.Second)//确保时候goroutine都执行完成
	fmt.Println(a)//输出一个奇奇怪怪的数字，每次还不一样，但是一定小于100000
}
```
这里我们可以发现，Add2的结果和预期的不同，这就是并发的一个典型的问题。  
从代码来看，a++是一次操作，不就是a+1嘛，但是实际上，对于cpu来说这是一个多部操作：
1. cpu将a的值从内存中取出放入cpu寄存器
2. cpu从寄存器中取出+1后放回寄存器
3. cpu将值重新放回内存中a的地址

我们可以看出，如果cpu的不同核心同时执行上述操作时，可能会出现，核心1读到a为1，核心2也读到1，然后它们各自+1，再写入2，聪明的你已经发现了，虽然执行了两次+1操作，但是实际只增加了1。
为了解决这个问题，我们就需要使用互斥锁或者原子操作。

### 原子操作
针对上述的赋值操作，go标准库提供了原子操作库`sync/atomic`，对于一些变量的赋值操作，可以原子性的进行修改：
```go
func AddAtomic(){
	var a atomic.Int64// 0值可用
	a.Store(0)// 赋值为0
	for range 10{
		go func(){
		    for range 10000{
				a.Add(1)
            }
		}()
    }
	time.Sleep(5*time.Second)
	fmt.Println(a.Load())// 输出 100000
}
```
如果只是对于单个变量来说，Atomic基本已经够用了，但是如果我们需要确保多部操作的原子性，我们就需要引入`sync`包的互斥锁。

### 锁
#### 互斥锁 sync.Mutex
Atomic只能支持部分类型，用起来往往不够灵活，比如我们想限制某一些操作在同一时刻只能有一个线程执行，Atomic便无法实现，于是我们引入了互斥锁。
先来个简单的例子：
```go
func AddMutex(){
    a:=0
    b:=100000
    var mu sync.Mutex// 0值可用
    for range 10{
        go func(){
            for range 10000{
                mu.Luck()
                a++
                b--
                mu.Unluck()
            }
        }()
    }
    time.Sleep(5*time.Second)
    fmt.Println(a,b)// 输出 100000 0
}
```
在这个例子中，我们可以使用sync.Mutex保护a+1和b-1两个操作。那么为什么锁可以保护a，b两个变量？

同一个时刻只能有一个线程能够拿到同一把锁，然后它才能执行接下来的代码。如果锁被其他线程持有的时候调用Lock(),由于无法拿到锁，也就无法向下继续执行代码，这个goroutine会在自旋几次后被挂起，不消耗cpu资源(对底层感兴趣的可以看看Lock源码和调度器)。
当调用UnLock()释放锁时，相关的线程会被唤醒并争抢锁。这里涉及到一个概念叫做**惊群效应**，感兴趣可以了解一下。

#### 读写锁 sync.RWMutex
很多时候我们仅仅只是需要读取某个数据但是不用修改它的值，如果使用sync.Mutex轮流读取的话就会很慢，于是对于读多写少的场景来说，这是一个浪费，所以我们可以使用读写锁对相应的代码块进行管理。

相比普通的锁只有一把，读写锁有两种锁：只有一把的写锁以及想要多少有多少的读锁。使用读写锁可以做到读和写分离。
读锁的获取和释放为RLock()RUnlock()方法，写锁和互斥锁一样，是Lock()Unlock()方法。
当读锁被持有的时候，其他进程可以正常获取读锁。
当写锁被持有的时候，任何锁的获取都会被阻塞。
其中有个特殊的情况，读锁被持有的时候，写锁被尝试获取，那么之后进来的读锁也会被阻塞，说白了就是写锁具有优先级。

来个例子吧，go的map如果并发读写会触发panic，正好我们可以使用新学的读写锁有效避免这个问题：
```go
type Hub struct {
    Data map[string]string
    RwMu sync.RWMutex
}

func NewHub(keys []string, values []string) *Hub {
    newHub := &Hub{
        Data: make(map[string]string),
    } // 同样的，读写锁也是0值可用
    //这里是个demo就不防御性编程了
    for i := range keys {
        newHub.Data[keys[i]] = values[i]
    }
    return newHub
}
func (h *Hub) Set(key string, value string) {
    h.RwMu.Lock()
    defer h.RwMu.Unlock() // 这个defer后面详细讲讲，这里知道是延迟执行即可
    time.Sleep(1 * time.Second)
    h.Data[key] = value
}
func (h *Hub) Get(key string) (value string, exist bool) {
    h.RwMu.RLock()
    defer h.RwMu.RUnlock()
    value, exist = h.Data[key]
    time.Sleep(500 * time.Millisecond)
    return // 等价于return value,exist ，因为返回参数有变量名，编译器会自动补全，感兴趣自行了解
}
func RwMutex() {
    dataHub := NewHub([]string{"第一条消息", "第二条消息"}, []string{"1旧数据", "2旧数据"})
    get := func(id int, key string) {
        data, exist := dataHub.Get(key)
        fmt.Println(id, "get: ", key, ":", data, ":", exist)
    }
    set := func(id int, key string, value string) {
        dataHub.Set(key, value)
        fmt.Printf("%d set:  %s已更新为%s\n", id, key, value)
    }
    go get(1, "第一条消息")
    go get(2, "第二条消息")
    go set(3, "第一条消息", "1新数据")
    time.Sleep(100 * time.Millisecond)
    go get(4, "第一条消息")
    go get(5, "第二条消息")
    time.Sleep(5 * time.Second)
}
```
由于runtime的调度和Println清空缓冲的问题，输出的结果比较随机，但是可以发现4和5永远是在3号set之后运行，这正好说明等待持有写锁期间是不会发放读锁的，并且程序也没有panic，很好的保护了我们的map。
当然go官方库中也有sync.Map并发安全的map，我们后面再说。

## 通道 channel
go语言有个经典的话 **“不要通过共享内存通信，而要通过通信来共享内存”**,这句话是go语言的并发哲学。
这句话说的是，两个goroutine希望共享一个变量，不应该通过一个为外部的变量进行加锁而是应该通过 channel 将 goroutine A 中的变量传递给 goroutine B。

实际上我觉得，对于对象本身的数据或者需要在内存中长期储存的数据来说，加锁读取是最好的方法。后面我们会说channel只能做数据的转移，取了就没了，多个goroutine如果都需要这份数据的话还是老实加锁吧。
实际使用的时候还是看情况，后面我会分享几个go比较成熟的并发模型，可以体会一下goroutine的用处。

那我们先从基础说起。
### channel的创建
channel和map、切片一样，是一种引用类型，可以直接在函数之间传递。
当然，创建的时候也需要使用make()
```go
ch1:=make(chan int,5)// 变量ch1的类型是“ch int",缓冲区的大小为5，可以放5个int类型的数据
ch2:=make(chan myType)// 这里也可以写自定义的类型，第二个参数不写默认为0，表示这是一个无缓冲的通道
var ch3 chan string// 未初始化，ch3为默认的nil
```
### channel的使用
#### 发送、接收
我们可以向channel中传入或者取出数据：
```go
func Channel1(){
	ch1:=make(chan int,5)
	ch1 <- -1 //使用 ch<- 向通道中发送数据，因为缓冲区没有满，所以不会阻塞
	go func(){
		for i:=range 5{
			ch1 <- i
			fmt.Println("已放入",i)
        }
		close(ch1)//关闭通道，关闭后的通道不嫩再发送数据，会panic
    }()
	time.Sleep(1*time.Second)
	for out:=range ch1{
	    fmt.Println("out:",out)	
	}
	fmt.Println("END")
}
```
`for range ch1`会一直尝试读取数据，直到close(ch1)
#### close
接下来我们说一下close的一些问题：
##### 操作已关闭的channel
```go
func Channel2(){
	ch2:=make(chan int)
    go func(){
		time.Sleep(time.Second)
	    close(ch2)
    }()
	data,ok:=<-ch2// 对的，读取操作也可以有两个接收值，ok表示的是通道是否正常，正常读取时都为true，关闭后才是false
	fmt.Println(data,ok)
}
```
我们会在一秒之后打印 0 false，表示读取已关闭的channel不会阻塞，并且阻塞中的读取操作也会直接接收到零值.

```go
func Channel3(){
	ch3:=make(chan int)
	go func(){
	    time.Sleep(time.Second)
		close(ch3)
    }()
	ch3<- 1
	fmt.Println("END")
}
```
这里会直接触发panic，说明已关闭的通道是不可写入的，所以我们在使用的时候一般是由发送方管理channel的close，避免别的地方一个close，发送方直接panic。

##### 重复close
```go
func Channel4(){
	ch4:=make(chan int)
	close(ch4)
	close(ch4)
    fmt.Println("END")
}
```
这里会触发panic，提示通道不能够重复关闭，所以多个发送方管理channel的时候一定要确保close只能够执行一次，

### 一些channel的补充
#### 缓冲区以及是否初始化的区别
先说重要的，如果一个channel不初始化会怎么样：
```go
func Channel5(){
	var ch5 chan string
	go func(){
        ch5<-"望海夯"
        fmt.Println("成功发送")
    }
	go func(){
		data,ok:=<-ch5
		if ok{
			fmt.Println("成功接收")
        }else{
			fmt.Println("接收失败")
        }
    }
	time.Sleep(5*time.Second)
}
```
我们会发现等了半天，一行输出都没有，这说明无论是发送还是接收都是永久阻塞，当然不是说不能channel赋值为nil，后面学习了select之后我们就可以通过nil值的channel关闭特定的case，算是高阶操作了~~反正我没用过~~。

然后是无缓冲通道：
```go
func Channel6(){
	ch6:=make(chan string)
	go func(){
		ch6<-"魂星"
        fmt.Println("成功发送")
    }()
	time.Sleep(2*time.Second)
    go func(){
        data,ok:=<-ch6
        if ok{
            fmt.Println("第一个goroutine接收成功：",data)
        }else{
            fmt.Println("第一个goroutine未成功接收")
        }
	}()
}
func Channel7(){
	ch7:=make(chan string)
	go func(){
        data,ok:=<-ch7
        if ok{
            fmt.Println("第一个goroutine接收成功：",data)
        }else{
            fmt.Println("第一个goroutine未成功接收")
        }
    }()
    time.Sleep(2*time.Second)
    go func(){
        ch7<-"魂星"
        fmt.Println("成功发送")
	}()
}
```
我们会发现两个例子中两个goroutine的输出总是同时出现，相当于是无论是先接收还是先发送，都会阻塞，直到有对应的goroutine执行发送或者接收，两个goroutine才会继续执行。

最后是有缓冲的通道，大家可以就上面的两个例子改一下。
简单来说就是发送的时候，只有在缓冲区满了的情况下才会阻塞。
而接收的时候，只有在缓冲区为空的时候才会阻塞。
#### <-ch 类型和ch<-类型
前面我们的channel的类型都是`ch type`，但是有的时候我们只希望这个channel只能发送或者接受，这时候就可以通过`<-ch type`和`ch<- type`进行限制：
```go
func Channel8(){
    test1:=func(ch chan<- int){
        ch <- 1
        //<-ch
    }
    test2:=func(ch <-chan int){
        <-ch
        //ch<-1	
    }
    ch:=make(chan int,1)
    go test1(ch)
    go test2(ch)
    time.Sleep(1 * time.Second)
}
```
## WaitGroup
前面我们为了确保goroutine能够运行完成是采用的在主线程中time.Sleep的方式，这种方式一方面是不够灵活，另一方面是设置的时间不一定能保证goroutine执行完成。
所以sync包中有一个WaitGroup用于帮助我们解决这一个问题，例如前面100000次累加就可以升级一下：
```go
func Add3(){
	a:=0
	var mu sync.Mutex
	var wg sync.WaitGroup // 依旧0值可用
	for range 10{
		wg.Add(1)//也可以直接在循环外Add(10)，反正必须在goroutine外面Add，并且数量得一致
        go func(){
			for range 10000{
				mu.Lock()
				a++
				mu.Unlock()
			}
			wg.Done()// 其实等价于wg.Add(-1)，源码就是这样的
        }
    }
	wg.Wait()
	fmt.Println(a)
}
```
WaitGroup本质上就是一个计数器，虽然我们也可以使用锁手搓一个，但是有现成的工具不是更好吗。

此外你是不是觉得每次都需要Add(),Done()很麻烦，没事，go 1.25添加了一个Go(func())方法，我们只需要将需要执行的代码封装为无传参无返参的闭包函数传进去即可：
```go
func Add4(){
	a:=0
	var mu sync.Mutex
	var wg sync.WaitGroup
	for range 10{
		wg.Go(func(){
            for range 10000{
                mu.Lock()
                a++
                mu.Unlock()
            }
        })
    }
	wg.Wait()
	fmt.Println(a)
}
```
## 一些关键字
到目前为止，我们已经学会了go的并发编程，但是对于前面提到过的`defer` `select`关键字还很陌生，那么我们现在讲解一下这最后的两块拼图
### defer
#### 基础用法
前面的例子中：
```go
func Afunc(){
	var mu sync.Mutex
	go func(){
		mu.Luck()
		defer mu.Unluck()
		......
    }()
}
```
当时的意思是使用defer延迟执行释放锁这个操作，事实上，`go`和`defer`的使用方式相同，只是含义不同。
`go`是启动一个goroutine异步执行后面的函数，而`defer`是将后面的函数注册，在当前函数return之后执行。
```go
func Defer1(){
	defer func(){
		fmt.Println("defer1")
    }()
	a:=rand.Int()%2
	if a==1{
		return
    }
	fmt.Println("a不是1")
}
```
我们会发现无论随机进去哪一个分支，什么时候return，defer中的函数都被执行了。

即使程序panic了defer也会依次执行，事实上panic不一定会使程序崩溃，正好和defer一起讲了，算是很常见的用法了：
```go
func Defer2(){
	defer func(){
		if err:=recover();err!=nil{
			fmt.Println("捕获到错误：",err.Error())
        }
    }()
	panic("一个panic")
	fmt.Println("我还活着")
    panic("第二个panic")
}
```
这里的recover就可以捕获到逐层展开的panic并终止展开行为，并且recover只能在defer中使用，如果在正常函数中使用就只会返回nil。
所以很多时候，例如执行文件操作，网络操作等时IDE经常弹出`defer Close()`的提示，就是为了避免程序中途某个分支退出或者panic的时候忘记Close()释放资源，造成文件句柄泄漏等问题。

然后一个函数中如果注册了多个defer，则会按照后注册先执行的顺序依次执行：
```go
func Defer3(){
	defer fmt.Println("第一个defer")
    defer fmt.Println("第二个defer")
    defer fmt.Println("第三个defer")
}
```
最后执行的结果就是 三，二，一 的顺序
#### 进阶(八股)
我们直到，在go中我们是可以给返回的参数命名的，例如`func example()(data int,err error)`,那么如果结合defer使用就有一些细节问题：
```go
func Defer4(){
	test1:=func()(int){
		a:=1
		defer func(){
		    a++	
        }()
		return a
    }
	test2:=func()(a int){
		a=1//这里不用:=是因为在返回参数那里已经声明了
		defer func(){
			a++
        }()
		return //等价于return a，相当于是默认return你返回参数中声明的变量
    }
	fmt.Println(test1())//输出 1
    fmt.Println(test2())//输出 2
}
```
可以理解为。当你不声明返回参数的名字的时候，go随便给你生成了一个名字，然后在return的时候将值赋给了返回参数，这时候你再用defer修改就无事发生了。
但是test2中因为本身有变量名 a 了，return就相当于直接将值赋给了a，所以defer就可以修改return的值。
反正现在先了解一下吧，遇到了至少知道有这么一回事。
#### 坑
前面的Add中你可能发现了我没有使用defer mu.Unlock(),而是老老实实的mu.Unlock()手动管理。这是因为这里是在循环内部，如果在循环内部执行defer的话，每一个循环都注册一个defer函数，那开销也不小了。
而且在刚才的例子中，如果通过defer释放锁，那么Unlock就会在循环外面才执行，实际上在第二次循环的时候就直接因为锁未释放就死锁了。
所以有的时候虽然defer很方面，但是该手动管理的时候也别想着偷懒了。

### select
前面提到channel的发送和接收是会发生阻塞的，但是实际情况下不会让你一直阻塞在哪里，这时候我们可以使用`select`关键字通过监听多个channel，只需要至多一个channel就绪就可以接着执行。
这里直接说可能不太清楚，先来个例子：
```go
func Select1() {
    ch1 := make(chan int, 1)
    ch2 := make(chan struct{}, 1)
    timer := time.NewTimer(5 * time.Second)
    beginner := make(chan struct{})
    go func() {
        <-beginner
        time.Sleep(2 * time.Second)
        ch1 <- 1
        close(ch1)
    }()
    go func() {
        <-beginner
        time.Sleep(2 * time.Second)
        <-ch2
        close(ch2)
    }()
    close(beginner)
    select {
    case data := <-ch1:
        fmt.Println("ch1接收到数据:", data)
    case ch2<-struct{}{}:
        fmt.Println("ch2接收到数据")
    case <-timer.C:
        fmt.Println("已超时")
    }
}
```
我们可以通过修改两个goroutine的Sleep参数控制channel，可以发现select(没使用default)的一个性质：
1. 所有的channel未准备完成的时候会阻塞
2. 阻塞状态下哪一个case的channel先准备好就执行对应的代码
3. 当多个case同时就绪时会随机选择一个执行

当然有个时候我们不希望select阻塞，就可以使用`default`：
```go
func Select2(){
	ch1:=make(chan struct{},1)
    //ch1<-struct{}{}
	select{
	case <-ch1:
        fmt.Println("ch1接收到数据")
    default:
		fmt.Println("未接收数据，继续执行")
    }
}
```
可以看到原本应该阻塞在`<-ch1`处的代码直接执行了default

## context
在前面的例子中，一个goroutine就做一件事情，做完了就返回。
实际上在项目中goroutine中还可能会开新的goroutine，一旦上层的goroutine因为某些原因，例如对端关闭网络连接等情况中断执行，那下层的goroutine也应该终止而不是无视风险继续执行。
为了解决这个问题，刚好官方就提供了context包用于管理
### 基础
首先是一些默认的开发习惯，context.Context一般是作为函数的第一个参数传入，变量名字为**ctx**，**绝对，绝对不要放在结构体里面**。  

#### 创建
首先是新建一个原始的context，只需要调用`context.Background()`即可，

#### 基础方法
context.Context本质上是一个接口：
```go
type Context interface{
	Deadline()(deadline time.Time,ok bool)//获取截止时间以及是否超时
	Done()<-chan struct{}//取消的时候会close这个通道，后面会讲用法
	Err()error//返回取消的原因，一般是：Canceled/DeadlineExceeded
	Value(key any)any//获取值，和map差不多
}
```
实际上我们直接用官方的实现就行。
##### 超时和取消
前面说到context需要控制goroutine的生命周期，所以超时和取消是少不了的
```go
ctx:=context.Background()

ctx,cancel:=context.WithCancel(ctx)//设置手动取消，cancel必须调用，防止goroutine泄漏
ctx,cancel:=context.WithTimeout(ctx,5*time.Second)//设置超时时间，cancel同上
ctx,cancel:=context.WithDeadline(ctx,time.Now().Add(5*time.Second))//设置截止时间，cancel同上
ctx=context.WithValue(ctx,"key","value")//储存数据，不过不要放业务数据，一般只存请求id等请求级别的数据
```
此外cancel()不像close()一样反复调用就死给你看，cancel可以多次调用，但是只有第一次调用的时候会取消
##### 判断超时或者是否取消
前面接口中看到，Done方法会返回一个只读的channel，所以用法其实很简单，比如我们可以把前面的那个超时的例子改一下：
```go
func Context1(){
	ctx,cancel:=context.WithTimeout(context.Background(),2*time.Second)
	defer cancel()
	ch:=make(chan int,10)
	var wg sync.WaitGroup
	wg.Go(func(){
		ticker:=time.NewTicker(1*time.Second)
        a:=0
        for{
            select{
            case <-ticker.C:
                ch<-a
                a++
            case <-ctx.Done():
                ticker.Stop()
                close(ch)
                fmt.Println("1 已退出发送")
                return
            }
        }
    })
	
	wg.Go(func(){
	    for{
			select{
			case data:=<-ch:
				fmt.Println("2 已接收数据：",data)
            case <-ctx.Done():
				fmt.Println("2 已退出接收")
                return
			}
        }
	})
	wg.Wait()
}

func Context2(){
	ctx,cancel:=context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Go(func(){
		ticker:=time.NewTicker(500*time.Millisecond)
		fmt.Println("准备就绪")
		for{
			select{
			case <-ticker.C:
				fmt.Println("ping")
			case <-ctx.Done():
				ticker.Stop()
				fmt.Println("已退出")
				return
            }
        }
    })
	time.Sleep(3*time.Second)
	cancel()
	wg.Wait()
}
```
可以看到，不管是手动退出还是超时退出都会导致这个context关闭

### 链式取消
上面的例子中我们只是知道了context可以用于结合select管理goroutine，那么goroutine嵌套创建goroutine的问题还没有解决。
其实从底层来说context是一个树状结构，上层的context取消会使下层的子context全部取消，这也刚好满足goroutine的嵌套的情况。
来个例子：
```go
func Context3() {
    ctx1, cancel1 := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel1()
    var wg sync.WaitGroup
    wg.Go(func() {
        begin1 := time.Now()

        ctx2, cancel2 := context.WithTimeout(ctx1, 2*time.Second)
        defer cancel2()

        wg.Go(func() {
            begin2 := time.Now()

            ctx3, cancel3 := context.WithTimeout(ctx2, 3*time.Second)
            defer cancel3()

            wg.Go(func() {
                begin3 := time.Now()

                <-ctx3.Done()
                fmt.Println("3 退出，总耗时：", time.Since(begin3).Seconds())
            })

            <-ctx2.Done()
            fmt.Println("2 退出，总耗时：", time.Since(begin2).Seconds())
        })

        <-ctx1.Done()
        fmt.Println("1 退出，总耗时：", time.Since(begin1).Seconds())
    })
    wg.Wait()
}
```
这里可以发现2和3几乎同时退出，即使3的超时为3秒，但是因为2触发超时取消了，3也跟着取消，同时1没有受到影响