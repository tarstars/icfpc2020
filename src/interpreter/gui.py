import queue
import random
import time

from PyQt5 import QtCore
from PyQt5.QtCore import QRunnable, QThreadPool
from PyQt5.QtGui import QPainter, QColor
from PyQt5.QtWidgets import QApplication, QLabel, QMainWindow

from computer import Machine


class Worker(QRunnable):
    def __init__(self, queue, stop_queue, *args, **kwargs):
        super().__init__()
        self.queue = queue
        self.stop_queue = stop_queue

    def run(self):
        with open("/home/tass/database/icfpc2020/messages/galaxy.txt") as fd:
            all_lines = [line.strip() for line in fd.readlines()]
            machine = Machine.from_lines(all_lines)
            ans = machine.eval("ap list_debug ap ap ap interact galaxy nil ap ap cons 0 0")
            for point in ans.value[1][0].value:
                self.queue.put(point, block=False)

        # while True:
        #     time.sleep(.1)
        #     new_point = (random.randint(100, 500), random.randint(100, 500))
        #     print('new point', new_point)
        #     self.queue.put(new_point, block=False)
        #     if not self.stop_queue.empty():
        #         break


class MainWindow(QMainWindow):
    def __init__(self):
        super().__init__()
        side = 600
        self.setFixedWidth(side)
        self.setFixedHeight(side)

        self.threadpool = QThreadPool()
        self.input_queue = queue.Queue()
        self.stop_queue = queue.Queue()
        worker = Worker(self.input_queue, self.stop_queue)
        self.threadpool.start(worker)

        self.all_points = []

        self.update_timer = QtCore.QTimer(self)
        self.update_timer.timeout.connect(self.repaint)
        self.update_timer.start(1000)

    def mousePressEvent(self, e):
        print(e.pos())
        super().mousePressEvent(e)
        self.repaint()

    def paintEvent(self, e):
        super().paintEvent(e)
        painter = QPainter(self)
        while not self.input_queue.empty():
            point = self.input_queue.get(block=False)
            self.all_points.append(point)

        s_lt = -7, -7
        s_rb = 7, 7
        d_lt = 100, 100
        d_rb = 500, 500

        kx, ky = (d_rb[0] - d_lt[0]) / (s_rb[0] - s_lt[0]), (d_rb[1] - d_lt[1]) / (s_rb[1] - s_lt[1])
        tran_x = lambda x: d_lt[0] + (x - s_lt[0]) * kx
        tran_y = lambda y: d_lt[0] + (y - s_lt[0]) * ky

        painter.setPen(QColor(0, 0, 0, 0))

        painter.setBrush(QColor(0, 0, 0, 255))
        painter.drawRect(d_lt[0] - 10, d_lt[1] - 10, d_rb[0] - d_lt[0] + 20, d_rb[1] - d_lt[1] + 20)

        painter.setBrush(QColor(0, 255, 0, 80))
        painter.drawRect(d_lt[0], d_lt[1], d_rb[0] - d_lt[0], d_rb[1] - d_lt[1])

        painter.setBrush(QColor(255, 0, 0, 80))
        for point in self.all_points:
            painter.drawEllipse(tran_x(point[0]), tran_y(point[1]), 20, 20)

        painter.end()


def main():
    app = QApplication([])
    window = MainWindow()

    window.show()
    app.exec_()
    window.stop_queue.put('stop')


if __name__ == '__main__':
    main()
