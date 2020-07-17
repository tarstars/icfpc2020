import unittest


class Ap:
    def __repr__(self):
        return 'ap'


class Cons:
    def __init__(self):
        self.a1 = None
        self.a2 = None

    def __repr__(self):
        return f'cons({self.a1}, {self.a2})'

    def apply(self, arg):
        if self.a1 is None:
            self.a1 = arg
            return [self]
        if self.a2 is None:
            self.a2 = arg
            return [self]
        return [arg, self.a1, self.a2]


class Number:
    def __init__(self, val):
        self.val = val

    def __repr__(self):
        return f'Number({self.val})'

    def apply(self):
        raise ValueError('can\'t apply a number')


class Nil:
    def __repr__(self):
        return 'Nil'

    def apply(self):
        raise ValueError('can\'t apply nil')


class StoredValue:
    def __init__(self, vid):
        self.vid = vid

    def __repr__(self):
        return f'StoredValue({self.vid})'

    def apply(self):
        raise ValueError('it is a bad idea to apply a stored value')


class Neg:
    def __repr__(self):
        return f'Neg()'

    def apply(self, val):
        return [Number(-val)]


class CombC:
    def __init__(self):
        self.a1 = None
        self.a2 = None

    def __repr__(self):
        return f'CombC({self.a1}, {self.a2})'

    def apply(self, val):
        if self.a1 is None:
            self.a1 = val
        if self.a2 is None:
            self.a2 = val
        return [Ap(), Ap(), self.a1, val, self.a2]


class CombB:
    def __init__(self):
        self.a1 = None
        self.a2 = None

    def __repr__(self):
        return f'CombC({self.a1}, {self.a2})'

    def apply(self, val):
        if self.a1 is None:
            self.a1 = val
        if self.a2 is None:
            self.a2 = val
        return [Ap(), self.a1, Ap(), val]


class CombS:
    def __init__(self):
        self.a1 = None
        self.a2 = None

    def __repr__(self):
        return f'CombC({self.a1}, {self.a2})'

    def apply(self, val):
        if self.a1 is None:
            self.a1 = val
        if self.a2 is None:
            self.a2 = val
        return [Ap(), Ap(), self.a1, val, Ap(), self.a2, val]


class IsNil:
    def __repr__(self):
        return 'IsNil()'

    def apply(self, val):
        if isinstance(val, Nil):
            return [FunT()]
        return [FunF()]


class FunT:
    def __init__(self):
        self.a1 = None

    def __repr__(self):
        return f'FunT({self.a1})'

    def apply(self, val):
        if self.a1 is None:
            self.a1 = val
        return [self.a1]


class FunF:
    def __init__(self):
        self.a1 = None

    def __repr__(self):
        return f'FunF({self.a1})'

    def apply(self, val):
        if self.a1 is None:
            self.a1 = val
        return [val]


class Car:
    def __repr__(self):
        return 'Car()'

    def apply(self, val):
        return [Ap(), val, FunT()]


class Cdr:
    def __repr__(self):
        return 'Cdr()'

    def apply(self, val):
        return [Ap(), val, FunF()]


class Eq:
    def __init__(self):
        self.x1 = None

    def __repr__(self):
        return f'Eq({self.x1})'

    def apply(self, val):
        if self.x1 is None:
            self.x1 = val
        if not isinstance(self.x1, Number):
            raise ValueError('not number in eq')
        if not isinstance(val, Number):
            raise ValueError('not number in eq')
        return [FunT() if self.x1.val == val.val else FunF()]


class Mul:
    def __init__(self):
        self.x1 = None

    def __repr__(self):
        return f'Mul({self.x1})'

    def apply(self, val):
        if self.x1 is None:
            self.x1 = val
        if not isinstance(self.x1, Number):
            raise ValueError('not number in eq')
        if not isinstance(val, Number):
            raise ValueError('not number in eq')
        return [Number(self.x1.val * val.val)]


class Add:
    def __init__(self):
        self.x1 = None

    def __repr__(self):
        return f'Add({self.x1})'

    def apply(self, val):
        if self.x1 is None:
            self.x1 = val
        if not isinstance(self.x1, Number):
            raise ValueError('not number in eq')
        if not isinstance(val, Number):
            raise ValueError('not number in eq')
        return [Number(self.x1.val + val.val)]


class Div:
    def __init__(self):
        self.x1 = None

    def __repr__(self):
        return f'Div({self.x1})'

    def apply(self, val):
        if self.x1 is None:
            self.x1 = val
        if not isinstance(self.x1, Number):
            raise ValueError('not number in eq')
        if not isinstance(val, Number):
            raise ValueError('not number in eq')
        return [Number(val.val // self.x1.val)]


class Lt:
    def __init__(self):
        self.x1 = None

    def __repr__(self):
        return f'Lt({self.x1})'

    def apply(self, val):
        if self.x1 is None:
            self.x1 = val
        if not isinstance(self.x1, Number):
            raise ValueError('not number in eq')
        if not isinstance(val, Number):
            raise ValueError('not number in eq')
        return [FunT() if self.x1.val < val.val else FunF()]


def FuncI():
    def __repr__(self):
        return 'I()'

    def apply(self, val):
        return [val]


def is_number(str_token):
    try:
        int(str_token)
        return True
    except ValueError:
        return False


def tokens_factory(str_token):
    if str_token == 'ap':
        return Ap()
    elif str_token == 'cons':
        return Cons()
    elif is_number(str_token=str_token):
        return Number(int(str_token))
    elif str_token == 'nil':
        return Nil()
    elif str_token[0] == ':':
        return StoredValue(str_token)
    elif str_token == 'neg':
        return Neg()
    elif str_token == 'c':
        return CombC()
    elif str_token == 'b':
        return CombB()
    elif str_token == 's':
        return CombB()
    elif str_token == 'isnil':
        return IsNil()
    elif str_token == 'car':
        return Car()
    elif str_token == 'eq':
        return Eq()
    elif str_token == 'mul':
        return Mul()
    elif str_token == 'add':
        return Add()
    elif str_token == 'lt':
        return Lt()
    elif str_token == 'div':
        return Div()
    elif str_token == 'i':
        return FuncI()
    elif str_token == 't':
        return FunT()
    elif str_token == 'cdr':
        return Cdr()

    raise NotImplementedError(f'no token {str_token}')


class Machine:
    def __init__(self):
        pass

    def execute(self, line):
        new_id, operator = map(lambda x: x.strip(), line.split('='))

        input_stack = []
        for str_token in operator.split():
            try:
                input_stack.append(tokens_factory(str_token))
            except NotImplementedError as ne:
                raise NotImplementedError(str(ne) + ' in operator ' + operator)


def execute_program(lines):
    machine = Machine()
    for line in lines:
        machine.execute(line)


def main():
    with open('/home/tass/database/icfpc2020/messages/galaxy.txt') as fd:
        all_lines = [line.strip() for line in fd.readlines()]
        execute_program(all_lines)


if __name__ == '__main__':
    main()
