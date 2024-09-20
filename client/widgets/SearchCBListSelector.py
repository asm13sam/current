from PyQt6.QtCore import pyqtSignal

from widgets.SearchListSelector import SearchListSelector
from widgets.ComboBoxSelector import ComboBoxSelector

# values is list of value dict with at least two fields 'name'
class SearchCBListSelector(SearchListSelector):
    cbSelectionChanged = pyqtSignal(dict)
    def __init__(self, title: str='', cb_values: list=[], cb_title:str=''):
        super().__init__()
        top_widget = 0
        if title:
            top_widget = 1     
        self.cb_selector = ComboBoxSelector(cb_title, values=cb_values)
        self.main_layout.insertWidget(top_widget, self.cb_selector)
        self.cb_selector.selectionChanged.connect(self.cb_selection_changed)

    def cb_selection_changed(self, value):
        self.cbSelectionChanged.emit(value)


if __name__ == '__main__':
    import sys
    from PyQt6.QtWidgets import QApplication, QVBoxLayout, QWidget
    
    sl_values = [
        {'name': 'abcName A', 'id': 1, 'group_id': 1},
        {'name': 'defName B', 'id': 2, 'group_id': 2},
        {'name': 'ghiName C', 'id': 3, 'group_id': 3},
        {'name': 'klmName D', 'id': 4, 'group_id': 1},
        {'name': 'nopName E', 'id': 5, 'group_id': 2}, 
        {'name': 'Qwerty111', 'id': 11, 'is_active': True, 'group_id': 1},
        {'name': 'Asdfg222', 'id': 12, 'is_active': False, 'group_id': 2},
        {'name': 'Zxcvb333', 'id': 13, 'is_active': True, 'group_id': 3},
    ]

    cb_values = [
        {'name': 'Group1', 'id': 1},
        {'name': 'Group2', 'id': 2},
        {'name': 'Group3', 'id': 3},
    ]

    qt_app = QApplication(sys.argv)
    window = QWidget()
    layout = QVBoxLayout()
    window.setLayout(layout)
    
    sl = SearchCBListSelector('TestBlock', cb_values=cb_values)
    layout.addWidget(sl)
    sl.selectionChanged.connect(lambda value: print(value))
    layout.addStretch()
    
    def search_changed(text):
        values = [v for v in sl_values if text.lower() in v['name'].lower()]
        sl.reload(values)
    sl.searchStringChanged.connect(search_changed)

    def cb_changed(value):
        values = [v for v in sl_values if value['id'] == v['group_id']]
        sl.reload(values)
    sl.cbSelectionChanged.connect(cb_changed)

    window.show()
    sys.exit(qt_app.exec())