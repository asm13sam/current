import json
from PyQt6 import QtPrintSupport
from PyQt6.QtCore import (
    Qt,
    pyqtSignal,
    )
from PyQt6.QtGui import QKeySequence, QShortcut, QPainter, QStandardItem
from PyQt6.QtWidgets import (
    QVBoxLayout,
    QLabel,
    QPushButton,
    QWidget,
    QLineEdit,
    QFormLayout,
    QTabWidget,
    QHBoxLayout,
    QSplitter,
    QGridLayout,
    QGroupBox,
    
    )

import time
from datetime import datetime, date
from data.model import Item
from data.app import App
from data_widgets.CheckBoxUa import CheckBoxFS
from widgets.Dialogs import error, CustomDialog, ok_cansel_dlg, messbox
from widgets.Table import Table, TableModel
from widgets.Form import Selector, ContactSelector, TimeWidget, CustomFormDialog
from common.params import SORT_ROLE
from data_widgets.Calculation import ProductExtra, ProductView, MatherialExtra, OperationExtra
from data_widgets.ItemsToOrdering import ProductToOrderingForm
from data_widgets.Helpers import (
    OrdTree, 
    MatherialToOrderingForm, 
    OperationToOrderingForm, 
    ProductToOrderingForm,
    OrderingForm,
    )
from common.funcs import fake_id

class PMOTableModel(TableModel):
    def __init__(self, data_model: dict, field_names: list = []):
        super().__init__(data_model, field_names)

    def make_item(self, value, name):
        if 'product_extra' in value:
            if name == 'name':
                v = value['product_extra']['product']['name']
            elif name == 'comm':
                v = value['product_extra']['product_to_ordering']['info']
            elif name == 'type':
                v = 'product_to_ordering'
            else:
                v = value['product_extra']['product_to_ordering'][name]
        elif 'matherial' in value:
            if name == 'name':
                v = value['matherial']['name']
            elif name == 'type':
                v = 'matherial_to_ordering'
            else:
                v = value['matherial_to_ordering'][name]
        else:
            if name == 'name':
                v = value['operation']['name']
            elif name == 'type':
                v = 'operation_to_ordering'
            elif name == 'persent' or name == 'profit':
                v = 0.0
            else:
                v = value['operation_to_ordering'][name]
        
        if type(v) == float:
            v = round(v, 2)
        item = QStandardItem(str(v))
        item.setData(v, SORT_ROLE)
        item.setEditable(False)
        return item
    
    
class PMOTable(Table):
    def __init__(self):
        data_model = {
            "name": {"def": "", "hum": "Назва"},
            "number": {"def": 0.0, "hum": "Кількість"},
            "price": {"def": "date", "hum": "Ціна"},
            "cost": {"def": 0.0, "hum": "Вартість"},
            "persent": {"def": 0.0, "hum": "Націнка,%"},
            "profit": {"def": 0.0, "hum": "Націнка,грн."},
            "comm": {"def": "", "hum": "Коментар"},
            "type": {"def": "", "hum": "Тип"}
        }
        fields = ["name", "number", "price", "cost", "persent", "profit", "comm"]
        self.sum_field_num = 3 # cost - field to order sum recalc
        super().__init__(data_model, table_fields=fields)
        self._model = PMOTableModel(data_model, fields)
        self.setModel(self._model)
        self.dataset = {}

    def add_item(self, item: ProductExtra | MatherialExtra | OperationExtra):
        self._model.append(item.value)
        
    def clear(self):
        self._model.clear()
        self.dataset = {}

    def recalc(self):
        order_sum = 0
        for i in range(self._model.rowCount()):
            order_sum += self._model.item(i, self.sum_field_num).data(SORT_ROLE)
        return order_sum


class OrderTableModel(TableModel):
    def __init__(self, data_model: dict, field_names: list = []):
        super().__init__(data_model, field_names)

    def make_item(self, value, name):
        if 'product_extra' in value:
            v = value['product_extra']['product_to_ordering'][name]
        elif 'matherial' in value:
            if name == 'product':
                v = value['matherial_to_ordering']['matherial']
            elif name == 'info':
                v = value['matherial_to_ordering']['comm']
            else:
                v = value['matherial_to_ordering'][name]
        if type(v) == float:
            v = round(v, 2)
        item = QStandardItem(str(v))
        item.setData(v, SORT_ROLE)
        item.setEditable(False)
        return item
    
    
class OrderTable(Table):
    def __init__(self, data_model: dict, sum_field_num: int, fields: list = []):
        super().__init__(data_model, table_fields=fields)
        self._model = OrderTableModel(data_model, fields)
        self.setModel(self._model)
        self.sum_field_num = sum_field_num
        self.dataset = {}

    def add_item(self, item: ProductExtra | MatherialExtra):
        self._model.append(item.value)
        
    def clear(self):
        self._model.clear()
        self.dataset = {}

    def recalc(self):
        order_sum = 0
        for i in range(self._model.rowCount()):
            order_sum += self._model.item(i, self.sum_field_num).data(SORT_ROLE)
        return order_sum


class CashForm(QWidget):
    def __init__(self, summa: float):
        super().__init__()
        self.total = summa
        self.form = QFormLayout()
        self.setLayout(self.form)
        self.cash = QLineEdit()
        self.summa = QLabel()
        self.summa.setText(str(summa))
        self.left = QLabel()
        self.form.addRow("До сплати", self.summa)
        self.form.addRow("Внесено", self.cash)
        self.form.addRow("Решта", self.left)

        self.cash.textChanged.connect(self.calc_left)

    def calc_left(self):
        try:
            self.left.setText(str(float(self.cash.text()) - self.total))
        except:
            return
         

class CashQRForm(QWidget):
    def __init__(self, summa: float):
        super().__init__()
        self.form = QFormLayout()
        self.setLayout(self.form)
        self.uid = QLineEdit()
        self.summa = QLabel()
        self.summa.setText(str(summa))
        self.form.addRow("До сплати", self.summa)
        self.form.addRow("Код транзакції", self.uid)


class CashPlusForm(QWidget):
    def __init__(self, summa: float):
        super().__init__()
        self.total = summa
        self.form = QFormLayout()
        self.setLayout(self.form)
        self.cash = QLineEdit()
        self.qr = QLineEdit()
        self.trans_uid = QLineEdit()
        self.summa = QLabel()
        self.summa.setText(str(summa))
        self.left = QLabel()
        self.form.addRow("До сплати", self.summa)
        self.form.addRow("Карта", self.qr)
        self.form.addRow("Код транзакції", self.trans_uid)
        self.form.addRow("Готівка", self.cash)
        self.form.addRow("Решта", self.left)

        self.cash.textChanged.connect(self.calc_left)
        self.qr.textChanged.connect(self.calc_left)

    def calc_left(self):
        try:
            paid = float(self.cash.text()) + float(self.qr.text()) - self.total
            self.left.setText(str(paid))
        except:
            return


class PMODialog(CustomFormDialog):
    def __init__(self, item, form, title):
        self.item = item
        self.form = form
        self.extend_form()
        super().__init__(title, self.form)
        self.form.widgets['number'].setFocus()

    def extend_form(self):
        rows = self.form.grid.rowCount()
        self.cash = QLineEdit()
        self.form.grid.addWidget(self.cash, rows, 1)
        cash_lbl = QLabel('Внесено')
        self.form.grid.addWidget(cash_lbl, rows, 0)
        
        self.rest = QLabel('0.00')
        self.form.grid.addWidget(self.rest, rows+1, 1)
        self.form.grid.addWidget(QLabel('Решта'), rows+1, 0)
        
        self.cash.textChanged.connect(self.cash_changed)
        
        for k in [f'{self.item.name}_id', 'price', 'cost']:
            self.form.widgets[k].setDisabled(True)
        if 'persent' in self.form.widgets:
            self.show_btn = QPushButton('[+]')
            self.show_btn.setCheckable(True)
            self.form.grid.addWidget(self.show_btn, 0, 2)
            self.show_btn.toggled.connect(self.show_persent)        
            self.form.widgets['persent'].setVisible(False)
            self.form.widgets['profit'].setVisible(False)
            self.form.labels['persent'].setVisible(False)
            self.form.labels['profit'].setVisible(False)
        if "equipment_cost" in self.form.widgets:
            self.form.widgets['user_sum'].setVisible(False)
            self.form.widgets['equipment_cost'].setVisible(False)
            self.form.labels['user_sum'].setVisible(False)
            self.form.labels['equipment_cost'].setVisible(False)
        style = '''
            text-align: center;
            border-style: outset;
            border-width: 2px;
            border-radius: 10px;
            border-color: #221115;
            font: bold 18px;
            padding: 6px;
        '''
        cost_style = style + '\nbackground-color: yellow;\ncolor: #992233;'     
        input_style = style + '\ncolor: yellow;'     
        self.form.widgets['cost'].setStyleSheet(cost_style)
        self.form.widgets['number'].setStyleSheet(input_style)
        self.cash.setStyleSheet(input_style)

    def show_persent(self):
        if self.show_btn.isChecked():
            self.form.widgets['persent'].setVisible(True)
            self.form.widgets['profit'].setVisible(True)
            self.form.labels['persent'].setVisible(True)
            self.form.labels['profit'].setVisible(True)
        else:
            self.form.widgets['persent'].setVisible(False)
            self.form.widgets['profit'].setVisible(False)
            self.form.labels['persent'].setVisible(False)
            self.form.labels['profit'].setVisible(False)
        
    def cash_changed(self, txt):
        try:
            num = float(txt)
        except:
            error('Неправильний формат суми')
            return
        cost = self.form.widgets['cost'].value()
        rest = num - cost
        self.rest.setText(str(rest))


class ProductFormDialog(PMODialog):
    def __init__(self, prod_view, product):
        self.pw = prod_view
        fields = ['product_id', 'info', 'price', 'persent', 
                  'profit', 'cost', 'width', 'length', 
                  'pieces', 'number']
        form = ProductToOrderingForm(
            fields=fields, 
            value=product.value['product_extra']['product_to_ordering']
        )
        super().__init__(product, form, product.value['product_extra']['product']['name'])
        
    def extend_form(self):
        super().extend_form()  
        self.form.grid.removeWidget(self.form.widgets['info'])
        self.form.widgets['info'].setFixedHeight(50)  
        rows = self.form.grid.rowCount()
        self.form.grid.addWidget(QLabel("Коментар"), rows, 0)
        self.form.grid.addWidget(self.form.widgets['info'], rows, 1)
        self.form.grid.addWidget(self.pw, 0, 3, 20, 1)
        self.pw.productChanged.connect(self.product_changed)
        self.form.sizeChanged.connect(self.form_size_changed)
        
    def product_changed(self):
        self.pw.reload()
        persent = self.form.widgets['persent'].value()
        self.form.reload(self.item.value['product_extra']['product_to_ordering'])
        self.form.widgets['persent'].set_value(persent)
        self.form.persent_changed()

    def form_size_changed(self, width, length, pieces, number):
        p2o = self.item.value['product_extra']['product_to_ordering']
        p2o['width'] = width
        p2o['length'] = length
        p2o['pieces'] = pieces
        p2o['number'] = number
        self.item.recalc_num()
        total, *_ = self.item.recalc()
        if total is not None:
            self.form.reload(self.item.value['product_extra']['product_to_ordering'])
            # self.form.widgets['cost'].setValue(round(total, 1))
            # self.form.persent_changed()


class MatherialFormDialog(PMODialog):
    def __init__(self, matherial):
        fields = ["matherial_id", "width", "length", "pieces", 
                      "color_id", "comm", "price", "persent", 
                      "profit", "number", "cost"]
        form = MatherialToOrderingForm(
            fields=fields, 
            value=matherial.value['matherial_to_ordering']
            )
        super().__init__(matherial, form, matherial.value['matherial']['name'])    
        
        
class OperationFormDialog(PMODialog):
    def __init__(self, operation):
        fields = ["operation_id", "comm", "price", "number", "equipment_cost", "user_sum", "cost"]
        form = OperationToOrderingForm(
            fields=fields, 
            value=operation.value['operation_to_ordering']
            )
        super().__init__(operation, form, operation.value['operation']['name'])    
        
        
class CasheQRFormDialog(CustomDialog):
    def __init__(self, summa):
        self.widget = CashQRForm(summa)
        title = 'Сплатити через QRcode'
        super().__init__(self.widget, title)

    def accept(self) -> None:
        if not self.widget.uid.text():
            return
        return super().accept()
    

class CashePlusFormDialog(CustomDialog):
    def __init__(self, summa):
        self.summa = summa
        self.widget = CashPlusForm(summa)
        title = 'Сплатити готівкою + QRcode'
        super().__init__(self.widget, title)

    def accept(self) -> None:
        try:
            cash_sum = float(self.widget.cash.text())
            qr_sum = float(self.widget.qr.text())
        except:
            error("Неправильний формат суми")
            return
        if cash_sum + qr_sum < self.summa:
            error('Недостатньо коштів')
            return
        return super().accept()


class CasheFormDialog(CustomDialog):
    def __init__(self, summa):
        self.summa = summa
        self.widget = CashForm(summa)
        title = 'Сплатити'
        super().__init__(self.widget, title)

    def accept(self) -> None:
        try:
            summa = float(self.widget.cash.text())
        except:
            error("Неправильний формат суми")
            return
        if summa < self.summa:
            error('Недостатньо коштів')
            return
        return super().accept()


class ProductGroupBox(QGroupBox):
    def __init__(self, name, values, callback):
        super().__init__(name)
        self.setStyleSheet('padding: 4px; background: #131511')
        grid = QGridLayout()
        self.setLayout(grid)
        
        width = 4
        x = 0
        y = 0
        for v in values:
            b = QPushButton(v['short_name'])
            b.clicked.connect(lambda _, value=v: callback(value))
            grid.addWidget(b, y, x)
            x += 1
            if x >= width:
                x = 0
                y += 1


class ProductsButtonsBlock(QWidget):
    p2oCreated = pyqtSignal(ProductExtra, str, float)
    def __init__(self):
        super().__init__()
        
        self.box = QVBoxLayout()
        self.setStyleSheet('padding: 1px')
        self.setLayout(self.box)
        
        app = App()
        self.groups = app.config['product_groups']
        self.item = Item('product')
        self.prod_groups = Item('product_group')
        
        err = self.prod_groups.get_all()
        if err:
            error(err)
            return

        for i in self.groups[0]:
            values = self.get_values(self.prod_groups.values[i]['id'])
            if values:
                group = ProductGroupBox(self.prod_groups.values[i]['name'], values, self.create_value)
                self.box.addWidget(group)

        tabs = QTabWidget()
        self.box.addWidget(tabs)

        for i in range(1, len(self.groups)):
                tab = QWidget()
                tl = QVBoxLayout()
                tab.setLayout(tl)
                for j in self.groups[i]:
                    values = self.get_values(self.prod_groups.values[j]['id'])
                    if not values:
                        continue
                    group = ProductGroupBox(self.prod_groups.values[j]['name'], values, self.create_value)
                    tl.addWidget(group)
                tl.addStretch()
                tabs.addTab(tab, f'Додаток {i}')

        if j+1 == len(self.prod_groups.values):
            return

        tab = QWidget()
        tl = QVBoxLayout()
        tab.setLayout(tl)

    def get_values(self, prod_group_id):
        err = self.item.get_filter_w('product_group_id', prod_group_id)
        if err:
            error(err)
            return
        return self.item.values

    def create_value(self, product_value):
        product = ProductExtra()
        if not product.set_value(product_value['id'], 0):
            return
        product.value['product_extra']['product_to_ordering']['name'] = product.value['product_extra']['product']['name']
        pw = ProductView(product)
        pw.reload()
        pw.product.recalc_num()
        pw.product.recalc()
        dlg = ProductFormDialog(pw, product)
        res = dlg.exec()
        if res:
            cash = dlg.cash.text()
            summa = dlg.form.widgets['cost'].value()
            self.p2oCreated.emit(product, cash, summa)


class PMOSelector(QTabWidget):
    p2oCreated = pyqtSignal(ProductExtra, str, float)
    m2oCreated = pyqtSignal(MatherialExtra, str, float)
    o2oCreated = pyqtSignal(OperationExtra, str, float)
    def __init__(self):
        super().__init__()
        self.current_value = {} #
        self.current_item = {} #
        
        self.prod_tree = OrdTree('product', 'Вироби')
        self.mat_tree = OrdTree('matherial', 'Матеріали')
        self.op_tree = OrdTree('operation', 'Операції')
        self.addTab(self.prod_tree, 'Вироби')
        self.addTab(self.mat_tree, 'Матеріали')
        self.addTab(self.op_tree, 'Операції')
        
        self.prod_tree.valueDoubleCklicked.connect(lambda v: self.show_form_dlg(v, 'product'))
        self.mat_tree.valueDoubleCklicked.connect(lambda v: self.show_form_dlg(v, 'matherial'))
        self.op_tree.valueDoubleCklicked.connect(lambda v: self.show_form_dlg(v, 'operation'))
        
    def reload(self):
        self.reload_selects()
    
    def reload_selects(self):
        self.prod_tree.reload()
        self.mat_tree.reload()
        self.op_tree.reload()

    def show_form_dlg(self, item_value, item_type):
        if not 'cost' in item_value:
            return
        self.pw = None
        if item_type == 'product':
            product = ProductExtra()
            if not product.set_value(item_value['id'], 0):
                return
            product.value['product_extra']['product_to_ordering']['name'] = product.value['product_extra']['product']['name']
            pw = ProductView(product)
            pw.reload()
            pw.product.recalc_num()
            pw.product.recalc()
            dlg = ProductFormDialog(pw, product)
            res = dlg.exec()
            if res:
                cash = dlg.cash.text()
                summa = dlg.form.widgets['cost'].value()
                self.p2oCreated.emit(product, cash, summa)
        # --    
        if item_type == 'matherial':
            matherial = MatherialExtra()
            matherial.set_value(item_value)
            dlg = MatherialFormDialog(matherial)
            res = dlg.exec()
            if res:
                cash = dlg.cash.text()
                summa = dlg.form.widgets['cost'].value()
                self.m2oCreated.emit(matherial, cash, summa)
            
        if item_type == 'operation':
            operation = OperationExtra()
            operation.set_value(item_value)
            dlg = OperationFormDialog(operation)
            res = dlg.exec()
            if res:
                cash = dlg.cash.text()
                summa = dlg.form.widgets['cost'].value()
                self.o2oCreated.emit(operation, cash, summa)
        

class ProductsTab(QWidget):
    def __init__(self, check: CheckBoxFS):
        super().__init__()
        self.item = Item('product')
        self.summa = 0.0
        self.order_counter = 0
        self.state = {"Поточне": []}
        self.check = check
        app = App()
        self.user = app.user
        self.current_ordering = None
        self.p2o_ready_status = app.config['product_to_ordering state ready']
        
        self.main_layout = QVBoxLayout()
        self.setLayout(self.main_layout)

        controls_widget = QWidget()
        self.controls = QHBoxLayout()
        controls_widget.setLayout(self.controls)
        self.main_layout.addWidget(controls_widget, stretch=0)
        
        self.start_btn = QPushButton('Відкрити зміну')
        self.start_btn.clicked.connect(self.start)
        self.controls.addWidget(self.start_btn)
        self.offline_lbl = QLabel("Офлайн")
        self.controls.addWidget(self.offline_lbl)
        reload_btn = QPushButton("Оновити вироби")
        self.controls.addWidget(reload_btn)
        reload_btn.clicked.connect(self.reload_selectors)
        self.new_order_name = QLineEdit()
        self.controls.addWidget(self.new_order_name)
        self.new_order_name.returnPressed.connect(self.add_tab)
        self.new_order_name.setMaximumWidth(100)
        new_btn = QPushButton("Нова вкладка")
        new_btn.clicked.connect(self.add_tab)
        self.controls.addWidget(new_btn)
        to_ord_btn = QPushButton("Виділити в замовлення")
        to_ord_btn.clicked.connect(self.to_ordering)
        self.controls.addWidget(to_ord_btn)
        
        self.controls.addStretch()
        
        del_btn = QPushButton("Видалити")
        del_btn.clicked.connect(self.del_products)
        self.controls.addWidget(del_btn)
        self.sum_lbl = QLabel()
        self.sum_lbl.setText(f'Загалом {self.summa} грн.')
        style = '''
            background-color: yellow;
            color: #992233;
            text-align: center;
            border-style: outset;
            border-width: 2px;
            border-radius: 10px;
            border-color: #221115;
            font: bold 18px;
            padding: 6px;
        '''  
        self.sum_lbl.setStyleSheet(style)
        self.controls.addWidget(self.sum_lbl, alignment=Qt.AlignmentFlag.AlignRight)
        pay_btn = QPushButton("Сплатити")
        pay_btn.clicked.connect(self.pay)
        self.controls.addWidget(pay_btn, alignment=Qt.AlignmentFlag.AlignRight)
        pay_btn_qr = QPushButton("QRCode")
        pay_btn_qr.clicked.connect(self.pay_qr)
        self.controls.addWidget(pay_btn_qr, alignment=Qt.AlignmentFlag.AlignRight)
        pay_btn_plus = QPushButton("QR+CASH")
        pay_btn_plus.clicked.connect(self.pay_plus)
        self.controls.addWidget(pay_btn_plus, alignment=Qt.AlignmentFlag.AlignRight)

        self.splitter = QSplitter(Qt.Orientation.Horizontal)
        self.main_layout.addWidget(self.splitter, stretch=10)
        
        self.selector_frame = QTabWidget()
        self.splitter.addWidget(self.selector_frame)
        self.prod_buttons = ProductsButtonsBlock()
        self.selector_frame.addTab(self.prod_buttons, 'Головне')
        self.pmo_selector = PMOSelector()
        self.selector_frame.addTab(self.pmo_selector, "Додатково")
        self.tabs = QTabWidget()
        self.splitter.addWidget(self.tabs)
        self.splitter.setStretchFactor(0, 1)
        self.splitter.setStretchFactor(1, 10)

        self.prod_buttons.p2oCreated.connect(self.add_item_to_order)
        self.pmo_selector.p2oCreated.connect(self.add_item_to_order)
        self.pmo_selector.m2oCreated.connect(self.add_item_to_order)
        self.pmo_selector.o2oCreated.connect(self.add_item_to_order)

        self.tabs.currentChanged.connect(self.recalc_sum)
        self.tabs.tabBarDoubleClicked.connect(self.close_tab)

        for k, v in {'1':1, '2':2}.items():
            sc = QShortcut(QKeySequence('Ctrl+' + k), self)
            i = Item('product')
            err = i.get_w(v)
            if err:
                error(err)
                return
            sc.activated.connect(lambda it=i: self.prod_buttons.create_value(it.value))
        self.load_state()
        self.set_controls_disabled(True)
        
    def reload(self):
        pass

    def reload_selectors(self):
        self.pmo_selector.reload()
        self.selector_frame.removeTab(0)
        self.prod_buttons = ProductsButtonsBlock()
        self.selector_frame.insertTab(0, self.prod_buttons, 'Головне')
        self.prod_buttons.p2oCreated.connect(self.add_item_to_order)
        self.selector_frame.setCurrentIndex(0)

    def set_controls_disabled(self, turn_on: bool):
        for i in range(self.controls.count()-1):
            w = self.controls.itemAt(i+1).widget()
            if w is not None:
                w.setDisabled(turn_on)
        self.pmo_selector.setDisabled(turn_on)
        self.prod_buttons.setDisabled(turn_on)
            
    def start(self):
        app = App()
        if self.check.is_signed:
            res = ok_cansel_dlg("Закриваємо зміну?")
            if not res:
                return
            self.start_btn.setText('Відкрити зміну')
            self.set_controls_disabled(True)
            res = self.check.shift_close()
            if 'error' in res:
                error(res['error'])
                return
            if not res:
                return
            with open ('check.json', "w") as f:
                f.write(json.dumps("{}"))
            self.current_ordering.value['state'] = app.config["ordering state ready"]
            err = self.current_ordering.save()
            if err:
                error(err)
            
        else:
            res = ok_cansel_dlg(f"Відкриваємо зміну?")
            if not res:
                return
            if not self.check.is_signed:
                res = self.check.sign_in()
                if not res:
                    error("Не можу під'єднатися до сервера CheckBox!")
                    return
                if 'error' in res:
                    error(f"При підключенні до Checkbox: {res['error']}")
                    return
            cash_state = self.check.get_cash_state()
            if not cash_state:
                error("Не можу отримати статус каси з сервера CheckBox!")
                return
            if 'error' in cash_state:
                error(f"Отримати статус каси з Checkbox: {cash_state['error']}")
                return    
            self.offline_lbl.setText("Онлайн")
            if cash_state['offline_mode']:
                res = self.check.work_online()
                if not res:
                    error("Не можу перейти в онлайн режим на CheckBox!")
                    self.offline_lbl.setText("Офлайн")
            if cash_state['has_shift']:
                messbox('Зміна вже відкрита, продовжуємо роботу')
                self.set_current_project()
            else:
                shift_data = self.check.shifts()
                if not shift_data:
                    error("Не можу відкрити зміну у CheckBox!")
                    return
                if 'error' in shift_data:
                    error(f"При відкритті зміни у Checkbox: {shift_data['error']}")
                    return 
                self.create_project(shift_data)
            
            self.start_btn.setText('Закрити зміну')
            self.set_controls_disabled(False)
        
    def save_state(self):
        with open ('state.json', "w") as f:
            f.write(json.dumps(self.state))

    def load_state(self):
        with open ('state.json', "r") as f:
            self.state = json.loads(f.read())
        if len(self.state) == 0:
            self.state = {"Поточне": []}
            self.save_state()
        for caption, values in self.state.items():
            products_to_order_list = PMOTable()
            self.tabs.addTab(products_to_order_list, caption)
            for v in values:
                if 'product_extra' in v:
                    item = ProductExtra()
                elif 'matherial' in v:
                    item = MatherialExtra()
                else:
                    item = OperationExtra()
                item.value = v
                products_to_order_list.add_item(item)
        self.recalc_sum()

    def create_project(self, shift_data):
        app = App()
        ordering = Item('ordering')
        ordering.create_default()
        ordering.value['name'] = f'Зміна {shift_data["serial"]}'
        
        ordering.value['user_id'] = self.user['id']
        ordering.value['contragent_id'] = app.config["contragent copycenter default"]
        ordering.value['contact_id'] = app.config["contact copycenter default"]
        ordering.value['info'] = f"id: {shift_data['id']}'\n'serial: {shift_data['serial']}"
        ordering.value['deadline_at'] = date.today().isoformat() + 'T23:59:59'
        err = ordering.save()
        if err:
            error(err)
            return
        err = ordering.create_dirs(ordering.value['id'])
        if err:
            error(err)    
        self.current_ordering = ordering
    
    def set_current_project(self):
        app = App()
        ordering = Item('ordering')
        date_from = date.today().isoformat()
        date_to = date_from + 'T23:59:59'
        err = ordering.get_between_w('created_at', date_from, date_to)
        if err:
            error(err)
            return
        for v in ordering.values[::-1]:
            if v['contragent_id'] == app.config['contragent copycenter default']:
                self.current_ordering = Item('ordering')
                self.current_ordering.value = v
                break
        else:
            shift_data = self.check.get_shift()
            if not shift_data:
                    error("Не можу відкрити зміну у CheckBox!")
                    return
            if 'error' in shift_data:
                error(f"При відкритті зміни у Checkbox: {shift_data['error']}")
                return 
            self.create_project(shift_data)

    def close_tab(self):
        i = self.tabs.currentIndex()
        if i == 0:
            return
        caption = self.tabs.tabText(i)
        del(self.state[caption])
        self.save_state()
        self.tabs.removeTab(i)
        self.recalc_sum()

    def del_products(self):
        caption = self.tabs.tabText(self.tabs.currentIndex())
        selected_rows = self.tabs.currentWidget().get_selected_rows()
        selected_rows.reverse()
        for row in selected_rows:
            del(self.state[caption][row])
        self.save_state()
        self.tabs.currentWidget().delete_values()
        self.recalc_sum()

    # 
    def add_item_to_order(self, item, cash, cost):
        if cash:
            try:
                cash = float(cash)
            except:
                error('Неправильний формат суми')
                return
            if cash >= cost:
                self.create_docs(item, summa=cash)
            else:    
                error('Недостатньо коштів')
            return
        self.tabs.currentWidget().add_item(item)
        caption = self.tabs.tabText(self.tabs.currentIndex())
        self.state[caption].append(item.value)
        self.save_state()
        self.recalc_sum()

    def recalc_sum(self):
        total = self.tabs.currentWidget().recalc()
        self.summa = total
        self.sum_lbl.setText(f'Загалом {self.summa} грн.')

    def clear(self):
        self.tabs.currentWidget().clear()
        self.summa = 0.0
        self.sum_lbl.setText(f'Загалом {self.summa} грн.')
        self.state = {"Поточне": []}
        self.save_state()

    def add_tab(self):
        name = self.new_order_name.text()
        if not name:
            self.order_counter += 1
            name = f'Замовлення {self.order_counter}'
        else:
            self.new_order_name.clear()
        table = PMOTable()
        self.tabs.addTab(table, name)
        self.tabs.setCurrentWidget(table)
        self.state[name] = []
        self.save_state()

    def remove_current_order(self):
        self.tabs.currentWidget().clear()
        i = self.tabs.currentIndex()
        self.summa = 0.0
        self.sum_lbl.setText(f'Загалом {self.summa} грн.')
        caption = self.tabs.tabText(i)
        if i != 0:
            del(self.state[caption])
            self.tabs.removeTab(i)
        else:
            self.state[caption]=[]
            
        self.save_state()

    def pay(self):
        dlg = CasheFormDialog(self.summa)
        res = dlg.exec()
        if res:
            summa = float(dlg.widget.cash.text())
            if self.create_docs(summa=summa):
                self.remove_current_order()
        
    def pay_qr(self):
        dlg = CasheQRFormDialog(self.summa)
        res = dlg.exec()
        if res:
            if self.create_docs(cash_qr=dlg.widget.uid.text(), summa_qr=self.summa):
                self.remove_current_order()

    def pay_plus(self):
        dlg = CashePlusFormDialog(self.summa)
        res = dlg.exec()
        if res:
            trans_uid=dlg.widget.trans_uid.text()
            cash_sum = float(dlg.widget.cash.text())
            qr_sum = float(dlg.widget.qr.text())
            if self.create_docs(cash_qr=trans_uid, summa=cash_sum, summa_qr=qr_sum):
                self.remove_current_order()

    def print_check(self, png):
        pinfo = QtPrintSupport.QPrinterInfo.printerInfo('58mm Series Printer')
        printer = QtPrintSupport.QPrinter(pinfo, QtPrintSupport.QPrinter.PrinterMode.HighResolution)
        painter = QPainter()
        painter.begin(printer)
        png = png.scaledToWidth(printer.width())
        painter.drawPixmap(0, 0, png)
        painter.end()
    
    def create_doc_cbox_check(self, item, summa, summa_qr):
        app = App()
        cbox_check = Item('cbox_check')
        cbox_check.create_default()
        cbox_check.value["user_id"] = self.user['id']
        cbox_check.value["created_at"] = datetime.now().isoformat(timespec='seconds')
        cbox_check.value["contragent_id"] = app.config["contragent copycenter default"]
        cbox_check.value["ordering_id"] = self.current_ordering.value['id']
        cbox_check.value["based_on"] = self.current_ordering.value['document_uid']
        
        receipt = self.create_receipt(item, summa=summa, summa_qr=summa_qr)
        if receipt is None:
            return
        res = self.check.create_receipt(receipt)
        if not res:
            error("Помилка при створенні чеку в Checkbox")
            return
        if 'error' in res:
            error(f"При створенні чеку в Checkbox: {res['error']}")
            return
        cbox_check.value["checkbox_uid"] = res['id']
        cbox_check.value["cash_sum"] = self.summa
        err = cbox_check.save()
        if err:
            error(err)
            return
        else:
            cbox_check.value['name'] = f'Чек {cbox_check.value["id"]}'
            err = cbox_check.save()
            if err:
                error(err)
        if item is None: #multiposition order
            self.create_check_items(cbox_check)
        else: #one position order
            self.create_check_item(item.value, cbox_check)
        return cbox_check

    def create_docs(self, item=None, cash_qr='', summa=0.0, summa_qr=0.0):
        cin_item_qr = None
        cin_item_cash = None
        if item is None: #multiposition order
            items_to_order = self.create_items_to_order()
            if not items_to_order:
                return
            if summa:
                sm = self.summa - summa_qr
                cin_item_cash = self.create_cash_in(summa=sm)
            if summa_qr:
                cin_item_qr = self.create_cash_in(cash_qr=cash_qr, summa=summa_qr)
        else: #one position order
            res = self.create_item_to_order(item.value)
            if not res:
                return
            if item.name == 'product':
                self.summa = res['product_extra']['product_to_ordering']['cost']
            else:
                self.summa = res[f'{item.name}_to_ordering']['cost']
            cin_item_cash = self.create_cash_in(self.summa, cash_qr=cash_qr)
            
        if not cin_item_cash and not cin_item_qr:
            return
        
        self.current_ordering.value['price'] += self.summa
        self.current_ordering.value['cost'] += self.summa
        err = self.current_ordering.save()
        if err:
            error(err)
        
        cbox_check = self.create_doc_cbox_check(item, summa, summa_qr)
        if cbox_check is None:
            return

        if cin_item_cash is not None:
            cin_item_cash.value['cbox_check_id'] = cbox_check.value["id"]
            err = cin_item_cash.save()
            if err:
                error(err)
        if cin_item_qr is not None:
            cin_item_qr.value['cbox_check_id'] = cbox_check.value["id"]
            err = cin_item_qr.save()
            if err:
                error(err)

        png_res = self.get_check_png(cbox_check.value['checkbox_uid'])
        fs_res = self.get_check_fs_code(cbox_check)        
        
        if not png_res:
            png_res = self.get_check_png(cbox_check.value['checkbox_uid'])
            if not png_res:
                error("Не можу отримати png чеку")
        if png_res and not fs_res:
            self.get_check_fs_code(cbox_check)

        return True        
    
    
    def get_check_fs_code(self, cbox_check):
        c = self.check.get_check(cbox_check.value['checkbox_uid'])
        if 'error' in c:
            error(f"При отриманні чеку з Checkbox: {c['error']}")
            return False
        if not c['fiscal_code']:
            error('Не отримали податковий номер')
            return False
        else:
            cbox_check.value["fs_uid"] = c['fiscal_code']
            cbox_check.value["created_at"] = c['fiscal_date']
        err = cbox_check.save()
        if err:
            error(err)
            return False
        return True

    def get_check_png(self, check_uid):
        i = 0
        png = ''
        while i < 3:
            png = self.check.get_png(check_uid)
            time.sleep(2)
            i += 1
            if png:
                break

        if not png:
            return False
        
        if type(png) == str:
            return False
        l = QLabel()
        lpng = png.scaledToHeight(700)
        l.setPixmap(lpng)
        dlg = CustomDialog(l, 'Друкувати чек?')
        if dlg.exec():
            self.print_check(png)
        return True

    def create_receipt(self, item=None, summa=0.0, summa_qr=0.0):
        payments = []
        if summa:
            payments.append({
                "type": 'CASH',
                "value": summa * 100,
                "label": 'Готівка',
            })

        if summa_qr:
            payments.append({
                "type": 'CASHLESS',
                "value": summa_qr * 100,
                "label": 'Картка',
            })
        
        receipt = {
            "department": "Копіцентр",
            "goods": [],
            "payments": payments,
        }
        
        if item is None:
            values_to_order = self.tabs.currentWidget()._model.values()
        else:
            values_to_order = [item.value,]
        
        for v in values_to_order:
            value_type = self.get_type_of_value(v)
            if value_type == 'product':
                value = v['product_extra']['product_to_ordering']
                item_value = v['product_extra']['product']
            else:
                item_value = v[value_type]
                value = v[f'{value_type}_to_ordering']
            if value_type == 'operation':
                price = int(value['cost'] / value['number'] * 100)
            else:
                price = int(value['price'] * 100)
            good = {
                        "good": {
                            "code": item_value['id'],
                            "name": item_value['name'],
                            "price": price,
                        },
                        "quantity": int(value['number']*1000),
                    }
            if 'profit' in value and value['profit']:
                dsc = round(value['profit']*value['number'], 0)
                good["discounts"] = self.make_discount(dsc)
            receipt['goods'].append(good)
        return receipt
    
    def make_discount(self, value):
        if value > 0:
            dsct_type = 'EXTRA_CHARGE'
            dsct_name = "Націнка"
        else:
            value = -value
            dsct_type = 'DISCOUNT'
            dsct_name = "Знижка"
        
        return [
            {
                "type": dsct_type,
                "mode": "VALUE",
                "value": value*100,
                "name": dsct_name
            },
        ]

    def create_used(self, value):
        p2o_id = value['product_extra']['product_to_ordering']['id']
        number = value['product_extra']['product_to_ordering']['number']
        status = value['product_extra']['product_to_ordering']['product_to_ordering_status_id']
        for k, val in value['matherial_extra'].items():
            for v in val:
                if v['matherial_to_product']['is_used']:
                    v['matherial_to_ordering']['ordering_id'] = self.current_ordering.value['id']
                    v['matherial_to_ordering']['number'] = number * v['matherial_to_product']['number']
                    v['matherial_to_ordering']['cost'] = number * v['matherial_to_ordering']['price']
                    v['matherial_to_ordering']['product_to_ordering_id'] = p2o_id
                    m2o = Item('matherial_to_ordering')
                    m2o.value = v['matherial_to_ordering']
                    m2o.value['id'] = 0
                    err = m2o.save()
                    if err:
                        error(err)
        for k, val in value['operation_extra'].items():
            for v in val:
                if v['operation_to_product']['is_used']:
                    v['operation_to_ordering']['ordering_id'] = self.current_ordering.value['id']
                    v['operation_to_ordering']['number'] = number * v['operation_to_product']['number']
                    v['operation_to_ordering']['cost'] = number * v['operation_to_ordering']['price']
                    v['operation_to_ordering']["user_sum"] = number * v['operation']['price'] * v['operation_to_product']['number']
                    v['operation_to_ordering']["equipment_cost"] = number * v['operation_to_product']['equipment_cost']
                    v['operation_to_ordering']['product_to_ordering_id'] = p2o_id
                    o2o = Item('operation_to_ordering')
                    o2o.value = v['operation_to_ordering']
                    o2o.value['id'] = 0
                    err = o2o.save()
                    if err:
                        error(err)
        for k, val in value['product_deep'].items():
            for v in val:
                if v['product_extra']['product_to_product']['is_used']:
                    v['product_extra']['product_to_ordering']['ordering_id'] = self.current_ordering.value['id']
                    v['product_extra']['product_to_ordering']['number'] = number * v['product_extra']['product_to_product']['number']
                    v['product_extra']['product_to_ordering']['cost'] = number * v['product_extra']['product_to_ordering']['price']
                    v['product_extra']['product_to_ordering']['product_to_ordering_id'] = p2o_id
                    v['product_extra']['product_to_ordering']["name"] = v['product_extra']['product']['name']
                    v['product_extra']['product_to_ordering']['product_to_ordering_status_id'] = status
                    p2o = Item('product_to_ordering')
                    p2o.value = v['product_extra']['product_to_ordering']
                    p2o.value['id'] = 0
                    err = p2o.create_p2o_defaults()
                    if err:
                        error(err)
                        return
                    v['product_extra']['product_to_ordering'] = p2o.value
                    self.create_used(v)

    def create_product_to_order(self, product_value, order_id=0):
        if order_id:
            product_value['product_extra']['product_to_ordering']['ordering_id'] = order_id
        else:
            product_value['product_extra']['product_to_ordering']['ordering_id'] = self.current_ordering.value['id']
        product_to_order = Item('product_to_ordering')
        product_to_order.value = product_value['product_extra']['product_to_ordering']
        product_to_order.value['product_to_ordering_status_id'] = self.p2o_ready_status

        err = product_to_order.create_p2o_defaults()
        if err:
            error(err)
            return
        
        product_value['product_extra']['product_to_ordering'] = product_to_order.value
        self.create_used(product_value)
        return product_value
    
    def create_item_to_order(self, value, order_id=0):
        item_type = self.get_type_of_value(value)
        if not item_type:
            return
        if item_type == 'product':
            return self.create_product_to_order(value, order_id)
        if order_id:
            value[f'{item_type}_to_ordering']['ordering_id'] = order_id
        else:
            value[f'{item_type}_to_ordering']['ordering_id'] = self.current_ordering.value['id']
        item_to_order = Item(f'{item_type}_to_ordering')
        item_to_order.value = value[f'{item_type}_to_ordering']
        err = item_to_order.save()
        if err:
            error(err)
            return
        value[f'{item_type}_to_ordering'] = item_to_order.value
        return value

    def create_items_to_order(self, order_id=0):
        item_values = self.tabs.currentWidget()._model.values()
        i2os = []
        for value in item_values:
            res = self.create_item_to_order(value, order_id)
            if res is None:
                item_name = self.get_name_of_value(value)
                error(f"Не вдалося створити {item_name}")
                continue
            i2os.append(res)
        return i2os

    def get_name_of_value(self, value):
        if 'product_extra' in value:
            return value['product_extra']['product']['name']
        if 'matherial' in value:
            return value['matherial']['name']
        if 'operation' in value:
            return value['operation']['name']
        return ''
    
    def get_type_of_value(self, value):
        if 'product_extra' in value:
            return 'product'
        if 'matherial' in value:
            return 'matherial'
        if 'operation' in value:
            return 'operation'
        return ''

    def create_check_item(self, value, check):
        item_type = self.get_type_of_value(value)
        if item_type == 'product':
            item_value = value['product_extra']['product']
            i2o_value = value['product_extra']['product_to_ordering']
        else:
            item_value = value[item_type]
            i2o_value = value[f'{item_type}_to_ordering']
        check_item = Item('item_to_cbox_check')
        check_item.create_default()
        check_item.value["name"] = item_value['name']
        check_item.value["cbox_check_id"] = check.value['id']
        check_item.value["number"] = i2o_value['number']
        check_item.value["measure_id"] = item_value['measure_id']
        check_item.value["price"] = i2o_value['price']
        if 'profit' in i2o_value:
            check_item.value["discount"] = -i2o_value['profit']
        check_item.value["cost"] = i2o_value['cost']
        check_item.value["item_code"] = f"{item_type}-{item_value['id']}"
        err = check_item.save()
        if err:
            error(err)
            return
        return check_item
    
    def create_check_items(self, check):
        values = self.tabs.currentWidget()._model.values()
        cis = []
        for v in values:
            res = self.create_check_item(v, check)
            if res is None:
                item_name = self.get_name_of_value(v)
                error(f"Не вдалося створити у чеку {check.value['id']} {item_name}")
                continue
            cis.append(res)
        return cis

    def create_whs_out(self, products: list):
        app = App()
        wo_item = Item('whs_out')
        wo_item.create_default()
        wo_item.value["name"] = f"ВН по замовленню {self.current_ordering.value['id']}"
        wo_item.value["whs_id"] = app.config["contragent copycenter default"]
        wo_item.value['based_on'] = self.current_ordering.value['document_uid']
        wo_item.value["contragent_id"] = app.config["contragent copycenter default"]
        wo_item.value["contact_id"] = app.config["contact copycenter default"]
        wo_item.value["user_id"] = self.user['id']
        err = wo_item.save()
        if err:
            error(err)
            return
        wo_id = wo_item.value['id']
        summa = 0

        for pe in products:
            for k, val in pe['matherial_extra'].items():
                for v in val:
                    if k == 'default' or v['matherial_to_product']['is_used']:
                        mwo = self.make_matherial_to_whs_out(v, pe['product_extra']['product_to_ordering']['number'], wo_id)
                        err = mwo.save()
                        if err:
                            error(err)
                        summa += mwo.value['cost']
            
        wo_item.value['whs_sum'] = summa
        err = wo_item.save()
        if err:
            error(err)
            return
        return wo_item

    def make_matherial_to_whs_out(self, mat, p2o_number, wo_id):
        mwo = Item('matherial_to_whs_out')
        mwo.create_default()
        mwo.value["matherial_id"] = mat['matherial']['id']
        mwo.value["whs_out_id"] = wo_id
        mwo.value["number"] = p2o_number * mat['matherial_to_product']['number']
        mwo.value["price"] = mat['matherial']['price']
        mwo.value["cost"] = mwo.value["number"] * mwo.value["price"]
        mwo.value["color_id"] = mat['matherial_to_ordering']["color_id"]
        mwo.value["width"] = mat['matherial_to_ordering']["width"]
        mwo.value["length"] = mat['matherial_to_ordering']["length"]
        mwo.value["pieces"] = mat['matherial_to_ordering']["pieces"]
        return mwo
    
    def create_cash_in(self, summa = 0.0, cash_qr=''):
        app = App()
        order_id = self.current_ordering.value['id']
        cin = Item('cash_in')
        cin.create_default()
        cin.value["contragent_id"] = app.config["contragent copycenter default"]
        cin.value["contact_id"] = app.config["contact copycenter default"]
        cin.value["name"] = f"ПКО до Замовлення {order_id}"
        cin.value["cash_id"] = self.user['cash_id']
        cin.value['based_on'] = self.current_ordering.value['document_uid']
        cin.value["user_id"] = self.user['id']
    
        if summa:
            cin.value['cash_sum'] = summa
        else:
            cin.value['cash_sum'] = self.summa
        cin.value['comm'] = 'Auto'
        if cash_qr:
            cash = Item('cash')
            cash.get_all()
            for c in cash.values:
                if c['comm'] == 'qr':
                    cin.value['cash_id'] = c['id']
                    break
            cin.value['comm'] = cash_qr
        
        err = cin.save()
        if err:
            error(err)
            return None
            
        cin.value['name'] = f"ПКО {cin.value['id']} до Замовлення {order_id}"
        err = cin.save()
        if err:
            error(err)
            return
        return cin
    
    def to_ordering(self):
        i = self.tabs.currentIndex()
        caption = self.tabs.tabText(i)
        ordering = Item('ordering')
        ordering.create_default()
        ordering.value['cost'] = self.summa
        ordering.value['price'] = self.summa
        ordering.value['name'] = caption
        ordering.value['user_id'] = self.user['id']
        order_id = self.current_ordering.value['id']
        order_create = self.current_ordering.value['created_at']
        ordering.value['info'] = f'Зі зміни № {order_id} за {order_create}'
        ordering.value['deadline_at'] = date.today().isoformat() + 'T23:59:59'
        form = OrderingForm(value=ordering.value)
        dlg = CustomFormDialog('Нове замовлення', form)
        res = dlg.exec()
        if not (res and dlg.value):
            return
        err = ordering.save()
        if err:
            error(err)
            return
        print(ordering.value['id'])
        self.create_items_to_order(ordering.value['id'])
        err = ordering.create_dirs(ordering.value['id'])
        if err:
            error(err)
            return
        self.remove_current_order() 
        