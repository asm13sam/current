from PyQt6.QtCore import Qt
from PyQt6.QtWidgets import (
    QWidget,
    QVBoxLayout,
    QPushButton,
    QGridLayout,
    QLabel,
    QSplitter,
    QComboBox,
    QColorDialog,
    QFileDialog,
    QTabWidget,
    QGroupBox,
    )
from PyQt6.QtGui import QColor

import json

from data.model import Item
from data.app import App
from data.conn import Data
from widgets.Dialogs import error, askdlg, messbox, ok_cansel_dlg, CustomDialog
from widgets.Form import (
    LineEditWidget,
    IntWidget,
    Selector,
    PersentWidget,
    )
from widgets.Table import TableWControls
from data_widgets.Helpers import OrdTree


class PathSelector(QPushButton):
    def __init__(self):
        super().__init__()
        self.path = ''
        self.clicked.connect(self.select_path)

    def set_value(self, value):
        self.path = value
        self.setText(str(value))

    def select_path(self):
        dlg = QFileDialog()
        dir = dlg.getExistingDirectory(directory=self.path)
        if dir:
            self.set_value(dir)

    def value(self):
        return self.path


class ColorSelector(QPushButton):
    def __init__(self):
        super().__init__()
        self.color = ''
        self.clicked.connect(self.select_color)

    def set_value(self, value):
        self.color = value
        self.setStyleSheet(f'color:{self.color}')
        self.setText('███████ ' + str(value))

    def select_color(self):
        dlg = QColorDialog()
        color = dlg.getColor(initial=QColor(self.color))
        if color.isValid():
            self.set_value(color.name())

    def value(self):
        return self.color


class ProductGroupBox(QGroupBox):
    def __init__(self, name, values):
        super().__init__(name)
        self.setStyleSheet('padding: 4px; background: #131511')
        grid = QGridLayout()
        self.setLayout(grid)
        self.setCheckable(True)
        self.setChecked(False)
        
        width = 4
        x = 0
        y = 0
        for v in values:
            b = QPushButton(v['short_name'])
            # b.clicked.connect(lambda _, value=v: callback(value))
            grid.addWidget(b, y, x)
            x += 1
            if x >= width:
                x = 0
                y += 1


class ProductsGroupTab(QWidget):
    def __init__(self, groups: list):
        super().__init__()
        self.box = QVBoxLayout()
        self.setLayout(self.box)

        group = Item('product_group')
        product = Item('product')
        for group_id in groups:
            group_id += 1 # TODO виправити в Orders->ProductsButtonsBlock
            err = group.get(group_id)
            if err:
                error(f'Не знайдена група id={group_id}: {err}')
                continue
            err = product.get_filter('product_group_id', group_id)
            if err:
                error(f'Не знайдені вироби до групи {group.value["name"]} id={group_id}: {err}')
                continue
            self.box.addWidget(ProductGroupBox(group.value['name'], product.values))
        self.box.addStretch()

    def get_checked(self):
        return [i for i in range(self.box.count()-1) if self.box.itemAt(i).widget().isChecked()]
        

class ProductPlacesWidget(QWidget):
    def __init__(self, places: list=[]):
        super().__init__()
        self.places = places
        self.grid = QGridLayout()
        self.setLayout(self.grid)
        self.grid.setContentsMargins(0, 0, 0, 0)
        self.grid.setVerticalSpacing(1)

        add_group_btn = QPushButton('Додати групу до поточної вкладки')
        self.grid.addWidget(add_group_btn, 0, 0)
        add_group_btn.clicked.connect(self.add_group)
        add_tab_btn = QPushButton('Додати вкладку')
        self.grid.addWidget(add_tab_btn, 0, 1)
        add_tab_btn.clicked.connect(self.add_tab)
        del_tab_btn = QPushButton('Видалити поточну вкладку')
        self.grid.addWidget(del_tab_btn, 0, 2)
        del_tab_btn.clicked.connect(self.del_tab)
        del_group_btn = QPushButton('Видалити групу')
        self.grid.addWidget(del_group_btn, 0, 3)
        del_group_btn.clicked.connect(self.del_group)
        self.prod_tree = OrdTree('product', 'Вироби')
        self.tabs = QTabWidget()
        self.grid.addWidget(self.prod_tree, 1, 0, 1, 2)
        self.grid.addWidget(self.tabs, 1, 2, 1, 2)

        if places:
            self.reload()

    def add_group(self):
        v = self.prod_tree.get_selected_value()
        if not v or 'cost' in v:
            error('Оберіть групу!')
        i = self.tabs.currentIndex()
        self.places[i].append(v['id']-1)
        self.reload()

    def add_tab(self):
        self.places.append([])
        self.reload()

    def del_tab(self):
        i = self.tabs.currentIndex()
        self.places.pop(i)
        self.reload()

    def del_group(self):
        index = self.tabs.currentIndex()
        checked = self.tabs.currentWidget().get_checked()
        new_row = []
        for i, p in enumerate(self.places[index]):
            if i not in checked:
                new_row.append(p)
        self.places[index] = new_row
        self.reload()

    def value(self):
        return self.places

    def reload(self, places=[]):
        if places:
            self.places = places
        index = self.tabs.currentIndex()
        self.tabs.clear()
        for i, groups in enumerate(self.places):
            title = 'Головна' if i == 0 else f'Додаток {i}'
            w = ProductsGroupTab(groups)
            self.tabs.addTab(w, title)
        self.tabs.setCurrentIndex(index)
        

class ProductPlacesConfig(QPushButton):
    def __init__(self):
        super().__init__()
        self.places = []
        self.clicked.connect(self.select_places)

    def set_value(self, value):
        self.places = value
        self.setText(f'Вкладок: {len(value)}')

    def select_places(self):
        w = ProductPlacesWidget(self.places)
        dlg = CustomDialog(w, 'Редагувати розташування')
        res = dlg.exec()
        if res:
            self.set_value(w.value())

    def value(self):
        return self.places
        

class ThemeSelector(QComboBox):
    def __init__(self):
        super().__init__()
        self.addItem('Темна')
        self.addItem('Світла')

    def set_value(self, value):
        self.setCurrentText('Темна' if value == 'dark' else 'Світла')

    def value(self):
        print('self.currentText', self.currentText)
        return 'dark' if self.currentText() == 'Темна' else 'light'


class ConfigEditor(QWidget):
    def __init__(self) -> None:
        super().__init__()
        self.cfg = {}
        self.hum = {}
        self.widgets = {}
        self.grid = QGridLayout()
        self.setLayout(self.grid)
        self.grid.setContentsMargins(0, 0, 0, 0)
        self.grid.setVerticalSpacing(1)

    def make_gui(self):
        self.make_widgets()
        self.make_form()

    def read_config(self):
        with open ('config.json', "r") as f:
            self.cfg = json.loads(f.read())
        with open ('config_hum.json', "r") as f:
            self.hum = json.loads(f.read())

    def make_widgets(self):
        self.widgets["host"] = LineEditWidget()
        self.widgets["port"] = LineEditWidget()

        self.widgets["makets_path"] = PathSelector()
        self.widgets["bs_makets_path"] = PathSelector()
        self.widgets["new_makets_path"] = PathSelector()
        self.widgets["program"] = LineEditWidget()

        self.widgets["theme"] = ThemeSelector()
        self.widgets["font_size"] = IntWidget()
        self.widgets["color"] = ColorSelector()
        self.widgets["form_names_color"] = ColorSelector()
        self.widgets["form_values_color"] = ColorSelector()
        self.widgets["form_bg_color"] = ColorSelector()
        self.widgets["info_names_color"] = ColorSelector()
        self.widgets["info_values_color"] = ColorSelector()
        self.widgets["info_bg_color"] = ColorSelector()

        self.widgets["license_key"] = LineEditWidget()
        self.widgets["cashier_pin"] = LineEditWidget()
        self.widgets["checkbox_url"] = LineEditWidget()

        self.widgets["warehouse persent"] = PersentWidget()
        self.widgets["copycenter warehouse id"] = Selector('whs')
        self.widgets["measure pieces"] = Selector('measure')
        self.widgets["measure linear"] = Selector('measure')
        self.widgets["measure square"] = Selector('measure')
        self.widgets["contragent to production"] = Selector('contragent')
        self.widgets["contact to production"] = Selector('contact')
        self.widgets["contragent copycenter default"] = Selector('contragent')
        self.widgets["contact copycenter default"] = Selector('contact')
        self.widgets["ordering state in work"] = Selector('ordering_status')
        self.widgets["ordering state ready"] = Selector('ordering_status')
        self.widgets["ordering state taken"] = Selector('ordering_status')
        self.widgets["ordering state canceled"] = Selector('ordering_status')
        self.widgets["product_to_ordering state ready"] = Selector('ordering_status')

        self.widgets["product_groups"] = ProductPlacesConfig()

    def make_form(self):
        row = 0
        for k, v in self.cfg.items():
            self.grid.addWidget(QLabel(self.hum[k]), row, 0)
            self.grid.addWidget(self.widgets[k], row, 1)
            self.widgets[k].set_value(v)
            row += 1
        
        cancel_btn = QPushButton('Скинути')
        self.grid.addWidget(cancel_btn, row, 0)
        cancel_btn.clicked.connect(self.restore_config)
        save_btn = QPushButton('Зберегти')
        self.grid.addWidget(save_btn, row, 1)
        save_btn.clicked.connect(self.save_config)

    def restore_config(self):
        for k, v in self.cfg.items():
            self.widgets[k].set_value(v)
            
    def save_config(self):
        new_cfg = {}
        for k in self.cfg:
            new_cfg[k] = self.widgets[k].value()
        with open ('prew_config.json', "w") as f:
            f.write(json.dumps(self.cfg))
        with open ('config.json', "w") as f:
            f.write(json.dumps(new_cfg))
        self.cfg = new_cfg


class AdminTab(QSplitter):
    def __init__(self):
        super().__init__(Qt.Orientation.Horizontal)
        app = App()
        self.repo: Data = app.repository
        data_model = {
            "id": {"def": 0, "hum": "Номер", "form": 0},
            "name": {"def": "", "hum": "Назва", "form": 2},
            "created_at": {"def": "", "hum": "Створена", "form": 1},
        }
        btns = {'reload':'Оновити', 'create':'Створити', 'delete':'Видалити'}
        self.backups = TableWControls(data_model, ['id', 'name', 'created_at'], buttons=btns)
        self.addWidget(self.backups)
        self.backups.actionInvoked.connect(self.action)
        rest_btn = QPushButton('Відновити')
        self.backups.hbox.insertWidget(0, rest_btn)
        rest_btn.clicked.connect(self.restore_base)
        config_editor = ConfigEditor()
        self.addWidget(config_editor)
        config_editor.read_config()
        config_editor.make_gui()
        self.setStretchFactor(0, 10)
        self.setStretchFactor(1, 1)

    def restore_base(self):
        current = self.backups.table.get_selected_value()
        if not current:
            error('Оберіть збережену базу')
            return
        q = ok_cansel_dlg(f'Ви дійсно хочете відновити базу {current["name"]} за {current["created_at"]}?')
        if not q:
            return
        res = self.repo.restore_base(f'{current["name"]}_{current["created_at"]}')
        if res['error']:
            error(res['error'])
            return
        messbox(f'База під назвою {current["name"]} успішно відновлена')

    def delete_base(self):
        current = self.backups.table.get_selected_value()
        if not current:
            error('Оберіть збережену базу')
            return False
        q = ok_cansel_dlg(f'Ви дійсно хочете видалити базу {current["name"]} за {current["created_at"]}?')
        if not q:
            return False
        res = self.repo.delete_base(f'{current["name"]}_{current["created_at"]}')
        if res['error']:
            error(res['error'])
            return False
        messbox(f'База під назвою {current["name"]} успішно видалена')
        return True

    def create(self):
        name = askdlg('Задайте назву для збереженої бази')
        if not name:
            return False
        name = name.strip().replace(' ', '-').replace('_', '-')
        res = self.repo.create_base_backup(name)
        if res['error']:
            error(res['error'])
            return
        messbox(f'База успішно збережена під назвою {name}')
        return True
    
    def action(self, action_name, value):
        result = False
        if action_name == 'create':
            result = self.create()                
        if action_name == 'delete':
            result = self.delete_base()
        if action_name == 'reload' or result:
            self.reload()

    def reload(self):
        res = self.repo.get_base_backups()
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


