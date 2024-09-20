# -*- coding: utf-8 -*-

from PyQt6.QtGui import QPixmap
import requests
import json
#from PIL import Image

from io import BytesIO

# ddiz61578@gmail.com
# Nhfrnjh55



# BASE_URL = 'https://api.checkbox.ua/api/v1/'

class CheckBoxFS:
    def __init__(self, base_url, license_key='', cashier_pin=''):
        self.base_url = base_url
        self.license_key = license_key
        self.cashier_pin = cashier_pin
        self.auth_key = ''
        self.session = requests.Session()
        self.set_session_default()
        self.is_signed = False

    def set_session_default(self):
        self.session.headers.update({
            'accept': 'application/json',
            'X-Client-Name': 'TgtCopyCenter',
            'X-Client-Version': '0.1.1',
            'X-License-Key': self.license_key,
            'Content-Type': 'application/json',
        })
        # if self.auth_key:
        #     self.session.headers.update({'Authorization': self.auth_key})

    def close(self):
        self.session.close()

    # Авторизація.
    def sign_in(self):
        with open ('check.json', "r") as f:
            cfg = f.read()
        if len(cfg) > 5:
            cfg = json.loads(cfg)
            self.auth_key = "Bearer " + cfg['access_token']
            self.session.headers.update({'Authorization': "Bearer " + cfg["access_token"]})
        else:
            url = self.base_url + 'cashier/signinPinCode'

            payload = {
                "pin_code": self.cashier_pin,
            }

            r = self.session.post(url, data=json.dumps(payload))
            mess = r.json()
            print(mess)
            if r.status_code == requests.codes.ok:
                self.auth_key = "Bearer " + mess["access_token"]
                with open ('check.json', "w") as f:
                    f.write(json.dumps(mess))
            else:
                self.is_signed = False
                if 'message' in mess:
                    return{'error': mess['message']}
                return {}
            self.session.headers.update({'Authorization': "Bearer " + mess["access_token"]})
        self.is_signed = True
        return {'signed': True}

    # Попередній вибір режиму створення транзакцій (онлайн/офлайн).
    # Отримання інформації про стан каси.
    def get_cash_state(self):
        url = self.base_url + 'cash-registers/info'
        r = self.session.get(url)
        if r.status_code > 299:
            res = r.json()
            print(res)
            if 'message' in res:
                return{'error': res['message']}
            return {}
        return r.json()

    # Перехід у потрібний режим роботи (онлайн/офлайн).
    # 4.2. Ініціалізація переходу ПРРО у онлайн режим.

    def work_online(self):
        url = self.base_url + 'cash-registers/go-online'
        r = self.session.post(url)

        if r.status_code != requests.codes.ok:
            return {}
        return r.json()

    # Перевірка зв'язку з сервером ДПС.
    def ping_FS(self):
        url = self.base_url + 'cash-registers/ping-tax-service'

        r = self.session.post(url)
        res = r.json()
        if r.status_code > 299:
            # print(r.status_code)
            # print(res)
            return False
        # print(res)
        return res['status']

        

    # 4.3. Отримання нових офлайн кодів.
    # Отримання нових офлайн фіскальних кодів з сервера ДПС
    def ask_offline_codes(self, count=2000, sync=True):
        sync_str = 'true' if sync else 'false'
        url = f'cash-registers/ask-offline-codes?count={count}&sync={sync_str}'
        url = self.base_url + url
        r = self.session.get(url)

        # print(r.status_code)
        # print(r.status_code == requests.codes.ok)
        # print(r.json())

    # Отримання списку невикористаних офлайн фіскальних кодів з сервера Checkbox
    def get_offline_codes(self, count=1000):
        url = self.base_url + f'cash-registers/ask-offline-codes?count={count}'

        r = self.session.get(url)

        # print(r.status_code)
        # print(r.status_code == requests.codes.ok)
        # print(r.json())

    # Ініціалізація переходу ПРРО у офлайн режим.
    def go_offline(self, date='', fiscal_code=''):
        url = self.base_url + 'cash-registers/go-offline'

        # payload = {
            # "go_offline_date": "2022-06-03T14:26:45+03:00",
            # "fiscal_code": "lRvyZSJLrC4",
        # }

        payload = {
            "go_offline_date": date,
            "fiscal_code": fiscal_code,
        }


        r = self.session.post(url, data=json.dumps(payload))
        mess = r.json()
        # print(mess)

    # Відкриття зміни.
    def shifts(self):
        url = self.base_url + 'shifts'

        # payload = {
            # "id": "",
            # "fiscal_code": "",
            # "fiscal_date": "",
        # }

        r = self.session.post(url)
        if r.status_code > 299:
            res = r.json()
            print(res)
            if 'message' in res:
                return{'error': res['message']}
            return {}
        return r.json()

    def shifts_offline(self):
        url = self.base_url + 'shifts'

        # payload = {
            # "id": "унікальний ідентифікатор зміни у форматі UUID",
            # "fiscal_code": "<фіскальний код>",
            # "fiscal_date": "<фіскальна дата у форматі ISO 8601 за шаблоном YYYY-MM-DDThh:mm:ss.ssssss±hh:mm>",
        # }

        payload = {
            "id": "унікальний ідентифікатор зміни у форматі UUID",
            "fiscal_code": "<фіскальний код>",
            "fiscal_date": "<фіскальна дата у форматі ISO 8601 за шаблоном YYYY-MM-DDThh:mm:ss.ssssss±hh:mm>",
        }


        r = self.session.post(url, data=json.dumps(payload))
        mess = r.json()
        # print(mess)

    # Перевірка статусу поточної зміни
    def get_shift(self):
        url = self.base_url + 'cashier/shift'

        r = self.session.get(url)

        if r.status_code > 299:
            # print(r.json)
            return {}
        return r.json()

    # Створення службового чеку внесення готівки (за необхідності).
    def cash_service(self, summa):
        url = self.base_url + 'receipts/service'

        # payload = {
            # "payment": {
                # "type": "CASH",
                # "value": <сума у копійках, для створення чеку службового вилучення перед сумою має бути - >,
                # "label": "Готівка"
            # }
        # }

        payload = {
            "payment": {
                "type": "CASH",
                "value": summa,
                "label": "Готівка",
            },
        }


        r = self.session.post(url, data=json.dumps(payload))
        mess = r.json()
        # print(mess)

    # Створення чеку/чеків.
    def create_receipt(self, receipt):
        url = self.base_url + 'receipts/sell'

        r = self.session.post(url, data=json.dumps(receipt))
        # print('status', r.status_code)
        if r.status_code > 299:
            res = r.json()
            print(res)
            if 'message' in res:
                return{'error': res['message']}
            return {}
        return r.json()
    
    # Створення чеку передоплати.
    def create_pre_receipt(self, receipt):
        url = self.base_url + 'prepayment-receipts'
        r = self.session.post(url, data=json.dumps(receipt))
        # print('status', r.status_code)
        if r.status_code > 299:
            res = r.json()
            if 'message' in res:
                return{'error': res['message']}
            return {}
        return r.json()
    
    # Створення чеку післяплати.
    def create_post_receipt(self, receipt, relation_id):
        url = self.base_url + f'prepayment-receipts/{relation_id}'

        r = self.session.post(url, data=json.dumps(receipt))
        # print('status', r.status_code)
        if r.status_code > 299:
            res = r.json()
            if 'message' in res:
                return{'error': res['message']}
            return {}
        return r.json()
    
    # Створення чеку/чеків from json
    def create_receipt_json(self, receipt):
        url = self.base_url + 'receipts/sell'

        r = self.session.post(url, data=receipt)
        if r.status_code > 299:
            res = r.json()
            if 'message' in res:
                return{'error': res['message']}
            return {}
        return r.json()
    
#/api/v1/receipts/{receipt_id}
    def get_check(self, receipt_id):
        url = self.base_url + f'receipts/{receipt_id}'

        r = self.session.get(url)
        if r.status_code > 299:
            res = r.json()
            # print(res)
            if 'message' in res:
                return{'error': res['message']}
            return {}
        return r.json()
        
#/api/v1/receipts/{receipt_id}/png
    def get_png(self, receipt_id):
        url = self.base_url + f'receipts/{receipt_id}/png'

        r = self.session.get(url)
        # r = self.session.get(url, timeout=5)
        
        if r.status_code > 299:
            # if r.status_code == 424:
            #     res = self.ping_FS()
            #     print(res)
            #     if res:
            #         return self.get_png(receipt_id)
            # print(r.json())
            return None
        pixmap = QPixmap()
        res = pixmap.loadFromData(r.content)
        if res:
            return pixmap
        
        return None

    def get_html(self, receipt_id):
        url = self.base_url + f'receipts/{receipt_id}/html'

        r = self.session.get(url)

        # print(r.status_code)
        # print(r.status_code == requests.codes.ok)
        #i = Image.open(BytesIO(r.content))
        r.encoding = 'utf-8'
        # print(r.text)
        with open ('check.html', "w") as f:
            f.write(r.text)

    def get_txt(self, receipt_id):
        url = self.base_url + f'receipts/{receipt_id}/text'

        r = self.session.get(url)
        if r.status_code > 299:
            if r.status_code == 424:
                res = self.ping_FS()
                # print(res)
                if res:
                    return self.get_txt(receipt_id)
            # print(r.json())
            return ''
        
        # print(r.status_code)
        # print(r.status_code == requests.codes.ok)
        #i = Image.open(BytesIO(r.content))
        r.encoding = 'utf-8'
        return r.text
        # print(r.text)
        # with open ('check.txt', "w") as f:
        #     f.write(r.text)

    def shift_close(self):
        url = self.base_url + 'shifts/close'

        r = self.session.post(url)
        if r.status_code > 299:
            res = r.json()
            # print(res)
            if 'message' in res:
                return{'error': res['message']}
            return {}
        self.is_signed = False
        return r.json()

# {
  # "id": "<UUID чека>",
  # "cashier_name": "<Ім'я касира>",
  # "departament": "<Назва відділу>",
  # "goods":
  # [
    # {
      # "good":
       # {
        # "code": "<Код товару>",
        # "name": "<Назва товару>",
        # "barcode": "<Штрих-код товару>",
        # "excise_barcode": "<цифрове позначення штрих-коду акцизної марки>",
        # "excise_barcodes": ["<цифрове позначення штрих-коду акцизної марки 1>","<цифрове позначення штрих-коду акцизної марки 2>"],
        # "header": "<Хедер товару 1>",
        # "footer": "<Футер товару 1>",
        # "price": <ціна  у копійках>,
        # "tax": [<цифровий або літерний код ставки податку (попередньо програмується у особистому кабінеті). Якщо до товару потрібно застосувати декілька податків - вказати через кому>],
        # "uktzed": "<код УКТЗЕД>"
      # },
      # "good_id": "<UUID товару>",
      # "quantity": <кількість у тисячах, 1 шт = 1000>,
      # "is_return": <флаг true/false, що визначає, чи це чек повернення>,
      # "discounts":
      # [
       # {
       # "type": "<тип знижки - "DISCOUNT"/"EXTRA_CHARGE" (ЗНИЖКА/НАДБАВКА)>",
       # "mode": "<режим знижки "VALUE"/"PERCENT" (АБСОЛЮТНЕ ЗНАЧЕННЯ/ВІДСОТКОВА ЗНИЖКА)>",
       # "value": <значення знижки>,
       # "tax_code": [<код податку, який застосовується для товару. Потрібно вказувати для коректного обрахунку знижки, якщо товар має податкову ставку>],
       # "tax_codes": [<коди податкових ставок, що застосовуються для товару (якщо їх >1). Потрібно вказувати через кому для коректного обрахунку знижки, якщо товар має податкову ставку>],
       # "name": "<назва знижки або надбавки>"
       # }
      # ]
   # },
   # {
     # "good":
     # {<блок з даними про другий товар, за структурою аналогічний попередньому. блоки good потрібно повторювати стільки разів, скільки у вас товарів у чеку>}
   # }
  # ],
  # "delivery":
  # {
    # "email": "<e-mail клієнта для відправки копії чека>",
    # "emails":
    # ["<e-mail клієнта для відправки копії чека 1>","<e-mail клієнта для відправки копії чека 2>"],
    # "phone": "<номер телефона клієнта для відправки копії чека по SMS/Viber (для роботи функції має бути налаштована та підключена відповідна послуга)>. Формат 380..."
  # },
  # "discounts":
      # [
       # {
       # "type": "<тип знижки - "DISCOUNT"/"EXTRA_CHARGE" (ЗНИЖКА/НАДБАВКА)>",
       # "mode": "<режим знижки "VALUE"/"PERCENT" (АБСОЛЮТНЕ ЗНАЧЕННЯ/ВІДСОТКОВА ЗНИЖКА)>",
       # "value": <значення знижки>,
       # "tax_code": [<код податку, який застосовується для товару. Потрібно вказувати для коректного обрахунку знижки, якщо товар має податкову ставку>],
       # "tax_codes":
        # [<коди податкових ставок через кому, що застосовуються для товару (якщо їх >1). Потрібно вказувати для коректного обрахунку знижки, якщо товар має податкову ставку>],
       # "name": "<назва знижки або надбавки>"
       # }
      # ]
    # },
  # "payments":
  # [
   # {
   # "type": "<"CASH"/"CASHLESS" (ГОТІВКА/БЕЗГОТІВКОВИЙ РОЗРАХУНОК (картка, сертифікати, бонуси тощо))>",
   # "pawnshop_is_return": <true/false Ознака того, що даний чек є видатковим чеком ломбарду. Для звичайного чека параметр вказувати не потрібно>,
   # "value": <сума оплати у копійках>,
   # "label": "<текстова назва форми оплати>",
   # "code": <номер оплати (тільки для безготівкових платежів)>,
   # "commission": <комісія (тільки для безготівкових платежів)>,
   # "card_mask": "<маска карти (не більше 19 символів) (тільки для безготівкових платежів)>",
   # "bank_name": "<назва банку (тільки для безготівкових платежів)>",
   # "auth_code": "<код авторизації банківської операції (тільки для безготівкових платежів)>",
   # "rrn": "<Reference Retrieval Number - унікальний ідентифікатор банківської транзакції (тільки для безготівкових платежів)>",
   # "payment_system": "<платіжна система (тільки для безготівкових платежів)>",
   # "owner_name": "<ім'я власника електронного платіжного засобу (тільки для безготівкових платежів)>",
   # "terminal": "<інформація про платіжний термінал (тільки для безготівкових платежів)>",
   # "acquirer_and_seller": "<ідентифікатор еквайра та торгівця, або інші реквізити, що дають змогу їх ідентифікувати (тільки для безготівкових платежів)>",
   # "receipt_no": "<номер банківського чека (тільки для безготівкових платежів)>",
   # "signature_required": <true/false флаг, який визначає, чи має бути доступною графа для підпису власника картки та касира>
   # },
    # {
       # <блок з даними по додатковій формі оплати за шаблоном, який описаний вище (якщо в чеку декілька форм оплати)>
    # }
  # ],
  # "rounding": <true/false активація режиму заокруглення (для його роботи у блоці payments має бути вказана сума вже заокруглена за правилами НБУ>,
  # "header": "<хедер чека>",
  # "footer": "<футер чека>",
  # "barcode": "<штрих-код чека>",
  # "order_id": "<UUID замовлення (вказується у випадку роботи з API у режимі замовлень)>",
  # "related_receipt_id": "<UUID пов'язаного чека>",
  # "previous_receipt_id": "<UUID попереднього чека>",
  # "technical_return": <true/false флаг, яким можна позначити, що даний чек є чеком технічного (помилкового) повернення>,
  # "context": {
    # "additionalProp1": "<додаткова властивість 1>",
    # "additionalProp2": "<додаткова властивість 2>",
    # "additionalProp3": "<додаткова властивість 3>"
  # },
  # "is_pawnshop": <true/false флаг, яким можна позначити, що даний чек є чеком ломбарду>,
  # "custom": {
    # "html_global_header": "<глобальний хедер для чеків html/pdf візуалізацій>",
    # "html_global_footer": "<глобальний футер для чеків html/pdf візуалізацій>",
    # "html_body_style": "<фон сторінки з чеком>",
    # "html_receipt_style": "<стиль блоку з чеком>",
    # "html_ruler_style": "<стиль роздільника з зірочками між блоками чеку",
    # "html_light_block_style": "<стиль світлих блоків, це весь підвал чеку та клітинки з вартістю та кількістю>",
    # "text_global_header": "<глобальний хедер для чеків png/txt візуалізацій>",
    # "text_global_footer": "<глобальний футер для чеків png/txt візуалізацій>"
  # }
# }



# Перевірка статусу чеку.
# Відправка чека у ДПС (якщо чек створений у офлайн режимі та ще не був доставлений у податкову).
# Створення службового чеку винесення готівки (за необхідності).
# Перевірка балансу зміни за допомогою X-Звіту (за необхідності). Закриття зміни (створення Z - звіту).
# Перевірка статусу каси (онлайн/офлайн).
# Переведення каси у онлайн режим (якщо каса у офлайні), що означатиме, що всі офлайн транзакції успішно передані до податкової.
# Завершення сесії (за необхідності).




# cb = CheckBoxFS(BASE_URL, 'test3e3f6dc789547d23a26003a5', "9207718872", a_key)
# cb.sign_in()
# print(cb.session.headers)
# print('Отримання інформації про стан каси.')
# cb.get_cash_state()
#print('Ініціалізація переходу ПРРО у онлайн режим.')
#cb.work_online()
#print('Перевірка зв`язку з сервером ДПС.')
#cb.ping_FS()
#print('Відкриття зміни.')
#cb.shifts()
# print('Перевірка статусу поточної зміни')
# cb.get_shift()
# print("Створення службового чеку внесення готівки")
# cb.cash_service(50000)
# print('Створення чеку/чеків')
# cb.create_receipt()
# print('get png')
#cb.get_png('2023da54-c39c-4305-bd78-074afc520700')
#cb.shift_close()
# cb.get_txt('2023da54-c39c-4305-bd78-074afc520700')

