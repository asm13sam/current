import json

from data.app import App
from data.conn import Data
from gui import MainWindow

with open ('config.json', "r") as f:
    cfg = json.loads(f.read())

# with open ('models.json', "r") as f:
#     models = json.loads(f.read())

project_cfg = {
            "host": "192.168.0.4",
            "port": "8085",
        }
app = App()
app.set_params(cfg, Data(cfg), Data(project_cfg))
w = MainWindow()
w.run()
