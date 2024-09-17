package ttl_test

import (
	"fmt"
	"github.com/eth-blockchain-parser/utils/cache/ttl"
	"testing"
	"time"
)

func TestTTL(t *testing.T) {
	cache1 := ttl.New[string, string](time.Minute * 1)
	cache2 := ttl.New[string, string](time.Minute * 2)
	cache3 := ttl.New[string, string](time.Minute * 3)

	list := make(map[string]*ttl.Cache[string, string])

	addr1 := "0xasfafafs1"
	addr2 := "0xasfafafs2"
	addr3 := "0xasfafafs3"
	cache1.Put(addr1, time.Now().Format(time.DateTime))
	cache2.Put(addr2, time.Now().Add(1*time.Minute).Format(time.DateTime))
	cache3.Put(addr3, time.Now().Add(2*time.Minute).Format(time.DateTime))
	list[addr1] = cache1
	list[addr2] = cache2
	list[addr3] = cache3

	go timerDelete(list, addr1)
	go timerDelete(list, addr2)
	go timerDelete(list, addr3)

	fmt.Println("Print all")
	print(list, addr1, addr2, addr3)
	time.Sleep(time.Minute * 1)

	fmt.Println("result delete 1st")
	print(list, addr1, addr2, addr3)
	time.Sleep(time.Minute * 1)

	fmt.Println("result delete 2nd")
	print(list, addr1, addr2, addr3)
	time.Sleep(time.Minute * 1)

	fmt.Println("result delete all")
	print(list, addr1, addr2, addr3)
	/*time.Sleep(time.Minute * 2)

	fmt.Println("6 ============= 6")
	for _, i := range list {
		for _, v := range vList {
			rst, ok := i.Get(v)
			if !ok {
				fmt.Println(fmt.Sprintf("%v. get fail %v", i, v))
			}
			fmt.Println(fmt.Sprintf("%v. rst %v", i, rst))
		}
	}
	time.Sleep(time.Minute * 1)*/

	/*t.Log("set key 1")
	time.Sleep(time.Second * 2)
	t.Log("set key 2")
	cache.Put(2, "2")
	time.Sleep(time.Second * 3)
	t.Log("set key 3")
	cache.Put(3, "3")
	time.Sleep(time.Second * 4)
	t.Log("set key 4")
	cache.Put(4, "4")
	time.Sleep(time.Second * 5)
	t.Log("set key 5")
	cache.Put(5, "5")
	t.Log(cache.Get(1))
	time.Sleep(time.Second)
	t.Log(cache.Get(2))
	time.Sleep(time.Second)
	t.Log(cache.Get(3))
	time.Sleep(time.Second)
	t.Log(cache.Get(4))
	time.Sleep(time.Second)
	t.Log(cache.Get(5))*/
}

func print(list map[string]*ttl.Cache[string, string], addr ...string) {
	for k, i := range addr {
		rst, ok := list[i]
		if !ok {
			fmt.Println(fmt.Sprintf("%v. get fail %v", k, i))
		} else {
			rst2, ok2 := rst.Get(i)
			if !ok2 {
				fmt.Println(fmt.Sprintf("%v. get fail %v", k, i))
			} else {
				fmt.Println(fmt.Sprintf("%v. get success %v - %v", k, i, rst2))
			}
		}
	}
}
func deleteCache(list map[string]*ttl.Cache[string, string], addr string) bool {
	rst, ok := list[addr]
	if !ok {
		fmt.Println(fmt.Sprintf("no such addr %v", addr))
		return true
	}
	data, _ := rst.Get(addr)
	if data == "" {
		delete(list, addr)
		fmt.Println(fmt.Sprintf("delete done %v", addr))
		return true
	}
	return false
}
func timerDelete(list map[string]*ttl.Cache[string, string], addr string) {
	for {
		if deleteCache(list, addr) {
			break
		}
		time.Sleep(5 * time.Second)
	}
	fmt.Println(fmt.Sprintf("timer delete done %v", addr))
}
