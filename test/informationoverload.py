import army
from threading import Thread

for i in range(2):
    Thread(target=army.main).start()