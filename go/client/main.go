package main

import (
	"errors"
	"fmt"
	"io"
	"time"

	"golang.org/x/net/context"

	"google.golang.org/grpc"
)

const (
	address = "localhost:9527"
)

var opcodeMap = map[string]CalculatorParameter_OperatorCode{
	"+": CalculatorParameter_ADD,
	"-": CalculatorParameter_MINUS,
	"*": CalculatorParameter_MULTIPLY,
	"/": CalculatorParameter_DIVIDE,
}

type client struct {
	conn *grpc.ClientConn
	sc   CalculateServicerClient
}

func timeInterceptor(ctx context.Context, method string, req, reply interface{},
	cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	now := time.Now()
	err := invoker(ctx, method, req, reply, cc, opts...)
	elapse := time.Since(now)
	fmt.Printf("rpc call %s elapsed %s\n", method, elapse)
	return err
}

func newClient(adress string) *client {
	conn, err := grpc.Dial(address, grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(timeInterceptor))
	if err != nil {
		panic(fmt.Sprintf("grpc dial to %s failed with error %v", address, err))
	}

	sc := NewCalculateServicerClient(conn)

	return &client{conn, sc}
}

func (c *client) calculate(op1, op2 float64, opcode string, timeout time.Duration) (float64, error) {
	if timeout <= 0 {
		timeout = time.Hour * 24 * 30
	}
	ctx, cancle := context.WithTimeout(context.Background(), timeout)
	defer cancle()

	if _, ok := opcodeMap[opcode]; !ok {
		return 0, errors.New("Unsupported opcode")
	}
	eOpcode := opcodeMap[opcode]

	param := CalculatorParameter{OperandA: op1, OperandB: op2, OperCode: eOpcode}
	result, err := c.sc.Calculate(ctx, &param)
	if err != nil {
		return 0, err
	} else {
		return result.Result, err
	}
}

func (c *client) close() {
	c.conn.Close()
}

func main() {
	c := newClient(address)
	defer c.close()
	fmt.Println(c.calculate(1, 2, "+", 5*time.Millisecond))
	fmt.Println(c.calculate(1, 2, "-", 0))
	fmt.Println(c.calculate(1, 2, "*", 0))
	fmt.Println(c.calculate(1, 2, "/", 0))

	stream, err := c.sc.StreamCalculate(context.Background())
	if err != nil {
		panic(fmt.Sprintf("StreamCalculate error: %v", err))
	}

	for i := 0; i < 10; i++ {
		oprand := float64(i)
		err := stream.Send(&CalculatorParameter{
			OperandA: oprand, OperandB: oprand, OperCode: CalculatorParameter_ADD})
		if err != nil {
			panic(fmt.Sprintf("Stream send error: %v", err))
		}
	}

	err = stream.CloseSend()
	if err != nil {
		panic(fmt.Sprintf("Stream close error: %v", err))
	}

	for {
		result, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			panic(fmt.Sprintf("Stream recv error: %v", err))
		}
		fmt.Println(result.Result)
	}
}
