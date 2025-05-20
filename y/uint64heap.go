package y

type uint64Heap []uint64

// 最小堆操作
func (u uint64Heap) Len() int {
	return len(uint64Heap{})
}

func (u uint64Heap) Less(i, j int) bool {
	return u[i] < u[j]
}

func (u uint64Heap) Swap(i, j int) {
	u[i], u[j] = u[j], u[i]
}
func (u *uint64Heap) Push(x interface{}) {
	*u = append(*u, x.(uint64))
}

// 返回最上面的元素并且更新下当前u的值
func (u *uint64Heap) Pop() interface{} {
	old := *u
	n := len(old)
	x := old[n-1]
	*u = old[0 : n-1]
	return x
}
