from PyQt6.QtWidgets import (
    QWidget,
    QHBoxLayout,
    QPushButton,
    )

from data.model import Item
from data.app import App
from data.conn import Data
from widgets.Dialogs import error, askdlg, messbox, ok_cansel_dlg
from widgets.Form import ItemTable
from widgets.Table import TableWControls 

class AdminTab(QWidget):
    def __init__(self):
        super().__init__()
        app = App()
        self.repo: Data = app.repository
        box = QHBoxLayout()
        self.setLayout(box)
        data_model = {
            "id": {"def": 0, "hum": "Номер", "form": 0},
            "name": {"def": "", "hum": "Назва", "form": 2},
            "created_at": {"def": "", "hum": "Створена", "form": 1},
        }
        btns = {'reload':'Оновити', 'create':'Створити', 'delete':'Видалити'}
        self.backups = TableWControls(data_model, ['id', 'name', 'created_at'], buttons=btns)
        box.addWidget(self.backups)
        self.backups.actionInvoked.connect(self.action)
        rest_btn = QPushButton('Відновити')
        self.backups.hbox.insertWidget(0, rest_btn)
        rest_btn.clicked.connect(self.restore_base)

    def restore_base(self):
        current = self.backups.table.get_selected_value()
        if not current:
            error('Оберіть збережену базу')
            return
        q = ok_cansel_dlg(f'Ви дійсно хочете відновити базу {current["name"]} за {current["created_at"]}?')
        if q:
            res = self.repo.restore_base(f'{current["name"]}_{current["created_at"]}')
            if res['error']:
                error(res['error'])
                return
            if res['value']:
                print(res['value'])
                messbox(f'База під назвою {current["name"]} успішно відновлена')

    def delete_base(self):
        current = self.backups.table.get_selected_value()
        if not current:
            error('Оберіть збережену базу')
            return
        q = ok_cansel_dlg(f'Ви дійсно хочете видалити базу {current["name"]} за {current["created_at"]}?')
        if q:
            res = self.repo.delete_base(f'{current["name"]}_{current["created_at"]}')
            if res['error']:
                error(res['error'])
                return
            if res['value']:
                print(res['value'])
                messbox(f'База під назвою {current["name"]} успішно видалена')

    def action(self, action_name, value):
        print('action', action_name, value)
        if action_name == 'create':
            name = askdlg('Задайте назву для збереженої бази')
            if name:
                res = self.repo.create_base_backup(name)
                if res['error']:
                    error(res['error'])
                    return
                if res['value']:
                    print(res['value'])
                    messbox(f'База успішно збережена під назвою {name}')
                    self.reload()
        if action_name == 'reload':
            self.reload()
        if action_name == 'delete':
            self.delete_base()
            self.reload()

    def reload(self):
        res = self.repo.get_base_backups()
        print(res)
        if res['error']:
            error(res['error'])
            return
        if not res['value']['base_names']:
            return
        values = []
        for i, v in enumerate(res['value']['base_names']):
            name, date = v.split('_')
            values.append({'id': i, 'name': name, 'created_at': date})
        self.backups.table.reload(values)

