import math
from PyQt6.QtWidgets import (
    QPushButton,
    QWidget,
    QHBoxLayout,
    QVBoxLayout,
    )

from data.model import Item
from widgets.Dialogs import error, messbox
from common.params import TABLE_BUTTONS
from widgets.Form import (
    ItemTableWithDetails, 
    MainItemTable, 
    DetailsItemTable, 
    FormDialog, 
    CustomFormDialog,
    CustomForm,
    QSplitter,
    QTabWidget,
    )

from data_widgets.Documents import DocsTable
from data_widgets.Matherial import MatherialTab
from data.app import App

class MatherialToWhsInForm(CustomForm):
    def __init__(self, fields: list = [], value: dict = {}):
        self.item = Item('matherial_to_whs_in')
        super().__init__(self.item.model, fields, value)
        self.widgets['color_id'].setVisible(False)
        self.widgets['width'].setVisible(False)
        self.widgets['length'].setVisible(False)
        self.labels['color_id'].setVisible(False)
        self.labels['width'].setVisible(False)
        self.labels['length'].setVisible(False)
        self.widgets['matherial_id'].valChanged.connect(self.matherial_selected)

    def matherial_selected(self):
        matherial_value = self.widgets['matherial_id'].full_value()
        app = App()
        if matherial_value['measure'] == app.config['measure linear']:
            self.labels['length'].setVisible(True)
            self.widgets['length'].setVisible(True)
        if matherial_value['measure'] == app.config['measure square']:
            self.labels['length'].setVisible(True)
            self.labels['width'].setVisible(True)
            self.widgets['width'].setVisible(True)
            self.widgets['length'].setVisible(True)
        if matherial_value['color_group_id']:
            self.labels['color_id'].setVisible(True)
            self.widgets['color_id'].setVisible(True)
            self.widgets['color_id'].group_id = matherial_value['color_group_id']
            

class DetailsMatherialToWhsInTable(DetailsItemTable):
    def __init__(self, fields: list = [], values: list = None, buttons=TABLE_BUTTONS, group_id=0):
        super().__init__('matherial_to_whs_in', '', fields, values, buttons, group_id)
        bottom_hbox = QWidget()
        self.bottom_controls = QHBoxLayout()
        self.bottom_controls.setContentsMargins(0, 0, 0, 0)
        bottom_hbox.setLayout(self.bottom_controls)
        reprice_btn = QPushButton('Переоцінка')
        self.bottom_controls.addStretch()
        self.bottom_controls.addWidget(reprice_btn)
        self.addWidget(bottom_hbox)
        reprice_btn.clicked.connect(self.reprice)

    def reprice(self):
        values = self.table.table.get_selected_values()
        matherial = Item('matherial')
        mess = ''
        for v in values:
            err = matherial.get(v['matherial_id'])
            if err:
                error(err)
                continue
            new_price = round(v['price'], 2)
            add_mess = f"{matherial.value['name']} з {matherial.value['price']} грн. на {new_price} грн.\n"
            matherial.value['price'] = new_price
            err = matherial.save()
            if err:
                error(err)
                continue
            mess += add_mess
        messbox(mess[:-1], 'Переоцінка')
    
    def dialog(self, value, title):
        i = Item(self.item.name)
        form = MatherialToWhsInForm(fields=i.columns, value=value)
        dlg = CustomFormDialog(title, form)
        res = dlg.exec()
        if res and dlg.value:
            i.value = self.prepare_value_to_save(dlg.value)
            err = i.save()
            if err:
                error(err)
                return
            self.reload()
            self.actionResolved.emit()
        
    def prepare_value_to_save(self, value):
        value['cost'] = value['number'] * value['price']
        return value


class MatherialToWhsOutForm(CustomForm):
    def __init__(self, fields: list = [], value: dict = {}):
        self.item = Item('matherial_to_whs_out')
        super().__init__(self.item.model, fields, value)
        self.widgets['color_id'].setVisible(False)
        self.widgets['width'].setVisible(False)
        self.widgets['length'].setVisible(False)
        self.labels['color_id'].setVisible(False)
        self.labels['width'].setVisible(False)
        self.labels['length'].setVisible(False)
        self.widgets['matherial_id'].valChanged.connect(self.matherial_selected)

    def matherial_selected(self):
        matherial_value = self.widgets['matherial_id'].full_value()
        self.widgets['price'].setValue(matherial_value['price'])
        app = App()
        if matherial_value['measure'] == app.config['measure linear']:
            self.widgets['length'].setVisible(True)
            self.labels['length'].setVisible(True)
        if matherial_value['measure'] == app.config['measure square']:
            self.widgets['width'].setVisible(True)
            self.widgets['length'].setVisible(True)
            self.labels['width'].setVisible(True)
            self.labels['length'].setVisible(True)
        if matherial_value['color_group_id']:
            self.widgets['color_id'].setVisible(True)
            self.labels['color_id'].setVisible(True)
            self.widgets['color_id'].group_id = matherial_value['color_group_id']


class DetailsMatherialToWhsOutTable(DetailsItemTable):
    def __init__(self, fields: list = [], values: list = None, buttons=TABLE_BUTTONS, group_id=0):
        super().__init__('matherial_to_whs_out', '', fields, values, buttons, group_id)

    def dialog(self, value, title):
        i = Item(self.item.name)
        form = MatherialToWhsOutForm(value=value, fields=i.columns)
        dlg = CustomFormDialog(title, form)
        res = dlg.exec()
        if res and dlg.value:
            i.value = self.prepare_value_to_save(dlg.value)
            err = i.save()
            if err:
                error(err)
                return
            self.reload()
            self.actionResolved.emit()
        
    def prepare_value_to_save(self, value):
        value['cost'] = value['number'] * value['price']
        return value


class WhsItemTableWithDetails(QSplitter):
    def __init__(
            self, 
            main_table: MainItemTable,
            details_table: DetailsItemTable,
            matherial_table: MatherialTab,
            ):
        super().__init__()
        self.table = main_table
        self.details = details_table
        self.table.set_detais_table(self.details)
        self.details.set_main_table(self.table)
        tabs = QTabWidget()
        tabs.addTab(self.table, "Накладна")
        self.matherials = matherial_table
        tabs.addTab(self.matherials, "Обрати матеріали")
        self.addWidget(tabs)
        self.addWidget(self.details)
        self.table.table.table.valueSelected.connect(self.reload_details)
        
    def reload_details(self, value):
        self.details.item.get_filter_w(self.table.item.name + '_id', value['id'])
        self.details.reload(self.details.item.values)

    def reload(self, values=None):
        self.table.reload(values)
        if values:
            self.reload_details(values[0])


class WhsTab(WhsItemTableWithDetails):
    def __init__(self, main_table, details_table):
        self.main_table = main_table
        self.details_table = details_table
        matherials = MatherialTab()
        super().__init__(self.main_table, self.details_table, matherials)
        self.details_table.actionResolved.connect(self.update_sum)
        self.setStretchFactor(0, 2)
        self.setStretchFactor(1, 3)
    
    def reload(self, values=None):
        self.main_table.reload(values)
        if values:
            self.main_table.info.reload(values[0])
            self.reload_details(values[0])

    def calc_sum(self):
        return self.details_table.calc_sum('cost')
    
    def update_sum(self):
        rows = self.main_table.table.table.get_selected_rows()
        if not rows:
            return
        row = rows[0]
        self.main_table.reload()
        index = self.main_table.table.table._model.createIndex(row, 0)
        self.main_table.table.table.setCurrentIndex(index)


class WhsInTab(WhsTab):
    def __init__(self):
        main_table = MainItemTable(
            'whs_in', 
            fields=['id', 'name', 'whs_sum', 'contragent', 'whs'], 
            is_vertical_inner=False, 
            releazed_buttons=True,
            )
        details_table = DetailsMatherialToWhsInTable(
            fields=['id', 'matherial', 'number', 'price', 'cost', 'color'],
            )
        super().__init__(main_table, details_table)

        self.doc_table = DocsTable('whs_in')
        self.main_table.add_doc_table(self.doc_table)

        cash_out_btn = QPushButton('Створити ВКО')
        cash_in_btn = QPushButton('Створити ПКО')
        whs_out_btn = QPushButton('Створити ВН')
        cash_out_delivery_btn = QPushButton('ВКО на доставку')
        delivery_btn = QPushButton('Поділити доставку')
        row = self.main_table.info.grid.rowCount()
        col = self.main_table.info.grid.columnCount()
        self.main_table.info.grid.addWidget(cash_out_btn, row-4, col-1)
        self.main_table.info.grid.addWidget(cash_in_btn, row-3, col-1)
        self.main_table.info.grid.addWidget(whs_out_btn, row-3, col-2)
        self.main_table.info.grid.addWidget(cash_out_delivery_btn, row-2, col-1)
        self.main_table.info.grid.addWidget(delivery_btn, row-2, col-2)
        cash_out_btn.clicked.connect(self.create_cash_out)
        cash_in_btn.clicked.connect(self.create_cash_in)
        whs_out_btn.clicked.connect(self.create_whs_out)
        cash_out_delivery_btn.clicked.connect(self.create_cash_out_delivery)
        delivery_btn.clicked.connect(self.distrib_delivery)
        self.matherials.remove_dblclick_cb()
        self.matherials.set_dblclick_cb(self.add_matherial)

    def add_matherial(self, value):
        cur_value = self.main_table.table.table.get_selected_value()
        if not cur_value:
            error("Оберіть накладну!")
            return
        i = Item('matherial_to_whs_in')
        i.create_default_w()
        i.value['matherial_id'] = value['id']
        i.value['price'] = value['price']
        i.value['whs_in_id'] = cur_value['id']
        self.details.dialog(i.value, "Додати матеріал")

    def calc_sum(self):
        cur_value = self.main_table.table.table.get_selected_value()
        if not cur_value:
            return
        return super().calc_sum() + cur_value["delivery"]

    def create_cash_out(self):
        cur_value = self.main_table.table.table.get_selected_value()
        if not cur_value:
            return
        
        payed = self.doc_table.calc_by_type('cash_out')
        cash_out = Item('cash_out')
        cash_out.create_default()

        if cur_value['based_on']:
            cash_out.value["based_on"] = cur_value['based_on']    
        else:
            cash_out.value["based_on"] = f"whs_in.{cur_value['id']}"

        cash_out.value["contragent_id"] = cur_value["contragent_id"]
        cash_out.value["contact_id"] = cur_value["contact_id"]
        cash_out.value["cash_sum"] = cur_value["whs_sum"] - payed
        cash_out.value["comm"] = f'авт. до {cur_value["name"]}'
        app = App()
        cash_out.value["cash_id"] = app.config["whs_in cash id"]
        cash_out.value["user_id"] = app.user['id']

        dlg = FormDialog('Створити ВКО', cash_out.model, cash_out.columns, cash_out.value)
        res = dlg.exec()
        if res:
            err = cash_out.save()
            if err:
                error(err)
                return
            self.doc_table.reload(cur_value)

    def create_cash_in(self):
        cur_value = self.main_table.table.table.get_selected_value()
        if not cur_value:
            return
        cash_in = Item('cash_in')
        cash_in.create_default()
        cash_in.value["based_on"] = f"whs_in.{cur_value['id']}"
        cash_in.value["contragent_id"] = cur_value["contragent_id"]
        cash_in.value["contact_id"] = cur_value["contact_id"]
        cash_in.value["comm"] = f'авт. до {cur_value["name"]}'
        app = App()
        cash_in.value["cash_id"] = app.config["whs_in cash id"]
        cash_in.value["user_id"] = app.user['id']

        m2wi_values = self.details_table.table.table.get_selected_values()
        for v in m2wi_values:
            cash_in.value["cash_sum"] += v['cost']

        dlg = FormDialog('Створити ПКО', cash_in.model, cash_in.columns, cash_in.value)
        res = dlg.exec()
        if res:
            err = cash_in.save()
            if err:
                error(err)
                return
            self.doc_table.reload(cur_value)

    def create_whs_out(self):
        cur_value = self.main_table.table.table.get_selected_value()
        if not cur_value:
            return
        whs_out = Item('whs_out')
        whs_out.create_default()
        if cur_value['based_on']:
            whs_out.value["based_on"] = cur_value['based_on']    
        else:
            whs_out.value["based_on"] = f"whs_in.{cur_value['id']}"
        whs_out.value["contragent_id"] = cur_value["contragent_id"]
        whs_out.value["contact_id"] = cur_value["contact_id"]
        whs_out.value["comm"] = f'авт. до {cur_value["name"]}'
        app = App()
        whs_out.value["user_id"] = app.user['id']

        m2wi_values = self.details_table.table.table.get_selected_values()
        if not m2wi_values:
            m2wi_values = self.details_table.values()
        for v in m2wi_values:
            whs_out.value["whs_sum"] += v['cost']
        
        dlg = FormDialog(
            'Створити ВН', 
            whs_out.model, 
            whs_out.columns, 
            whs_out.value,
            )
        res = dlg.exec()
        if not res:
            return
        whs_out.value["whs_sum"] = 0
        err = whs_out.save()
        if err:
            error(err)
            return
        for v in m2wi_values:
            m2wo = Item('matherial_to_whs_out')
            m2wo.create_default()
            m2wo.value["whs_out_id"] = whs_out.value['id']
            m2wo.value["matherial_id"] = v["matherial_id"]
            m2wo.value["number"] = v["number"]
            m2wo.value["price"] = v["price"]
            m2wo.value["cost"] = v["cost"]
            m2wo.value["width"] = v["width"]
            m2wo.value["length"] = v["length"]
            m2wo.value["color_id"] = v["color_id"]
            err = m2wo.save()
            if err:
                error(err)
            
        self.doc_table.reload(cur_value)

    def create_cash_out_delivery(self):
        cur_value = self.main_table.table.table.get_selected_value()
        if not cur_value:
            return
        app = App()
        cash_out = Item('cash_out')
        cash_out.create_default()
        cash_out.value["based_on"] = f"whs_in.{cur_value['id']}"
        cash_out.value["contragent_id"] = app.config["contragent for delivery"]
        cash_out.value["contact_id"] = app.config["contact for delivery"]
        cash_out.value["cash_sum"] = cur_value["delivery"]
        cash_out.value["comm"] = f'авт. до {cur_value["name"]}'
        cash_out.value["cash_id"] = app.config["whs_in cash id"]
        cash_out.value["user_id"] = app.user['id']
        err = cash_out.save()
        if err:
            error(err)
            return
        self.doc_table.reload(cur_value)
    
    def correct_add_sum(self, value, add_sum):
        sum = value['cost'] + add_sum
        price = math.ceil(sum/value['number'] * 100) / 100
        sum = round(price * value['number'], 2)
        add_sum = sum - value['cost']
        return add_sum, sum, price

    def distrib_delivery(self):
        cur_value = self.main_table.table.table.get_selected_value()
        if not cur_value:
            return
        k = cur_value['delivery']/cur_value['whs_sum']
        delivery_sum = cur_value['delivery']
        m2wi_values = self.details_table.table.table.get_selected_values()
        if not m2wi_values:
            m2wi_values = self.details_table.values()
        else:
            ksum = 0
            for v in m2wi_values:
                ksum += v['cost']
            k = cur_value['delivery']/ksum
        m2wi = Item('matherial_to_whs_in')
        for v in m2wi_values:
            add_sum = round(v['cost'] * k)
            add_sum, sum, price = self.correct_add_sum(v, add_sum)
            
            if delivery_sum - add_sum < 0:
                add_sum = delivery_sum
                add_sum, sum, price = self.correct_add_sum(v, add_sum)
            
            delivery_sum -= add_sum
            
            v['price'] = price
            v['cost'] = sum
            m2wi.value = v
            err = m2wi.save()
            if err:
                error(err)
        whs_in = Item('whs_in')
        cur_value['whs_sum'] += cur_value['delivery']
        cur_value['delivery'] = 0
        whs_in.value = cur_value
        err = whs_in.save()
        if err:
            error(err)
        self.reload()
                    

class WhsOutTab(WhsTab):
    def __init__(self):
        main_table = MainItemTable(
            'whs_out', 
            fields=['id', 'name', 'whs_sum', 'contragent', 'whs'], 
            releazed_buttons=True,
            )
        details_table = DetailsMatherialToWhsOutTable(fields=['id', 'matherial', 'number', 'price', 'cost', 'color'])
        super().__init__(main_table, details_table)
        self.matherials.remove_dblclick_cb()
        self.matherials.set_dblclick_cb(self.add_matherial)

    def add_matherial(self, value):
        cur_value = self.main_table.table.table.get_selected_value()
        if not cur_value:
            error("Оберіть накладну!")
            return
        i = Item('matherial_to_whs_out')
        i.create_default_w()
        i.value['matherial_id'] = value['id']
        i.value['price'] = value['price']
        i.value['whs_out_id'] = cur_value['id']
        self.details.dialog(i.value, "Додати матеріал")
        

class WhsesTab(QWidget):
    def __init__(self):
        super().__init__()
        self.box = QVBoxLayout()
        self.setLayout(self.box)
        self.whs_table = MainItemTable('whs', is_info_bottom=True, is_vertical_inner=True)
        docs = DocsTable(
                'whs',
                # main_key='id', 
                # doc_key='whs_id', 
                docs=('whs_in', 'whs_out'),
                controls=True,
                )
        self.whs_table.add_doc_table(docs, True)
        self.whs_table.setStretchFactor(0, 1)
        self.whs_table.setStretchFactor(2, 10)
        self.mat_table = DetailsItemTable('wmc_number')
        self.main = ItemTableWithDetails(self.whs_table, self.mat_table)
        self.box.addWidget(self.main)
        self.whs_table.show_info(False)

    def reload(self):
        self.main.reload()
        