import webbrowser

from PyQt6 import QtGui
from PyQt6.QtCore import pyqtSignal
from PyQt6.QtWidgets import (
    QVBoxLayout,
    QLabel,
    QPushButton,
    QWidget,
    QHBoxLayout,
    QSplitter,
    QListWidget,
    QListWidgetItem,
    )

from widgets.Table import Table
from data.model import Item
from widgets.Dialogs import error
from widgets.Form import PeriodWidget, FormDialog, ItemTable
from data_widgets.ProjectFilter import ByContragentFilter
from common.params import FULL_VALUE_ROLE


class ContragentFilter(ByContragentFilter):
    contactChanged = pyqtSignal(dict)
    legalChanged = pyqtSignal(dict)
    def __init__(self, title: str = ''):
        self.contragent_group = Item('contragent_group')
        self.contragent_group.get_all_w()
        self.contact = Item('contact')
        self.legal = Item('legal')
        super().__init__(title, self.contragent_group.values)
        self.contragent_selector.valueDoubleClicked.connect(self.edit_contragent)

    def edit_contragent(self, value):
        dlg = FormDialog('Редагувати контрагента', self.contragent.model_w, self.contragent.columns, value)
        res = dlg.exec()
        if not res:
            return
        self.contragent.value = dlg.value
        err = self.contragent.save()
        if err:
            error(err)
    
    def append_widget(self):
        self.contacts_list = QListWidget()
        self.contacts_list.currentItemChanged.connect(self.contact_selection_changed)
        self.box.addWidget(self.contacts_list, 2)
        self.legals_list = QListWidget()
        self.legals_list.currentItemChanged.connect(self.legal_selection_changed)
        self.box.addWidget(self.legals_list, 1)

    def reload(self, values: list, list_widget: QListWidget):    
        list_widget.clear()
        for v in values:
            fi = QListWidgetItem(v['name'], list_widget)
            fi.setData(FULL_VALUE_ROLE, v)

    def contact_selection_changed(self, current, previous):
        if current is None:
            return
        value = current.data(FULL_VALUE_ROLE)
        self.contactChanged.emit(value)

    def legal_selection_changed(self, current, previous):
        if current is None:
            return
        value = current.data(FULL_VALUE_ROLE)
        self.legalChanged.emit(value)

    def contragent_changed(self, value):
        self.contact.get_filter_w('contragent_id', value['id'])
        self.reload(self.contact.values, self.contacts_list)
        self.legal.get_filter_w('contragent_id', value['id'])
        self.reload(self.legal.values, self.legals_list)
        return super().contragent_changed(value)


class ContragentsTab(QWidget):
    def __init__(self):
        super().__init__()
        self.current_contragent = {}
        self.current_contact = {}
        self.current_legal = {}
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
        viber_btn = QPushButton('Viber')
        self.hbox.addWidget(viber_btn)
        viber_btn.clicked.connect(self.to_viber)
        tg_btn = QPushButton('Telegram')
        self.hbox.addWidget(tg_btn)
        tg_btn.clicked.connect(self.to_telegram)
        self.period = PeriodWidget()
        self.hbox.addWidget(self.period)
        self.period.periodSelected.connect(self.period_selected)

        reload = QPushButton()
        reload.setIcon(QtGui.QIcon(f'images/icons/reload.png'))
        reload.setToolTip('Оновити')
        self.hbox.addWidget(reload)
        reload.clicked.connect(self.period_selected)

        reload_all = QPushButton('Оновити для всіх')
        self.hbox.addWidget(reload_all)
        reload_all.clicked.connect(lambda: self.period_selected(all_contragents=True))
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
        self.contragent = Item('contragent')
               
        self.contragent_filter = ContragentFilter()
        self.splitter.addWidget(self.contragent_filter)
        self.contragent_filter.contragentChanged.connect(self.current_contragent_changed)
        self.contragent_filter.contactChanged.connect(self.current_contact_changed)
        self.contragent_filter.legalChanged.connect(self.current_legal_changed)

        self.ordering = Item('ordering')
        self.orderings = ItemTable('ordering', buttons=[], show_period=False, is_info_bottom=True)
        self.splitter.addWidget(self.orderings)
        self.orderings.remove_dblclick_cb()
        
        model = {
            "id": {"def": 0, "hum": "Номер"},
            "name": {"def": "Замовлення", "hum": "Назва"},
            "created_at": {"def": "date", "hum": "Створений"},
            "cost": {"def": 0.0, "hum": "Вартість"},
            "type_hum": {"def": "", "hum": "Тип"},
        }
        self.docs_in_table = Table(model)
        self.splitter.addWidget(self.docs_in_table)
        self.whs_in = Item('whs_in')
        self.cash_in = Item('cash_in')
        
        self.cash_out = Item('cash_out')
        self.whs_out = Item('whs_out')
        self.invoice = Item('invoice')
        self.doc_out_table = Table(model)
        self.splitter.addWidget(self.doc_out_table)
            
    def current_contact_changed(self, contact_value):
        self.current_contact = contact_value
        self.period_selected()

    def current_legal_changed(self, legal_value):
        self.current_legal = legal_value
        self.period_selected()
    
    def current_contragent_changed(self, contragent_value):
        self.current_contragent = contragent_value
        self.current_contact = {}
        self.period_selected()

    def period_selected(self, date_from='', date_to = '', all_contragents=False):
        if date_from and date_to:
            self.date_from = date_from
            self.date_to = date_to
        elif not (self.date_from and self.date_to):
            self.date_from, self.date_to = self.period.get_period()

        doc_in_values = []
        self.cash_in.get_between_w('created_at', self.date_from, self.date_to)
        cash_in_values = self.cash_in.values
        if not all_contragents:
            if self.current_contragent:
                cash_in_values = [v for v in cash_in_values if v['contragent_id'] == self.current_contragent['id']]
            if self.current_contact:
                cash_in_values = [v for v in cash_in_values if v['contact_id'] == self.current_contact['id']]
            if self.current_legal:
                cash_in_values = [v for v in cash_in_values if v['legal_id'] == self.current_legal['id']]
        self.whs_in.get_between_w('created_at', self.date_from, self.date_to)
        whs_in_values = self.whs_in.values
        if not all_contragents:
            if self.current_contragent:
                whs_in_values = [
                    v for v in whs_in_values 
                    if v['contragent_id'] == self.current_contragent['id']
                    ]
            if self.current_contact:
                whs_in_values = [
                    v for v in whs_in_values 
                    if v['contact_id'] == self.current_contact['id']
                    ]
            if self.current_legal:
                whs_in_values = [
                    v for v in whs_in_values 
                    if v['legal_id'] == self.current_legal['id']
                    ]

        for v in cash_in_values:
            v['cost'] = v['cash_sum']
            v['type_hum'] = 'ПКО'
            v['type'] = 'cash_in'
            doc_in_values.append(v)

        for v in whs_in_values:
            v['cost'] = v['whs_sum']
            v['type_hum'] = 'ПН'
            v['type'] = 'whs_in'
            doc_in_values.append(v)
        
        sum_in = 0
        for v in doc_in_values:
            sum_in += v['cost'] if v['is_realized'] else 0
        self.in_sum.setText(str(round(sum_in, 2)))
        self.docs_in_table.reload(doc_in_values)
        
        doc_out_values = []
        self.cash_out.get_between_w('created_at', self.date_from, self.date_to)
        cash_out_values = self.cash_out.values
        if not all_contragents:
            if self.current_contragent:
                cash_out_values = [v for v in cash_out_values if v['contragent_id'] == self.current_contragent['id']]
            if self.current_contact:
                cash_out_values = [v for v in cash_out_values if v['contact_id'] == self.current_contact['id']]
            if self.current_legal:
                cash_out_values = [v for v in cash_out_values if v['legal_id'] == self.current_legal['id']]
        
        self.invoice.get_between_w('created_at', self.date_from, self.date_to)
        invoice_values = self.invoice.values
        if not all_contragents:
            if self.current_contragent:
                invoice_values = [v for v in invoice_values if v['contragent_id'] == self.current_contragent['id']]
            if self.current_contact:
                invoice_values = [v for v in invoice_values if v['contact_id'] == self.current_contact['id']]
            if self.current_legal:
                invoice_values = [v for v in invoice_values if v['legal_id'] == self.current_legal['id']]

        self.whs_out.get_between_w('created_at', self.date_from,self. date_to)
        whs_out_values = self.whs_out.values
        if not all_contragents:
            if self.current_contragent:
                whs_out_values = [v for v in whs_out_values if v['contragent_id'] == self.current_contragent['id']]
            if self.current_contact:
                whs_out_values = [v for v in whs_out_values if v['contact_id'] == self.current_contact['id']]
            if self.current_legal:
                whs_out_values = [v for v in whs_out_values if v['legal_id'] == self.current_legal['id']]

        for v in cash_out_values:
            v['cost'] = v['cash_sum']
            v['type_hum'] = 'ВКО'
            v['type'] = 'cash_out'
            doc_out_values.append(v)

        for v in invoice_values:
            v['cost'] = v['cash_sum']
            v['type_hum'] = 'Рахунок'
            v['type'] = 'invoice'
            doc_out_values.append(v)

        for v in whs_out_values:
            v['cost'] = v['whs_sum']
            v['type_hum'] = 'ВН'
            v['type'] = 'whs_out'
            doc_out_values.append(v)

        sum_out = 0
        for v in doc_out_values:
            sum_out += v['cost'] if v['is_realized'] else 0
        self.out_sum.setText(str(round(sum_out, 2)))
        self.doc_out_table.reload(doc_out_values)
        total = sum_in - sum_out
        self.total_sum.setText(str(round(total, 2)))
        if not all_contragents:
            if self.current_contact:
                sum_in_before, sum_out_before = self.calc_sums("contact_id", self.current_contact['id'], self.date_from)
            elif self.current_legal:
                sum_in_before, sum_out_before = self.calc_sums("legal_id", self.current_legal['id'], self.date_from)
            elif self.current_contragent:
                sum_in_before, sum_out_before = self.calc_sums("contragent_id", self.current_contragent['id'], self.date_from)
            else:
                return    
            total_before = sum_in_before - sum_out_before
            self.last_sum.setText(str(round(total_before, 2)))
            self.total_sum.setText(str(round(total + total_before, 2)))

        self.ordering.get_between_w('created_at', self.date_from, self. date_to)
        ordering_values = self.ordering.values
        if not all_contragents:
            if self.current_contragent:
                ordering_values = [v for v in ordering_values if v['contragent_id'] == self.current_contragent['id']]
            if self.current_contact:
                ordering_values = [v for v in ordering_values if v['contact_id'] == self.current_contact['id']]
            if self.current_legal:
                ordering_values = [v for v in ordering_values if v['legal_id'] == self.current_legal['id']]
        self.orderings.reload(ordering_values)
    
    def calc_sums(self, field, id, date_from):
        sum_in = 0
        sum_out = 0
        sum_in += self.calc_sum(self.cash_in, "cash_sum", field, id, date_from)
        sum_in += self.calc_sum(self.whs_in, "whs_sum", field, id, date_from)
        sum_out += self.calc_sum(self.cash_out, "cash_sum", field, id, date_from)
        sum_out += self.calc_sum(self.whs_out, "whs_sum", field, id, date_from)
        sum_out += self.calc_sum(self.invoice, "cash_sum", field, id, date_from)
        return sum_in, sum_out

    def calc_sum(self, table, sum_field, by_field, id, date_from):
        res = table.get_sum_before(sum_field, by_field, id, date_from)
        if res['error']:
            return 0
        else:
            return res['value']['sum']
        

    def reload(self, values=None):
        self.period_selected(all_contragents=True)

    def to_viber(self):
        if not self.current_contact:
            error("Оберіть контакт")
            return
        url = ''
        if self.current_contact['viber']:
            url = f'viber://chat?number={self.current_contact["viber"]}'
        elif len(self.current_contact['phone']) == 10:
            phone = self.current_contact['phone'][1:]
            url = f'viber://chat?number={phone}'
        if url:
            webbrowser.open(url)

    def to_telegram(self):
        if not self.current_contact:
            error("Оберіть контакт")
            return
        url = ''
        if self.current_contact['telegram']:
            url = f'tg://resolve?domain={self.current_contact["telegram"]}'
        elif len(self.current_contact['phone']) == 10:
            phone = '38' + self.current_contact['phone']
            url = f'tg://resolve?phone={phone}'
        if url:
            webbrowser.open(url)

    