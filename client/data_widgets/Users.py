from PyQt6.QtWidgets import QLabel

from widgets.Dialogs import error
from widgets.Form import DetailsItemTable


class UserOrderingDetailsTable(DetailsItemTable):
    def __init__(self, values: list = None):
        fields = [
            "id",
            "ordering",
            "is_done",
            "user_sum",
            "operation",
            "number",
            "price",
            "comm",
        ]
        buttons = {'reload':'Оновити'}
        super().__init__('operation_to_ordering', '', fields, values, buttons, show_period=True)
        user_sum_caption = QLabel("Загалом")
        self.table.hbox.insertWidget(1, user_sum_caption)
        self.user_sum = QLabel("0.00")
        self.table.hbox.insertWidget(2, self.user_sum)
    
    def period_changed(self, date_from, date_to):
        err = self.item.get_between_up_w('created_at', date_from, date_to)
        if err:
            error(err)
            return
        self.table.table.reload(self.item.values)

    def get_values(self):
        date_from, date_to = self.period.get_period()
        err = self.item.get_between_up_w('created_at', date_from, date_to)
        if err:
            error(err)
            return
        values = self.item.values
        return values
    
    def reload(self, values=None):
        super().reload(values)
        s = round(self.calc_sum('user_sum'), 2)
        self.user_sum.setText(str(s))

    def calc_sum(self, field: str):
        res = 0
        for i in range(self.table.table._model.rowCount()):
            row = self.table.table._model.get_row_value(i)
            if row['is_done']:
                res += row[field]
        return res
    



