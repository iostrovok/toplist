package main

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/iostrovok/toplist"
	"github.com/iostrovok/toplist/internal/queue"
)

const (
	checkCount = 50_000
	goCount    = 5
)

type one struct {
	Action queue.Action
	Index  int64
}

func main() {
	tl := toplist.New()
	start := make(chan struct{})
	gw := sync.WaitGroup{}

	allData := makeData()
	for i := range goCount {
		data := make([]one, len(allData))
		copy(data, allData)
		rand.Shuffle(len(data), func(i, j int) {
			data[i], data[j] = data[j], data[i]
		})

		gw.Add(1)
		go func(i int) {
			defer gw.Done()
			<-start

			for i := range data {
				if data[i].Action == queue.DeleteAction {
					_ = tl.Delete(data[i].Index)
				} else if data[i].Action == queue.SaveAction {
					if err := tl.Save(data[i].Index, data[i].Index); err != nil {
						panic(err)
					}
				} else if data[i].Action == queue.InsertAction {
					if err := tl.Insert(data[i].Index, data[i].Index); err != nil {
						panic(err)
					}
				}
			}
		}(i)
	}

	// waiting for all goroutines are ready
	time.Sleep(100 * time.Millisecond)

	close(start)

	gw.Wait()

	time.Sleep(100 * time.Millisecond)

	//fmt.Printf("\nDebug out - 0:\n")
	//for i, s := range tl.ToStringList() {
	//	fmt.Printf("Level [%d]: %s\n", i, s)
	//}

	deleted := tl.Clean()
	fmt.Printf("\ndeleted: %d\n", deleted)

	tl.Build()

	fmt.Printf("\nDebug out:\n")
	//for i, s := range tl.ToStringList() {
	//	fmt.Printf("Level [%d]: %s\n", i, s)
	//}

	resultHash := tl.DebugMap()

	keys := make([]int64, 0, len(resultHash))
	for id := range resultHash {
		keys = append(keys, id)
	}

	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })

	for i, id := range keys {
		action := resultHash[id]
		if i%1000 == 0 {
			fmt.Printf("%d: %s, 	", id, action.String())
		}

		//if action != queue.DeleteAction {
		//	continue
		//}

		res, find := tl.Find(id)
		if !find && action == queue.DeleteAction {
			continue
		}
		if find && (action == queue.SaveAction || action == queue.InsertAction) {
			continue
		}

		fmt.Printf("\nBad result: id: %v, action: %v, res: %v, find: %v\n", id, action.String(), res, find)
		//break
	}

	fmt.Printf("\n\n")

	debugStr, debugIds := tl.CheckBase()

	fmt.Printf("%s\n\n", strings.Join(debugStr, "\n"))
	for _, id := range debugIds {
		fmt.Printf("--- debugHash[%d]:\n%v\n", id, resultHash[id])
	}

	fmt.Printf("\n\nb.CountRun:::: %d\n\n", tl.Queue.Build.CountRun)

	duplicates := tl.CheckDuplicate()
	fmt.Printf("\n\nb.CheckDuplicate from %d::::\n\n", len(duplicates))

	for id, count := range duplicates {
		if count > 1 {
			fmt.Printf("--- duplicates[%d] - %v\n", id, count)
		}
	}
}

func makeData() []one {
	out := make([]one, checkCount)
	for i := range out {
		out[i].Action = queue.SaveAction
		out[i].Index = int64(i)
		switch t := rand.Intn(3); t {
		case 1:
			out[i].Action = queue.SaveAction
		case 2:
			out[i].Action = queue.InsertAction
		case 0:
			out[i].Action = queue.DeleteAction
		}
	}

	return out
}
