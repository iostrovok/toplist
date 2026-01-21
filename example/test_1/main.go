package main

import (
	"fmt"
	"strings"

	"github.com/iostrovok/toplist"
)

func main() {
	tl := toplist.New()

	fmt.Printf("\nEmpty list\n%s\n", strings.Join(tl.ToStringList(), " - "))

	err := tl.BaseSave(1, 1)
	fmt.Printf("BaseSave(1, 1) erros: %+v\n", err)
	fmt.Printf("\nSave(1, 1) list\n%s\n", strings.Join(tl.ToStringList(), " - "))

	err = tl.Save(2, 2)
	fmt.Printf("Save(2, 2) erros: %+v\n", err)
	fmt.Printf("\nSave(2, 2) list\n%s\n", strings.Join(tl.ToStringList(), " - "))
	//
	err = tl.Delete(2)
	fmt.Printf("Delete(2) erros: %+v\n", err)
	fmt.Printf("\nDelete(2) list\n%s\n", strings.Join(tl.ToStringList(), " - "))
}
