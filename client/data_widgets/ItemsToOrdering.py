from datetime import datetime
import webbrowser
import os
import subprocess

from PyQt6.QtCore import (
    Qt,
    pyqtSignal,
    )
from PyQt6.QtWidgets import (
    QLabel,
    QPushButton,
    QWidget,
    QCheckBox,
    QVBoxLayout,
    QHBoxLayout,
    QSplitter,
    QTabWidget,
    QGridLayout,
    )


from data.model import Item
from data.app import App
from widgets.Dialogs import CustomDialog, error, DeleteDialog, askdlg
from common.params import TABLE_BUTTONS, VIRTUAL
from widgets.Form import (
    MainItemTable, 
    DetailsItemTable, 
    FormDialog, 
    CustomFormDialog,
    CustomForm,
    Selector,
    ContactSelectDialog,
    CheckWidget,
    InfoBlock,
    )
from widgets.ButtonsBlock import ButtonsBlock
from widgets.Table import TableWControls, Table
from widgets.Tree import Tree
from data_widgets.Documents import DocsTable
from common.excel import ExcelDoc
from data_widgets.Contragents import ContragentFilter 

from data_widgets.Helpers import (
    PMOInfo, 
    MatherialToOrderingForm, 
    OperationToOrderingForm, 
    ProductToOrderingForm,
    ComboBox,
    ComplexList,
    OrderingForm,
    )
from data_widgets.Calculation import CalculatorTab


class OrderingTable(MainItemTable):
    def __init__(self, fields: list = [], values: list = None, buttons=TABLE_BUTTONS, group_id=0):
        super().__init__('ordering', '', fields, values, buttons, group_id, releazed_buttons=True)
        self.pos_btn = QPushButton('_') #'[|]'
        self.table.hbox.addWidget(self.pos_btn)
        self.is_info_position_vertical = True
        self.pos_btn.clicked.connect(self.change_info_position)

    def change_info_position(self):
        if self.is_info_position_vertical:
            self.widget(1).deleteLater()
            self.tabs = QTabWidget()
            
            self.doc_table = self.inner.replaceWidget(1, self.tabs)
            self.info = InfoBlock(self.item.model_w, self.item.columns, columns=2)
            
            self.tabs.addTab(self.info, "Детально")
            self.tabs.addTab(self.doc_table, "Документи")
            
            self.inner.setStretchFactor(0, 10)
            self.inner.setStretchFactor(1, 1)
            
            self.is_info_position_vertical = False
            self.pos_btn.setText('|')
        else:
            self.tabs.setCurrentIndex(1)
            self.doc_table = self.tabs.widget(1)
            self.tabs.deleteLater()
            self.inner.addWidget(self.doc_table)
            self.inner.setStretchFactor(0, 10)
            self.inner.setStretchFactor(1, 1)
            self.info = InfoBlock(self.item.model_w, self.item.columns)
            self.addWidget(self.info)

            self.setStretchFactor(0, 3)
            self.setStretchFactor(1, 1)
            
            self.is_info_position_vertical = True
            self.pos_btn.setText('_')
    
    def prepare_value_to_save(self, value, prev_state):
        app=App()
        new_state = value['ordering_status_id']
        if (
            prev_state != new_state
            and (
                new_state == app.config['ordering state taken']    
                or  new_state == app.config['ordering state canceled']
            )
            ):
            value['finished_at'] = datetime.now().isoformat(timespec='seconds')
        return value
    
    def dialog(self, value, title):
        creation = not value['id']
        prev_state = value['ordering_status_id']
        i = Item('ordering')
        form = OrderingForm(value=value)
        dlg = CustomFormDialog(title, form)
        res = dlg.exec()
        if res and dlg.value:
            i.value = self.prepare_value_to_save(dlg.value, prev_state)
            err = i.save()
            if err:
                error(err)
                return
            if creation:
                err = i.create_dirs(i.value['id'])
                if err:
                    error(err)    
            self.reload()
            self.actionResolved.emit()


class DetailsMatherialToOrderingTable(DetailsItemTable):
    def __init__(self, fields: list = [], values: list = None, buttons=TABLE_BUTTONS, group_id=0):
        super().__init__('matherial_to_ordering', '', fields, values, buttons, group_id)
        
    def prepare_value_to_save(self, value):
        return value
    
    def dialog(self, value, title):
        i = Item(self.item.name)
        form = MatherialToOrderingForm(value=value)
        dlg = CustomFormDialog(title, form)
        form.widgets['number'].setFocus()
        res = dlg.exec()
        if res and dlg.value:
            i.value = self.prepare_value_to_save(dlg.value)
            err = i.save()
            if err:
                error(err)
                return
            self.reload()
            self.actionResolved.emit()
            

class DetailsOperationToOrderingTable(DetailsItemTable):
    def __init__(self, fields: list = [], values: list = None, buttons=TABLE_BUTTONS, group_id=0):
        super().__init__('operation_to_ordering', '', fields, values, buttons, group_id)
        
    def dialog(self, value, title):
        i = Item(self.item.name)
        form = OperationToOrderingForm(value=value)
        dlg = CustomFormDialog(title, form)
        res = dlg.exec()
        if res and dlg.value:
            i.value = dlg.value
            err = i.save()
            if err:
                error(err)
                return
            self.reload()
            self.actionResolved.emit()


class ProductView(QWidget):
    def __init__(self, product_value: dict):
        super().__init__()
        self.prod_val = product_value
        self.box = QVBoxLayout()
        self.setLayout(self.box)
        self.selects = {}
        self.multiselects = {}
        self.mat_lists, self.mat_mlists = self.get_lists('matherial_to_product')
        self.make_select_lists('matherial_to_product', self.mat_lists)
        self.make_multiselect_lists('matherial_to_product', self.mat_mlists)
        self.op_lists, self.op_mlists = self.get_lists('operation_to_product')
        self.make_select_lists('operation_to_product', self.op_lists)
        self.make_multiselect_lists('operation_to_product', self.op_mlists)
        self.pr_lists, self.pr_mlists = self.get_lists('product_to_product')
        self.make_select_lists('product_to_product', self.pr_lists)
        self.make_multiselect_lists('product_to_product', self.pr_mlists)
        
    def has_lists(self):
        return len(self.selects) and len(self.multiselects)
    
    def make_select_lists(self, item_name, lists):
        base_name = item_name.split('_to_')[0]
        if base_name == 'product':
            base_name = 'product2'
        for name, slist in lists.items():
            listname = '@'.join((item_name, name))
            self.selects[listname] = ComboBox()
            
            for i in range(len(slist)):
                self.selects[listname].addItem(slist[i][base_name], userData=slist[i])

            self.box.addWidget(QLabel(name))
            self.box.addWidget(self.selects[listname])
            self.selects[listname].selectionChanged.connect(
                lambda i, prew, n=listname, cb=self.selects[listname]: self.list_selected(i, prew, n, cb)
                )

    def list_selected(self, i, prew, listname, cb):
        item_name, list_name = listname.split('@')
        index = self.box.indexOf(cb)
        

    def make_multiselect_lists(self, item_name, mlists):
        base_name = item_name.split('_to_')[0]
        if base_name == 'product':
            base_name = 'product2'
        for name, mslist in mlists.items():
            listname = '@'.join((item_name, name))
            self.multiselects[listname] = ComplexList()
            dataset = []
            for v in mslist:
                dataset.append([v, v[base_name], False])
            self.multiselects[listname].setDataset(dataset)

            self.box.addWidget(QLabel(name))
            self.box.addWidget(self.multiselects[listname])
            
            self.multiselects[listname].itemChanged.connect(
                lambda item, n=listname: self.multilist_changed(item, n)
                )

    def multilist_changed(self, item, name):
        item_name, list_name = name.split('@')
        index = self.multiselects[name].indexFromItem(item).row()

    def get_lists(self, name):
        n2p = Item(name)
        err = n2p.get_filter_w('product_id', self.prod_val['id'])
        if err:
            error(err)
            return None, None
        
        lists = {}
        mlists = {}
        for v in n2p.values:
            if v['list_name'] == 'default':
                continue
            if v['is_multiselect']:
                if v['list_name'] not in mlists:
                    mlists[v['list_name']] = []    
                mlists[v['list_name']].append(v)
            else:
                if v['list_name'] not in lists:
                    lists[v['list_name']] = []    
                lists[v['list_name']].append(v)
        return lists, mlists
                    
    def get_applied(self):
        applied = {
            'matherial_to_product':[], 
            'operation_to_product':[], 
            'product_to_product':[],
            }
        for k, s in self.selects.items():
            item_name = k.split('@')[0]
            applied[item_name].append(s.currentData())
        for k, s in self.multiselects.items():
            item_name = k.split('@')[0]
            applied[item_name] += s.get_checked()
        return applied


class ColorDialog(CustomDialog):
    def __init__(self, matherial_value):
        w = Selector('color', title=f'Колір {matherial_value["name"]}', group_id=matherial_value['color_group_id'])
        super().__init__(w, "Обрати")        


class DetailsProductToOrderingTable(DetailsItemTable):
    actionResolved = pyqtSignal(dict, dict, bool, bool)
    def __init__(self, fields: list = [], values: list = None, buttons=TABLE_BUTTONS, group_id=0):
        super().__init__('product_to_ordering', '', fields, values, buttons, group_id)
        
    def prepare_value_to_save(self, value):
        return value
    
    def dialog(self, value, title):
        is_new = value['id'] == 0
        i = Item(self.item.name)
        num = value['number'] if 'number' in value else 0
        form = ProductToOrderingForm(value=value)
        dlg = CustomFormDialog(title, form)
        res = dlg.exec()
        if res and dlg.value:
            i.value = self.prepare_value_to_save(dlg.value)
            if dlg.form.product_value['min_cost'] * i.value['pieces'] > i.value['cost']:
                i.value['cost'] = dlg.form.product_value['min_cost'] * i.value['pieces']
            err = i.save()
            if err:
                error(err)
                return
            self.reload()
            num_changed = num and (num != value['number'])
            self.actionResolved.emit(i.value, form.product_value, num_changed, is_new)

    def del_dialog(self, value):
        dlg = DeleteDialog(value)
        res = dlg.exec()
        if res:
            i = Item(self.item.name)
            i.value = value
            cause = dlg.entry.text()
            if cause:
                if 'comm' in value:
                    i.value['comm'] = f'del: {cause}' + value['comm']
                    i.save()
                elif 'info' in value:
                    i.value['info'] = f'Причина видалення:\n{cause}\n' + value['info']
                    i.save()
            err = i.delete(value['id'])
            if err:
                error(err)
                return
            self.reload()
            self.actionResolved.emit(i.value, {}, False, False)


class ItemToCheckForm(CustomForm):
    def __init__(self, fields: list = [], value: dict = {}):
        self.item = Item('item_to_cbox_check')
        super().__init__(self.item.model, fields, value)
        self.widgets['price'].valChanged.connect(self.price_changed)
        self.widgets['number'].valChanged.connect(self.price_changed)
        self.widgets['cost'].valChanged.connect(self.cost_changed)

    def price_changed(self):
        number = self.widgets['number'].value()
        if not number:
            return
        cost = self.widgets['price'].value() * number
        self.widgets['cost'].set_value(cost)

    def cost_changed(self):
        number = self.widgets['number'].value()
        if not number:
            return
        price = round(self.widgets['cost'].value() / number, 3)
        self.widgets['price'].set_value(price)


class CheckFormDialog(CustomFormDialog):
    def __init__(self, ordering_value, m2os, o2os, p2os):
        self.counter = 1
        self.ordering_value = ordering_value
        table = self.create_table(m2os, o2os, p2os)
        form = self.create_form()
        super().__init__("Створити чек", form, table)
        table.actionInvoked.connect(self.action)
        
    def recalc(self):
        total = self.table.table.recalc('cost')
        self.form.widgets['cash_sum'].set_value(total)
    
    def action(self, action:str, value:dict=None):
        if action == 'delete':
            self.table.table.delete_values()
            self.recalc()
            return
        title = 'Створити'
        if action == 'copy':
            value['id'] = 0
        self.dialog(value, title)

    def dialog(self, value, title):
        i = Item('item_to_cbox_check')
        if not value:
            i.create_default_w()
            value = i.value
        form = ItemToCheckForm(value=value)
        dlg = CustomFormDialog(title, form)
        res = dlg.exec()
        if res and dlg.value:
            dlg.value["item_code"] = f"04-{self.counter}"
            dlg.value['measure'] = dlg.form.widgets['measure_id'].full_value()['name']
            self.counter += 1 
            self.table.table.add_value(dlg.value)
            self.recalc()

    def create_form(self):
        cbox_check = Item('cbox_check')
        cbox_check.create_default()
        price = round(self.ordering_value['price'])
        discount = round(self.ordering_value['profit'])
        cbox_check.value['cash_sum'] = price
        cbox_check.value['contragent_id'] = self.ordering_value['contragent_id']
        cbox_check.value['ordering_id'] = self.ordering_value['id']
        cbox_check.value['based_on'] = self.ordering_value['document_uid']
        cbox_check.value['name'] = f'Чек до зам. {self.ordering_value["id"]}'
        cbox_check.value['discount'] = discount
        check_form = CustomForm(cbox_check.model, value=cbox_check.value, fields=cbox_check.columns)
        is_prepay = CheckWidget()
        row = check_form.grid.rowCount() + 1
        check_form.grid.addWidget(QLabel("Передоплата"), row, 0)
        check_form.grid.addWidget(is_prepay, row, 1)
        check_form.is_prepay = is_prepay
        return check_form
        
    def create_table(self, m2os, o2os, p2os):
        app = App()
        price = round(self.ordering_value['price'])
        check_item = Item('item_to_cbox_check')
        fields = ["name", "number", "measure", "price", "discount", "cost", "item_code"]
        values = []
        m = Item('matherial')
        for m2o in m2os:
            if m2o["product_to_ordering_id"]:
                continue
            err = m.get_w(m2o['matherial_id'])
            if err:
                error(err)
                continue
            v = {}
            v["id"] = 0
            v["name"] = m2o['matherial'] 
            v["cbox_check_id"] = 0
            v["cbox_check"] = ''
            if m2o['width']:
                v["number"] = m2o['pieces']
                v["measure_id"] = app.config["measure pieces"]
                v["measure"] = 'шт.'
                size_txt = f" {m2o['width']}x{m2o['length']}мм"
                v['name'] += size_txt
                v['price'] = round(m2o['cost'] / v['number'])
                v["cost"] = round(m2o['cost'])
                v["discount"] = -m2o['profit'] / v['number']    
            else:                
                v["number"] = m2o['number']
                v["measure_id"] = m.value['measure_id']
                v["measure"] = m.value['measure']
                v["price"] = m2o['cost'] / m2o['number']
                v["discount"] = -m2o['profit']
                v["cost"] = m2o['cost']
            
            v["item_code"] = f"01-{m2o['matherial_id']}" 
            v["is_active"] = True
            values.append(v)

        for o2o in o2os:
            if o2o["product_to_ordering_id"]:
                continue
            v = {}
            v["id"] = 0
            v["name"] = o2o['operation']
            v["cbox_check_id"] = 0
            v["cbox_check"] = ''
            v["number"] = o2o['number']
            v["price"] = round(o2o['cost'] / o2o['number'], 2)
            v["discount"] = 0.0
            v["cost"] = o2o['cost']
            v["item_code"] = f"02-{o2o['operation_id']}" 
            v["is_active"] = True
            values.append(v)
        
        p = Item('product')
        for p2o in p2os:
            if p2o["product_to_ordering_id"]:
                continue
            err = p.get_w(p2o['product_id'])
            if err:
                error(err)
                continue
            v = {}
            v["id"] = 0
            v["name"] = p2o['product']
            v["cbox_check_id"] = 0
            v["cbox_check"] = ''
            v["number"] = p2o['number']
            if p2o['width']:
                v["number"] = p2o['pieces']
                v["measure_id"] = app.config["measure pieces"]
                v["measure"] = 'шт.'
                size_txt = f"{p2o['width']}x{p2o['length']}"
                if not size_txt in v['name']:
                    v['name'] += f" {size_txt}мм"
                v["discount"] = -p2o['profit'] / v['number']
                v['price'] = round(p2o['cost'] / v['number'] + v["discount"])
                v["cost"] = round(p2o['cost'])
                v["discount"] = -p2o['profit'] / v['number'] 
            else:    
                v["number"] = p2o['number']
                v["measure"] = p.value['measure']
                v["measure_id"] = p.value['measure_id']
                v["discount"] = -p2o['profit']
                v["price"] = p2o['cost'] / p2o['number'] + v["discount"]
                
                v["cost"] = p2o['cost']
            
            v["item_code"] = f"03-{p2o['product_id']}" 
            v["is_active"] = True
            values.append(v)

        for v in values:
            price -= v['cost']
        if price:
            values[-1]['cost'] += price
            values[-1]['price'] = round(values[-1]['cost'] / values[-1]['number'], 3)
        
        check_table = TableWControls(
            check_item.model_w, 
            values=values, 
            table_fields=fields,
            buttons={'create':'Створити', 'copy':'Копіювати', 'delete':'Видалити'},
            )
        return check_table


class ItemToInvoiceForm(CustomForm):
    def __init__(self, fields: list = [], value: dict = {}):
        self.item = Item('item_to_invoice')
        super().__init__(self.item.model, fields, value)
        self.widgets['price'].valChanged.connect(self.price_changed)
        self.widgets['number'].valChanged.connect(self.price_changed)
        self.widgets['cost'].valChanged.connect(self.cost_changed)

    def price_changed(self):
        number = self.widgets['number'].value()
        if not number:
            return
        cost = self.widgets['price'].value() * number
        self.widgets['cost'].set_value(cost)

    def cost_changed(self):
        number = self.widgets['number'].value()
        if not number:
            return
        price = round(self.widgets['cost'].value() / number, 3)
        self.widgets['price'].set_value(price)


class InvoiceFormDialog(CustomFormDialog):
    def __init__(self, ordering_value, m2os, o2os, p2os):
        self.ordering_value = ordering_value
        table = self.create_table(m2os, o2os, p2os)
       
        form = self.create_form()
        self.is_signed = QCheckBox('З підписом')
        self.with_date = QCheckBox('З датою')
        super().__init__("Створити рахунок", form, table)
        form.grid.addWidget(self.with_date)
        form.grid.addWidget(self.is_signed)
        table.actionInvoked.connect(self.action)
        
    def recalc(self):
        total = self.table.table.recalc('cost')
        self.form.widgets['cash_sum'].set_value(total)
    
    def action(self, action:str, value:dict=None):
        if action == 'delete':
            self.table.table.delete_values()
            self.recalc()
            return
        title = 'Створити'
        if action == 'copy':
            value['id'] = 0
        self.dialog(value, title)

    def dialog(self, value, title):
        i = Item('item_to_invoice')
        if not value:
            i.create_default_w()
            value = i.value
        form = ItemToInvoiceForm(value=value, fields=i.columns)
        dlg = CustomFormDialog(title, form)
        res = dlg.exec()
        if res and dlg.value:
            dlg.value['measure'] = dlg.form.widgets['measure_id'].full_value()['name']
            self.table.table.add_value(dlg.value)
            self.recalc()

    def create_form(self):
        invoice = Item('invoice')
        invoice.create_default()
        invoice.value['cash_sum'] = round(self.ordering_value['cost'])
        invoice.value['contragent_id'] = self.ordering_value['contragent_id']
        invoice.value['contact_id'] = self.ordering_value['contact_id']
        invoice.value['based_on'] = self.ordering_value['document_uid']
        invoice.value['ordering_id'] = self.ordering_value['id']
        invoice.value['name'] = f'Рахунок до зам. {self.ordering_value["id"]}'
        invoice.value['discount'] = -self.ordering_value['profit']
        form = CustomForm(invoice.model, value=invoice.value, fields=invoice.columns)
        return form
        
    def create_table(self, m2os, o2os, p2os):
        app = App()
        cost = round(self.ordering_value['cost'])
        price = self.ordering_value['price']
        
        k = cost/price
        invoice_item = Item('item_to_invoice')
        fields = ["name", "number", "measure", "price", "cost"]
        values = []
        m = Item('matherial')
        for m2o in m2os:
            if m2o["product_to_ordering_id"]:
                continue
            err = m.get_w(m2o['matherial_id'])
            if err:
                error(err)
                continue
            v = {}
            v["id"] = 0
            v["name"] = m2o['matherial']
            v["invoice_id"] = 0
            v["invoice"] = ''
            
            if m2o['width']:
                v["number"] = m2o['pieces']
                v["measure_id"] = app.config["measure pieces"]
                v["measure"] = 'шт.'
                size_txt = f" {m2o['width']}x{m2o['length']}мм"
                v['name'] += size_txt
            else:                
                v["number"] = m2o['number']
                v["measure_id"] = m.value['measure_id']
                v["measure"] = m.value['measure']
            v["cost"] = round(m2o['cost'] * k)
            v["price"] = round(v['cost'] / v['number'], 3)
            v["is_active"] = True
            values.append(v)

        o = Item('operation')
        for o2o in o2os:
            if o2o["product_to_ordering_id"]:
                continue
            err = o.get_w(o2o['operation_id'])
            if err:
                error(err)
                continue
            v = {}
            v["id"] = 0
            v["name"] = o2o['operation']
            v["invoice_id"] = 0
            v["invoice"] = ''
            v["number"] = o2o['number']
            v["measure"] = o.value['measure']
            v["measure_id"] = o.value['measure_id']
            v["cost"] = o2o['cost'] * k
            v["price"] = round(v['cost'] / o2o['number'], 2)
            v["is_active"] = True
            values.append(v)

        p = Item('product')
        for p2o in p2os:
            if p2o["product_to_ordering_id"]:
                continue
            err = p.get_w(p2o['product_id'])
            if err:
                error(err)
                continue
            v = {}
            v["id"] = 0
            v["name"] = p2o['name']
            v["invoice_id"] = 0
            v["invoice"] = ''
            if p2o['width']:
                v["number"] = p2o['pieces']
                v["measure_id"] = app.config["measure pieces"]
                v["measure"] = 'шт.'
                size_txt = f"{p2o['width']}x{p2o['length']}"
                if not size_txt in v['name']:
                    v['name'] += f" {size_txt}мм"
            else:    
                v["number"] = p2o['number']
                v["measure"] = p.value['measure']
                v["measure_id"] = p.value['measure_id']
            v["cost"] = p2o['cost'] * k
            v["price"] = round(v['cost'] / v['number'], 3)
            v["is_active"] = True
            values.append(v)
                    
        for v in values:
            cost -= v['cost']
        if cost:
            values[-1]['cost'] += cost
            values[-1]['price'] = round(values[-1]['cost'] / values[-1]['number'], 3)
        table = TableWControls(
            invoice_item.model_w, 
            values=values, 
            table_fields=fields,
            buttons={'create':'Створити', 'copy':'Копіювати', 'delete':'Видалити'},
            )
        return table


class ContragentGreateFormDialog(CustomDialog):
    def __init__(self, contragent: Item, contact: Item) -> None:
        self.contact_value = {}
        self.value = {}
        self.box = QVBoxLayout()
        self.forms = QWidget()
        self.forms.setLayout(self.box)
        self.box.addWidget(QLabel('Контрагент'))
        self.contragent = contragent
        contragent_fields = [
            "name",
            "contragent_group_id",
            "phone",
            "email",
            "web",
            "comm",
            "dir_name",
            "full_name",
            "edrpou",
            "ipn",
            "iban",
            "bank",
            "mfo",
            "fop",
            "address",
        ]
        self.contragent_form = CustomForm(
            data_model=self.contragent.model,
            fields=contragent_fields,
            value=self.contragent.value,
            columns=2,
            )
        self.box.addWidget(self.contragent_form)
        self.contragent_form.hide_save_btn()
        self.box.addWidget(QLabel('Контакт'))
        self.contact = contact
        contact_fields = [
            "name",
            "phone",
            "email",
            "viber",
            "telegram",
            "comm",
        ]
        self.contact_form = CustomForm(
            data_model=self.contact.model,
            fields=contact_fields,
            value=self.contact.value,
            )
        self.box.addWidget(self.contact_form)
        self.contact_form.hide_save_btn()
        
        super().__init__(self.forms, 'Створити контрагента і контакт')
        self.contact_form.saveRequested.connect(self.get_contact_value)
        self.contragent_form.saveRequested.connect(self.get_contragent_value)

    def get_contact_value(self, value):
        self.contact_value = value
        self.contact_form.set_changed(False)
        
    def get_contragent_value(self, value):
        self.value = value
        self.contragent_form.set_changed(False)

    def accept(self) -> None:
        if self.contragent_form.get_value() and self.contact_form.get_value():
            return super().accept()


class WhsInFormDialog(CustomFormDialog):
    def __init__(self, cur_value, m2o_values):
        whs_in = Item('whs_in')
        whs_in.create_default()
        whs_in.value["based_on"] = cur_value['document_uid']
        whs_in.value["contragent_id"] = cur_value["contragent_id"]
        whs_in.value["contact_id"] = cur_value["contact_id"]
        whs_in.value["comm"] = f'авт. до {cur_value["name"]}'
                
        self.m2o = Item('matherial_to_ordering')
        fields = ('matherial', 'number', 'price', 'cost', 'color')
        self.m2o.model.update(self.m2o.model_w)
        table = TableWControls(self.m2o.model, fields, m2o_values, {'edit':'Редагувати', 'delete':'Видалити'})
        whs_in.value['whs_sum'] = table.table._model.calc_sum('cost')
        table.actionInvoked.connect(self.apply)
        form = CustomForm(whs_in.model, whs_in.columns, whs_in.value)
        super().__init__('Створити ПН', form, table)

    def apply(self, action, value):
        if not value:
            return
        if action == 'delete':
            self.table.table.delete_values()
        elif action == 'edit':
            new_cost = askdlg('Вкажіть вартість:')
            if not new_cost:
                return
            try:
                cost = float(new_cost)
            except:
                error('Це не число, спробуйте ще!')
                return
            value['cost'] = cost
            value['price'] = round(cost / value['number'], 2)
            self.m2o.value = value
            err = self.m2o.save()
            if err:
                error(err)
                return
            self.table.table.delete_values()
            self.table.table._model.append(value)
        else:
            return

        cost = self.table.table._model.calc_sum('cost')
        self.form.widgets['whs_sum'].set_value(cost)


class ItemsToOrdering(QSplitter):
    viewContragentReport = pyqtSignal(int)
    def __init__(self, checkbox) -> None:
        super().__init__()
        self.check = checkbox
        self.current_ordering_value = {}
        self.current_contragent = {}
        self.current_contact = {}
        self.ordering = Item('ordering')
        order_state = Item('ordering_status')
        err = order_state.get_all_w()
        if err:
            error(err)
            return
        
        aside = QWidget()
        self.addWidget(aside)
        self.aside_box = QVBoxLayout()
        self.aside_box.setContentsMargins(0, 0, 0, 0)
        aside.setLayout(self.aside_box)
        
        self.side_controls = QTabWidget()
        self.aside_box.addWidget(self.side_controls)
        self.buttons = ButtonsBlock('За статусом', order_state.values)
        self.side_controls.addTab(self.buttons, "За станом")
        self.buttons.buttonClicked.connect(self.on_state_selected)
        self.contragent_filter = ContragentFilter()
        self.side_controls.addTab(self.contragent_filter, "За контрагентом")
        self.contragent_filter.contragentChanged.connect(self.current_contragent_changed)
        self.contragent_filter.contactChanged.connect(self.current_contact_changed)
        
        to_folder = QPushButton('До теки')
        self.buttons.box.addWidget(to_folder)
        to_folder.clicked.connect(self.to_folder)
        to_contragent = QPushButton('До контрагенту')
        self.buttons.box.addWidget(to_contragent)
        to_contragent.clicked.connect(self.to_contragent)
        
        cont_box = QHBoxLayout()
        self.contragent_filter.box.addLayout(cont_box)
        contragent_plus_btn = QPushButton('+Контрагент')
        cont_box.addWidget(contragent_plus_btn)
        contragent_plus_btn.clicked.connect(self.add_contragent)
        contact_plus_btn = QPushButton('+Контакт')
        cont_box.addWidget(contact_plus_btn)
        contact_plus_btn.clicked.connect(self.add_contact)
        
        messengers_box = QHBoxLayout()
        self.aside_box.addLayout(messengers_box)
        viber_btn = QPushButton('Viber')
        messengers_box.addWidget(viber_btn)
        viber_btn.clicked.connect(self.to_viber)
        tg_btn = QPushButton('Telegram')
        messengers_box.addWidget(tg_btn)
        tg_btn.clicked.connect(self.to_telegram)
        rep_contragent = QPushButton('Звіт контрагента')
        self.aside_box.addWidget(rep_contragent)
        rep_contragent.clicked.connect(self.to_contragent_report)
        btns_grid = QWidget()
        self.grid = QGridLayout()
        self.grid.setContentsMargins(0, 0, 0, 0)
        self.grid.setVerticalSpacing(1)
        btns_grid.setLayout(self.grid)
        self.aside_box.addWidget(btns_grid)

        fields = [
            "id",
            "name",
            "ordering_state",
            "contragent",
            "user",
            "deadline_at",
            "ordering_status",
        ]
        self.orderings = OrderingTable(fields=fields)
        self.addWidget(self.orderings)
        
        self.tabs = QTabWidget()
        self.addWidget(self.tabs)
        
        self.details_tree = TreeItemToOrdering(is_info_bottom=True)
        self.tabs.addTab(self.details_tree, "Детально")
        self.details_tree.orderingChanged.connect(self.update_sum_by_tree)

        fields = [
            "matherial",
            "number",
            "price",
            "cost",
            "color",
            "product_to_ordering",
        ]
        self.m2os = DetailsMatherialToOrderingTable(fields=fields)
        self.m2os.set_main_table(self.orderings)
        self.tabs.addTab(self.m2os, "Матеріал")
        
        fields = [
            "operation",
            "cost",
            "user",
            "comm",
            "product_to_ordering",
        ]
        self.o2os = DetailsOperationToOrderingTable(fields=fields)
        self.o2os.set_main_table(self.orderings)
        self.tabs.addTab(self.o2os, "Операція")
        
        fields = [
            "product",
            "number",
            "price",
            "cost",
            "product_to_ordering_status",
        ]
        self.p2os = DetailsProductToOrderingTable(fields=fields)
        self.p2os.set_main_table(self.orderings)
        self.tabs.addTab(self.p2os, "Виріб")
        
        self.setStretchFactor(0, 2)
        self.setStretchFactor(1, 3)
        
        self.orderings.table.table.valueSelected.connect(self.on_ordering_selected)
        self.m2os.actionResolved.connect(self.update_sum)
        self.o2os.actionResolved.connect(self.update_sum)
        self.p2os.actionResolved.connect(self.on_product_changed)

        self.doc_table = DocsTable()
        self.orderings.add_doc_table(self.doc_table)

        cash_out_btn = QPushButton('Створити ВКО')
        cash_in_btn = QPushButton('Створити ПКО')
        whs_out_btn = QPushButton('Створити ВН')
        whs_in_btn = QPushButton('Створити ПН')
        check_btn = QPushButton('Створити чек')
        invoice_btn = QPushButton('Створити рахунок')
        self.grid.addWidget(cash_out_btn, 0, 0)
        self.grid.addWidget(cash_in_btn, 0, 1)
        self.grid.addWidget(whs_out_btn, 1, 0)
        self.grid.addWidget(whs_in_btn, 1, 1)
        self.grid.addWidget(check_btn, 2, 0)
        self.grid.addWidget(invoice_btn, 2, 1)
        cash_out_btn.clicked.connect(self.create_cash_out)
        cash_in_btn.clicked.connect(self.create_cash_in)
        whs_out_btn.clicked.connect(self.create_whs_out)
        whs_in_btn.clicked.connect(self.create_whs_in1)
        check_btn.clicked.connect(self.create_check)
        invoice_btn.clicked.connect(self.create_invoice)

        self.total_cash = QLabel()
        self.total_whs = QLabel()
        self.grid.addWidget(QLabel("Сплачено"), 3, 0)
        self.grid.addWidget(QLabel("По складу"), 4, 0)
        self.grid.addWidget(self.total_cash, 3, 1)
        self.grid.addWidget(self.total_whs, 4, 1)
    
    def to_folder(self):
        if not self.current_ordering_value:
            error('Оберіть замовлення')
            return
        app = App()
        base_dir = app.config['new_makets_path']
        contragent_id = self.current_ordering_value['contragent_id']
        contragent = Item('contragent')
        err = contragent.get(contragent_id)
        if err:
            error(err)
            return
        path = os.path.join(
                base_dir,
                contragent.value['dir_name'],
                str(self.current_ordering_value['id']),
            )
        
        subprocess.run([app.config['program'], path])
    
    def to_contragent(self):
        contragent_id = self.current_ordering_value['contragent_id']
        self.contragent_filter.set_contragent_by_id(contragent_id)
        self.side_controls.setCurrentIndex(1)

    def get_current_contragent_id(self):
        if self.side_controls.currentIndex() == 1 and self.current_contragent:
            return self.current_contragent['id']
        elif self.current_ordering_value:
            return self.current_ordering_value['contragent_id']
        else:
            error('Оберіть замовлення або контрагента')
            return 0
        
    def get_current_contact_id(self):
        if self.side_controls.currentIndex() == 1 and self.current_contact:
            return self.current_contact['id']
        elif self.current_ordering_value:
            return self.current_ordering_value['contact_id']
        else:
            error('Оберіть замовлення або контакт')
            return 0
        
    def get_current_contact_value(self):
        if self.side_controls.currentIndex() == 1 and self.current_contact:
            return self.current_contact
        elif self.current_ordering_value:
            contact_id = self.current_ordering_value['contact_id']
            if not contact_id:
                error('Оберіть замовлення або контакт')
                return {}
            contact = Item('contact')
            err = contact.get(contact_id)
            if err:
                error(err)
                return {}
            return contact.value

    def to_contragent_report(self):
        contragent_id = self.get_current_contragent_id()
        if contragent_id:
            self.viewContragentReport.emit(contragent_id)
    
    def to_viber(self):
        contact_value = self.get_current_contact_value()
        if not contact_value:
            return
        url = ''
        if contact_value['viber']:
            url = f'viber://chat?number={contact_value["viber"]}'
        else:
            phone = ''.join(i for i in contact_value['phone'] if i.isdigit())
            if len(phone) >= 9:
                phone = phone[-9:]
            else:
                return
            url = f'viber://chat?number={phone}'
        if url:
            webbrowser.open(url)

    def to_telegram(self):
        contact_value = self.get_current_contact_value()
        if not contact_value:
            return
        url = ''
        if contact_value['telegram']:
            url = f'tg://resolve?domain={contact_value["telegram"]}'
        else:
            phone = ''.join(i for i in contact_value['phone'] if i.isdigit())
            if len(phone) >= 9:
                phone = '380' + phone[-9:]
            else:
                return
            url = f'tg://resolve?phone={phone}'
        if url:
            webbrowser.open(url)

    def add_contragent(self):
        contragent = Item('contragent')
        contragent.create_default()
        contact = Item('contact')
        contact.create_default()
        dlg = ContragentGreateFormDialog(contragent, contact)
        res = dlg.exec()
        if not res:
            return
        contragent.value = dlg.value
        if not contragent.value:
            return
        err = contragent.save()
        if err:
            error(err)
            return
        contact.value = dlg.contact_value
        if not contact.value:
            return
        contact.value['contragent_id'] = contragent.value['id']  
        err = contact.save()
        if err:
            error(err)
            return
        self.contragent_filter.set_contragent_by_id(contragent.value['id'])

    def add_contact(self):
        contragent_id = self.get_current_contragent_id()
        if not contragent_id:
            return
        dlg = ContactSelectDialog(contragent_id=contragent_id)
        dlg.exec()

    def current_contact_changed(self, contact_value):
        self.current_contact = contact_value
        self.ordering.get_filter_w('contact_id', contact_value['id'])
        self.reload(self.ordering.values)
    
    def current_contragent_changed(self, contragent_value):
        self.current_contragent = contragent_value
        self.current_contact = {}
        self.ordering.get_filter_w('contragent_id', contragent_value['id'])
        self.reload(self.ordering.values)
    
    def reload(self, values=None):
        if values is None and self.buttons.current_value:
            err = self.ordering.get_filter_w('ordering_status_id', self.buttons.current_value['id'])
            if not err:
                self.orderings.reload(self.ordering.values)
                return
            else:
                error(err)
        self.orderings.reload(values)
        self.doc_table.reload([])
        self.details_tree.reload()

    def on_state_selected(self, state_value):
        app = App()
        if state_value['id'] == app.config["ordering state taken"] or state_value['id'] == app.config["ordering state canceled"]:
            date_from, date_to = self.orderings.period.get_period()
            err = self.ordering.get_between_w('created_at', date_from, date_to)
            if err:
                error(err)
                return
            values = [v for v in self.ordering.values if v['ordering_status_id'] == state_value['id']]
            self.reload(values)
            return
        err = self.ordering.get_filter_w('ordering_status_id', state_value['id'])
        if err:
            error(err)
            return
        self.reload(self.ordering.values)

    def create_invoice(self):
        cur_value = self.orderings.table.table.get_selected_value()
        if not cur_value:
            return
        if not cur_value['price']:
            error("Ціна не може дорівнювати нулю")
            return
        m2os = self.m2os.values() 
        o2os = self.o2os.values() 
        p2os = self.p2os.values()
        if not (len(m2os) + len(o2os) + len(p2os)):
            error("Неможливо створити рахунок без виробів")
            return
        dlg = InvoiceFormDialog(cur_value, m2os, o2os, p2os)
        if not dlg:
            return
        res = dlg.exec()
        if res:
            invoice = Item('invoice')
            invoice.value = dlg.value
            invoice.value['cash_sum'] = 0
            err = invoice.save()
            if err:
                error(err)
                return
            invoice_item = Item('item_to_invoice')
            for v in dlg.table.table.values():
                v['invoice_id'] = invoice.value['id']
                invoice_item.value = v
                err = invoice_item.save()
                if err:
                    error(err)
            owner = dlg.form.widgets['owner_id'].full_value()
            contragent = dlg.form.widgets['contragent_id'].full_value()
            err = invoice_item.get_filter_w('invoice_id', invoice.value['id'])
            if err:
                error(err)
                return
            is_signed = dlg.is_signed.checkState() == Qt.CheckState.Checked
            with_date = dlg.with_date.checkState() == Qt.CheckState.Checked
            app = App()
            base_dir = app.config['new_makets_path']
            path = os.path.join(
                base_dir,
                contragent['dir_name'],
                str(self.current_ordering_value['id']),
                'documents',
                f"{invoice.value['name']}.xlsx"
            )
            doc = ExcelDoc(owner, invoice.value, invoice_item.values, contragent)
            doc.create_docs(path, is_signed, with_date)

    def create_cash_out(self):
        cur_value = self.orderings.table.table.get_selected_value()
        if not cur_value:
            return
        
        cash_out = Item('cash_out')
        cash_out.create_default()
        cash_out.value["based_on"] = cur_value['document_uid']
        cash_out.value["contragent_id"] = cur_value["contragent_id"]
        cash_out.value["contact_id"] = cur_value["contact_id"]
        cash_out.value["cash_sum"] = cur_value["cost"]
        cash_out.value["comm"] = f'авт. до {cur_value["name"]}'

        dlg = FormDialog('Створити ВКО', cash_out.model, cash_out.columns, cash_out.value)
        res = dlg.exec()
        if res:
            err = cash_out.save()
            if err:
                error(err)
                return
            self.doc_table.reload(cur_value)

    def create_cash_in(self):
        cur_value = self.orderings.table.table.get_selected_value()
        if not cur_value:
            return
        payed = self.doc_table.calc_by_type('cash_in')
        cash_in = Item('cash_in')
        cash_in.create_default()
        cash_in.value["based_on"] = cur_value['document_uid']
        cash_in.value["contragent_id"] = cur_value["contragent_id"]
        cash_in.value["contact_id"] = cur_value["contact_id"]
        cash_in.value["comm"] = f'авт. до {cur_value["name"]}'
        cash_in.value["cash_sum"] = round(cur_value["cost"] - payed)

        dlg = FormDialog('Створити ПКО', cash_in.model, cash_in.columns, cash_in.value)
        res = dlg.exec()
        if res:
            err = cash_in.save()
            if err:
                error(err)
                return
            self.doc_table.reload(cur_value)
        
    def create_whs_out(self):
        app = App()
        cur_value = self.orderings.table.table.get_selected_value()
        if not cur_value:
            return
        whs_out = Item('whs_out')
        whs_out.create_default()
        whs_out.value["based_on"] = cur_value['document_uid']
        whs_out.value["contragent_id"] = app.config["contragent to production"]
        whs_out.value["contact_id"] = app.config["contact to production"]
        whs_out.value["comm"] = f'авт. до {cur_value["name"]}'
        dlg = FormDialog('Створити ВН', whs_out.model, whs_out.columns, whs_out.value)
        res = dlg.exec()
        if not res:
            return
        err = whs_out.save()
        if err:
            error(f'При збереженні накладної:\n{err}')
            return
        m2o_values = self.m2os.table.table.get_selected_values()
        if not m2o_values:
            m2o_values = self.m2os.values()
        grp_values = self.group_matherials(m2o_values)
        whs_sum = 0
        mat = Item('matherial')
        for v in grp_values.values():
            err = mat.get(v['matherial_id'])
            if err:
                error(f'При завантаженні матеріалу:\n{err}')
            if mat.value['count_type_id'] == VIRTUAL:
                continue
            m2wo = Item('matherial_to_whs_out')
            m2wo.create_default()
            m2wo.value["whs_out_id"] = whs_out.value['id']
            m2wo.value["matherial_id"] = v["matherial_id"]
            m2wo.value["number"] = v["number"]
            m2wo.value["price"] = mat.value["price"]
            m2wo.value["cost"] = round(v["number"] * mat.value["price"], 2)
            m2wo.value["width"] = v["width"]
            m2wo.value["length"] = v["length"]
            m2wo.value["color_id"] = v["color_id"]
            err = m2wo.save()
            if err:
                error(f'При збереженні матеріалу до ПН:\n{err}')
                continue
            whs_sum += m2wo.value["cost"]
        whs_out.value['whs_sum'] = whs_sum
        err = whs_out.save()
        if err:
            error(f'При збереженні суми накладної:\n{err}')
        self.doc_table.reload(cur_value)

    def group_matherials(self, m2o_values):
        res_values = {}
        for v in m2o_values:
            if (not v['matherial_id'] in res_values) or v["length"] or v["color_id"]:
                res_values[v['matherial_id']] = v.copy()
            else:
                res_values[v['matherial_id']]["number"] += v["number"]
        return res_values
    
    def filter_virtual_materials(self, m2o_values):
        mat = Item('matherial')
        res = []
        for v in m2o_values:
            err = mat.get(v['matherial_id'])
            if err:
                error(f'При завантаженні матеріалу:\n{err}')
                continue
            if mat.value['count_type_id'] == VIRTUAL:
                continue
            res.append(v)
        return res

    def create_whs_in1(self):
        cur_value = self.orderings.table.table.get_selected_value()
        if not cur_value:
            return
        m2o_values = self.m2os.table.table.get_selected_values()
        if not m2o_values:
            m2o_values = self.m2os.values()
        m2o_values = self.filter_virtual_materials(m2o_values)
        dlg = WhsInFormDialog(cur_value, m2o_values)
        res = dlg.exec()
        if not res:
            return
        whs_in = Item('whs_in')
        whs_in.value = dlg.value
        whs_in.value['whs_sum'] = 0
        err = whs_in.save()
        if err:
            error(f'При збереженні накладної:\n{err}')
            return
        m2o_values = dlg.table.table.values()
        print(m2o_values)
        if not m2o_values:
            return
        
        for v in m2o_values:
            m2wi = Item('matherial_to_whs_in')
            m2wi.create_default()
            m2wi.value["whs_in_id"] = whs_in.value['id']
            m2wi.value["matherial_id"] = v["matherial_id"]
            m2wi.value["number"] = v["number"]
            price = v['price']
            m2wi.value["price"] = price
            m2wi.value["cost"] = round(v["number"] * price, 2)
            m2wi.value["width"] = v["width"]
            m2wi.value["length"] = v["length"]
            m2wi.value["color_id"] = v["color_id"]
            err = m2wi.save()
            if err:
                error(f'При збереженні матеріалу до ПН:\n{err}')
            
        self.doc_table.reload(cur_value)

    def create_whs_in(self):
        cur_value = self.orderings.table.table.get_selected_value()
        if not cur_value:
            return

        whs_in = Item('whs_in')
        whs_in.create_default()
        whs_in.value["based_on"] = cur_value['document_uid']
        whs_in.value["contragent_id"] = cur_value["contragent_id"]
        whs_in.value["contact_id"] = cur_value["contact_id"]
        whs_in.value["comm"] = f'авт. до {cur_value["name"]}'
        
        dlg = FormDialog('Створити ПН', whs_in.model, whs_in.columns, whs_in.value)
        res = dlg.exec()
        if not res:
            return
        err = whs_in.save()
        if err:
            error(f'При збереженні накладної:\n{err}')
            return
        m2o_values = self.m2os.table.table.get_selected_values()
        if not m2o_values:
            m2o_values = self.m2os.values()
        whs_sum = 0
        for v in m2o_values:
            mat = Item('matherial')
            err = mat.get(v['matherial_id'])
            if err:
                error(f'При завантаженні матеріалу:\n{err}')
            
            m2wi = Item('matherial_to_whs_in')
            m2wi.create_default()
            m2wi.value["whs_in_id"] = whs_in.value['id']
            m2wi.value["matherial_id"] = v["matherial_id"]
            m2wi.value["number"] = v["number"]
            price = mat.value["price"] if mat.value['price'] else v['price']
            m2wi.value["price"] = price
            m2wi.value["cost"] = round(v["number"] * price, 2)
            m2wi.value["width"] = v["width"]
            m2wi.value["length"] = v["length"]
            m2wi.value["color_id"] = v["color_id"]
            err = m2wi.save()
            if err:
                error(f'При збереженні матеріалу до ПН:\n{err}')
            whs_sum += m2wi.value["cost"]
        whs_in.value['whs_sum'] = whs_sum
        err = whs_in.save()
        if err:
            error(f'При збереженні суми накладної:\n{err}')
            
        self.doc_table.reload(cur_value)

    def create_check(self):
        cur_value = self.orderings.table.table.get_selected_value()
        if not cur_value:
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
        if not cash_state['has_shift']:
            shift_data = self.check.shifts()
            if not shift_data:
                error("Не можу відкрити зміну у CheckBox!")
                return
            if 'error' in shift_data:
                error(f"При відкритті зміни у Checkbox: {shift_data['error']}")
                return 
        
        dlg = CheckFormDialog(cur_value, self.m2os.values(), self.o2os.values(), self.p2os.values())
        res = dlg.exec()
        if res:
            check = Item('cbox_check')
            check.value = dlg.value
            if dlg.form.is_prepay.value():
                check.value['comm'] = 'Передоплата'
            err = check.save()
            if err:
                error(err)
                return
            check_item = Item('item_to_cbox_check')
            for v in dlg.table.table.values():
                v['cbox_check_id'] = check.value['id']
                check_item.value = v
                err = check_item.save()
                if err:
                    error(err)
                    return
            
            receipt = self.create_receipt(check)
            if receipt is None:
                return
            if dlg.form.is_prepay.value():
                relation_id = "{price: 011n}".format(price = cur_value['id'])
                receipt['custom_relation_id'] = relation_id
                res = self.check.create_pre_receipt(receipt)

            else:
                res = self.check.create_receipt(receipt)
            if not res:
                error("Помилка при створенні чеку в Checkbox")
                return
            if 'error' in res:
                error(f"При створенні чеку в Checkbox: {res['error']}")
                return
            check.value["checkbox_uid"] = res['id']
            err = check.save()
            if err:
                error(err)
                return
        self.reload()
        

    # add_cashless_sum used if combined cash and cashless, 
    # is_cash must be True, add_cashless_sum included in cash_sum
    def create_receipt(self, check:Item, add_cashless_sum: float=0.0):
        payments = []
        if check.value['is_cash']:
            if add_cashless_sum:
                payments.append({
                  "type": 'CASHLESS',
                    "value": int(add_cashless_sum * 100),
                    "label": 'Картка',
                })    
            payments.append({
                "type": 'CASH',
                "value": int((check.value['cash_sum'] - add_cashless_sum) * 100),
                "label": 'Готівка',
            })
        else:
            payments.append({
                "type": 'CASHLESS',
                "value": int(check.value['cash_sum'] * 100),
                "label": 'Картка',
            })
        
        receipt = {
            "department": "Копіцентр",
            "goods": [],
            "payments": payments,
        }
        check_item = Item('item_to_cbox_check')
        err = check_item.get_filter_w('cbox_check_id', check.value['id'])
        if err:
            error(err)
            return

        for check_item_value in check_item.values:
            good = {
                        "good": {
                            "code": check_item_value['item_code'],
                            "name": check_item_value['name'],
                            "price": int(check_item_value['price'] * 100),
                        },
                        "quantity": int(check_item_value['number']*1000),
                    }
            
            
            if check_item_value['discount'] > 0.01:
                dsc = round(check_item_value['discount']*check_item_value['number'], 0)
                good["discounts"] = self.make_discount(dsc)
            receipt['goods'].append(good)
        
        if check.value['discount'] > 0.01:
            receipt["discounts"] = self.make_discount(check.value['discount'])
            
        return receipt
    
    def make_discount(self, value):
        if value < 0:
            value = -value
            dsct_type = 'EXTRA_CHARGE'
            dsct_name = "Націнка"
        else:
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

    def on_ordering_selected(self, value):
        self.current_ordering_value = value
        self.details_tree.reload(value['id'])
        m2o = Item('matherial_to_ordering')
        err = m2o.get_filter_w('ordering_id', value['id'])
        if err:
            error(err)
            return
        self.m2os.reload(m2o.values)
        o2o = Item('operation_to_ordering')
        err = o2o.get_filter_w('ordering_id', value['id'])
        if err:
            error(err)
            return
        self.o2os.reload(o2o.values)
        p2o = Item('product_to_ordering')
        err = p2o.get_filter_w('ordering_id', value['id'])
        if err:
            error(err)
            return
        self.p2os.reload(p2o.values)
        self.total_whs.setNum(self.doc_table.calc_by_type('whs_in') - self.doc_table.calc_by_type('whs_out'))
        self.total_cash.setNum(self.doc_table.calc_by_type('cash_in') - self.doc_table.calc_by_type('cash_out')) 

    def update_sum(self):
        m2o_values = self.m2os.values()
        total = 0
        for v in m2o_values:
            if not v['product_to_ordering_id']:
                total += v["cost"]
        o2o_values = self.o2os.values()
        for v in o2o_values:
            if not v['product_to_ordering_id']:
                total += v["cost"]
        self.p2os.reload()
        p2o_values = self.p2os.values()
        for v in p2o_values:
            if not v['product_to_ordering_id']:
                total += v["cost"]
        self.current_ordering_value['price'] = total
        self.current_ordering_value['cost'] = total
        self.current_ordering_value['persent'] = 0
        self.current_ordering_value['profit'] = 0
        order = Item('ordering')
        order.value = self.current_ordering_value
        err = order.save()
        if err:
            error(err)
            return
        rows = self.orderings.table.table.get_selected_rows()
        if not rows:
            return
        row = rows[0]
        self.orderings.reload()
        index = self.orderings.table.table._model.createIndex(row, 0)
        self.orderings.table.table.setCurrentIndex(index)

    def update_sum_by_tree(self):
        new_sum = self.details_tree.deep_calc_sum()
        self.current_ordering_value['price'] = new_sum
        self.current_ordering_value['profit'] = round(new_sum * self.current_ordering_value['persent']/100, 2)
        self.current_ordering_value['cost'] = round(new_sum + self.current_ordering_value['profit'], 2)
        order = Item('ordering')
        order.value = self.current_ordering_value
        err = order.save()
        if err:
            error(err)
            return
        rows = self.orderings.table.table.get_selected_rows()
        row = rows[0] if rows else 0
        self.orderings.reload()
        index = self.orderings.table.table._model.createIndex(row, 0)
        self.orderings.table.table.setCurrentIndex(index)


    def clear_p2o_positions(self, p2o_id):
        for name in ('product', 'operation', 'matherial'):
            name += '_to_ordering'
            item = Item(name)
            err = item.get_filter_w('product_to_ordering_id', p2o_id)
            if not err:
                for v in item.values:
                    item.value = v
                    err = item.delete(v['id'])
                    if err:
                        error(err)
                        continue

    def create_m2o(self, ms2p, m2o, p2o_value):
        total = 0.0
        matherial = Item('matherial')
        for m in ms2p:
            m2o.value["matherial_id"] = m['matherial_id']
            m2o.value["number"] = m['number'] * p2o_value['number']
            m2o.value["cost"] = m['cost'] * p2o_value['number']
            m2o.value["price"] = m['cost'] / m['number']
            err = matherial.get(m['matherial_id'])
            if err:
                error(err)
                continue
            if matherial.value['color_group_id']:
                cdlg = ColorDialog(matherial_value=matherial.value)
                res = cdlg.exec()
                if res:
                    m2o.value['color_id'] = cdlg.widget.value()
            m2o.value['width'] = p2o_value['width']
            m2o.value['length'] = p2o_value['length']
            m2o.value['pieces'] = p2o_value['pieces']
            
            err = m2o.save()
            if err:
                error(err)
            m2o.value['id'] = 0
            total += m2o.value["cost"]
        return total

    def create_o2o(self, os2p, o2o, p2o_value):
        total = 0.0
        operation = Item('operation')
        for o in os2p:
            err = operation.get_w(o['operation_id'])
            if err:
                error(err)
                continue
            o2o.value["operation_id"] = o['operation_id']
            o2o.value["equipment_id"] = o['equipment_id']
            if p2o_value['width'] and operation.value['measure'] == 'мп.':
                o2o.value["number"] = p2o_value['pieces'] * (p2o_value['width'] + p2o_value['length'])*2/1000
            else:
                o2o.value["number"] = o['number'] * p2o_value['number']
            o2o.value["cost"] = o['cost'] * o2o.value['number']
            o2o.value["equipment_cost"] = o['equipment_cost'] * o2o.value['number']
            o2o.value["price"] = operation.value['price']
            o2o.value["user_id"] = o['user_id']
            o2o.value["user_sum"] = o2o.value['price'] * o2o.value['number']
            err = o2o.save()
            if err:
                error(err)
            o2o.value['id'] = 0
            total += o2o.value["cost"]
        return total

    def create_p2o(self, ps2p, p2o, p2o_value):
        total = 0.0
        product = Item('product')
        for p in ps2p:
            err = product.get_w(p['product2_id'])
            if err:
                error(err)
                continue
            p2o.value["product_id"] = p['product2_id']
            if p2o_value['width'] and product.value['measure'] == 'мп.':
                p2o.value["number"] = p2o_value['pieces'] * (p2o_value['width'] + p2o_value['length'])*2/1000
            else:
                p2o.value["number"] = p['number'] * p2o_value['number']
            p2o.value["cost"] = p['cost'] * p2o_value['number']
            p2o.value["price"] = p['cost']
            err = p2o.save()
            if err:
                error(err)
                continue
            
            total += self.make_product_value(p2o.value, product.value)
            p2o.value['id'] = 0
        return total

    def on_product_changed(self, p2o_value, product_value, num_changed, is_new):
        # if deleted
        if not p2o_value["is_active"]:
            self.update_sum()
            return 0
        if num_changed:
            self.clear_p2o_positions(p2o_value["id"])
            self.make_product_value(p2o_value, product_value)
        elif is_new:
            self.make_product_value(p2o_value, product_value)

        self.update_sum()
        
    def make_product_value(self, p2o_value, product_value):
        applied = {}
        if product_value:
            w = ProductView(product_value)
            if w.selects or w.multiselects:
                dlg = CustomDialog(w, 'Product')
                res = dlg.exec()
                if res:
                    applied = w.get_applied()
        
        total = 0.0
        m2o = Item("matherial_to_ordering")
        m2o.create_default()
        m2o.value["ordering_id"] = self.current_ordering_value['id']
        m2o.value["product_to_ordering_id"] = p2o_value['id']
        m2o.value["user_id"] = p2o_value['user_id']
        m2o.value["comm"] = "Auto"
        
        ms2p = Item('matherial_to_product')
        err = ms2p.get_filter_w('product_id', p2o_value['product_id'])
        if err:
            error(err)
            return
        default_ms2p = [v for v in ms2p.values if v['list_name']=='default']
        total += self.create_m2o(default_ms2p, m2o, p2o_value)
        if applied:
            total += self.create_m2o(applied['matherial_to_product'], m2o, p2o_value)
            
        o2o = Item("operation_to_ordering")
        o2o.create_default()
        o2o.value["ordering_id"] = self.current_ordering_value['id']
        o2o.value["product_to_ordering_id"] = p2o_value['id']
        o2o.value["comm"] = "Auto"
        
        os2p = Item('operation_to_product')
        err = os2p.get_filter_w('product_id', p2o_value['product_id'])
        if err:
            error(err)
            return
        default_os2p = [v for v in os2p.values if v['list_name']=='default']
        total += self.create_o2o(default_os2p, o2o, p2o_value)
        if applied:
            total += self.create_o2o(applied['operation_to_product'], o2o, p2o_value)
            
        p2o = Item("product_to_ordering")
        p2o.create_default()
        p2o.value["ordering_id"] = self.current_ordering_value['id']
        p2o.value["product_to_ordering_id"] = p2o_value['id']
        p2o.value["user_id"] = p2o_value['user_id']
        p2o.value["comm"] = "Auto"

        ps2p = Item('product_to_product')
        err = ps2p.get_filter_w('product_id', p2o_value['product_id'])
        if err:
            error(err)
            return
        default_ps2p = [v for v in ps2p.values if v['list_name']=='default']
        total += self.create_p2o(default_ps2p, p2o, p2o_value)
        if applied:
            total += self.create_p2o(applied['product_to_product'], p2o, p2o_value)
        
        p2o_value['price'] = round(total / p2o_value['number'], 2)
        p2o_value['cost'] = round(total + p2o_value['profit'] * p2o_value['number'], 2)
        p2o_item = Item('product_to_ordering')
        p2o_item.value = p2o_value
        err = p2o_item.save()
        if err:
            error(err)
        return p2o_value['cost']
        
       
class TreeItemToOrdering(QSplitter):
    orderingChanged = pyqtSignal()
    def __init__(self, show_info=True, is_info_bottom=False, buttons={}):
        super().__init__(Qt.Orientation.Vertical if is_info_bottom else Qt.Orientation.Horizontal)
        self.ordering_id = 0
        self.current_value = {}
        self.dataset = (
            ('matherial', 'Матеріал'),
            ('operation', 'Операція'),
            ('product', 'Виріб'),
        )
        self.tree_box = QSplitter(Qt.Orientation.Vertical)
        self.tree_box.setContentsMargins(0, 0, 0, 0)
        self.addWidget(self.tree_box)
        self.tree = Tree('product_to_ordering', fields=['name', 'cost'], headers=['Назва', 'Вартість'])
        
        self.tree.valueDoubleCklicked.connect(self.on_double_click)
        self.tree_box.addWidget(self.tree)
        if not show_info:
            return
        self.tree.itemSelected.connect(self.item_selected)
        self.info = PMOInfo(is_info_bottom)
        self.tree_box.addWidget(self.info)
        self.tree_box.setStretchFactor(0, 10)
        self.tree_box.setStretchFactor(1, 1)
        
        controls = QWidget()
        self.controls_box = QHBoxLayout()
        self.controls_box.setContentsMargins(0, 0, 0, 0)
        self.tree_box.insertWidget(0, controls)
        controls.setLayout(self.controls_box)
        self.controls_box.addStretch()
        del_btn = QPushButton('Видалити')
        self.controls_box.addWidget(del_btn)
        del_btn.clicked.connect(self.del_dialog)
        create_btn = QPushButton('Додати')
        self.controls_box.addWidget(create_btn)
        create_btn.clicked.connect(self.add_items)
        
        self.setStretchFactor(0, 10)
        self.setStretchFactor(1, 3)
        self.tree_box.setStretchFactor(0, 1)
        self.tree_box.setStretchFactor(1, 10)
        self.tree_box.setStretchFactor(2, 1)
        
    def add_items(self):
        creator = CalculatorTab()
        ordering = Item('ordering')
        err = ordering.get(self.ordering_id)
        if err:
            error(err)
            return
        creator.set_ordering(ordering)
        dlg = CustomDialog(creator, 'Додати позиції', 1250, 600)
        res = dlg.exec()
        if res:
            creator.make_it()
        self.reload(self.ordering_id)
        self.orderingChanged.emit()
    
    def reload(self, ordering_id=0):
        if ordering_id:
            self.ordering_id = ordering_id
        else:
            self.tree.reload()
            return
        values = []
        for data in self.dataset:
            res = self.get_item_values(data[0], data[1], self.ordering_id)
            if res:
                values += res
        self.tree.reload(values)        

    def get_item_values(self, name, hum, ordering_id):    
        values = []
        name_type = name + '_to_ordering'
        item = Item(name_type)
        err = item.get_filter_w('ordering_id', ordering_id)
        if err:
            error(err)
            return
        for v in item.values:
            if name == 'operation':
                cost = v["cost"]
            else:
                cost = v["cost"]
            values.append({
                "id": v["id"],
                "product_to_ordering_id": v['product_to_ordering_id'],
                "name": v['name'] if 'name' in v else v[name],
                "cost": round(cost, 2),
                "type": name_type,
                "type_hum": hum,
                "value": v,
            })
        return values
    
    def add_value(self, item):
        name = item.name.split('_to_')[0]
        v = {
                "id": item.value["id"],
                "product_to_ordering_id": item.value['product_to_ordering_id'],
                "name": item.value['name'] if 'name' in item.value else item.value[name],
                "cost": item.value["cost"],
                "type": item.name,
                "type_hum": item.hum,
                "value": item.value,
            }
        self.tree.add_value(v)
    
    def add_values(self, values):
        self.tree.add_values(values)

    def item_selected(self, value):
        self.current_value = value
        self.info.reload(value['value'], value['type'])
        
    def action(self, action, value):
        item = Item(value['type'])
        item.value = value['value']
        if action == 'delete':
            self.del_dialog(item)
            return
        self.dialog(item, 'Редагувати')
    
    def dialog(self, item, title):
        if item.name == 'matherial_to_ordering':
            form = MatherialToOrderingForm(value=item.value)
            form.widgets['matherial_id'].setDisabled(True)
        elif item.name == 'operation_to_ordering':
            form = OperationToOrderingForm(value=item.value)
            form.widgets['operation_id'].setDisabled(True)
        else:
            form = ProductToOrderingForm(value=item.value)
            form.widgets['product_id'].setDisabled(True)
        dlg = CustomFormDialog(title, form)

        res = dlg.exec()
        if res and dlg.value:
            item.value = self.prepare_value_to_save(dlg.value)
            err = item.save()
            if err:
                error(err)
                return
            self.reload(item.value['ordering_id'])
            self.orderingChanged.emit()
    
    def prepare_value_to_save(self, value):
        return value
    
    def del_dialog(self, value={}):
        if not value:
            if self.current_value:
                value = self.current_value
            else:
                error('Оберіть позицію!')
                return
        item = Item(value['type'])
        item.value = value['value']
        
        dlg = DeleteDialog(item.value)
        res = dlg.exec()
        if res:
            cause = dlg.entry.text()
            if cause:
                if 'comm' in item.value:
                    item.value['comm'] = f'del: {cause}' + item.value['comm']
                    item.save()
                elif 'info' in item.value:
                    item.value['info'] = f'Причина видалення:\n{cause}\n' + item.value['info']
                    item.save()
                
            err = item.delete(item.value['id'], cause)
            if err:
                error(err)
                return
            self.reload(self.ordering_id)
            self.orderingChanged.emit()

    def on_double_click(self, value):
        self.action('edit', value)

    def calc_sum(self):
        total = 0
        for child in self.tree.dataset[0]:
            total += child['cost']
        return total
    
    def deep_calc_sum(self, parent_id=0):
        total = 0
        if not len(self.tree.dataset):
            return 0
        for child in self.tree.dataset[parent_id]:
            if child['type'] == 'product_to_ordering' and child['id'] in self.tree.dataset:
                child_total = self.deep_calc_sum(child['id'])
                if child_total != child['cost']:
                    old_cost = child['cost']
                    child['cost'] = round(child_total*(1 + child['value']['persent']/100))
                    child['value']['cost'] = child['cost']
                    prod_item = Item('product_to_ordering')
                    prod_item.value = child['value']
                    err = prod_item.save()
                    if err: 
                        error(err)
                        total += old_cost
                        continue
            total += child['cost']
        return total
    
    
class ToOrdering(QSplitter):
    def __init__(self):
        super().__init__(Qt.Orientation.Horizontal)
        self.ordering = Item('ordering')
        fields = [
            "id",
            "name",
            "contragent",
            "user",
            "deadline_at",
        ]
        order_state = Item('ordering_status')
        err = order_state.get_all_w()
        if err:
            error(err)
            return
        
        self.buttons = ButtonsBlock('За статусом', order_state.values)
        self.ordering_table = OrderingTable(fields=fields)
        self.details_tree = TreeItemToOrdering(is_info_bottom=False)
        self.ordering_table.table.table.valueSelected.connect(self.ordering_selected)
        self.addWidget(self.buttons)
        self.addWidget(self.ordering_table)
        self.addWidget(self.details_tree)

        self.buttons.buttonClicked.connect(self.on_state_selected)
        self.details_tree.orderingChanged.connect(self.update_sum)
        self.ordering_table.actionResolved.connect(self.reload)
        
        self.doc_table = DocsTable()
        self.ordering_table.add_doc_table(self.doc_table)

    def reload(self, values=None):
        if values is None and self.buttons.current_value:
            err = self.ordering.get_filter_w('ordering_status_id', self.buttons.current_value['id'])
            if not err:
                self.ordering_table.reload(self.ordering.values)
                return
            else:
                error(err)
        self.ordering_table.reload(values)
        self.doc_table.reload([])
        self.details_tree.reload()

    def ordering_selected(self, value):
        self.details_tree.reload(value['id'])

    def on_state_selected(self, state_value):
        err = self.ordering.get_filter_w('ordering_status_id', state_value['id'])
        if err:
            error(err)
            return
        self.reload(self.ordering.values)

    def update_sum(self):
        self.details_tree.reload()
        new_sum = self.details_tree.deep_calc_sum()
        self.ordering.value = self.ordering_table.current_value
        self.ordering.value['price'] = new_sum
        self.ordering.value['profit'] = round(new_sum * self.ordering.value['persent']/100, 2)
        self.ordering.value['cost'] = round(new_sum + self.ordering.value['profit'], 2)
        err = self.ordering.save()
        if err:
            error(err)
            return
        rows = self.ordering_table.table.table.get_selected_rows()
        if not rows:
            return
        row = rows[0]
        self.ordering_table.reload()
        index = self.ordering_table.table.table._model.createIndex(row, 0)
        self.ordering_table.table.table.setCurrentIndex(index)
        