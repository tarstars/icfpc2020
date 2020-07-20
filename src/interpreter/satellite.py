import math

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


class MainWindow(QMainWindow):
    def __init__(self):
        super().__init__()
        side = 800
        self.setFixedWidth(side)
        self.setFixedHeight(side)

        self.update_timer = QtCore.QTimer(self)
        self.update_timer.timeout.connect(self.repaint)
        self.update_timer.start(100)

        self.state = np.array([100, 0, 1, 0])

    def mousePressEvent(self, e):
        super().mousePressEvent(e)

    def paintEvent(self, e):
        super().paintEvent(e)
        painter = QPainter(self)

        sx = int(self.x)
        sy = int(self.y)
        painter.drawEllipse(500 + sx - 5, 500 - sx - 5, 10, 10)
        self.next()
        painter.end()

    def next(self):
        rt_d = self.x
        lt_d = self.y

        dt = 0.01
        G = 1e-4
        sign = lambda x: -1 if x < 0 else 1
        a_rt = - sign(rt_d) * G * rt_d ** 4
        a_lt = - sign(lt_d) * G * lt_d ** 4

        self.x += self.vx * dt
        self.y += self.vy * dt

        self.vx += a_rt * dt
        self.vy += a_lt * dt


def main():
    app = QApplication([])
    window = MainWindow()

    window.show()
    app.exec_()


if __name__ == "__main__":
    main()
