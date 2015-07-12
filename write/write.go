// Copyright 2015, Xavier Henner
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
The write tool write new tags to a file
*/
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/xhenner/tag"
)

func init() {
}

func main() {
	flag.Parse()

	if flag.NArg() != 1 {
		fmt.Printf("usage: %v filename\n", os.Args[0])
		return
	}

	f, err := os.OpenFile(flag.Arg(0), os.O_RDWR, 0666)
	if err != nil {
		fmt.Printf("error loading file: %v", err)
		return
	}
	defer f.Close()

	w, err := tag.WriteTo(f)
	if err != nil {
		fmt.Printf("error reading file: %v\n", err)
		return
	}
	if err := w.Title("totoéé"); err != nil {
		fmt.Println(err)
	}
	w.Title("totoéé")
	w.Artist("Air")
	if err := w.Commit(); err != nil {
		fmt.Println(err)
	}
}
