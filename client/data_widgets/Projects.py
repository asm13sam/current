import os
from PyQt6.QtGui import QPixmap
from PyQt6.QtCore import Qt, pyqtSignal
from PyQt6.QtWidgets import (
    QSplitter,
    QTabWidget,
    QTextEdit,
    QScrollArea,
    QWidget,
    QVBoxLayout,
    QLabel
    )

from data.model import Item
from data.app import App
from data_widgets.ProjectViewer import ProjectViewer
from widgets.Form import InfoBlock
from widgets.Dialogs import error, CustomDialog


class Projects(QSplitter):
    orderingRequired = pyqtSignal(dict)
    def __init__(self):
        super().__init__(Qt.Orientation.Horizontal)
        self.setStyleSheet('padding: 1px; margin: 0px;')
        self.project_viewer = ProjectViewer()
        self.tabs = ProjectTabs(self.project_viewer)
        self.addWidget(self.tabs)
        
        self.insertWidget(0, self.project_viewer)
        self.project_viewer.projectSelected.connect(self.tabs.current_reload)
        self.tabs.projectSaved.connect(self.project_viewer.reload)
        self.tabs.contragentRequired.connect(self.project_viewer.to_client)
        self.tabs.orderingRequired.connect(lambda project_value: self.orderingRequired.emit(project_value))
        self.setStretchFactor(0, 2)
        self.setStretchFactor(1, 5)

    def reload(self):
        pass


class ProjectTabs(QTabWidget):
    projectSaved = pyqtSignal()
    contragentRequired = pyqtSignal(int)
    orderingRequired = pyqtSignal(dict)
    def __init__(self, project_viewer):
        super().__init__()
        self.counter = 0
        self.project_current = ProjectCurrent(project_viewer)
        self.addTab(self.project_current, 'Поточний')
        
    def current_reload(self, project_value):
        self.project_current.reload(project_value)
        self.setCurrentWidget(self.project_current)


class ProjectCurrent(QSplitter):
    def __init__(self, project_viewer):
        super().__init__(Qt.Orientation.Horizontal)
        self.project = Item('project')
        self.project_viewer = project_viewer
        info_fields = [
            "id",
        #    "document_uid",
            "name",
            "project_group",
            "user",
            "contragent",
            "contact",
            "phone",
            "email",
            "cost",
            "cash_sum",
            "whs_sum",
            "project_type",
        #    "type_dir",
            "project_status",
        #    "number_dir",
        #    "info",
            "created_at",
            "is_in_work",
        #    "is_active",
        ]
        data_model = self.project.model_w.copy()
        contact = Item('contact')
        data_model['phone'] = contact.model['phone']
        data_model['email'] = contact.model['email']
        self.info_block = InfoBlock(data_model=data_model, field_names=info_fields)
        scroll = QScrollArea()
        scroll.setWidget(self.info_block)
        scroll.setWidgetResizable(True)
        self.addWidget(scroll)
        self.setSizes([100, 400])

        self.tabs = QSplitter(Qt.Orientation.Horizontal)
        self.addWidget(self.tabs)
        self.info = QTextEdit()
        self.tabs.addWidget(self.info)
        self.info.setReadOnly(True)
        self.order = Item('ordering')
        pic_scrol = QScrollArea()
        thbs = QWidget()
        self.thumbs = QVBoxLayout()
        thbs.setLayout(self.thumbs)
        pic_scrol.setWidget(thbs)
        pic_scrol.setWidgetResizable(True)
        self.tabs.addWidget(pic_scrol)
        
    def reload(self, value=None):
        if value is not None:
            self.project.value = value
        else:
            if self.project.value:
                id = self.project.value['id']
                err = self.project.get_w(id)
                if err:
                    error(err)
                    return
                value = self.project.value
            else:
                return
        contact = Item('contact')
        err = contact.get(value['contact_id'])
        if err:
            error(f'On get contact id: {err}')
            value['email'] = ''
            value['phone'] = ''
        else:
            value['email'] = contact.value['email']
            value['phone'] = contact.value['phone']
        self.info_block.reload(value)
        self.info.setText(value['info'])
        res = self.show_thumbs(self.project.value)
        if not res:
            self.create_thumbs(self.project.value)

    def create_thumbs(self, value):
        # self.clear_thumbs()
        path = self.project_viewer.make_project_path(value)
        th_files = self.get_thumbs_files(path)
        if not th_files:
            return
        counter = 0
        for tf in th_files:
            pixes = self.extract_previews(tf)
            for p in pixes:
                label = QLabel()
                pixmap = QPixmap(p)
                label.setPixmap(pixmap)
                self.thumbs.addWidget(label)
                counter += 1
                if counter >= 10:
                    return

    def get_thumbs_files(self, path):
        res = []
        for root, _, files in os.walk(path):
            for file_name in files:
                if file_name == 'Thumbs.db':
                    res.append(os.path.join(root, file_name))
        return res
    
    def extract_previews(self, th_file):
        res = []
        separator = b'\xff\xd8\xff\xe0'
        with open(th_file, 'rb') as f:
            content = f.read()
            thumbs = content.split(separator)
            for thumb in thumbs:
                pix = QPixmap()
                loaded = pix.loadFromData(separator + thumb, 'JPG')
                if loaded:
                    res.append(pix)
                    if len(res) >= 10:
                        return res
        return res
    
    def clear_thumbs(self):
        while self.thumbs.count():
            child = self.thumbs.takeAt(0)
            if child.widget():
                child.widget().deleteLater()

    def show_thumbs(self, value):
        count = 0
        self.clear_thumbs()
        path = self.project_viewer.make_project_path(value)
        prew = os.path.join(path, 'preview')
        if os.path.exists(prew):
            files = [f for f in os.listdir(prew) if os.path.isfile(os.path.join(prew, f))]
            for f in files:
                label = QLabel()
                pixmap = QPixmap()
                res = pixmap.load(os.path.join(prew, f))
                if res:
                    if pixmap.width() > 300:
                        pixmap = pixmap.scaledToWidth(256)
                    label.setPixmap(pixmap)
                    self.thumbs.addWidget(label)
                    count += 1
        return count        
            
    def doc_dblclicked(self, value):
        app = App()
        model = {k: v for k, v in app.model_w[value['type']].items() 
                    if not (k.endswith('_id') or k.endswith('_uid') or k == 'comm')}
        if value['type'].startswith('cash_'):
            value['cash_sum'] = value['cost']
        elif value['type'].startswith('whs_'):
            value['whs_sum'] = value['cost']
        info = InfoBlock(model, value=value)
        dlg = CustomDialog(widget=info, title='Document')
        dlg.exec()

