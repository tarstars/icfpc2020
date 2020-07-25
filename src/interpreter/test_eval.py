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


class CNode(IdOperations):
    def __init__(self, left=None, right=None):
        super().__init__(local_hash=hash((left, right)))
        self.left = left
        self.right = right

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


class S:
    def apply(self, arg_cache, arg):
        return S1(arg)

    def __repr__(self):
        return "S"


c_s = CLeaf(S())


class S1:
    def __init__(self, arg):
        self.arg1 = arg

    def apply(self, arg_cache, arg):
        return S2(self.arg1, arg)

    def __repr__(self):
        return f"S1 [{self.arg1}]"


class S2:
    def __init__(self, arg1, arg2):
        self.arg1 = arg1
        self.arg2 = arg2

    def __repr__(self):
        return f"S2 [{self.arg1}, {self.arg2}] "

    def apply(self, cache, arg):
        x = cache[self.arg1]
        xz = cache[(x, arg)]
        yz = Node(self.arg2, arg)

        return cache[(xz, yz)]


class I:
    def apply(self, cache, arg):
        return cache[arg]

    def __repr__(self):
        return f"I"


class OurKeyError(KeyError):
    pass


c_i = CLeaf(I())

s1 = Node(c_s, c_i)
s2 = Node(s1, c_i)

s3 = Node(c_s, c_i)
s4 = Node(s3, c_i)

root = Node(s2, s4)

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
