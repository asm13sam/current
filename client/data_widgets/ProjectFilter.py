from PyQt6.QtCore import pyqtSignal
from PyQt6.QtWidgets import (
    QVBoxLayout,
    QWidget,
    QLabel,
    QLineEdit,
    QPushButton,
    QTabWidget,
    QHBoxLayout,
    )

from data.model import Item, ProjectItem
from widgets.ButtonsBlock import ButtonsBlock
from widgets.ComboBoxSelector import ComboBoxSelector
from widgets.SearchCBListSelector import SearchCBListSelector
from widgets.Dialogs import error
from widgets.Form import SingleSelectDialog
from common.params import MIN_SEARCH_STR_LEN


class ByProjectFilter(QWidget):
    statusChanged = pyqtSignal(dict)
    groupChanged = pyqtSignal(dict)
    idEntered = pyqtSignal(int)
    searchChanged = pyqtSignal(str)
    
    def __init__(self, title: str='', status_values: list=[], group_values: list=[]):
        super().__init__()
        layout = QVBoxLayout()
        self.setLayout(layout)
        layout.setContentsMargins(0, 0, 0, 0)
        if title:
            layout.addWidget(QLabel(title))
        
        self.project_status = ButtonsBlock(values=status_values)
        layout.addWidget(self.project_status)
        self.project_status.buttonClicked.connect(self.project_status_changed)

        self.project_group = ComboBoxSelector('Група', group_values)
        layout.addWidget(self.project_group)
        self.project_group.selectionChanged.connect(self.project_group_changed)

        self.project_id = QLineEdit()
        self.project_id.setPlaceholderText('За номером')
        layout.addWidget(self.project_id)
        self.project_id.returnPressed.connect(self.project_id_entered)

        self.project_search = QLineEdit()
        self.project_search.setPlaceholderText('Пошук')
        layout.addWidget(self.project_search)
        self.project_search.textChanged.connect(self.project_search_changed)
        layout.addStretch()

    def project_status_changed(self, value):
        self.statusChanged.emit(value)
        
    def project_group_changed(self, value):
        if not value['id']:
            return
        self.groupChanged.emit(value)

    def project_id_entered(self, id=0):
        if not id:
            text = self.project_id.text()
            try:
                id = int(text)
            except:
                return
        self.idEntered.emit(id)
        
    def project_search_changed(self, text):
        if len(text) < MIN_SEARCH_STR_LEN:
            return
        self.searchChanged.emit(text)


class ByContragentFilter(QWidget):
    contragentChanged = pyqtSignal(dict)
    actionInvoked = pyqtSignal(str, dict)
    contactsRequired = pyqtSignal(dict)
    def __init__(self, title: str='', group_values: list=[]):
        super().__init__()
        self.contragent = Item('contragent')
        self.box = QVBoxLayout()
        self.setLayout(self.box)
        self.box.setContentsMargins(0, 0, 0, 0)
        self.current_contragent_value = {}

        if title:
            self.box.addWidget(QLabel(title))

        ext_box = QHBoxLayout()
        self.box.addLayout(ext_box)
        from_table = QPushButton('З таблиці')
        ext_box.addWidget(from_table)
        from_table.clicked.connect(self.get_from_table)
        reload = QPushButton('Оновити')
        ext_box.addWidget(reload)
        reload.clicked.connect(self.refresh)
        self.contragent_id = QLineEdit()
        self.contragent_id.setPlaceholderText('За номером')
        self.box.addWidget(self.contragent_id)
        self.contragent_id.returnPressed.connect(self.contragent_id_entered)
        
        self.contragent_selector = SearchCBListSelector(cb_title='Група', cb_values=group_values)
        self.box.addWidget(self.contragent_selector, 10)
        self.contragent_selector.cbSelectionChanged.connect(self.contragent_group_changed)
        self.contragent_selector.searchStringChanged.connect(self.contragent_search_changed)
        self.contragent_selector.selectionChanged.connect(self.contragent_changed)
        self.append_widget()
        
    def append_widget(self):    
        cont_table = QPushButton('Контакти')
        self.box.addWidget(cont_table)
        cont_table.clicked.connect(self.contacts_required)
        self.edit_btn = QPushButton('Редагувати')
        self.box.addWidget(self.edit_btn)
        self.edit_btn.clicked.connect(lambda: self.action_invoked('edit'))
        
    def contragent_group_changed(self, value):
        if not value['id']:
            return
        err = self.contragent.get_filter_w('contragent_group_id', value['id'])
        if err:
            return
        self.contragent_selector.reload(self.contragent.values)

    def contragent_id_entered(self):
        text = self.contragent_id.text()
        try:
            id = int(text)
        except:
            return
        self.set_contragent_by_id(id)
        
    def set_contragent_by_id(self, id:int):
        err = self.contragent.get_w(id)
        if err:
            return
        self.contragent_selector.reload([self.contragent.value])
        self.contragentChanged.emit(self.contragent.value)
        
    def contragent_search_changed(self, text):
        if len(text) < MIN_SEARCH_STR_LEN:
            return
        err = self.contragent.get_find_w(self.contragent.find[0], text.lower())
        if err:
            return
        self.contragent_selector.reload(self.contragent.values)

    def contragent_changed(self, value):
        self.contragentChanged.emit(value)
        
    def get_from_table(self):
        dlg = SingleSelectDialog('contragent')
        res = dlg.exec()
        if not res:
            return
        v = dlg.value
        if v:
            self.contragent_selector.reload([v])
            self.contragentChanged.emit(v)

    def action_invoked(self, action:str):
        value = self.contragent_selector.get_current_value()
        if value or action == 'create':
            self.actionInvoked.emit(action, value)

    def contacts_required(self):
        value = self.contragent_selector.get_current_value()
        if value:
            self.contactsRequired.emit(value)

    def refresh(self):
        self.contragent_selector.reload()


class ProjectFilter(QTabWidget):
    valuesChanged = pyqtSignal(list)
    valuesReloaded = pyqtSignal(list)
    def __init__(self) -> None:
        super().__init__()
        self.project = ProjectItem()
        self.current_filter = {}
        
        self.project_status = Item('project_status')
        err = self.project_status.get_all_w()
        if err:
            return
        self.project_group = Item('project_group')
        err = self.project_group.get_all_w()
        if err:
            return
        
        self.project_filter = ByProjectFilter(
            group_values=self.project_group.values, 
            status_values=self.project_status.values,
            )
        self.addTab(self.project_filter, 'Проект')
        self.project_filter.statusChanged.connect(self.project_status_changed)
        self.project_filter.groupChanged.connect(self.project_group_changed)
        self.project_filter.idEntered.connect(self.project_id_entered)
        self.project_filter.searchChanged.connect(self.project_search_changed)
        
        self.contragent_group = Item('contragent_group')
        err = self.contragent_group.get_all_w()
        if err:
            error('On get contragent groups: ', err)
            return
        self.contragent_filter = ByContragentFilter(group_values=self.contragent_group.values)
        self.addTab(self.contragent_filter, 'Контрагент')
        self.contragent_filter.contragentChanged.connect(self.contragent_changed)


    def project_status_changed(self, value):
        if not value['id']:
            return
        self.current_filter = {'field': 'project_status_id', 'id':value['id']}
        err = self.project.get_filter_w('project_status_id', value['id'])
        if err:
            return
        self.values_changed(self.project.values)

    def project_group_changed(self, value):
        if not value['id']:
            return
        self.current_filter = {'field': 'project_group_id', 'id':value['id']}
        err = self.project.get_filter_w('project_group_id', value['id'])
        if err:
            return
        self.values_changed(self.project.values)

    def project_id_entered(self, id=0):
        if not id:
            text = self.project_filter.project_id.text()
            try:
                id = int(text)
            except:
                return
        self.current_filter = {'id': id}
        err = self.project.get_w(id)
        if err:
            return
        self.values_changed([self.project.value])
        
    def project_search_changed(self, text):
        if len(text) < MIN_SEARCH_STR_LEN:
            return
        self.current_filter = {'find': self.project.find[0], 'text':text}
        err = self.project.get_find_w(self.project.find[0], text)
        if err:
            return
        self.valuesReloaded.emit(self.project.values)

    def contragent_changed(self, value):
        self.current_filter = {'field': 'contragent_id', 'id':value['id']}
        err = self.project.get_filter_w('contragent_id', value['id'])
        if err:
            return
        self.values_changed(self.project.values)
    
    def reload(self):
        if 'find' in self.current_filter:
            err = self.project.get_find_w(self.current_filter['find'], self.current_filter['text'])
        elif 'field' in self.current_filter:
            err = self.project.get_filter_w(self.current_filter['field'], self.current_filter['id'])
        else:
            err = self.project.get_w(self.current_filter['id'])
        if err:
            return
        self.valuesReloaded.emit(self.project.values)

    
    def values_changed(self, values: list):
        self.valuesChanged.emit(values)
    