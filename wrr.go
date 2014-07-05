package wrr

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

var DefaultWeight = 100
var BTreeBorder = 10

type WeightedRoundRobin struct {
	rrList        DataSlice
	weights       int
	defaultWeight int
	bTreeBorder   int
	rand          Rand
}

type Data struct {
	Key    interface{}
	Value  interface{}
	Weight int
	rng    int
}

type Option struct {
	DefaultWeight int
	BTreeBorder   int
}

type Rand interface {
	Intn(int) int
}

type DataSlice []*Data

func (d DataSlice) Len() int {
	return len(d)
}

func (d DataSlice) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

func (d DataSlice) Less(i, j int) bool {
	is, js := fmt.Sprintf("%v", d[i].Key), fmt.Sprintf("%v", d[j].Key)
	return is < js
}

func New(rrList DataSlice, args ...Option) *WeightedRoundRobin {
	opt := Option{DefaultWeight, BTreeBorder}
	if 0 < len(args) {
		opt = args[0]
		if opt.DefaultWeight == 0 {
			opt.DefaultWeight = DefaultWeight
		}
		if opt.BTreeBorder == 0 {
			opt.BTreeBorder = BTreeBorder
		}
	}
	wrr := &WeightedRoundRobin{
		rrList:        DataSlice{},
		weights:       0,
		defaultWeight: opt.DefaultWeight,
		bTreeBorder:   opt.BTreeBorder,
		rand:          rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	wrr.Set(rrList)
	return wrr
}

func (wrr *WeightedRoundRobin) normalize(data *Data) *Data {
	if data.Value == nil {
		return nil
	}
	if data.Weight < 0 {
		return nil
	}
	if data.Weight == 0 {
		data.Weight = wrr.defaultWeight
	}
	if data.Key == nil {
		data.Key = data.Value
	}
	return data
}

func (wrr *WeightedRoundRobin) Set(list DataSlice) bool {
	normalized := map[interface{}]*Data{}
	for _, data := range list {
		data = wrr.normalize(data)
		if data == nil {
			continue
		}
		normalized[data.Key] = data
	}

	sortedList := DataSlice{}
	for _, val := range normalized {
		sortedList = append(sortedList, val)
	}
	sort.Sort(sortedList)

	rrList := DataSlice{}
	weights := 0
	for _, val := range sortedList {
		rrList = append(DataSlice{&Data{
			Key:    val.Key,
			Value:  val.Value,
			Weight: val.Weight,
			rng:    weights,
		}}, rrList...)
		weights += val.Weight
	}

	wrr.rrList = rrList
	wrr.weights = weights

	return true
}

func equal(a, b *Data) bool {
	return fmt.Sprintf("%v", a.Key) == fmt.Sprintf("%v", b.Key)
}

func (wrr *WeightedRoundRobin) Add(value *Data) bool {
	value = wrr.normalize(value)
	if value == nil {
		return false
	}

	for _, data := range wrr.rrList {
		if equal(data, value) {
			return false
		}
	}

	wrr.rrList = append(wrr.rrList, value)
	wrr.Set(wrr.rrList)

	return true
}

func (wrr *WeightedRoundRobin) Replace(value *Data) bool {
	value = wrr.normalize(value)
	if value == nil {
		return false
	}

	for i, data := range wrr.rrList {
		if equal(data, value) {
			wrr.rrList[i] = value
			wrr.Set(wrr.rrList)
			return true
		}
	}

	return false
}

func (wrr *WeightedRoundRobin) Remove(value interface{}) bool {
	for i, data := range wrr.rrList {
		if fmt.Sprintf("%v", data.Key) == fmt.Sprintf("%v", value) {
			wrr.rrList = append(wrr.rrList[:i], wrr.rrList[i+1:]...)
			wrr.Set(wrr.rrList)
			return true
		}
	}

	return false
}

func (wrr *WeightedRoundRobin) Next() interface{} {
	if len(wrr.rrList) == 0 {
		// empty
		return nil
	}

	rweights := wrr.rand.Intn(wrr.weights)
	if len(wrr.rrList) < wrr.bTreeBorder {
		// linder
		for _, rr := range wrr.rrList {
			if rweights >= rr.rng {
				return rr.Value
			}
		}
	} else {
		start, end := 0, len(wrr.rrList)
		for start < end {
			mid := int((start + end) / 2)
			if wrr.rrList[mid].rng <= rweights {
				end = mid
			} else {
				start = mid + 1
			}
		}
		// b-tree
		return wrr.rrList[start].Value
	}

	// never reach
	return nil
}
