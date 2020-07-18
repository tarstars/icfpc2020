import logging
import unittest
import sys
from collections import deque

sys.setrecursionlimit(100000)


class Ap:
    def __repr__(self):
        return "ap"


class Cons:
    def __init__(self):
        self.a1 = None
        self.a2 = None

    def __repr__(self):
        return f"Cons({self.a1}, {self.a2})"

    def apply(self, arg):
        if self.a1 is None:
            self.a1 = arg
            return [self]
        if self.a2 is None:
            self.a2 = arg
            return [self]
        return [Ap(), Ap(), arg, self.a1, self.a2]


class Number:
    def __init__(self, val):
        self.val = val

    def __repr__(self):
        return f"Number({self.val})"

    def apply(self):
        raise ValueError("can't apply a number")


class Nil:
    def __repr__(self):
        return "Nil"

    def apply(self):
        raise ValueError("can't apply nil")


class StoredValue:
    def __init__(self, vid: str):
        self.vid = vid

    def __repr__(self):
        return f"StoredValue({self.vid})"

    def apply(self):
        raise ValueError("it is a bad idea to apply a stored value")


class Neg:
    def __repr__(self):
        return f"Neg()"

    def apply(self, val):
        return [Number(-val.val)]


class CombC:
    def __init__(self):
        self.a1 = None
        self.a2 = None

    def __repr__(self):
        return f"CombC({self.a1}, {self.a2})"

    def apply(self, val):
        if self.a1 is None:
            self.a1 = val
            return [self]

        if self.a2 is None:
            self.a2 = val
            return [self]

        return [Ap(), Ap(), self.a1, val, self.a2]


class CombB:
    def __init__(self):
        self.a1 = None
        self.a2 = None

    def __repr__(self):
        return f"CombC({self.a1}, {self.a2})"

    def apply(self, val):
        if self.a1 is None:
            self.a1 = val
            return [self]

        if self.a2 is None:
            self.a2 = val
            return [self]

        return [Ap(), self.a1, Ap(), self.a2, val]


class CombS:
    def __init__(self):
        self.a1 = None
        self.a2 = None

    def __repr__(self):
        return f"CombC({self.a1}, {self.a2})"

    def apply(self, val):
        if self.a1 is None:
            self.a1 = val
            return [self]

        if self.a2 is None:
            self.a2 = val
            return [self]

        return [Ap(), Ap(), self.a1, val, Ap(), self.a2, val]


class IsNil:
    def __repr__(self):
        return "IsNil()"

    def apply(self, val):
        if isinstance(val, Nil):
            return [FunT()]
        return [FunF()]


class FunT:
    def __init__(self):
        self.a1 = None

    def __repr__(self):
        return f"FunT({self.a1})"

    def apply(self, val):
        if self.a1 is None:
            self.a1 = val
            return [self]

        return [self.a1]


class FunF:
    def __init__(self):
        self.a1 = None

    def __repr__(self):
        return f"FunF({self.a1})"

    def apply(self, val):
        if self.a1 is None:
            self.a1 = val
            return [self]

        return [val]


class Car:
    def __repr__(self):
        return "Car()"

    def apply(self, val):
        return [Ap(), val, FunT()]


class Cdr:
    def __repr__(self):
        return "Cdr()"

    def apply(self, val):
        return [Ap(), val, FunF()]


class Inc:
    def __repr__(self):
        return "Inc()"

    def apply(self, val):
        return [Number(val.val + 1)]


class Dec:
    def __repr__(self):
        return "Dec()"

    def apply(self, val):
        return [Number(val.val - 1)]


class Eq:
    def __init__(self):
        self.a1 = None

    def __repr__(self):
        return f"Eq({self.a1})"

    def apply(self, val):
        if self.a1 is None:
            self.a1 = val
            return [self]

        if not isinstance(self.a1, Number):
            raise ValueError("not number in eq")
        if not isinstance(val, Number):
            raise ValueError("not number in eq")

        return [FunT() if self.a1.val == val.val else FunF()]


class Mul:
    def __init__(self):
        self.x1 = None

    def __repr__(self):
        return f"Mul({self.x1})"

    def apply(self, val):
        if self.x1 is None:
            self.x1 = val
            return [self]

        if not isinstance(self.x1, Number):
            raise ValueError("not number in eq")
        if not isinstance(val, Number):
            raise ValueError("not number in eq")

        return [Number(self.x1.val * val.val)]


class Add:
    def __init__(self):
        self.x1 = None

    def __repr__(self):
        return f"Add({self.x1})"

    def apply(self, val):
        if self.x1 is None:
            self.x1 = val
            return [self]

        if not isinstance(self.x1, Number):
            raise ValueError("not number in eq")
        if not isinstance(val, Number):
            raise ValueError("not number in eq")

        return [Number(self.x1.val + val.val)]


def sign(x: int):
    return 1 if x >= 0 else -1


class Div:
    def __init__(self):
        self.x1 = None

    def __repr__(self):
        return f"Div({self.x1})"

    def apply(self, val):
        if self.x1 is None:
            self.x1 = val
            return [self]

        if not isinstance(self.x1, Number):
            raise ValueError("not number in eq")
        if not isinstance(val, Number):
            raise ValueError("not number in eq")

        return [
            Number(
                sign(self.x1.val) * sign(val.val) * (abs(self.x1.val) // abs(val.val))
            )
        ]


class Lt:
    def __init__(self):
        self.x1 = None

    def __repr__(self):
        return f"Lt({self.x1})"

    def apply(self, val):
        if self.x1 is None:
            self.x1 = val
            return [self]

        if not isinstance(self.x1, Number):
            raise ValueError("not number in eq")
        if not isinstance(val, Number):
            raise ValueError("not number in eq")

        return [FunT() if self.x1.val < val.val else FunF()]


def FuncI():
    def __repr__(self):
        return "I()"

    def apply(self, val):
        return [val]


def is_number(str_token):
    try:
        int(str_token)
        return True
    except ValueError:
        return False


class LineNotProcessed(BaseException):
    pass


def tokens_factory(str_token, definitions):
    if str_token == "ap":
        return Ap()
    elif str_token == "cons":
        return Cons()
    elif is_number(str_token=str_token):
        return Number(int(str_token))
    elif str_token == "nil":
        return Nil()
    elif str_token[0] == ":" or str_token in definitions:
        return StoredValue(str_token)
    elif str_token == "neg":
        return Neg()
    elif str_token == "c":
        return CombC()
    elif str_token == "b":
        return CombB()
    elif str_token == "s":
        return CombS()
    elif str_token == "isnil":
        return IsNil()
    elif str_token == "car":
        return Car()
    elif str_token == "eq":
        return Eq()
    elif str_token == "mul":
        return Mul()
    elif str_token == "add":
        return Add()
    elif str_token == "lt":
        return Lt()
    elif str_token == "div":
        return Div()
    elif str_token == "i":
        return FuncI()
    elif str_token == "t":
        return FunT()
    elif str_token == "cdr":
        return Cdr()
    elif str_token == "inc":
        return Inc()
    elif str_token == "dec":
        return Dec()

    raise NotImplementedError(f"no token {str_token}")


class Machine:
    def __init__(self, definitions):
        self.definitions = definitions

    @staticmethod
    def parse_line(str_operator, definitions=None):
        if definitions is None:
            definitions = {}
        operators = []
        for str_token in str_operator.split():
            try:
                operators.append(tokens_factory(str_token=str_token, definitions=definitions))
            except NotImplementedError as ne:
                raise NotImplementedError(str(ne) + " in operator " + str_operator)
        return operators

    @classmethod
    def from_lines(cls, lines):
        definitions = {}
        for one_line in lines:
            new_id, str_operator = map(lambda x: x.strip(), one_line.split("="))
            definitions[new_id] = cls.parse_line(str_operator)
        return cls(definitions=definitions)

    def eval(self, line):
        output_stack = []
        input_stack = self.parse_line(line, self.definitions)
        meter = 0
        while input_stack:
            if meter % 10000000 == 0:
                logging.warning(f'iteration {meter} input stack len {len(input_stack)}')
            meter += 1
            current_operation = input_stack.pop()
            if isinstance(current_operation, StoredValue):
                try:
                    input_stack.extend(self.definitions[current_operation.vid])
                except KeyError:
                    raise LineNotProcessed
            elif isinstance(current_operation, Ap):
                funct = output_stack.pop()
                operand = output_stack.pop()
                result = funct.apply(operand)
                input_stack.extend(result)
            else:
                output_stack.append(current_operation)
        return output_stack

    def __getitem__(self, item):
        return self.definitions[item]


class CanNotPerformProgram(BaseException):
    pass


def main():
    logging.basicConfig(level=logging.INFO)
    with open("/home/tass/database/icfpc2020/messages/galaxy.txt") as fd:
        all_lines = [line.strip() for line in fd.readlines()]
        machine = Machine.from_lines(all_lines)
        ans = machine.eval('galaxy')
        print(ans)


class TestMachine(unittest.TestCase):
    def test_apply_sum(self):
        m = Machine.from_lines([":1 = ap ap add 8 5"])
        self.assertEqual("[Number(13)]", str(m.eval(":1")))

    def test_apply_mul(self):
        m = Machine.from_lines([":1 = ap ap mul 17 59"])
        self.assertEqual("[Number(1003)]", str(m.eval(":1")))

    def test_workspace(self):
        m = Machine.from_lines([":1 = ap ap add 8 5", ":2 = ap ap add ap neg 7 :1"])
        self.assertEqual("[Number(6)]", str(m.eval(":2")))

    def test_integer_division(self):
        m = Machine.from_lines(
            [":1 = ap ap div 4 2", ":2 = ap ap div 6 -2", ":3 = ap ap div 5 -3"]
        )

        self.assertEqual("[Number(2)]", str(m.eval(":1")))
        self.assertEqual("[Number(-3)]", str(m.eval(":2")))
        self.assertEqual("[Number(-1)]", str(m.eval(":3")))

    def test_cons(self):
        m = Machine.from_lines([":1 = ap ap ap cons 4 2 div "])
        self.assertEqual("[Number(2)]", str(m.eval(":1")))

    def test_car(self):
        m = Machine.from_lines([":1 = ap car ap ap cons 111 7"])
        self.assertEqual("[Number(111)]", str(m.eval(":1")))

    def test_conditional(self):
        m = Machine.from_lines(
            [
                ":1 = ap ap ap ap eq 777 777 car cons",
                ":2 = ap ap ap ap eq 777 778 car cons",
            ]
        )

        self.assertEqual("[Car()]", str(m.eval(":1")))
        self.assertEqual("[Cons(None, None)]", str(m.eval(":2")))

    def test_lt(self):
        m = Machine.from_lines(
            [
                ":1 = ap ap ap ap lt 777 777 car cons",
                ":2 = ap ap ap ap lt 777 778 car cons",
            ]
        )

        self.assertEqual("[Cons(None, None)]", str(m.eval(":1")))
        self.assertEqual("[Car()]", str(m.eval(":2")))

    def test_inc(self):
        m = Machine.from_lines([":1 = ap inc 1"])

        self.assertEqual("[Number(2)]", str(m.eval(":1")))

    def test_s(self):
        m = Machine.from_lines([":1 = ap ap ap s add inc 1"])

        self.assertEqual("[Number(3)]", str(m.eval(":1")))

    def test_b(self):
        m = Machine.from_lines([":1 = ap ap ap b inc dec 111"])

        self.assertEqual("[Number(111)]", str(m.eval(":1")))

    def test_c(self):
        m = Machine.from_lines([":1 = ap ap ap c div 1 20"])

        self.assertEqual("[Number(20)]", str(m.eval(":1")))


if __name__ == "__main__":
    main()
