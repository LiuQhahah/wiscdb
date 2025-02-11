package y

type uint64Heap []uint64

func (u uint64Heap) Len() int {
	return 0

}

func (u uint64Heap) Less(i, j int) bool {
	return false
}

func (u uint64Heap) Swap(i, j int) {

}
func (u *uint64Heap) Push(x interface{}) {

}

func (u *uint64Heap) Pop() interface{} {

	return nil
}
