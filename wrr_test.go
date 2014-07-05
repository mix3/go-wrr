package wrr

import (
	"reflect"
	"testing"
)

func expect(t *testing.T, a interface{}, b interface{}) {
	defer func() {
		if err := recover(); err != nil {
			if !reflect.DeepEqual(a, b) {
				t.Errorf("Expected %#v (type %v) - Got %#v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
			}
		}
	}()
	if a != b {
		t.Errorf("Expected %#v (type %v) - Got %#v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func expectRRList(t *testing.T, a DataSlice, b DataSlice) {
	expect(t, len(a), len(b))
	for i := 0; i < len(a); i++ {
		expect(t, *a[i], *b[i])
	}
}

func TestNew(t *testing.T) {
	func() {
		rr := New(DataSlice{})
		expect(t, rr.rrList, DataSlice{})
		expect(t, rr.weights, 0)
		expect(t, rr.defaultWeight, 100)
	}()
	func() {
		rr := New(DataSlice{
			&Data{Value: "foo"},
			&Data{Value: "bar"},
		})
		expectRRList(t, rr.rrList, DataSlice{
			&Data{Key: "foo", Value: "foo", Weight: 100, rng: 100},
			&Data{Key: "bar", Value: "bar", Weight: 100, rng: 0},
		})
		expect(t, rr.weights, 200)
		expect(t, rr.defaultWeight, 100)
	}()
	func() {
		rr := New(DataSlice{
			&Data{Value: "foo", Weight: 50},
			&Data{Value: "bar", Weight: 100},
		})
		expectRRList(t, rr.rrList, DataSlice{
			&Data{Key: "foo", Value: "foo", Weight: 50, rng: 100},
			&Data{Key: "bar", Value: "bar", Weight: 100, rng: 0},
		})
		expect(t, rr.weights, 150)
		expect(t, rr.defaultWeight, 100)
	}()
	func() {
		rr := New(DataSlice{
			&Data{Value: "foo"},
			&Data{Value: "bar"},
		}, Option{
			DefaultWeight: 20,
		})
		expectRRList(t, rr.rrList, DataSlice{
			&Data{Key: "foo", Value: "foo", Weight: 20, rng: 20},
			&Data{Key: "bar", Value: "bar", Weight: 20, rng: 0},
		})
		expect(t, rr.weights, 40)
		expect(t, rr.defaultWeight, 20)
	}()
	func() {
		DefaultWeightBak := DefaultWeight
		DefaultWeight = 20
		defer func() {
			DefaultWeight = DefaultWeightBak
		}()
		rr := New(DataSlice{
			&Data{Value: "foo"},
			&Data{Value: "bar"},
		})
		expectRRList(t, rr.rrList, DataSlice{
			&Data{Key: "foo", Value: "foo", Weight: 20, rng: 20},
			&Data{Key: "bar", Value: "bar", Weight: 20, rng: 0},
		})
		expect(t, rr.weights, 40)
		expect(t, rr.defaultWeight, 20)
	}()
}

type MockRand []int

func (mr *MockRand) Intn(n int) int {
	slice := *mr
	ret := slice[0]
	slice = slice[1:]
	*mr = slice
	return ret
}

func TestNext(t *testing.T) {
	for _, bTreeBorder := range []int{10, 0} {
		BTreeBorderBak := BTreeBorder
		BTreeBorder = bTreeBorder
		defer func() {
			BTreeBorder = BTreeBorderBak
		}()
		func() {
			rr := New(DataSlice{})
			expect(t, rr.Next(), nil)
			expect(t, rr.Next(), nil)
			expect(t, rr.Next(), nil)
		}()
		func() {
			rr := New(DataSlice{
				&Data{Value: "foo"},
			})
			expect(t, rr.Next(), "foo")
			expect(t, rr.Next(), "foo")
			expect(t, rr.Next(), "foo")
		}()
		func() {
			rr := New(DataSlice{
				&Data{Value: "foo"},
				&Data{Value: "bar"},
			})
			rr.rand = &MockRand{0, 100, 10, 200}
			expect(t, rr.Next(), "bar")
			expect(t, rr.Next(), "foo")
			expect(t, rr.Next(), "bar")
			expect(t, rr.Next(), "foo")
		}()
		func() {
			rr := New(DataSlice{
				&Data{Value: "foo", Weight: 50},
				&Data{Value: "bar", Weight: 100},
			})
			rr.rand = &MockRand{0, 100, 99, 150}
			expect(t, rr.Next(), "bar")
			expect(t, rr.Next(), "foo")
			expect(t, rr.Next(), "bar")
			expect(t, rr.Next(), "foo")
		}()
		func() {
			rr := New(DataSlice{
				&Data{Value: "foo", Weight: 50},
				&Data{Value: "bar", Weight: 100},
				&Data{Value: "baz", Weight: 20},
			})
			rr.rand = &MockRand{0, 100, 99, 150, 110, 120}
			expect(t, rr.Next(), "bar")
			expect(t, rr.Next(), "baz")
			expect(t, rr.Next(), "bar")
			expect(t, rr.Next(), "foo")
			expect(t, rr.Next(), "baz")
			expect(t, rr.Next(), "foo")
		}()
	}
}

func TestSet(t *testing.T) {
	func() {
		rr := New(DataSlice{})
		expect(t, rr.Set(DataSlice{}), true)
		expectRRList(t, rr.rrList, DataSlice{})
		expect(t, rr.weights, 0)
	}()
	func() {
		rr := New(DataSlice{})
		expect(t, rr.Set(DataSlice{
			&Data{Value: "foo"},
			&Data{Value: "bar"},
		}), true)
		expectRRList(t, rr.rrList, DataSlice{
			&Data{Key: "foo", Value: "foo", Weight: 100, rng: 100},
			&Data{Key: "bar", Value: "bar", Weight: 100, rng: 0},
		})
		expect(t, rr.weights, 200)
	}()
	func() {
		rr := New(DataSlice{
			&Data{Value: "foo"},
			&Data{Value: "bar"},
		})
		expect(t, rr.Set(DataSlice{
			&Data{Value: "hoge", Weight: 50},
			&Data{Value: "fuga", Weight: 100},
		}), true)
		expectRRList(t, rr.rrList, DataSlice{
			&Data{Key: "hoge", Value: "hoge", Weight: 50, rng: 100},
			&Data{Key: "fuga", Value: "fuga", Weight: 100, rng: 0},
		})
		expect(t, rr.weights, 150)
	}()
}

func TestAdd(t *testing.T) {
	func() {
		rr := New(DataSlice{})
		expect(t, rr.Add(&Data{}), false)
		expectRRList(t, rr.rrList, DataSlice{})
		expect(t, rr.weights, 0)
	}()
	func() {
		rr := New(DataSlice{})
		expect(t, rr.Add(
			&Data{Value: "foo"},
		), true)
		expectRRList(t, rr.rrList, DataSlice{
			&Data{Key: "foo", Value: "foo", Weight: 100, rng: 0},
		})
		expect(t, rr.weights, 100)
	}()
	func() {
		rr := New(DataSlice{})
		expect(t, rr.Add(
			&Data{Value: "foo", Weight: 10},
		), true)
		expectRRList(t, rr.rrList, DataSlice{
			&Data{Key: "foo", Value: "foo", Weight: 10, rng: 0},
		})
		expect(t, rr.weights, 10)
	}()
	func() {
		rr := New(DataSlice{
			&Data{Value: "foo"},
		})
		expect(t, rr.Add(
			&Data{Value: "foo", Weight: 10},
		), false)
		expectRRList(t, rr.rrList, DataSlice{
			&Data{Key: "foo", Value: "foo", Weight: 100, rng: 0},
		})
		expect(t, rr.weights, 100)
	}()
}

func TestReplace(t *testing.T) {
	func() {
		rr := New(DataSlice{})
		expect(t, rr.Replace(&Data{}), false)
		expectRRList(t, rr.rrList, DataSlice{})
		expect(t, rr.weights, 0)
	}()
	func() {
		rr := New(DataSlice{})
		expect(t, rr.Replace(&Data{Value: "foo"}), false)
		expectRRList(t, rr.rrList, DataSlice{})
		expect(t, rr.weights, 0)
	}()
	func() {
		rr := New(DataSlice{
			&Data{Value: "foo"},
		})
		expect(t, rr.Replace(
			&Data{Value: "foo"},
		), true)
		expectRRList(t, rr.rrList, DataSlice{
			&Data{Key: "foo", Value: "foo", Weight: 100, rng: 0},
		})
		expect(t, rr.weights, 100)
	}()
	func() {
		rr := New(DataSlice{
			&Data{Value: "foo"},
		})
		expect(t, rr.Replace(
			&Data{Value: "foo", Weight: 50},
		), true)
		expectRRList(t, rr.rrList, DataSlice{
			&Data{Key: "foo", Value: "foo", Weight: 50, rng: 0},
		})
		expect(t, rr.weights, 50)
	}()
}

func TestRemove(t *testing.T) {
	func() {
		rr := New(DataSlice{})
		expect(t, rr.Remove(nil), false)
		expectRRList(t, rr.rrList, DataSlice{})
		expect(t, rr.weights, 0)
	}()
	func() {
		rr := New(DataSlice{})
		expect(t, rr.Remove("foo"), false)
		expectRRList(t, rr.rrList, DataSlice{})
		expect(t, rr.weights, 0)
	}()
	func() {
		rr := New(DataSlice{
			&Data{Value: "foo"},
		})
		expect(t, rr.Remove("foo"), true)
		expectRRList(t, rr.rrList, DataSlice{})
		expect(t, rr.weights, 0)
	}()
}
