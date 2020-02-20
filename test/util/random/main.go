package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().Unix())

	arr := []int{1, 2, 3, 4, 5, 6, 7, 8}

	for i := 0; i < 3; i++ {
		rand.Shuffle(len(arr), func(i, j int) {
			arr[i], arr[j] = arr[j], arr[i]
		})
		fmt.Println(arr)
	}

}
