package main

import (
	"fmt"
	"net"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	address = ":9527"
)

type server struct{}

func (s *server) Calculate(ctx context.Context, p *CalculatorParameter) (*CalculatorResult, error) {
	op1 := p.OperandA
	op2 := p.OperandB
	opcode := p.OperCode
	var result float64
	switch opcode {
	case CalculatorParameter_ADD:
		result = op1 + op2
	case CalculatorParameter_MINUS:
		result = op1 - op2
	case CalculatorParameter_MULTIPLY:
		result = op1 * op2
	case CalculatorParameter_DIVIDE:
		if op2 == 0 {
			result = 0
		} else {
			result = op1 / op2
		}
	default:
		result = 0
	}
	return &CalculatorResult{Result: result}, nil
}

func startServer(address string) {
	l, err := net.Listen("tcp", address)
	if err != nil {
		panic(fmt.Sprintf("listen %s error %v", address, err))
	}
	gs := grpc.NewServer()
	RegisterCalculateServicerServer(gs, &server{})
	reflection.Register(gs)
	gs.Serve(l)
}

func main() {
	startServer(address)
}
