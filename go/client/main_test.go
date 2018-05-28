package main

import (
	"testing"
)

var c *client

func init() {
	c = newClient(address)
}

func BenchmarkCalculatorSingleAdd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		c.calculate(1.2, 3.2, "+", 0)
	}
}

func BenchmarkCalculatorSingleMinus(b *testing.B) {
	for i := 0; i < b.N; i++ {
		c.calculate(1.2, 3.2, "-", 0)
	}
}

func BenchmarkCalculatorSingleDivide(b *testing.B) {
	for i := 0; i < b.N; i++ {
		c.calculate(1.2, 3.2, "/", 0)
	}
}

func BenchmarkCalculatorSingleMultiply(b *testing.B) {
	for i := 0; i < b.N; i++ {
		c.calculate(1.2, 3.2, "*", 0)
	}
}
