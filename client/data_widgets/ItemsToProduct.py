from PyQt6.QtWidgets import (
    QPushButton,
    QWidget,
    QVBoxLayout,
    QHBoxLayout,
    QSplitter,
    QTabWidget,
    QStyle,
    )

from data.model import Item
from widgets.Dialogs import error, askdlg, CustomDialog, messbox, ok_cansel_dlg
from common.params import TABLE_BUTTONS
from widgets.Form import (
    MainItemTable, 
    DetailsItemTable, 
    CustomFormDialog,
    CustomForm,
    SelectDialog,
    )
from data_widgets.Calculation import ProductWidget


class ProductForm(CustomForm):
    def __init__(self, fields: list = [], value: dict = {}):
        self.item = Item('product')
        super().__init__(self.item.model, fields, value, 2)
        

class ProductTable(MainItemTable):
    def __init__(self, fields: list = [], values: list = None, buttons=TABLE_BUTTONS, group_id=0):
        super().__init__('product', 'name', fields, values, buttons, group_id, is_info_bottom=True)
        test_btn = QPushButton('Тест')
        self.table.hbox.addWidget(test_btn)
        test_btn.clicked.connect(self.calc_dialog)

    def calc_dialog(self):
        value = self.table.table.get_selected_value()
        if value is None:
            return
        calc = ProductWidget(value)
        calc.form.hide_save_btn()
        for name in ('name', 'profit', 'user_id', 'info'):
            if name in calc.form.labels:
                calc.form.labels[name].setVisible(False)
            if name in calc.form.widgets:
                calc.form.widgets[name].setVisible(False)
        dlg = CustomDialog(calc, 'Тест вироба')
        dlg.exec()


    def dialog(self, value, title):
        i = Item(self.item.name)
        # print(self.item.name, '====>', i.model.keys())
        form = ProductForm(value=value)
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


class ProductToProductForm(CustomForm):
    def __init__(self, fields: list = [], value: dict = {}):
        self.item = Item('product_to_product')
        self.price = 0
        super().__init__(self.item.model, fields, value)
        product =  self.widgets['product2_id'].full_value()
        if product:
            self.price = product['cost']
        self.widgets['product2_id'].valChanged.connect(self.product_selected)
        self.widgets['cost'].valChanged.connect(self.cost_changed)
        self.widgets['number'].valChanged.connect(self.number_changed)
        # self.widgets['coeff'].valChanged.connect(self.number_changed)

        self.labels['coeff'].setVisible(False)
        self.labels['is_multiselect'].setVisible(False)
        self.widgets['coeff'].setVisible(False)
        self.widgets['is_multiselect'].setVisible(False)
        
        if value['list_name'] == 'default':
            self.labels['is_used'].setVisible(False)
            self.labels['list_name'].setVisible(False)
            self.widgets['is_used'].setVisible(False)
            self.widgets['list_name'].setVisible(False)
        else:    
            self.widgets['list_name'].setDisabled(True)
        

    def cost_changed(self):
        number = self.widgets['number'].value()
        coeff = (self.widgets['cost'].value() / number) / self.price
        self.widgets['coeff'].set_value(coeff)

    def number_changed(self):
        number = self.widgets['number'].value()
        if not number:
            # error("Кількість має бути більшою за 0!")
            return
        coeff = self.widgets['coeff'].value()
        self.widgets['cost'].set_value(round(self.price * number * coeff, 2))

    def product_selected(self):
        product_value = self.widgets['product2_id'].full_value()
        self.price = product_value['cost']
        number = self.widgets['number'].value()
        coeff = self.widgets['coeff'].value()
        self.widgets['cost'].set_value(round(self.price * number * coeff, 2))



class DetailsProductToProductTable(DetailsItemTable):
    def __init__(self, fields: list = [], values: list = None, list_name='default', is_multiselect=False):
        self.list_name = list_name
        self.is_multiselect = is_multiselect
        super().__init__('product_to_product', '', fields, values)
        btn = QPushButton()
        pixmapi = QStyle.StandardPixmap.SP_FileDialogListView
        icon = self.style().standardIcon(pixmapi)
        btn.setIcon(icon)
        btn.setToolTip('Додати декілька')
        btn.clicked.connect(lambda _,action='add_many': self.action(action))
        self.table.hbox.addWidget(btn)
        
    def action(self, action: str, value: dict=None):
        if action == 'add_many':
            self.add_many_dialog(value)
            return
        return super().action(action, value)
    
    def prepare_value_to_action(self):
        value = super().prepare_value_to_action()
        value['list_name'] = self.list_name
        value['is_multiselect'] = self.is_multiselect
        return value
    
    def reload(self, values=None):
        if values is None:
            item = Item('product_to_product')
            product_value = self.main_table.current_value
            err = item.get_filter_w('product_id', product_value['id'])
            if err:
                error(err)
                return
            values = [v for v in item.values if v['list_name'] == self.list_name]
        return super().reload(values)
            
    def prepare_value_to_save(self, value):
        return value
    
    def dialog(self, value, title):
        i = Item(self.item.name)
        # print(self.item.name, '====>', i.model.keys())
        form = ProductToProductForm(value=value)
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

    def add_many_dialog(self, value):
        i = Item(self.item.name)
        dlg = SelectDialog('product')
        res = dlg.exec()
        if res:
            values = dlg.widget.table.table.get_selected_values()
            for v in values:
                i.value = self.prepare_value_to_action()
                i.value['product2_id'] = v['id']
                i.value['cost'] = v['cost']
                err = i.save()
                if err:
                    error(err)
                    continue
            self.reload()
            self.actionResolved.emit()



class MatherialToProductForm(CustomForm):
    def __init__(self, fields: list = [], value: dict = {}):
        self.item = Item('matherial_to_product')
        self.price = 0
        # if value:
        #     if value['matherial_id']:
        #         m = Item('matherial')
        #         err = m.get(value['matherial_id'])
        #         if err:
        #             error(err)
        #             return
        #         self.price = m.value['cost']
        super().__init__(self.item.model, fields, value)
        mat =  self.widgets['matherial_id'].full_value()
        if mat:
            self.price = mat['cost']
        self.widgets['matherial_id'].valChanged.connect(self.matherial_selected)
        self.widgets['cost'].valChanged.connect(self.cost_changed)
        self.widgets['number'].valChanged.connect(self.number_changed)
        # self.widgets['coeff'].valChanged.connect(self.number_changed)
        
        self.labels['coeff'].setVisible(False)
        self.labels['is_multiselect'].setVisible(False)
        self.widgets['coeff'].setVisible(False)
        self.widgets['is_multiselect'].setVisible(False)
        
        if value['list_name'] == 'default':
            self.labels['is_used'].setVisible(False)
            self.labels['list_name'].setVisible(False)
            self.widgets['is_used'].setVisible(False)
            self.widgets['list_name'].setVisible(False)
        else:    
            self.widgets['list_name'].setDisabled(True)
        

    def cost_changed(self):
        number = self.widgets['number'].value()
        coeff = (self.widgets['cost'].value() / number) / self.price
        self.widgets['coeff'].set_value(coeff)

    def number_changed(self):
        number = self.widgets['number'].value()
        if not number:
            # error("Кількість має бути більшою за 0!")
            return
        coeff = self.widgets['coeff'].value()
        self.widgets['cost'].set_value(round(self.price * number * coeff, 2))

    def matherial_selected(self):
        matherial_value = self.widgets['matherial_id'].full_value()
        self.price = matherial_value['cost']
        number = self.widgets['number'].value()
        coeff = self.widgets['coeff'].value()
        self.widgets['cost'].set_value(round(self.price * number * coeff, 2))


class DetailsMatherialToProductTable(DetailsItemTable):
    def __init__(self, fields: list = [], values: list = None, list_name='default', is_multiselect=False):
        self.list_name = list_name
        self.is_multiselect = is_multiselect
        super().__init__('matherial_to_product', '', fields, values)
        btn = QPushButton()
        pixmapi = QStyle.StandardPixmap.SP_FileDialogListView
        icon = self.style().standardIcon(pixmapi)
        btn.setIcon(icon)
        btn.setToolTip('Додати декілька')
        btn.clicked.connect(lambda _,action='add_many': self.action(action))
        self.table.hbox.addWidget(btn)
        
    def action(self, action: str, value: dict=None):
        if action == 'add_many':
            self.add_many_dialog(value)
            return
        return super().action(action, value)
    
    def prepare_value_to_action(self):
        value = super().prepare_value_to_action()
        value['list_name'] = self.list_name
        value['is_multiselect'] = self.is_multiselect
        return value

    def prepare_value_to_save(self, value):
        return value
    
    def reload(self, values=None):
        if values is None:
            item = Item('matherial_to_product')
            product_value = self.main_table.current_value
            err = item.get_filter_w('product_id', product_value['id'])
            if err:
                error(err)
                return
            values = [v for v in item.values if v['list_name'] == self.list_name]
        return super().reload(values)
    
    def dialog(self, value, title):
        i = Item(self.item.name)
        form = MatherialToProductForm(value=value)
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

    def add_many_dialog(self, value):
        i = Item(self.item.name)
        dlg = SelectDialog('matherial')
        res = dlg.exec()
        if res:
            values = dlg.widget.table.table.get_selected_values()
            for v in values:
                i.value = self.prepare_value_to_action()
                i.value['matherial_id'] = v['id']
                i.value['cost'] = v['cost']
                err = i.save()
                if err:
                    error(err)
                    continue
            self.reload()
            self.actionResolved.emit()




class OperationToProductForm(CustomForm):
    def __init__(self, fields: list = [], value: dict = {}):
        self.item = Item('operation_to_product')
        self.price = 0
        self.eq_price = 0
        super().__init__(self.item.model, fields, value)
        operation_value =  self.widgets['operation_id'].full_value()
        if operation_value:
            self.price = operation_value['cost']
            self.eq_price = operation_value['equipment_price']
        
        self.widgets['operation_id'].valChanged.connect(self.operation_selected)
        self.widgets['cost'].valChanged.connect(self.cost_changed)
        self.widgets['equipment_cost'].valChanged.connect(self.cost_changed)
        self.widgets['number'].valChanged.connect(self.number_changed)
        # self.widgets['coeff'].valChanged.connect(self.number_changed)

        self.labels['coeff'].setVisible(False)
        self.labels['is_multiselect'].setVisible(False)
        self.widgets['coeff'].setVisible(False)
        self.widgets['is_multiselect'].setVisible(False)
        
        if value['list_name'] == 'default':
            self.labels['is_used'].setVisible(False)
            self.labels['list_name'].setVisible(False)
            self.widgets['is_used'].setVisible(False)
            self.widgets['list_name'].setVisible(False)
        else:    
            self.widgets['list_name'].setDisabled(True)
        

    def operation_selected(self):
        operation_value = self.widgets['operation_id'].full_value()
        self.price = operation_value['cost']
        self.eq_price = operation_value['equipment_price']
        number = self.widgets['number'].value()
        coeff = self.widgets['coeff'].value()
        self.widgets['cost'].set_value(round(self.price * number * coeff, 2))
        self.widgets['equipment_cost'].set_value(round(self.eq_price * number, 2))
        self.widgets['equipment_id'].set_value(operation_value['equipment_id'])

    def cost_changed(self):
        number = self.widgets['number'].value()
        cost = self.widgets['cost'].value() 
        coeff = (cost / number) / self.price if self.price else 1
        self.widgets['coeff'].set_value(coeff)

    def number_changed(self):
        number = self.widgets['number'].value()
        if not number:
            # error("Кількість має бути більшою за 0!")
            return
        coeff = self.widgets['coeff'].value()
        self.widgets['cost'].set_value(round(self.price * number * coeff, 2))
        self.widgets['equipment_cost'].set_value(round(self.eq_price * number, 2))

    
    
class DetailsOperationToProductTable(DetailsItemTable):
    def __init__(self, fields: list = [], values: list = None, list_name='default', is_multiselect=False):
        self.list_name = list_name
        self.is_multiselect = is_multiselect
        super().__init__('operation_to_product', '', fields, values)
        btn = QPushButton()
        pixmapi = QStyle.StandardPixmap.SP_FileDialogListView
        icon = self.style().standardIcon(pixmapi)
        btn.setIcon(icon)
        btn.setToolTip('Додати декілька')
        btn.clicked.connect(lambda _,action='add_many': self.action(action))
        self.table.hbox.addWidget(btn)

    def action(self, action: str, value: dict=None):
        if action == 'add_many':
            self.add_many_dialog(value)
            return
        return super().action(action, value)

    def prepare_value_to_action(self):
        value = super().prepare_value_to_action()
        value['list_name'] = self.list_name
        value['is_multiselect'] = self.is_multiselect
        return value
    
    def reload(self, values=None):
        if values is None:
            item = Item('operation_to_product')
            product_value = self.main_table.current_value
            err = item.get_filter_w('product_id', product_value['id'])
            if err:
                error(err)
                return
            values = [v for v in item.values if v['list_name'] == self.list_name]
        return super().reload(values)
            
    def dialog(self, value, title):
        i = Item(self.item.name)
        form = OperationToProductForm(value=value)
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

    def add_many_dialog(self, value):
        i = Item(self.item.name)
        dlg = SelectDialog('operation')
        res = dlg.exec()
        if res:
            values = dlg.widget.table.table.get_selected_values()
            for v in values:
                i.value = self.prepare_value_to_action()
                i.value['operation_id'] = v['id']
                i.value['cost'] = v['cost']
                err = i.save()
                if err:
                    error(err)
                    continue
            self.reload()
            self.actionResolved.emit()


class DetailsNumbersToProductTable(DetailsItemTable):
    def __init__(self, fields: list = [], values: list = None, list_name='default', is_multiselect=False):
        self.list_name = list_name
        self.is_multiselect = is_multiselect
        super().__init__('numbers_to_product', '', fields, values)
        btn = QPushButton()
        pixmapi = QStyle.StandardPixmap.SP_FileDialogListView
        icon = self.style().standardIcon(pixmapi)
        btn.setIcon(icon)
        btn.setToolTip('Додати декілька')
        btn.clicked.connect(lambda _,action='add_many': self.action(action))
        self.table.hbox.addWidget(btn)
        
    def action(self, action: str, value: dict=None):
        if action == 'add_many':
            self.add_many_dialog(value)
            return
        return super().action(action, value)
    
    def prepare_value_to_action(self):
        value = super().prepare_value_to_action()
        value['list_name'] = self.list_name
        value['is_multiselect'] = self.is_multiselect
        return value

    def prepare_value_to_save(self, value):
        return value
    
    def reload(self, values=None):
        if values is None:
            item = Item('numbers_to_product')
            product_value = self.main_table.current_value
            err = item.get_filter_w('product_id', product_value['id'])
            if err:
                error(err)
                return
            
        return super().reload(item.values)
    
    def dialog(self, value, title):
        i = Item(self.item.name)
        # i.create_default_w()
        # form = MatherialToProductForm(value=value)
        form = CustomForm(i.model, value=value)
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

    def add_many_dialog(self, value):
        i = Item(self.item.name)
        dlg = SelectDialog('matherial')
        res = dlg.exec()
        if res:
            values = dlg.widget.table.table.get_selected_values()
            for v in values:
                i.value = self.prepare_value_to_action()
                i.value['matherial_id'] = v['id']
                i.value['cost'] = v['cost']
                err = i.save()
                if err:
                    error(err)
                    continue
            self.reload()
            self.actionResolved.emit()

class DetailsToProductsTab(QWidget):
    def __init__(self, item: Item, main_table=None, fields=[]):
        super().__init__()
        self.main_table = main_table
        self.box = QVBoxLayout()
        self.box.setContentsMargins(0,0,0,0)
        self.setLayout(self.box)
        self.controls = QHBoxLayout()
        self.controls.setContentsMargins(0,0,0,0)
        top = QWidget()
        top.setLayout(self.controls)
        self.box.addWidget(top, 1)
        self.tabs = QTabWidget()
        self.box.addWidget(self.tabs, 10)
        self.controls.addStretch()
        add_list_btn = QPushButton('Додати список')
        self.controls.addWidget(add_list_btn)
        add_list_btn.clicked.connect(self.add_list)
        add_multi_btn = QPushButton('Додати мультивибір')
        self.controls.addWidget(add_multi_btn)
        add_multi_btn.clicked.connect(self.add_multi_list)
        self.item = item
        self.fields = fields
        
    def add_list(self):
        res = askdlg("Назва списку:")
        if res:
            table = create_details_table(
                self.item.name, 
                list_name=res,
                fields=self.fields,
                )
            table.set_main_table(self.main_table)
            self.tabs.addTab(table, res)
            self.tabs.setCurrentWidget(table)

    def add_multi_list(self):
        res = askdlg("Назва списку:")
        if res:
            table = create_details_table(
                self.item.name, 
                list_name=res, 
                is_multiselect=True,
                fields=self.fields,
                )
            table.set_main_table(self.main_table)
            self.tabs.addTab(table, '[+]'+res)
            self.tabs.setCurrentWidget(table)

    def reload(self, product_id):    
        err = self.item.get_filter_w('product_id', product_id)
        if err:
            error(err)
            return
        lists = {'default': []}
        if self.item.values and 'is_multiselect' in self.item.values[0]:
            for v in self.item.values:
                if v['is_multiselect']:
                    list_name = '[+]' + v['list_name']
                else:
                    list_name = v['list_name']   
                if list_name not in lists:
                    lists[list_name] = []
                lists[list_name].append(v)
        self.tabs.clear()
        for list_name in lists:
            list_tab_name = list_name[3:] if list_name.startswith('[+]') else list_name
            table = create_details_table(
                self.item.name, 
                values=lists[list_name], 
                list_name=list_tab_name,
                fields=self.fields,
                )
            table.set_main_table(self.main_table)
            self.tabs.addTab(table, list_name)
            table.reload()


def create_details_table(item_name: str, **kvargs):
    if item_name == 'operation_to_product':
        return DetailsOperationToProductTable(**kvargs)
    if item_name == 'matherial_to_product':
        return DetailsMatherialToProductTable(**kvargs)
    if item_name == 'product_to_product':
        return DetailsProductToProductTable(**kvargs)
    if item_name == 'numbers_to_product':
        return DetailsNumbersToProductTable(**kvargs)
    


class ItemsToProduct(QSplitter):
    def __init__(self) -> None:
        super().__init__()
        fields = [
            "id",
            "name",
            "measure",
            "cost",
            "width",
            "length",
        ]
        self.products = ProductTable(fields=fields)
        self.tabs = QTabWidget()
        self.addWidget(self.products)
        self.addWidget(self.tabs)
        fields = [
            "matherial",
            "number",
            "coeff",
            "cost",
            "list_name",
            "is_multiselect",
        ]
        self.m2ps = DetailsToProductsTab(
            Item('matherial_to_product'), 
            main_table=self.products,
            fields=fields,
            )
        self.tabs.addTab(self.m2ps, "Матеріал")
        
        fields = [
            "operation",
            "number",
            "coeff",
            "cost",
            "list_name",
            "is_multiselect",
            "equipment_cost",
        ]
        
        self.o2ps = DetailsToProductsTab(
            Item('operation_to_product'), 
            main_table=self.products,
            fields=fields,
            )
        self.tabs.addTab(self.o2ps, "Операція")
        
        fields = [
            "product2",
            "number",
            "coeff",
            "cost",
            "list_name",
            "is_multiselect",
        ]

        self.p2ps = DetailsToProductsTab(
            Item('product_to_product'),
            main_table=self.products,
            fields=fields,
            )
        self.tabs.addTab(self.p2ps, "Виріб")

        fields = [
            "number",
            "pieces",
            "size",
            "persent",
        ]

        self.n2ps = DetailsToProductsTab(
            Item('numbers_to_product'),
            main_table=self.products,
            fields=fields,
            )
        self.tabs.addTab(self.n2ps, "Кількість")
        
        self.setStretchFactor(0, 2)
        self.setStretchFactor(1, 3)
        self.products.table.table.valueSelected.connect(self.on_product_selected)
        recalc_btn = QPushButton('Перерахувати')
        row = self.products.info.grid.rowCount()
        col = self.products.info.grid.columnCount()
        self.products.info.grid.addWidget(recalc_btn, row-1, col-4)
        recalc_btn.clicked.connect(self.recalc_product)
        reverse_price_btn = QPushButton('Реверс ціни')
        self.products.info.grid.addWidget(reverse_price_btn, row-1, col-3)
        reverse_price_btn.clicked.connect(self.revers_product_price)

        update_prices_persent_btn = QPushButton('Оновити ціни на відсоток')
        self.products.info.grid.addWidget(update_prices_persent_btn, row-1, col-2)
        update_prices_persent_btn.clicked.connect(self.update_pricing_percent)
        update_list_prices_btn = QPushButton('Оновити ціни в списках')
        self.products.info.grid.addWidget(update_list_prices_btn, row-1, col-1)
        update_list_prices_btn.clicked.connect(self.update_prices_in_lists)
        update_prices_btn = QPushButton('Оновити ціни')
        self.products.info.grid.addWidget(update_prices_btn, row, col-1)
        update_prices_btn.clicked.connect(self.update_product_prices)
                
    def reload(self):
        self.products.reload()

    def update_pricing_percent(self, field):
        res = askdlg("Вкажіть відсоток")
        try:
            persent = int(res)
        except:
            return
        if not persent:
            return
        res = askdlg("Вкажіть точність до якого знаку після коми")
        try:
            approx = int(res)
        except:
            return
        
        item = Item('product')
        values = self.products.table.table.get_selected_values()
        if not values:
            error("Оберіть позиції")
            return
        for v in values:
            v['cost'] = round(v['cost'] + v['cost'] * persent / 100, approx)
            item.value = v
            err = item.save()
            if err:
                error(f'Не можу оновити {v["name"]}[{v["id"]}]:\n {err}')
                continue
            self.revers_product_price(v)
        
    def update_prices(self, name, with_defaults=False) -> int:
        item = Item(name)
        err = item.get_all()
        if err:
            error(err)
            return 0
        prices = {v['id']: v['cost'] for v in item.values}
        
        i2p = Item(f'{name}_to_product')
        err = i2p.get_all()
        if err:
            error(err)
            return 0
        
        update_counter = 0
        for v in i2p.values:
            if not with_defaults and v['list_name'] == 'default':
                continue
            id_name = f'product2_id' if name == 'product' else f'{name}_id'
            cost = round(v['coeff'] * prices[v[id_name]] * v['number'], 1)
            if cost == v['cost']:
                continue
            v['cost'] = cost
            i2p.value = v
            err = i2p.save()
            if err:
                error(err)
                continue
            update_counter += 1
            # print('updated', item.hum, v[id_name], 'to ordering', v['id'])
        messbox(f"{item.hum.title()} - ціни оновлено у {update_counter} позиціях")
        return update_counter
    
    def update_prices_in_lists(self):
        self.update_prices('matherial')
        self.update_prices('operation')
        self.update_prices('product')

    def update_product_prices(self):
        self.update_prices('matherial', with_defaults=True)
        self.update_prices('operation', with_defaults=True)
        self.update_prices('product', with_defaults=True)
        product = Item('product')
        err = product.get_all()
        if err:
            error(err)
            return 0
        prod_updated = 0
        for v in product.values:
            new_price = self.update_product_price(v)
            # print(v['name'], v['cost'], new_price)
            if not new_price or v['cost'] == new_price:
                continue
            v['cost'] = new_price
            product.value = v
            err = product.save()
            if err:
                error(err)
                return
            prod_updated += 1
        if not prod_updated:
            return 
        p2p_updated = self.update_prices('product', with_defaults=True)
        if not p2p_updated:
            return
        res = ok_cansel_dlg("Перервати?")
        if res:
            return
        self.update_product_prices()

    def on_product_selected(self, value):
        self.m2ps.reload(value['id'])
        self.o2ps.reload(value['id'])
        self.p2ps.reload(value['id'])
        self.n2ps.reload(value['id'])
        # o2p = Item('operation_to_ordering')
        # err = o2p.get_filter_w('ordering_id', value['id'])
        # if err:
        #     error(err)
        #     return
        # self.o2ps.reload(o2p.values)
        
    def update_sum(self):
        rows = self.products.table.table.get_selected_rows()
        if not rows:
            return
        row = rows[0]
        self.products.reload()
        index = self.products.table.table._model.createIndex(row, 0)
        self.products.table.table.setCurrentIndex(index)

    def update_product_price(self, prod_value):
        m2p = Item('matherial_to_product')
        err = m2p.get_filter_w('product_id', prod_value['id'])
        if err:
            error(err)
            return 0
        prod_sum = 0
        for v in m2p.values:
            if v['list_name'] == 'default':
                prod_sum += v['cost']
        # print('mat', mat_sum)
        o2p = Item('operation_to_product')
        err = o2p.get_filter_w('product_id', prod_value['id'])
        if err:
            error(err)
            return 0
        for v in o2p.values:
            if v['list_name'] == 'default':
                prod_sum += v['cost']
        # print('+oper', mat_sum)
        p2p = Item('product_to_product')
        err = p2p.get_filter_w('product_id', prod_value['id'])
        if err:
            error(err)
            return 0
        for v in p2p.values:
            if v['list_name'] == 'default':
                prod_sum += v['cost']
        return prod_sum
    
    def recalc_product(self):
        prod_value = self.products.current_value
        if not prod_value:
            return
        prod_sum = self.update_product_price(prod_value)
        if not prod_sum:
            return
        prod = Item('product')
        prod.value = prod_value
        prod.value['cost'] = prod_sum
        err = prod.save()
        if err:
            error(err)
            return
        self.update_sum()

    def revers_product_price(self, value=None):
        if not value:
            prod_value = self.products.current_value
            # print('prod_value', prod_value)
        else:
            prod_value = value
        
        if not prod_value:
            error('Оберіть позицію')
            return
        price = prod_value['cost']
        calc_price = self.update_product_price(prod_value)
        if not calc_price or price == calc_price:
            return
        k = price / calc_price
        # print('k=', k)
        for name in ('matherial_to_product', 'operation_to_product', 'product_to_product'):
            i2p = Item(name)
            err = i2p.get_filter_w('product_id', prod_value['id'])
            if err:
                error(err)
                return 0
            for v in i2p.values:
                if v['list_name'] == 'default':
                    v['cost'] = round(v['cost'] * k, 1)
                    v['coeff'] *= k
                    i2p.value = v
                    err = i2p.save()
                    if err:
                        error(err)
                        continue
