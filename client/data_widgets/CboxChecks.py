from datetime import datetime

from PyQt6 import QtPrintSupport
from PyQt6.QtGui import QPainter
from PyQt6.QtWidgets import QLabel, QPushButton

from data.model import Item
from widgets.Form import CustomForm, DetailsItemTable, ItemTableWithDetails, MainItemTable, CustomFormDialog
from widgets.Dialogs import error, CustomDialog


class CboxCheck(ItemTableWithDetails):
    def __init__(self, checkbox):
        self.check = checkbox
        self.main_table = MainItemTable('cbox_check') 
        self.details_table = DetailsItemTable('item_to_cbox_check')
        super().__init__(self.main_table, self.details_table)
        self.main_table.remove_dblclick_cb()
        self.main_table.table.table.valueDoubleCklicked.connect(self.on_doubleclicked)
        post_pay_btn = QPushButton("Доплата")
        self.main_table.info.grid.addWidget(post_pay_btn)
        post_pay_btn.clicked.connect(self.create_post_pay)

    def create_post_pay(self):
        value = self.main_table.current_value
        if not value:
            return
        
        if not self.check.is_signed:
            res = self.check.sign_in()
            if not res:
                error("Не можу під'єднатися до сервера CheckBox!")
                return
            if 'error' in res:
                error(f"При підключенні до Checkbox: {res['error']}")
                return
        
        cash_state = self.check.get_cash_state()
        if not cash_state:
            error("Не можу отримати статус каси з сервера CheckBox!")
            return
        if 'error' in cash_state:
            error(f"Отримати статус каси з Checkbox: {cash_state['error']}")
            return    
        if not cash_state['has_shift']:
            shift_data = self.check.shifts()
            if not shift_data:
                error("Не можу відкрити зміну у CheckBox!")
                return
            if 'error' in shift_data:
                error(f"При відкритті зміни у Checkbox: {shift_data['error']}")
                return 
            
        c = self.check.get_check(value['checkbox_uid'])
        if 'error' in c:
            error(f"При отриманні чеку з Checkbox: {c['error']}")
            return
        if c['pre_payment_relation_id']:
            ordering = Item('ordering')
            err = ordering.get(value['ordering_id'])
            if err:
                error(err)
                return
            pay_sum = ordering.value['cost'] - c['total_sum']/100
            i = Item('cbox_check')
            i.value = value
            i.value["id"] = 0
            i.value["name"] = 'Доплата до ' + i.value["name"]
            i.value["cash_sum"] = pay_sum
            i.value['comm'] = 'Доплата'
            i.value["fs_uid"] = '0'
            i.value["checkbox_uid"] = '0'
            i.value["created_at"] = datetime.now().isoformat(timespec='seconds')
            
            form = CustomForm(i.model, value=i.value)
            dlg = CustomFormDialog('Чек доплати', form)
            res = dlg.exec()
            if res and dlg.value:
                i.value = dlg.value
                err = i.save()
                if err:
                    error(err)
                    return
                payments = []
                if i.value['is_cash']:
                    payments.append({
                        "type": 'CASH',
                        "value": int((i.value['cash_sum']) * 100),
                        "label": 'Готівка',
                    })
                else:
                    payments.append({
                        "type": 'CASHLESS',
                        "value": int(i.value['cash_sum'] * 100),
                        "label": 'Картка',
                    })
                
                receipt = {
                    # "cashier_name": "Олег",
                    "department": "Копіцентр",
                    "goods": [],
                    "payments": payments,
                }
                res = self.check.create_post_receipt(receipt, c['pre_payment_relation_id'])
                if not res:
                    error("Помилка при створенні чеку в Checkbox")
                    return
                if 'error' in res:
                    error(f"При створенні чеку в Checkbox: {res['error']}")
                    return
                i.value["checkbox_uid"] = res['id']
                err = i.save()
                if err:
                    error(err)
                    return
                self.main_table.reload()
            

    def on_doubleclicked(self, value):
        if value["checkbox_uid"] == '0':
            self.create_checkbox_check(value)
            return
        if value["fs_uid"] == '0':
            res = self.get_check_fs_code(value)
            if not res:
                return
        self.get_check_png(value['checkbox_uid'])

    def create_checkbox_check(self, value):
        if not self.check.is_signed:
            res = self.check.sign_in()
            if not res:
                error("Не можу під'єднатися до сервера CheckBox!")
                return
            if 'error' in res:
                error(f"При підключенні до Checkbox: {res['error']}")
                return
        
        cash_state = self.check.get_cash_state()
        if not cash_state:
            error("Не можу отримати статус каси з сервера CheckBox!")
            return
        if 'error' in cash_state:
            error(f"Отримати статус каси з Checkbox: {cash_state['error']}")
            return    
        if not cash_state['has_shift']:
            shift_data = self.check.shifts()
            if not shift_data:
                error("Не можу відкрити зміну у CheckBox!")
                return
            if 'error' in shift_data:
                error(f"При відкритті зміни у Checkbox: {shift_data['error']}")
                return 
        check_item = Item('cbox_check')
        check_item.value = value
        receipt = self.create_receipt(check_item)
        if receipt is None:
            return
        res = self.check.create_receipt(receipt)
        if not res:
            error("Помилка при створенні чеку в Checkbox")
            return
        if 'error' in res:
            error(f"При створенні чеку в Checkbox: {res['error']}")
            return
        check_item.value["checkbox_uid"] = res['id']
        err = check_item.save()
        if err:
            error(err)
            return
        self.main_table.reload()

    # add_cashless_sum used if combined cash and cashless, 
    # is_cash must be True, add_cashless_sum included in cash_sum
    def create_receipt(self, check:Item, add_cashless_sum: float=0.0):
        payments = []
        if check.value['is_cash']:
            if add_cashless_sum:
                payments.append({
                  "type": 'CASHLESS',
                    "value": int(add_cashless_sum * 100),
                    "label": 'Картка',
                })    
            payments.append({
                "type": 'CASH',
                "value": int((check.value['cash_sum'] - add_cashless_sum) * 100),
                "label": 'Готівка',
            })
        else:
            payments.append({
                "type": 'CASHLESS',
                "value": int(check.value['cash_sum'] * 100),
                "label": 'Картка',
            })
        
        receipt = {
            "department": "Копіцентр",
            "goods": [],
            "payments": payments,
        }
        check_item = Item('item_to_cbox_check')
        err = check_item.get_filter_w('cbox_check_id', check.value['id'])
        if err:
            error(err)
            return

        for check_item_value in check_item.values:
            good = {
                        "good": {
                            "code": check_item_value['item_code'],
                            "name": check_item_value['name'],
                            "price": int(check_item_value['price'] * 100),
                        },
                        "quantity": int(check_item_value['number']*1000),
                    }
            
            
            if check_item_value['discount'] >= 0.01:
                dsc = round(check_item_value['discount']*check_item_value['number'], 0)
                good["discounts"] = self.make_discount(dsc)
            receipt['goods'].append(good)
        
        if check.value['discount']  >= 0:
            receipt["discounts"] = self.make_discount(check.value['discount'])
            
        return receipt
    
    def make_discount(self, value):
        if value < 0:
            value = -value
            dsct_type = 'EXTRA_CHARGE'
            dsct_name = "Націнка"
        else:
            dsct_type = 'DISCOUNT'
            dsct_name = "Знижка"
        
        return [
            {
                "type": dsct_type,
                "mode": "VALUE",
                "value": value*100,
                "name": dsct_name
            },
        ]
     
    def get_check_fs_code(self, check_value):
        c = self.check.get_check(check_value['checkbox_uid'])
        if 'error' in c:
            error(f"При отриманні чеку з Checkbox: {c['error']}")
            return False
        if not c['fiscal_code']:
            error('Не отримали податковий номер')
            return False
        else:
            check_value["fs_uid"] = c['fiscal_code']
            check_value["created_at"] = c['fiscal_date']
        cbox_check = Item('cbox_check')
        cbox_check.value = check_value
        err = cbox_check.save()
        if err:
            error(err)
            return False
        return True

    def get_check_png(self, check_uid):
        png = ''
        png = self.check.get_png(check_uid)
        if not png:
            return False
        if type(png) == str:
            error(png)
            return False
        l = QLabel()
        lpng = png.scaledToHeight(700)
        l.setPixmap(lpng)
        dlg = CustomDialog(l, 'Друкувати чек?')
        if dlg.exec():
            self.print_check(png)
        return True
    
    def print_check(self, png):
        pinfo = QtPrintSupport.QPrinterInfo.printerInfo('58mm Series Printer')
        printer = QtPrintSupport.QPrinter(pinfo, QtPrintSupport.QPrinter.PrinterMode.HighResolution)
        painter = QPainter()
        painter.begin(printer)
        png = png.scaledToWidth(printer.width())
        painter.drawPixmap(0, 0, png)
        painter.end()