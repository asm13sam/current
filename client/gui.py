from PyQt6.QtWidgets import (
    QVBoxLayout,
    QWidget,
    QApplication,
    QTabWidget,
    QLabel,
    QHBoxLayout,
    )
from PyQt6.QtGui import QKeySequence, QShortcut, QFont

import sys
import qdarktheme

from data.app import App
from widgets.ComboBoxSelector import ComboBoxSelector
from data.model import Item
from data_widgets.Orders import ProductsTab
from data_widgets.CheckBoxUa import CheckBoxFS
from widgets.Dialogs import error
from widgets.Form import TabGroup, ItemTable, ItemTableWithDetails, MainItemTable, DetailsItemTable
from data_widgets.Whs import WhsInTab, WhsOutTab, WhsesTab
from data_widgets.Projects import Projects
from data_widgets.Cashes import CashesTab
from data_widgets.Contragents import ContragentsTab
from data_widgets.ItemsToOrdering import ItemsToOrdering, ToOrdering
from data_widgets.ItemsToProduct import ItemsToProduct
from data_widgets.Users import UserOrderingDetailsTable
from data_widgets.CboxChecks import CboxCheck
from data_widgets.Calculation import CalculatorTab
from data_widgets.Dashboard import Dashboard
from common.funcs import prepare_search_string
from data_widgets.Matherial import MatherialTab
from data_widgets.Operation import OperationTab
from data_widgets.Invoice import InvoiceTab
from data_widgets.Admin import AdminTab


class Window(QWidget):
    def __init__(self):
        super().__init__()
        self.setWindowTitle("Таргет")
        app = App()
        cfg = app.config
        self.check = CheckBoxFS(cfg['checkbox_url'], cfg['license_key'], cfg['cashier_pin'])
        self.main_layout = QVBoxLayout()
        self.setLayout(self.main_layout)
        self.main_layout.setContentsMargins(0, 0, 0, 0)
        self.main_tabs = QTabWidget()
        self.main_layout.addWidget(self.main_tabs, stretch=10)
        self.main_tabs.currentChanged.connect(self.reload_tab)
        self.make_gui()

    def make_gui(self):
        self.controls = QWidget()
        self.docs = {}
        self.controls_layout = QHBoxLayout()
        self.controls.setLayout(self.controls_layout)
        self.controls_layout.setContentsMargins(0, 5, 0, 0)
        self.main_layout.insertWidget(0, self.controls, stretch=0)
        # user = Item('user')
        # user.get_all()
        user_values = [
            {'id': 2, 'login': 'Serhii' , 'name': 'Сергій'},
            {'id': 3, 'login': 'Svitlana' , 'name': 'Світлана'},
            {'id': 5, 'login': 'Vadim' , 'name': 'Вадим'},
        ]
        self.users = ComboBoxSelector(values=user_values)
        self.users.selectionChanged.connect(self.change_user)
        self.controls_layout.addWidget(QLabel('Менеджер:'))
        self.controls_layout.addWidget(self.users)
        self.controls_layout.addStretch()
        self.main_tabs.setDisabled(True)

    def make_add_gui(self):
        sc = QShortcut(QKeySequence('Ctrl+f'), self)
        sc.activated.connect(self.find_from_clipboard)
        self.ordering = ItemsToOrdering(self.check)
        self.main_tabs.addTab(self.ordering, 'Замовити')
        self.ordering.doc_table.valueDoubleCklicked.connect(self.open_document)
        self.ordering.viewContragentReport.connect(self.open_contragent_report)
        self.docs['ordering'] = self.ordering

        calculator = CalculatorTab()
        self.main_tabs.addTab(calculator, 'Калькулятор')
        copycenter = ProductsTab(self.check)
        self.main_tabs.addTab(copycenter, 'Копіцентр')

        self.doc_tabs = TabGroup()
        self.main_tabs.addTab(self.doc_tabs, 'Документи')

        to_ord = ToOrdering()
        self.doc_tabs.addTab(to_ord, 'Замовлення')
        to_ord.doc_table.valueDoubleCklicked.connect(self.open_document)

        fields = [
            "id",
            "name",
            "cash",
            "cash_sum",
            "contragent",
            "contact",
            "created_at",
            "user",
            "comm",
        ]

        cash_in = ItemTable('cash_in', 'name', fields=fields)
        self.doc_tabs.addTab(cash_in, 'ПКО')
        self.docs['cash_in'] = cash_in
        cash_out = ItemTable('cash_out', 'name', fields=fields)
        self.doc_tabs.addTab(cash_out, 'ВКО')
        self.docs['cash_out'] = cash_out

        whs_in = WhsInTab()
        self.doc_tabs.addTab(whs_in, 'ПН')
        self.docs['whs_in'] = whs_in
        whs_in.doc_table.valueDoubleCklicked.connect(self.open_document)

        whs_out = WhsOutTab()
        self.doc_tabs.addTab(whs_out, 'ВН')
        self.docs['whs_out'] = whs_out

        invoice = InvoiceTab()
        self.doc_tabs.addTab(invoice, 'Рахунки')
        self.docs['invoice'] = invoice

        checks = CboxCheck(self.check)
        self.doc_tabs.addTab(checks, 'Чеки')
        self.docs['cbox_check'] = checks

        self.rep_tabs = TabGroup()
        self.main_tabs.addTab(self.rep_tabs, 'Звіти')

        cashes = CashesTab()
        self.rep_tabs.addTab(cashes, 'Каси')
        cashes.cash_in_table.table.table.valueDoubleCklicked.connect(
            lambda v: self.open_document(v, doc_type='cash_in')
            )
        cashes.cash_out_table.table.table.valueDoubleCklicked.connect(
            lambda v: self.open_document(v, doc_type='cash_out')
            )
        self.contragents = ContragentsTab()
        self.rep_tabs.addTab(self.contragents, 'Контрагенти')
        self.contragents.docs_in_table.valueDoubleCklicked.connect(self.open_document)
        self.contragents.doc_out_table.valueDoubleCklicked.connect(self.open_document)
        self.contragents.orderings.table.table.valueDoubleCklicked.connect(
            lambda v: self.open_document(v, doc_type='ordering')
            )

        whses = WhsesTab()
        self.rep_tabs.addTab(whses, 'Склади')
        whses.whs_table.doc_table.valueDoubleCklicked.connect(self.open_document)
        users = ItemTableWithDetails(MainItemTable('user'), UserOrderingDetailsTable())
        self.rep_tabs.addTab(users, 'Користувачі')
        counters = ItemTableWithDetails(MainItemTable('counter'), DetailsItemTable('record_to_counter'))
        self.rep_tabs.addTab(counters, 'Лічильники')

        self.cat_tabs = TabGroup()
        self.main_tabs.addTab(self.cat_tabs, 'Довідники')
        maths = MatherialTab()
        self.cat_tabs.addTab(maths, 'Матеріали')
        colors = ItemTable('color', 'name')
        self.cat_tabs.addTab(colors, 'Кольори')
        opers = OperationTab()
        self.cat_tabs.addTab(opers, 'Операції')
        equips = ItemTable('equipment', 'name')
        self.cat_tabs.addTab(equips, 'Обладнання')
        prods = ItemsToProduct()
        self.cat_tabs.addTab(prods, 'Вироби')
        owner = ItemTable('owner')
        self.cat_tabs.addTab(owner, 'Власні рахунки')

        self.projects_tab = Projects()
        self.main_tabs.addTab(self.projects_tab, 'Проекти')

        self.test_tabs = TabGroup()
        self.main_tabs.addTab(self.test_tabs, 'Test')
        # calc1 = CalculatorTab()
        # self.test_tabs.addTab(calc1, 'Калькуляtor')
        # orders = ProductsTab(self.check)
        # self.test_tabs.addTab(orders, 'Замовлення')
        # ord2_cr = OrderingSimpleCreator()
        # self.test_tabs.addTab(ord2_cr, 'Замовлення2')
        p2os_work = Dashboard()
        self.test_tabs.addTab(p2os_work, 'В роботі')
        admin = AdminTab()
        self.test_tabs.addTab(admin, "Адміністрування")


    def reload_tab(self, index):
        w = self.main_tabs.widget(index)
        w.reload()

    def open_document(self, value, doc_type=''):
        if doc_type:
            dtype = doc_type
        elif 'type' in value:
            dtype = value['type']
        else:
            return
        if dtype != 'ordering':
            self.main_tabs.setCurrentWidget(self.doc_tabs)
            self.doc_tabs.setCurrentWidget(self.docs[dtype])
        else:
            self.main_tabs.setCurrentWidget(self.docs['ordering'])
        item = Item(dtype)
        err = item.get_w(value['id'])
        if err:
            error(err)
            return
        self.docs[dtype].reload([item.value])

    def open_contragent_report(self, contragent_id):
        self.main_tabs.setCurrentWidget(self.rep_tabs)
        self.rep_tabs.setCurrentWidget(self.contragents)
        self.contragents.contragent_filter.set_contragent_by_id(contragent_id)

    def find_from_clipboard(self):
        clipboard = QApplication.clipboard()
        s = prepare_search_string(clipboard.text())
        self.ordering.side_controls.setCurrentIndex(1)
        self.ordering.contragent_filter.contragent_selector.search_entry.setText(s)
        self.projects_tab.project_viewer.filter.contragent_filter.contragent_selector.search_entry.setText(s)
        self.projects_tab.project_viewer.filter.setCurrentIndex(1)

    def change_user(self, value):
        app = App()
        print(value['login'])
        res = app.repository.signing(value['login'], '123')
        print('login', res)
        if res['error']:
            return
        user = Item('user')
        err = user.get(value['id'])
        if err:
            error(err)
            return
        app.user = user.value
        self.main_tabs.setDisabled(False)
        self.users.setDisabled(True)
        self.controls.setVisible(False)
        self.setWindowTitle(f"Таргет - {value['name']}")
        self.make_add_gui()


class MainWindow():
    def __init__(self):
        self.qt_app = QApplication(sys.argv)
        # self.app = App()
        # res = self.app.repository.signing('Vadim', '123')
        # print('login', res)
        screen = self.qt_app.primaryScreen()
        print('Screen: %s' % screen.name())
        size = screen.size()
        print('Size: %d x %d' % (size.width(), size.height()))
        rect = screen.availableGeometry()
        print('Available: %d x %d' % (rect.width(), rect.height()))
        self.window = Window()
        # self.window.showFullScreen()
        # self.window.showMaximized()
        color = "#99BCBC"
        qss = """
        QToolTip {
            background-color: black;
            color: white;
            border: black solid 1px
                }
        """
        qdarktheme.setup_theme(custom_colors={'primary': color}, additional_qss=qss)
        font = QFont()

        QApplication.instance().setFont(font)
        font.setPointSize(10)


    def run(self):
        self.window.show()
        sys.exit(self.qt_app.exec())


if __name__ == '__main__':
    w = MainWindow()
    w.run()
