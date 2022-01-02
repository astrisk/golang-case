package main

import (
	"fmt"
	"github.com/golang-case/week04/pkg"
	"time"
)

func main() {

	count := pkg.NewCount()
	for _, x := range []float64{0.5, 0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.1} {
		count.Add(x)
		fmt.Println("add count:", x)
		time.Sleep(time.Second)
	}

	result1 := count.Sum(time.Now())
	fmt.Println(result1)

	time.Sleep(10 * time.Second)
	//count.Add(1.0)
	fmt.Println(count.Sum(time.Now()))

}
