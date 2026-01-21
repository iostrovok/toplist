package main

import (
	"fmt"
	"github.com/iostrovok/toplist"
	"math/rand"
)

const (
	size = int64(10_000)
)

func main() {
	fmt.Println("Start")

	data := make([]int64, size)
	for i := range size {
		data[i] = i
	}

	rand.Shuffle(len(data), func(i, j int) {
		data[i], data[j] = data[j], data[i]
	})

	tl := toplist.New()
	for _, i := range data {
		if err := tl.Insert(i, i); err != nil {
			panic(err)
		}
	}

	rand.Shuffle(len(data), func(i, j int) {
		data[i], data[j] = data[j], data[i]
	})

	for _, i := range data {
		if err := tl.Delete(i); err != nil {
			panic(err)
		}
	}

	rand.Shuffle(len(data), func(i, j int) {
		data[i], data[j] = data[j], data[i]
	})

	for _, i := range data {
		_, find := tl.Find(i)
		if find {
			panic("found")
		}
	}

	rand.Shuffle(len(data), func(i, j int) {
		data[i], data[j] = data[j], data[i]
	})

	for _, i := range data {
		if err := tl.Save(i, i); err != nil {
			panic(err)
		}
	}

	rand.Shuffle(len(data), func(i, j int) {
		data[i], data[j] = data[j], data[i]
	})

	for _, i := range data {
		if i%3 == 0 {
			if err := tl.Delete(i); err != nil {
				panic(err)
			}
		}
	}

	rand.Shuffle(len(data), func(i, j int) {
		data[i], data[j] = data[j], data[i]
	})

	for _, i := range data {
		_, find := tl.Find(i)
		if 0 == i%3 && find {
			panic("found")
		}

		if 0 != i%3 && !find {
			panic("not found")
		}
	}

	fmt.Println("Done")
}
