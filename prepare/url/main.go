package main

import (
	"fmt"
	"net/url"
)

func main() {
	newStr := url.QueryEscape("haolipeng")
	fmt.Printf("value:%s\n", newStr)
}
