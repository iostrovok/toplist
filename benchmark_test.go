package toplist_test

import (
	"math/rand"

	"github.com/iostrovok/toplist"
)

type one struct {
	Action toplist.Action
	Index  int64
}

func makeData(size int, actions []toplist.Action) []one {
	out := make([]one, size)
	ln := len(actions)
	for i := range size {
		out[i].Action = actions[rand.Intn(ln)]
		out[i].Index = int64(i)
	}

	rand.Shuffle(len(out), func(i, j int) {
		out[i], out[j] = out[j], out[i]
	})

	return out
}

//func BenchmarkCAS(b *testing.B) {
//	initialVal := int64(0)
//	initialValRef := &initialVal
//
//	tl := toplist.New()
//
//	b.SetParallelism(5)
//	b.ResetTimer()
//	var compares atomic.Int64
//	b.RunParallel(func(pb *testing.PB) {
//		k := atomic.AddInt64(initialValRef, 1)
//
//		counter := int64(0)
//		counterRef := &counter
//
//		for pb.Next() {
//			compares.Add(1)
//
//			if err := tl.Save(k, k); err != nil {
//				b.Errorf("Save error for count %d: %s", k, err.Error())
//			}
//			if err := tl.Delete(k); err != nil {
//				b.Errorf("Delete error for count %d: %s", k, err.Error())
//			}
//
//			if atomic.AddInt64(counterRef, 1) > 10 {
//				return
//			}
//		}
//	})
//
//}

//func BenchmarkAll(b *testing.B) {
//	benchmarks := []struct {
//		name   string
//		pools  int
//		size   int
//		action []toplist.Action
//	}{
//		{"1000-5-Save", 5, 1000, []toplist.Action{toplist.SaveAction}},
//		{"1000-5-Insert", 5, 1000, []toplist.Action{toplist.InsertAction}},
//		{"1000-5-All", 5, 1000, []toplist.Action{toplist.SaveAction, toplist.InsertAction, toplist.DeleteAction}},
//
//		{"10_000-5-Save", 5, 10_000, []toplist.Action{toplist.SaveAction}},
//		{"10_000-5-Insert", 5, 10_000, []toplist.Action{toplist.InsertAction}},
//		{"10_000-5-All", 5, 10_000, []toplist.Action{toplist.SaveAction, toplist.InsertAction, toplist.DeleteAction}},
//
//		{"20_000-5-All", 5, 20_000, []toplist.Action{toplist.SaveAction, toplist.InsertAction, toplist.DeleteAction}},
//		{"20_000-10-All", 10, 20_000, []toplist.Action{toplist.SaveAction, toplist.InsertAction, toplist.DeleteAction}},
//
//		{"40_000-10-All", 10, 40_000, []toplist.Action{toplist.SaveAction, toplist.InsertAction, toplist.DeleteAction}},
//	}
//
//	tl := toplist.New()
//
//	b.ResetTimer()
//	//b.StopTimer()
//	for _, bm := range benchmarks {
//		makeData := makeData(bm.size, bm.action)
//		rand.Shuffle(len(makeData), func(i, j int) {
//			makeData[i], makeData[j] = makeData[j], makeData[i]
//		})
//
//		b.Run(bm.name, func(b *testing.B) {
//			for b.Loop() {
//				wg := &sync.WaitGroup{}
//				startCh := make(chan struct{}, 1)
//				for _ = range bm.pools {
//					wg.Go(func() {
//						<-startCh
//
//						for _, item := range makeData {
//							tl.Run(item.Action, item.Index, item.Index, nil)
//						}
//					})
//				}
//
//				b.StartTimer()
//				close(startCh)
//				wg.Wait()
//				b.StopTimer()
//			}
//		})
//	}
//}
//
//func BenchmarkFind(b *testing.B) {
//	b.Skip("ddd")
//	benchmarks := []struct {
//		name string
//		size int
//	}{
//		//{"1000", 1000},
//		{"10_000", 10_000},
//		//{"50_000", 50_000},
//	}
//
//	comparesTotal := int64(0)
//	NanosecondsTotal := int64(0)
//
//	for _, bm := range benchmarks {
//		fmt.Printf("\n")
//
//		all := make([]int64, bm.size)
//		tl := toplist.New()
//		for i := range bm.size {
//			all[i] = int64(i)
//
//			if err := tl.Save(int64(i), int64(i)); err != nil {
//				fmt.Printf("Insert error for %d, ERROR: %+v\n", i, err)
//			}
//
//			if i%10000 == 0 {
//				fmt.Printf("Insert %d...", i)
//			}
//		}
//
//		tl.Build()
//		fmt.Printf("\nBuild is done\n")
//
//		rand.Shuffle(len(all), func(i, j int) {
//			all[i], all[j] = all[j], all[i]
//		})
//
//		b.SetParallelism(5)
//		b.ResetTimer()
//
//		compares := atomic.Int64{}
//
//		b.RunParallel(func(pb *testing.PB) {
//			step := rand.Intn(100)
//			i := 0
//			for pb.Next() {
//				for range 10_000 {
//					point := all[i]
//					compares.Add(1)
//
//					res, find := tl.Find(point)
//					if !find {
//						fmt.Printf("point: %d, res: %d, find: %t\n", point, res, find)
//					} else if res.Index != point {
//						fmt.Printf("point: %d, res: %d, find: %t\n", point, res, find)
//					}
//
//					i += step
//					if bm.size <= i {
//						i = 0
//					}
//				}
//			}
//		})
//
//		nanoseconds := b.Elapsed().Nanoseconds()
//		fmt.Printf("(compares.Load(): %d, Nanoseconds: %d [Microsecond: %d, Millisecond: %d], b.N: %d\n",
//			compares.Load(), nanoseconds, b.Elapsed().Truncate(time.Microsecond), b.Elapsed().Truncate(time.Millisecond), b.N)
//
//		comparesTotal += compares.Load()
//		NanosecondsTotal += nanoseconds
//	}
//
//	fmt.Printf("%d / %d => %f compares/op\n", comparesTotal, NanosecondsTotal, float64(comparesTotal)/float64(NanosecondsTotal))
//	fmt.Printf("%d / %d => %f op/compares\n\n", NanosecondsTotal, comparesTotal, float64(NanosecondsTotal)/float64(comparesTotal))
//}
