package main

import (
	"fmt"
	"github.com/iostrovok/toplist"
	"math/rand"
	"time"
)

const (
	size = int64(50_000)
)

func main() {
	data := make([]int64, size)
	for i := range size {
		data[i] = i
	}

	rand.Shuffle(len(data), func(i, j int) {
		data[i], data[j] = data[j], data[i]
	})

	fmt.Printf("Start\n")
	start := time.Now()

	tl := toplist.New()
	localStart := time.Now()
	for _, i := range data {
		if err := tl.Insert(i, i); err != nil {
			panic(err)
		}
	}

	fmt.Printf("Insert done %s, %d ns\n", time.Since(start), int64(time.Since(localStart).Truncate(time.Nanosecond))/size)

	rand.Shuffle(len(data), func(i, j int) {
		data[i], data[j] = data[j], data[i]
	})

	localStart = time.Now()
	for _, i := range data {
		if err := tl.Delete(i); err != nil {
			panic(err)
		}
	}
	fmt.Printf("Delete done %s, %d ns\n", time.Since(start), int64(time.Since(localStart).Truncate(time.Nanosecond))/size)

	rand.Shuffle(len(data), func(i, j int) {
		data[i], data[j] = data[j], data[i]
	})

	localStart = time.Now()
	for _, i := range data {
		_, find := tl.Find(i)
		if find {
			panic("found")
		}
	}

	fmt.Printf("Find done %s, %d ns\n", time.Since(start), int64(time.Since(localStart).Truncate(time.Nanosecond))/size)

	rand.Shuffle(len(data), func(i, j int) {
		data[i], data[j] = data[j], data[i]
	})

	localStart = time.Now()
	for _, i := range data {
		if err := tl.Save(i, i); err != nil {
			panic(err)
		}
	}
	fmt.Printf("Save done %s, %d ns\n", time.Since(start), int64(time.Since(localStart).Truncate(time.Nanosecond))/size)

	rand.Shuffle(len(data), func(i, j int) {
		data[i], data[j] = data[j], data[i]
	})

	localStart = time.Now()
	for _, i := range data {
		if i%3 == 0 {
			if err := tl.Delete(i); err != nil {
				panic(err)
			}
		}
	}
	fmt.Printf("Delete done %s, %d ns\n", time.Since(start), int64(time.Since(localStart).Truncate(time.Nanosecond))/size)

	rand.Shuffle(len(data), func(i, j int) {
		data[i], data[j] = data[j], data[i]
	})

	tl.Build()
	localStart = time.Now()
	for n, i := range data {
		_, find := tl.Find(i)
		if 0 == i%3 && find {
			panic("found")
		}

		if 0 != i%3 && !find {
			panic("not found")
		}

		if n > 10_000 {
			break
		}
	}

	fmt.Printf("Find done %s, %d ns\n", time.Since(start), int64(time.Since(localStart).Truncate(time.Nanosecond))/10_000)

	fmt.Printf("All done\n")
}
