package cache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v any) *ListItem
	PushBack(v any) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value any
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	frontItem *ListItem // Первый элемент двусвязного списка.
	backItem  *ListItem // Последний элемент двусвязного списка.
	len       int       // Количество элементов в двусвязном списке.
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.frontItem
}

func (l *list) Back() *ListItem {
	return l.backItem
}

func (l *list) PushFront(v any) *ListItem {
	newFrontItem := &ListItem{Value: v}

	if l.frontItem == nil {
		l.frontItem = newFrontItem
		l.backItem = newFrontItem
	} else {
		newFrontItem.Next = l.Front()
		l.frontItem.Prev = newFrontItem

		l.frontItem = newFrontItem
	}

	l.len++

	return newFrontItem
}

func (l *list) PushBack(v any) *ListItem {
	newBackItem := &ListItem{Value: v}

	if l.backItem == nil {
		l.frontItem = newBackItem
		l.backItem = newBackItem
	} else {
		newBackItem.Prev = l.Back()
		l.backItem.Next = newBackItem

		l.backItem = newBackItem
	}

	l.len++

	return newBackItem
}

func (l *list) Remove(i *ListItem) {
	switch {
	case i.Prev == nil && i.Next == nil:
		l.frontItem = nil
		l.backItem = nil
	case i.Prev == nil:
		i.Next.Prev = nil
		l.frontItem = i.Next
	case i.Next == nil:
		i.Prev.Next = nil
		l.backItem = i.Prev
	default:
		i.Prev.Next = i.Next
		i.Next.Prev = i.Prev
	}

	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	l.PushFront(i.Value)
	l.Remove(i)
}

func NewList() List {
	return new(list)
}
