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
    )

from data_widgets.Documents import DocsTable

# class MatherialToWhsInFormDialog(FormDialog):
#     def __init__(self, title, data_model, value):
#         super().__init__(title, data_model, value)

class MatherialToWhsInForm(CustomForm):
    def __init__(self, fields: list = [], value: dict = {}):
        self.item = Item('matherial_to_whs_in')
        super().__init__(self.item.model, fields, value)
        self.widgets['color_id'].setVisible(False) #.setDisabled(True)
        self.widgets['width'].setVisible(False) #.setDisabled(True)
        self.widgets['length'].setVisible(False) #.setDisabled(True)
        self.labels['color_id'].setVisible(False)
        self.labels['width'].setVisible(False)
        self.labels['length'].setVisible(False)
        self.widgets['matherial_id'].valChanged.connect(self.matherial_selected)

    def matherial_selected(self):
        matherial_value = self.widgets['matherial_id'].full_value()
        if matherial_value['measure'] == 'мп.':
            self.labels['length'].setVisible(True)
            self.widgets['length'].setVisible(True) #.setDisabled(False)
        if matherial_value['measure'] == 'м2':
            self.labels['length'].setVisible(True)
            self.labels['width'].setVisible(True)
            self.widgets['width'].setVisible(True) #.setDisabled(False)
            self.widgets['length'].setVisible(True) #.setDisabled(False)
        if matherial_value['color_group_id']:
            self.labels['color_id'].setVisible(True)
            self.widgets['color_id'].setVisible(True) #.setDisabled(False)
            self.widgets['color_id'].group_id = matherial_value['color_group_id']
            


class DetailsMatherialToWhsInTable(DetailsItemTable):
    def __init__(self, fields: list = [], values: list = None, buttons=TABLE_BUTTONS, group_id=0):
        super().__init__('matherial_to_whs_in', 'name', fields, values, buttons, group_id)
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
            add_mess = f"{matherial.value['name']} з {matherial.value['price']} грн. на {v['price']} грн.\n"
            matherial.value['price'] = v['price']
            err = matherial.save()
            if err:
                error(err)
                continue
            mess += add_mess
        messbox(mess[:-1], 'Переоцінка')
    
    def dialog(self, value, title):
        i = Item(self.item.name)
        # print(self.item.name, '====>', i.model.keys())
        form = MatherialToWhsInForm(value=value)
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
        self.widgets['color_id'].setVisible(False) #.setDisabled(True)
        self.widgets['width'].setVisible(False) #.setDisabled(True)
        self.widgets['length'].setVisible(False) #.setDisabled(True)
        self.widgets['matherial_id'].valChanged.connect(self.matherial_selected)

    def matherial_selected(self):
        matherial_value = self.widgets['matherial_id'].full_value()
        self.widgets['price'].setValue(matherial_value['price'])
        if matherial_value['measure'] == 'мп.':
            self.widgets['length'].setVisible(True) #.setDisabled(False)
        if matherial_value['measure'] == 'м2':
            self.widgets['width'].setVisible(True) #.setDisabled(False)
            self.widgets['length'].setVisible(True) #.setDisabled(False)
        if matherial_value['color_group_id']:
            self.widgets['color_id'].setVisible(True) #.setDisabled(False)
            self.widgets['color_id'].group_id = matherial_value['color_group_id']



class DetailsMatherialToWhsOutTable(DetailsItemTable):
    def __init__(self, fields: list = [], values: list = None, buttons=TABLE_BUTTONS, group_id=0):
        super().__init__('matherial_to_whs_out', 'name', fields, values, buttons, group_id)

    def dialog(self, value, title):
        i = Item(self.item.name)
        form = MatherialToWhsOutForm(value=value)
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


class WhsTab(ItemTableWithDetails):
    def __init__(self, name, main_table, details_table):
        self.main_table = main_table
        self.details_table = details_table
        super().__init__(self.main_table, self.details_table)
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
        main_table = MainItemTable('whs_in', fields=['id', 'name', 'whs_sum', 'contragent', 'whs'], is_vertical_inner=False)
        details_table = DetailsMatherialToWhsInTable(
            fields=['id', 'matherial', 'number', 'price', 'cost', 'color'],
            )
        super().__init__('whs_in', main_table, details_table)

        self.doc_table = DocsTable()
        self.main_table.add_doc_table(self.doc_table)

        cash_out_btn = QPushButton('Створити ВКО')
        cash_in_btn = QPushButton('Створити ПКО')
        whs_out_btn = QPushButton('Створити ВН')
        row = self.main_table.info.grid.rowCount()
        col = self.main_table.info.grid.columnCount()
        self.main_table.info.grid.addWidget(cash_out_btn, row-4, col-1)
        self.main_table.info.grid.addWidget(cash_in_btn, row-3, col-1)
        self.main_table.info.grid.addWidget(whs_out_btn, row-3, col-2)
        cash_out_btn.clicked.connect(self.create_cash_out)
        cash_in_btn.clicked.connect(self.create_cash_in)
        whs_out_btn.clicked.connect(self.create_whs_out)

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
        cash_out.value["based_on"] = cur_value['document_uid']
        cash_out.value["contragent_id"] = cur_value["contragent_id"]
        cash_out.value["contact_id"] = cur_value["contact_id"]
        cash_out.value["cash_sum"] = cur_value["whs_sum"] + cur_value["delivery"] - payed
        cash_out.value["comm"] = f'авт. до {cur_value["name"]}'

        dlg = FormDialog('Створити ВКО', cash_out.model, cash_out.value)
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
        cash_in.value["based_on"] = cur_value['document_uid']
        cash_in.value["contragent_id"] = cur_value["contragent_id"]
        cash_in.value["contact_id"] = cur_value["contact_id"]
        cash_in.value["comm"] = f'авт. до {cur_value["name"]}'

        m2wi_values = self.details_table.table.table.get_selected_values()
        for v in m2wi_values:
            cash_in.value["cash_sum"] += v['cost']

        dlg = FormDialog('Створити ПКО', cash_in.model, cash_in.value)
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
        whs_out.value["based_on"] = cur_value['document_uid']
        whs_out.value["contragent_id"] = cur_value["contragent_id"]
        whs_out.value["contact_id"] = cur_value["contact_id"]
        whs_out.value["comm"] = f'авт. до {cur_value["name"]}'

        m2wi_values = self.details_table.table.table.get_selected_values()
        for v in m2wi_values:
            whs_out.value["whs_sum"] += v['cost']
        
        dlg = FormDialog('Створити ВН', whs_out.model, whs_out.value)
        res = dlg.exec()
        if not res:
            return
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

    # def reload_details(self, value):
    #     super().reload_details(value)
    #     self.doc_table.reload_docs(value)
                
        
    # def calc_sum(self):
    #     cur_value = self.main_table.table.table.get_selected_value()
    #     new_sum = self.details_table.calc_sum('cost')
    #     return new_sum + cur_value['delivery']
        


class WhsOutTab(WhsTab):
    def __init__(self):
        main_table = MainItemTable('whs_out', fields=['id', 'name', 'whs_sum', 'contragent', 'whs'])
        details_table = DetailsMatherialToWhsOutTable(fields=['id', 'matherial', 'number', 'price', 'cost', 'color'])
        super().__init__('whs_out', main_table, details_table)
        

class WhsesTab(QWidget):
    def __init__(self):
        super().__init__()
        self.box = QVBoxLayout()
        self.setLayout(self.box)
        self.whs_table = MainItemTable('whs', is_info_bottom=True, is_vertical_inner=True)
        self.whs_table.add_doc_table(DocsTable(main_key='id', doc_key='whs_id', docs=('whs_in', 'whs_out')), True)
        self.mat_table = DetailsItemTable('wmc_number')
        
        self.main = ItemTableWithDetails(self.whs_table, self.mat_table)
        self.box.addWidget(self.main)

    def reload(self):
        self.main.reload()
        