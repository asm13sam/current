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
        form = ProductForm(fields=i.columns, value=value)
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
        super().__init__('product_to_product', '', fields, values, deleted_buttons=False)
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
        form = ProductToProductForm(fields=i.columns, value=value)
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
        super().__init__(self.item.model, fields, value)
        mat =  self.widgets['matherial_id'].full_value()
        if mat:
            self.price = mat['cost']
        self.widgets['matherial_id'].valChanged.connect(self.matherial_selected)
        self.widgets['cost'].valChanged.connect(self.cost_changed)
        self.widgets['number'].valChanged.connect(self.number_changed)
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
        super().__init__('matherial_to_product', '', fields, values, deleted_buttons=False)
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
        form = MatherialToProductForm(fields=i.columns, value=value)
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
            return
        coeff = self.widgets['coeff'].value()
        self.widgets['cost'].set_value(round(self.price * number * coeff, 2))
        self.widgets['equipment_cost'].set_value(round(self.eq_price * number, 2))

    
class DetailsOperationToProductTable(DetailsItemTable):
    def __init__(self, fields: list = [], values: list = None, list_name='default', is_multiselect=False):
        self.list_name = list_name
        self.is_multiselect = is_multiselect
        super().__init__('operation_to_product', '', fields, values, deleted_buttons=False)
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
        form = OperationToProductForm(fields=i.columns, value=value)
        dlg = CustomFormDialog(title, form)
        res = dlg.exec()
        if res and dlg.value:
            i.value = dlg.value
            i.value['is_multiselect'] = self.is_multiselect
            print('is_multiselect', self.is_multiselect)
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
        super().__init__('numbers_to_product', '', fields, values, deleted_buttons=False)
        
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
        form = CustomForm(i.model, fields=i.columns, value=value)
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


class DetailsToProductsTab(QWidget):
    def __init__(self, item: Item, main_table=None, fields=[], list_btns=True):
        super().__init__()
        self.main_table = main_table
        self.item = item
        self.fields = fields
        self.product_id = 0

        self.box = QVBoxLayout()
        self.box.setContentsMargins(0,0,0,0)
        self.setLayout(self.box)
        if list_btns:
            self.controls = QHBoxLayout()
            self.controls.setContentsMargins(0,0,0,0)
            top = QWidget()
            top.setLayout(self.controls)
            self.controls.addStretch()
            self.box.addWidget(top, 1)
            add_list_btn = QPushButton('Додати список')
            self.controls.addWidget(add_list_btn)
            add_list_btn.clicked.connect(self.add_list)
            add_multi_btn = QPushButton('Додати мультивибір')
            self.controls.addWidget(add_multi_btn)
            add_multi_btn.clicked.connect(self.add_multi_list)

        self.tabs = QTabWidget()
        self.box.addWidget(self.tabs, 30)
        self.tabs.tabBarDoubleClicked.connect(self.edit_tab_title)

    def edit_tab_title(self):
        i = self.tabs.currentIndex()
        title = self.tabs.tabText(i)
        list_tab_name = title[3:] if title.startswith('[+]') else title
        res = askdlg(f"Вкажіть нову назву для списка '{list_tab_name}':")
        if not res:
            return
        err = self.item.get_filter_w('product_id', self.product_id)
        if err:
            error(err)
            return
        for v in self.item.values:
            if v['list_name'] == list_tab_name:
                v['list_name'] = res
                self.item.value = v
                err = self.item.save()
                if err:
                    error(err)
        self.reload()
        
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

    def reload(self, product_id=None):
        if product_id is None:
            if self.product_id:
                product_id = self.product_id
            else:
                return
        else:
            self.product_id = product_id    
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
                is_multiselect = list_name.startswith('[+]'),
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
            list_btns=False,
            )
        self.tabs.addTab(self.n2ps, "Кількість")
        
        self.setStretchFactor(0, 2)
        self.setStretchFactor(1, 3)
        self.products.table.table.valueSelected.connect(self.on_product_selected)
        
        row = self.products.info.grid.rowCount()
        col = self.products.info.grid.columnCount()
        recalc_btn = QPushButton('Перерахувати')
        self.products.info.grid.addWidget(recalc_btn, row-1, col-4)
        recalc_btn.clicked.connect(self.recalc_products)
        reverse_price_btn = QPushButton('Реверс ціни')
        self.products.info.grid.addWidget(reverse_price_btn, row-1, col-3)
        reverse_price_btn.clicked.connect(self.revers_product_prices)
        discard_coeff_btn = QPushButton('Скинути коефіцієнти')
        self.products.info.grid.addWidget(discard_coeff_btn, row-1, col-2)
        discard_coeff_btn.clicked.connect(self.discard_coefficients)
        update_prices_persent_btn = QPushButton('Оновити ціни на відсоток')
        self.products.info.grid.addWidget(update_prices_persent_btn, row-1, col-1)
        update_prices_persent_btn.clicked.connect(self.update_pricing_percent)
        update_prices_btn = QPushButton('Оновити ціни')
        self.products.info.grid.addWidget(update_prices_btn, row, col-1)
        update_prices_btn.clicked.connect(self.reload_prices)
                
    def reload(self):
        self.products.reload()

    def update_pricing_percent(self, field):
        values = self.products.table.table.get_selected_values()
        if not values:
            error("Оберіть позиції")
            return
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
        for v in values:
            v['cost'] = round(v['cost'] + v['cost'] * persent / 100, approx)
            item.value = v
            err = item.save()
            if err:
                error(f'Не можу оновити {v["name"]}[{v["id"]}]:\n {err}')
                continue
        self.revers_product_prices(values)
        self.reload()
        
    def on_product_selected(self, value):
        self.m2ps.reload(value['id'])
        self.o2ps.reload(value['id'])
        self.p2ps.reload(value['id'])
        self.n2ps.reload(value['id'])
        
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
        o2p = Item('operation_to_product')
        err = o2p.get_filter_w('product_id', prod_value['id'])
        if err:
            error(err)
            return 0
        for v in o2p.values:
            if v['list_name'] == 'default':
                prod_sum += v['cost']
        p2p = Item('product_to_product')
        err = p2p.get_filter_w('product_id', prod_value['id'])
        if err:
            error(err)
            return 0
        for v in p2p.values:
            if v['list_name'] == 'default':
                prod_sum += v['cost']
        return prod_sum
    
    def recalc_products(self):
        values = self.products.table.table.get_selected_values()
        if not values:
            error("Оберіть позиції")
            return
        prod = Item('product')
        for v in values:
            prod.value = v
            self.recalc_product(prod)
        self.update_sum()    

    def recalc_product(self, prod: Item):
        prod_sum = self.update_product_price(prod.value)
        if not prod_sum:
            return
        prod.value['cost'] = round(prod_sum)
        err = prod.save()
        if err:
            error(err)
            
    def revers_product_prices(self, _, values:dict=None):
        print(values)
        if values is None:
            values = self.products.table.table.get_selected_values()
        if not values:
            error("Оберіть позиції")
            return
        for v in values:
            self.revers_product_price(v)
        self.reload_details()
        
    def revers_product_price(self, prod_value):
        price = prod_value['cost']
        calc_price = self.update_product_price(prod_value)
        if not calc_price or price == calc_price:
            return
        k = price / calc_price
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

    def reload_details(self):
        self.m2ps.reload()
        self.o2ps.reload()
        self.p2ps.reload()
        self.n2ps.reload()

    def discard_coefficients(self):
        values = self.products.table.table.get_selected_values()
        if not values:
            error("Оберіть позиції")
            return
        for v in values:
            self.discard_coefficient(v)
        self.reload_details()
        
    def discard_coefficient(self, prod_value):
        for name in ('matherial_to_product', 'operation_to_product', 'product_to_product'):
            i2p = Item(name)
            err = i2p.get_filter_w('product_id', prod_value['id'])
            if err:
                error(err)
                return 0
            for v in i2p.values:
                v['cost'] = round(v['cost'] / v['coeff'], 1)
                v['coeff'] = 1
                i2p.value = v
                err = i2p.save()
                if err:
                    error(err)

    def reload_prices(self, _, values=None):
        if values is None:
            values = self.products.table.table.get_selected_values()
        if not values:
            error("Оберіть позиції")
            return
        
        for v in values:
            self.reload_price(v)
            
        self.reload()

    def reload_price(self, value):
        p2p = Item('product_to_product')
        err = p2p.get_filter_w('product_id', value['id'])
        if err:
            error(err)
            return
        prod = Item('product')
        for v in p2p.values:
            err = prod.get(v['product2_id'])
            if err:
                error(err)
                continue
            self.reload_price(prod.value)
        prod_sum = 0
        prod_sum += self.reload_price_for('matherial', value['id'])
        prod_sum += self.reload_price_for('operation', value['id'])
        prod_sum += self.reload_price_for('product', value['id'])
        value['cost'] = round(prod_sum)
        prod.value = value
        err = prod.save()
        if err:
            error(err)
                
    def reload_price_for(self, name, product_id):    
        item2p = Item(name + '_to_product')
        item = Item(name)
        res_sum = 0
        id_name = f'product2_id' if name == 'product' else f'{name}_id'
        err = item2p.get_filter_w('product_id', product_id)
        if err:
            error(err)
            return 0
        for v in item2p.values:
            err = item.get(v[id_name])
            if err:
                error(err)
                return 0
            v['cost'] = round(v['coeff'] * item.value['cost'] * v['number'], 2)
            item2p.value = v
            err = item2p.save()
            if err:
                error(err)
                return 0
            if v['list_name'] == 'default':
                res_sum += v['cost']
        return res_sum