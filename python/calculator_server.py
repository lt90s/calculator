import time
import grpc
from concurrent import futures

import calculator_pb2
import calculator_pb2_grpc

class CalculatorService(calculator_pb2_grpc.CalculateServicerServicer):
    def __init__(self, *args, **kwargs):
        self.op1 = 0
        self.op2 = 0
        self.opcode = calculator_pb2.CalculatorParameter.ADD
        super(CalculatorService, self).__init__(*args, *kwargs)

    def calculate(self, request, context):
        self.op1 = request.operandA
        self.op2 = request.operandB
        self.opcode = request.operCode
        result = self._do_calculate()
        return calculator_pb2.CalculatorResult(result=result)

    def _do_calculate(self):
        if self.opcode == calculator_pb2.CalculatorParameter.ADD:
            return self.op1 + self.op2
        elif self.opcode == calculator_pb2.CalculatorParameter.MINUS:
            return self.op1 - self.op2
        elif self.opcode == calculator_pb2.CalculatorParameter.MULTIPLY:
            return self.op1 * self.op2
        elif self.opcode == calculator_pb2.CalculatorParameter.DIVIDE:
            if self.op2 == 0:
                return 0
            return self.op1 / self.op2
        return 0


def startService():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=5))
    calculator_pb2_grpc.add_CalculateServicerServicer_to_server(CalculatorService(), server)
    server.add_insecure_port(':9527')
    server.start()
    try:
        while True:
            time.sleep(60*60*24)
    except KeyboardInterrupt:
        server.stop(0)

if __name__ == '__main__':
    startService()

