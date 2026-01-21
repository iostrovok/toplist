package queue

type Item struct {
	// type of item
	action Action

	// in
	index int64
	body  any

	// out
	ResultFunc ResultFunction // result
}

func (item *Item) Run(save, insert SaveFunc, del DeleteFunc) {
	switch item.action {
	case InsertAction:
		item.ResultFunc(InsertAction, item.index, insert(item.index, item.body))
	case SaveAction:
		item.ResultFunc(SaveAction, item.index, save(item.index, item.body))
	case DeleteAction:
		item.ResultFunc(DeleteAction, item.index, del(item.index))
	}
}
