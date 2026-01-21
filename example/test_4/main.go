package main

import (
	"fmt"
	"strings"

	"github.com/iostrovok/toplist"
)

func main() {
	tl := toplist.New()

	syncCh := make(chan struct{}, 0)
	resultFunc := func(action toplist.Action, index int64, err error) {
		fmt.Printf("%s(%d) erros: %+v\n", action, index, err)
		syncCh <- struct{}{}
	}

	fmt.Printf("\nEmpty list\n%s\n", strings.Join(tl.ToStringList(), " - "))

	tl.Run(toplist.SaveAction, 1, "just body FOR 1", resultFunc)
	<-syncCh
	fmt.Printf("\ntoplist.SaveAction list\n%s\n", strings.Join(tl.ToStringList(), " - "))

	tl.Run(toplist.InsertAction, 2, "just body FOR 2", resultFunc)
	<-syncCh
	fmt.Printf("\ntoplist.InsertAction list\n%s\n", strings.Join(tl.ToStringList(), " - "))

	tl.Run(toplist.DeleteAction, 2, nil, resultFunc)
	<-syncCh
	fmt.Printf("\ntoplist.DeleteAction list\n%s\n", strings.Join(tl.ToStringList(), " - "))
}
