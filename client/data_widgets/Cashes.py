from PyQt6.QtGui import QIcon
from PyQt6.QtWidgets import (
    QVBoxLayout,
    QLabel,
    QPushButton,
    QWidget,
    QLineEdit,
    QHBoxLayout,
    QSplitter,
    )

from data.model import Item
from data.app import App
from widgets.Dialogs import error
from widgets.Form import PeriodWidget, Selector, CustomDialog, ItemTable

class CashesTab(QWidget):
    def __init__(self):
        super().__init__()
        self.current_cash = {}
        self.date_from = ''
        self.date_to = ''
        
        self.box = QVBoxLayout()
        self.box.setContentsMargins(0,0,0,0)
        self.setLayout(self.box)
        
        controls = QWidget()
        self.box.addWidget(controls, 0)
        self.hbox = QHBoxLayout()
        self.hbox.setContentsMargins(0,0,0,0)
        controls.setLayout(self.hbox)
        self.period = PeriodWidget()
        self.hbox.addWidget(self.period)
        self.period.periodSelected.connect(self.period_selected)
        reload = QPushButton()
        reload.setIcon(QIcon(f'images/icons/reload.png'))
        reload.setToolTip('Оновити')
        self.hbox.addWidget(reload)
        reload.clicked.connect(self.period_selected)

        reload_all = QPushButton('Оновити для всіх')
        self.hbox.addWidget(reload_all)
        reload_all.clicked.connect(lambda: self.period_selected(all_cashes=True))
        self.hbox.addStretch()

        cashes_move = QPushButton('Перевести до каси')
        self.hbox.addWidget(cashes_move)
        cashes_move.clicked.connect(self.cash_move)
        self.hbox.addStretch()
        
        self.hbox.addWidget(QLabel('Залишок'))
        self.last_sum = QLabel()
        self.hbox.addWidget(self.last_sum)
        self.hbox.addWidget(QLabel('Прибуток'))
        self.in_sum = QLabel()
        self.hbox.addWidget(self.in_sum)
        self.hbox.addWidget(QLabel('Видаток'))
        self.out_sum = QLabel()
        self.hbox.addWidget(self.out_sum)
        self.hbox.addWidget(QLabel('Загалом'))
        self.total_sum = QLabel()
        self.hbox.addWidget(self.total_sum)
        
        self.splitter = QSplitter()
        self.box.addWidget(self.splitter, 10)
        self.cash = Item('cash')
        err = self.cash.get_all_w()
        if err:
            error(err)
            self.cash.values = []
        self.cash_table = ItemTable(
            'cash',
            values=self.cash.values,
            fields=["name", "total", "comm"],
            is_info_bottom=True,
        )
        self.splitter.addWidget(self.cash_table)
        self.cash_table.table.table.valueSelected.connect(self.current_cash_changed)

        self.cash_in = Item('cash_in')
        self.cash_in_table = ItemTable(
            'cash_in',
            fields=["id", "name", "cash_sum", "comm"],
            buttons=[],
            is_info_bottom=True,
            show_period=False,
            )
        self.splitter.addWidget(self.cash_in_table)
        self.cash_in_table.remove_dblclick_cb()
        
        self.cash_out = Item('cash_out')
        self.cash_out_table = ItemTable(
            'cash_in',
            fields=["id", "name", "cash_sum", "comm"],
            buttons=[],
            is_info_bottom=True,
            show_period=False,
            )
        self.splitter.addWidget(self.cash_out_table)
        self.cash_out_table.remove_dblclick_cb()

    def reload(self):
        self.cash_table.reload()
        self.period_selected(all_cashes=True)
    
    def cash_move(self):
        w = QWidget()
        box = QVBoxLayout()
        w.setLayout(box)
        box.addWidget(QLabel(f"З каси: {self.current_cash['name']}"))
        box.addWidget(QLabel(f"В наявності: {self.current_cash['total']}"))
        box.addWidget(QLabel("Вкажіть суму переказу:"))
        trans_sum = QLineEdit()
        box.addWidget(trans_sum)
        box.addWidget(QLabel("Оберіть касу:"))
        cash_sel = Selector('cash')
        box.addWidget(cash_sel)
        dlg = CustomDialog(w, 'Переказ між касами')
        res = dlg.exec()
        if not res:
            return
        app = App()
        try:
            t_sum = int(trans_sum.text())
        except:
            error("Неправильний формат суми")
            return
        cash_to_value = cash_sel.full_value()

        cash_out = Item('cash_out')
        cash_out.create_default()
        cash_out.value['cash_id'] = self.current_cash['id']
        cash_out.value['cash_sum'] = t_sum
        cash_out.value['user_id'] = app.user['id']
        cash_out.value['comm'] = 'Переказ з цієї каси'
        cash_out.value['name'] += f" до каси {cash_to_value['name']}"
        err = cash_out.save()
        if err:
            error(err)
            return
        cash_in = Item('cash_in')
        cash_in.create_default()
        cash_in.value['cash_id'] = cash_sel.value()
        cash_in.value['cash_sum'] = t_sum
        cash_in.value['user_id'] = app.user['id']
        cash_in.value['comm'] = 'Переказ до цієї каси'
        cash_in.value['name'] += f" з каси {self.current_cash['name']}"
        err = cash_in.save()
        if err:
            error(err)
            return
        self.cash_table.reload()
        
    
    def current_cash_changed(self, cash_value):
        self.current_cash = cash_value
        self.period_selected()

    def period_selected(self, date_from='', date_to = '', all_cashes=False):
        if date_from and date_to:
            self.date_from = date_from
            self.date_to = date_to
        elif not (self.date_from and self.date_to):
            self.date_from, self.date_to = self.period.get_period()

        self.cash_in.get_between_w('created_at', self.date_from, self.date_to)
        cash_in_values = self.cash_in.values
        if not all_cashes:
            if self.current_cash:
                cash_in_values = [v for v in cash_in_values if v['cash_id'] == self.current_cash['id']]
        sum_in = 0
        for v in cash_in_values:
            sum_in += v['cash_sum']
        self.in_sum.setText(str(sum_in))
        self.cash_in_table.reload(cash_in_values)
        
        self.cash_out.get_between_w('created_at', self.date_from, self.date_to)
        cash_out_values = self.cash_out.values
        if not all_cashes:
            if self.current_cash:
                cash_out_values = [v for v in cash_out_values if v['cash_id'] == self.current_cash['id']]
        sum_out = 0
        for v in cash_out_values:
            sum_out += v['cash_sum']
        self.out_sum.setText(str(sum_out))
        self.cash_out_table.reload(cash_out_values)
        self.total_sum.setText(str(sum_in - sum_out))
        if not all_cashes:
            if self.current_cash:
                sum_in_res = self.cash_in.get_sum_before("cash_sum", "cash_id", self.current_cash['id'], self.date_from)
                sum_out_res = self.cash_out.get_sum_before("cash_sum", "cash_id", self.current_cash['id'], self.date_from)
                self.last_sum.setText(str(sum_in_res['value']['sum'] - sum_out_res['value']['sum']))
    
    
    