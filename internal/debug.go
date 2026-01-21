package internal

import (
	"unsafe"
)

func ToListLevel(head unsafe.Pointer) []*Node {
	out := make([]*Node, 0)

	preHead := loadUnmarkPointer(head)
	if preHead == nil {
		// false start
		return out
	}

	out = append(out, preHead)
	currPtr := preHead.next

	for {
		currHead := loadUnmarkPointer(currPtr)
		out = append(out, currHead)

		if currHead == nil || currHead.next == nil {
			return out
		}

		currPtr = currHead.next
	}
}

func ToList(head unsafe.Pointer) [][]*Node {
	out := make([][]*Node, 0)
	for head != nil {
		out = append(out, ToListLevel(head))
		currHead := loadUnmarkPointer(head)
		if currHead == nil {
			break
		}

		head = currHead.down
	}

	return out
}

func CheckDuplicate(head unsafe.Pointer) map[int64]int {
	out := map[int64]int{}

	head = GetBaseLavelHead(head)
	if head == nil {
		// false start
		return out
	}

	preHead := loadUnmarkPointer(head)
	if preHead == nil {
		// false start
		return out
	}

	out[preHead.Index] = 1
	currPtr := preHead.next

	for {
		currHead := loadUnmarkPointer(currPtr)
		out[currHead.Index] = out[currHead.Index] + 1

		if currHead == nil || currHead.next == nil {
			return out
		}

		currPtr = currHead.next
	}
}

func ToStringListLevel(head unsafe.Pointer) []string {
	out := make([]string, 0)
	preHead := loadUnmarkPointer(head)
	if preHead == nil {
		// false start
		return out
	}

	out = append(out, preHead.String())
	currPtr := preHead.next

	for {
		currHead := loadUnmarkPointer(currPtr)
		out = append(out, currHead.String())

		if currHead == nil || currHead.next == nil {
			return out
		}

		currPtr = currHead.next
	}
}

func ToStringList(head unsafe.Pointer) [][]string {
	out := make([][]string, 0)

	for head != nil {
		out = append(out, ToStringListLevel(head))
		currHead := loadUnmarkPointer(head)
		if currHead == nil {
			break
		}
		head = currHead.down
	}

	return out
}

func CheckBase(head unsafe.Pointer) ([]string, []int64) {
	baseNode := GetBaseLavelHead(head)
	return CheckBaseLevel(baseNode)
}

func CheckBaseLevel(head unsafe.Pointer) ([]string, []int64) {
	out := make([]string, 0)
	ids := make([]int64, 0)
	predHead := loadUnmarkPointer(head)
	head = predHead.next

	for head != nil {
		currHead := loadUnmarkPointer(head)
		if currHead == nil {
			break
		}

		if predHead.Index > currHead.Index {
			out = append(out, predHead.String()+" --> "+currHead.String())
			ids = append(ids, predHead.Index)
			ids = append(ids, currHead.Index)
		}

		predHead = currHead
		head = currHead.next
	}

	return out, ids
}
