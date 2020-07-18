import logging
import unittest
import sys
from collections import deque
from typing import List

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
    def __init__(self, vid: str, machine: 'Machine'):
        self.vid = vid
        self.machine = machine

    def __repr__(self):
        return f"StoredValue({self.vid})"

    def apply(self):
        return self.expand().apply()

    def expand(self):
        return self.machine.definitions[self.vid]


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
    def __repr__(self):
        return f"Mul()"

    def apply(self, val):
        return Mul1(val)


class Mul1:
    def __init__(self, x1):
        self.x1 = x1

    def __repr__(self):
        return f"Mul({self.x1})"

    def apply(self, val):
        arg1 = self.x1.eval()
        arg2 = val.eval()

        if not isinstance(self.x1, Number):
            raise ValueError("not number in eq")
        if not isinstance(val, Number):
            raise ValueError("not number in eq")

        return Number(arg1.val * arg2.val)


class BinaryFunction:
    def __init__(self, func_name: str, func_to_apply):
        self.func_name = func_name
        self.func_to_apply = func_to_apply

    def __repr__(self):
        return f"{self.func_name}()"

    def apply(self, val):
        return BinaryFunction1(self.func_name, self.func_to_apply, val)


class BinaryFunction1:
    def __init__(self, func_name, func_to_apply, val):
        self.x1 = val
        self.func_name = func_name
        self.func_to_apply = func_to_apply

    def __repr__(self):
        return f"{self.func_name}({self.x1})"

    def apply(self, val):
        arg1 = self.x1.eval()
        arg2 = val.eval()

        if not isinstance(arg1, Number):
            raise ValueError("not number in eq")
        if not isinstance(arg2, Number):
            raise ValueError("not number in eq")

        return Number(self.func_to_apply(arg1.val, arg2.val))


class UnaryFunction:
    def __init__(self, function_name, function_to_apply):
        self.function_name = function_name
        self.function_to_apply = function_to_apply

    def __repr__(self):
        return f"{self.function_name}()"

    def apply(self, val):
        if not isinstance(val, Number):
            raise ValueError("not number in eq")

        return self.function_to_apply(val.eval())


class Add:
    def __repr__(self):
        return f"Add()"

    def apply(self, val):
        return Add1(val)


class Add1:
    def __init__(self, val):
        self.x1 = val

    def __repr__(self):
        return f"Add({self.x1})"

    def apply(self, val):
        arg1 = self.x1.eval()
        arg2 = val.eval()

        if not isinstance(arg1, Number):
            raise ValueError("not number in eq")
        if not isinstance(arg2, Number):
            raise ValueError("not number in eq")

        return Number(arg1.val + arg2.val)


def sign(x: int):
    return 1 if x >= 0 else -1






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


class Evaluable:
    """
    E -> raw_object | ap E E
    raw_function -> mul | inc ...
    raw_obect -> raw_function | raw_number
    """

    def __init__(self, func, func_arg, atomic):
        self.func = func
        self.func_arg = func_arg
        self.atomic = atomic

    @classmethod
    def from_tokens_list(cls, tokens: List):
        if isinstance(tokens[0], Ap):
            func, gobbled_first = Evaluable.from_tokens_list(tokens[1:])
            func_arg, gobbled_second = Evaluable.from_tokens_list(tokens[1 + gobbled_first:])
            return Evaluable(func=func, func_arg=func_arg, atomic=None), gobbled_first + gobbled_second + 1
        else:
            return Evaluable(None, None, tokens[0]), 1

    def is_atomic(self):
        return self.atomic is not None

    def is_function(self):
        return self.atomic is None

    def eval(self):
        if self.is_atomic():
            if hasattr(self.atomic, 'expand'):
                return self.atomic.expand().eval()
            else:
                return self.atomic
        else:
            return self.func.apply(self.func_arg)

    def apply(self, arg):
        return self.eval().apply(arg)


def div_to_zero(x, y):
    return sign(x) * sign(y) * (abs(x) // abs(y))


class Machine:
    def __init__(self, definitions):
        self.definitions = definitions

    def parse_line(self, str_operator):
        operators = []
        for str_token in str_operator.split():
            try:
                operators.append(self.tokens_factory(str_token=str_token))
            except NotImplementedError as ne:
                raise NotImplementedError(str(ne) + " in operator " + str_operator)
        return operators

    @classmethod
    def from_lines(cls, lines):
        definitions = {}
        self = cls(definitions=definitions)
        for one_line in lines:
            new_id, str_operator = map(lambda x: x.strip(), one_line.split("="))
            tokens = self.parse_line(str_operator)
            new_evaluable, gobbled = Evaluable.from_tokens_list(tokens=tokens)
            if gobbled != len(tokens):
                raise IndexError('wrong amount of tokens gobbled')
            definitions[new_id] = new_evaluable
        return self

    def eval(self, line):
        evaluable = Evaluable.from_tokens_list(self.parse_line(line))[0]
        return evaluable.eval()

    def __getitem__(self, item):
        return self.definitions[item]

    def tokens_factory(self, str_token):
        if str_token == "ap":
            return Ap()
        elif str_token == "cons":
            return Cons()
        elif is_number(str_token=str_token):
            return Number(int(str_token))
        elif str_token == "nil":
            return Nil()
        elif str_token[0] == ":" or str_token in self.definitions:
            return StoredValue(str_token, machine=self)
        elif str_token == "neg":
            return UnaryFunction(function_name='Neg', function_to_apply=lambda x: -x)
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
            return BinaryFunction(func_name='Eq', func_to_apply=lambda x, y: x == y)
        elif str_token == "mul":
            return BinaryFunction(func_name='Mul', func_to_apply=lambda x, y: x * y)
        elif str_token == "add":
            return BinaryFunction(func_name='Add', func_to_apply=lambda x, y: x + y)
        elif str_token == "lt":
            return BinaryFunction(func_name='Add', func_to_apply=lambda x, y: x < y)
        elif str_token == "div":
            return BinaryFunction(func_name='Div', func_to_apply=div_to_zero)
        elif str_token == "i":
            return FuncI()
        elif str_token == "t":
            return FunT()
        elif str_token == "cdr":
            return Cdr()
        elif str_token == "inc":
            return UnaryFunction(function_name='Inc', function_to_apply=lambda x: x + 1)
        elif str_token == "dec":
            return UnaryFunction(function_name='Dec', function_to_apply=lambda x: x - 1)


        raise NotImplementedError(f"no token {str_token}")


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
        result = m.eval(":1")
        self.assertEqual("Number(13)", str(result))

    def test_apply_mul(self):
        m = Machine.from_lines([":1 = ap ap mul 17 59"])
        self.assertEqual("Number(1003)", str(m.eval(":1")))

    def test_workspace(self):
        m = Machine.from_lines([":1 = ap ap add 8 5", ":2 = ap ap add ap neg 7 :1"])
        self.assertEqual("Number(6)", str(m.eval(":2")))

    def test_integer_division(self):
        m = Machine.from_lines(
            [":1 = ap ap div 4 2",
             ":2 = ap ap div 6 -2",
             ":3 = ap ap div 5 -3"]
        )

        self.assertEqual("Number(2)", str(m.eval(":1")))
        self.assertEqual("Number(-3)", str(m.eval(":2")))
        self.assertEqual("Number(-1)", str(m.eval(":3")))

    def test_cons(self):
        m = Machine.from_lines([":1 = ap ap ap cons 4 2 div "])
        self.assertEqual("Number(2)", str(m.eval(":1")))

    def test_car(self):
        m = Machine.from_lines([":1 = ap car ap ap cons 111 7"])
        self.assertEqual("Number(111)", str(m.eval(":1")))

    def test_conditional(self):
        m = Machine.from_lines(
            [
                ":1 = ap ap ap ap eq 777 777 car cons",
                ":2 = ap ap ap ap eq 777 778 car cons",
            ]
        )

        self.assertEqual("Car()", str(m.eval(":1")))
        self.assertEqual("Cons(None, None)", str(m.eval(":2")))

    def test_lt(self):
        m = Machine.from_lines(
            [
                ":1 = ap ap ap ap lt 777 777 car cons",
                ":2 = ap ap ap ap lt 777 778 car cons",
            ]
        )

        self.assertEqual("Cons(None, None)", str(m.eval(":1")))
        self.assertEqual("Car()", str(m.eval(":2")))

    def test_inc(self):
        m = Machine.from_lines([":1 = ap inc 1"])

        self.assertEqual("Number(2)", str(m.eval(":1")))

    def test_s(self):
        m = Machine.from_lines([":1 = ap ap ap s add inc 1"])

        self.assertEqual("Number(3)", str(m.eval(":1")))

    def test_b(self):
        m = Machine.from_lines([":1 = ap ap ap b inc dec 111"])

        self.assertEqual("Number(111)", str(m.eval(":1")))

    def test_c(self):
        m = Machine.from_lines([":1 = ap ap ap c div 1 20"])

        self.assertEqual("Number(20)", str(m.eval(":1")))


if __name__ == "__main__":
    unittest.main()
