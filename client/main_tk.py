#!/usr/bin/python3
# -*- coding: utf-8 -*-

import json

from data.app import App
from data.conn import Data
from data.model import Item

import tkinter as tk
from tkinter.ttk import Style
import ttkbootstrap as ttk
import ttkbootstrap.dialogs.dialogs as dlg
from ttkbootstrap.constants import *

class Table(ttk.Treeview):
    def __init__(self, parent, data_name: str, columns=[], column_widths=[]):
        self.data_model = Item(data_name)
        self.data = {}
        self.columns = (self.data_model.columns + self.data_model.columns_w) if not columns else columns
        super().__init__(parent, columns=self.columns, show='headings', bootstyle=INFO)
        self.set_headers()
        if column_widths:
            self.set_column_widths(column_widths)
    
        self.set_data()

    #for dialogs
    def accept(self):
        return self.get_selected() 
   
    def get_selected(self):
        sel = self.selection()
        if sel:
            return self.data[sel[0]]
        return None
    
    def make_sortable(self, columns=[]):    
        if not columns:
            columns = self.columns
        for column in columns:
            self.heading(column, command=lambda c=column: self.sort(c, False))
    
    def set_headers(self):    
        for column in self.columns:
            if column in self.data_model.model_w:
                header = self.data_model.model_w[column]['hum']
            else:
                header = self.data_model.model[column]['hum']
            self.heading(column, text=header)
            
    def set_column_widths(self, widths):
        for column, width in list(zip(self.columns, widths)):
            self.column(column, width=width, minwidth=width*2//3)
        
    def set_data(self):
        self.clear()
        err = self.data_model.get_all_w()
        if err:
            dlg.Messagebox.show_error(err, "Помилка")
        for v in self.data_model.values:
            self.set_item(v)
            
    def set_item(self, value):
        values = list([value[c] for c in self.columns])
        self.insert("", END, value['id'], values=values)

    def clear(self):
        self.delete(*self.get_children(""))
        
    def reload(self):
        sel_item = self.get_selected()
        self.clear()
        self.set_data()
        if sel_item:
            str_uid = str(sel_item['uid'])
            if self.exists(str_uid):
                self.selection_set(str_uid)

    
    # def make_info(self, key):
    #     res = ''
    #     for k, v in self.data[key].items():
    #         res += k + " -> " + str(v) + "\n"
    #     return res
        
    def sort(self, column, reverse):
        self.data_model.values.sort(key=lambda v: v[column], reverse=reverse)
        self.clear()
        for v in self.data_model.values:
            self.set_item(v)
        #  reverse sort next time
        self.heading(column, command=lambda: self.sort(column, not reverse))

class Tree(ttk.Treeview):
    def __init__(self, parent, group_name: str):
        self.gitem = Item(group_name)
        self.gkey = group_name + '_id'
        self.data = {}
        super().__init__(parent, show='tree', bootstyle=INFO)
        err = self.gitem.get_all()
        if err:
            dlg.Messagebox.show_error(err, "Помилка")
        self.set_data()
        
    def clear(self):
        self.delete(*self.get_children("")) 
    
    def set_data(self):
        for v in self.gitem.values:
            # if v[self.gkey] in self.data:
            #     self.data[self.gkey].append(v)
            # else:
            #     self.data[self.gkey] = [v,]
            self.data[v['id']] = v
        for k in sorted(self.data.keys()):
            if not self.exists(k):
                self.set_item(self.data[k])
    
    def set_item(self, item):
        if item[self.gkey] == 0:
            item[self.gkey] = ""
        str_guid = str(item[self.gkey])
        if not self.exists(str_guid):
            self.set_item(self.data[item[self.gkey]])
        self.insert(item[self.gkey], 'end', item['id'], text=item['name'])
    
    def reload(self):
        sel_item = self.get_selected()
        self.clear()
        if sel_item:
            str_uid = str(sel_item['id'])
            if self.exists(str_uid):
                self.see(str_uid)
                self.selection_set(str_uid)
        
    def expand_all(self):
        for c in self.get_children(""):
            self.item(c, open=True)
            
    def get_selected(self):
        sel = self.selection()
        if sel:
            return self.data[sel[0]]
        return None
    
    #for dialogs
    def accept(self):
        return self.get_selected() 

if __name__ == '__main__':
    
    root = ttk.Window(themename="darkly")

    style = Style()
    style.configure("Treeview", font=(None, 10))
    user_values = [
            {'id': 2, 'login': 'Serhii' , 'name': 'Сергій'},
            {'id': 3, 'login': 'Svitlana' , 'name': 'Світлана'},
            {'id': 5, 'login': 'Vadim' , 'name': 'Вадим'},
        ]

    with open ('config.json', "r") as f:
        cfg = json.loads(f.read())

    app = App()
    app.set_params(cfg, Data(cfg))
    res = app.repository.signing(user_values[1]['login'], '123')
    print('login', res)
    if res['error']:
        dlg.Messagebox.show_error(title='Помилка', message=res)
    res = app.repository.get_models()
    if res['error']:
        dlg.Messagebox.show_error(title='Помилка', message=res)
    app.set_models(res['value'])
    user = Item('user')
    err = user.get(user_values[1]['id'])
    if err:
        dlg.Messagebox.show_error(title='Помилка', message=res)
    app.user = user.value
    
    tree = Tree(root, "matherial_group")
    tree.pack(expand=True, fill='both', side="left")
    table = Table(root, 'matherial', ['id', 'name', 'price', 'cost'], (40, 250, 40, 40))
    table.pack(expand=True, fill='both', side="right")
    table.make_sortable()
    # ttk.Button(root, text="Submit", bootstyle=(SUCCESS)).pack(side='left', padx=5, pady=10)
    # ttk.Button(root, text="Submit", bootstyle=(PRIMARY, OUTLINE)).pack(side='left', padx=5, pady=10)
    # columns = ('last_name', 'email')
    # cb = ttk.Treeview(root, columns=columns, bootstyle=INFO)
    # cb.pack(expand=True, fill='both')
    # cb.heading('#0', text='First Name')
    # cb.heading('last_name', text='Last Name')
    # cb.heading('email', text='Email')
    # contacts = []
    # for n in range(0, 100):
    #     contacts.append((f'first {n}', f'last {n}', f'email{n}@example.com', n))

    # # add data to the treeview
    # for i, contact in enumerate(contacts):
    #     prt = str(i//10*10)
    #     if i%10 == 0:
    #         prt = ''
    #     cb.insert(prt, tk.END, iid=str(i), text=contact[0], values=contact[1:])
        
    # def item_selected(event):
    #     for selected_item in cb.selection():
    #         item = cb.item(selected_item)
    #         record = item['values']
    #         # show a message
    #         print(item)
    #         # dlg.Messagebox.show_info(title='Information', message=','.join(record))


    # cb.bind('<<TreeviewSelect>>', item_selected)


    root.mainloop()