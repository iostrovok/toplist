package toplist_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/iostrovok/toplist"
)

type TestSuite struct {
	suite.Suite
}

func Test(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (s *TestSuite) TestNew() {
	tl := toplist.New()
	s.NotNil(tl)
	tl.PrintList()

	err := tl.Save(1, 1)
	tl.PrintList()
	s.Nil(err)

	err = tl.Save(2, 2)
	tl.PrintList()
	s.Nil(err)

	err = tl.Delete(2)
	tl.PrintList()
	s.Nil(err)
}

//func (s *TestSuite) TestLongTest() {
//	tl := toplist.New()
//	s.NotNil(tl)
//	tl.PrintList()
//
//	start := make(chan struct{})
//	resultHash := &sync.Map{}
//	for i := range 10 {
//
//	}
//
//	close(start)
//
//	err := tl.Save(1, 1)
//	tl.PrintList()
//	s.Nil(err)
//
//	err = tl.Save(2, 2)
//	tl.PrintList()
//	s.Nil(err)
//
//	err = tl.Delete(2)
//	tl.PrintList()
//	s.Nil(err)
//}
//
//func (s *TestSuite) OneSave(tl *toplist.List, data []int64, resultHash *sync.Map) {
//	for
//}
