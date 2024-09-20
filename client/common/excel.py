import xlsxwriter

from common.funcs import thousands_to_words, dataiso_to_words

class ExcelDoc:
    def __init__(self, owner, invoice, items_to_invoice, contragent) -> None:
        self.owner = owner
        self.invoice = invoice
        self.items_to_invoice = items_to_invoice
        self.contragent = contragent
        
    def create_sheet(self):
        self.sheet = self.book.add_worksheet()
        self.sheet.set_paper(11)
        self.sheet.set_portrait()
        self.sheet.set_margins(0.2, 0.2, 0.2, 0.2)
        self.sheet.fit_to_pages(1, 1)
        self.sheet.set_column("A:A", 2)
        self.sheet.set_row(1, 22)

    def create_docs(self, name, with_sign=False, with_date=False):
        self.book = xlsxwriter.Workbook(name)
        self.make_formats()
        self.create_sheet()
        if with_date:
            date = dataiso_to_words(self.invoice["created_at"])
            title = f'Рахунок на оплату № {self.invoice["id"]} від {date}'
        else:
            date = dataiso_to_words(self.invoice["created_at"], only_year=True)
            title = f'Рахунок на оплату № {self.invoice["id"]} від __________________ {date}'    
        self.create_top(title)
        row = self.create_table()
        self.create_invoice_bottom(row, with_sign)

        self.create_sheet()
        date = dataiso_to_words(self.invoice["created_at"], only_year=True)
        title = f'Накладна № {self.invoice["id"]} від __________________ {date}'
        self.create_top(title)
        row = self.create_table()
        self.create_whs_out_bottom(row, with_sign)
        
        self.book.close()

    def make_formats(self):
        self.title_format = self.book.add_format(
            {
                "bold": 1,
                "bottom": 2,
                "font_size": 14,
                "font_name": 'Arial',
                "align": "left",
                "valign": "vcenter",
            }
        )

        self.main_format = self.book.add_format(
            {
                "font_size": 9,
                "font_name": 'Arial',
                "align": "left",
            }
        )

        self.bold_format = self.book.add_format(
            {
                "bold": 1,
                "font_size": 9,
                "font_name": 'Arial',
                "align": "left",
            }
        )

        self.tab_tl_format = self.book.add_format(
            {
                "bold": 1,
                "bottom": 1,
                "top": 2,
                "left": 2,
                "right": 1,
                "font_size": 9,
                "font_name": 'Arial',
                "align": "center",
                "valign": "vcenter",
                "fg_color": "#fffbcc",
            }
        )

        self.tab_tc_format = self.book.add_format(
            {
                "bold": 1,
                "bottom": 1,
                "top": 2,
                "left": 1,
                "right": 1,
                "font_size": 9,
                "font_name": 'Arial',
                "align": "center",
                "valign": "vcenter",
                "fg_color": "#fffbcc",
            }
        )

        self.tab_tr_format = self.book.add_format(
            {
                "bold": 1,
                "bottom": 1,
                "top": 2,
                "left": 1,
                "right": 2,
                "font_size": 9,
                "font_name": 'Arial',
                "align": "center",
                "valign": "vcenter",
                "fg_color": "#fffbcc",
            }
        )

        self.tab_cl_format = self.book.add_format(
            {
                "bottom": 1,
                "top": 1,
                "left": 2,
                "right": 1,
                "font_size": 9,
                "font_name": 'Arial',
                "align": "center",
                "valign": "vcenter",
            }
        )

        self.tab_cc_format = self.book.add_format(
            {
                "bottom": 1,
                "top": 1,
                "left": 1,
                "right": 1,
                "font_size": 9,
                "font_name": 'Arial',
                "align": "left",
                "valign": "vcenter",
            }
        )

        self.tab_num_format = self.book.add_format(
            {
                "bottom": 1,
                "top": 1,
                "left": 1,
                "right": 1,
                "font_size": 9,
                "font_name": 'Arial',
                "align": "right",
                "valign": "vcenter",
                'num_format': '#,##0'
            }
        )

        self.tab_price_format = self.book.add_format(
            {
                "bottom": 1,
                "top": 1,
                "left": 1,
                "right": 1,
                "font_size": 9,
                "font_name": 'Arial',
                "align": "right",
                "valign": "vcenter",
                'num_format': '0.000'
            }
        )

        self.tab_cost_format = self.book.add_format(
            {
                "bottom": 1,
                "top": 1,
                "left": 1,
                "right": 2,
                "font_size": 9,
                "font_name": 'Arial',
                "align": "right",
                "valign": "vcenter",
                'num_format': '0.00'
            }
        )

        self.tab_cr_format = self.book.add_format(
            {
                "bottom": 1,
                "top": 1,
                "left": 1,
                "right": 2,
                "font_size": 9,
                "font_name": 'Arial',
                "align": "center",
                "valign": "vcenter",
            }
        )

        self.tab_bl_format = self.book.add_format(
            {
                "bottom": 2,
                "top": 1,
                "left": 2,
                "right": 1,
                "font_size": 9,
                "font_name": 'Arial',
                "align": "center",
                "valign": "vcenter",
            }
        )

        self.tab_bc_format = self.book.add_format(
            {
                "bottom": 2,
                "top": 1,
                "left": 1,
                "right": 1,
                "font_size": 9,
                "font_name": 'Arial',
                "align": "left",
                "valign": "vcenter",
            }
        )

        self.tab_br_format = self.book.add_format(
            {
                "bottom": 2,
                "top": 1,
                "left": 1,
                "right": 2,
                "font_size": 9,
                "font_name": 'Arial',
                "align": "center",
                "valign": "vcenter",
            }
        )

        self.total_format = self.book.add_format(
            {
                "bold": 1,
                "font_size": 10,
                "font_name": 'Arial',
                "align": "center",
                "valign": "vcenter",
                'num_format': '#,##0.00',
            }
        )
        self.bc_format = self.book.add_format({"bottom": 2})

    def create_top(self, title):
        self.sheet.merge_range("B2:L2", title, self.title_format)
        self.sheet.merge_range("B4:C4", 'Відпущено:', self.main_format)
        self.sheet.merge_range("D4:K4", self.owner['full_name'], self.bold_format)
        txt = f'П/р {self.owner["iban"]}, Банк {self.owner["bank"]}\n'
        txt += f'МФО {self.owner["mfo"]}\n'
        txt += f'{self.owner["address"]}\n'
        txt += f'код за ЄДРПОУ {self.owner["edrpou"]}\n'
        txt += f'{self.owner["fop"]}'
        self.sheet.set_row(4, 80)
        self.sheet.merge_range("D5:L5", '')
        self.sheet.insert_textbox(4, 3, txt, {'width': 512, 'height': 100, 'line': {'none': True}})
        
        self.sheet.merge_range("B6:C6", 'Покупець:', self.main_format)
        self.sheet.merge_range("D6:K6", self.contragent['name'], self.bold_format)
        self.sheet.merge_range("B8:C8", 'Договір:', self.main_format)
        self.sheet.merge_range("D8:K8", '', self.bold_format)

    def create_table(self):
        self.sheet.set_row(9, 30)
        self.sheet.write('B10', '№', self.tab_tl_format)
        self.sheet.merge_range("C10:G10", 'Товари (роботи, послуги)', self.tab_tc_format)
        self.sheet.write('H10', 'Кіл-сть', self.tab_tc_format)
        self.sheet.write('I10', 'Од.', self.tab_tc_format)
        self.sheet.write('J10', 'Ціна', self.tab_tc_format)
        self.sheet.merge_range("K10:L10", 'Сума', self.tab_tr_format)

        row = 10
        total = 0
        for i, item in enumerate(self.items_to_invoice):
            self.sheet.write(row, 1, str(i+1), self.tab_cl_format)
            self.sheet.merge_range(row, 2, row, 6, item['name'], self.tab_cc_format)
            self.sheet.write(row, 7, item['number'], self.tab_num_format)
            self.sheet.write(row, 8, item['measure'], self.tab_cc_format)
            self.sheet.write(row, 9, item['price'], self.tab_price_format)
            self.sheet.merge_range(row, 10, row, 11, item['cost'], self.tab_cost_format)
            total += item['cost']
            row += 1

        self.sheet.write(row, 1, '', self.tab_bl_format)
        self.sheet.merge_range(row, 2, row, 6, '', self.tab_bc_format)
        self.sheet.write(row, 7, '', self.tab_bc_format)
        self.sheet.write(row, 8, '', self.tab_bc_format)
        self.sheet.write(row, 9, '', self.tab_bc_format)
        self.sheet.merge_range(row, 10, row, 11, '', self.tab_br_format)

        row += 2
        self.sheet.write(row, 9, 'Всього:', self.bold_format)
        self.sheet.merge_range(row, 10, row, 11, total, self.total_format)
        row += 2
        txt = f'Всього найменувань {i+1}, на суму {total} грн.'
        self.sheet.merge_range(row, 1, row, 10, txt, self.main_format)
        row += 1
        total_words = thousands_to_words(int(total))
        self.sheet.merge_range(row, 1, row, 10, total_words.capitalize(), self.bold_format)
        row += 1
        self.sheet.write(row, 1, 'Без ПДВ', self.main_format)
        row += 1
        self.sheet.merge_range(row, 1, row, 11, "", self.bc_format)
        return row

    def create_invoice_bottom(self, row, with_sign):
        row += 3
        bc_format = self.book.add_format({"bottom": 1})
        bold_format = self.book.add_format({"bold": 1, "font_size": 9, "font_name": 'Arial', "align": "left"})
        self.sheet.merge_range(row, 6, row, 7, 'Виписав(ла):', bold_format)
        self.sheet.merge_range(row, 8, row, 9, "", bc_format)
        self.sheet.write(row, 10, self.owner['name'], bold_format)
        if with_sign:
            sign_image_filename = 'images/signs/' + self.owner['sign']
            self.sheet.insert_image(row-3, 8, sign_image_filename)
        row += 2
        self.sheet.write(row, 9, 'Б.П.', bold_format)

    def create_whs_out_bottom(self, row, with_sign):
        row += 3
        bc_format = self.book.add_format({"bottom": 1})
        bold_format = self.book.add_format({"bold": 1, "font_size": 9, "font_name": 'Arial', "align": "left"})
        bold_right_format = self.book.add_format({"bold": 1, "font_size": 9, "font_name": 'Arial', "align": "right"})
        self.sheet.merge_range(row, 1, row, 2, 'Від постачальника:', bold_right_format)
        self.sheet.merge_range(row, 3, row, 4, "", bc_format)
        self.sheet.write(row, 5, self.owner['name'], bold_format)
        self.sheet.merge_range(row, 7, row, 8, 'Отримав(ла):', bold_right_format)
        self.sheet.merge_range(row, 9, row, 11, "", bc_format)
        if with_sign:
            sign_image_filename = 'images/signs/' + self.owner['sign']
            self.sheet.insert_image(row-3, 3, sign_image_filename)
        row += 2
        self.sheet.write(row, 4, 'Б.П.', bold_format)
        self.sheet.merge_range(row, 7, row, 8, 'За довіреністю №', bold_right_format)
        self.sheet.write(row, 9, "", bc_format)
        self.sheet.write(row, 10, "від", bold_right_format)
        self.sheet.write(row, 11, '', bc_format)
        


