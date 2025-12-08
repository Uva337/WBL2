package main

import (
	"fmt"
	"github.com/beevik/ntp"
	"os"
)

func main() {
	time, err := ntp.Time("0.beevik-ntp.pool.ntp.org")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	} else {
		fmt.Println(time)

	}

}
