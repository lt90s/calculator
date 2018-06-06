package main

import (
	"fmt"
	"io"
	"net"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	address = ":9527"
)

type server struct{}

func (s *server) doCalculate(op1, op2 float64, opcode CalculatorParameter_OperatorCode) float64 {
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
	return result
}

func (s *server) Calculate(ctx context.Context, p *CalculatorParameter) (*CalculatorResult, error) {
	op1 := p.OperandA
	op2 := p.OperandB
	opcode := p.OperCode
	result := s.doCalculate(op1, op2, opcode)
	return &CalculatorResult{Result: result}, nil
}

func (s *server) StreamCalculate(stream CalculateServicer_StreamCalculateServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		op1 := in.OperandA
		op2 := in.OperandB
		opcode := in.OperCode
		result := s.doCalculate(op1, op2, opcode)
		out := &CalculatorResult{Result: result}
		if err := stream.Send(out); err != nil {
			return nil
		}
	}
}

func timeInterceptor(ctx context.Context, req interface{},
	info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	now := time.Now()
	rsp, err := handler(ctx, req)
	elapse := time.Since(now)
	fmt.Printf("rpc calls %s elpased %s\n", info.FullMethod, elapse)
	return rsp, err
}

func startServer(address string) {
	l, err := net.Listen("tcp", address)
	if err != nil {
		panic(fmt.Sprintf("listen %s error %v", address, err))
	}
	gs := grpc.NewServer(grpc.UnaryInterceptor(timeInterceptor))
	RegisterCalculateServicerServer(gs, &server{})
	reflection.Register(gs)
	gs.Serve(l)
}

func main() {
	startServer(address)
}
