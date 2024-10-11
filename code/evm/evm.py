class Context:
    def __init__(self, data):
        self.stack = []
        self.pc = 0
        self.data = data

class EVM:
    def __init__(self):
        self.opcodes = {
            1: self._push,
            2: self._sub,
            3: self._add,
            4: self._return,
            5: self._eq,
            6: self._jumpi
        }

    def _eq(self, context):
        a = context.stack.pop()
        b = context.stack.pop()
        if a == b:
            context.stack.append(1)
        else :
            context.stack.append(0)

    def _jumpi(self, context):
        location = context.stack.pop()
        condition = context.stack.pop()
        if condition == 1:
            context.pc = location

    def _push(self, context):
        value_to_push = context.data[context.pc]
        context.pc += 1
        context.stack.append(value_to_push)

    def _sub(self, context):
        a = context.stack.pop()
        b = context.stack.pop()
        context.stack.append(b - a)

    def _add(self, context):
        a = context.stack.pop()
        b = context.stack.pop()
        context.stack.append(a + b)

    def _return(self, context):
        return context.stack.pop()

    def execute(self, data):
        context = Context(data)

        while context.pc < len(data):
            opcode = data[context.pc]
            executor = self.opcodes[int(opcode)]

            if not executor:
                raise RuntimeError("unknown opcode {}".format(opcode))

            context.pc += 1

            to_return = executor(context)
            if to_return:
                return to_return

        raise RuntimeError("no return")

if __name__ == '__main__2':
    data = "0101 0101 03 0102 05 010E 06 0103 04 0107 04"
    transaction = bytes.fromhex(data.replace(" ", ""))

    evm = EVM()
    print(evm.execute(transaction))

if __name__ == '__main__':
    transaction = bytes.fromhex("010201070101020304")

    evm = EVM()
    print(evm.execute(transaction))