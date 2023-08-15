package chainservice

import (
	"github.com/ethereum/go-ethereum/core/types"
)

type EventQueue []types.Log

func (q EventQueue) Len() int { return len(q) }
func (q EventQueue) Less(i, j int) bool {
	if q[i].BlockNumber == q[j].BlockNumber {
		return i < j
	}
	return q[i].BlockNumber < q[j].BlockNumber
}

func (q EventQueue) Swap(i, j int) { q[i], q[j] = q[j], q[i] }

func (q *EventQueue) Push(x interface{}) {
	*q = append(*q, x.(types.Log))
}

func (q *EventQueue) Pop() interface{} {
	old := *q
	n := len(old)
	x := old[n-1]
	*q = old[0 : n-1]
	return x
}
