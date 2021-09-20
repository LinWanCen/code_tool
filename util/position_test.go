package util

import "fmt"

func ExamplePosition_t0() {
	fmt.Println(Position(`12
345
67`, 0))

	// Output:
	// 12
	// 345
	// 67
}

func ExamplePosition_t1() {
	fmt.Println(Position(`12
345
67`, 1))

	// Output:
	// 12
	// ^
	// 345
	// 67
}

func ExamplePosition_t3() {
	fmt.Println(Position(`12
345
67`, 3))

	// Output:
	// 12
	//   ^
	// 345
	// 67
}

func ExamplePosition_t4() {
	fmt.Println(Position(`12
345
67`, 4))

	// Output:
	// 12
	// 345
	// ^
	// 67
}

func ExamplePosition_t8() {
	fmt.Println(Position(`12
345
67`, 8))

	// Output:
	// 12
	// 345
	// 67
	// ^
}

func ExamplePosition_t9() {
	fmt.Println(Position(`12
345
67`, 9))

	// Output:
	// 12
	// 345
	// 67
	//   ^
}

func ExamplePosition_t10() {
	fmt.Println(Position(`12
345
67`, 10))

	// Output:
	// 12
	// 345
	// 67
}
