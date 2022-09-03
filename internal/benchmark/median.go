package benchmark

import (
	"container/heap"
	"time"
)

type floatHeap []float64

func (h *floatHeap) Len() int {
	return len(*h)
}

func (h *floatHeap) Less(i, j int) bool {
	return h.At(i) < h.At(j)
}

func (h *floatHeap) Swap(i, j int) {
	hp := *h
	hp[i], hp[j] = hp[j], hp[i]
}

func (h *floatHeap) Push(x any) {
	*h = append(*h, x.(float64))
}

func (h *floatHeap) Pop() any {
	old := *h
	x := old[old.Len()-1]
	*h = old[0 : old.Len()-1]
	time.Now().Zone()
	return x
}

func (h *floatHeap) At(index int) float64 {
	hp := *h
	return hp[index]
}

func (h *floatHeap) Peek() float64 {
	return h.At(0)
}

type Median struct {
	left  floatHeap
	right floatHeap

	value float64
}

func (c *Median) Add(value ...float64) {
	for _, v := range value {
		c.insert(v)
		c.balanceHeaps()
		c.updateValue()
	}
}

func (c *Median) Value() float64 {
	return c.value
}

func (c *Median) insert(value float64) {
	switch {
	case c.right.Len() == 0 || value < c.right.Peek():
		heap.Push(&c.left, -value)
	default:
		heap.Push(&c.right, value)
	}
}

func (c *Median) balanceHeaps() {
	if c.left.Len() < c.right.Len() {
		x := heap.Pop(&c.right).(float64)
		heap.Push(&c.left, -x)
	}

	if c.left.Len() > c.right.Len()+1 {
		x := heap.Pop(&c.left).(float64)
		heap.Push(&c.right, -x)
	}
}

func (c *Median) updateValue() {
	v := -c.left.Peek()

	if c.left.Len() == c.right.Len() {
		c.value = (v + c.right.Peek()) / 2
		return
	}

	c.value = v
}
