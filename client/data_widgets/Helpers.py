from PyQt6.QtCore import (
    Qt,
    pyqtSignal,
    )

from PyQt6.QtWidgets import (
    QWidget,
    QVBoxLayout,
    QComboBox,
    QListWidget,
    QListWidgetItem,
    QSizePolicy,
    )

from widgets.Form import (
    CustomForm,
    InfoBlock,
    )
from data.app import App
from data.model import Item
from widgets.Dialogs import error
from common.params import FULL_VALUE_ROLE
from widgets.Tree import ExtTree


class OrdTree(ExtTree):
    def __init__(self, name: str, title: str = ''):
        self.item = Item(name)
        group_name = name + '_group'
        self.group_item = Item(group_name)
        super().__init__(group_name, title)
        self.reload()
        
    def reload(self):    
        err = self.group_item.get_all_w()
        if err:
            error(err)
            return
        err = self.item.get_all_w()
        if err:
            error(err)
            return
        super().reload(self.group_item.values, self.item.values, self.group_item.name)


class PMOInfo(QWidget):
    def __init__(self, is_info_bottom=False):
        super().__init__()
        info_box = QVBoxLayout()
        self.setLayout(info_box)
        
        columns = 2 if is_info_bottom else 1
        self.product = Item('product_to_ordering')
        self.prod_info = InfoBlock(self.product.model_w, self.product.columns, columns=columns)
        info_box.addWidget(self.prod_info)
        self.prod_info.setVisible(False)

        self.operation = Item('operation_to_ordering')
        self.op_info = InfoBlock(self.operation.model_w, self.operation.columns, columns=columns)
        info_box.addWidget(self.op_info)
        self.op_info.setVisible(False)

        self.matherial = Item('matherial_to_ordering')
        self.mat_info = InfoBlock(self.matherial.model_w, self.matherial.columns, columns=columns)
        info_box.addWidget(self.mat_info)
        self.mat_info.setVisible(False)

    def reload(self, value, value_type):
        if value_type == 'matherial_to_ordering':
            self.op_info.setVisible(False)
            self.prod_info.setVisible(False)
            self.mat_info.setVisible(True)
            self.mat_info.reload(value)
        elif value_type == 'operation_to_ordering':
            self.mat_info.setVisible(False)
            self.prod_info.setVisible(False)
            self.op_info.setVisible(True)
            self.op_info.reload(value)
        elif value_type == 'product_to_ordering':
            self.op_info.setVisible(False)
            self.mat_info.setVisible(False)
            self.prod_info.setVisible(True)
            self.prod_info.reload(value)


class MatherialToOrderingForm(CustomForm):
    def __init__(self, fields: list = [], value: dict = {}):
        self.item = Item('matherial_to_ordering')
        fields = fields if fields else self.item.columns
        super().__init__(self.item.model, fields, value)
        self.widgets['color_id'].setVisible(False)
        self.widgets['width'].setVisible(False)
        self.widgets['length'].setVisible(False)
        self.widgets['pieces'].setVisible(False)
        self.labels['color_id'].setVisible(False)
        self.labels['width'].setVisible(False)
        self.labels['length'].setVisible(False)
        self.labels['pieces'].setVisible(False)
        if 'product_to_ordering_id' in self.widgets:
            self.widgets['product_to_ordering_id'].setDisabled(True)
        
        self.widgets['matherial_id'].valChanged.connect(self.matherial_selected)
        self.widgets['price'].valChanged.connect(self.price_changed)
        self.widgets['number'].valChanged.connect(self.price_changed)
        self.widgets['persent'].valChanged.connect(self.price_changed)
        self.widgets['profit'].valChanged.connect(self.profit_changed)
        self.widgets['width'].valChanged.connect(self.size_changed)
        self.widgets['length'].valChanged.connect(self.size_changed)
        self.widgets['pieces'].valChanged.connect(self.size_changed)
        self.matherial_selected()

    def size_changed(self):
        matherial_value = self.widgets['matherial_id'].full_value()
        length = self.widgets['length'].value()
        pieces = self.widgets['pieces'].value()
        if matherial_value['measure'] == 'мп.':
            self.widgets['number'].setValue(length*pieces/1000)
            return
        width = self.widgets['width'].value()
        self.widgets['number'].setValue(length*width*pieces/1000000)

    def price_changed(self):
        number = self.widgets['number'].value()
        if not number:
            return
        cost = self.widgets['price'].value() * number
        profit = round(cost * self.widgets['persent'].value() / 100 / number, 2)
        self.widgets['profit'].set_value(profit)
        self.widgets['cost'].set_value(cost + profit*number)

    def profit_changed(self):
        number = self.widgets['number'].value()
        if not number:
            return
        cost = self.widgets['price'].value() * number
        profit = self.widgets['profit'].value() * number
        persent = profit * 100 / cost
        self.widgets['persent'].set_value(persent)
        self.widgets['cost'].set_value(cost + profit)

    def matherial_selected(self):
        matherial_value = self.widgets['matherial_id'].full_value()
        if not matherial_value:
            return
        if not self.widgets['price'].value():
            self.widgets['price'].setValue(matherial_value['cost'])
        self.measure = matherial_value['measure']
        self.widgets['number'].setSuffix(' ' + self.measure)
        app = App()
        if self.measure == app.config['measure linear']:
            self.labels['length'].setVisible(True)
            self.widgets['length'].setVisible(True)
            self.labels['pieces'].setVisible(True)
            self.widgets['pieces'].setVisible(True)
        if self.measure == app.config['measure square']:
            self.labels['length'].setVisible(True)
            self.labels['width'].setVisible(True)
            self.widgets['width'].setVisible(True)
            self.widgets['length'].setVisible(True)
            self.labels['pieces'].setVisible(True)
            self.widgets['pieces'].setVisible(True)
        if matherial_value['color_group_id']:
            self.labels['color_id'].setVisible(True)
            self.widgets['color_id'].setVisible(True)
            self.widgets['color_id'].group_id = matherial_value['color_group_id']


class OperationToOrderingForm(CustomForm):
    def __init__(self, fields: list = [], value: dict = {}):
        self.item = Item('operation_to_ordering')
        self.price = 0
        self.cost = 0
        self.eq_price = 0
        fields = fields if fields else self.item.columns
        super().__init__(self.item.model, fields, value)
        operation_value = self.widgets['operation_id'].full_value()
        if operation_value:
            self.measure = operation_value['measure']
            self.widgets['number'].setSuffix(' ' + self.measure)
            self.cost = operation_value['cost']
            self.price = operation_value['price']
            self.eq_price = operation_value['equipment_price']

        self.widgets['operation_id'].valChanged.connect(self.operation_selected)
        self.widgets['number'].valChanged.connect(self.price_changed)
        self.widgets['price'].valChanged.connect(self.price_changed)
        self.widgets['cost'].valChanged.connect(self.price_changed)
        if 'product_to_ordering_id' in self.widgets:
            self.widgets['product_to_ordering_id'].setDisabled(True)
        
        self.operation_selected()

    def set_widget_value(self, widget_name, value):
        if widget_name in self.widgets:
            self.widgets[widget_name].set_value(value)

    def operation_selected(self):
        operation_value = self.widgets['operation_id'].full_value()
        if not operation_value:
            return
        self.cost = operation_value['cost']
        self.price = operation_value['price']
        self.eq_price = operation_value['equipment_price']
        self.set_widget_value('price', operation_value['price'])
        self.set_widget_value('user_sum', operation_value['price'])
        self.set_widget_value('user_id', operation_value['user_id'])
        self.set_widget_value('equipment_id', operation_value['equipment_id'])
        self.set_widget_value('equipment_cost', operation_value['equipment_price'])
        self.widgets['cost'].setValue(operation_value['cost'])
        self.measure = operation_value['measure']
        self.widgets['number'].setSuffix(' ' + self.measure)

    def price_changed(self):
        number = self.widgets['number'].value()
        if not number:
            return
        cost = self.cost * number
        price = self.price * number
        equipment_cost = self.eq_price * number
        self.widgets['cost'].set_value(cost)
        self.set_widget_value('equipment_cost', equipment_cost)
        self.set_widget_value('user_sum', price)

    
class ProductToOrderingForm(CustomForm):
    sizeChanged = pyqtSignal(float, float, int, float)
    def __init__(self, fields: list = [], value: dict = {}):
        self.item = Item('product_to_ordering')
        fields = fields if fields else self.item.columns
        super().__init__(self.item.model, fields, value)
        self.numbers_to_product = Item('numbers_to_product')
        err = self.numbers_to_product.get_filter('product_id', value['product_id'])
        if err:
            error(err)
        else:
            self.numbers_to_product.values.sort(key=lambda v: v['number'], reverse=True)
            
        self.widgets['width'].setVisible(False) 
        self.widgets['length'].setVisible(False)
        self.labels['width'].setVisible(False)
        self.labels['length'].setVisible(False)
        self.labels['pieces'].setVisible(False)
        self.widgets['pieces'].setVisible(False)
        if 'product_to_ordering_id' in self.widgets:
            self.widgets['product_to_ordering_id'].setDisabled(True)

        self.widgets['product_id'].valChanged.connect(self.product_selected)
        self.widgets['price'].valChanged.connect(self.price_changed)
        self.widgets['number'].valChanged.connect(self.price_changed)
        self.widgets['persent'].valChanged.connect(self.persent_changed)
        self.widgets['profit'].valChanged.connect(self.profit_changed)
        self.widgets['width'].valChanged.connect(self.size_changed)
        self.widgets['length'].valChanged.connect(self.size_changed)
        self.widgets['pieces'].valChanged.connect(self.size_changed)
        if value:
            self.product_selected(is_new=False)

    def size_changed(self):
        self.product_value = self.widgets['product_id'].full_value()
        self.numbers_to_product.values.sort(key=lambda v: v['size'], reverse=True)
        length = self.widgets['length'].value()
        pieces = self.widgets['pieces'].value()
        length_m = length/1000
        app = App()
        if self.product_value['measure_id'] == app.config['measure linear']:
            self.widgets['number'].setValue(length_m*pieces)
        elif self.product_value['measure_id'] == app.config['measure square']:
            width = self.widgets['width'].value()
            width_m = width/1000
            square = length_m * width_m
            self.widgets['number'].setValue(square * pieces)
                   
        self.sizeChanged.emit(
            self.widgets['width'].value(),
            self.widgets['length'].value(),
            self.widgets['pieces'].value(),
            self.widgets['number'].value(),
        )

    def price_changed(self):
        number = self.widgets['number'].value()
        if not number:
            return
                
        cost = self.widgets['price'].value() * number
        profit = round(cost * self.widgets['persent'].value() / 100 / number, 1)
        self.widgets['profit'].set_value(profit)
        cost_value = round(cost + profit*number, 1)
        self.widgets['cost'].set_value(cost_value)
        
        self.sizeChanged.emit(
            self.widgets['width'].value(),
            self.widgets['length'].value(),
            self.widgets['pieces'].value(),
            self.widgets['number'].value(),
        )
        
    def persent_changed(self):
        number = self.widgets['number'].value()
        if not number:
            return
        cost = self.widgets['price'].value() * number
        profit = round(cost * self.widgets['persent'].value() / 100 / number, 2)
        self.widgets['profit'].set_value(profit)
        self.widgets['cost'].set_value(round(cost + profit*number, 1))
        
    def profit_changed(self):
        number = self.widgets['number'].value()
        if not number:
            return
        cost = self.widgets['price'].value() * number
        profit = self.widgets['profit'].value() * number
        persent = profit * 100 / cost
        self.widgets['persent'].set_value(persent)
        self.widgets['cost'].set_value(round(cost + profit, 1))

    def product_selected(self, is_new=True):
        self.product_value = self.widgets['product_id'].full_value()
        if not self.product_value:
            return
        self.measure = self.product_value['measure']
        self.measure_id = self.product_value['measure_id']
        self.widgets['number'].setSuffix(' ' + self.measure)
        if 'name' in self.widgets:
            self.widgets['name'].setValue(self.product_value['name'])
        if is_new:
            self.widgets['product_id'].setDisabled(True)
            self.widgets['price'].setValue(self.product_value['cost'])
            self.widgets['name'].setValue(f"{self.product_value['short_name']} до зам.{self.value['ordering_id']}") 
        app = App()
        if self.measure_id == app.config['measure linear']:
            self.labels['length'].setVisible(True)
            self.widgets['length'].setVisible(True) 
            self.labels['pieces'].setVisible(True)
            self.widgets['pieces'].setVisible(True)
        if self.measure_id == app.config['measure square']:
            self.labels['length'].setVisible(True)
            self.labels['width'].setVisible(True)
            self.widgets['width'].setVisible(True)
            self.widgets['length'].setVisible(True)
            self.labels['pieces'].setVisible(True)
            self.widgets['pieces'].setVisible(True)


class ComboBox(QComboBox):
    selectionChanged = pyqtSignal(int, int)
    def __init__(self):
        super().__init__()
        self.prewious = -1
        self.currentIndexChanged.connect(self.changed)

    def setCurrentIndex(self, index: int) -> None:
        self.prewious = self.currentIndex()
        return super().setCurrentIndex(index)

    def changed(self, index):
        prew = self.prewious
        self.prewious = index
        self.selectionChanged.emit(index, prew)


class ComplexList(QListWidget):
    def __init__(self):
        super().__init__()
        self.setSizePolicy(QSizePolicy(QSizePolicy.Policy.Minimum, QSizePolicy.Policy.Minimum))
        self.setMaximumHeight(30)

    def setDataset(self, dataset):
        self.clear()
        self.setMaximumHeight(25*len(dataset))
        for v in dataset:
            listItem = QListWidgetItem(v[1], self)
            listItem.setData(FULL_VALUE_ROLE, v[0])
            listItem.setCheckState(Qt.CheckState.Checked if v[2] else Qt.CheckState.Unchecked)

    def get_checked(self):
        res = []
        for i in range(self.count()):
            li = self.item(i)
            if li.checkState() == Qt.CheckState.Checked:
                res.append(li.data(FULL_VALUE_ROLE))
        return res


class OrderingForm(CustomForm):
    def __init__(self, fields: list = [], value: dict = {}):
        self.item = Item('ordering')
        fields = fields if fields else self.item.columns
        super().__init__(self.item.model, fields, value, 2)
        self.widgets['price'].valChanged.connect(self.price_changed)
        self.widgets['persent'].valChanged.connect(self.price_changed)
        self.widgets['profit'].valChanged.connect(self.profit_changed)

    def price_changed(self):
        cost = self.widgets['price'].value()
        profit = round(cost * self.widgets['persent'].value() / 100, 2)
        self.widgets['profit'].set_value(profit)
        self.widgets['cost'].set_value(cost + profit)

    def profit_changed(self):
        cost = self.widgets['price'].value()
        if cost == 0:
            return
        profit = self.widgets['profit'].value()
        persent = profit * 100 / cost
        self.widgets['persent'].set_value(persent)
        self.widgets['cost'].set_value(cost + profit)