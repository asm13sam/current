import json
import sqlite3
import shutil
import os
import sys

import create_srvr
import create_sql

with open ('changes.json', "r") as f:
    changes = json.loads(f.read())

with open ('models_bk.json', "r") as f:
    model_bk = json.loads(f.read())

with open ('models.json', "r") as f:
    model = json.loads(f.read())

tables = list(model['models'].keys())

if len(sys.argv) == 1 or sys.argv[1] == 'tables':
    shutil.copyfile('../target_cc', 'base.db')
    shutil.copyfile('../base.db', 'update.db')
    shutil.copyfile('base.db', 'base.bk')

    con_from = sqlite3.connect('base.bk')
    con_from.row_factory = sqlite3.Row
    con_update = sqlite3.connect('update.db')
    con_update.row_factory = sqlite3.Row
    con_to = sqlite3.connect('base.db')
    cur_from = con_from.cursor()
    cur_update = con_update.cursor()
    cur_to = con_to.cursor()

    sql_creator = create_sql.SqlCreator(changes, model, model_bk)
    sql_creator.recreate_tables(cur_to, cur_from, cur_update)

    con_to.commit()
    shutil.copyfile('base.db', '../server/base.db')
    shutil.copyfile('models.json', '../server/models.json')

if len(sys.argv) == 1 or sys.argv[1] == 'srv':
    create_srvr.create_go_models(model, tables)

    os.chdir('../golang')
    os.system('go fmt .')
    os.system('go build')

    shutil.copyfile('tgtsrv.exe', '../server/tgtsrv.exe')

print("All done")
