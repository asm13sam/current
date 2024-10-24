import os
import subprocess

from PyQt6.QtCore import Qt, pyqtSignal
from PyQt6.QtWidgets import (
    QSplitter,
    QPushButton,
    QWidget,
    QVBoxLayout,
    QLabel,
)

from data.model import Item
from widgets.Dialogs import error
from data.app import App
from data_widgets.ProjectFilter import ProjectFilter
from widgets.Table import TableWControls
from widgets.Dialogs import CustomDialog, ok_cansel_dlg, messbox
from widgets.Form import ContactSelectDialog, FormDialog, CustomForm


class ContragentCreateFormDialog(CustomDialog):
    def __init__(self, contragent: Item, contact: Item) -> None:
        self.contact_value = {}
        self.value = {}
        self.box = QVBoxLayout()
        self.forms = QWidget()
        self.forms.setLayout(self.box)
        self.box.addWidget(QLabel('Контрагент'))
        self.contragent = contragent
        self.contragent_form = CustomForm(
            data_model=self.contragent.model,
            fields=['name', 'contragent_group_id', 'phone', 'email', 'comm', 'dir_name'],
            value=self.contragent.value,
            )
        self.box.addWidget(self.contragent_form)
        self.contragent_form.hide_save_btn()
        self.box.addWidget(QLabel('Контакт'))
        self.contact = contact
        self.contact_form = CustomForm(
            data_model=self.contact.model,
            fields=['name', 'phone', 'email'],
            value=self.contact.value,
            )
        self.box.addWidget(self.contact_form)
        
        super().__init__(self.forms, 'Створити контрагента і контакт')
        self.contact_form.saveRequested.connect(self.get_contact_value)
        self.contragent_form.saveRequested.connect(self.get_contragent_value)

    def get_contact_value(self, value):
        self.contact_value = value
        self.contragent_form.get_value()
        
    def get_contragent_value(self, value):
        self.value = value
        self.contragent_form.set_changed(False)

    def accept(self) -> None:
        if self.contragent_form.changed():
            mess = "Відкинути не збережені зміни?"
            if not ok_cansel_dlg(mess):
                return
        return super().accept()    


class ProjectViewer(QSplitter):
    projectSelected = pyqtSignal(dict)
    projectDoubleClicked = pyqtSignal(dict)
    doAction = pyqtSignal(str, dict)
    def __init__(self):
        super().__init__(Qt.Orientation.Horizontal)
        self.project = Item('project')
        self.filter = ProjectFilter()
        self.addWidget(self.filter)
        self.filter.valuesChanged.connect(self.filter_changed)
        self.filter.valuesReloaded.connect(self.filter_reloaded)
        
        self.fields = ['id', 'contragent', 'project_type', 'name', 'created_at']
        self.by_client_fields = ['id', 'project_status', 'project_type', 'name', 'created_at']
        buttons=['create', 'copy', 'edit']
        self.table = TableWControls(self.project.model_w, self.fields, buttons=buttons)
        self.addWidget(self.table)
        folder_btn = QPushButton('До теки')
        self.table.hbox.insertWidget(0, folder_btn)
        folder_btn.clicked.connect(self.to_folder)
        client_btn = QPushButton('До клієнта')
        self.table.hbox.insertWidget(1, client_btn)
        client_btn.clicked.connect(self.to_client)
        move_btn = QPushButton('Перемістити')
        self.table.hbox.insertWidget(1, move_btn)
        move_btn.clicked.connect(self.move_dir)
        
        self.table.table.valueSelected.connect(self.project_selected)
        self.table.table.valueDoubleCklicked.connect(self.project_opened)
        
        self.table.actionInvoked.connect(self.action)
        self.filter.contragent_filter.contactsRequired.connect(self.show_contacts)
        self.filter.contragent_filter.actionInvoked.connect(self.act_contragent)
        self.setSizes([100, 400])

    def filter_changed(self, values):
        if self.filter.currentIndex() == 0:
            self.table.table.set_fields(self.fields)
        else:
            self.table.table.set_fields(self.by_client_fields)
            
        self.table.table.reload(values)
        index = self.table.table._model.createIndex(0, 0)
        self.table.table.setCurrentIndex(index)
        self.table.table.setFocus()

    def filter_reloaded(self, values):
        self.table.table.reload(values)
        
    def project_selected(self, value):
        self.projectSelected.emit(value)

    def project_opened(self, value):
        self.projectDoubleClicked.emit(value)

    def action(self, action, value):
        self.doAction.emit(action, value)
    
    def reload(self):
        self.filter.reload()

    def move_dir(self):
        project_value = self.table.table.get_selected_value()
        if not project_value:
            error('Оберіть проект')
            return
        
        if project_value['project_group_id'] == 4 or project_value['project_group_id'] == 5:
            res = ok_cansel_dlg("Перенести теку на новий сервер?")
            if not res:
                return
            project_value['project_group_id'] = 3
        else:
            res = ok_cansel_dlg("Перенести теку на старий сервер?")
            if not res:
                return
            project_value['project_group_id'] = 5
        
        project = Item('project')
        project.value = project_value
        err = project.copy_dirs()
        if err:
            error(err)
            return
        err = project.save()
        if err:
            error(err)
            return
        
        messbox(f'{project.value["name"]} скопійовано')

    def to_folder(self):
        project = self.table.table.get_selected_value()
        if not project:
            error('Оберіть проект')
            return
        app = App()
        path = self.make_project_path(project)
        subprocess.run([app.config['program'], path])
    
    def make_project_path(self, project):    
        app = App()
        base_dir = app.config['bs_makets_path']
        if project['project_group_id'] == 4 or project['project_group_id'] == 5:
            base_dir = app.config['makets_path']
        contragent = Item('contragent')
        err = contragent.get(project['contragent_id'])
        if err:
            error(err)
            return
        path = os.path.join(
                base_dir,
                contragent.value['dir_name'],
                project['type_dir'],
                project['number_dir'],
            )
        return path

    def to_client(self, id=0):
        if not id:
            project = self.table.table.get_selected_value()
            if not project:
                error('Оберіть проект')
                return
            id = project['contragent_id']
        contragent = Item('contragent')
        err = contragent.get_w(id)
        if err:
            error(err)
            return
        self.filter.setCurrentIndex(1)
        self.filter.contragent_filter.contragent_selector.reload([contragent.value])
        self.filter.contragent_filter.contragent_changed(contragent.value)
        
    def show_contacts(self, contragent_value):
        dlg = ContactSelectDialog(contragent_id=contragent_value['id'])
        dlg.exec()

    def act_contragent(self, action, value):
        contragent = Item('contragent')
        if action == 'edit':
            dlg = FormDialog('Редагувати', contragent.model, value)
        else:
            return
        res = dlg.exec()
        if not res:
            return
    
        contragent.value = dlg.value
        if not contragent.value:
            return
        err = contragent.save()
        if err:
            error(err)
            