import logging
from typing import Dict, Union, Any


class IdGenerator:
    def __init__(self):
        self.meter = 0

    def get_id(self):
        result = self.meter
        self.meter += 1
        return result


generator = IdGenerator()
parents = {}


class IdOperations:
    def __init__(self, local_hash):
        self.local_id = generator.get_id()
        self.local_hash = local_hash

    def get_id(self):
        return self.local_id

    def __hash__(self):
        return self.local_hash

    def __eq__(self, other):
        if not isinstance(other, IdOperations):
            return NotImplemented
        return self.get_id() == other.get_id()

    def __mul__(self, other):
        if isinstance(other, CNode) or isinstance(other, CLeaf):
            return Node(self, other)
        if isinstance(other, int):
            return Node(self, CLeaf(other))
        return NotImplemented


class CNode(IdOperations):
    def __init__(self, left=None, right=None):
        super().__init__(local_hash=hash((left, right)))
        self.left = left  # a function
        self.right = right  # an argument

    def __repr__(self):
        if isinstance(self.right, CNode):
            return f"{self.left}({self.right})"
        return f"{self.left} {self.right}"


class CLeaf(IdOperations):
    def __init__(self, val):
        super().__init__(hash(val))
        self.val = val

    def __repr__(self):
        return repr(self.val)


def Node(left, right):
    p = (left.get_id(), right.get_id())

    if p not in parents:
        parents[p] = CNode(left=left, right=right)

    return parents[p]


class NArn:
    def __init__(self, n_args, func, name, args=()):
        self.n_args = n_args
        self.func = func
        self.name = name
        self.args = args

    def __repr__(self):
        args_descr = ",".join(str(current_arg) for current_arg in self.args)
        if args_descr:
            args_descr = "[" + args_descr + "]"
        return f"{self.name}" + args_descr

    def apply(self, arg_cache, x):
        new_args = self.args + (x,)

        if len(new_args) == self.n_args:
            return self.func(arg_cache, *new_args)

        return NArn(args=new_args, n_args=self.n_args, name=self.name, func=self.func)


def Leaf(n_args, func, name):
    return CLeaf(NArn(n_args=n_args, func=func, name=name))


def get_leaf_name(y):
    return y.val.name if isinstance(y, CLeaf) and isinstance(y.val, NArn) else None


def s_apply(ap, x, y, z):
    x = ap[x]
    xz = ap[(x, z)]
    yz = z if get_leaf_name(y) == "I" else Node(y, z)

    return ap[(xz, yz)]


def i_apply(ap, x):
    return ap[x]


c_s = Leaf(n_args=3, func=s_apply, name="S")
c_i = Leaf(n_args=1, func=i_apply, name="I")
c_b = Leaf(n_args=3, func=lambda ap, x, y, z: ap[(ap[x], Node(y, z))], name="B")
c_c = Leaf(n_args=3, func=lambda ap, x, y, z: ap[(ap[(ap[x], z)], y)], name="C")
c_t = Leaf(n_args=2, func=lambda ap, x, y: ap[x], name="T")
c_f = Leaf(n_args=2, func=lambda ap, x, y: ap[y], name="F")
c_car = Leaf(n_args=1, func=lambda ap, x: ap[(ap[x], c_t)], name="car")
c_cdr = Leaf(n_args=1, func=lambda ap, x: ap[(ap[x], c_f)], name="cdr")
c_cons = Leaf(n_args=3, func=lambda ap, x, y, z: ap[(ap[(ap[z], x)], y)], name="cons")
c_nil = Leaf(n_args=1, func=lambda ap, x: c_t, name="nil")
c_neg = Leaf(n_args=1, func=lambda ap, x: -ap[x], name="neg")
c_isnil = Leaf(n_args=1, func=lambda ap, x: c_t if get_leaf_name(ap[x]) == "nil" else c_f, name="isnil")
c_eq = Leaf(n_args=2, func=lambda ap, x, y: c_t if ap[x] == ap[y] else c_f, name="eq")
c_lt = Leaf(n_args=2, func=lambda ap, x, y: c_t if ap[x] < ap[y] else c_f, name="lt")
c_mul = Leaf(n_args=2, func=lambda ap, x, y: ap[x] * ap[y], name="mul")
c_add = Leaf(n_args=2, func=lambda ap, x, y: ap[x] + ap[y], name="add")
c_div = Leaf(n_args=2, func=lambda ap, x, y: int(ap[x] / ap[y]), name="add")
c_inc = Leaf(n_args=1, func=lambda ap, x: ap[x] + 1, name="inc")
c_dec = Leaf(n_args=1, func=lambda ap, x: ap[x] - 1, name="dec")


class OurKeyError(KeyError):
    pass


def create_sii_test():
    s1 = Node(c_s, c_i)
    s2 = Node(s1, c_i)

    s3 = Node(c_s, c_i)
    s4 = Node(s3, c_i)

    root = Node(s2, s4)
    return root


def create_s_inc_test():
    return c_s * c_add * c_inc * 1


root = create_s_inc_test()
todo = [root]
todo_set = set(todo)
cache: Dict[Union[CNode, tuple, CLeaf], Any] = {}

while todo:
    try:
        current = todo[-1]
        if isinstance(current, CNode):
            left = current.left

            try:
                apply_task = (cache[left], current.right)
                task_result = cache[apply_task]
            except KeyError as ke:
                raise OurKeyError(ke)

            cache[current] = task_result
            todo.pop()
        elif isinstance(current, tuple):
            left, right = current
            try:
                value = left.apply(cache, right)
            except KeyError as ke:
                raise OurKeyError(ke)
            cache[current] = value
            logging.warning(f"{current} -> {value}")
            todo.pop()
        elif isinstance(current, CLeaf):
            cache[current] = current.val
            todo.pop()
    except OurKeyError as oke:
        vta = oke.args[0].args[0]
        todo.append(vta)
        if vta in todo_set:
            raise ValueError("Simple infinite recursion detected")
        todo_set.add(vta)
        continue

print(cache[root])
