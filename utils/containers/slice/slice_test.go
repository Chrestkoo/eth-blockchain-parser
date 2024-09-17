package slice_test

import (
	"fmt"
	"github.com/eth-blockchain-parser/utils/containers/slice"
	"testing"
)

func TestMutex(t *testing.T) {
	m := slice.New[int](new(slice.Mutex[int]))
	m.Push(1)
	fmt.Println(m.Index(0))
	m.Push(2)
	m.Push(3)
	m.Delete(3)
	fmt.Println(m.Len())
	err := m.Range(func(v int) error {
		fmt.Println("V=>", v)
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
	t.Log(m)
}

func TestRW(t *testing.T) {
	m := slice.New[int](new(slice.RWMutex[int]))
	for i := 0; i < 10; i++ {
		m.Push(i)
	}
	fmt.Println(m.Index(0))
	err := m.Range(func(v int) error {
		fmt.Println("a: ", v)
		return nil
	})
	m.Delete(5)
	m.Delete(5)
	fmt.Println("========")
	err = m.Range(func(v int) error {
		fmt.Println("b: ", v)
		return nil
	})
	m.Delete(5)
	m.Delete(4)
	fmt.Println("========")
	err = m.Range(func(v int) error {
		fmt.Println("c: ", v)
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
	t.Log(m)
}
