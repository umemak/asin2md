package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/umemak/asin2md"
)

func main() {
	flag.Parse()
	err := run(flag.Args())
	if err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	var asin, buy, from, to string
	l := len(args)
	if l == 0 {
		return fmt.Errorf("no args error")
	}
	if l >= 1 {
		asin = args[0]
	}
	if l >= 2 {
		buy = args[1]
	}
	if l >= 3 {
		from = args[2]
	}
	if l >= 4 {
		to = args[3]
	}
	res, err := asin2md.Get(asin, buy, from, to)
	if err != nil {
		return fmt.Errorf("getamazoninfo.Get: %w", err)
	}
	// fmt.Println(res)
	os.WriteFile(asin+".md", []byte(res), 0666)
	return nil
}
