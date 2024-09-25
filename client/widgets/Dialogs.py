from PyQt6.QtCore import Qt
from PyQt6.QtWidgets import (
    QVBoxLayout,
    QWidget,
    QLabel,
    QDialogButtonBox as DBB,
    QMessageBox as MB,
    QDialog,
    QLineEdit,
    )

def error(text):
    MB(MB.Icon.Critical, "Помилка", text).exec()

def messbox(text, title="Повідомлення"):
    MB(MB.Icon.Information, title, text).exec()

def askdlg(question: str):
    dlg = AskDialog(question)
    res = dlg.exec()
    if not res:
        return ''
    return dlg.entry.text()

def ok_cansel_dlg(question: str, title: str='Запитання'):
    dlg = CustomDialog(QLabel(question), title)
    return dlg.exec()


# widget - widget to show on dialog
class CustomDialog(QDialog):
    def __init__(self, widget: QWidget, title: str, width=0, height=0):
        super().__init__()

        self.setWindowTitle(title)
        self.setWindowFlags(Qt.WindowType.WindowMinimizeButtonHint | Qt.WindowType.WindowMaximizeButtonHint | Qt.WindowType.WindowCloseButtonHint)
        self.widget = widget
        widget.setMinimumWidth(width)
        widget.setMinimumHeight(height)

        self.layout = QVBoxLayout()
        self.setLayout(self.layout)

        QBtn = DBB.StandardButton.Ok | DBB.StandardButton.Cancel
        self.buttonBox = DBB(QBtn)
        self.buttonBox.accepted.connect(self.accept)
        self.buttonBox.rejected.connect(self.reject)

        self.layout.addWidget(self.widget)
        self.layout.addWidget(self.buttonBox)


# Dialog with QLineEdit widget
class AskDialog(CustomDialog):
    def __init__(self, question: str):
        question_label = QLabel(question)
        self.entry = QLineEdit()
        w = QWidget()
        wl = QVBoxLayout()
        wl.addWidget(question_label)
        wl.addWidget(self.entry)
        w.setLayout(wl)
        super().__init__(w, 'Запитання')


class DeleteDialog(CustomDialog):
    def __init__(self, value, with_prompt=True):
        name = value["name"] if 'name' in value else value["id"]
        label_text = f'Ви дійсно бажаєте видалити {name}?'
        w = QWidget()
        wl = QVBoxLayout()
        wl.addWidget(QLabel(label_text))
        if with_prompt:
            wl.addWidget(QLabel('Причина:'))
            self.entry = QLineEdit()
            wl.addWidget(self.entry)
        w.setLayout(wl)
        super().__init__(w, "Видалити")

    