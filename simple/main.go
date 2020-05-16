//
// Simple FizzBuzz in Go
//

package main

import (
	"fmt"
	"strconv"
)

// func fizzbuzz_print(x int) {
// 	if x%15 == 0 {
// 		fmt.Printf("FizzBuzz\n")
// 	} else if x%3 == 0 {
// 		fmt.Printf("Fizz\n")
// 	} else if x%5 == 0 {
// 		fmt.Printf("Buzz\n")
// 	} else {
// 		fmt.Printf("%d\n", x)
// 	}
// }

func main() {
	for i := 0; i < 20; i++ {
		result := fizzbuzz(i)
		fmt.Printf("%d --> %v\n", i, result)
		
		//fizzbuzz_print(i)
	}
} 


func fizzbuzz(x int) string {
	var s string
	if x%15 == 0 {
		s = "FizzBuzz"
	} else if x%3 == 0 {
		s = "Fizz"
	} else if x%5 == 0 {
		s = "Buzz"
	} else {
		s = strconv.Itoa(x)
	}

	return s
}
