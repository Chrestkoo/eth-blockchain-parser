package maps_test

import (
	"fmt"
	"github.com/eth-blockchain-parser/utils/containers/maps"
	"testing"
)

func TestMutex(t *testing.T) {
	m := maps.New[int, int](new(maps.Mutex[int, int]))
	m.Set(1, 2)
	fmt.Println(m.Get(1))
	m.Set(2, 3)
	m.Set(3, 4)
	m.Delete(1)
	fmt.Println(m.Get(1))
	err := m.Range(func(k int, v int) error {
		fmt.Println(k, "=>", v)
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
}

func TestRW(t *testing.T) {
	m := maps.New[int, int](new(maps.RWMutex[int, int]))
	m.Set(1, 2)
	fmt.Println(m.Get(1))
	m.Set(2, 3)
	m.Set(3, 4)
	m.Delete(1)
	fmt.Println(m.Get(1))
	err := m.Range(func(k int, v int) error {
		fmt.Println(k, "=>", v)
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
}

func TestShards(t *testing.T) {
	m := maps.New[int, int](new(maps.Shards[int, int]))
	count := 999999999
	for i := 0; i < count; i++ {
		m.Set(i, i)
	}
	//fmt.Println(m.Get(1))
	//m.Set(2, 3)
	//m.Set(3, 4)
	//m.Delete(1)
	//fmt.Println(m.Get(1))
	tmp := make(map[int]int)
	err := m.Range(func(k int, v int) error {
		//fmt.Println(k, "=>", v)
		tmp[k] = v
		return nil
	})
	fmt.Println(len(tmp))
	if err != nil {
		fmt.Println(err)
	}
}

func TestSortByVal(t *testing.T) {
	var m = map[string]int{"A": 1, "B": 2, "C": 3}
	fmt.Printf("%#v", maps.SortByVal(m))
}

func TestSortByValReverse(t *testing.T) {
	var m = map[string]int{"A": 1, "B": 2, "C": 3}
	fmt.Printf("%#v", maps.SortByValReverse(m))
}

func TestSortByKey(t *testing.T) {
	var m = map[string]int{"A": 1, "B": 2, "C": 3}
	fmt.Printf("%#v", maps.SortByKey(m))
}

func TestSortByKeyReverse(t *testing.T) {
	var m = map[string]int{"A": 1, "B": 2, "C": 3}
	fmt.Printf("%#v", maps.SortByKeyReverse(m))
}

func TestValues(t *testing.T) {
	t.Log(maps.Values(map[string]int{"A": 1, "B": 2, "C": 3}))
	type T struct {
		A string
		B int
		C float32
	}
	t.Log(maps.Values(map[string]T{"A": {
		A: "A",
		B: 123,
		C: 0.00007,
	}, "B": {
		A: "B",
		B: 456,
		C: 0.00014,
	}, "C": {
		A: "C",
		B: 789,
		C: 0.00028,
	}}))
}
