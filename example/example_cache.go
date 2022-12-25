package main

import (
	"fmt"
	"github.com/lfu"
)

func main() {
	cache := lfu.NewLFU(4)
	cache.Set("Key_1", 1)
	cache.Set("Key_2", 2)
	cache.Set("Key_3", 3)
	cache.Set("Key_4", 4)

	data, err := cache.Get("Key_2")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Key_2=", data)

	cache.Set("Key_2", "NewValue")

	data, err = cache.Get("Key_2")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Key_2=", data)
	cache.Set("Key_5", 5)

	data, err = cache.Get("Key_5")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Key_5=", data)

}
