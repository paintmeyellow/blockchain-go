package main

import "fmt"

func main() {
out:
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			fmt.Println(i+j)
			if i+j == 20 {
				fmt.Println("i+j==20")
				break out
			}
		}
	}
}
