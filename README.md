[![Build Status](https://travis-ci.org/mix3/go-wrr.svg?branch=master)](https://travis-ci.org/mix3/go-wrr)
[![Coverage Status](https://coveralls.io/repos/mix3/go-wrr/badge.png?branch=master)](https://coveralls.io/r/mix3/go-wrr?branch=master)

# go-wrr

from cpan module [Data::WeightedRoundRobin](https://metacpan.org/pod/Data::WeightedRoundRobin)

# SYNOPSIS

```
package main

import "github.com/mix3/go-wrr"

func main() {
	rr := wrr.New(wrr.DataSlice{
		&wrr.Data{Value: "foo"},
		&wrr.Data{Value: "bar"},
		&wrr.Data{Value: "baz", Weight: 50},
		&wrr.Data{Key: "hoge", Value: []string{"fuga", "piyo"}, Weight: 120},
	})
	rr.Next() // 'foo' : 'bar' : 'baz' : []string{"fuga", "piyo"} = 100 : 100 : 50 : 120
}
```
