from data.model import Item
from widgets.Table import Table, TableWControls
from widgets.Dialogs import error
from widgets.Form import PeriodWidget


class DocsTable(TableWControls):
    def __init__(self, 
                 main_key='document_uid',
                 doc_key='based_on',
                 values=[],
                 docs=('cash_in', 'cash_out', 'whs_in', 'whs_out', 'invoice', 'cbox_check'),
                 controls=False,
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
        self.cur_value = {}
        btns = {}
        if controls:
            btns = {'reload':'Оновити'}
        super().__init__(model, table_fields, values, btns)
        self.period = None
        if controls:
            self.period = PeriodWidget()
            self.hbox.insertWidget(0, self.period)
            self.period.periodSelected.connect(lambda:self.reload(self.cur_value))

    def reload(self, value=None):
        values = []
        if not value:
            self.clear()
            self.cur_value = {}
            return
        self.cur_value = value
        for doc_name in self.docs:
            doc = Item(doc_name)
            if not value[self.main_key]:
                continue
            if not self.period:
                err = doc.get_filter_w(self.doc_key, value[self.main_key])
                if err:
                    error(err)
                    continue
            else:
                date_from, date_to = self.period.get_period()
                err = doc.get_between_w('created_at', date_from, date_to)
                if err:
                    error(err)
                    continue
                doc.values = [v for v in doc.values if v[self.doc_key]==value[self.main_key]]
            if not doc.values:
                continue
            for val in doc.values:
                doc_val = {}
                doc_val['id'] = val['id']
                doc_val['name'] = val['name']
                doc_val['created_at'] = val['created_at']
                doc_val['is_realized'] = val['is_realized'] if 'is_realized' in val else True
                doc_val['type'] = doc.name
                doc_val['type_hum'] = doc.hum
                doc_val['comm'] = val['comm'] if 'comm' in val else ''
                doc_val['doc_sum'] = val['whs_sum'] if doc.name.startswith('whs') else val['cash_sum']
                values.append(doc_val) 
        return super().reload(values)
    
    def calc_by_type(self, type_name):
        values = self.values()
        res = 0
        for v in values:
            if v['type'] == type_name:
                res += v['doc_sum'] if v['is_realized'] else 0
        return res