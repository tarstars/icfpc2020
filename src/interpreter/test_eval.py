import logging
from typing import Dict, Union, Any


class IdGenerator:
    def __init__(self):
        self.meter = 0

    def get_id(self):
        result = self.meter
        self.meter += 1
        return result


class IdOperations:
    def __init__(self, local_hash, solver):
        self.solver = solver
        self.local_id = solver.generator.get_id()
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
            return self.solver.Node(self, other)
        if isinstance(other, int):
            return self.solver.Node(self, CLeaf(other, solver=self.solver))
        return NotImplemented


class CNode(IdOperations):
    def __init__(self,  solver, left=None, right=None):
        super().__init__(local_hash=hash((left, right)), solver=solver)
        self.left = left  # a function
        self.right = right  # an argument

    def __repr__(self):
        if isinstance(self.right, CNode):
            return f"{self.left}({self.right})"
        return f"{self.left} {self.right}"


class CLeaf(IdOperations):
    def __init__(self, val, solver):
        super().__init__(local_hash=hash(val), solver=solver)
        self.val = val

    def __repr__(self):
        return repr(self.val)


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


def get_leaf_name(y):
    return y.val.name if isinstance(y, CLeaf) and isinstance(y.val, NArn) else None


class OurKeyError(KeyError):
    pass


class Solver:
    def __init__(self):
        self.generator = IdGenerator()
        self.parents = {}

        self.c_s = self.Leaf(
            n_args=3,
            func=lambda ap, x, y, z: ap[
                (ap[(ap[x], z)], z if get_leaf_name(y) == "I" else self.Node(y, z))
            ],
            name="S",
        )
        self.c_i = self.Leaf(n_args=1, func=lambda ap, x: ap[x], name="I")
        self.c_b = self.Leaf(
            n_args=3, func=lambda ap, x, y, z: ap[(ap[x], self.Node(y, z))], name="B"
        )
        self.c_c = self.Leaf(
            n_args=3, func=lambda ap, x, y, z: ap[(ap[(ap[x], z)], y)], name="C"
        )
        self.c_t = self.Leaf(n_args=2, func=lambda ap, x, y: ap[x], name="T")
        self.c_f = self.Leaf(n_args=2, func=lambda ap, x, y: ap[y], name="F")
        self.c_car = self.Leaf(
            n_args=1, func=lambda ap, x: ap[(ap[x], self.c_t)], name="car"
        )
        self.c_cdr = self.Leaf(
            n_args=1, func=lambda ap, x: ap[(ap[x], self.c_f)], name="cdr"
        )
        self.c_cons = self.Leaf(
            n_args=3, func=lambda ap, x, y, z: ap[(ap[(ap[z], x)], y)], name="cons"
        )
        self.c_nil = self.Leaf(n_args=1, func=lambda ap, x: self.c_t, name="nil")
        self.c_neg = self.Leaf(n_args=1, func=lambda ap, x: -ap[x], name="neg")
        self.c_isnil = self.Leaf(
            n_args=1,
            func=lambda ap, x: self.c_t if get_leaf_name(ap[x]) == "nil" else self.c_f,
            name="isnil",
        )
        self.c_eq = self.Leaf(
            n_args=2,
            func=lambda ap, x, y: self.c_t if ap[x] == ap[y] else self.c_f,
            name="eq",
        )
        self.c_lt = self.Leaf(
            n_args=2,
            func=lambda ap, x, y: self.c_t if ap[x] < ap[y] else self.c_f,
            name="lt",
        )
        self.c_mul = self.Leaf(
            n_args=2, func=lambda ap, x, y: ap[x] * ap[y], name="mul"
        )
        self.c_add = self.Leaf(
            n_args=2, func=lambda ap, x, y: ap[x] + ap[y], name="add"
        )
        self.c_div = self.Leaf(
            n_args=2, func=lambda ap, x, y: int(ap[x] / ap[y]), name="add"
        )
        self.c_inc = self.Leaf(n_args=1, func=lambda ap, x: ap[x] + 1, name="inc")
        self.c_dec = self.Leaf(n_args=1, func=lambda ap, x: ap[x] - 1, name="dec")

    def Node(self, left, right):
        p = (left.get_id(), right.get_id())

        if p not in self.parents:
            self.parents[p] = CNode(solver=self, left=left, right=right)

        return self.parents[p]

    def Leaf(self, n_args, func, name):
        return CLeaf(
            NArn(n_args=n_args, func=func, name=name), solver=self
        )

    def create_sii_test(self):
        return (self.c_s * self.c_i * self.c_i) * (self.c_s * self.c_i * self.c_i)

    def create_s_inc_test(self):
        return self.c_s * self.c_add * self.c_inc * 1

    def evaluate(self, root):
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

        return cache[root]


def main():
    solver = Solver()
    print(solver.evaluate(solver.create_s_inc_test()))


if __name__ == "__main__":
    main()
