package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/umemak/asin2md"
)

func main() {
	flag.Parse()
	err := run(flag.Arg(0))
	if err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}
}

func run(asin string) error {
	res, err := asin2md.Get(asin)
	if err != nil {
		return fmt.Errorf("getamazoninfo.Get: %w", err)
	}
	fmt.Println(res)
	return nil
}
