import grpc

import calculator_pb2
import calculator_pb2_grpc

class ExOperCode(Exception):
    pass

class CalculatorClient(object):
    VALID_OPCODE_MAPPINGS = {
        '+': calculator_pb2.CalculatorParameter.ADD,
        '-': calculator_pb2.CalculatorParameter.MINUS,
        '*': calculator_pb2.CalculatorParameter.MULTIPLY,
        '/': calculator_pb2.CalculatorParameter.DIVIDE
    }

    VALID_OPCODE_CHARS = VALID_OPCODE_MAPPINGS.keys()

    def __init__(self, addr):
        channel = grpc.insecure_channel(addr)
        self.stub = calculator_pb2_grpc.CalculateServicerStub(channel)

    def calculate(self, op1, opcode, op2):
        if opcode not in self.VALID_OPCODE_CHARS:
            raise ExOperCode('Invalid operator code: {}, choices are {}'
                             .format(opcode, self.VALID_OPCODE_CHARS))

        opcode = self.VALID_OPCODE_MAPPINGS[opcode]

        # omit parameter sanity checks
        parameter = calculator_pb2.CalculatorParameter(operandA=op1, operandB=op2, operCode=opcode)
        result = self.stub.calculate(parameter)
        return result.result


def calculate():
    client = CalculatorClient('localhost:9527')
    print('1 + 2 = {}'.format(client.calculate(1, '+', 2)))
    print('1 - 12 = {}'.format(client.calculate(1, '-', 12)))
    print('2 * 22.4 = {}'.format(client.calculate(2, '*', 22.4)))
    print('2 / 13 = {}'.format(client.calculate(2, '/', 13)))


if __name__ == '__main__':
    calculate()

