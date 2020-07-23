class Node:
    def __init__(self, left=None, right=None):
        self.left = left
        self.right = right


class S:
    def apply(self, cache, arg):
        return S1(arg)


class S1:
    def __init__(self, arg):
        self.arg1 = arg

    def apply(self, cache, arg):
        return S2(self.arg1, arg)


class S2:
    def __init__(self, arg1, arg2):
        self.arg1 = arg1
        self.arg2 = arg2

    def apply(self, cache, arg):
        if self.arg1 not in cache:
            return [self.arg1], None
        x = cache[self.arg1]
        if (x, arg) not in cache:
            return [(x, arg)], None
        xz = cache[(x, arg)]
        return []


class I:
    pass


s1 = Node(S(), I())
s2 = Node(s1, I())

s3 = Node(S(), I())
s4 = Node(s3, I())

root = Node(s2, s4)

todo = [root]
cache = {}

while todo:
    current = todo[-1]
    if not isinstance(current.left, Node):
        cache[current] = current
        todo.pop()
        continue

    if current.left not in cache:
        todo.append(current.left)
        continue

    new_todo, value = cache[current.left].apply(cache, current.right)
    if new_todo:
        todo.extend(new_todo)
        continue

    cache[current] = value
    todo.pop()
