import os
import subprocess
from PyQt6.QtWidgets import (
    QWidget,
    QHBoxLayout,
    QPushButton,
    )

from data.app import App
from data.model import Item
from widgets.Dialogs import error
from widgets.Form import ItemTableWithDetails, MainItemTable, DetailsItemTable


class InvoiceTab(ItemTableWithDetails):
    def __init__(self):
        super().__init__(MainItemTable('invoice'), DetailsItemTable('item_to_invoice'))
        app = App()
        self.base_dir = app.config['new_makets_path']
        self.file_mngr = app.config['program'] 
        to_folder_btn = QPushButton("До теки")
        open_doc_btn = QPushButton("Відкрити рахунок")
        
        self.table.info.grid.addWidget(to_folder_btn)
        self.table.info.grid.addWidget(open_doc_btn)
        
        to_folder_btn.clicked.connect(self.to_folder)
        open_doc_btn.clicked.connect(self.open_doc)

    def make_path(self):
        if not self.table.current_value:
            error('Оберіть рахунок')
            return
        invoice = self.table.current_value
        contragent_id = invoice['contragent_id']
        contragent = Item('contragent')
        err = contragent.get(contragent_id)
        if err:
            error(err)
            return
        path = os.path.join(
                self.base_dir,
                contragent.value['dir_name'],
                str(invoice['ordering_id']),
                'documents',
            )
        return path
    
    def to_folder(self):    
        doc_path = self.make_path()
        print(doc_path)
        subprocess.run([self.file_mngr, doc_path])

    def open_doc(self):
        invoice = self.table.current_value
        doc_path = self.make_path() 
        doc_path = os.path.join(doc_path, f"Рахунок до зам. {invoice['ordering_id']}.xlsx")
        print(doc_path)
        subprocess.run(['c:\\Program Files\\LibreOffice 5\\program\\scalc.exe', doc_path])