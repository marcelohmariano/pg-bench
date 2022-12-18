package bench

import (
	"container/heap"
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
	return x
}

func (h *floatHeap) At(index int) float64 {
	hp := *h
	return hp[index]
}

func (h *floatHeap) Peek() float64 {
	return h.At(0)
}

type medianCalculator struct {
	left  floatHeap
	right floatHeap

	value float64
}

func (mc *medianCalculator) Calculate(values ...float64) float64 {
	for _, v := range values {
		mc.insert(v)
		mc.balanceHeaps()
		mc.updateValue()
	}
	return mc.value
}

func (mc *medianCalculator) Value() float64 {
	return mc.value
}

func (mc *medianCalculator) insert(value float64) {
	switch {
	case mc.right.Len() == 0 || value < mc.right.Peek():
		heap.Push(&mc.left, -value)
	default:
		heap.Push(&mc.right, value)
	}
}

func (mc *medianCalculator) balanceHeaps() {
	if mc.left.Len() < mc.right.Len() {
		x := heap.Pop(&mc.right).(float64)
		heap.Push(&mc.left, -x)
	}

	if mc.left.Len() > mc.right.Len()+1 {
		x := heap.Pop(&mc.left).(float64)
		heap.Push(&mc.right, -x)
	}
}

func (mc *medianCalculator) updateValue() {
	v := -mc.left.Peek()

	if mc.left.Len() == mc.right.Len() {
		mc.value = (v + mc.right.Peek()) / 2
		return
	}

	mc.value = v
}
