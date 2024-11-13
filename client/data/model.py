import json
from barcode import EAN13
from datetime import datetime, date
from typing import Self

from data.app import App

class Item:
    def __init__(self, name: str):
        app = App()
        self.name = name
        self.hum: str = app.model['models'][name]['hum']
        # print('hum', self.hum)
        self.find: list = app.model['models'][name]['find'] if 'find' in app.model['models'][name] else []
        self.op_model: dict = app.model['models'][name]
        self.repo = app.repository
        self.model: dict = app.model['models'][name]['model']
        # print('model', self.model)
        self.value: dict = {}
        self.values: list = []
        self.model_w: dict = app.model['models'][name]['w_model']
        self.model_w.update(self.model)
        # print('model_w', self.model_w)
        self.columns: list = app.model['models'][name]['columns']
        if app.user:
            self.user_id = app.user['id']
            self.username = app.user['name']

    # return error str if error, else empty string
    # value in self.value(s)

    def get(self, id: int) -> str:
        res = self.repo.get(self.name, id)
        return self.process_result(res)

    def get_product_deep(self, id: int) -> str:
        res = self.repo.get_product_deep(id)
        return self.process_result(res)

    def get_all(self) -> str:
        res = self.repo.get_all(self.name)
        return self.process_result(res)

    def get_all_with_deleted(self) -> str:
        res = self.repo.get_all(self.name, all='all')
        return self.process_result(res)

    def get_all_deleted_only(self) -> str:
        res = self.repo.get_all(self.name, all='deleted')
        return self.process_result(res)

    def get_filter(self, field: str, value: str|int|bool) -> str:
        res = self.repo.get_filter(self.name, field, value)
        return self.process_result(res)

    def get_filter_with_deleted(self, field: str, value: str|int|bool) -> str:
        res = self.repo.get_filter(self.name, field, value, all='all')
        return self.process_result(res)

    def get_filter_deleted_only(self, field: str, value: str|int|bool) -> str:
        res = self.repo.get_filter(self.name, field, value, all='deleted')
        return self.process_result(res)

    def get_between(self, field: str, value1: str|int|bool, value2: str|int|bool) -> str:
        res = self.repo.get_between(self.name, field, value1, value2)
        return self.process_result(res)

    def get_between_with_deleted(self, field: str, value1: str|int|bool, value2: str|int|bool) -> str:
        res = self.repo.get_between(self.name, field, value1, value2, all='all')
        return self.process_result(res)

    def get_between_deleted_only(self, field: str, value1: str|int|bool, value2: str|int|bool) -> str:
        res = self.repo.get_between(self.name, field, value1, value2, all='deleted')
        return self.process_result(res)

    def get_find(self, find: list, value: str|int|bool) -> str:
        res = self.repo.get_find(self.name, find, value)
        return self.process_result(res)

    def get_sum_before(self, sum_field: str, field: str, id: int, date_before: str) -> str:
        res = self.repo.get_sum_before(self.name, sum_field, field, id, date_before)
        return res

    def get_sum_filter(self, field1: str, id1: int, field2: str='-', id2: int=0) -> str:
        res = self.repo.get_sum_filter(self.name, field1, id1, field2, id2)
        return res

    def save(self) -> str:
        if self.name == 'contragent' or self.name == 'contact':
            self.update_search()

        if self.value['id']:
            return self.update()
        else:
            return self.create()

    def update_search(self):
        search_params = [
                    self.value['name'].lower(),
                    self.value['phone'],
                    self.value['email'].lower(),
                ]

        # Viber +380681317555
        # Telegram +380 98 742 78 18 or @Bondar_Yana
        if 'viber' in self.value:
            search_params.append(self.value['viber'])
        if 'telegram' in self.value:
            search_params.append(self.value['telegram'].lower())
        self.value['search'] = '_'.join(search_params)

    def add_barcode(self, code: str):
        bcode = EAN13('%s%06d' % (code, self.value['id']))
        bcode.build()
        self.value['barcode'] = bcode.get_fullcode()
        res = self.repo.update(self.name, self.value['id'], self.value)
        return self.process_result(res)

    def create(self) -> str:
        if not self.value:
            return 'Неможливо створити пустий елемент.'
        res = self.repo.create(self.name, self.value)
        
        err = self.process_result(res)

        if err:
            return err
        elif self.op_model['message']:
            mess = f"створив {self.hum} {self.value['name'] if 'name' in self.value else ''} [#{self.value['id']}]"
            self.repo.send_ws_message(json.dumps({'user_id': self.user_id, 'username': self.username, 'message': mess}))
        
        if self.name == 'matherial':
            return self.add_barcode('222222')
        if self.name == 'operation':
            return self.add_barcode('222223')
        if self.name == 'product':
            return self.add_barcode('222224')

        return ''

    def update(self) -> str:
        if not self.value:
            return 'Неможливо оновити пустий елемент.'
        res = self.repo.update(self.name, self.value['id'], self.value)
        if not res['error'] and self.op_model['message']:
            mess = f"змінив {self.hum} {self.value['name'] if 'name' in self.value else ''} [#{self.value['id']}]"
            self.repo.send_ws_message(json.dumps({'user_id': self.user_id, 'username': self.username, 'message': mess}))
        return self.process_result(res)

    def create_p2o_defaults(self) -> str:
        res = self.repo.create_default(self.value)
        return self.process_result(res)

    def delete(self, id: int, cause: str='') -> str:
        res = self.repo.delete(self.name, id)
        if not res['error'] and self.op_model['message']:
            mess = f"видалив {self.hum} {self.value['name'] if 'name' in self.value else ''} [#{id}]"
            if cause:
                mess += f'\nПричина: {cause}'
            self.repo.send_ws_message(json.dumps({'user_id': self.user_id, 'username': self.username, 'message': mess}))
        return self.process_result(res)
    
    def unrealize(self, id: int, cause: str='') -> str:
        res = self.repo.unrealize(self.name, id)
        if not res['error']:
            mess = f"відмінив проведення {self.hum} {self.value['name'] if 'name' in self.value else ''} [#{id}]"
            if cause:
                mess += f'\nПричина: {cause}'
            self.repo.send_ws_message(json.dumps({'user_id': self.user_id, 'username': self.username, 'message': mess}))
        return self.process_result(res)
    
    def realize(self, id: int, cause: str='') -> str:
        res = self.repo.realized(self.name, id)
        if not res['error']:
            mess = f"провів {self.hum} {self.value['name'] if 'name' in self.value else ''} [#{id}]"
            if cause:
                mess += f'\nПричина: {cause}'
            self.repo.send_ws_message(json.dumps({'user_id': self.user_id, 'username': self.username, 'message': mess}))
        return self.process_result(res)

    def process_result(self, result: dict) -> str:
        if result['error']:
            return result['error']

        if type(result['value']) == dict:
            self.value = result['value']
            self.values = []
        else:
            self.values = result['value']
            self.value = {}

        return ''

    # ----------------- WWWWWWWWWWWWW --------------------
    def get_w(self, id: int) -> str:
        res = self.repo.get_w(self.name, id)
        return self.process_result(res)

    def get_all_w(self) -> str:
        res = self.repo.get_all_w(self.name)
        return self.process_result(res)

    def get_all_w_with_deleted(self) -> str:
        res = self.repo.get_all_w(self.name, all='all')
        return self.process_result(res)

    def get_all_w_deleted_only(self) -> str:
        res = self.repo.get_all_w(self.name, all='deleted')
        return self.process_result(res)

    def get_filter_w(self, field: str, value: str|int|bool) -> str:
        res = self.repo.get_filter_w(self.name, field, value)
        return self.process_result(res)

    def get_filter_w_with_deleted(self, field: str, value: str|int|bool) -> str:
        res = self.repo.get_filter_w(self.name, field, value, all='all')
        return self.process_result(res)

    def get_filter_w_deleted_only(self, field: str, value: str|int|bool) -> str:
        res = self.repo.get_filter_w(self.name, field, value, all='deleted')
        return self.process_result(res)

    def get_between_w(self, field: str, value1: str|int|bool, value2: str|int|bool) -> str:
        res = self.repo.get_between_w(self.name, field, value1, value2)
        return self.process_result(res)

    def get_between_w_with_deleted(self, field: str, value1: str|int|bool, value2: str|int|bool) -> str:
        res = self.repo.get_between_w(self.name, field, value1, value2, all='all')
        return self.process_result(res)

    def get_between_w_deleted_only(self, field: str, value1: str|int|bool, value2: str|int|bool) -> str:
        res = self.repo.get_between_w(self.name, field, value1, value2, all='deleted')
        return self.process_result(res)

    def get_between_up_w(self, field: str, value1: str|int|bool, value2: str|int|bool) -> str:
        res = self.repo.get_between_up_w(self.name, field, value1, value2)
        return self.process_result(res)

    def get_between_up_w_with_deleted(self, field: str, value1: str|int|bool, value2: str|int|bool) -> str:
        res = self.repo.get_between_up_w(self.name, field, value1, value2, all='all')
        return self.process_result(res)

    def get_between_up_uw_deleted_only(self, field: str, value1: str|int|bool, value2: str|int|bool) -> str:
        res = self.repo.get_between_up_w(self.name, field, value1, value2, all='deleted')
        return self.process_result(res)

    def get_filter_by_doc_w(self, field: str, value: str|int|bool) -> str:
        res = self.repo.get_filter_by_doc_w(self.name, field, value)
        return self.process_result(res)

    def get_find_w(self, find: list, value: str|int|bool) -> str:
        res = self.repo.get_find_w(self.name, find, value)
        return self.process_result(res)

    def create_default(self):
        app = App()
        for k, v in self.model.items():
            if k == 'user_id':
                self.value[k] = app.user['id']
            elif k == 'created_at':
                self.value[k] = datetime.now().isoformat(timespec='seconds')
            elif k == 'deadline_at':
                self.value[k] = date.today().isoformat() + 'T23:59:59'
            else:
                self.value[k] = v['def']

    def create_default_w(self):
        self.create_default()
        for k, v in self.model_w.items():
            if k in self.value:
                continue
            self.value[k] = v['def']

    def cross(self, other_item: Self):
       for k, v in self.model.items():
            if k != 'id' and v['def'] != 'date' and k in other_item.value:
                self.value[k] = other_item.value[k]

    def __str__(self) -> str:
        res = ''
        res += self.name + '\n'
        res += self.hum + '\n'
        if self.value:
            for k, vi in self.value.items():
                res += f'{k} -> {vi}\n'
        if self.values:
            for val in self.values:
                for k, vi in val.items():
                    res += f'{k} -> {vi}\n'
        return res

    def create_dirs(self, id:int) -> str:
        if self.name != 'ordering':
            return self.process_result({'error': "Неможливо створити теку не для замовлення", 'value': None})
        res = self.repo.create_dir(id)
        if res['error']:
            return res['error']
        return ''
    
    def copy_dirs(self, id: int=0) -> str:
        if self.name != 'project':
            return self.process_result({'error': "Неможливо копіювати теку не для проекту", 'value': None})
        res = self.repo.copy_project_dir(self.value['id'] if not id else id)
        if res['error']:
            return res['error']
        return ''

    # def get_url(self, url):
    #     res = self.fdata.get(url)
    #     if 'detail' in res:
    #         return {'error': res['detail'][0]['msg']}
    #     else:
    #         return res
