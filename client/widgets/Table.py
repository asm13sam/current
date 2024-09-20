from PyQt6 import QtGui
from PyQt6.QtCore import pyqtSignal, QModelIndex, Qt
from PyQt6.QtGui import QStandardItemModel, QStandardItem
from PyQt6.QtWidgets import (
    QVBoxLayout,
    QHBoxLayout,
    QWidget,
    QAbstractItemView,
    QTableView,
    QHeaderView,
    QPushButton,
    QLineEdit,
    )

from widgets.Dialogs import error
from common.params import SORT_ROLE, FULL_VALUE_ROLE, TABLE_BUTTONS


# field_names specified fields of table to be showed,
# if it is empty - all from data model
class TableModel(QStandardItemModel):
    def __init__(self, data_model: dict, field_names: list = []):
        super().__init__()
        self.data_model = data_model
        self.set_fields(field_names)
        self.setHorizontalHeaderLabels(self.headers)
        self.setSortRole(SORT_ROLE)

    def set_fields(self, field_names: str):
        self.field_names = field_names if field_names else self.data_model.keys()
        self.field_names = [name for name in self.field_names if not name in self.data_model or (not name.endswith('_id') and self.data_model[name]['def'] != [])]
        header_fields = [name[:-1] if name.endswith('2') else name for name in self.field_names]
        self.headers = [self.data_model[name]['hum'] for name in header_fields if name in self.data_model]
        
    def reload(self, values):
        self.clear()
        self.setHorizontalHeaderLabels(self.headers)
        for value in values:
            self.append(value)

    def make_item(self, value, name):
        v = value[name]
        if type(v) == float:
            v = round(v, 2)
        item = QStandardItem(str(v))
        item.setData(v, SORT_ROLE)
        item.setEditable(False)
        return item

    def append(self, value):
        row = []
        for name in self.field_names:
            row.append(self.make_item(value, name))
        row[0].setData(value, FULL_VALUE_ROLE)    
        self.appendRow(row)

    def get_row_value(self, index):
        row = self.item(index, 0)
        if not row:
            return
        row_value = row.data(FULL_VALUE_ROLE)
        return row_value

    def values(self):
        values = []
        for i in range(self.rowCount()):
            row = self.get_row_value(i)
            values.append(row)
        return values

    def calc_sum(self, field: str):
        res = 0
        for i in range(self.rowCount()):
            row = self.get_row_value(i)
            res += row[field]
        return res


# simple table without links to item
# multiselect enabled
# values gets from outer code
# data model example:
# {
#     "id": {"def": 0, "hum": "Номер"},
#     "name": {"def": "", "hum": "Назва"},
#     "full_name": {"def": "", "hum": "Повна назва"},
#     "matherial_group_id": {"def": 0, "hum": "Група"},
#     "measure_id": {"def": 0, "hum": "Од. виміру"},
#     "cost": {"def": 0.0, "hum": "Ціна"},
#     "is_active": {"def": true, "hum": "Діючий"}
# }
class Table(QTableView):
    valueSelected = pyqtSignal(dict)
    valueDoubleCklicked = pyqtSignal(dict)
    tableChanged = pyqtSignal()
    def __init__(self, data_model: dict, table_fields: list=[], values=[]):
        super().__init__()
        self.horizontalHeader().setSectionResizeMode(QHeaderView.ResizeMode.Interactive)
        self.horizontalHeader().setStretchLastSection(True)

        self._model = TableModel(data_model, table_fields)
        self.setModel(self._model)
        self.doubleClicked[QModelIndex].connect(self.value_dblclicked)
        self.setSortingEnabled(True)
        self.setSelectionBehavior(QAbstractItemView.SelectionBehavior.SelectRows)

        if values:
            self.reload(values)

    def set_fields(self, fields):
        self._model.set_fields(fields)

    def add_value(self, value):
        self._model.append(value)
    
    def reload(self, values):
        self._model.reload(values)
        for i in range(self._model.columnCount()):
            self.resizeColumnToContents(i)
        self.sortByColumn(0, Qt.SortOrder.DescendingOrder)
        self.tableChanged.emit()

    def clear(self):
        self._model.clear()
        
    def keyPressEvent(self, e: QtGui.QKeyEvent) -> None:
        if e.key() == Qt.Key.Key_Enter or e.key() == Qt.Key.Key_Return:
            indexes = self.selectedIndexes()
            if indexes:
                self.value_dblclicked(indexes[0])
        super().keyPressEvent(e)

    def recalc(self, field):
        values = self._model.values()
        res = 0
        for v in values:
            res += v[field]
        return res
            
    def delete_values(self):
        while True:
            selected_rows = self.get_selected_rows()
            if not selected_rows:
                return
            self._model.removeRow(selected_rows[0])
    
    def get_selected_ids(self):
        rows = self.get_selected_rows()
        if not rows:
            return []
        ids = [int(self._model.item(row).text()) for row in rows]
        return ids

    def get_selected_values(self):
        rows = self.get_selected_rows()
        if not rows:
            return []
        values = [self._model.get_row_value(row) for row in rows]
        return values
        
    def get_selected_value(self):
        selected_values = self.get_selected_values()
        if not selected_values or len(selected_values) > 1:
            return
        return selected_values[0]

    def get_selected_rows(self) -> list:
        indexes = self.selectedIndexes()
        if not indexes:
            return []
        selected_rows = list(set(index.row() for index in indexes))
        return selected_rows

    def value_selected(self, index):
        pass

    def value_dblclicked(self, index):
        value = self._model.get_row_value(index.row())
        if value:
            self.valueDoubleCklicked.emit(value)

    def currentChanged(self, current, previous) -> None:
        if current == previous:
            return
        value = self._model.get_row_value(current.row())
        if value:
            self.valueSelected.emit(value)

    def values(self):
        return self._model.values()


class TableWControls(QWidget):
    actionInvoked = pyqtSignal(str, dict)
    searchChanged = pyqtSignal(str)
    def __init__(
            self, 
            data_model: dict, 
            table_fields: list = [], 
            values=[],
            buttons = TABLE_BUTTONS,
            with_search=False,
            ):
        super().__init__()
        self.with_search = with_search
        self.box = QVBoxLayout()
        self.box.setContentsMargins(0,0,0,0)
        self.setLayout(self.box)
        self.table = Table(data_model, table_fields, values)
        self.box.addWidget(self.table)
        if buttons:
            self.add_buttons(buttons)
            
    def add_buttons(self, buttons):
        controls = QWidget()
        self.box.insertWidget(0, controls)
        self.hbox = QHBoxLayout()
        self.hbox.setContentsMargins(0,0,0,0)
        controls.setLayout(self.hbox)
        if self.with_search:
            self.search_entry = QLineEdit()
            self.search_entry.setPlaceholderText('Пошук')
            self.search_entry.textChanged.connect(self.search_changed)
            self.search_entry.setMinimumWidth(80)
            self.search_entry.setMaximumWidth(150)
            self.hbox.addWidget(self.search_entry)
        self.hbox.addStretch()
        for b in TABLE_BUTTONS:
            if b in buttons:
                btn = QPushButton()
                btn.setIcon(QtGui.QIcon(f'images/icons/{b}.png'))
                btn.setToolTip(TABLE_BUTTONS[b])
                btn.clicked.connect(lambda _,action=b: self.action(action))
                self.hbox.addWidget(btn)

    def action(self, action):
        if action == 'create' or action == 'reload':
            self.actionInvoked.emit(action, {})
        else:
            value = self.table.get_selected_value()
            if not value:
                return
            self.actionInvoked.emit(action, value)
    
    def search_changed(self, text:str):
        self.searchChanged.emit(text)

    def values(self):
        return self.table.values()
        

class EditTableModel(TableModel):
    itemChanged = pyqtSignal(tuple)
    def __init__(self, data_model: dict, field_names: list = [], edit_fields: list = []):
        super().__init__(data_model, field_names)
        self.edit_fields = edit_fields

        self.dataChanged.connect(self.data_changed)

    def data_changed(self, topLeft:QModelIndex, bottomRight:QModelIndex):
        row, column = topLeft.row(), topLeft.column()
        item = self.item(row, column)
        item0 = self.item(row, 0)
        row_value = item0.data(FULL_VALUE_ROLE)
        t = type(item.data(SORT_ROLE))
        if t == bool:
            data = item.checkState() == Qt.CheckState.Checked
        elif t == int or t == float:
            try:
                if t == int:
                    data = int(item.text())
                else:
                    data = float(item.text())    
            except ValueError:
                error('Некорректне значення')
                item.setText(str(item.data(SORT_ROLE)))
                return    
        else:    
            data = item.text()
        
        row_value[self.field_names[column]] = data
        item.setData(data, SORT_ROLE)
        item0.setData(row_value, FULL_VALUE_ROLE)
        self.itemChanged.emit((self.field_names[column], data, row))

    def make_item(self, value, name):
        v = value[name]
        if type(v) == float:
            v = round(v, 2)
        item = QStandardItem(str(v))
        item.setData(v, SORT_ROLE)
        if name in self.edit_fields:
            if type(v) == bool:
                item.setCheckable(True)
                item.setEditable(False)
                item.setCheckState(Qt.CheckState.Checked if v else Qt.CheckState.Unchecked)
            else:
                item.setEditable(True)
        else:
            item.setEditable(False)
        return item

    
class EditTable(QWidget):
    itemEdited = pyqtSignal(tuple)
    def __init__(self, data_model: dict, table_fields: list = [], edit_fields: list = []):
        super().__init__()
        self.box = QVBoxLayout()
        self.setLayout(self.box)
        self.table = Table(data_model, table_fields, values=[])
        self.box.addWidget(self.table)
            
        controls = QWidget()
        self.box.insertWidget(0, controls)
        self.hbox = QHBoxLayout()
        controls.setLayout(self.hbox)
        self.hbox.addStretch()
        del_btn = QPushButton("Видалити")
        del_btn.clicked.connect(self.delete_row)
        self.hbox.addWidget(del_btn)

        self.table._model = EditTableModel(data_model, table_fields, edit_fields)
        self.table.setModel(self.table._model)
        self.table._model.itemChanged.connect(lambda changes: self.itemEdited.emit(changes))
        
    def add_value(self, value):
        self.table._model.append(value)

    def delete_row(self):
        self.table.delete_values()

    def values(self):
        return self.table.values()    
    