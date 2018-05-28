package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"google.golang.org/grpc"
)

const (
	address = "localhost:9527"
)

type client struct {
	conn *grpc.ClientConn
	sc   CalculateServicerClient
}

func newClient(adress string) *client {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
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

	var eOpcode CalculatorParameter_OperatorCode
	switch opcode {
	case "+":
		eOpcode = CalculatorParameter_ADD
	case "-":
		eOpcode = CalculatorParameter_MINUS
	case "*":
		eOpcode = CalculatorParameter_MULTIPLY
	case "/":
		eOpcode = CalculatorParameter_DIVIDE
	default:
		return 0, errors.New("Unsupported opcode")
	}
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
}
