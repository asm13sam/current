from PyQt6.QtCore import (
    Qt,
    QDate,
    pyqtSignal,
    QEvent,
    )
from PyQt6.QtGui import QIcon
from PyQt6.QtWidgets import (
    QLabel,
    QPushButton,
    QWidget,
    QLineEdit,
    QTextEdit,
    QCheckBox,
    QSpinBox,
    QDoubleSpinBox,
    QVBoxLayout,
    QHBoxLayout,
    QCalendarWidget,
    QSplitter,
    QApplication,
    QGridLayout,
    QTabWidget,
    )

from datetime import datetime, date
from data.model import Item
from data.app import App
from widgets.Dialogs import CustomDialog, error, DeleteDialog
from widgets.Table import TableWControls
from widgets.ComboBoxSelector import ComboBoxSelector
from common.params import MIN_SEARCH_STR_LEN, TABLE_BUTTONS
from widgets.Tree import TreeWControls

from common.funcs import phonef

class TabGroup(QTabWidget):
    def __init__(self):
        super().__init__()
        self.currentChanged.connect(self.reload_tab)

    def reload(self):
        pass

    def reload_tab(self, index):
        w = self.widget(index)
        w.reload()

class InfoBlock(QWidget):
    def __init__(self, data_model: dict, field_names: list = [], value:dict={}, columns=1, is_full=False):
        super().__init__()
        app = App()
        self.cfg = app.config
        fields = field_names if field_names else data_model.keys()
        if not is_full:
            fields = [name for name in fields if not name.endswith('_id')]
        headers = [data_model[name]['hum'] for name in fields]
        self.lables: dict[str, QLabel] = {}

        self.grid = QGridLayout()
        self.grid.setContentsMargins(0, 0, 0, 0)
        self.grid.setVerticalSpacing(1)
        # self.grid.setLabelAlignment( Qt.AlignmentFlag.AlignRight | Qt.AlignmentFlag.AlignTop )
        self.setLayout(self.grid)

        info_names_color = self.cfg['info_names_color'] if 'info_names_color' in self.cfg else 'lightgreen'
        self.setStyleSheet(f'color: {info_names_color}; padding: 0px;')
        info_values_color = self.cfg['info_values_color'] if 'info_values_color' in self.cfg else 'yellow'
        info_bg_color = self.cfg['info_bg_color'] if 'info_bg_color' in self.cfg else 'black'
        row = 0
        col = 0
        k = len(fields) // columns
        if len(fields) % columns:
            k += 1
        # if 'info' in fields:
        #     k -= 1
        
        align = Qt.AlignmentFlag.AlignRight | Qt.AlignmentFlag.AlignTop
        for field, header in zip(fields, headers):
            field_name = QLabel(header)
            if data_model[field]['def'] == '\n':
                self.lables[field] = TextWidget()
                self.lables[field].set_value("---" if not value else value[field] if field in value else '')
                # self.grid.addWidget(field_name, row, col, align)
                self.grid.addWidget(self.lables[field], k, 0, 1, 2*columns)
                self.grid.setRowStretch(k, 10)
                self.lables[field].setStyleSheet(f'color:{info_values_color}; padding: 2px; background-color: {info_bg_color};')
                row -= 1
                self.lables[field].setReadOnly(True)
            else:
                ttp = ''
                if not value:
                    txt = "---" 
                elif field in value:
                    txt, ttp = self.prep_value(value[field])  
                else:
                    txt = ''
                l = QLabel(txt)
                self.lables[field] = l
                if ttp:
                    l.setToolTip(ttp)
                self.grid.addWidget(field_name, row, col, align)
                self.grid.addWidget(self.lables[field], row, col+1)
                self.lables[field].setStyleSheet(f'min-width: 5em;color:{info_values_color}; padding: 2px; background-color: {info_bg_color};')
                self.lables[field].setWordWrap(True)
            self.lables[field].installEventFilter(self)
            row += 1
            if row >= k:
                row = 0
                col += 2
        # self.grid.addWidget(QLabel(''), k+2, 0)
        if not 'info' in fields:
            self.grid.setRowStretch(k+5, 10)
            
    def reload(self, value=None):
        if not value:
            self.clear()
            return
        for field, label in self.lables.items():
            if type(value[field]) == bool:
                label.setText('Так' if value[field] else 'Ні')
            elif field == 'info':
                label.setText(value[field])
                label.setToolTip(value[field])    
            else:
                txt, ttp = self.prep_value(value[field])
                label.setText(txt)
                if ttp:
                    label.setToolTip(ttp)
        if 'phone' in self.lables.keys():
            self.lables['phone'].setText(phonef(value['phone']))

    def clear(self):
        for label in self.lables.values():
            label.setText('--')

    def eventFilter(self, source, event: QEvent):
        if event.type() == QEvent.Type.MouseButtonDblClick:
            clipboard = QApplication.clipboard()
            clipboard.setText(source.text())
        return super(InfoBlock, self).eventFilter(source, event)
    
    @staticmethod
    def prep_value(value):
        ttp = ''
        if type(value) == float:
            v = str(round(value, 2))
        else:
            v = str(value)
        if len(v) > 35:
            ttp = v
            v = v[:35] + '>'
        return v, ttp



class CustomForm(QWidget):
    formChanged = pyqtSignal(bool)
    saveRequested = pyqtSignal(dict)
    def __init__(self, 
                 data_model: dict,
                 fields: list=[],
                 value: dict={},
                 columns: int = 1,
                 measure: str = '',
                ):
        super().__init__()
        app = App()
        self.cfg = app.config
        self.value = value
        self.widgets = {}
        self.labels = {}
        self.measure = measure
        self.model = data_model
        self.fields = fields if fields else self.model.keys()
        self.fields = [field for field in self.fields if self.model[field]['def'] != []]
        
        self.grid = QGridLayout()
        self.setLayout(self.grid)
        self.grid.setContentsMargins(0, 0, 0, 0)
        form_names_color = self.cfg['form_names_color'] if 'form_names_color' in self.cfg else 'lightgreen'
        self.setStyleSheet(f'color: {form_names_color}; padding: 0px;')
        rows_number = self.create_widgets(columns)

        self.save_but = QPushButton("Зберегти")
        self.grid.addWidget(self.save_but, rows_number+1, columns*2-1)
        self.save_but.clicked.connect(self.get_value)
        self.is_changed = False

    def hide_save_btn(self):
        self.save_but.setVisible(False)

    def set_changed(self, changed=True):
        # print('calling set changed', changed)
        if changed == self.is_changed:
            return
        self.is_changed = changed
        # print('==call form set changed==')
        self.formChanged.emit(changed)

    def changed(self):
        return self.is_changed

    def create_widgets(self, columns):
        form_values_color = self.cfg['form_values_color'] if 'form_values_color' in self.cfg else 'yellow'
        form_bg_color = self.cfg['form_bg_color'] if 'form_bg_color' in self.cfg else 'black'
        align = Qt.AlignmentFlag.AlignRight
        row = 0
        col = 0
        k = len([f for f in self.fields if self.model[f]['form'] and self.model[f]['def'] != '\n']) // columns
        for field in self.fields:
            if not self.model[field]['form']:
                continue
            def_value = self.model[field]['def']
            value = self.value[field] if self.value else def_value
            if 'group_id' in self.model[field]:
                group_id = self.model[field]['group_id']
            else:
                group_id=0
            w = self.create_widget(field, def_value, value, group_id)
            if w is None:
                continue
            w.valChanged.connect(self.set_changed)
            w.setStyleSheet(f'color:{form_values_color}; padding: 2px; background-color: {form_bg_color};')
            self.widgets[field] = w
            self.labels[field] = QLabel(self.model[field]['hum'])
            
            if def_value == '\n':
                self.grid.addWidget(self.widgets[field], k, 0, 1, 2*columns)
                # self.grid.addWidget(self.labels[field], 0, 2*columns, align)
                # self.grid.addWidget(self.widgets[field], 0, 2*columns+1, k+1, 1)
                continue    

            self.grid.addWidget(self.labels[field], row, col, align)
            self.grid.addWidget(self.widgets[field], row, col+1)
            row += 1
            if row > k:
                row = 0
                col += 2
        self.grid.addWidget(QLabel(''), k+2, 0)
        self.grid.setRowStretch(k+5, 10)
        for w in self.widgets.values():
            if type(w) == Selector and w.item_name in ('matherial', 'operation', 'product') and 'number' in self.widgets:
                w.set_number_widget(self.widgets['number'])
        return row
            

    def create_widget(self, field, def_value, cur_value, group_id=0):
        t = type(def_value)
        # date, text or string
        if t == str:
            return self.create_str_widget(def_value, cur_value)
        # id or int
        if t == int:
            return self.create_int_widget(cur_value, field, group_id)
        # sum of money, number or float
        if t == float:
            return self.create_float_widget(cur_value, field)
        if t == bool:
            return self.create_bool_widget(cur_value)
        
    def create_str_widget(self, def_value, cur_value):
        if def_value == 'date':
            w = TimeWidget()
        elif '\n' in def_value:
            w = TextWidget()
        else:
            w = LineEditWidget()
        w.setValue(cur_value)
        return w

    def create_int_widget(self, cur_value, field, group_id=0):
        if field == 'id' or field == 'document_uid':
            w = IdWidget()
        elif field.endswith('_id'):
            if field == 'contact_id':
                w = ContactSelector(self.widgets['contragent_id'] if 'contragent_id' in self.widgets else None)
                w.setValue(cur_value)
                return w
            name = '_'.join(field.split('_')[:-1])   # "matherial_id" split to material and id
            # as it is foring key, need to select other item wich id is here
            w = Selector(name, group_id=group_id, form_value=self.value)
        else:
            # if it is simple integer value
            w = IntWidget()
            
        w.setValue(cur_value)
        return w

    def create_float_widget(self, cur_value, field):
        if field == 'number':
            w = NumWidget(self.measure)
        if field == 'price':
            w = PriceWidget()
        elif field == 'persent':
            w = PersentWidget()
        elif field == 'profit':
            w = ProfitWidget()
        else:
            w = CostWidget()
        # elif field == 'cost':
        #     w = CostWidget()
        # else:
        #     w = QDoubleSpinBox()
        w.setValue(cur_value)
        return w

    def create_bool_widget(self, cur_value):
        w = CheckWidget()
        w.setValue(cur_value)
        return w

    def get_value(self):
        for k in self.widgets:
            v = self.widgets[k].value()
            # print(self.model[k]['form'], v)
            if self.model[k]['form'] == 2 and not v:
                error(f'Не обрано {self.model[k]["hum"]}')
                return False
            self.value[k] = v
            # print('>|', k, self.value[k])
        self.saveRequested.emit(self.value)
        return True
        
    def reload(self, value):
        # print('start reload')
        self.value = value
        # print(value)
        for k, w in self.widgets.items():
            # print(k, value[k])
            w.setValue(value[k])
            
        self.set_changed(False)


class TimeWidget(QWidget):
    valChanged = pyqtSignal()
    def __init__(self, value='date', with_time=True):
        super().__init__()
        if value != 'date':
            self.date = value
        elif with_time:
            self.date = datetime.now().isoformat(timespec='seconds')
        else:
            self.date = date.today().isoformat() #.now().isoformat(timespec='hours')
        self.layout = QHBoxLayout()
        self.setLayout(self.layout)
        # self.label = QLabel()
        self.btn = QPushButton('Обрати')
        self.btn.clicked.connect(self.act)
        # self.layout.addWidget(self.label)
        self.layout.addWidget(self.btn)
        self.btn.setText(self.date)

    def setValue(self, value):
        self.date = value if value != 'date' else datetime.now().isoformat(timespec='seconds')
        self.btn.setText(self.date)
        
    def set_value(self, value):
        self.setValue(value)
        self.valChanged.emit()

    def value(self):
        return self.date

    def act(self):
        dlg = TimeDialog(self.date)
        res = dlg.exec()
        if not res:
            return
        self.set_value(dlg.widget.selectedDate().toString(Qt.DateFormat.ISODate))


class TimeDialog(CustomDialog):
    def __init__(self, value=''):
        cal = QCalendarWidget()
        cal.setGridVisible(True)
        if value:
            cal.setSelectedDate(QDate.fromString(value, Qt.DateFormat.ISODate))
        super().__init__(cal, "Оберіть дату")


class TextWidget(QTextEdit):
    valChanged = pyqtSignal()
    def __init__(self):
        super().__init__()
        self.textChanged.connect(lambda: self.valChanged.emit())

    def setValue(self, value: str):
        value = value.strip()
        self.setText(value)
    
    def set_value(self, value):
        self.setValue(value)
        self.valChanged.emit()

    def value(self):
        return self.toPlainText()


class CheckWidget(QCheckBox):
    valChanged = pyqtSignal()
    def __init__(self):
        super().__init__()
        self.stateChanged.connect(lambda: self.valChanged.emit())

    def setValue(self, value: bool):
        self.setCheckState(Qt.CheckState.Checked if value else Qt.CheckState.Unchecked)
        
    def set_value(self, value):
        self.setValue(value)
        self.valChanged.emit()

    def value(self):
        return self.checkState() == Qt.CheckState.Checked


class LineEditWidget(QLineEdit):
    valChanged = pyqtSignal()
    def __init__(self):
        super().__init__()
        self.textChanged.connect(lambda: self.valChanged.emit())

    def setValue(self, value: str):
        self.setText(value)

    def set_value(self, value):
        self.setValue(value)
        self.valChanged.emit()

    def value(self):
        return self.text()


class IdWidget(QLabel):
    valChanged = pyqtSignal()
    def __init__(self):
        super().__init__()

    def setValue(self, id, name=''):
        self.id = id
        self.name = name
        if name:
            self.setText(name)
        else:
            self.setText(str(id))
        
    def set_value(self, id, name=''):
        self.setValue(id, name)
        self.valChanged.emit()

    def value(self):
        return self.id


class Selector(QWidget):
    valChanged = pyqtSignal()
    def __init__(self, field, title='Обрати', group_id=0, form_value={}):
        super().__init__()
        self.form_value = form_value
        self.item_name = field[:-1] if field.endswith('2') else field
        self.val = {}
        self.group_id = group_id
        self.number_widget = None
        self.box = QVBoxLayout()
        self.setLayout(self.box)
        self.box.setContentsMargins(0, 0, 0, 0)
        self.btn = QPushButton(title)
        self.btn.clicked.connect(self.act)
        self.btn.setStyleSheet('text-align:left;')
        self.box.addWidget(self.btn)
        self.setStyleSheet('padding: 1px; margin: 0px;')

    def set_number_widget(self, number_widget):
        self.number_widget = number_widget
        if self.val:
            self.number_widget.setSuffix(' ' + self.val['measure'])
    
    def act(self):
        dlg = SingleSelectDialog(self.item_name, group_id=self.group_id, form_value=self.form_value)
        res = dlg.exec()
        if not res:
            return
        self.set_value(dlg.value)

    def value(self):
        return self.val['id'] if self.val else 0
    
    def full_value(self):
        return self.val

    def setValue(self, value=None):
        self.val = {}
        if not value:
            self.btn.setText('')
            return
        if type(value) != dict:
            value = {'id': value}
            item = Item(self.item_name)
            err = item.get_w(value['id'])
            if err:
                error(err)
                return
            self.val = item.value
        else:
            self.val = value
        
        if 'name' in self.val:
            self.btn.setText(self.val['name'])
        else:
            self.btn.setText(str(self.val['id']))
        if self.number_widget:
            self.number_widget.setSuffix(' ' + self.val['measure'])
        
    def set_value(self, value=None):
        self.setValue(value)    
        self.valChanged.emit()


class ContactSelector(Selector):
    def __init__(self, contragent_widget=None, title='Обрати'):
        self.contragent_widget = contragent_widget
        super().__init__('contact', title)
        self.label = QLabel()
        self.box.addWidget(self.label)
        self.label.installEventFilter(self)

    def act(self):
        contragent_id = 0
        if self.contragent_widget:
            contragent_id = self.contragent_widget.value()
            if not contragent_id:
                error('Оберіть контрагента')
                return
        dlg = ContactSelectDialog(contragent_id)
        res = dlg.exec()
        if not res:
            return
        self.set_value(dlg.value)
        
    def setValue(self, value=None):
        super().setValue(value)
        if 'phone' in self.val:
            self.label.setText(self.val['phone'])

    def eventFilter(self, source, event):
        if event.type() == QEvent.Type.MouseButtonDblClick:
            clipboard = QApplication.clipboard()
            clipboard.setText(source.text())
        return super().eventFilter(source, event)
    
    def remove_phone(self):
        self.label.hide()
       

class FormDialog(CustomDialog):
    def __init__(self, title, data_model, value):
        self.value = {}
        self.form = CustomForm(data_model=data_model, value=value)
        super().__init__(self.make_widget(self.form), title)
        self.form.saveRequested.connect(self.get_value)
        self.form.setMinimumWidth(300)
        self.form.hide_save_btn()

    def make_widget(self, form):
        return form

    def get_value(self, value):
        self.value = value
        self.form.set_changed(False)

    def accept(self) -> None:
        # if self.form.changed():
        if not self.form.get_value():
            return
            # mess = "Відкинути не збережені зміни?"
            # if not ok_cansel_dlg(mess):
            #     return
        return super().accept()


class CustomFormDialog(CustomDialog):
    def __init__(self, title, form, table=None):
        self.value = {}
        self.form = form
        self.table = table
        self.split = QSplitter(Qt.Orientation.Horizontal)
        self.split.addWidget(self.form)
        min_width = 300
        if table:
            self.split.addWidget(self.table)
            min_width = 900
            self.split.setStretchFactor(0, 1)
            self.split.setStretchFactor(1, 3)
            
        super().__init__(self.split, title)
        self.form.saveRequested.connect(self.get_value)
        self.split.setMinimumWidth(min_width)
        self.form.hide_save_btn()

    def get_value(self, value):
        self.value = value
        self.form.set_changed(False)

    def accept(self) -> None:
        # if self.form.changed():
        if not self.form.get_value():
            return
            # mess = "Відкинути не збережені зміни?"
            # if not ok_cansel_dlg(mess):
            #     return
        return super().accept()


class SelectDialog(CustomDialog):
    def __init__(self, item_name: str, group_id=0, form_value={}):
        self.value = None
        self.form_value = form_value
        self.item = Item(item_name)
        if group_id:
            values = None
        else:
            values = self.set_values()
        
        self.search_field = ''
        if item_name == 'contragent':
            self.search_field = 'search'
                
        # if self.item.name + '_id' in self.item.model:
        #     widget = GroupTree(self.item.name, self.item.hum, values, group_id)
        #     if self.form_value:
        #         widget.tree.delete_by_id(self.form_value['id'])
        # else:
        widget = ItemTable(
            self.item.name, 
            search_field=self.search_field, 
            values=values,
            group_id=group_id,
        )

        widget.setMinimumWidth(1000)
        widget.setMinimumHeight(600)
        
        super().__init__(widget, f"Обрати {self.item.hum}")
        
    def set_values(self):    
        err = self.item.get_all_w()
        if err:
            error(err)
            return
        return self.item.values

    
class ContactSelectDialog(SelectDialog):
    def __init__(self, contragent_id):
        self.contragent_id = contragent_id
        # print('contragent id in contact selector', contragent_id)
        super().__init__('contact')
        self.widget.table.table.valueDoubleCklicked.disconnect()
        self.widget.table.table.valueDoubleCklicked.connect(self.accept_value)
        self.widget.setMinimumWidth(1000)

    def do_action(self, action: str, value: dict):
        if action == 'create':
            action = 'copy'
            item = Item(self.item.name)
            item.create_default()
            item.value['contragent_id'] = self.contragent_id
            value = item.value
        return super().do_action(action, value)
    
    def set_values(self):
        if self.contragent_id:
            err = self.item.get_filter_w('contragent_id', self.contragent_id)
        else:
            err = self.item.get_all_w()
        if err:
            error(err)
            return
        return self.item.values
    
    def accept_value(self, value):
        self.value = value
        return super().accept()
    
    def accept(self):
        self.value = self.widget.table.table.get_selected_value()
        return super().accept()


class SingleSelectDialog(SelectDialog):
    def __init__(self, item_name: str, group_id=0, form_value={}):
        super().__init__(item_name, group_id, form_value)
        self.widget.remove_dblclick_cb()
        self.widget.set_dblclick_cb(self.accept_value)

    def accept_value(self, value):
        self.value = value
        return super().accept()
    
    def accept(self):
        self.value = self.widget.value()
        return super().accept()


class MultiSelectDialog(SelectDialog):
    def __init__(self, item_name: str):
        super().__init__(item_name)
        self.values = []

    def accept(self) -> None:
        self.values = self.widget.table.table.get_selected_values()
        return super().accept()


class NumWidget(QDoubleSpinBox):
    valChanged = pyqtSignal()
    def __init__(self, measure: str=''):
        super().__init__()
        self.setRange(0.0, 10000.0)
        self.valueChanged.connect(lambda: self.valChanged.emit())
        if measure:
            self.setSuffix(' ' + measure)

    def setValue(self, val: float) -> None:
        val = round(val, 2)
        return super().setValue(val)    
    
    def set_value(self, value):
        self.blockSignals(True)
        self.setValue(value)
        self.blockSignals(False)
        

# price do not sets by user, only calculates, so it is label with interface of spinbox
# class PriceWidget(QLabel):
#     valChanged = pyqtSignal()
#     def __init__(self):
#         super().__init__()
#         self.val = 0.0
#         self.init_val = 0.0
#         #self.setSuffix(measure)

#     def setValue(self, value):
#         self.val = value
#         self.setText(str(self.val))
        
#     def set_value(self, value):
#         self.setValue(value)
#         self.valChanged.emit()

#     def value(self):
#         return self.val

class PriceWidget(QDoubleSpinBox):
    valChanged = pyqtSignal()
    def __init__(self):
        super().__init__()
        self.setRange(-100000000.0, 100000000.0)
        self.valueChanged.connect(lambda: self.valChanged.emit())

    def setValue(self, val: float) -> None:
        val = round(val, 2)
        return super().setValue(val)
    
    def set_value(self, value):
        self.blockSignals(True)
        self.setValue(value)
        self.blockSignals(False)
        # self.valChanged.emit()


class PersentWidget(QDoubleSpinBox):
    valChanged = pyqtSignal()
    def __init__(self):
        super().__init__()
        self.setRange(-100000000.0, 100000000.0)
        self.valueChanged.connect(lambda: self.valChanged.emit())
        
    def setValue(self, val: float) -> None:
        val = round(val, 2)
        return super().setValue(val)
    
    def set_value(self, value):
        self.blockSignals(True)
        self.setValue(value)
        self.blockSignals(False)
        

class ProfitWidget(QDoubleSpinBox):
    valChanged = pyqtSignal()
    def __init__(self):
        super().__init__()
        self.setRange(-100000000.0, 100000000.0)
        self.valueChanged.connect(lambda: self.valChanged.emit())

    def setValue(self, val: float) -> None:
        val = round(val, 2)
        return super().setValue(val)
    
    def set_value(self, value):
        self.blockSignals(True)
        self.setValue(value)
        self.blockSignals(False)
        

class CostWidget(QDoubleSpinBox):
    valChanged = pyqtSignal()
    def __init__(self):
        super().__init__()
        self.setRange(-100000000.0, 100000000.0)
        self.valueChanged.connect(lambda: self.valChanged.emit())

    def setValue(self, val: float) -> None:
        val = round(val, 2)
        return super().setValue(val)

    def set_value(self, value):
        self.blockSignals(True)
        self.setValue(value)
        self.blockSignals(False)
        

class IntWidget(QSpinBox):
    valChanged = pyqtSignal()
    def __init__(self):
        super().__init__()
        self.setRange(-1000000000, 1000000000)
        self.valueChanged.connect(lambda: self.valChanged.emit())

    def set_value(self, value):
        self.blockSignals(True)
        self.setValue(value)
        self.blockSignals(False)
        

class PeriodWidget(QWidget):
    periodSelected = pyqtSignal(str, str)
    def __init__(self):
        super().__init__()
        box = QHBoxLayout()
        self.setLayout(box)
        box.setContentsMargins(0,0,0,0)
        start_date = date.today().replace(day=1)
                
        self.date_from = TimeWidget(start_date.isoformat(), with_time=False)
        self.date_to = TimeWidget(with_time=False)
        box.addWidget(QLabel("З:"))
        box.addWidget(self.date_from)
        box.addWidget(QLabel("до:"))
        box.addWidget(self.date_to)
        self.date_from.valChanged.connect(self.period_changed)
        self.date_to.valChanged.connect(self.period_changed)

    def period_changed(self):
        date_from, date_to = self.get_period()
        if date_from > date_to:
            error('Неправильний порядок дат')
            return
        self.periodSelected.emit(date_from, date_to)

    def get_period(self):
        date_from = self.date_from.value()[0:10]
        date_to = self.date_to.value()[0:10] + 'T23:59:59'
        return date_from, date_to


class ItemTable(QSplitter):
    actionResolved = pyqtSignal()
    def __init__(
            self, 
            item_name:str, 
            search_field:str='', 
            fields:list=[], 
            values:list=None,
            buttons = TABLE_BUTTONS,
            group_id=0,
            is_vertical_inner=True,
            is_info_bottom=False,
            show_period=True,
            ):
        if is_vertical_inner:
            super().__init__()
            self.inner = QSplitter(Qt.Orientation.Vertical)
            self.inner.setStretchFactor(0, 10)
        else:
            super().__init__(Qt.Orientation.Vertical)
            self.inner = QSplitter()
        self.addWidget(self.inner)
        self.setStretchFactor(0, 10)
        
        self.item = Item(item_name)
        self.table = TableWControls(
            data_model=self.item.model_w, 
            with_search=bool(search_field),
            table_fields=fields,
            buttons=buttons,
        )
        self.search_field = search_field
        if search_field:
            self.table.searchChanged.connect(self.search_changed)
        
        self.group_name = ''
        group_name = item_name + '_group' 
        if group_name + '_id' in self.item.model:
            self.make_group_selector(group_name, group_id)
            # reload_group_btn = QPushButton()
            # reload_group_btn.setIcon(QtGui.QIcon(f'images/icons/reload.png'))
            # reload_group_btn.setToolTip('Оновити групи')
            # reload_group_btn.clicked.connect(self.reload_group)
            # self.table.hbox.insertWidget(0, reload_group_btn)

        self.period = None
        # if 'created_at' in self.item.model and show_period:
        if ('between' in self.item.op_model or 'between_up' in self.item.op_model) and show_period:
            self.period = PeriodWidget()
            self.table.hbox.insertWidget(0, self.period)
            self.period.periodSelected.connect(self.period_changed)
                
        self.inner.addWidget(self.table)
        self.inner.setStretchFactor(0, 4)

        if is_info_bottom:
            self.info = InfoBlock(self.item.model_w, columns=2)
            if is_vertical_inner:
                self.inner.addWidget(self.info)
                self.inner.setStretchFactor(0, 3)
                self.inner.setStretchFactor(1, 1)
            else:
                self.addWidget(self.info)
                self.setStretchFactor(0, 3)
                self.setStretchFactor(1, 1)
        else:
            self.info = InfoBlock(self.item.model_w)
            if not is_vertical_inner:
                self.inner.addWidget(self.info)
            else:
                self.addWidget(self.info)
        
        self.table.actionInvoked.connect(self.action)
        self.table.table.valueDoubleCklicked.connect(self.on_double_click)
        self.table.table.valueSelected.connect(self.on_value_selected)

        if buttons:
            all_btn = QPushButton('Усі')
            self.table.hbox.addWidget(all_btn)
            all_btn.clicked.connect(self.show_all)
            del_btn = QPushButton('Видалені')
            self.table.hbox.addWidget(del_btn)
            del_btn.clicked.connect(self.show_deleted)

        if values:
            self.reload(values)

    def show_all(self):
        values = self.get_values('all')
        self.reload(values)

    def show_deleted(self):
        values = self.get_values('deleted')
        self.reload(values)

    def on_value_selected(self, value):
        self.info.reload(value)
    
    def calc_sum(self, field):
        return self.table.table._model.calc_sum(field)
    
    def value(self):
        return self.table.table.get_selected_value()
    
    def values(self):
        return self.table.values()
    
    def set_dblclick_cb(self, cb):
        self.table.table.valueDoubleCklicked.connect(cb)
        
    def remove_dblclick_cb(self):
        self.table.table.valueDoubleCklicked.disconnect()
    
    def make_group_selector(self, group_name, group_id):
        group = Item(group_name)
        err = group.get_all()
        if err:
            error(err)
            return
        if group_name + '_id' not in group.model:
            self.groups = ComboBoxSelector(values=group.values)
            self.table.hbox.insertWidget(0, QLabel('Група'))
            self.table.hbox.insertWidget(1, self.groups)
            self.groups.selectionChanged.connect(self.group_changed)
            self.group_name = group_name
        else:    
            self.setOrientation(Qt.Orientation.Horizontal)
            self.groups = GroupTree(group_name, group.hum, group.values)
            self.groups.tree.tree.itemSelected.connect(self.group_changed)
            self.group_name = group_name
            self.groups.setMaximumWidth(300)
            self.insertWidget(0, self.groups)
        self.groups.set_current_id(group_id)
        # self.setStretchFactor(0, 1)
        # self.setStretchFactor(1, 10)
        
    # def reload_group(self):
    #     group_value = self.groups.value()

    #     group = Item(self.group_name)
    #     err = group.get_all()
    #     if err:
    #         error(err)
    #         return
    #     self.groups.reload(group.values)
    #     if group_value:
    #         self.groups.set_current_id(group_value['id'])
    
    def get_group_values(self, all=''):
        group = self.groups.value()
        if group and group['id']:
            if all == 'all':
                return self.item.get_filter_w_with_deleted(self.group_name + '_id', group['id'])
            if all == 'deleted':
                return self.item.get_filter_w_deleted_only(self.group_name + '_id', group['id'])
            return self.item.get_filter_w(self.group_name + '_id', group['id'])
        if all == 'all':
            return self.item.get_all_w_with_deleted()
        if all == 'deleted':
            return self.item.get_all_w_deleted_only()
        return self.item.get_all_w()
                    
    def get_values(self, all=''):
        # group
        if self.group_name:
            err = self.get_group_values(all)
        elif self.period is not None:
            date_from, date_to = self.period.get_period()
            if all == 'all':
                err = self.item.get_between_w_with_deleted('created_at', date_from, date_to)
            elif all == 'deleted':    
                err = self.item.get_between_w_deleted_only('created_at', date_from, date_to)
            else:
                err = self.item.get_between_w('created_at', date_from, date_to)
        else:
            if all == 'all':
                err = self.item.get_all_w_with_deleted()
            elif all == 'deleted':
                err = self.item.get_all_w_deleted_only()
            else:
                err = self.item.get_all_w()
        if err:
            error(err)
            return
        
        values = self.item.values
        
        # search
        if self.search_field:
            text = self.table.search_entry.text()
            if len(text) >= MIN_SEARCH_STR_LEN:
                values = [v for v in values if text in v[self.search_field].lower()]
        return values
    
    def reload(self, values=None):
        if values is None:
            values = self.get_values()
        self.table.table.reload(values)
        self.info.reload()

    def search_changed(self, text):
        if len(text) < MIN_SEARCH_STR_LEN:
            values = self.item.values
        else:    
            text = text.lower()
            values = [v for v in self.item.values if text in v[self.search_field].lower()]
        self.table.table.reload(values)

    def group_changed(self, value):
        if not value['id']:
            err = self.item.get_all_w()
        else:    
            err = self.item.get_filter_w(self.group_name + '_id', value['id'])
        if err:
            error(err)
            return
        self.table.table.reload(self.item.values)

    def period_changed(self, date_from, date_to):
        err = self.item.get_between_w('created_at', date_from, date_to)
        if err:
            error(err)
            return
        self.table.table.reload(self.item.values)

    def action(self, action:str, value:dict=None):
        # print('action', action)
        if action == 'reload':
            self.reload()
            return
        if action == 'delete':
            self.del_dialog(value)
            return
        title = 'Створити'
        if action == 'copy':
            value['id'] = 0
        elif action == 'create':
            value = self.prepare_value_to_action()
            if value is None:
                return
        elif action == 'edit':
            title = 'Редагувати'
        self.dialog(value, title)

    def prepare_value_to_action(self):
        item = Item(self.item.name)
        item.create_default()
        return item.value

    def dialog(self, value, title):
        i = Item(self.item.name)
        # print(self.item.name, '====>', i.model.keys())
        dlg = FormDialog(title, i.model, value)
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
        return value
    
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
                    i.value['info'] = f'\nПричина видалення:\n{cause}\n' + value['info']
                    i.save()
                
            err = i.delete(value['id'], cause)
            if err:
                error(err)
                return
            self.reload()
            self.actionResolved.emit()

    def on_double_click(self, value):
        self.action('edit', value)
    

class MainItemTable(ItemTable):
    def __init__(self, 
                item_name: str, 
                search_field: str = '', 
                fields: list = [], 
                values: list = None, 
                buttons=TABLE_BUTTONS, 
                group_id=0,
                is_vertical_inner=True,
                is_info_bottom=False,
                show_period=True,
                ):
        super().__init__(
            item_name, 
            search_field, 
            fields, 
            values, 
            buttons, 
            group_id, 
            is_vertical_inner, 
            is_info_bottom,
            show_period,
            )
        self.details_table = None
        self.doc_table = None
        self.is_vertical_inner = is_vertical_inner
        self.is_info_bottom = is_info_bottom
        self.current_value = None
        
    def set_detais_table(self, table: ItemTable):
        self.details_table = table

    def del_dialog(self, value):
        # if self.details_table.
        return super().del_dialog(value)
    
    def add_doc_table(self, doc_table, is_bottom=False):
        self.doc_table = doc_table
        if is_bottom:
            self.inner.addWidget(doc_table)
            return
        if self.is_info_bottom:
            if not self.is_vertical_inner:
                self.inner.addWidget(doc_table)    
            else:
                self.addWidget(doc_table)
        else:
            if self.is_vertical_inner:
                self.inner.addWidget(doc_table)    
            else:
                self.addWidget(doc_table)
    
    def on_value_selected(self, value):
        # print('on value main table selected')
        self.current_value = value
        if self.doc_table:
            self.doc_table.reload(value)
        # print('on value main table selected end')
        return super().on_value_selected(value)
    
    def reload(self, values=None):
        if self.details_table:
            self.details_table.reload()
        if self.doc_table:
            self.doc_table.reload()
        return super().reload(values)


class DetailsItemTable(ItemTable):
    def __init__(self, item_name: str, search_field: str = '', fields: list = [], values: list = None, buttons=TABLE_BUTTONS, group_id=0, show_period=False):
        super().__init__(item_name, search_field, fields, values, buttons, group_id, False, True, show_period)
        self.main_table = None

    def set_main_table(self, table: ItemTable):
        self.main_table = table
    
    def prepare_value_to_action(self):
        value = super().prepare_value_to_action()
        if self.main_table is not None:
            cur_value = self.main_table.table.table.get_selected_value()
            if not cur_value:
                return
            value[self.main_table.item.name + '_id'] = cur_value['id']
        return value
    
    def prepare_value_to_save(self, value):
        if 'cost' in value and 'number' in value and 'price' in value:
            value['cost'] = value['number'] * value['price']
        return value
    
    def reload(self, values=None):
        if values is None:
            cur_value = self.main_table.table.table.get_selected_value()
            if cur_value:
                self.item.get_filter_w(self.main_table.item.name + '_id', cur_value['id'])
                values = self.item.values
            else:
                values = self.get_values()
        self.table.table.reload(values)


class ItemTableWithDetails(QSplitter):
    def __init__(
            self, 
            main_table: MainItemTable,
            details_table: DetailsItemTable,
            ):
        super().__init__()
        self.table = main_table
        self.details = details_table
        self.table.set_detais_table(self.details)
        self.details.set_main_table(self.table)
        self.addWidget(self.table)
        self.addWidget(self.details)
        self.table.table.table.valueSelected.connect(self.reload_details)
        
    def reload_details(self, value):
        self.details.item.get_filter_w(self.table.item.name + '_id', value['id'])
        self.details.reload(self.details.item.values)

    def reload(self, values=None):
        self.table.reload(values)
        if values:
            self.reload_details(values[0])


class GroupItemTable(ItemTable):
    def __init__(
                self, 
                item_name: str, 
                fields: list = [], 
                values: list = None, 
                buttons=TABLE_BUTTONS, 
                group_id=0,
                current_id=0,
                ):
        super().__init__(
            item_name, 
            '', 
            fields, 
            values, 
            buttons, 
            group_id, 
            True,
            False,
            False,
            )
        self.current_id = current_id
        self.id_name = item_name + '_id'

    def reload(self, values=None):
        if values is None:
            values = self.get_values()
        if self.current_id:
            values = self.filter_childs(values, self.current_id)
        self.table.table.reload(values)

    def filter_childs(self, values, id):
        res = []
        for v in values:
            if v[self.id_name] != id:
                res.append(v)


class ItemTree(QSplitter):
    actionResolved = pyqtSignal()
    def __init__(self, 
                 item_name: str,
                 name: str, 
                 title: str = '', 
                 values: list = None, 
                 fields: list = None, 
                 headers: list = None, 
                 buttons=TABLE_BUTTONS,
                 is_info_bottom = False,
                 show_info=False,
                 ):
        self.item = Item(item_name)
        super().__init__(Qt.Orientation.Vertical if is_info_bottom else Qt.Orientation.Horizontal)
        self.tree = TreeWControls(name, title, values, fields, headers, buttons)
        self.tree.actionInvoked.connect(self.perform_action)
        self.addWidget(self.tree)
        self.current_value = {}
        if not show_info:
            return
        self.info = InfoBlock(self.item.model_w, columns=2 if is_info_bottom else 1)
        self.addWidget(self.info)
        self.tree.tree.itemSelected.connect(self.item_selected)
        self.setStretchFactor(0, 5)
        self.setStretchFactor(1, 2)

    def item_selected(self, value):
        self.current_value = value
        self.info.reload(value)

    def get_values(self):
        err = self.item.get_all_w()
        if err:
            error(err)
            return
        return self.item.values
    
    def reload(self, values=None):
        if values is None:
            values = self.get_values()
        self.tree.tree.reload(values)

    def perform_action(self, action:str, value:dict=None):
        # print('action', action, value)
        if action == 'reload':
            self.reload()
            return
        if action == 'delete':
            self.del_dialog(value)
            return
        title = 'Створити'
        if action == 'copy':
            value['id'] = 0
        elif action == 'create':
            value = self.prepare_value_to_action()
            if value is None:
                return
        elif action == 'edit':
            title = 'Редагувати'
        self.dialog(value, title)

    def prepare_value_to_action(self):
        item = Item(self.item.name)
        item.create_default()
        return item.value

    def dialog(self, value, title):
        i = Item(self.item.name)
        dlg = FormDialog(title, i.model, value)
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
        return value
    
    def del_dialog(self, value):
        cur = self.tree.tree.currentItem()
        if cur.childCount():
            error("Не можна видалити елемент, що має дочірні елементи")
            return
        dlg = DeleteDialog(value, with_prompt=False)
        res = dlg.exec()
        if res:
            i = Item(self.item.name)
            i.value = value
            cause = ''
            if hasattr(dlg, 'ehtry'):
                cause = dlg.entry.text()
            if cause:
                if 'comm' in value:
                    i.value['comm'] = f'del: {cause}' + value['comm']
                    i.save()
                elif 'info' in value:
                    i.value['info'] = f'\nПричина видалення:\n{cause}\n' + value['info']
                    i.save()
                    
            err = i.delete(value['id'], cause)
            if err:
                error(err)
                return
            self.reload()
            self.actionResolved.emit()
        

class GroupTree(ItemTree):
    def __init__(self, 
                 item_name: str, 
                 title: str = '', 
                 values: list = None, 
                 group_id = None,
                 ):
        btns = {'reload':'Оновити', 'create':'Створити', 'edit':'Редагувати', 'delete':'Видалити'}
        super().__init__(item_name, item_name, title, values, buttons=btns)
        top_btn = QPushButton()
        top_btn.setIcon(QIcon(f'images/icons/top.png'))
        top_btn.setToolTip('На верхній рівень')
        top_btn.clicked.connect(lambda _,action='top': self.perform_action(action))
        self.tree.hbox.addWidget(top_btn)
        if group_id is not None:
            self.set_current_id(group_id)
        
    def perform_action(self, action: str, value: dict = None):
        if action != 'top':
            return super().perform_action(action, value)
        i = Item(self.item.name)
        if not value:
            value = self.tree.tree.get_selected_value()
        i.value = value
        i.value[self.tree.tree.key_name] = 0
        err = i.save()
        if err:
            error(err)
            return
        self.reload()
        self.actionResolved.emit()

    def set_current_id(self, group_id):
        self.tree.tree.set_current_id(group_id)

    def value(self):
        return self.tree.value()
    
    def del_dialog(self, value):
        i = Item(self.item.name.split('_')[0])
        err = i.get_filter(self.item.name + '_id', value['id'])
        if err:
            error(err)
            return
        if i.values:
            error("Не можна видалити групу, що має дочірні елементи")
            return
        return super().del_dialog(value)

        