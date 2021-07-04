package main

import (
	"fmt"
	"time"
)

func main() {
	nowSec := time.Now().Unix()

	fmt.Println(nowSec)
	nowConverted := time.Unix(nowSec, 0)
	fmt.Println(nowConverted.String())
}
