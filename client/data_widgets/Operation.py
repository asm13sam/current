from PyQt6.QtWidgets import (
    QWidget,
    QHBoxLayout,
    QPushButton,
    )

from data.model import Item
from widgets.Dialogs import error, askdlg
from widgets.Form import ItemTable


class OperationTab(ItemTable):
    def __init__(self):
        super().__init__('operation', 'name')
        hbox = QHBoxLayout()
        recalc_widget = QWidget()
        recalc_widget.setLayout(hbox)
        hbox.setContentsMargins(0,0,0,0)
        self.inner.addWidget(recalc_widget)
        self.inner.setStretchFactor(self.inner.count(), 1)
        hbox.addStretch()
        up_cost_btn = QPushButton("Оновити відпускні ціни")
        up_price_btn = QPushButton("Оновити тариф")
        up_am_btn = QPushButton("Оновити амортизацію")
        hbox.addWidget(up_cost_btn)
        hbox.addWidget(up_price_btn)
        hbox.addWidget(up_am_btn)
        up_cost_btn.clicked.connect(lambda: self.update_pricing('cost'))
        up_price_btn.clicked.connect(lambda: self.update_pricing('price'))
        up_am_btn.clicked.connect(lambda: self.update_pricing('equipment_price'))

    def update_pricing(self, field):
        res = askdlg("Вкажіть відсоток")
        try:
            persent = int(res)
        except:
            return
        if not persent:
            return
        res = askdlg("Вкажіть точність до якого знаку після коми")
        try:
            approx = int(res)
        except:
            return
        
        item = Item('operation')
        values = self.table.table.get_selected_values()
        if not values:
            error("Оберіть позиції")
            return
        for v in values:
           v[field] = round(v[field] + v[field] * persent / 100, approx)
           item.value = v
           err = item.save()
           if err:
               error(f'Не можу оновити {v["name"]}[{v["id"]}]:\n {err}')
               continue 