package main

import (
	"context"
	"fmt"

	"github.com/libdns/westcn"
)

func main() {
	p := westcn.Provider{
		Username:    "YOUR_USERNAME",
		APIPassword: "YOUR_API_PASSWORD",
	}

	ret, err := p.GetRecords(context.TODO(), "your-domain")

	fmt.Println("Result:", ret)
	fmt.Println("Error:", err)
}
