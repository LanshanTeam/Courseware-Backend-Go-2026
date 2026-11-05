package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

func main() {

	time.Sleep(1 * time.Second)
}

func Demo1() {
	fmt.Println(time.Now())
	go time.Sleep(5 * time.Second)
	fmt.Println(time.Now())
} //这是直接调用已有函数
func Demo2(a *int) {
	go func() {
		time.Sleep(5 * time.Second)
		(*a)++
	}()
} //这是启动了一个闭包函数

func Add1() {
	a := 0
	for range 100000 {
		a++
	}
	fmt.Println(a) //输出 100000 很合理
}

func Add2() {
	a := 0
	for range 10 {
		go func() {
			for range 10000 {
				a++
			}
		}() //启动10个goroutine进行计算
	}
	time.Sleep(5 * time.Second) //确保所以goroutine都执行完成
	fmt.Println(a)              //输出一个奇奇怪怪的数字，每次还不一样，但是一定小于100000
}
func AddAtomic() {
	var a atomic.Int64 // 0值可用
	a.Store(0)         // 赋值为0
	for range 10 {
		go func() {
			for range 10000 {
				a.Add(1)
			}
		}()
	}
	time.Sleep(5 * time.Second)
	fmt.Println(a.Load()) // 输出 100000
}
func AddMutex() {
	a := 0
	b := 100000
	var mu sync.Mutex // 0值可用
	for range 10 {
		go func() {
			for range 10000 {
				mu.Lock()
				a++
				b--
				mu.Unlock()
			}
		}()
	}
	time.Sleep(5 * time.Second)
	fmt.Println(a, b) // 输出 100000 0
}

func Channel1() {
	ch1 := make(chan int, 5)
	ch1 <- -1 //使用 ch<- 向通道中发送数据，因为缓冲区没有满，所以不会阻塞
	go func() {
		for i := range 5 {
			ch1 <- i
			fmt.Println("已放入", i)
		}
		close(ch1) //关闭通道，关闭后的通道不能再发送数据，会panic
	}()
	time.Sleep(1 * time.Second)
	for out := range ch1 {
		fmt.Println("out:", out)
	}
	fmt.Println("END")
}
func Channel2() {
	ch2 := make(chan int)
	go func() {
		time.Sleep(time.Second)
		close(ch2)
	}()
	data, ok := <-ch2 // 对的，读取操作也可以有两个接收值，ok表示的是读取的数据是否有效，正常读取时都为true，关闭后并且数据读取完才是false
	fmt.Println(data, ok)
}
func Channel3() {
	ch3 := make(chan int)
	go func() {
		time.Sleep(time.Second)
		close(ch3)
	}()
	ch3 <- 1
	fmt.Println("END")
}
func Channel4() {
	ch4 := make(chan int)
	close(ch4)
	close(ch4)
	fmt.Println("END")
}
func Channel5() {
	var ch5 chan string
	go func() {
		ch5 <- "望海夯"
		fmt.Println("成功发送")
	}()
	go func() {
		_, ok := <-ch5
		if ok {
			fmt.Println("成功接收")
		} else {
			fmt.Println("接收失败")
		}
	}()
	time.Sleep(5 * time.Second)
}
func Channel6() {
	ch6 := make(chan string)
	go func() {
		ch6 <- "魂星"
		fmt.Println("成功发送")
	}()
	time.Sleep(2 * time.Second)
	go func() {
		data, ok := <-ch6
		if ok {
			fmt.Println("第一个goroutine接收成功：", data)
		} else {
			fmt.Println("第一个goroutine未成功接收")
		}
	}()
}
func Channel7() {
	ch7 := make(chan string)
	go func() {
		data, ok := <-ch7
		if ok {
			fmt.Println("第一个goroutine接收成功：", data)
		} else {
			fmt.Println("第一个goroutine未成功接收")
		}
	}()
	time.Sleep(2 * time.Second)
	go func() {
		ch7 <- "魂星"
		fmt.Println("成功发送")
	}()
}
func Channel8() {
	test1 := func(ch chan<- int) {
		ch <- 1
		//<-ch
	}
	test2 := func(ch <-chan int) {
		<-ch
		//ch<-1
	}
	ch := make(chan int, 1)
	go test1(ch)
	go test2(ch)
	time.Sleep(1 * time.Second)
}

func Add3() {
	a := 0
	var mu sync.Mutex
	var wg sync.WaitGroup // 依旧0值可用
	for range 10 {
		wg.Add(1) //也可以直接在循环外Add(10)，反正必须在goroutine外面Add，并且数量得一致
		go func() {
			for range 10000 {
				mu.Lock()
				a++
				mu.Unlock()
			}
			wg.Done() // 其实等价于wg.Add(-1)，源码就是这样的
		}()
	}
	wg.Wait()
	fmt.Println(a)
}
func Add4() {
	a := 0
	var mu sync.Mutex
	var wg sync.WaitGroup
	for range 10 {
		wg.Go(func() {
			for range 10000 {
				mu.Lock()
				a++
				mu.Unlock()
			}
		})
	}
	wg.Wait()
	fmt.Println(a)
}

func Defer1() {
	defer func() {
		fmt.Println("defer1")
	}()
	a := rand.Int() % 2
	if a == 1 {
		return
	}
	fmt.Println("a不是1")
}
func Defer2() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("捕获到panic：", err)
		}
	}()
	panic("一个panic")
	fmt.Println("我还活着")
	panic("第二个panic")
}
func Defer3() {
	defer fmt.Println("第一个defer")
	defer fmt.Println("第二个defer")
	defer fmt.Println("第三个defer")
}
func Defer4() {
	test1 := func() int {
		a := 1
		defer func() {
			a++
		}()
		return a
	}
	test2 := func() (a int) {
		a = 1 //这里不用:=是因为在返回参数那里已经声明了
		defer func() {
			a++
		}()
		return //等价于return a，相当于是默认return你返回参数中声明的变量
	}
	fmt.Println(test1()) //输出 1
	fmt.Println(test2()) //输出 2
}

func Select1() {
	ch1 := make(chan int, 1)
	ch2 := make(chan struct{})
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
	case ch2 <- struct{}{}:
		fmt.Println("ch2已发送数据")
	case <-timer.C:
		fmt.Println("已超时")
	}
}
func Select2() {
	ch1 := make(chan struct{}, 1)
	//ch1<-struct{}{}
	select {
	case <-ch1:
		fmt.Println("ch1接收到数据")
	default:
		fmt.Println("未接收数据，继续执行")
	}
}

func Context1() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	ch := make(chan int, 10)
	var wg sync.WaitGroup
	wg.Go(func() {
		ticker := time.NewTicker(1 * time.Second)
		a := 0
		for {
			select {
			case <-ticker.C:
				ch <- a
				a++
			case <-ctx.Done():
				ticker.Stop()
				close(ch)
				fmt.Println("1 已退出发送")
				return
			}
		}
	})

	wg.Go(func() {
		for {
			select {
			case data, ok := <-ch:
				if !ok {
					fmt.Println("2 已退出接收")
					return
				}
				fmt.Println("2 已接收数据：", data)
			case <-ctx.Done():
				fmt.Println("2 已退出接收")
				return
			}
		}
	})
	wg.Wait()
}
func Context2() {
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Go(func() {
		ticker := time.NewTicker(500 * time.Millisecond)
		fmt.Println("准备就绪")
		for {
			select {
			case <-ticker.C:
				fmt.Println("ping")
			case <-ctx.Done():
				ticker.Stop()
				fmt.Println("已退出")
				return
			}
		}
	})
	time.Sleep(3 * time.Second)
	cancel()
	wg.Wait()
}
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

func Sync1() {
	var init sync.Once //依旧是0值可用
	num := 0
	for range 10 {
		init.Do(func() {
			num++
		})
	}
	fmt.Println(num) //输出 1
}
func Sync2() {
	var m sync.Map //还是0值可用，貌似sync包的常用的都是0值可用，没了解过hh
	//1. 存：Store(key,value)
	m.Store("name", "望海夯")
	m.Store("age", 18)

	//2. 取：Load(key)(value any,ok bool)
	name, ok := m.Load("name")
	if ok {
		fmt.Println(name)
	}

	//3. 存在就返回，不存在就写入：LoadOrStore
	v, loaded := m.LoadOrStore("gender", "male")
	fmt.Println(v, loaded) // loaded=true代表原本就存在

	// 4. 删除：Delete
	m.Delete("age")

	// 5. 遍历：Range
	m.Range(func(k, v any) bool {
		fmt.Printf("key:%v, val:%v\n", k, v)
		return true // 返回true继续遍历，false停止
	})
}
