import numpy as np
import os
import random
import re
import simplejson
import subprocess
import time
import queue

from PyQt5 import QtCore, Qt
from PyQt5.QtCore import QRunnable, QThreadPool
from PyQt5.QtGui import QPainter, QColor
from PyQt5.QtWidgets import QApplication, QLabel, QMainWindow

from computer import Machine


def game_state():
    return "ap ap cons 6 ap ap cons ap ap cons 0 ap ap cons 12 ap ap cons 270608505102339400 ap ap cons 2 ap ap cons 0 ap ap cons 2 ap ap cons ap ap cons 0 ap ap cons 1 ap ap cons ap ap cons -16 47 ap ap cons ap ap cons 0 -1 ap ap cons ap ap cons 510 ap ap cons 0 ap ap cons 0 ap ap cons 1 nil ap ap cons 0 ap ap cons 64 ap ap cons 1 nil ap ap cons nil ap ap cons 4 ap ap cons ap ap cons 1 ap ap cons ap ap cons 16 ap ap cons 128 nil ap ap cons ap ap cons ap ap cons ap ap cons 1 ap ap cons 0 ap ap cons ap ap cons 15 -48 ap ap cons ap ap cons -1 0 ap ap cons ap ap cons 205 ap ap cons 30 ap ap cons 10 ap ap cons 1 nil ap ap cons 0 ap ap cons 64 ap ap cons 1 nil ap ap cons ap ap cons ap ap cons 0 ap ap cons ap ap cons 1 1 nil nil nil ap ap cons ap ap cons ap ap cons 0 ap ap cons 1 ap ap cons ap ap cons -16 47 ap ap cons ap ap cons 0 -1 ap ap cons ap ap cons 510 ap ap cons 0 ap ap cons 0 ap ap cons 1 nil ap ap cons 0 ap ap cons 64 ap ap cons 1 nil ap ap cons nil nil nil nil ap ap cons nil ap ap cons ap ap cons ap ap cons 16 ap ap cons 128 nil ap ap cons ap ap cons ap ap cons 0 ap ap cons ap ap cons ap ap cons ap ap cons 1 ap ap cons 0 ap ap cons ap ap cons 16 -48 ap ap cons ap ap cons 0 0 ap ap cons ap ap cons 206 ap ap cons 30 ap ap cons 10 ap ap cons 1 nil ap ap cons 0 ap ap cons 64 ap ap cons 1 nil ap ap cons nil nil ap ap cons ap ap cons ap ap cons 0 ap ap cons 1 ap ap cons ap ap cons -16 48 ap ap cons ap ap cons 0 0 ap ap cons ap ap cons 510 ap ap cons 0 ap ap cons 0 ap ap cons 1 nil ap ap cons 0 ap ap cons 64 ap ap cons 1 nil ap ap cons nil nil nil nil ap ap cons ap ap cons 1 ap ap cons ap ap cons ap ap cons ap ap cons 1 ap ap cons 0 ap ap cons ap ap cons 15 -48 ap ap cons ap ap cons -1 0 ap ap cons ap ap cons 205 ap ap cons 30 ap ap cons 10 ap ap cons 1 nil ap ap cons 0 ap ap cons 64 ap ap cons 1 nil ap ap cons ap ap cons ap ap cons 0 ap ap cons ap ap cons 1 1 nil nil nil ap ap cons ap ap cons ap ap cons 0 ap ap cons 1 ap ap cons ap ap cons -16 47 ap ap cons ap ap cons 0 -1 ap ap cons ap ap cons 510 ap ap cons 0 ap ap cons 0 ap ap cons 1 nil ap ap cons 0 ap ap cons 64 ap ap cons 1 nil ap ap cons nil nil nil nil ap ap cons ap ap cons 2 ap ap cons ap ap cons ap ap cons ap ap cons 1 ap ap cons 0 ap ap cons ap ap cons 13 -48 ap ap cons ap ap cons -2 0 ap ap cons ap ap cons 204 ap ap cons 30 ap ap cons 10 ap ap cons 1 nil ap ap cons 0 ap ap cons 64 ap ap cons 1 nil ap ap cons ap ap cons ap ap cons 0 ap ap cons ap ap cons 1 1 nil nil nil ap ap cons ap ap cons ap ap cons 0 ap ap cons 1 ap ap cons ap ap cons -16 45 ap ap cons ap ap cons 0 -2 ap ap cons ap ap cons 510 ap ap cons 0 ap ap cons 0 ap ap cons 1 nil ap ap cons 0 ap ap cons 64 ap ap cons 1 nil ap ap cons nil nil nil nil ap ap cons ap ap cons 3 ap ap cons ap ap cons ap ap cons ap ap cons 1 ap ap cons 0 ap ap cons ap ap cons 10 -47 ap ap cons ap ap cons -3 1 ap ap cons ap ap cons 203 ap ap cons 30 ap ap cons 10 ap ap cons 1 nil ap ap cons 0 ap ap cons 64 ap ap cons 1 nil ap ap cons ap ap cons ap ap cons 0 ap ap cons ap ap cons 1 0 nil nil nil ap ap cons ap ap cons ap ap cons 0 ap ap cons 1 ap ap cons ap ap cons -16 42 ap ap cons ap ap cons 0 -3 ap ap cons ap ap cons 510 ap ap cons 0 ap ap cons 0 ap ap cons 1 nil ap ap cons 0 ap ap cons 64 ap ap cons 1 nil ap ap cons nil nil nil nil ap ap cons ap ap cons 4 ap ap cons ap ap cons ap ap cons ap ap cons 1 ap ap cons 0 ap ap cons ap ap cons 6 -45 ap ap cons ap ap cons -4 2 ap ap cons ap ap cons 202 ap ap cons 30 ap ap cons 10 ap ap cons 1 nil ap ap cons 0 ap ap cons 64 ap ap cons 1 nil ap ap cons ap ap cons ap ap cons 0 ap ap cons ap ap cons 1 0 nil nil nil ap ap cons ap ap cons ap ap cons 0 ap ap cons 1 ap ap cons ap ap cons -16 38 ap ap cons ap ap cons 0 -4 ap ap cons ap ap cons 510 ap ap cons 0 ap ap cons 0 ap ap cons 1 nil ap ap cons 0 ap ap cons 64 ap ap cons 1 nil ap ap cons nil nil nil nil ap ap cons ap ap cons 5 ap ap cons ap ap cons ap ap cons ap ap cons 1 ap ap cons 0 ap ap cons ap ap cons 2 -42 ap ap cons ap ap cons -4 3 ap ap cons ap ap cons 202 ap ap cons 30 ap ap cons 10 ap ap cons 1 nil ap ap cons 0 ap ap cons 64 ap ap cons 1 nil ap ap cons nil nil ap ap cons ap ap cons ap ap cons 0 ap ap cons 1 ap ap cons ap ap cons -16 33 ap ap cons ap ap cons 0 -5 ap ap cons ap ap cons 510 ap ap cons 0 ap ap cons 0 ap ap cons 1 nil ap ap cons 0 ap ap cons 64 ap ap cons 1 nil ap ap cons nil nil nil nil ap ap cons ap ap cons 6 ap ap cons ap ap cons ap ap cons ap ap cons 1 ap ap cons 0 ap ap cons ap ap cons -2 -38 ap ap cons ap ap cons -4 4 ap ap cons ap ap cons 202 ap ap cons 30 ap ap cons 10 ap ap cons 1 nil ap ap cons 0 ap ap cons 64 ap ap cons 1 nil ap ap cons nil nil ap ap cons ap ap cons ap ap cons 0 ap ap cons 1 ap ap cons ap ap cons -16 27 ap ap cons ap ap cons 0 -6 ap ap cons ap ap cons 510 ap ap cons 0 ap ap cons 0 ap ap cons 1 nil ap ap cons 0 ap ap cons 64 ap ap cons 1 nil ap ap cons nil nil nil nil ap ap cons ap ap cons 7 ap ap cons ap ap cons ap ap cons ap ap cons 1 ap ap cons 0 ap ap cons ap ap cons -6 -33 ap ap cons ap ap cons -4 5 ap ap cons ap ap cons 202 ap ap cons 30 ap ap cons 10 ap ap cons 1 nil ap ap cons 0 ap ap cons 64 ap ap cons 1 nil ap ap cons nil nil ap ap cons ap ap cons ap ap cons 0 ap ap cons 1 ap ap cons ap ap cons -16 20 ap ap cons ap ap cons 0 -7 ap ap cons ap ap cons 510 ap ap cons 0 ap ap cons 0 ap ap cons 1 nil ap ap cons 0 ap ap cons 64 ap ap cons 1 nil ap ap cons nil nil nil nil ap ap cons ap ap cons 8 ap ap cons ap ap cons ap ap cons ap ap cons 1 ap ap cons 0 ap ap cons ap ap cons -10 -27 ap ap cons ap ap cons -4 6 ap ap cons ap ap cons 202 ap ap cons 30 ap ap cons 10 ap ap cons 1 nil ap ap cons 20 ap ap cons 64 ap ap cons 1 nil ap ap cons ap ap cons ap ap cons 2 ap ap cons ap ap cons -16 12 ap ap cons 30 ap ap cons 25 ap ap cons 4 nil nil nil ap ap cons ap ap cons ap ap cons 0 ap ap cons 1 ap ap cons ap ap cons -16 12 ap ap cons ap ap cons 0 -8 ap ap cons ap ap cons 0 ap ap cons 0 ap ap cons 0 ap ap cons 0 nil ap ap cons 0 ap ap cons 64 ap ap cons 1 nil ap ap cons nil nil nil nil nil nil ap ap cons nil nil ap ap cons 1 ap ap cons nil nil"

def game_menu():
    return "ap car ap ap ap interact :1338 ap ap cons 1 ap ap cons ap ap cons 11 nil ap ap cons 0 ap ap cons nil nil ap ap vec 0 4"

def get_draw_data(galaxy_coords, state):
    icfpc_base = "/home/tass/go/src/github.com/tarstars/icfpc2020/"
    computer_partial_path = "diseaz/cmd/galaxy-eval/galaxy-eval"
    galaxy_path = "/home/tass/database/icfpc2020/messages/galaxy_fixed.txt"

    print("galaxy coords = ", galaxy_coords)
    to_evaluate = f"ap car ap ap ap interact :1338 {state} ap ap vec {galaxy_coords[0]} {galaxy_coords[1]}"
    with open(galaxy_path) as source_fd:
        with open('local_galaxy.txt', 'w') as destination_fd:
            destination_fd.write('\n'.join([source_fd.read(), to_evaluate]))

    incoming_process = subprocess.Popen(
        [
            os.path.join(icfpc_base, computer_partial_path),
            "-key=faa0647bb89f42d6a0a1850cf1b71954",
            'local_galaxy.txt',
        ],
        stdout=subprocess.PIPE,
    )
    incoming = incoming_process.communicate()[0].decode('utf8')
    raw_result = simplejson.loads(incoming)
    print(raw_result)
    return [(v['X'], v['Y'], i) for i, inner_list in enumerate(raw_result['Picture']['Draw']) for v in
            inner_list], raw_result


class Worker(QRunnable):
    def __init__(self, queue, stop_queue, set_coord_queue, rewind_queue, *args, **kwargs):
        super().__init__()
        self.queue = queue
        self.stop_queue = stop_queue
        self.set_coord_queue = set_coord_queue
        self.rewind_queue = rewind_queue

    def run(self):
        state_stack = []

        galaxy_coords = 0, 0  # ap ap cons nil ap ap cons nil ap ap cons nil ap ap cons nil ap ap cons nil ap ap cons 22288 nil ap ap cons 0 ap ap cons nil nil
        state = game_menu() # "ap ap cons 1 ap ap cons ap ap cons 6 nil ap ap cons 0 ap ap cons nil nil"  # state = 'ap ap cons 0 ap ap cons ap ap cons 0 nil ap ap cons 0 ap ap cons nil nil'
        while True:
            draw_data, raw_output = get_draw_data(galaxy_coords, state=state)
            state = raw_output['Results'][0]
            state_stack.append(state)
            print('states in stack', len(state_stack))
            self.queue.put({'draw_data': draw_data}, block=False)  # , 'special_points': special_points
            if not self.stop_queue.empty():
                break
            if not self.rewind_queue.empty():
                self.rewind_queue.get()
                if state_stack:
                    state_stack.pop()
                if state_stack:
                    state = state_stack.pop()
            galaxy_coords = self.set_coord_queue.get()


class MainWindow(QMainWindow):
    def __init__(self):
        super().__init__()
        side = 900
        self.setFixedWidth(side)
        self.setFixedHeight(side)

        self.threadpool = QThreadPool()
        self.input_queue = queue.Queue()
        self.set_coord_queue = queue.Queue()
        self.stop_queue = queue.Queue()
        self.rewind_queue = queue.Queue()
        worker = Worker(queue=self.input_queue, stop_queue=self.stop_queue, set_coord_queue=self.set_coord_queue, rewind_queue=self.rewind_queue)
        self.threadpool.start(worker)

        self.all_points = []
        # self.special_points = []

        self.update_timer = QtCore.QTimer(self)
        self.update_timer.timeout.connect(self.repaint)
        self.update_timer.start(1000)
        self.set_scale()

    def set_scale(self):
        if not self.all_points:
            s_lt = -20, -20
            s_rb = 20, 20
        else:
            min_x = min(v[0] for v in self.all_points)
            min_y = min(v[1] for v in self.all_points)
            max_x = max(v[0] for v in self.all_points)
            max_y = max(v[1] for v in self.all_points)
            rx = max_x - min_x
            ry = max_y - min_y
            r = max(rx, ry) / 2 + 1
            cx, cy = (min_x + max_x) / 2, (min_y + max_y) / 2
            s_lt = cx - r, cy - r
            s_rb = cx + r, cy + r

        d_lt = 50, 50
        d_rb = 800, 800

        self.s_lt = s_lt
        self.s_rb = s_rb

        self.d_lt = d_lt
        self.d_rb = d_rb

        kx, ky = (
            (d_rb[0] - d_lt[0]) / (s_rb[0] - s_lt[0]),
            (d_rb[1] - d_lt[1]) / (s_rb[1] - s_lt[1]),
        )
        self.tran_x = lambda x: d_lt[0] + (x - s_lt[0]) * kx
        self.tran_y = lambda y: d_lt[1] + (y - s_lt[1]) * ky

        self.inv_tran_x = lambda x: (x - d_lt[0]) / kx + s_lt[0]
        self.inv_tran_y = lambda y: (y - d_lt[1]) / ky + s_lt[1]

    def mousePressEvent(self, e):
        super().mousePressEvent(e)
        self.repaint()

        if e.button() == 1:  # left button
            print(e.pos())
            s_x = int(np.round(self.inv_tran_x(e.pos().x())))
            s_y = int(np.round(self.inv_tran_y(e.pos().y())))
            self.set_coord_queue.put((s_x, s_y))

        if e.button() == 2:  # right button
            self.rewind_queue.put(True)

    def paintEvent(self, e):
        super().paintEvent(e)
        painter = QPainter(self)
        while not self.input_queue.empty():
            draw_data = self.input_queue.get(block=False)
            self.all_points = draw_data['draw_data']
            self.set_scale()
            # self.special_points = draw_data['special_points']

        painter.setPen(QColor(0, 0, 0, 0))

        painter.setBrush(QColor(0, 0, 0, 255))

        d_lt = self.d_lt
        d_rb = self.d_rb

        painter.drawRect(
            d_lt[0] - 10, d_lt[1] - 10, d_rb[0] - d_lt[0] + 20, d_rb[1] - d_lt[1] + 20
        )

        painter.setBrush(QColor(240, 240, 240, 100))
        painter.drawRect(d_lt[0], d_lt[1], d_rb[0] - d_lt[0], d_rb[1] - d_lt[1])

        pallet = [QColor(255, 0, 0, 150), QColor(0, 255, 0, 150), QColor(0, 0, 255, 150),
                  QColor(255, 0, 255, 150), QColor(0, 255, 255, 150), QColor(0, 0, 255, 150)]

        # figure
        for point in self.all_points:
            color = point[2]
            if color >= len(pallet):
                print('color out of range!!!!!!!!')
                color = color % len(pallet)
            painter.setBrush(pallet[color])
            painter.drawRect(int(self.tran_x(point[0])) - 2, int(self.tran_y(point[1])) - 2, 4, 4)

        # # state changing points
        # painter.setBrush(QColor(255, 255, 255, 255))
        # for point in self.special_points:
        #     painter.drawEllipse(int(self.tran_x(point[0])) - 3, int(self.tran_y(point[1])) - 3, 6, 6)

        # coordinate grid
        # painter.setBrush(QColor(80, 80, 80, 255))
        # for x in range(-10, 11, 1):
        #     for y in range(-10, 11, 1):
        #         xs = int(self.tran_x(x) + 0.5)
        #         ys = int(self.tran_y(y) + 0.5)
        #         delta = 0
        #         if x % 5 == 0 and y % 5 == 0:
        #             delta = 2
        #         painter.drawEllipse(xs - 2 - delta, ys - 2 - delta, 4 + 2 * delta, 4 + 2 * delta)

        # origin point
        # xs = int(self.tran_x(0) + 0.5)
        # ys = int(self.tran_y(0) + 0.5)

        # painter.setBrush(QColor(255, 255, 0, 255))
        # painter.drawEllipse(xs - 4, ys - 4, 8, 8)

        painter.end()


def main():
    app = QApplication([])
    window = MainWindow()

    window.show()
    app.exec_()
    window.stop_queue.put("stop")


if __name__ == "__main__":
    main()
