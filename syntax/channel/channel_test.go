package channel

import (
	"fmt"
	"testing"
)

func TestChannel(t *testing.T) {
	//声明了，但是 还未初始化,这样会报错
	ch := make(chan int, 3)
	ch <- 1
	ch <- 2
	ch <- 3
	ch <- 4
	fmt.Println(ch)
}
