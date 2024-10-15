import math
from PyQt6.QtCore import (
    Qt,
    pyqtSignal,
    )
from PyQt6.QtGui import QIcon
from PyQt6.QtWidgets import (
    QLabel,
    QPushButton,
    QWidget,
    QVBoxLayout,
    QHBoxLayout,
    QSplitter,
    QTabWidget,
    QScrollArea,
    QTreeWidgetItem,
    )

from data.app import App
from data.model import Item
from widgets.Dialogs import error
from common.params import FULL_VALUE_ROLE
from widgets.Form import CustomFormDialog, Selector
from widgets.ButtonsBlock import ButtonsBlock
from widgets.Tree import Tree
from data_widgets.Helpers import (
    OrdTree, 
    PMOInfo, 
    MatherialToOrderingForm, 
    OperationToOrderingForm, 
    ProductToOrderingForm,
    ComboBox,
    ComplexList,
    OrderingForm,
    )
from common.funcs import fake_id


class MatherialExtra:
    def __init__(self) -> None:
        self.name = 'matherial'
        self.value = {
            'matherial': {},
            'matherial_to_ordering': {},
        }

    def set_value(self, matherial:int|dict, ordering_id=0):
        if type(matherial) == int:                              # by id
            mat = Item('matherial')
            err = mat.get_w(matherial)
            if err:
                error(err)
                return False
            self.value['matherial'] = mat.value
        else:                                                   # by value
            self.value['matherial'] = matherial
        
        self.value['matherial_to_ordering'] = self.make_m2o(ordering_id).value
                
    def make_m2o(self, ordering_id):
        m2o = Item('matherial_to_ordering')
        m2o.create_default_w()
        m2o.value["id"] = next(fake_id)
        m2o.value["ordering_id"] = ordering_id
        m2o.value["matherial_id"] = self.value['matherial']['id']
        m2o.value["matherial"] = self.value['matherial']['name']
        m2o.value["price"] = self.value['matherial']['cost']
        m2o.value["cost"] = self.value['matherial']['cost']
        m2o.value["product_to_ordering_id"] = 0
        return m2o

    def get_cost(self):
        return self.value['matherial_to_ordering']['cost']

    def save_2o(self, ordering_id, p2o_id=0):
        self.value["matherial_to_ordering"]["ordering_id"] = ordering_id
        self.value["matherial_to_ordering"]["id"] = 0
        self.value["matherial_to_ordering"]["product_to_ordering_id"] = p2o_id
        m2o = Item('matherial_to_ordering')
        m2o.value = self.value["matherial_to_ordering"]
        err = m2o.save()
        if err:
            error(err)
            return
        self.value["matherial_to_ordering"] = m2o.value
        return m2o.value


class OperationExtra:
    def __init__(self) -> None:
        self.name = 'operation'
        self.value = {
            'operation': {},
            'operation_to_ordering': {},
        }

    def set_value(self, operation:int|dict, ordering_id:int=0):
        if type(operation) == int:                              # by id
            op = Item('operation')
            err = op.get_w(operation)
            if err:
                error(err)
                return False
            self.value['operation'] = op.value
        else:                                                   # by value
            self.value['operation'] = operation
        
        self.value['operation_to_ordering'] = self.make_o2o(ordering_id).value
                
    def make_o2o(self, ordering_id:int):
        o2o = Item('operation_to_ordering')
        o2o.create_default_w()
        o2o.value["id"] = next(fake_id)
        o2o.value["ordering_id"] = ordering_id
        o2o.value["operation_id"] = self.value['operation']['id']
        o2o.value["operation"] = self.value['operation']['name']
        o2o.value["price"] = self.value['operation']['cost']
        o2o.value["user_id"] = self.value['operation']['user_id']
        o2o.value["user"] = self.value['operation']['user']
        o2o.value["user_sum"] = self.value['operation']['price']
        o2o.value["cost"] = self.value['operation']['cost']
        o2o.value["equipment_id"] = self.value['operation']['equipment_id']
        o2o.value["equipment"] = self.value['operation']['equipment']
        o2o.value["equipment_cost"] = self.value['operation']['equipment_price']
        o2o.value["product_to_ordering_id"] = 0
        return o2o
    
    def get_cost(self):
        return self.value['operation_to_ordering']['cost']    
    
    def save_2o(self, ordering_id, p2o_id=0):
        self.value["operation_to_ordering"]["ordering_id"] = ordering_id
        self.value["operation_to_ordering"]["id"] = 0
        self.value["operation_to_ordering"]["product_to_ordering_id"] = p2o_id
        o2o = Item('operation_to_ordering')
        o2o.value = self.value["operation_to_ordering"]
        err = o2o.save()
        if err:
            error(err)
            return
        self.value["operation_to_ordering"] = o2o.value
        return o2o.value


class ProductExtra:
    def __init__(self) -> None:
        self.name = 'product'
        self.value = {}
        self.ordering_id = 0
        
    def set_value(self, product:dict|int, ordering_id):
        self.ordering_id = ordering_id
        if type(product) == dict:
            self.value = product
            return True
        prod = Item('product')
        err = prod.get_product_deep(product)
        if err:
            error(err)
            return False
        self.value = prod.value
        self.make_to_ordering_items()
        return True
    
    def get_cost(self):
        return self.value['product_extra']['product_to_ordering']['cost']
    
    def make_to_ordering_items(self, value:dict=None):
        if value is None:
            value = self.value
        p2o = self.make_p2o(value['product_extra'])
        value['product_extra']['product_to_ordering'] = p2o.value
        
        prod_name = value['product_extra']['product']['name']
        p2o_id = p2o.value['id']

        for val in value['matherial_extra'].values():
            for v in val:
                v['matherial_to_ordering'] = self.make_m2o(v, prod_name, p2o.value).value

        for val in value['operation_extra'].values():
            for v in val:
                v['operation_to_ordering'] = self.make_o2o(v, prod_name, p2o_id).value

        for val in value['product_deep'].values():
            for v in val:
                v['product_extra']['product_to_ordering'] = self.make_p2o(v['product_extra'], p2o_id).value
                self.make_to_ordering_items(v)
    
    def make_p2o(self, v, p2o_id=0):
        p2o = Item('product_to_ordering')
        p2o.create_default_w()
        p2o.value["id"] = next(fake_id)
        p2o.value["ordering_id"] = self.ordering_id
        p2o.value["product_id"] = v['product']['id']
        p2o.value["product"] = v['product']['name']
        p2o.value["user_id"] = v['product']["user_id"]
        p2o.value["price"] = v['product']['cost']
        p2o.value["cost"] = v['product']['cost']
        p2o.value["product_to_ordering_id"] = p2o_id
        return p2o
        
    def make_m2o(self, v, product_name='', p2o_value=None):
        m2o = Item('matherial_to_ordering')
        m2o.create_default_w()
        m2o.value["id"] = next(fake_id)
        m2o.value["ordering_id"] = 1
        m2o.value["matherial_id"] = v['matherial']['id']
        m2o.value["matherial"] = v['matherial']['name']
        m2o.value["user_id"] = p2o_value["user_id"]
        p2o_num = p2o_value['number']
        m2o.value["number"] = v['matherial_to_product']['number'] * p2o_num
        m2o.value["price"] = v['matherial_to_product']['cost']
        m2o.value["cost"] = m2o.value["price"] * p2o_num
        m2o.value["comm"] = f'auto to product {product_name}'
        m2o.value["product_to_ordering_id"] = p2o_value['id']
        return m2o
    
    def make_o2o(self, v, product_name='', p2o_id=0):
        o2o = Item('operation_to_ordering')
        o2o.create_default_w()
        o2o.value["id"] = next(fake_id)
        o2o.value["ordering_id"] = 1
        o2o.value["operation_id"] = v['operation']['id']
        o2o.value["operation"] = v['operation']['name']
        o2o.value["number"] = v['operation_to_product']['number']
        o2o.value["price"] = v['operation_to_product']['cost']
        o2o.value["user_id"] = v['operation_to_product']['user_id']
        o2o.value["user_sum"] = v['operation']['price'] * v['operation_to_product']['number']
        o2o.value["cost"] = v['operation_to_product']['cost']
        o2o.value["equipment_id"] = v['operation_to_product']['equipment_id']
        o2o.value["equipment_cost"] = v['operation_to_product']['equipment_cost']
        o2o.value["comm"] = f'auto to product {product_name}'
        o2o.value["product_to_ordering_id"] = p2o_id
        return o2o

    def get_by_uid(self, uid, value=None):
        if value is None:
            value = self.value
        if value['uid'] == uid:
            return value
        v = value['product_extra']
        if v['uid'] == uid:
            return v
        for val in value['matherial_extra'].values():
            for v in val:
                if v['uid'] == uid:
                    return v
        for val in value['operation_extra'].values():
            for v in val:
                if v['uid'] == uid:
                    return v
        for val in value['product_deep'].values():
            for v in val:
                res = self.get_by_uid(uid, v)
                if res is not None:
                    return res
                
    def replace(self, uid, new_value, product_value=None):
        if self.value['uid'] == uid:
            self.value = new_value
            return True
        
        if product_value is None:
            product_value = self.value
        
        if product_value['product_extra']['uid'] == uid:
            product_value['product_extra']['product_to_ordering'] = new_value
            return True
        
        is_find = False
        for k, val in product_value['matherial_extra'].items():
            for i, v in enumerate(val):
                if v['uid'] == uid:
                    is_find = True
                    break
        if is_find:
            product_value['matherial_extra'][k][i]['matherial_to_ordering'] = new_value
            return True
        
        for k, val in product_value['operation_extra'].items():
            for i, v in enumerate(val):
                if v['uid'] == uid:
                    is_find = True
                    break
        if is_find:
            product_value['operation_extra'][k][i]['operation_to_ordering'] = new_value
            return True
        
        for k, val in product_value['product_deep'].items():
            for i, v in enumerate(val):
                if v['uid'] == uid:
                    is_find = True
                    break

        if is_find:
            product_value['product_deep'][k][i]['product_extra']['product_to_ordering'] = new_value
            return True
        
        for val in product_value['product_deep'].values():
            for v in val:
                res = self.replace(uid, new_value, v) 
                if res:
                    return True
        return False
                
    def recalc_num(self, value=None):
        app = App()
        if value is None:
            value = self.value
        width = value['product_extra']['product_to_ordering']['width']
        length = value['product_extra']['product_to_ordering']['length']
        pieces = value['product_extra']['product_to_ordering']['pieces']
        number = value['product_extra']['product_to_ordering']['number']
        for k, val in value['matherial_extra'].items():
            for v in val:
                if k == 'default' or v['matherial_to_product']['is_used']:
                    if v['matherial']['measure_id'] == app.config["measure linear"] and width:
                        mat_number = (width + length) * 2 * pieces / 1000
                    else:
                        mat_number = number
                    if v['matherial']['measure_id'] == app.config["measure square"] and width:
                        k = math.sqrt(v['matherial_to_product']['number'])
                        v['matherial_to_ordering']['width'] = int(width * k)
                        v['matherial_to_ordering']['length'] = int(length * k)
                        v['matherial_to_ordering']['pieces'] = pieces
                    v['matherial_to_ordering']['number'] = mat_number * v['matherial_to_product']['number']
                    v['matherial_to_ordering']['cost'] = mat_number * v['matherial_to_ordering']['price']

        for k, val in value['operation_extra'].items():
            for v in val:
                if k == 'default' or v['operation_to_product']['is_used']:
                    if v['operation']['measure_id'] == app.config["measure linear"] and width:
                        op_number = (width + length) * 2 * pieces / 1000
                    else:
                        op_number = number
                    v['operation_to_ordering']['number'] = op_number * v['operation_to_product']['number']
                    v['operation_to_ordering']['cost'] = op_number * v['operation_to_ordering']['price']
                    v['operation_to_ordering']["user_sum"] = op_number * v['operation']['price'] * v['operation_to_product']['number']
                    v['operation_to_ordering']["equipment_cost"] = op_number * v['operation_to_product']['equipment_cost']

        for k, val in value['product_deep'].items():
            for v in val:
                if k == 'default' or v['product_extra']['product_to_product']['is_used']:
                    if v['product_extra']['product']['measure_id'] == app.config["measure linear"] and width:
                        prod_number = (width + length) * 2 * pieces / 1000
                    else:
                        prod_number = number
                    if v['product_extra']['product']['measure_id'] == app.config["measure square"] and width:
                        k = math.sqrt(v['product_extra']['product_to_product']['number'])
                        v['product_extra']['product_to_ordering']['width'] = int(width * k)
                        v['product_extra']['product_to_ordering']['length'] = int(length * k)
                        v['product_extra']['product_to_ordering']['pieces'] = pieces
                    v['product_extra']['product_to_ordering']['number'] = prod_number * v['product_extra']['product_to_product']['number']
                    v['product_extra']['product_to_ordering']['cost'] = prod_number * v['product_extra']['product_to_ordering']['price']
                    self.recalc_num(v)

    def get_persent_by_number(self, value=None):
        if value is None:
            value = self.value
        numbers_to_product = Item('numbers_to_product')
        err = numbers_to_product.get_filter('product_id', value['product_extra']['product']['id'])
        if err:
            error(err)
            return 0
        if not numbers_to_product.values:
            return 0
        
        app = App()
        is_linear = value['product_extra']['product']['measure_id'] == app.config['measure linear']
        is_square = value['product_extra']['product']['measure_id'] == app.config['measure square']
        if is_linear or is_square:
            size = 0.0
            persent = 0
            numbers_to_product.values.sort(key=lambda v: v['size'], reverse=True)
            size_m = value['product_extra']['product_to_ordering']['length']/1000
            if is_square:
                width_m = value['product_extra']['product_to_ordering']['width']/1000
                size_m *= width_m
        
            for v in numbers_to_product.values:
                if v['size'] and size_m >= v['size']:
                    persent = v['persent']
                    size = v['size']
                    break
            pieces_values = [v for v in numbers_to_product.values if v['size'] == size]
            pieces_values.sort(key=lambda v: v['pieces'], reverse=True)
            pieces = value['product_extra']['product_to_ordering']['pieces']
            for v in pieces_values:
                if v['pieces'] and pieces >= v['pieces']:
                    return v['persent']
            return persent        
        
        numbers_to_product.values.sort(key=lambda v: v['number'], reverse=True)
        number = value['product_extra']['product_to_ordering']['number']
        for v in numbers_to_product.values:
            if v['number'] and number >= v['number']:
                return v['persent']
        return 0

    def recalc(self, value=None):
        is_top = False
        if value is None:
            value = self.value
            is_top = True
        number = value['product_extra']['product_to_ordering']['number']
        if not number:
            return (None,)
        total = 0
        total_oper = 0
        matherials_price = 0
        operations_price = 0
        amortisation = 0
        for k, val in value['matherial_extra'].items():
            for v in val:
                if k == 'default' or v['matherial_to_product']['is_used']:
                    total += v['matherial_to_ordering']['price'] 
                    matherials_price += v['matherial']['price'] * number * v['matherial_to_product']['number']
        for k, val in value['operation_extra'].items():
            for v in val:
                if k == 'default' or v['operation_to_product']['is_used']:
                    total_oper = v['operation_to_ordering']['cost']
                    total += total_oper / number
                    operations_price += v['operation']['price'] * number * v['operation_to_product']['number']
                    amortisation += v['operation']['equipment_price'] * number * v['operation_to_product']['number']
        for k, val in value['product_deep'].items():
            for v in val:
                if k == 'default' or v['product_extra']['product_to_product']['is_used']:
                    total_sum, mat_price, op_price, amort = self.recalc(v)
                    persent = self.get_persent_by_number(v)
                    total += (total_sum + total_sum * persent / 100) / number
                    matherials_price += mat_price
                    operations_price += op_price
                    amortisation += amort
        value['product_extra']['product_to_ordering']['price'] = total
        
        if is_top:
            if value['product_extra']['product_to_ordering']['length']:
                total_price = total * number / value['product_extra']['product_to_ordering']['pieces']
                if total_price < value['product_extra']['product']['min_cost']:
                    total_price = value['product_extra']['product']['min_cost']
                    cost = total_price * value['product_extra']['product_to_ordering']['pieces']
                    value['product_extra']['product_to_ordering']['cost'] = cost
                    return cost, matherials_price, operations_price, amortisation
        persent = self.get_persent_by_number()
        price = round(total + total * persent / 100, 2)
        value['product_extra']['product_to_ordering']['price'] = price
        cost = price * number
        value['product_extra']['product_to_ordering']['cost'] = cost
        return cost, matherials_price, operations_price, amortisation
    
    def save_2o(self, ordering_id, p2o_id=0):
        self.value['product_extra']['product_to_ordering']["ordering_id"] = ordering_id
        self.value['product_extra']['product_to_ordering']["id"] = 0
        self.value['product_extra']['product_to_ordering']["product_to_ordering_id"] = p2o_id
        p2o = Item('product_to_ordering')
        p2o.value = self.value['product_extra']['product_to_ordering']
        err = p2o.create_p2o_defaults()
        if err:
            error(err)
            return
        self.value['product_extra']['product_to_ordering'] = p2o.value
        self.create_used(ordering_id)
        return p2o.value
                    
    def create_used(self, ordering_id, value=None):
        if value is None:
            value = self.value
        p2o_id = value['product_extra']['product_to_ordering']['id']
        for val in value['matherial_extra'].values():
            for v in val:
                if v['matherial_to_product']['is_used']:
                    v['matherial_to_ordering']['ordering_id'] = ordering_id
                    v['matherial_to_ordering']['product_to_ordering_id'] = p2o_id
                    m2o = Item('matherial_to_ordering')
                    m2o.value = v['matherial_to_ordering']
                    m2o.value['id'] = 0
                    err = m2o.save()
                    if err:
                        error(err)
        for val in value['operation_extra'].values():
            for v in val:
                if v['operation_to_product']['is_used']:
                    v['operation_to_ordering']['ordering_id'] = ordering_id
                    v['operation_to_ordering']['product_to_ordering_id'] = p2o_id
                    o2o = Item('operation_to_ordering')
                    o2o.value = v['operation_to_ordering']
                    o2o.value['id'] = 0
                    err = o2o.save()
                    if err:
                        error(err)
        for k, val in value['product_deep'].items():
            for v in val:
                if k == 'default' or v['product_extra']['product_to_product']['is_used']:
                    v['product_extra']['product_to_ordering']['ordering_id'] = ordering_id
                    v['product_extra']['product_to_ordering']['product_to_ordering_id'] = p2o_id
                    v['product_extra']['product_to_ordering']["name"] = v['product_extra']['product']['name']
                    p2o = Item('product_to_ordering')
                    p2o.value = v['product_extra']['product_to_ordering']
                    p2o.value['id'] = 0
                    err = p2o.create_p2o_defaults()
                    if err:
                        error(err)
                        return
                    v['product_extra']['product_to_ordering'] = p2o.value
                    self.create_used(ordering_id, v)        


class ProductView(QWidget):
    productChanged = pyqtSignal()
    def __init__(self, product: ProductExtra) -> None:
        super().__init__()
        self.box = QVBoxLayout()
        self.setLayout(self.box)
        self.product = product
        self.selects = {}

    def clear_layout(self):
        while self.box.count():
            child = self.box.takeAt(0)
            if child.widget():
                child.widget().deleteLater()

    def reload(self):
        self.clear_layout()
        for list_name, list_items in self.product.value['matherial_extra'].items():
            
            if list_name == 'default':
                for v in list_items:
                    if v['matherial']['color_group_id']:
                        w = self.create_color_selector(v)
                        self.box.addWidget(w)
                continue
            # matherial multiselect
            if list_items[0]['matherial_to_product']['is_multiselect']:
                self.selects[list_name] = ComplexList()
                dataset = []
                color_selectors = []
                for v in list_items:
                    dataset.append([v['uid'], v['matherial']['name'], v['matherial_to_product']['is_used']])
                    if v['matherial_to_product']['is_used'] and v['matherial']['color_group_id']:
                        color_selectors.append(self.create_color_selector(v))
                
                self.selects[list_name].setDataset(dataset)
                self.box.addWidget(QLabel(list_name))
                self.box.addWidget(self.selects[list_name])
                self.selects[list_name].itemChanged.connect(
                    lambda item, name='matherial': self.multilist_changed(item, name)
                    )
                
                for cs in color_selectors:
                    self.box.addWidget(cs)
    
            # matherial select
            else:
                self.selects[list_name] = ComboBox()
                used_index = 0
                for i, v in enumerate(list_items):
                    self.selects[list_name].addItem(v['matherial']['name'], userData=v['uid'])
                    if v['matherial_to_product']['is_used'] and not used_index:
                        used_index = i
                        self.selects[list_name].setCurrentIndex(used_index)
    
                if used_index == 0:
                    list_items[0]['matherial_to_product']['is_used'] = True
    
                self.box.addWidget(QLabel(list_name))
                self.box.addWidget(self.selects[list_name])
                self.selects[list_name].selectionChanged.connect(
                    lambda cur, prew, cb=self.selects[list_name], name='matherial': self.list_selected(
                        cur, prew, cb, name
                        )
                )

                # color selector
                if list_items[used_index]['matherial']['color_group_id']:
                    w = self.create_color_selector(list_items[used_index])
                    self.box.addWidget(w)
    
            
        for list_name, list_items in self.product.value['operation_extra'].items():
            if list_name == 'default':
                continue
            if list_items[0]['operation_to_product']['is_multiselect']:
                self.selects[list_name] = ComplexList()
                dataset = []
                for v in list_items:
                    dataset.append([v['uid'], v['operation']['name'], v['operation_to_product']['is_used']])
                self.selects[list_name].setDataset(dataset)
                self.selects[list_name].itemChanged.connect(
                    lambda item, name='operation': self.multilist_changed(item, name)
                )
            else:
                self.selects[list_name] = ComboBox()
                used_index = 0
                for i, v in enumerate(list_items):
                    self.selects[list_name].addItem(v['operation']['name'], userData=v['uid'])
                    if v['operation_to_product']['is_used']:
                        self.selects[list_name].setCurrentIndex(i)
                        used_index = i
                if used_index == 0:
                    list_items[0]['operation_to_product']['is_used'] = True
                self.selects[list_name].selectionChanged.connect(
                    lambda cur, prew, cb=self.selects[list_name], name='operation': self.list_selected(
                        cur, prew, cb, name
                        )
                )

            self.box.addWidget(QLabel(list_name))
            self.box.addWidget(self.selects[list_name])

        # products lists
        for list_name, list_items in self.product.value['product_deep'].items():
            if list_name == 'default':
                for v in list_items:
                    pe = ProductExtra()
                    pe.set_value(v, self.product.ordering_id)
                    w = ProductView(pe)
                    n = w.reload()
                    if n:
                        self.box.addWidget(w) 
                continue
            if list_items[0]['product_extra']['product_to_product']['is_multiselect']:
                self.selects[list_name] = ComplexList()
                dataset = []
                prod_views = []
                for v in list_items:
                    dataset.append([
                        v['product_extra']['uid'], 
                        v['product_extra']['product']['name'], 
                        v['product_extra']['product_to_product']['is_used'],
                        ])
                    if v['product_extra']['product_to_product']['is_used']: 
                        pe = ProductExtra()
                        pe.set_value(v, self.product.ordering_id)
                        w = ProductView(pe)
                        n = w.reload()
                        if n:
                            prod_views.append(w) 

                self.selects[list_name].setDataset(dataset)
                self.box.addWidget(QLabel(list_name))
                self.box.addWidget(self.selects[list_name])
                self.selects[list_name].itemChanged.connect(
                    lambda item, name='product': self.multilist_changed(item, name)
                )
                for pw in prod_views:
                    self.box.addWidget(pw)
                    pw.productChanged.connect(self.subproduct_changed)

            else:
                self.selects[list_name] = ComboBox()
                used_index = 0
                for i, v in enumerate(list_items):
                    self.selects[list_name].addItem(v['product_extra']['product']['name'], userData=v['product_extra']['uid'])
                    if v['product_extra']['product_to_product']['is_used']:
                        self.selects[list_name].setCurrentIndex(i)
                        used_index = i
                if used_index == 0:
                    list_items[0]['product_extra']['product_to_product']['is_used'] = True
                self.box.addWidget(QLabel(list_name))
                self.box.addWidget(self.selects[list_name])
                self.selects[list_name].selectionChanged.connect(
                    lambda cur, prew, cb=self.selects[list_name], name='product': self.list_selected(
                        cur, prew, cb, name
                        )
                )
                pe = ProductExtra()
                pe.set_value(list_items[used_index], self.product.ordering_id)
                w = ProductView(pe)
                n = w.reload()
                if n:
                    self.box.addWidget(w)
                    w.productChanged.connect(self.subproduct_changed)
    
            
        return self.box.count()

    def create_color_selector(self, value):
        matherial = value['matherial']
        w = Selector(
                    'color', 
                    title=f'Колір {matherial["name"]}', 
                    group_id=matherial['color_group_id'],
                    )
        color_id = value['matherial_to_ordering']['color_id']
        if color_id:
            w.setValue(color_id)
        w.valChanged.connect(lambda w=w, value=value: self.color_changed(w, value))
        return w

    def color_changed(self, w: Selector, value):
        color = w.full_value()
        value['matherial_to_ordering']['color_id'] = color['id']
        value['matherial_to_ordering']['color'] = color['name']

    def multilist_changed(self, item, name):
        v_uid = item.data(FULL_VALUE_ROLE)
        v = self.product.get_by_uid(v_uid)
        if v is None:
            return
        
        v[f'{name}_to_product']['is_used'] = (item.checkState() == Qt.CheckState.Checked) 
        self.product.recalc_num()
        self.product.recalc()
        self.reload()
        self.productChanged.emit()

    def list_selected(self, curr, prew, cb, name):
        v_uid = cb.itemData(curr)
        v_prew_uid = cb.itemData(prew)
        v = self.product.get_by_uid(v_uid)
        if v is None:
            return
        v[f'{name}_to_product']['is_used'] = True
        v_prew = self.product.get_by_uid(v_prew_uid)
        if v_prew is None:
            return
        v_prew[f'{name}_to_product']['is_used'] = False
        self.product.recalc_num()
        self.product.recalc()
        self.reload()
        self.productChanged.emit()

    def subproduct_changed(self):
        self.product.recalc_num()
        self.product.recalc()
        self.reload()
        self.productChanged.emit()


class ProductWidget(QWidget):
    def __init__(self, product_value: dict):
        super().__init__()
        form_box = QVBoxLayout()
        self.setLayout(form_box)
        self.form_plaseholder = QScrollArea()
        self.form_plaseholder.setWidgetResizable(True)
        form_box.addWidget(self.form_plaseholder)
        self.price_info_text = "Матеріали:\t{mat:.2f}\nОперації:\t{op:.2f}\nАмортизація:\t{am:.2f}\n-----------\nЗагалом:\t{tot:.2f}"
        self.price_info = QLabel(self.price_info_text.format(mat=0.00, op=0.00, am=0.00, tot=0.00))
        form_box.addWidget(self.price_info)

        self.product = ProductExtra()
        if not self.product.set_value(product_value['id'], 0):
            return
        self.product.value['product_extra']['product_to_ordering']['name'] = self.product.value['product_extra']['product']['name']
        
        self.pw = ProductView(self.product)
        self.pw.reload()
        self.pw.product.recalc_num()
        res = self.pw.product.recalc()
        self.reload_price_info(res)
        self.add_product()
        self.pw.productChanged.connect(self.product_changed)

        w = QWidget()
        box = QVBoxLayout()
        w.setLayout(box)
        box.addWidget(self.pw)
        box.addWidget(self.form)
        self.form_plaseholder.setWidget(w)

    def add_product(self):
        fields = ["name", "product_id", "width", "length", "pieces", "number", 
                "price", "persent", "profit", "cost", "info", "user_id"]
        self.form = ProductToOrderingForm(
            fields=fields, 
            value=self.product.value['product_extra']['product_to_ordering']
            )
        self.form.sizeChanged.connect(self.form_size_changed)

    def reload_price_info(self, result):
        if result[0] is None:
            return
        _, mat, op, am = result
        tot = mat + op + am
        self.price_info.setText(
            self.price_info_text.format(mat=mat, op=op, am=am, tot=tot)
            )
        
    def form_size_changed(self, width, length, pieces, number):
        p2o = self.product.value['product_extra']['product_to_ordering']
        p2o['width'] = width
        p2o['length'] = length
        p2o['pieces'] = pieces
        p2o['number'] = number
        self.product.recalc_num()
        res = self.product.recalc()
        if res[0] is None:
            return
        self.reload_price_info(res)
        if res is not None:
            total = res[0]
            self.form.widgets['price'].setValue(total/number)
            profit = round(total * self.form.widgets['persent'].value() / 100 / number, 2)
            self.form.widgets['profit'].set_value(profit)
            self.form.widgets['cost'].set_value(total + profit*number)

    def product_changed(self):
        self.pw.reload()
        persent = self.form.widgets['persent'].value()
        self.form.reload(self.product.value['product_extra']['product_to_ordering'])
        self.form.widgets['persent'].set_value(persent)
        self.form.persent_changed()


class Calculator(QSplitter):
    positionAdded = pyqtSignal()
    def __init__(self):
        super().__init__(Qt.Orientation.Horizontal)
        self.select_tabs = QTabWidget()
        self.form = None
        self.current_value = {} #
        self.current_item = {} #
        self.ordering_id = next(fake_id)
        self.pw = None
        
        self.prod_tree = OrdTree('product', 'Вироби')
        self.mat_tree = OrdTree('matherial', 'Матеріали')
        self.op_tree = OrdTree('operation', 'Операції')
        self.select_tabs.addTab(self.prod_tree, 'Вироби')
        self.select_tabs.addTab(self.mat_tree, 'Матеріали')
        self.select_tabs.addTab(self.op_tree, 'Операції')
        self.details_tree = CalcTree(fields=['name', 'cost'], headers=['Назва', 'Вартість'])
        form_frame = QWidget()
        form_box = QVBoxLayout()
        form_frame.setLayout(form_box)

        self.form_plaseholder = QScrollArea()
        self.form_plaseholder.setWidgetResizable(True)
        self.tree_info = PMOInfo()
        self.prod_tree.itemSelected.connect(lambda v: self.add_form(v, 'product'))
        self.mat_tree.itemSelected.connect(lambda v: self.add_form(v, 'matherial'))
        self.op_tree.itemSelected.connect(lambda v: self.add_form(v, 'operation'))
        self.addWidget(self.select_tabs)
        self.addWidget(form_frame)
        form_box.addWidget(self.form_plaseholder)
        self.addWidget(self.details_tree)
        self.addWidget(self.tree_info)
        self.details_tree.itemSelected.connect(self.show_info)
        self.price_info_text = "Матеріали:\t{mat:.2f}\nОперації:\t{op:.2f}\nАмортизація:\t{am:.2f}\n-----------\nЗагалом:\t{tot:.2f}"
        self.price_info = QLabel(self.price_info_text.format(mat=0.00, op=0.00, am=0.00, tot=0.00))
        form_box.addWidget(self.price_info)

    def reload(self):
        pass
    
    def delete_selected_from_ordering(self):
        self.details_tree.delete_current()
    
    def reload_selects(self):
        self.prod_tree.reload()
        self.mat_tree.reload()
        self.op_tree.reload()

    def reload_price_info(self, result):
        if result[0] is None:
            return
        _, mat, op, am = result
        tot = mat + op + am
        self.price_info.setText(
            self.price_info_text.format(mat=mat, op=op, am=am, tot=tot)
            )
    
    def add_form(self, item_value, item_type):
        if not 'cost' in item_value:
            return
        self.pw = None
        if item_type == 'product':
            product = ProductExtra()
            if not product.set_value(item_value['id'], self.ordering_id):
                return
            product.value['product_extra']['product_to_ordering']['name'] = product.value['product_extra']['product']['name']
            self.current_item = product
            self.pw = ProductView(product)
            self.pw.reload()
            self.pw.product.recalc_num()
            res = self.pw.product.recalc()
            self.reload_price_info(res)
            self.add_product()
            self.pw.productChanged.connect(self.product_changed)
            
        if item_type == 'matherial':
            matherial = MatherialExtra()
            matherial.set_value(item_value)
            self.current_item = matherial
            self.add_matherial()
            
        if item_type == 'operation':
            operation = OperationExtra()
            operation.set_value(item_value)
            self.current_item = operation
            self.add_operation()
        
        if self.form:
            w = QWidget()
            box = QVBoxLayout()
            w.setLayout(box)
            if self.pw:
                box.addWidget(self.pw)
            box.addWidget(self.form)
            self.form_plaseholder.setWidget(w)
            self.form.saveRequested.connect(
                lambda v:self.add_position(v, self.form.item.name)
            )

    def add_edit_form(self, item_value, type_name, uid=0):
        self.pw = None
       
        if type_name == 'matherial_to_ordering':
            self.add_matherial_to_ordering(item_value)

        if type_name == 'operation_to_ordering':
            self.add_operation_to_ordering(item_value)

        if type_name == 'product_to_ordering':
            self.add_product_to_ordering(item_value)

        if self.form:
            w = QWidget()
            box = QVBoxLayout()
            w.setLayout(box)
            if self.pw:
                box.addWidget(self.pw)
            box.addWidget(self.form)
            self.form_plaseholder.setWidget(w)
            self.form.saveRequested.connect(
                lambda v, uid=uid:self.edit_position(v, self.form.item.name, uid)
            )

    def form_size_changed(self, width, length, pieces, number):
        p2o = self.current_item.value['product_extra']['product_to_ordering']
        p2o['width'] = width
        p2o['length'] = length
        p2o['pieces'] = pieces
        p2o['number'] = number
        self.current_item.recalc_num()
        res = self.current_item.recalc()
        if res[0] is None:
            return
        self.reload_price_info(res)
        if res is not None:
            total = res[0]
            self.form.widgets['price'].setValue(total/number)
            profit = round(total * self.form.widgets['persent'].value() / 100 / number, 2)
            self.form.widgets['profit'].set_value(profit)
            self.form.widgets['cost'].set_value(total + profit*number)

    def product_changed(self):
        self.pw.reload()
        persent = self.form.widgets['persent'].value()
        self.form.reload(self.current_item.value['product_extra']['product_to_ordering'])
        self.form.widgets['persent'].set_value(persent)
        self.form.persent_changed()

    def add_position(self, value, name):
        self.details_tree.dataset.append(self.current_item)
        self.details_tree.reload()
        self.positionAdded.emit()

    def edit_position(self, form_value, type_name, uid):
        if not uid:
            self.current_item.value[type_name] = form_value
        else:
            res = self.current_item.replace(uid, form_value)
            if not res:
                return
            if not form_value['product_to_ordering_id']:
                self.current_item.recalc_num()
            result = self.current_item.recalc()
            self.reload_price_info(result)
        self.details_tree.reload()

    def add_matherial(self):
        fields = ["matherial_id", "width", "length", "pieces", 
                      "color_id", "number", "price", "persent", 
                      "profit", "cost", "comm", "user_id"]
        self.form = MatherialToOrderingForm(
            fields=fields, 
            value=self.current_item.value['matherial_to_ordering']
            )
        
    def add_matherial_to_ordering(self, value):
        fields = ["matherial_id", "width", "length", "pieces", 
                      "color_id", "number", "price", "persent", 
                      "profit", "cost", "comm", "user_id"]
        self.form = MatherialToOrderingForm(
            fields=fields, 
            value=value,
            )

    def add_operation(self):
        fields = ["operation_id", "number", 
                      "price", "user_sum", "cost", "equipment_id", 
                      "equipment_cost", "comm", "user_id"]
        self.form = OperationToOrderingForm(
            fields=fields, 
            value=self.current_item.value['operation_to_ordering']
            )
    
    def add_operation_to_ordering(self, value):
        fields = ["operation_id", "number", 
                      "price", "user_sum", "cost", "equipment_id", 
                      "equipment_cost", "comm", "user_id"]
        self.form = OperationToOrderingForm(
            fields=fields, 
            value=value,
            )

    def add_product(self):
        fields = ["name", "product_id", "width", "length", "pieces", "number", 
                "price", "persent", "profit", "cost", "info", "user_id"]
        self.form = ProductToOrderingForm(
            fields=fields, 
            value=self.current_item.value['product_extra']['product_to_ordering']
            )
        self.form.sizeChanged.connect(self.form_size_changed)

    def add_product_to_ordering(self, value):
        fields = ["name", "product_id", "width", "length", "pieces", "number", 
                "price", "persent", "profit", "cost", "info", "user_id"]
        self.form = ProductToOrderingForm(
            fields=fields, 
            value=value,
            )
        self.form.sizeChanged.connect(self.form_size_changed)

    def show_info(self, tree_value):
        self.current_item = tree_value['item']
        type_name = tree_value['type']
        uid = tree_value['uid']
        if not uid:
            self.tree_info.reload(self.current_item.value[type_name], type_name)
            self.add_edit_form(self.current_item.value[type_name], type_name)
        else:
            value = self.current_item.get_by_uid(uid)
            self.tree_info.reload(value[type_name], type_name)
            self.add_edit_form(value[type_name], type_name, uid)
    
    def get_total_cost(self):
        return self.details_tree.recalc_cost()


class CalcTree(Tree):
    def __init__(self, 
                 items: list[MatherialExtra | OperationExtra | ProductExtra] = None, 
                 fields: list = None, 
                 headers: list = None
                 ):
        super().__init__('product_to_ordering', '', items, fields, headers)
        
        self.dataset = items if items is not None else []

    def reload(self, items: list[MatherialExtra | OperationExtra | ProductExtra]=[]):
        self.clear()
        if items:
            self.dataset = items
        elif not self.dataset:
            return
        self.add_childs()
                
    def add_childs(self):
        parent = self.invisibleRootItem()
        for item in self.dataset:
            self.add_child(item.value, parent)
            
    def add_child(self, item_value, parent: QTreeWidgetItem=None):        
        if not 'product_extra' in item_value:
            self.add_other(item_value, parent)
            return
        
        prod_extra_value = item_value['product_extra']
        prod_tree_item = self.make_tree_item(prod_extra_value, 'product')
        parent.addChild(prod_tree_item)

        for k, val in item_value['matherial_extra'].items():
            for v in val:
                if k == 'default' or v['matherial_to_product']['is_used']:
                    tree_item = self.make_tree_item(v, 'matherial')
                    prod_tree_item.addChild(tree_item)
        
        for k, val in item_value['operation_extra'].items():
            for v in val:
                if k == 'default' or v['operation_to_product']['is_used']:
                    tree_item = self.make_tree_item(v, 'operation')
                    prod_tree_item.addChild(tree_item)
        
        for k, val in item_value['product_deep'].items():
            for v in val:
                if k == 'default' or v['product_extra']['product_to_product']['is_used']:
                    self.add_child(v, prod_tree_item)

    def add_other(self, item_value, parent: QTreeWidgetItem):
        if 'matherial' in item_value:
            type_name = 'matherial'
        elif 'operation' in item_value:
            type_name = 'operation'
        else:
            return
        tree_item = self.make_tree_item(item_value, type_name)
        parent.addChild(tree_item)

    def make_tree_item(self, value, type_name):
        tree_item = QTreeWidgetItem()
        to_name = f'{type_name}_to_ordering'
        if not self.fields:
            tree_item.setText(0, value[type_name]['name'])
        else:
            for i, f in enumerate(self.fields):
                if f == 'name':
                    tree_item.setText(i, str(value[type_name][f]))
                else:    
                    tree_item.setText(i, str(value[to_name][f]))
        uid = value['uid'] if 'uid' in value else 0
        tree_item.setData(1, FULL_VALUE_ROLE, {'type': to_name, 'uid': uid})
        return tree_item

    def cur_changed(self, current, previous):
        if not current or not len(self.dataset):
            return
        i = self.get_toplevel_index(current)
        val = current.data(1, FULL_VALUE_ROLE)
        self.itemSelected.emit({'item': self.dataset[i], 'type': val['type'], 'uid': val['uid']})

    def get_toplevel_index(self, item: QTreeWidgetItem):
        i = self.indexOfTopLevelItem(item)
        if i != -1:
            return i
        i = self.indexFromItem(item)
        return self.get_toplevel_index(self.itemFromIndex(i.parent()))
    
    def get_toplevel_item(self, item: QTreeWidgetItem):
        i = self.indexOfTopLevelItem(item)
        if i != -1:
            return item
        return self.get_toplevel_item(item.parent())
    
    def delete_current(self):
        cur = self.currentItem()
        top_item = self.get_toplevel_item(cur)
        i = self.indexOfTopLevelItem(top_item)
        root = self.invisibleRootItem()
        root.removeChild(top_item)
        self.dataset.pop(i)

    def recalc_cost(self):
        res = 0.0
        for v in self.dataset:
            res += v.get_cost()
        return res


class CalculatorTab(QWidget):
    def __init__(self):
        super().__init__()
        self.ordering = None        
        self.box = QVBoxLayout()
        self.box.setContentsMargins(0,0,0,0)
        self.setLayout(self.box)
        
        controls = QWidget()
        self.box.addWidget(controls, 0)
        self.hbox = QHBoxLayout()
        self.hbox.setContentsMargins(0,0,0,0)
        controls.setLayout(self.hbox)

        reload = QPushButton()
        reload.setIcon(QIcon(f'images/icons/reload.png'))
        reload.setToolTip('Оновити')
        self.hbox.addWidget(reload)
        reload.clicked.connect(self.reload_selects)

        self.hbox.addStretch()
        self.hbox.addWidget(QLabel('Загалом'))
        self.total_sum = QLabel()
        self.hbox.addWidget(self.total_sum)

        self.hbox.addStretch()
        delete = QPushButton()
        delete.setIcon(QIcon(f'images/icons/delete.png'))
        delete.setToolTip('Видалити')
        self.hbox.addWidget(delete)
        delete.clicked.connect(self.delete_from_order)
        # create_btn = QPushButton('Створити')
        # self.hbox.addWidget(create_btn)
        # create_btn.clicked.connect(self.make_it)
        new_btn = QPushButton('Нове замовлення')
        self.hbox.addWidget(new_btn)
        new_btn.clicked.connect(self.make_new)


        self.calc = Calculator()
        self.box.addWidget(self.calc, 10)
        self.calc.positionAdded.connect(self.reload_total_sum)

    def reload(self):
        pass

    def set_ordering(self, ordering: Item):
        self.ordering = ordering
    
    def reload_selects(self):
        self.calc.reload_selects()

    def delete_from_order(self):
        self.calc.delete_selected_from_ordering()
        self.reload_total_sum()
        
    def reload_total_sum(self):
        cost = self.calc.get_total_cost()
        self.total_sum.setText(str(round(cost, 2)))

    def make_it(self):
        if self.ordering is None:
            return
        for v in self.calc.details_tree.dataset:
            v.save_2o(self.ordering.value['id'])

    def make_new(self):
        ordering = self.create_ordering()
        if ordering is None:
            return
        for v in self.calc.details_tree.dataset:
            v.save_2o(self.ordering.value['id'])

    def create_ordering(self):
        ordering = Item('ordering')
        ordering.create_default_w()
        cost = float(self.total_sum.text())
        ordering.value['price'] = cost 
        ordering.value['cost'] = cost
        form = OrderingForm(value = ordering.value)
        dlg = CustomFormDialog('Створити замовлення', form)
        res = dlg.exec()
        if not res:
            return
        err = ordering.save()
        if err:
            error(err)
            return
        err = ordering.create_dirs(ordering.value['id'])
        if err:
            error(err)
        self.ordering = ordering
        return ordering
