import json

from PyQt6 import QtWebSockets, QtNetwork
from PyQt6.QtCore import (
    Qt,
    pyqtSignal,
    QUrl,
    QByteArray,
    )
from PyQt6.QtWidgets import (
    QLabel,
    QPushButton,
    QWidget,
    QVBoxLayout,
    QHBoxLayout,
    QSplitter,
    QTabWidget,
    QListWidget,
    QListWidgetItem,
    QTextEdit,
    )

from data.model import Item
from data.app import App
from widgets.Form import ItemTree
from widgets.Dialogs import error


class ProductsInWork(ItemTree):
    def __init__(self, 
                 fields: list = [], 
                 values: list = None, 
                 ):
        self.ordering = Item('ordering')
        
        super().__init__('product_to_ordering', 
                         'product_to_ordering', 
                         fields=fields, 
                         values=values, 
                         buttons={'reload':'Оновити', 'edit':'Редагувати'},
                         show_info=True, 
                         )
        
    def reload(self, values=None):
        app = App()
        if values is None:
            err = self.ordering.get_filter('ordering_status_id',  app.config["ordering state in work"])
            if err:
                error(err)
                self.ordering.values = []
            values = []
            p2o = Item('product_to_ordering')
            for v in self.ordering.values:
                err = p2o.get_filter_w('ordering_id', v['id'])
                if err:
                    error(err)
                    continue
                values += p2o.values
        return super().reload(values)
    

class ChatMessagesList(QListWidget):
    def __init__(self):
        super().__init__()
        # self.setSizePolicy(QSizePolicy(QSizePolicy.Policy.Minimum, QSizePolicy.Policy.Minimum))
        # self.setMaximumHeight(30)

    # def setDataset(self, dataset):
    #     self.clear()
    #     self.setMaximumHeight(25*len(dataset))
    #     for v in dataset:
    #         listItem = QListWidgetItem(v[1], self)
    #         listItem.setData(FULL_VALUE_ROLE, v[0])
    #         listItem.setCheckState(Qt.CheckState.Checked if v[2] else Qt.CheckState.Unchecked)

    def add_message(self, message):
        listItem = QListWidgetItem(message, self)


    # def get_checked(self):
    #     res = []
    #     for i in range(self.count()):
    #         li = self.item(i)
    #         if li.checkState() == Qt.CheckState.Checked:
    #             res.append(li.data(FULL_VALUE_ROLE))
    #     return res


class Chat(QWidget):
    # TextRecieved = pyqtSignal(str)
    def __init__(self):
        super().__init__()
        app = App()
        self.user = app.user
        url = f'ws://{app.config["host"]}:{app.config["port"]}/ws'
        self.box = QVBoxLayout()
        self.setLayout(self.box)
        self.mes_list = ChatMessagesList()
        self.box.addWidget(self.mes_list, 10)
        self.text_box = QTextEdit()
        self.box.addWidget(self.text_box)
        send_btn = QPushButton("Відправити")
        self.box.addWidget(send_btn)
        send_btn.clicked.connect(self.send_message)

        self.wsapp = QtWebSockets.QWebSocket("", QtWebSockets.QWebSocketProtocol.Version(13), None)
        self.wsapp.error.connect(self.error)
        req = QtNetwork.QNetworkRequest(QUrl(url))
        # req.setCookieHeader()
        req.setHeader(
            QtNetwork.QNetworkRequest.KnownHeaders.CookieHeader, 
            QtNetwork.QNetworkCookie(
                name=QByteArray('sess'.encode()), 
                value=QByteArray(app.repository.cookie_jar.get('sess').encode()),
                ),
            )
        self.wsapp.open(req)
        self.wsapp.pong.connect(self.onPong)
        self.wsapp.textMessageReceived.connect(self.on_message)
        app.repository.wsconn = self.wsapp

    def send(self, mess):
        self.wsapp.sendTextMessage(mess)
        print('send', mess)
        if mess == 'quit':
            self.wsapp.close()
    
    def do_ping(self):
        print("client: do_ping")
        self.wsapp.ping(b"foo")

    def onPong(self, elapsedTime, payload):
        print("onPong - time: {} ; payload: {}".format(elapsedTime, payload))

    def error(self, error_code):
        print("error code: {}".format(error_code))
        print(self.wsapp.errorString())

    def close(self):
        self.wsapp.close()

    def send_message(self):
        user = self.user['name']
        user_id = self.user['id']
        mess = self.text_box.toPlainText()
        self.mes_list.add_message(f'>> {mess}')
        self.text_box.clear()
        self.send(json.dumps({'user_id': user_id, 'username': user, 'message': mess}))

    def on_message(self, message):
        print('wrec', message)
        mess = json.loads(message)
        self.mes_list.add_message(f'{mess["username"]} >> {mess["message"]}')
        # self.TextRecieved.emit(message)


class Dashboard(QWidget):
    def __init__(self):
        super().__init__()
        self.box = QVBoxLayout()
        self.setLayout(self.box)
        self.main = QSplitter(Qt.Orientation.Horizontal)
        self.box.addWidget(self.main)
        self.piw = ProductsInWork()
        chat = Chat()
        self.main.addWidget(self.piw)
        self.main.addWidget(chat)

    def reload(self):
        self.piw.reload()
