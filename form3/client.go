package form3

import "fmt"

type Client struct{}

func (c Client) HelloWorld() {
	fmt.Println("Hello World!")
}
