from PyQt6.QtCore import pyqtSignal
from PyQt6.QtWidgets import (
    QVBoxLayout,
    QWidget,
    QLabel,
    QLineEdit,
    QListWidget,
    QListWidgetItem,
    )

from common.params import  FULL_VALUE_ROLE, MIN_SEARCH_STR_LEN
# FULL_VALUE_ROLE = 121
# MIN_SEARCH_STR_LEN = 3

# values is list of value dict with at least two fields 'name'
class SearchListSelector(QWidget):
    searchStringChanged = pyqtSignal(str)
    selectionChanged = pyqtSignal(dict)
    valueDoubleClicked = pyqtSignal(dict)
    def __init__(self, title: str=''):
        super().__init__()
        self.main_layout = QVBoxLayout()
        self.setLayout(self.main_layout)
        self.main_layout.setContentsMargins(0, 0, 0, 0)
        if title:
            self.main_layout.addWidget(QLabel(title))
        self.search_entry = QLineEdit()
        self.search_entry.setPlaceholderText('Пошук')
        self.main_layout.addWidget(self.search_entry)
        self.search_entry.textChanged.connect(self.search_entry_changed)
        self.list = QListWidget()
        self.main_layout.addWidget(self.list, stretch=10)
        self.list.currentItemChanged.connect(self.list_selection_changed)
        self.list.itemDoubleClicked.connect(self.value_double_clicked)

    def reload(self, values: list=None):
        if values is None:
            current_value = self.get_current_value()
            if not current_value:
                return
            values = [current_value,]
        self.list.clear()
        for v in values:
            fi = QListWidgetItem(v['name'], self.list)
            fi.setData(FULL_VALUE_ROLE, v)
        
    def search_entry_changed(self, text):
        if len(text) < MIN_SEARCH_STR_LEN:
            self.list.clear()
            return
        self.searchStringChanged.emit(text)

    def list_selection_changed(self, current, previous):
        if current is None:
            return
        value = current.data(FULL_VALUE_ROLE)
        self.selectionChanged.emit(value)

    def value_double_clicked(self, current):
        if current is None:
            return
        value = current.data(FULL_VALUE_ROLE)
        self.valueDoubleClicked.emit(value)

    def get_current_value(self):
        current = self.list.currentItem()
        if not current:
            return {}
        value = current.data(FULL_VALUE_ROLE)
        return value



if __name__ == '__main__':
    import sys
    from PyQt6.QtWidgets import QApplication
    
    sl_values = [
        {'name': 'abcName A', 'id': 1},
        {'name': 'defName B', 'id': 2},
        {'name': 'ghiName C', 'id': 3},
        {'name': 'klmName D', 'id': 4},
        {'name': 'nopName E', 'id': 5}, 
        {'name': 'Qwerty111', 'id': 11, 'is_active': True},
        {'name': 'Asdfg222', 'id': 12, 'is_active': False},
        {'name': 'Zxcvb333', 'id': 13, 'is_active': True},
    ]

    qt_app = QApplication(sys.argv)
    window = QWidget()
    layout = QVBoxLayout()
    window.setLayout(layout)
    
    sl = SearchListSelector('TestBlock')
    layout.addWidget(sl)
    sl.selectionChanged.connect(lambda value: print(value))
    layout.addStretch()
    
    def search_changed(text):
        values = [v for v in sl_values if text.lower() in v['name'].lower()]
        sl.reload(values)
    sl.searchStringChanged.connect(search_changed)

    window.show()
    sys.exit(qt_app.exec())