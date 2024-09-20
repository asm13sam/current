import json

from data.app import App
from data.conn import Data
from gui import MainWindow

with open ('config.json', "r") as f:
    cfg = json.loads(f.read())

with open ('models.json', "r") as f:
    models = json.loads(f.read())

app = App()
app.set_params(cfg, models, Data(cfg))
w = MainWindow()
w.run()
