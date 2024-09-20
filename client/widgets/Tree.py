import sys
from PyQt6.QtWidgets import (
    QWidget, 
    QTreeWidget, 
    QHeaderView, 
    QTreeWidgetItem, 
    QTreeWidgetItemIterator, 
    QStyle,
    QVBoxLayout,
    QHBoxLayout,
    QPushButton,
    )
from PyQt6.QtCore import Qt, pyqtSignal
from PyQt6.QtGui import QBrush, QColor, QIcon
from widgets.Dialogs import error

if __name__ != '__main__':
    from common.params import FULL_VALUE_ROLE, TABLE_BUTTONS, IS_EXT_ROLE, ID_ROLE

class Tree(QTreeWidget):
    valueDoubleCklicked = pyqtSignal(dict)
    itemSelected = pyqtSignal(dict)
    def __init__(self, name: str, title: str='', values: list=None, fields: list=None, headers: list=None):
        super().__init__()
        self.fields = fields
        self.setColumnCount(1 if not fields else len(fields))
        if not fields:
            self.setHeaderLabels([title])
        else:
            self.setHeaderLabels(headers)
        header = self.header()
        header.setSectionResizeMode(QHeaderView.ResizeMode.ResizeToContents)
        self.currentItemChanged.connect(self.cur_changed)
        self.itemDoubleClicked.connect(self.value_dblclicked)
        self.name = name
        self.key_name = name + '_id'
        self.dataset = {}
        if values:
            self.reload(values)

    def reload(self, values: list=[]):
        self.clear()
        self.dataset = {}
        if not values:
            return
        
        for v in values:
            if v[self.key_name] not in self.dataset:
                self.dataset[v[self.key_name]] = []    
            self.dataset[v[self.key_name]].append(v)
        self.add_childs(0)
    
    def add_childs(self, group_id, parent_item=None):
        for td in self.dataset[group_id]:
            if parent_item is None:
                parent_item = self.invisibleRootItem()
            data_item = QTreeWidgetItem()
            if not self.fields:
                data_item.setText(0, td['name'])
            else:
                for i, f in enumerate(self.fields):
                    data_item.setText(i, str(td[f]))
            data_item.setData(1, FULL_VALUE_ROLE, td)
            if parent_item is not None:
                parent_item.addChild(data_item)
            if 'type' in td and td['type'] != self.name or td['id'] not in self.dataset:
                continue    
            self.add_childs(td['id'], data_item)

    def add_value(self, value):
        self.clear()
        if value[self.key_name] not in self.dataset:
            self.dataset[value[self.key_name]] = []    
        self.dataset[value[self.key_name]].append(value)
        self.add_childs(0)

    def add_values(self, values):
        self.clear()
        for value in values:
            if value[self.key_name] not in self.dataset:
                self.dataset[value[self.key_name]] = []    
            self.dataset[value[self.key_name]].append(value)
        self.add_childs(0)

    def cur_changed(self, current, previous):
        if not current:
            return
        value = current.data(1, FULL_VALUE_ROLE)
        self.itemSelected.emit(value)

    def delete_current(self, id=None):
        cur = self.currentItem()
        if id is None:
            id = cur.data(1, FULL_VALUE_ROLE)['id']
        root = self.invisibleRootItem()
        (cur.parent() or root).removeChild(cur)
        for i, v in enumerate(self.dataset[0]):
            if v['id'] == id:
                break
        self.dataset[0].pop(i)
        
    def value(self):
        i = self.currentItem()
        if not i:
            return
        value = i.data(1, FULL_VALUE_ROLE)
        return value
    
    def set_dblclick_cb(self, cb):
        self.valueDoubleCklicked.connect(cb)
        
    def remove_dblclick_cb(self):
        try: 
            self.valueDoubleCklicked.disconnect()
        except Exception: 
            pass

    def value_dblclicked(self, index):
        value = self.currentItem().data(1, FULL_VALUE_ROLE)
        if value:
            self.valueDoubleCklicked.emit(value)

    def item_by_key(self, key, value):
        iterator = QTreeWidgetItemIterator(self)
        while iterator.value():
            item = iterator.value()
            if item.data(1, FULL_VALUE_ROLE)[key] == value:
                return item
            iterator += 1

    def set_current_id(self, id: int):
        item = self.item_by_key('id', id)
        self.setCurrentItem(item)
        self.scrollToItem(item)

    def delete_by_id(self, id):
        item = self.item_by_key('id', id)
        if not item:
            return
        root = self.invisibleRootItem()
        (item.parent() or root).removeChild(item)
        for i, v in enumerate(self.dataset[0]):
            if v['id'] == id:
                break
        self.dataset[0].pop(i)

    def get_selected_value(self):
        value = self.value()
        if not value:
            error('Оберіть один елемент')
            return
        return value
        
    
class ExtTree(Tree):
    def __init__(self, name: str, title: str = '', values: list = None, fields: list = None, headers: list = None):
        super().__init__(name, title, values, fields, headers)

    def reload(self, values: list, append_values: list, key_name: str):
        super().reload(values)
        if not append_values:
            return
        self.append_dataset = {}
        key_name += '_id'
        for v in append_values:
            if v[key_name] not in self.append_dataset:
                self.append_dataset[v[key_name]] = []    
            self.append_dataset[v[key_name]].append(v)
        self.append_childs()

    def append_childs(self):
        pixmapi = QStyle.StandardPixmap.SP_FileDialogListView
        icon = self.style().standardIcon(pixmapi)
                
        parents = {}
        for key in self.append_dataset:
            parent_item = self.item_by_key('id', key)
            parents[key] = parent_item
        for key, parent_item in parents.items():
            for v in self.append_dataset[key]:
                data_item = QTreeWidgetItem(parent_item)
                if not self.fields:
                    title = v['short_name'] if 'short_name' in v else v['name']
                    data_item.setText(0, title)
                else:
                    for i, f in enumerate(self.fields):
                        data_item.setText(i, str(v[f]))
                data_item.setData(1, FULL_VALUE_ROLE, v)
                data_item.setData(1, IS_EXT_ROLE, True)
                # Linux
                # data_item.setIcon(0, QIcon.fromTheme("list-add"))
                data_item.setIcon(0, icon)

    def value_dblclicked(self, index):
        is_ext = self.currentItem().data(1, IS_EXT_ROLE)
        if not is_ext:
            return 
        value = self.currentItem().data(1, FULL_VALUE_ROLE)
        if value:
            self.valueDoubleCklicked.emit(value)


class TreeWControls(QWidget):
    actionInvoked = pyqtSignal(str, dict)

    def __init__(self, 
                 name: str, 
                 title: str = '', 
                 values: list = None, 
                 fields: list = None, 
                 headers: list = None, 
                 buttons = TABLE_BUTTONS,
                 ):
        super().__init__()
        
        self.box = QVBoxLayout()
        self.box.setContentsMargins(0,0,0,0)
        self.setLayout(self.box)
        self.tree = Tree(name, title, values, fields, headers)
        self.box.addWidget(self.tree)
        if buttons:
            self.add_buttons(buttons)
            
    def add_buttons(self, buttons):
        controls = QWidget()
        self.box.insertWidget(0, controls)
        self.hbox = QHBoxLayout()
        self.hbox.setContentsMargins(0,0,0,0)
        controls.setLayout(self.hbox)
        self.hbox.addStretch()
        for b in TABLE_BUTTONS:
            if b in buttons:
                btn = QPushButton()
                btn.setIcon(QIcon(f'images/icons/{b}.png'))
                btn.setToolTip(TABLE_BUTTONS[b])
                btn.clicked.connect(lambda _,action=b: self.action(action))
                self.hbox.addWidget(btn)
        
    def action(self, action):
        if action == 'create' or action == 'reload':
            self.actionInvoked.emit(action, {})
        else:
            value = self.tree.get_selected_value()
            if not value:
                return
            self.actionInvoked.emit(action, value)

    def value(self):
        return self.tree.value()
    
    def set_dblclick_cb(self, cb):
        self.tree.set_dblclick_cb(cb)
        
    def remove_dblclick_cb(self):
        self.tree.remove_dblclick_cb()


class DatasetItem:
    def __init__(self, value: dict, item_type: str, type_hum: str) -> None:
        self.type = item_type
        self.type_hum = type_hum
        self.value = value


class Dataset:
    def __init__(self, tree_key_name: str, root_key_value: int|str=0) -> None:
        self.data: dict[int|str, list[DatasetItem]] = {}
        self.tree_key_name = tree_key_name
        self.root_key_value = root_key_value

    def add_item(self, item: DatasetItem):
        tree_key = item.value[self.tree_key_name]
        if tree_key not in self.data:
            self.data[tree_key] = []
        self.data[tree_key].append(item)

    def get(self, tree_key_value, item_id, item_type) -> DatasetItem | None:
        for item in self.data[tree_key_value]:
            if item.value['id'] == item_id and item.type == item_type:
                return item
            
    def delete(self, tree_key_value, item_id, item_type) -> DatasetItem | None:
        for i, item in enumerate(self.data[tree_key_value]):
            if item.value['id'] == item_id and item.type == item_type:
                return item

    def clear(self):
        self.data = {}
        

class DatasetTree(QTreeWidget):
    valueDoubleCklicked = pyqtSignal(dict)
    itemSelected = pyqtSignal(dict)
    def __init__(self, title: str='', fields: list=None, headers: list=None):
        super().__init__()
        self.fields = fields
        self.setColumnCount(1 if not fields else len(fields))
        if not fields:
            self.setHeaderLabels([title])
        else:
            self.setHeaderLabels(headers)
        header = self.header()
        header.setSectionResizeMode(QHeaderView.ResizeMode.ResizeToContents)
        self.currentItemChanged.connect(self.current_changed)
        self.itemDoubleClicked.connect(self.value_dblclicked)
        self.dataset: Dataset = None
        
    def reload(self, dataset: Dataset=None):
        self.clear()
        if dataset is None:
            return
        self.dataset = dataset
        self.add_childs(self.dataset.root_key_value, self.invisibleRootItem())
    
    def add_childs(self, current_tree_key_value: int|str, parent_item: QTreeWidgetItem):
        for item in self.dataset.data[current_tree_key_value]:
            item_widget = QTreeWidgetItem()
            if not self.fields:
                item_widget.setText(0, item.value['name'])
            else:
                for i, f in enumerate(self.fields):
                    item_widget.setText(i, str(item.value[f]))
            tree_key = self.dataset.tree_key_name
            item_widget.setData(1, ID_ROLE, (item.value[tree_key], item.value['id'], item.type))
            parent_item.addChild(item_widget)
            
            if tree_key.endswith['_id']:
                tree_key = tree_key[:3]
            if item.type != tree_key or item.value['id'] not in self.dataset.data:
                continue    
            self.add_childs(item.value['id'], item_widget)

    def _add_value(self, value, item_type, type_hum):
        tree_key = self.dataset.tree_key_name
        if value[tree_key] not in self.dataset.data:
            self.dataset.data[value[tree_key]] = []    
        self.dataset.data[value[tree_key]].append(DatasetItem(value, item_type, type_hum))

    def add_value(self, value, item_type, type_hum):
        self.clear()
        self._add_value(value, item_type, type_hum)
        self.add_childs(self.dataset.root_key_value)

    def add_values(self, values, item_type, type_hum):
        self.clear()
        for value in values:
            self._add_value(value, item_type, type_hum)
        self.add_childs(self.dataset.root_key_value)

    def current_changed(self, current, previous):
        if not current:
            return
        current_dataset_item = self.dataset.get(current.data(1, ID_ROLE))
        self.itemSelected.emit(current_dataset_item.value)

    def value(self):
        i = self.currentItem()
        if not i:
            return
        value = i.data(1, FULL_VALUE_ROLE)
        return value
    
    def set_dblclick_cb(self, cb):
        self.valueDoubleCklicked.connect(cb)
        
    def remove_dblclick_cb(self):
        try: 
            self.valueDoubleCklicked.disconnect()
        except Exception: 
            pass
        
    def value_dblclicked(self, index):
        value = self.currentItem().data(1, FULL_VALUE_ROLE)
        if value:
            self.valueDoubleCklicked.emit(value)

    def item_by_key(self, key, value):
        iterator = QTreeWidgetItemIterator(self)
        while iterator.value():
            item = iterator.value()
            if item.data(1, FULL_VALUE_ROLE)[key] == value:
                return item
            iterator += 1

    def set_current_id(self, id: int):
        item = self.item_by_key('id', id)
        self.setCurrentItem(item)
        self.scrollToItem(item)

    def delete_by_id(self, id):
        item = self.item_by_key('id', id)
        if not item:
            return
        root = self.invisibleRootItem()
        (item.parent() or root).removeChild(item)
        for i, v in enumerate(self.dataset[0]):
            if v['id'] == id:
                break
        self.dataset[0].pop(i)

    def get_selected_value(self):
        value = self.value()
        if not value:
            error('Оберіть один елемент')
            return
        return value


class TreeItem:
    def __init__(self, value: dict, item_type: str, type_hum: str, name_field: str) -> None:
        self.type = item_type
        self.type_hum = type_hum
        self.value = value
        self.name_field = name_field

    def __str__(self) -> str:
        return str(self.value)


class ItemsTreeDataset:
    def __init__(self, tree_key_name: str) -> None:
        self.tree_key_name = tree_key_name
        self.data = {}
        
    def add_item(self, item: TreeItem):
        self.data[(item.value[self.tree_key_name], item.value['id'], item.type)] = item

    def get(self, tree_key_value, item_id, item_type) -> TreeItem | None:
        return self.data.get((tree_key_value, item_id, item_type))
            
    def delete(self, tree_key_value, item_id, item_type) -> TreeItem | None:
        return self.data.pop((tree_key_value, item_id, item_type), None)
    
    def get_by_key(self, key_value):
        return (item for key, item in self.data.items() if key[0] == key_value)
    
    def clear(self):
        self.data = {}

    def __str__(self) -> str:
        res = ''
        for k, v in self.data.items():
            res += f'{k} {v}\n'
        return res


class DatasetItemsTree(QTreeWidget):
    valueDoubleCklicked = pyqtSignal(dict)
    itemSelected = pyqtSignal(TreeItem)
    def __init__(self, title: str='', fields: list=None, headers: list=None):
        super().__init__()
        self.fields = fields
        self.setColumnCount(1 if not fields else len(fields))
        if not fields:
            self.setHeaderLabels([title])
        else:
            self.setHeaderLabels(headers)
        header = self.header()
        header.setSectionResizeMode(QHeaderView.ResizeMode.ResizeToContents)
        self.currentItemChanged.connect(self.current_changed)
        self.dataset: ItemsTreeDataset = None
        
    def reload(self, dataset: ItemsTreeDataset=None):
        self.clear()
        if dataset is None:
            return
        self.dataset = dataset
        self.add_childs(0, self.invisibleRootItem())
    
    def add_childs(self, current_tree_key_value: int|str, parent_item: QTreeWidgetItem):
        for item in self.dataset.get_by_key(current_tree_key_value):
            item_widget = QTreeWidgetItem()
            if not self.fields:
                item_widget.setText(0, item.value[item.name_field])
            else:
                for i, f in enumerate(self.fields):
                    item_widget.setText(i, str(item.value[f]))
            tree_key = self.dataset.tree_key_name
            dataset_key = (item.value[tree_key], item.value['id'], item.type)
            item_widget.setData(1, ID_ROLE, dataset_key)
            parent_item.addChild(item_widget)
            
            if tree_key.endswith('_id'):
                tree_key = tree_key[:-3]
            if item.type != tree_key or dataset_key not in self.dataset.data:
                continue    
            self.add_childs(item.value['id'], item_widget)

    def add_value(self, value, item_type, type_hum):
        self.clear()
        self.dataset.add_item(TreeItem(value, item_type, type_hum))
        self.add_childs(0, self.invisibleRootItem())

    def add_values(self, values, item_type, type_hum):
        self.clear()
        for value in values:
            self.dataset.add_item(TreeItem(value, item_type, type_hum))
        self.add_childs(0, self.invisibleRootItem())

    def current_changed(self, current, previous):
        if not current:
            return
        current_dataset_item = self.dataset.get(*current.data(1, ID_ROLE))
        self.itemSelected.emit(current_dataset_item)


if __name__ == '__main__':
    
    import sys
    from PyQt6.QtWidgets import QApplication

    FULL_VALUE_ROLE = 177

    nomenclature = [
        {'id': 1, 'name': 'group1', 'nomenclature_id': 13},
        {'id': 2, 'name': 'group2', 'nomenclature_id': 1},
        {'id': 3, 'name': 'group3', 'nomenclature_id': 0},
        {'id': 4, 'name': 'group4', 'nomenclature_id': 7},
        {'id': 5, 'name': 'group5', 'nomenclature_id': 0},
        {'id': 6, 'name': 'group6', 'nomenclature_id': 3},
        {'id': 7, 'name': 'group7', 'nomenclature_id': 3},
        {'id': 8, 'name': 'group8', 'nomenclature_id': 5},
        {'id': 9, 'name': 'group9', 'nomenclature_id': 8},
        {'id': 10, 'name': 'group10', 'nomenclature_id': 8},
        {'id': 11, 'name': 'group11', 'nomenclature_id': 5},
        {'id': 12, 'name': 'group12', 'nomenclature_id': 5},
        {'id': 13, 'name': 'group13', 'nomenclature_id': 0},
    ]

    app = QApplication(sys.argv)
    window = Tree('nomenclature', 'Номенклатура')
    window.reload(nomenclature)
    window.itemSelected.connect(lambda v: print(f'Name: {v["name"]} ID: {v["id"]}'))
    window.show()
    sys.exit(app.exec())