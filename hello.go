package main

import (
	"errors"
	"fmt"
	"math"
)

type person struct {
	name string
	age  int
}

func main() {
	// Variables
	var x int = 5
	y := 7
	var total int = x + y

	fmt.Println(total)

	// If statements
	if total > 6 {
		fmt.Println("More than 6")
	} else if total == 6 {
		fmt.Println("Equal to 6")
	} else {
		fmt.Println("Less than 6")
	}

	// Arrays
	var a [5]int
	b := [5]int{5, 4, 3, 2, 1}
	c := []int{5, 4, 3, 2, 1}
	a[2] = 7
	c = append(c, 13)
	fmt.Println(a)
	fmt.Println(b)
	fmt.Println(c)

	// Maps
	vertices := make(map[string]int)

	vertices["triangle"] = 2
	vertices["square"] = 3
	vertices["dodecagon"] = 12

	delete(vertices, "square")

	fmt.Println(vertices)
	fmt.Println(vertices["triangle"])

	// Loops
	for i := 0; i < 5; i++ {
		fmt.Println(i)
	}

	i := 0
	for i < 5 {
		fmt.Println(i)
		i++
	}

	arr := []string{"a", "b", "c"}
	for index, value := range arr {
		fmt.Println("index:", index, "value:", value)
	}

	m := make(map[string]string)
	m["a"] = "alpha"
	m["b"] = "beta"
	for key, value := range arr {
		fmt.Println("key:", key, "value:", value)
	}

	// Functions
	result := sum(2, 3)
	fmt.Println(result)

	result2, err := sqrt(16)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(result2)
	}

	// Structs
	p := person{name: "Jake", age: 23}
	fmt.Println(p)
	fmt.Println(p.age)

	// Pointers
	num := 7
	inc(&num)
	fmt.Println(num)
	fmt.Println(&num)
}

func sum(x int, y int) int {
	return x + y
}

func sqrt(x float64) (float64, error) {
	if x < 0 {
		return 0, errors.New("undefined for negative numbers")
	}

	return math.Sqrt(x), nil
}

func inc(x *int) {
	*x++
}
