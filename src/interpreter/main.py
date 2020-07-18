import logging
import unittest
import sys
from collections import deque
from typing import List

import simplejson

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
    def __init__(self, vid: str, machine: "Machine"):
        self.vid = vid
        self.machine = machine

    def __repr__(self):
        return f"StoredValue({self.vid})"

    def apply(self):
        return self.expand().apply()

    def expand(self):
        return self.machine.definitions[self.vid]


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


class BinaryFunction:
    def __init__(self, func_name: str, func_to_apply, numerical_function=True):
        self.func_name = func_name
        self.func_to_apply = func_to_apply
        self.numerical_function = numerical_function

    def __repr__(self):
        return f"{self.func_name}()"

    def apply(self, val):
        return BinaryFunction1(
            self.func_name,
            self.func_to_apply,
            val,
            numerical_function=self.numerical_function,
        )

    def eval(self):
        return self


class BinaryFunction1:
    def __init__(self, func_name, func_to_apply, val, numerical_function):
        self.x1 = val
        self.func_name = func_name
        self.func_to_apply = func_to_apply
        self.numerical_function = numerical_function

    def __repr__(self):
        return f"{self.func_name}({self.x1})"

    def apply(self, val):
        if self.numerical_function:
            arg1 = self.x1.eval()
            arg2 = val.eval()

            if not isinstance(arg1, Number):
                raise ValueError(f"not number1 in eq for {self.func_name}")
            if not isinstance(arg2, Number):
                raise ValueError(f"not number2 in eq for {self.func_name}")
            return Number(self.func_to_apply(arg1.val, arg2.val))
        return self.func_to_apply(self.x1, val).eval()


class TernaryFunction:
    def __init__(self, func_name: str, func_to_apply):
        self.func_name = func_name
        self.func_to_apply = func_to_apply

    def __repr__(self):
        return f"{self.func_name}()"

    def apply(self, val):
        return TernaryFunction1(self.func_name, self.func_to_apply, val)


class TernaryFunction1:
    def __init__(self, func_name, func_to_apply, val):
        self.x1 = val
        self.func_name = func_name
        self.func_to_apply = func_to_apply

    def __repr__(self):
        return f"{self.func_name}_1({self.x1})"

    def apply(self, val):
        return TernaryFunction2(self.func_name, self.func_to_apply, self.x1, val)


class TernaryFunction2:
    def __init__(self, func_name, func_to_apply, x1, val):
        self.x1 = x1
        self.x2 = val
        self.func_name = func_name
        self.func_to_apply = func_to_apply

    def __repr__(self):
        return f"{self.func_name}_2({self.x1}, {self.x2})"

    def apply(self, val):
        arg1 = self.x1
        arg2 = self.x2
        arg3 = val

        return self.func_to_apply(arg1, arg2, arg3).eval()


class UnaryFunction:
    def __init__(self, function_name, function_to_apply, numeric=True):
        self.function_name = function_name
        self.function_to_apply = function_to_apply
        self.numeric = numeric

    def __repr__(self):
        return f"{self.function_name}()"

    def apply(self, val):
        if self.numeric:
            result = val.eval()
            if not isinstance(result, Number):
                raise ValueError(f"not number in eq for {self.function_name}")
            return Number(self.function_to_apply(result.val))
        return self.function_to_apply(val).eval()


def sign(x: int):
    return 1 if x >= 0 else -1


def is_number(str_token):
    try:
        int(str_token)
        return True
    except ValueError:
        return False


class LineNotProcessed(BaseException):
    pass


def comb_b(a1, a2, a3):
    return Evaluable.from_function(
        func=a1, arg=Evaluable.from_function(func=a2, arg=a3)
    )


def comb_c(a1, a2, a3):
    return Evaluable.from_function(
        func=Evaluable.from_function(func=a1, arg=a3), arg=a2
    )


def comb_s(a1, a2, a3):
    return Evaluable.from_function(
        func=Evaluable.from_function(func=a1, arg=a3),
        arg=Evaluable.from_function(func=a2, arg=a3),
    )


func_t = BinaryFunction(
    func_name="FuncT", func_to_apply=lambda x, y: x, numerical_function=False,
)

func_f = BinaryFunction(
    func_name="FuncF", func_to_apply=lambda x, y: y, numerical_function=False,
)


def comb_car(a1):
    return Evaluable.from_function(func=a1, arg=func_t)


def comb_cdr(a1):
    return Evaluable.from_function(func=a1, arg=func_f)


def comb_cons(a1, a2, a3):
    return Evaluable.from_function(
        func=Evaluable.from_function(func=a3, arg=a1), arg=a2
    )


def comb_is_nil(arg):
    """
    ap nil x0 = t
    ap isnil nil   =   t
    ap isnil ap ap cons x0 x1   =   f

    :param arg:
    :return:
    """
    val = arg.eval()
    cond = isinstance(val, UnaryFunction) and val.function_name == "Nil"
    return func_t if cond else func_f


def func_lt(a1, a2):
    return func_t if a1.eval().val < a2.eval().val else func_f


def is_cons(val_eval):
    return isinstance(val_eval, TernaryFunction2) and val_eval.func_name == "Cons"


class DebugList:
    def __init__(self, value, is_cons):
        self.value = value
        self.is_cons = is_cons

    def __repr__(self):
        return f"{self.value}"

    def eval(self):
        return self


def comb_debug(val):
    val_eval = val.eval()
    if is_cons(val_eval):
        x1, x2 = val_eval.x1, val_eval.x2
        rec_x1 = comb_debug(x1)
        rec_x2 = comb_debug(x2)
        if rec_x2.is_cons:
            return DebugList([rec_x1.value] + rec_x2.value, True)
        else:
            if rec_x2.value == "Nil":
                return DebugList([rec_x1.value], True)
            return DebugList([rec_x1.value, rec_x2.value], True)
    else:
        if isinstance(val_eval, Number):
            return DebugList(val_eval.val, False)
        if isinstance(val_eval, UnaryFunction) and val_eval.function_name == "Nil":
            return DebugList("Nil", False)
        return DebugList(val_eval, False)


class Evaluable:
    """
    E -> raw_object | ap E E
    raw_function -> mul | inc ...
    raw_obect -> raw_function | raw_number
    """

    def __init__(self, func, func_arg, atomic):
        if not (
            ((func is not None) == (func_arg is not None))
            and ((func is not None) != (atomic is not None))
        ):
            print()
        self.func = func
        self.func_arg = func_arg
        self.atomic = atomic

    def __repr__(self):
        if self.atomic is not None:
            return f"atomic {self.atomic}"
        return f"func {self.func} ({self.func_arg})"

    @classmethod
    def from_function(cls, func, arg):
        return cls(func=func, func_arg=arg, atomic=None)

    @classmethod
    def from_atomic(cls, atomic):
        return cls(func=None, func_arg=None, atomic=atomic)

    @classmethod
    def from_tokens_list(cls, tokens: List):
        if isinstance(tokens[0], Ap):
            func, gobbled_first = Evaluable.from_tokens_list(tokens[1:])
            assert func is not None
            func_arg, gobbled_second = Evaluable.from_tokens_list(
                tokens[1 + gobbled_first :]
            )
            return (
                Evaluable.from_function(func=func, arg=func_arg),
                gobbled_first + gobbled_second + 1,
            )
        else:
            return Evaluable.from_atomic(tokens[0]), 1

    def is_atomic(self):
        return self.atomic is not None

    def is_function(self):
        return self.atomic is None

    def eval(self):
        if self.is_atomic():
            if hasattr(self.atomic, "expand"):
                return self.atomic.expand().eval()
            else:
                return self.atomic
        else:
            return self.func.apply(self.func_arg)

    def apply(self, arg):
        # logging.warning(f"before {str(self)} with argument = {arg}")
        result = self.eval().apply(arg)
        # logging.warning(f"after {str(self)} with argument = {arg} -> {result}")
        return result


def div_to_zero(x, y):
    return sign(x) * sign(y) * (abs(x) // abs(y))


class Machine:
    def __init__(self, definitions):
        self.definitions = definitions
        self.not_found_list = set()

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
        # "ap ap s ap ap c ap eq 0 1 ap ap b ap mul 2 ap ap b pwr2 ap add -1"
        for one_line in lines + [
            "pwr2 = ap ap s ap ap c ap eq 0 1 ap ap b ap mul 2 ap ap b pwr2 dec",
            "checkerboard = ap ap s ap ap b s ap ap c ap ap b c ap ap b ap c ap c ap ap s ap ap b s ap ap b ap b ap ap s i i lt eq ap ap s mul i nil ap ap s ap ap b s ap ap b ap b cons ap ap s ap ap b s ap ap b ap b cons ap c div ap c ap ap s ap ap b b ap ap c ap ap b b add neg ap ap b ap s mul div ap ap c ap ap b b checkerboard ap ap c add 2",
        ]:
            new_id, str_operator = map(lambda x: x.strip(), one_line.split("="))
            definitions[new_id] = self.from_str_operator(str_operator=str_operator)
            if new_id in self.not_found_list:
                self.not_found_list.remove(new_id)
        if self.not_found_list:
            raise ValueError("Not found: " + ",".join(self.not_found_list))
        return self

    def from_str_operator(self, str_operator):
        tokens = self.parse_line(str_operator)
        new_evaluable, gobbled = Evaluable.from_tokens_list(tokens=tokens)
        if gobbled != len(tokens):
            raise IndexError("wrong amount of tokens gobbled")
        return new_evaluable

    def eval(self, line):
        evaluable = Evaluable.from_tokens_list(self.parse_line(line))[0]
        return evaluable.eval()

    def __getitem__(self, item):
        return self.definitions[item]

    def tokens_factory(self, str_token):
        if str_token == "ap":
            return Ap()
        elif str_token == "cons":
            return TernaryFunction(func_name="Cons", func_to_apply=comb_cons)
        elif is_number(str_token=str_token):
            return Number(int(str_token))
        elif str_token == "nil":
            return UnaryFunction(
                function_name="Nil", numeric=False, function_to_apply=lambda x: func_t
            )
        elif str_token == "neg":
            return UnaryFunction(function_name="Neg", function_to_apply=lambda x: -x)
        elif str_token == "c":
            return TernaryFunction(func_name="CombC", func_to_apply=comb_c)
        elif str_token == "b":
            return TernaryFunction(func_name="CombB", func_to_apply=comb_b)
        elif str_token == "s":
            return TernaryFunction(func_name="CombS", func_to_apply=comb_s)
        elif str_token == "isnil":
            return UnaryFunction(
                function_name="IsNil", function_to_apply=comb_is_nil, numeric=False
            )
        elif str_token == "car":
            return UnaryFunction(
                function_name="Car", function_to_apply=comb_car, numeric=False
            )
        elif str_token == "eq":
            return BinaryFunction(
                func_name="Eq",
                numerical_function=False,
                func_to_apply=lambda x, y: func_t
                if x.eval().val == y.eval().val
                else func_f,
            )
        elif str_token == "mul":
            return BinaryFunction(func_name="Mul", func_to_apply=lambda x, y: x * y)
        elif str_token == "add":
            return BinaryFunction(func_name="Add", func_to_apply=lambda x, y: x + y)
        elif str_token == "lt":
            return BinaryFunction(
                func_name="Lt", func_to_apply=func_lt, numerical_function=False
            )
        elif str_token == "div":
            return BinaryFunction(func_name="Div", func_to_apply=div_to_zero)
        elif str_token == "i":
            return UnaryFunction(
                function_name="Ident", numeric=False, function_to_apply=lambda x: x
            )
        elif str_token == "t":
            return func_t
        elif str_token == "f":
            return func_f
        elif str_token == "cdr":
            return UnaryFunction(
                function_name="Cdr", function_to_apply=comb_cdr, numeric=False
            )
        elif str_token == "inc":
            return UnaryFunction(function_name="Inc", function_to_apply=lambda x: x + 1)
        elif str_token == "dec":
            return UnaryFunction(function_name="Dec", function_to_apply=lambda x: x - 1)
        elif str_token == "list_debug":
            return UnaryFunction(
                function_name="ListDebug", numeric=False, function_to_apply=comb_debug
            )

        if str_token not in self.definitions:
            self.not_found_list.add(str_token)

        return StoredValue(str_token, machine=self)


class CanNotPerformProgram(BaseException):
    pass


def main():
    logging.basicConfig(level=logging.INFO)
    with open("/home/tass/database/icfpc2020/messages/galaxy.txt") as fd:
        all_lines = [line.strip() for line in fd.readlines()]
        machine = Machine.from_lines(all_lines)
        ans = machine.eval("galaxy")
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
            [":1 = ap ap div 4 2", ":2 = ap ap div 6 -2", ":3 = ap ap div 5 -3"]
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
        self.assertEqual("Cons()", str(m.eval(":2")))

    def test_eq(self):
        m = Machine.from_lines([":1 = ap ap eq 1 2", ":2 = ap ap eq 777 777",])

        self.assertEqual("FuncF()", str(m.eval(":1")))
        self.assertEqual("FuncT()", str(m.eval(":2")))

    def test_lt(self):
        m = Machine.from_lines(
            [
                ":1 = ap ap ap ap lt 777 777 car cons",
                ":2 = ap ap ap ap lt 777 778 car cons",
            ]
        )

        self.assertEqual("Cons()", str(m.eval(":1")))
        self.assertEqual("Car()", str(m.eval(":2")))

    def test_true_false(self):
        m = Machine.from_lines([":1 = ap ap t car cons", ":2 = ap ap f car cons",])

        self.assertEqual("Car()", str(m.eval(":1")))
        self.assertEqual("Cons()", str(m.eval(":2")))

    def test_inc_00(self):
        m = Machine.from_lines([":1 = ap inc 1"])

        self.assertEqual("Number(2)", str(m.eval(":1")))

    def test_inc_01(self):
        m = Machine.from_lines([":1 = ap inc -2"])

        self.assertEqual("Number(-1)", str(m.eval(":1")))

    def test_s(self):
        m = Machine.from_lines([":1 = ap ap ap s add inc 1"])

        self.assertEqual("Number(3)", str(m.eval(":1")))

    def test_b(self):
        m = Machine.from_lines([":1 = ap ap ap b inc dec 111"])

        self.assertEqual("Number(111)", str(m.eval(":1")))

    def test_c(self):
        m = Machine.from_lines([":1 = ap ap ap c div 1 20"])

        self.assertEqual("Number(20)", str(m.eval(":1")))

    def test_nil(self):
        m = Machine.from_lines([":1 = ap nil ap pwr2 1000000"])

        self.assertEqual("FuncT()", str(m.eval(":1")))

    def test_pwr2(self):
        for pow_ind in range(5):
            m = Machine.from_lines([f":1 = ap pwr2 {pow_ind}"])
            self.assertEqual(f"Number({2 ** pow_ind})", str(m.eval(":1")))

    def test_is_nil_00(self):
        m = Machine.from_lines([":1 = ap isnil nil"])

        self.assertEqual("FuncT()", str(m.eval(":1")))

    def test_is_nil_01(self):
        m = Machine.from_lines([":1 = ap isnil ap ap cons 10 nil"])

        self.assertEqual("FuncF()", str(m.eval(":1")))

    def test_is_nil_02(self):
        m = Machine.from_lines([":1 = ap isnil ap ap cons nil nil"])

        self.assertEqual("FuncF()", str(m.eval(":1")))

    def test_identity_00(self):
        m = Machine.from_lines([":1 = ap i 7"])

        self.assertEqual("Number(7)", str(m.eval(":1")))

    def test_checkerboard(self):
        m = Machine.from_lines([":1 = ap list_debug ap ap checkerboard 7 0"])

        val = m.eval(":1")

        # with open('/home/tass/database/icfpc2020/checker', 'w') as fd:
        #     simplejson.dump(val.value, fd)

        self.assertEqual(
            "[[0, 0], [0, 2], [0, 4], [0, 6], [1, 1], [1, 3], [1, 5], [2, 0], [2, 2], [2, 4], [2, 6], [3, 1], [3, 3], [3, 5], [4, 0], [4, 2], [4, 4], [4, 6], [5, 1], [5, 3], [5, 5], [6, 0], [6, 2], [6, 4], [6, 6]]",
            str(val),
        )


"""
ap ap ap b inc dec 111
ap inc ap dec 111
"""

if __name__ == "__main__":
    unittest.main()
