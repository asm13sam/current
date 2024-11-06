from data.model import Item
from widgets.Table import Table
from widgets.Dialogs import error


class DocsTable(Table):
    def __init__(self, 
                 main_key='document_uid',
                 doc_key='based_on',
                 values=[],
                 docs=('cash_in', 'cash_out', 'whs_in', 'whs_out', 'invoice', 'cbox_check'),
                 ):
        table_fields = ['id', 'name', 'type_hum', 'doc_sum', 'created_at', 'comm']
        model = {
            "id": {"def": 0, "hum": "Номер"},
            "name": {"def": "Замовлення", "hum": "Назва"},
            "created_at": {"def": "date", "hum": "Створений"},
            "is_realized": {"def": "date", "hum": "Проведений"},
            "doc_sum": {"def": 0.0, "hum": "Сума"},
            "type": {"def": "", "hum": "Тип code"},
            "type_hum": {"def": "", "hum": "Тип"},
            "comm": {"def": "", "hum": "Коментар"},
        }
        self.main_key = main_key
        self.doc_key = doc_key
        self.docs = docs
        super().__init__(model, table_fields, values)

    def reload(self, value=None):
        values = []
        if not value:
            self.clear()
            return
        for doc_name in self.docs:
            doc = Item(doc_name)
            if not value[self.main_key]:
                continue
            err = doc.get_filter_w(self.doc_key, value[self.main_key])
            if err:
                error(err)
                continue
            if not doc.values:
                continue
            for val in doc.values:
                doc_val = {}
                doc_val['id'] = val['id']
                doc_val['name'] = val['name']
                doc_val['created_at'] = val['created_at']
                doc_val['is_realized'] = val['is_realized']
                doc_val['type'] = doc.name
                doc_val['type_hum'] = doc.hum
                doc_val['comm'] = val['comm'] if 'comm' in val else ''
                doc_val['doc_sum'] = val['whs_sum'] if doc.name.startswith('whs') else val['cash_sum']
                values.append(doc_val) 
        return super().reload(values)
    
    def calc_by_type(self, type_name):
        values = self._model.values()
        res = 0
        for v in values:
            if v['type'] == type_name:
                res += v['doc_sum'] if v['is_realized'] else 0
        return res