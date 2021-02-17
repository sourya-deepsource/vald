package rand

import (
	"strconv"
	"testing"
)

func Benchmark_Uint32(b *testing.B) {
	type field struct {
	}
	type args struct {
	}
	type result struct {
		res1 uint32
	}
	type test struct {
		name       string
		field      field
		args       args
		result     result
		beforeFunc func()
		afterFunc  func()
	}

	tests := []test{
		{
			name: "test rand",
		},
	}
	for _, test := range tests {
		b.Run(test.name, func(b *testing.B) {
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}

			for i := 0; i < b.N; i++ {
				test.result.res1 = Uint32()
			}
		})
	}
	b.ResetTimer()

	type paralleTest struct {
		name       string
		field      field
		args       args
		result     result
		paralle    []int
		beforeFunc func()
		afterFunc  func()
	}
	ptests := []paralleTest{
		{
			name:    "test rand",
			paralle: []int{1, 2, 4, 6, 8, 16},
		},
	}
	for _, ptest := range ptests {
		test := ptest
		for _, p := range test.paralle {
			name := test.name + "-" + strconv.Itoa(p)
			b.Run(name, func(b *testing.B) {
				b.SetParallelism(p)
				b.ResetTimer()

				if test.beforeFunc != nil {
					test.beforeFunc()
				}
				if test.afterFunc != nil {
					defer test.afterFunc()
				}
				b.RunParallel(func(pb *testing.PB) {
					for pb.Next() {
						test.result.res1 = Uint32()
					}
				})
			})
		}
	}
}
