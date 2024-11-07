import requests


class Data:
    def __init__(self, cfg):
        self.base_url = f'http://{cfg["host"]}:{cfg["port"]}/'
        self.cookie_jar = None
        self.wsconn = None

    def signout(self):
        url = self.base_url + 'logout'
        response = requests.get(url, cookies=self.cookie_jar)
        if response.status_code != 200:
            return {'error': 'Помилка при виході, спробуйте пізніше.', 'value': None}
        self.cookie_jar = None
        return {'error': None, 'value': True}

    def signing(self, login, password):
        url = self.base_url + 'login'
        print(url)
        response = requests.post(url, json={'login': login, 'password': password})
        print(response.status_code)
        print(response.text)
        if response.status_code != 200:
            return {'error': 'Помилка при авторизації, спробуйте ще.', 'value': None}
        cookie = response.cookies['sess']
        self.cookie_jar = requests.cookies.RequestsCookieJar()
        self.cookie_jar.set('sess', cookie)
        return response.json()

    def get_models(self) -> dict:
        url = f'{self.base_url}models'
        print('GET MODELS', url)
        response = requests.get(url, cookies=self.cookie_jar)
        return self.format_response(response)

        # -----------------------------------------------------
    def send_ws_message(self, msg):
        if not self.wsconn:
            return {'error': "Не створено з'єднання чату" , 'value': None}
        res = self.wsconn.sendTextMessage(msg)
        if not res:
            return {'error': 'Не вдалося відправити повідомлення', 'value': None}
        return {'error': '', 'value': res}

        
        # -----------------------------------------------------

    def get(self, model_name: str, id: int) -> dict:
        url = f'{self.base_url}{model_name}/{id}'
        print('GET', url)
        response = requests.get(url, cookies=self.cookie_jar)
        return self.format_response(response)
    
    def get_product_deep(self, id: int) -> dict:
        url = f'{self.base_url}product_deep/{id}'
        print('GET', url)
        response = requests.get(url, cookies=self.cookie_jar)
        return self.format_response(response)

    def get_all(self, model_name: str, all: str='') -> dict:
        url = f'{self.base_url}{model_name}_get_all'
        if all:
            url += f'?all={all}'
        print('GET_ALL', url)
        response = requests.get(url, cookies=self.cookie_jar)
        return self.format_response(response)

    def get_filter(self, model_name: str, filter_field: str, filter_value: str|int|bool, all: str='') -> dict:
        filter_type = 'int' if type(filter_value) is int else 'str'
        url = f'{self.base_url}{model_name}_filter_{filter_type}/{filter_field}/{filter_value}'
        if all:
            url += f'?all={all}'
        print('GET_FILTER', url)
        response = requests.get(url, cookies=self.cookie_jar)
        return self.format_response(response)

    def get_between(self, model_name: str, field: str, value1: str|int|bool, value2: str|int|bool, all:str='') -> dict:
        url = f'{self.base_url}{model_name}_between_{field}/{value1}/{value2}'
        if all:
            url += f'?all={all}'
        print('GET_BETWEEN', url)
        response = requests.get(url, cookies=self.cookie_jar)
        return self.format_response(response)

    def get_sum_before(self, model_name: str, sum_field: str, field: str, id: int, date_before: str) -> dict:
        url = f'{self.base_url}{model_name}_of_{sum_field}_sum_before/{field}/{id}/{date_before}'
        print('GET_SUM_BEFORE', url)
        response = requests.get(url, cookies=self.cookie_jar)
        return self.format_response(response)

    def get_sum_filter(self, model_name: str, field1: str, id1: int, field2: str, id2: int) -> dict:
        url = f'{self.base_url}{model_name}_sum_filter_by/{field1}/{id1}/{field2}/{id2}'
        print('GET_SUM_FILTER', url)
        response = requests.get(url, cookies=self.cookie_jar)
        return self.format_response(response)

    # "find": [{"project": ["info"], "contragent":["-search"], "contact":["-search"]}]
    # "/find_project_project_info_contragent_no_search_contact_no_search/{fs}"
    def get_find(self, model_name: str, find: list, find_value: str|int) -> dict:
        find_pref = ''
        for k, v in find.items():
            find_pref += f"_{k}"
            for i in v:
                find_pref += f"_{'no_' + i[1:] if i.startswith('-') else i}"
        url = f'{self.base_url}find_{model_name}{find_pref}/{find_value}'
        print('GET_FIND', url)
        response = requests.get(url, cookies=self.cookie_jar)
        return self.format_response(response)

    def delete(self, model_name: str, id: int) -> dict:
        url = f'{self.base_url}{model_name}/{id}'
        print('DELETE', url)
        response = requests.delete(url, cookies=self.cookie_jar)
        return self.format_response(response)
    
    def realized(self, model_name: str, id: int) -> dict:
        url = f'{self.base_url}realized/{model_name}/{id}'
        print('REALIZED', url)
        response = requests.get(url, cookies=self.cookie_jar)
        return self.format_response(response)
    
    def unrealize(self, model_name: str, id: int) -> dict:
        url = f'{self.base_url}unrealize/{model_name}/{id}'
        print('UNREALIZED', url)
        response = requests.get(url, cookies=self.cookie_jar)
        return self.format_response(response)

    def create(self, model_name: str, data: dict) -> dict:
        url = f'{self.base_url}{model_name}'
        print('CREATE', url)
        response = requests.post(url, json=data, cookies=self.cookie_jar)
        return self.format_response(response)

    def update(self, model_name: str, id: int, data: dict) -> dict:
        url = f'{self.base_url}{model_name}/{id}'
        print('UPDATE', url)
        response = requests.put(url, json=data, cookies=self.cookie_jar)
        return self.format_response(response)

    def create_default(self, data: dict) -> dict:
        url = f'{self.base_url}product_to_ordering_default'
        print('CREATE_DEFAULT', url)
        response = requests.post(url, json=data, cookies=self.cookie_jar)
        return self.format_response(response)

    def format_response(self, response: requests.Response) -> dict:
        if response.status_code == 404:
            return {'error': 'Page not found', 'value': None}
        if response.status_code == 405:
            return {'error': 'Неприпустимий метод', 'value': None}
        j = response.json()
        #print("conn>> ", j)
        return j

# --------------------------WWWWWWWWWWWWWWW----------------------------

    def get_w(self, model_name: str, id: int) -> dict:
        url = f'{self.base_url}w_{model_name}/{id}'
        print('WGET', url)
        response = requests.get(url, cookies=self.cookie_jar)
        return self.format_response(response)

    def get_all_w(self, model_name: str, all: str="") -> dict:
        url = f'{self.base_url}w_{model_name}_get_all'
        if all:
            url += f'?all={all}'
        print('WGET_ALL', url)
        response = requests.get(url, cookies=self.cookie_jar)
        return self.format_response(response)

    def get_filter_w(self, model_name: str, filter_field:str, filter_value: str|int|bool, all: str='') -> dict:
        filter_type = 'int' if type(filter_value) is int else 'str'
        url = f'{self.base_url}w_{model_name}_filter_{filter_type}/{filter_field}/{filter_value}'
        if all:
            url += f'?all={all}'
        print('WGET_FILTER', url)
        response = requests.get(url, cookies=self.cookie_jar)
        return self.format_response(response)

    def get_between_w(self, model_name: str, field: str, value1: str|int|bool, value2: str|int|bool, all: str='') -> dict:
        url = f'{self.base_url}w_{model_name}_between_{field}/{value1}/{value2}'
        if all:
            url += f'?all={all}'
        print('WGET_BETWEEN', url)
        response = requests.get(url, cookies=self.cookie_jar)
        return self.format_response(response)

    def get_between_up_w(self, model_name: str, field: str, value1: str|int|bool, value2: str|int|bool, all: str='') -> dict:
        url = f'{self.base_url}w_{model_name}_between_up_{field}/{value1}/{value2}'
        if all:
            url += f'?all={all}'
        print('WGET_BETWEEN_UP', url)
        response = requests.get(url, cookies=self.cookie_jar)
        return self.format_response(response)

    # "find": [{"project": ["info"], "contragent":["-search"], "contact":["-search"]}]
    # "/find_project_project_info_contragent_no_search_contact_no_search/{fs}"
    def get_find_w(self, model_name: str, find: list, find_value: str|int|bool) -> dict:
        find_pref = ''
        for k, v in find.items():
            find_pref += f"_{k}"
            for i in v:
                find_pref += f"_{'no_' + i[1:] if i.startswith('-') else i}"
        url = f'{self.base_url}w_find_{model_name}{find_pref}/{find_value}'
        print('WGET_FIND', url)
        response = requests.get(url, cookies=self.cookie_jar)
        return self.format_response(response)

    def create_dir(self, id: int) -> dict:
        url = f'{self.base_url}project_dirs/{id}'
        print('CREATE_DIR', url)
        response = requests.get(url, cookies=self.cookie_jar)
        return self.format_response(response)
    
    def copy_project_dir(self, id: int) -> dict:
        url = f'{self.base_url}copy_project/{id}'
        print('COPY_DIR', url)
        response = requests.get(url, cookies=self.cookie_jar)
        return self.format_response(response)

    def create_base_backup(self, name: str) -> dict:
        url = f'{self.base_url}copy_base/{name}'
        print('CREATE_BCKP', url)
        response = requests.get(url, cookies=self.cookie_jar)
        return self.format_response(response)

    def get_base_backups(self) -> dict:
        url = f'{self.base_url}get_bases'
        print('GET_BCKPS', url)
        response = requests.get(url, cookies=self.cookie_jar)
        return self.format_response(response)
    
    def restore_base(self, name: str) -> dict:
        url = f'{self.base_url}restore_base/{name}'
        print('RESTORE_BCKP', url)
        response = requests.get(url, cookies=self.cookie_jar)
        return self.format_response(response)
    
    def delete_base(self, name: str) -> dict:
        url = f'{self.base_url}delete_base/{name}'
        print('DELETE_BCKP', url)
        response = requests.get(url, cookies=self.cookie_jar)
        return self.format_response(response)
