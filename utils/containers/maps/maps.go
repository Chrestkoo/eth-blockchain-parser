package maps

type T interface {
	comparable
}

type Inf[K T, V any] interface {
	new() Inf[K, V]
	Set(k K, v V)
	Get(k K) (V, bool)
	Delete(k K)
	Len() int
	Range(fn func(k K, v V) error) (err error)
	String() string
}

// New 创建指定的map
func New[K T, V any](i Inf[K, V]) Inf[K, V] {
	return i.new()
}

// Values 获取map的值并返回数组
func Values[K T, V any](m map[K]V) []V {
	vs := make([]V, len(m))
	var i = 0
	for k := range m {
		vs[i] = m[k]
		i++
	}
	return vs
}

// Iterator 构建一个包含 map 中所有值的切片，并对每个值应用给定的函数。
//
// 参数：
//   - m：要迭代的 map
//   - fn：应用于每个值的函数，参数为值的指针，返回一个布尔值
//
// 返回值：
//   - []V：包含 map 中所有值的切片
//
// 示例：
//
//	m := map[string]int{"foo": 1, "bar": 2, "baz": 3}
//	var fn = func(v *int) bool {
//	    *v = *v + 1
//	    return true
//	}
//	result := Iterator(m, fn)
//	fmt.Println(result)  // 输出：[2 3 4]
//
// 特定说明：
//   - 函数参数 fn 会修改原始 map 中值的内容
//   - 如果函数 fn 返回 false，Iterator 函数将不会添加当前值到结果切片中
//   - 如果参数 m 为 nil 或空 map，则返回一个空切片
//   - 如果 map 中的值为 nil，则在传递给参数 fn 时，其指针 *V 会为 nil
func Iterator[K T, V any](m map[K]V, fn func(v *V) bool) []V {
	var vs []V
	for k := range m {
		v := m[k]
		if !fn(&v) {
			continue
		}
		vs = append(vs, v)
	}
	return vs
}
