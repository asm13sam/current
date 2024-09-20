import json
import sqlite3
from barcode import EAN13

# with open ('changes.json', "r") as f:
#     changes = json.loads(f.read())

# with open ('models_bk.json', "r") as f:
#     model_bk = json.loads(f.read())

# with open ('models.json', "r") as f:
#     model = json.loads(f.read())

# tables_bk = list(model_bk['models'].keys())
# tables = list(model['models'].keys())

class SqlCreator:
    def __init__(self, changes, model, model_bk) -> None:
        self.changes = changes
        self.model = model
        self.model_bk = model_bk
        self.tables_bk = list(model_bk['models'].keys())
        self.tables = list(model['models'].keys())

    def table_is_changed(self, table_name):
        if len(self.model_bk[table_name]) != len(self.model[table_name]):
            return True
        for k in self.model[table_name]:
            if k not in self.model_bk[table_name]:
                return True
        return False

    def create_table(self, table_name: str):
        print('creating', table_name)
        s = f'CREATE TABLE IF NOT EXISTS {table_name}\n('
        for k, v in self.model[table_name].items():
            d = v['def']
            if k == 'id':
                s += '\n\tid INTEGER PRIMARY KEY AUTOINCREMENT,'
            elif type(d) == int:
                s += f'\n\t{k} INT NOT NULL,'
            elif type(d) == float:
                s += f'\n\t{k} REAL NOT NULL,'
            elif type(d) == bool:
                s += f'\n\t{k} BOOL NOT NULL,'
            else:
                s += f'\n\t{k} TEXT NOT NULL,'
        s = s[:-1] + '\n);\n'
        return s

    def make_sql(self):
        s = ''
        for table_name in self.tables_bk:
            if table_name not in self.tables:
                s += f'DROP TABLE IF EXISTS {table_name};\n'

        for table_name in self.tables:
            if (table_name not in self.tables_bk
                    or table_name in self.changes['tables_for_clearing']
                    or self.table_is_changed(table_name)):
                s += f'DROP TABLE IF EXISTS {table_name};\n'
                s += self.create_table(table_name)
        return s

    def reload_table(self, table_name: str, cur_from: sqlite3.Cursor, cur_to: sqlite3.Cursor):
        query_str = f"SELECT * FROM {table_name}"
        try:
            cur_from.execute(query_str)
        except:
            return
        items = cur_from.fetchall()

        templ = '?, ' * len(self.model[table_name])
        templ = templ[:-2]
        query_str = f"INSERT INTO {table_name} VALUES ({templ})"

        for item in items:
            m = dict(item)
            l = []
            for k in self.model[table_name]:
                if k in m:
                    l.append(m[k])
                else:
                    l.append(self.model[table_name][k]['def'])
            cur_to.execute(query_str, l)

    def clear_fields_in_table(self, table_name, fields, cur_to):
        fields_str = 'SET '
        for field in fields:
            fields_str += f'{field} = {self.model[table_name][field]["def"]}, '
        query_str = f'UPDATE {table_name} {fields_str[:-2]};'
        cur_to.execute(query_str)

    def recreate_tables(self, cur_to: sqlite3.Cursor, cur_from: sqlite3.Cursor, cur_update: sqlite3.Cursor):
        s = self.make_sql()
        cur_to.executescript(s)
        for table_name in self.tables:
            if (table_name in self.tables_bk
                    and table_name not in self.changes['tables_for_clearing']
                    and self.table_is_changed(table_name)):
                print('reloading', table_name)
                self.reload_table(table_name, cur_from, cur_to)

        for table_name in self.tables:
            if table_name in self.changes['update']:
                print('updating', table_name)
                s = ''
                s += f'DROP TABLE IF EXISTS {table_name};\n'
                s += self.create_table(table_name)
                cur_to.executescript(s)
                self.reload_table(table_name, cur_update, cur_to)

        # self.add_operations_barcodes(cur_to, cur_update)
        # self.add_products_barcodes(cur_to, cur_update)

        for table_name in self.tables:
            if table_name in self.changes['fields_for_clearing']:
                fields = self.changes['fields_for_clearing'][table_name]
                print('clearing', table_name, fields)
                self.clear_fields_in_table(table_name, fields, cur_to)

    def add_operations_barcodes(self, cur_to: sqlite3.Cursor, cur_from: sqlite3.Cursor):
        cur_from.execute("SELECT * FROM operation")
        operations = cur_from.fetchall()
        for operation in operations:
            o = dict(operation)
            my_code = EAN13('222223%06d' % o['id'])
            my_code.build()
            barcode = my_code.get_fullcode()
            cur_to.execute(
                "UPDATE operation SET barcode=? WHERE id=?", (barcode, o['id'])
            )

    def add_products_barcodes(self, cur_to: sqlite3.Cursor, cur_from: sqlite3.Cursor):
        cur_from.execute("SELECT * FROM product")
        products = cur_from.fetchall()
        for product in products:
            p = dict(product)
            my_code = EAN13('222224%06d' % p['id'])
            my_code.build()
            barcode = my_code.get_fullcode()
            cur_to.execute(
                "UPDATE product SET barcode=? WHERE id=?", (barcode, p['id'])
            )


# ~ def create_sql():
    # ~ s = ''
    # ~ for table in tables:
        # ~ # if table not in changes:
            # ~ # continue
        # ~ # print(table)
        # ~ # print(sts[table])
        # ~ s += f'DROP TABLE IF EXISTS {table};\nCREATE TABLE IF NOT EXISTS {table}\n('
        # ~ for k, v in model[table].items():
            # ~ #print(k, v)
            # ~ d = v['def']
            # ~ if k == 'id':
                # ~ s += '\n\tid INTEGER PRIMARY KEY AUTOINCREMENT,'
            # ~ elif type(d) == int:
                # ~ s += f'\n\t{k} INT NOT NULL,'
            # ~ elif type(d) == float:
                # ~ s += f'\n\t{k} REAL NOT NULL,'
            # ~ elif type(d) == bool:
                # ~ s += f'\n\t{k} BOOL NOT NULL,'
            # ~ else:
                # ~ s += f'\n\t{k} TEXT NOT NULL,'
        # ~ s = s[:-1] + '\n);\n'


    # ~ # with open ('sql/bases.sql', 'w') as f:
    # ~ #     f.write(s)

    # ~ with open ('sql/make_bases.py', 'w') as f:
        # ~ b = '''
# ~ import sqlite3

# ~ def recreate_table(cur):
    # ~ s = """
        # ~ '''
        # ~ e = '''
    # ~ """
    # ~ cur.executescript(s)


# ~ if __name__ == '__main__':
    # ~ import shutil

    # ~ shutil.copyfile('base.db', 'base.bk')

    # ~ con = sqlite3.connect('base.db')
    # ~ cur = con.cursor()

    # ~ recreate_table(cur)
        # ~ '''
        # ~ f.write(b + s + e)

# ~ if __name__ == '__main__':
    # ~ create_sql()
