import logging

from computer import Machine


def main():
    logging.basicConfig(level=logging.INFO)
    with open("/home/tass/database/icfpc2020/messages/galaxy.txt") as fd:
        all_lines = [line.strip() for line in fd.readlines()]
        machine = Machine.from_lines(all_lines)
        ans = machine.eval("ap list_debug ap ap ap interact galaxy nil ap ap cons 0 0")
        print(ans)


if __name__ == "__main__":
    main()
