import json


def to_go(k):
    return ''.join(s.title() for s in k.split("_"))

def create_go_models_head():

    g = f'''
package main

import (
    "database/sql"
    "errors"
    "fmt"
    "time"


    _ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func DBconnect(dbFile string) error {{
    var err error
    db, err = sql.Open("sqlite3", dbFile)
    return err
}}
func DbClose() error {{
    return db.Close()
}}
\n\n
'''
    return g

def create_go_main_head():
    g = '''
package main

import (
    "encoding/json"
    "log"
    "net/http"
    "os"
    "strings"

    "asm13sam/tg"

    "github.com/gorilla/mux"
)

func makeRouter() *mux.Router {
    r := mux.NewRouter()
    '''
    return g

def create_go_main_footer():
    g = '''
    r.HandleFunc("/upload/{id:[0-9]+}",
        WrapAuth(UploadFile, DOC_CREATE)).Methods("POST")

    r.HandleFunc("/product_to_ordering_default",
        WrapAuth(CreateProductToOrderingDefault, DOC_READ)).Methods("POST")

    r.HandleFunc("/product_deep/{id:[0-9]+}",
        WrapAuth(GetProductDeep, DOC_READ)).Methods("GET")

    r.HandleFunc("/product_complex/{id:[0-9]+}",
        WrapAuth(GetProductComplex, DOC_READ)).Methods("GET")

    r.HandleFunc("/project_dirs/{id:[0-9]+}",
        WrapAuth(CreateProjectDirs, DOC_CREATE)).Methods("GET")

    r.PathPrefix("/static/").Handler(http.StripPrefix("/static/",
        http.FileServer(http.Dir("./static/"))))

    r.HandleFunc("/login", WrapAuth(Login, LOGIN)).Methods("POST")
    r.HandleFunc("/logout", WrapAuth(Logout, LOGOUT)).Methods("GET")
    r.HandleFunc("/copy_project/{id:[0-9]+}", WrapAuth(CopyProject, ADMIN)).Methods("GET")
    r.HandleFunc("/copy_base/{fs}", WrapAuth(CopyBase, ADMIN)).Methods("GET")
    r.HandleFunc("/delete_base/{fs}", WrapAuth(DeleteBackupBase, ADMIN)).Methods("GET")
    r.HandleFunc("/get_bases", WrapAuth(GetBackupBases, ADMIN)).Methods("GET")
    r.HandleFunc("/restore_base/{fs}", WrapAuth(RestoreBaseFromBackup, ADMIN)).Methods("GET")
    r.HandleFunc("/ws", WrapAuth(UpgradeWS, WS_CONNECT)).Methods("GET")

    return r
}

type Config struct {
    Port          string   `json:"port"`
    BckpPath      string   `json:"bckp_path"`
    MaketsPath    string   `json:"makets_path"`
    OldMaketsPath string   `json:"old_makets_path"`
    NewMaketsPath string   `json:"new_makets_path"`
    MaketDirs     []string `json:"maket_dirs"`
    DBFile        string   `json:"db_file"`
}

var Cfg Config

func LoadConfig() error {
    data, err := os.ReadFile("config.json")
    if err == nil {
        decoder := json.NewDecoder(strings.NewReader(string(data)))
        err = decoder.Decode(&Cfg)
    }
    return err
}

func main() {

    if err := LoadConfig(); err != nil {
        log.Fatal("Can`t load config file!")
    }

    err := DBconnect(Cfg.DBFile)
    if err != nil {
        log.Fatal(err)
    }

    log.Print("running on port ", Cfg.Port)
    log.Println("makets on path ", Cfg.MaketsPath)
    r := makeRouter()
    //log.Fatal(http.ListenAndServe("0.0.0.0:"+*portFlag, r))

    Root, err := tg.InitRoot("Target", tg.Horizontal)
    if err != nil {
        log.Fatal(err)
    }

    tg.SetTheme("clam")
    tg.AddStatusField()
    tg.AddStatus(0, "This is test start.")

    makeGui(Root, Cfg.Port, r)

    tg.MainLoop()
}

func makeGui(Root tg.Container, port string, router *mux.Router) {
    mainFrame := tg.NewSplitter(tg.Expand | tg.Horizontal)
    Root.Add(mainFrame)

    MainPanel := tg.NewBox(tg.Expand)
    RightPanel := tg.NewBox(0)
    mainFrame.Add(MainPanel, RightPanel)
    block := tg.NewButton("Start Server", 0)
    RightPanel.Add(block)
    block.IfPressed(func(s string) {
        tg.MessageBox("Starting server on port " + port)
        go http.ListenAndServe("0.0.0.0:"+port, router)
        block.Disable()
        tg.AddStatus(0, "Starting server on port "+port)
        go handleMessages()
    })
}
    '''
    return g

def create_go_handler_head():
    g = '''
package main

import (
    "encoding/json"
)

func GetProductComplex(req Req) {
    req.Respond(ProductComplexGet(req.IntParam))
}

func GetProductDeep(req Req) {
    counter := 1
    req.Respond(ProductDeepGet(req.IntParam, &counter))
}


func CreateProductToOrderingDefault(req Req) {
    p, err := DecodeProductToOrdering(req)
    if err != nil {
        req.Respond(nil, err)
        return
    }
    req.Respond(ProductToOrderingCreateDefault(p))
}

    '''
    return g

def create_go_models(model, tables):
    g = create_go_models_head()
    h = create_go_handler_head()
    m = create_go_main_head()
    for table in tables:
        g += create_go_model(table, model)
        keys = model[table].keys()

        g1, h1, m1 = create_go_get(table, keys, model)
        g += g1
        h += h1
        m += m1

        g1, h1, m1 = create_go_gets(table, keys, model)
        g += g1
        h += h1
        m += m1

        g1, h1, m1 = create_go_create(table, keys, model)
        g += g1
        h += h1
        m += m1

        g1, h1, m1 = create_go_update(table, keys, model)
        g += g1
        h += h1
        m += m1

        g1, h1, m1 = create_go_delete(table, keys, model)
        g += g1
        h += h1
        m += m1

        g1, h1, m1 = create_go_filter_int(table, keys, model)
        g += g1
        h += h1
        m += m1

        g1, h1, m1 = create_go_filter_str(table, keys, model)
        g += g1
        h += h1
        m += m1

        h += create_go_decode(table)
        g += make_field_validation(table, keys)
        
        if table in model['documents'] and table != 'ordering': 
            g1, h1, m1 = create_go_realized(table, keys, model)
            g += g1
            h += h1
            m += m1

        if table in model['doc_table_items'] or table == 'item_to_invoice': 
            g1 = create_go_realized_item(table, keys, model)
            g += g1
            

        if 'find' in model['models'][table]:
            finds = model['models'][table]['find']
            for find in finds:
                g1, h1, m1 = create_go_find(table, keys, find, model)
                g += g1
                h += h1
                m += m1

        if 'between' in model['models'][table]:
            betweens = model['models'][table]['between']
            for between in betweens:
                g1, h1, m1 = create_go_between(table, keys, between, model)
                g += g1
                h += h1
                m += m1
        if 'sum' in model['models'][table]:
            sums = model['models'][table]['sum']
            for sum_field in sums:
                g1, h1, m1 = create_go_sum_before(table, keys, sum_field, model)
                g += g1
                h += h1
                m += m1
            sum_field = sums[0]
            g1, h1, m1 = create_go_sum_filter(table, keys, sum_field, model)
            g += g1
            h += h1
            m += m1

    g1, h1, m1 = create_go_models_w(model, tables)
    g += g1
    h += h1
    m += m1


    with open ('../golang/models.go', 'w') as f:
        f.write(g)

    with open ('../golang/handlers.go', 'w') as f:
        f.write(h)

    m += create_go_main_footer()
    with open ('../golang/main.go', 'w') as f:
        f.write(m)


def create_go_models_w(model, tables):
    g = ''
    h = ''
    m = ''
    for table in tables:
        keys = model[table].keys()
        g += create_go_model_w(table, model)

        g1, h1, m1 = create_go_get_w(table, keys, model)
        g += g1
        h += h1
        m += m1

        g1, h1, m1 = create_go_gets_w(table, keys, model)
        g += g1
        h += h1
        m += m1

        g1, h1, m1 = create_go_filter_w_int(table, keys, model)
        g += g1
        h += h1
        m += m1

        g1, h1, m1 = create_go_filter_w_str(table, keys, model)
        g += g1
        h += h1
        m += m1

        if 'find' in model['models'][table]:
            finds = model['models'][table]['find']
            for find in finds:
                g1, h1, m1 = create_go_find_w(table, keys, find, model)
                g += g1
                h += h1
                m += m1

        if 'between' in model['models'][table]:
            betweens = model['models'][table]['between']
            for between in betweens:
                g1, h1, m1 = create_go_between_w(table, keys, between, model)
                g += g1
                h += h1
                m += m1
        if 'between_up' in model['models'][table]:
            betweens = model['models'][table]['between_up']
            for between in betweens:
                g1, h1, m1 = create_go_between_up_w(table, keys, between, model)
                g += g1
                h += h1
                m += m1

    return g, h, m



def list_of_pointers(keys, gv):
    g = ''
    for k in keys:
        g += f'\n\t\t&{gv}.{to_go(k)},'
    return g[2:]

def list_of_vars(keys, gv):
    g = ''
    for k in keys:
        g += f'\n\t\t{gv}.{to_go(k)},'
    return g[2:]

def list_of_vars_no_id(keys, gv):
    g = ''
    for k in keys:
        if k != 'id':
            g += f'\n\t\t{gv}.{to_go(k)},'
    return g[2:]

def create_go_model(table, model):
    g = ''
    g += f'\n\ntype {to_go(table)} struct ' + '{'
    for k, v in model[table].items():
        d = v['def']
        if type(d) == int:
            g += f'\n\t{to_go(k)} int `json:"{k}"`'
        elif type(d) == float:
            g += f'\n\t{to_go(k)} float64 `json:"{k}"`'
        elif type(d) == bool:
            g += f'\n\t{to_go(k)} bool `json:"{k}"`'
        else:
            g += f'\n\t{to_go(k)} string `json:"{k}"`'
    g += '\n}'
    return g


def create_go_get(table, keys, model):
    right = model['models'][table]['rights'] + '_READ'
    gtype = to_go(table)
    gv = table[0]
    m = f'''
        r.HandleFunc("/{table}/{{id:[0-9]+}}",
        WrapAuth(Get{gtype}, {right})).Methods("GET")
    '''

    h = f'''
        func Get{gtype}(req Req) {{
            req.Respond({gtype}Get(req.IntParam, nil))
        }}
    '''

    g = f'''
        func {gtype}Get(id int, tx *sql.Tx) ({gtype}, error) {{
            var {gv} {gtype}
            var row *sql.Row
            if tx != nil {{
                row = tx.QueryRow("SELECT * FROM {table} WHERE id=?", id)
            }} else {{
                row = db.QueryRow("SELECT * FROM {table} WHERE id=?", id)
            }}

            err := row.Scan(
                {list_of_pointers(keys, gv)}
            )
            return {gv}, err
        }}
    '''
    return g, h, m


def create_go_gets(table, keys, model):
    right = model['models'][table]['rights'] + '_READ'
    gtype = to_go(table)
    gv = table[0]

    m = f'''
        r.HandleFunc("/{table}_get_all",
            WrapAuth(Get{gtype}All, {right})).Methods("GET")
    '''

    h = f'''
        func Get{gtype}All(req Req) {{
            req.Respond({gtype}GetAll(req.WithDeleted, req.DeletedOnly, nil))
        }}
    '''

    g = f'''
        func {gtype}GetAll(withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]{gtype}, error) {{
            var rows *sql.Rows
            var err error
            query := "SELECT * FROM {table}"
            if deletedOnly {{
                query += " WHERE is_active = 0"
            }} else if !withDeleted {{
                query += " WHERE is_active = 1"
            }}

            if tx != nil {{
                rows, err = tx.Query(query)
            }} else {{
                rows, err = db.Query(query)
            }}
            if err != nil {{
                return nil, err
            }}
            defer rows.Close()
            res := []{gtype}{{}}
            for rows.Next() {{
                var {gv} {gtype}
                if err := rows.Scan(
            {list_of_pointers(keys, gv)}
                ); err != nil {{
                    return nil, err
                }}
                res = append(res, {gv})
            }}
            return res, nil
        }}
        '''
    return g, h, m

def create_go_decode(type):
    gtype = to_go(type)
    gv = type[0]

    g = f'''
        func Decode{gtype}(req Req) ({gtype}, error) {{
            decoder := json.NewDecoder(req.R.Body)
            defer req.R.Body.Close()
            var {gv} {gtype}
            err := decoder.Decode(&{gv})
            return {gv}, err
        }}
        '''
    return g

def create_go_create_registers(register, gv):
    reg_table, reg_field = register["reg_field"].split(".")
    g_reg_field = to_go(reg_field)
    g_reg_table = to_go(reg_table)

    apply_func = ''
    for vf in register["val_field"]:
        apply_func += f"{reg_table}.{g_reg_field} {register['func']}= {gv}.{to_go(vf)}\n"

    r_get = f'''
    {reg_table}, err := {g_reg_table}Get({gv}.{g_reg_table}Id, tx)
    if err == nil {{
        {apply_func}
        _, err = {g_reg_table}Update({reg_table}, tx)
        if err != nil {{
            return {gv}, err
        }}
    }}
    '''
    return r_get

def create_go_create_complex_registers(register, gv, table):
    g_table = to_go(table)
    reg = f'''
        err = Create{g_table}ToNumber(&{gv}, tx)
        if err != nil {{
            return {gv}, err
        }}
    '''
    return reg


def create_go_create(table, keys, model):
    right = model['models'][table]['rights'] + '_CREATE'
    gtype = to_go(table)
    gv = table[0]

    m = f'''
        r.HandleFunc("/{table}",
        WrapAuth(Create{gtype}, {right})).Methods("POST")
    '''

    h = f'''
        func Create{gtype}(req Req){{
            {gv}, err := Decode{gtype}(req)
            if err != nil {{
                req.Respond(nil, err)
                return
            }}
            req.Respond({gtype}Create({gv}, nil))
        }}
    '''
    reg_get = ''
    if 'register' in model['models'][table]:
        registers = model['models'][table]['register']
        for register in registers:
            reg_get += create_go_create_registers(register, gv)

    create_doc = ''
    doc_update_name = ''
    doc_update_name_tx = ''
    if table in model['documents']:
        create_doc = f'''
            doc := Document{{Id:0, DocType: "{table}", IsActive:true}}
            doc, err = DocumentCreate(doc, tx)
            if err != nil {{
                return {gv}, err
            }}
            {gv}.DocumentUid = doc.Id
        '''
        if table != 'ordering':
            doc_update_name = f'''
                {gv}, err = {gtype}Update({gv}, tx)
                if err != nil {{
                    return {gv}, err
                }}
            '''
            doc_update_name_tx = f'''
                {gv}.Name = fmt.Sprintf("%s-%d", {gv}.Name, {gv}.Id)
            '''
    created_at = ''
    if 'created_at' in keys:
        created_at = f'''
        t := time.Now()
        {gv}.CreatedAt = t.Format("2006-01-02T15:04:05")
        '''

    g = f'''
    func {gtype}Create({gv} {gtype}, tx *sql.Tx) ({gtype}, error) {{
        var err error
        needCommit := false

        if tx == nil {{
            tx, err = db.Begin()
            if err != nil {{
                return {gv}, err
            }}
            needCommit = true
            defer tx.Rollback()
        }}
        {reg_get}{create_doc}{created_at}
        sql := `INSERT INTO {table}
            ({', '.join(list(keys)[1:])})
            VALUES({('?, '*(len(keys)-1))[:-2]});`
        res, err := tx.Exec(
            sql,
        {list_of_vars_no_id(keys, gv)}
        )
        if err != nil {{
            return {gv}, err
        }}
        last_id, err := res.LastInsertId()
        if err != nil {{
            return {gv}, err
        }}
        {gv}.Id = int(last_id){doc_update_name_tx}
        {doc_update_name}
        if needCommit {{
            err = tx.Commit()
            if err != nil {{
                return {gv}, err
            }}
        }}
        return {gv}, nil
    }}
    '''
    return g, h, m

def create_go_realized_relateds(related, gv):
    table = related['table']
    gtable = to_go(table)
    val = related['filter_value']
    gval = to_go(val)
    r = f'''
    {table}s, err := {gtable}GetByFilterInt("{related['filter']}", {gv}.{gval}, false, false, tx)
    if err != nil {{
        return {gv}, err
    }}
    for _, {table} := range {table}s {{
        _, err = {gtable}Realized({table}.Id, tx)
        if err != nil {{
            return {gv}, err
        }}
    }}
    '''
    return r

def create_go_realized(table, keys, model):
    right = model['models'][table]['rights'] + '_CREATE'
    gtype = to_go(table)
    gv = table[0]

    m = f'''
        r.HandleFunc("/realized/{table}/{{id:[0-9]+}}",
        WrapAuth(Realized{gtype}, {right})).Methods("GET")
    '''

    h = f'''
        func Realized{gtype}(req Req){{
            req.Respond({gtype}Realized(req.IntParam, nil))
        }}
    '''
    complex_reg = ''
    if 'complex_register' in model['models'][table]:
        complex_registers = model['models'][table]['complex_register']
        for register in complex_registers:
            complex_reg += create_go_create_complex_registers(register, gv, table)

    reg_get = ''
    if 'rz_register' in model['models'][table]:
        registers = model['models'][table]['rz_register']
        for register in registers:
            reg_get += create_go_create_registers(register, gv)
    rel_realized = ''
    if 'related' in model['models'][table]:
        relateds = model['models'][table]['related']
        for related in relateds:
            r = create_go_realized_relateds(related, gv)
            rel_realized += r
    
    g = f'''
    func {gtype}Realized(id int, tx *sql.Tx) ({gtype}, error) {{
        var err error
        needCommit := false
        var {gv} {gtype}
        if tx == nil {{
            tx, err = db.Begin()
            if err != nil {{
                return {gv}, err
            }}
            needCommit = true
            defer tx.Rollback()
        }}
        {gv}, err = {gtype}Get(id, tx)
            if err != nil {{
                return {gv}, err
            }}
        {complex_reg}{reg_get}{rel_realized}
        sql := `UPDATE {table} SET is_realized=1 WHERE id=?;`
        _, err = tx.Exec(sql, {gv}.Id)
        if err != nil {{
            return {gv}, err
        }}
        if needCommit {{
            err = tx.Commit()
            if err != nil {{
                return {gv}, err
            }}
        }}
        return {gv}, nil
    }}
    '''
    return g, h, m

def create_go_realized_item(table, keys, model):
    right = model['models'][table]['rights'] + '_CREATE'
    gtype = to_go(table)
    gv = table[0]
    
    complex_reg = ''
    if 'complex_register' in model['models'][table]:
        complex_registers = model['models'][table]['complex_register']
        for register in complex_registers:
            complex_reg += create_go_create_complex_registers(register, gv, table)

    reg_get = ''
    if 'rz_register' in model['models'][table]:
        registers = model['models'][table]['rz_register']
        for register in registers:
            reg_get += create_go_create_registers(register, gv)
    rel_realized = ''
    if 'related' in model['models'][table]:
        relateds = model['models'][table]['related']
        for related in relateds:
            r = create_go_realized_relateds(related, gv)
            rel_realized += r
    
    g = f'''
    func {gtype}Realized(id int, tx *sql.Tx) ({gtype}, error) {{
        var err error
        needCommit := false
        var {gv} {gtype}
        if tx == nil {{
            tx, err = db.Begin()
            if err != nil {{
                return {gv}, err
            }}
            needCommit = true
            defer tx.Rollback()
        }}
        {gv}, err = {gtype}Get(id, tx)
            if err != nil {{
                return {gv}, err
            }}
        {complex_reg}{reg_get}{rel_realized}
        
        if needCommit {{
            err = tx.Commit()
            if err != nil {{
                return {gv}, err
            }}
        }}
        return {gv}, nil
    }}
    '''
    return g


def create_go_update_complex_registers(register, gv, table):
    g_table = to_go(table)
    reg = f'''
        err = Update{g_table}ToNumber(&{gv}, {table}.Number, tx)
        if err != nil {{
            return {gv}, err
        }}
    '''
    return reg


def create_go_update_registers_get_previous(gv, table):
    gtable = to_go(table)

    r = f'''
    {table}, err := {gtable}Get({gv}.Id, tx)
    if err != nil {{
        return {gv}, err
    }}
    '''
    return r


def create_go_update_registers(register, gv, table):
    reg_table, reg_field = register["reg_field"].split(".")
    g_reg_field = to_go(reg_field)
    g_reg_table = to_go(reg_table)
    undo_func = '-' if register['func'] == '+' else '+'

    apply_func = ''
    apply_undo_func = ''
    for vf in register["val_field"]:
        apply_undo_func += f"{reg_table}.{g_reg_field} {undo_func}= {table}.{to_go(vf)}\n"
        apply_func += f"{reg_table}.{g_reg_field} {register['func']}= {gv}.{to_go(vf)}\n"


    r_get = f'''
    {reg_table}, err := {g_reg_table}Get({table}.{g_reg_table}Id, tx)
    if err == nil {{
        {apply_undo_func}
    }}

    if {table}.{g_reg_table}Id != {gv}.{g_reg_table}Id {{
        _, err = {g_reg_table}Update({reg_table}, tx)
        if err != nil {{
            return {gv}, err
        }}
        {reg_table}, err = {g_reg_table}Get({gv}.{g_reg_table}Id, tx)
        if err != nil {{
            return {gv}, err
        }}
    }}
    {apply_func}
    _, err = {g_reg_table}Update({reg_table}, tx)
    if err != nil {{
        return {gv}, err
    }}
    '''
    return r_get


def create_go_update(table, keys, model):
    right = model['models'][table]['rights'] + '_UPDATE'
    gtype = to_go(table)
    gv = table[0]
    m = f'''
        r.HandleFunc("/{table}/{{id:[0-9]+}}",
            WrapAuth(Update{gtype}, {right})).Methods("PUT")
    '''

    h = f'''
        func Update{gtype}(req Req) {{
            {gv}, err := Decode{gtype}(req)
            if err != nil {{
                req.Respond(nil, err)
                return
            }}
            req.Respond({gtype}Update({gv}, nil))
        }}
        '''

    complex_reg = ''
    if 'complex_register' in model['models'][table]:
        complex_registers = model['models'][table]['complex_register']
        for register in complex_registers:
            complex_reg += create_go_update_complex_registers(register, gv, table)

    reg_get = ''
    prew = ''
    if 'register' in model['models'][table]:
        prew = create_go_update_registers_get_previous(gv, table)
        reg_get += prew
        registers = model['models'][table]['register']
        for register in registers:
            reg_get += create_go_update_registers(register, gv, table)

    rz_reg_get = ''
    if 'rz_register' in model['models'][table]:
        if not prew:
            prew = create_go_update_registers_get_previous(gv, table)
            rz_reg_get += prew
        registers = model['models'][table]['rz_register']
        for register in registers:
            rz_reg_get += create_go_update_registers(register, gv, table)

    hooks_before = ''
    hooks_after = ''
    if 'hooks' in model['models'][table]:
        hooks = model['models'][table]['hooks']
        for hook in hooks:
            if hook['act'] == 'update':
                if hook['when'] == 'before':
                    hooks_before += f"{hook['func']}(&{gv})\n"
                else:
                    hooks_after += f"{hook['func']}(&{gv})\n"
    update = ''
    if 'updated_at' in keys:
        update = f'''
        t := time.Now()
        {gv}.UpdatedAt = t.Format("2006-01-02T15:04:05")
        '''
    realized = ''
    if table in model['documents'] and table != 'ordering':
        realized = f'''
            if {gv}.IsRealized{{
                {rz_reg_get}{complex_reg}
                }}
            '''
    elif table in model['doc_table_items']:
        to_table = table.split('_to_')[1]
        realized = f'''
            if {to_table}.IsRealized{{
                {rz_reg_get}{complex_reg}
                }}
                
            '''

    g = f'''
        func {gtype}Update({gv} {gtype}, tx *sql.Tx) ({gtype}, error) {{
            {hooks_before}var err error
            needCommit := false
            if tx == nil {{
                tx, err = db.Begin()
                if err != nil {{
                    return {gv}, err
                }}
                needCommit = true
                defer tx.Rollback()
            }}
            {reg_get}{update}
            {realized}
            sql := `UPDATE {table} SET
                    {'=?, '.join(list(keys)[1:])}=?
                    WHERE id=?;`

            _, err = tx.Exec(
                    sql,
                {list_of_vars(list(keys)[1:], gv)}
                    {gv}.Id,
                )
            if err != nil {{
                return {gv}, err
            }}
            if needCommit {{
                err = tx.Commit()
                if err != nil {{
                    return {gv}, err
                }}
            }}
            {hooks_after}return {gv}, nil
        }}
        '''
    return g, h, m

def create_go_delete_complex_registers(register, gv, table):
    g_table = to_go(table)
    reg = f'''
        err = Delete{g_table}ToNumber(&{gv}, tx)
        if err != nil {{
            return {gv}, err
        }}
    '''
    return reg


def create_go_delete_relateds(related, gv):
    table = related['table']
    gtable = to_go(table)
    val = related['filter_value']
    gval = to_go(val)
    r = f'''
    {table}s, err := {gtable}GetByFilterInt("{related['filter']}", {gv}.{gval}, false, false, tx)
    if err != nil {{
        return {gv}, err
    }}
    for _, {table} := range {table}s {{
        _, err = {gtable}Delete({table}.Id, tx, isUnRealize)
        if err != nil {{
            return {gv}, err
        }}
    }}
    '''
    return r

def create_go_delete_registers(register, gv):
    reg_table, reg_field = register["reg_field"].split(".")
    g_reg_field = to_go(reg_field)
    g_reg_table = to_go(reg_table)
    undo_func = '-' if register['func'] == '+' else '+'

    apply_func = ''
    for vf in register["val_field"]:
        apply_func += f"{reg_table}.{g_reg_field} {undo_func}= {gv}.{to_go(vf)}\n"


    r_get = f'''
    {reg_table}, err := {g_reg_table}Get({gv}.{g_reg_table}Id, tx)
    if err == nil {{
        {apply_func}
        _, err = {g_reg_table}Update({reg_table}, tx)
        if err != nil {{
            return {gv}, err
        }}
    }}
    '''
    return r_get


def create_go_delete(table, keys, model):
    right = model['models'][table]['rights'] + '_DELETE'
    gtype = to_go(table)
    gv = table[0]
    m = ''
    h = ''
    check_unrealize = f'''
            sql := `UPDATE {table} SET is_active=0 WHERE id=?;`
            _, err = tx.Exec(sql, {gv}.Id)
            if err != nil {{
                return {gv}, err
            }}
            '''
    if table in model['documents'] and model != 'ordering':
        m += f'''
            r.HandleFunc("/unrealize/{table}/{{id:[0-9]+}}",
                WrapAuth(UnRealize{gtype}, {right})).Methods("GET")
        '''

        h += f'''
            func UnRealize{gtype}(req Req) {{
                req.Respond({gtype}Delete(req.IntParam, nil, true))
            }}
            '''    
        check_unrealize = f'''
            sql := `UPDATE {table} SET is_active=0 WHERE id=?;`
            if isUnRealize {{
                sql = `UPDATE {table} SET is_realized=0 WHERE id=?;`
            }} 
            _, err = tx.Exec(sql, {gv}.Id)
            if err != nil {{
                return {gv}, err
            }}
            '''
    if table in model['doc_table_items'] or model == 'item_to_invoice':
        check_unrealize = f'''
            if !isUnRealize {{
                sql := `UPDATE {table} SET is_active=0 WHERE id=?;`
                _, err = tx.Exec(sql, {gv}.Id)
                if err != nil {{
                    return {gv}, err
                }}
            }} 
            '''

    m += f'''
        r.HandleFunc("/{table}/{{id:[0-9]+}}",
            WrapAuth(Delete{gtype}, {right})).Methods("DELETE")
    '''

    h += f'''
        func Delete{gtype}(req Req) {{
            req.Respond({gtype}Delete(req.IntParam, nil, false))
        }}
        '''
    complex_reg = ''
    if 'complex_register' in model['models'][table]:
        complex_registers = model['models'][table]['complex_register']
        for register in complex_registers:
            complex_reg += create_go_delete_complex_registers(register, gv, table)

    reg_get = ''
    if 'register' in model['models'][table]:
        registers = model['models'][table]['register']
        for register in registers:
            if not register['func']: #register = last value, no need to delete
                return '', '', ''
            reg_get += create_go_delete_registers(register, gv)
    if 'rz_register' in model['models'][table]:
        registers = model['models'][table]['rz_register']
        for register in registers:
            if not register['func']: #register = last value, no need to delete
                return '', '', ''
            reg_get += create_go_delete_registers(register, gv)
    rel_delete = ''
    if 'related' in model['models'][table]:
        relateds = model['models'][table]['related']
        for related in relateds:
            r = create_go_delete_relateds(related, gv)
            rel_delete += r

    g = f'''
        func {gtype}Delete(id int, tx *sql.Tx, isUnRealize bool) ({gtype}, error) {{
            needCommit := false
            var err error
            var {gv} {gtype}
            if tx == nil {{
                tx, err = db.Begin()
                if err != nil {{
                    return {gv}, err
                }}
                needCommit = true
                defer tx.Rollback()
            }}
            {gv}, err = {gtype}Get(id, tx)
            if err != nil {{
                return {gv}, err
            }}
            {complex_reg}{reg_get}{rel_delete}
            {check_unrealize}
            if needCommit {{
                err = tx.Commit()
                if err != nil {{
                    return {gv}, err
                }}
            }}
            {gv}.IsActive = false
            return {gv}, nil
        }}
    '''
    return g, h, m


def create_go_filter_int(table, keys, model):
    right = model['models'][table]['rights'] + '_READ'
    gtype = to_go(table)
    gv = table[0]

    m = f'''
        r.HandleFunc("/{table}_filter_int/{{fs}}/{{id:[0-9]+}}",
            WrapAuth(Get{gtype}ByFilterInt, {right})).Methods("GET")
        '''

    h = f'''
        func Get{gtype}ByFilterInt(req Req) {{
            req.Respond({gtype}GetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly, nil))
        }}
        '''

    g = f'''
        func {gtype}GetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]{gtype}, error) {{
            {create_go_filter(table, keys, gv, gtype)}
        }}
        '''
    return g, h, m


def create_go_filter_str(table, keys, model):
    right = model['models'][table]['rights'] + '_READ'
    gtype = to_go(table)
    gv = table[0]

    m = f'''
        r.HandleFunc("/{table}_filter_str/{{fs}}/{{fs2}}",
            WrapAuth(Get{gtype}ByFilterStr, {right})).Methods("GET")
        '''

    h = f'''
        func Get{gtype}ByFilterStr(req Req) {{
            req.Respond({gtype}GetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly, nil))
        }}
        '''

    g = f'''
        func {gtype}GetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool, tx *sql.Tx) ([]{gtype}, error) {{
            {create_go_filter(table, keys, gv, gtype)}
        }}
        '''
    return g, h, m

def make_field_validation(table, keys):
    gtype = to_go(table)
    s = '"'
    for k in keys:
        s += f'{k}", "'
    s = s[:-3]
    res = f'''
    func {gtype}TestForExistingField(fieldName string) bool {{
        fields := []string{{{s}}}
        for _, f := range fields {{
            if fieldName == f {{
                return true
            }}
        }}
        return false
    }}
    '''
    return res

def create_go_filter(table, keys, gv, gtype):
    f = f'''
    if !{gtype}TestForExistingField(field) {{
        return nil, errors.New("field not exist")
    }}
    var err error
    query := fmt.Sprintf("SELECT * FROM {table} WHERE %s=?", field)
    if deletedOnly {{
        query += "  AND is_active = 0"
    }} else if !withDeleted {{
        query += "  AND is_active = 1"
    }}

    var rows *sql.Rows
    if tx != nil {{
        rows, err = tx.Query(query, param)
    }} else {{
        rows, err = db.Query(query, param)
    }}
    if err != nil {{
        return nil, err
    }}
    defer rows.Close()
    res := []{gtype}{{}}
    for rows.Next() {{
        var {gv} {gtype}
        if err := rows.Scan(
    {list_of_pointers(keys, gv)}
        ); err != nil {{
            return nil, err
        }}
        res = append(res, {gv})
    }}
    return res, nil
    '''
    return f

# /field/id/date
def create_go_sum_before(table, keys, sum_field, model):
    right = model['models'][table]['rights'] + '_READ'
    gtype = to_go(table)
    gsum_field = to_go(sum_field)
    gv = table[0]

    m = f'''
        r.HandleFunc("/{table}_of_{sum_field}_sum_before/{{fs}}/{{id:[0-9]+}}/{{fs2}}",
            WrapAuth(Get{gtype}{gsum_field}SumBefore, {right})).Methods("GET")
    '''

    h = f'''
func Get{gtype}{gsum_field}SumBefore(req Req) {{
    req.Respond({gtype}{gsum_field}GetSumBefore(req.StrParam, req.IntParam, req.Str2Param))
}}
    '''

    g = f'''
func {gtype}{gsum_field}GetSumBefore(field string, id int, date string) (map[string]int, error) {{
    query := fmt.Sprintf("SELECT SUM({sum_field}) FROM {table} WHERE is_active = 1 AND %s = ? AND created_at <= ?", field)
    var sum int
    row := db.QueryRow(query, id, date)
    err := row.Scan(&sum)
    if err != nil {{
        return map[string]int{{"sum": 0}}, nil
    }}
    return map[string]int{{"sum": sum}}, nil
}}
'''
    return g, h, m

# /filter_field1/id1/filter_field2/id2
def create_go_sum_filter(table, keys, sum_field, model):
    right = model['models'][table]['rights'] + '_READ'
    gtype = to_go(table)
    gsum_field = to_go(sum_field)
    gv = table[0]

    m = f'''
        r.HandleFunc("/{table}_sum_filter_by/{{fs}}/{{id:[0-9]+}}/{{fs2}}/{{id2:[0-9]+}}",
            WrapAuth(Get{gtype}SumByFilter, {right})).Methods("GET")
    '''

    h = f'''
func Get{gtype}SumByFilter(req Req) {{
    req.Respond({gtype}GetSumByFilter(req.StrParam, req.IntParam, req.Str2Param, req.Int2Param))
}}
    '''

    g = f'''
func {gtype}GetSumByFilter(field string, id int, field2 string, id2 int) (map[string]int, error) {{
    query := ""
    var row *sql.Row
    if field2 == "-" && id2 == 0 {{
        query = fmt.Sprintf("SELECT SUM({sum_field}) FROM {table} WHERE is_active = 1 AND %s = ?", field)
        row = db.QueryRow(query, id)
    }} else {{
        query = fmt.Sprintf("SELECT SUM({sum_field}) FROM {table} WHERE is_active = 1 AND %s = ? AND %s = ?", field, field2)
        row = db.QueryRow(query, id, id2)
    }}
    var sum int
    err := row.Scan(&sum)
    if err != nil {{
        return map[string]int{{"sum": 0}}, nil
    }}
    return map[string]int{{"sum": sum}}, nil
}}
'''
    return g, h, m

def create_go_between(table, keys, between, model):
    right = model['models'][table]['rights'] + '_READ'
    gtype = to_go(table)
    gbetween = to_go(between)
    gv = table[0]
    d = model[table][between]['def']
    between_query1 = '{fs}'
    between_query2 = '{fs2}'
    between_param1 = 'req.StrParam'
    between_param2 = 'req.Str2Param'
    if type(d) == int:
        between_type = 'int'
        between_query1 = '{id:[0-9]+}'
        between_query2 = '{id:[0-9]+}'
        between_param1 = 'req.IntParam'
        between_param2 = 'req.Int2Param'
    else:
        between_type = 'string'

    m = f'''
        r.HandleFunc("/{table}_between_{between}/{between_query1}/{between_query2}",
            WrapAuth(Get{gtype}Between{gbetween}, {right})).Methods("GET")
    '''

    h = f'''
func Get{gtype}Between{gbetween}(req Req) {{
    req.Respond({gtype}GetBetween{gbetween}({between_param1}, {between_param2}, req.WithDeleted, req.DeletedOnly))
}}
    '''

    g = f'''

func {gtype}GetBetween{gbetween}({between}1, {between}2 {between_type}, withDeleted bool, deletedOnly bool) ([]{gtype}, error) {{
    query := "SELECT * FROM {table} WHERE {between} BETWEEN ? and ?"
    if deletedOnly {{
        query += "  AND is_active = 0"
    }} else if !withDeleted {{
        query += "  AND is_active = 1"
    }}

    rows, err := db.Query(query, {between}1, {between}2)
    if err != nil {{
        return nil, err
    }}
    defer rows.Close()
    res := []{gtype}{{}}
    for rows.Next() {{
        var {gv} {gtype}
        if err := rows.Scan(
    {list_of_pointers(keys, gv)}
        ); err != nil {{
            return nil, err
        }}
        res = append(res, {gv})
    }}
    return res, nil
}}
'''
    return g, h, m

def create_joins(table, find):
    if len(find.keys()) == 1:
        return ''
    g = ''
    for k in find.keys():
        if k != table:
            g += f'\t\t\t\tJOIN {k} on {table}.id = {k}.{table}_id\n'
    return g[4:]


def create_active(find):
    g = ''
    for k in find.keys():
        g += f'\t\t\t\tAND {k}.is_active=1\n'
    return 'WHERE ' + g[8:]

def create_search(find):
    g = ''
    i = 0
    for k, v in find.items():
        j = 0
        for field in v:
            if field.startswith('-'):
                continue
            else:
                g += f'\t\t\t\t{"OR" if i > 0 or j > 0 else ""} {k}.{field} LIKE ? \n'
            j += 1
        i += 1
    return 'AND (\n' + g + '\t\t\t)'

def create_not_search(find):
    g = ''
    i = 0
    for k, v in find.items():
        j = 0
        for field in v:
            if field.startswith('-'):
                g += f'\t\t\t\t{"OR" if i > 0 or j > 0 else ""} {k}.{field[1:]} LIKE ? \n'
            else:
                continue
            j += 1
        if j == 0:
            continue
        i += 1
    if g:
        return 'AND NOT (\n' + g + '\t\t\t)'
    else:
        return ''

finds = {}
finds['WProjectFindByProjectInfoContragentNoSearchContactNoSearch'] = '''
SELECT project.*, project_group.name, user.name, contragent.name, contact.name, project_type.name, project_status.name FROM project
    JOIN project_group ON project.project_group_id = project_group.id
    JOIN user ON project.user_id = user.id
    JOIN contragent on contragent.id = project.contragent_id
    JOIN contact on contact.id = project.contact_id
    JOIN project_type ON project.project_type_id = project_type.id
    JOIN project_status ON project.project_status_id = project_status.id
                WHERE project.is_active=1
                AND contragent.is_active=1
                AND contact.is_active=1
                AND (
                 project.info LIKE ?
            )
                AND NOT (
                 contragent.search LIKE ?
                OR contact.search LIKE ?
            );'''

finds['ProjectFindByProjectInfoContragentNoSearchContactNoSearch'] = '''
       SELECT project.* FROM project
        JOIN contragent on contragent.id = project.contragent_id
        JOIN contact on contact.id = project.contact_id
        WHERE project.info LIKE ?
        AND NOT (contragent.search LIKE ? OR contact.search LIKE ?);'''



finds['ContragentFindByContragentSearchContactSearch'] = '''
    SELECT DISTINCT contragent.* FROM contragent
    JOIN contact on contragent.id = contact.contragent_id
    WHERE
    contragent.search LIKE ?
    OR contact.search LIKE ?;'''



finds['WContragentFindByContragentSearchContactSearch'] = '''
        SELECT DISTINCT contragent.*, contragent_group.name FROM contragent
        JOIN contragent_group ON contragent.contragent_group_id =   contragent_group.id
        JOIN contact on contragent.id = contact.contragent_id
        WHERE contragent.is_active=1
        AND contact.is_active=1
        AND (contragent.search LIKE ? OR contact.search LIKE ?);'''

def create_go_find(table, keys, find, model):
    right = model['models'][table]['rights'] + '_READ'
    gtype = to_go(table)
    gfind = ''
    fnd = ''
    for k, v in find.items():
        gfind += to_go(k)
        fnd += f"_{k}"
        for i in v:
            gfind += ('No' + to_go(i[1:]) if i.startswith('-') else to_go(i))
            fnd += f"_{'no_' + i[1:] if i.startswith('-') else i}"

    gv = table[0]

    m = f'''
        r.HandleFunc("/find_{table}{fnd}/{{fs}}",
            WrapAuth(Get{gtype}FindBy{gfind}, {right})).Methods("GET")
    '''

    h = f'''
func Get{gtype}FindBy{gfind}(req Req) {{
    req.Respond({gtype}FindBy{gfind}(req.StrParam))
}}
    '''

    func_name = f'{gtype}FindBy{gfind}'
    find_query = finds[func_name]

    g = f'''

func {func_name}(fs string) ([]{gtype}, error) {{
    fs = "%" + fs + "%"

    query := `{find_query};`

    rows, err := db.Query(query{(', fs'*len(find))})

    if err != nil {{
        return nil, err
    }}
    defer rows.Close()

    res := []{gtype}{{}}
    for rows.Next() {{
        var {gv} {gtype}
        if err := rows.Scan(
    {list_of_pointers(keys, gv)}
        ); err != nil {{
            return nil, err
        }}
        res = append(res, {gv})
    }}
    return res, nil
}}
'''
    return g, h, m


# -------------------------WWWWWWWWWWWWWWWWWW-----------------------------

def create_go_model_w(table, model):
    g = ''
    g += f'\n\ntype W{to_go(table)} struct ' + '{'
    for k, v in model[table].items():
        d = v['def']
        if type(d) == int:
            g += f'\n\t{to_go(k)} int `json:"{k}"`'
        elif type(d) == float:
            g += f'\n\t{to_go(k)} float64 `json:"{k}"`'
        elif type(d) == bool:
            g += f'\n\t{to_go(k)} bool `json:"{k}"`'
        else:
            g += f'\n\t{to_go(k)} string `json:"{k}"`'

    for k in model[table].keys():
        if k.endswith('_id'):
            table_name = '_'.join(k.split('_')[:-1])
            g += f'\n\t{to_go(table_name)} string `json:"{table_name}"`'

    g += '\n}'
    return g

def create_add_joins(type, keys):
    add_sel = ''
    add_join = ''
    for k in keys:
        if k.endswith('_id'):
            table_name = '_'.join(k.split('_')[:-1])

            if table_name.endswith('2'):
                table_name_orig = table_name[:-1]
                add_sel += f', IFNULL({table_name}.name, "")'
                add_join += f'\n\tLEFT JOIN {table_name_orig} AS {table_name} ON {type}.{k} = {table_name}.id'

            elif table_name == type:
                add_sel += f', IFNULL({table_name[:2]}.name, "")'
                add_join += f'\n\tLEFT JOIN {table_name} AS {table_name[:2]} ON {type}.{k} = {table_name[:2]}.id'
            else:
                add_sel += f', IFNULL({table_name}.name, "")'
                add_join += f'\n\tLEFT JOIN {table_name} ON {type}.{k} = {table_name}.id'

    return add_sel, add_join


def list_of_pointers_w(keys, gv):
    g = ''

    for k in keys:
        g += f'\n\t\t&{gv}.{to_go(k)},'

    for k in keys:
        if k.endswith('_id'):
            table_name = '_'.join(k.split('_')[:-1])
            g += f'\n\t\t&{gv}.{to_go(table_name)},'


    return g[2:]


def create_go_get_w(table, keys, model):
    right = model['models'][table]['rights'] + '_READ'
    gtype = to_go(table)
    gv = table[0]
    m = f'''
        r.HandleFunc("/w_{table}/{{id:[0-9]+}}",
        WrapAuth(GetW{gtype}, {right})).Methods("GET")
    '''

    h = f'''
func GetW{gtype}(req Req) {{
    req.Respond(W{gtype}Get(req.IntParam))
}}
    '''
    add_sel, add_join = create_add_joins(table, keys)

    g = f'''

func W{gtype}Get(id int) (W{gtype}, error) {{
    var {gv} W{gtype}
    row := db.QueryRow(`SELECT {table}.*{add_sel} FROM {table}{add_join} WHERE {table}.id=?`, id)
    err := row.Scan(
    {list_of_pointers_w(keys, gv)}
    )
    return {gv}, err
}}
'''
    return g, h, m


def create_go_gets_w(table, keys, model):
    right = model['models'][table]['rights'] + '_READ'
    gtype = to_go(table)
    gv = table[0]

    m = f'''
        r.HandleFunc("/w_{table}_get_all",
            WrapAuth(GetW{gtype}All, {right})).Methods("GET")
    '''

    h = f'''
func GetW{gtype}All(req Req) {{
    req.Respond(W{gtype}GetAll(req.WithDeleted, req.DeletedOnly))
}}
    '''
    add_sel, add_join = create_add_joins(table, keys)
    g = f'''

func W{gtype}GetAll(withDeleted bool, deletedOnly bool) ([]W{gtype}, error) {{
    query := `SELECT {table}.*{add_sel} FROM {table}{add_join}`
    if deletedOnly {{
        query += "  WHERE {table}.is_active = 0"
    }} else if !withDeleted {{
        query += "  WHERE {table}.is_active = 1"
    }}

    rows, err := db.Query(query)
    if err != nil {{
        return nil, err
    }}
    defer rows.Close()
    res := []W{gtype}{{}}
    for rows.Next() {{
        var {gv} W{gtype}
        if err := rows.Scan(
    {list_of_pointers_w(keys, gv)}
        ); err != nil {{
            return nil, err
        }}
        res = append(res, {gv})
    }}
    return res, nil
}}
'''
    return g, h, m


def create_go_between_w(table, keys, between, model):
    right = model['models'][table]['rights'] + '_READ'
    gtype = to_go(table)
    gbetween = to_go(between)
    gv = table[0]
    d = model[table][between]['def']
    between_query1 = '{fs}'
    between_query2 = '{fs2}'
    between_param1 = 'req.StrParam'
    between_param2 = 'req.Str2Param'
    if type(d) == int:
        between_type = 'int'
        between_query1 = '{id:[0-9]+}'
        between_query2 = '{id:[0-9]+}'
        between_param1 = 'req.IntParam'
        between_param2 = 'req.Int2Param'
    else:
        between_type = 'string'

    m = f'''
        r.HandleFunc("/w_{table}_between_{between}/{between_query1}/{between_query2}",
            WrapAuth(GetW{gtype}Between{gbetween}, {right})).Methods("GET")
    '''

    h = f'''
func GetW{gtype}Between{gbetween}(req Req) {{
    req.Respond(W{gtype}GetBetween{gbetween}({between_param1}, {between_param2}, req.WithDeleted, req.DeletedOnly))
}}
    '''
    add_sel, add_join = create_add_joins(table, keys)
    g = f'''

func W{gtype}GetBetween{gbetween}({between}1, {between}2 {between_type}, withDeleted bool, deletedOnly bool) ([]W{gtype}, error) {{
    query := `SELECT {table}.*{add_sel} FROM {table}{add_join} WHERE ({table}.{between} BETWEEN ? AND ?)`
    if deletedOnly {{
        query += "  AND {table}.is_active = 0"
    }} else if !withDeleted {{
        query += "  AND {table}.is_active = 1"
    }}

    rows, err := db.Query(query, {between}1, {between}2)
    if err != nil {{
        return nil, err
    }}
    defer rows.Close()
    res := []W{gtype}{{}}
    for rows.Next() {{
        var {gv} W{gtype}
        if err := rows.Scan(
    {list_of_pointers_w(keys, gv)}
        ); err != nil {{
            return nil, err
        }}
        res = append(res, {gv})
    }}
    return res, nil
}}
'''
    return g, h, m


def create_go_between_up_w(table, keys, between_up, model):
    up_table, between = between_up.split('.')
    right = model['models'][table]['rights'] + '_READ'
    gtype = to_go(table)
    gbetween = to_go(between)
    gv = table[0]
    d = model[up_table][between]['def']
    between_query1 = '{fs}'
    between_query2 = '{fs2}'
    between_param1 = 'req.StrParam'
    between_param2 = 'req.Str2Param'
    if type(d) == int:
        between_type = 'int'
        between_query1 = '{id:[0-9]+}'
        between_query2 = '{id:[0-9]+}'
        between_param1 = 'req.IntParam'
        between_param2 = 'req.Int2Param'
    else:
        between_type = 'string'

    m = f'''
        r.HandleFunc("/w_{table}_between_up_{between}/{between_query1}/{between_query2}",
            WrapAuth(GetW{gtype}BetweenUp{gbetween}, {right})).Methods("GET")
    '''

    h = f'''
func GetW{gtype}BetweenUp{gbetween}(req Req) {{
    req.Respond(W{gtype}GetBetweenUp{gbetween}({between_param1}, {between_param2}, req.WithDeleted, req.DeletedOnly))
}}
    '''
    add_sel, add_join = create_add_joins(table, keys)
    g = f'''

func W{gtype}GetBetweenUp{gbetween}({between}1, {between}2 {between_type}, withDeleted bool, deletedOnly bool) ([]W{gtype}, error) {{
    query := `SELECT {table}.*{add_sel} FROM {table}{add_join}
                WHERE ({up_table}.{between} BETWEEN ? AND ?)`
    if deletedOnly {{
        query += "  AND {table}.is_active = 0"
    }} else if !withDeleted {{
        query += "  AND {table}.is_active = 1"
    }}

    rows, err := db.Query(query, {between}1, {between}2)
    if err != nil {{
        return nil, err
    }}
    defer rows.Close()
    res := []W{gtype}{{}}
    for rows.Next() {{
        var {gv} W{gtype}
        if err := rows.Scan(
    {list_of_pointers_w(keys, gv)}
        ); err != nil {{
            return nil, err
        }}
        res = append(res, {gv})
    }}
    return res, nil
}}
'''
    return g, h, m



def create_go_filter_w_int(table, keys, model):
    right = model['models'][table]['rights'] + '_READ'
    gtype = to_go(table)
    gv = table[0]

    m = f'''
        r.HandleFunc("/w_{table}_filter_int/{{fs}}/{{id:[0-9]+}}",
            WrapAuth(GetW{gtype}ByFilterInt, {right})).Methods("GET")
    '''

    h = f'''
func GetW{gtype}ByFilterInt(req Req) {{
    req.Respond(W{gtype}GetByFilterInt(req.StrParam, req.IntParam, req.WithDeleted, req.DeletedOnly))
}}
    '''
    add_sel, add_join = create_add_joins(table, keys)
    g = f'''

func W{gtype}GetByFilterInt(field string, param int, withDeleted bool, deletedOnly bool) ([]W{gtype}, error) {{
    {create_go_filter_w(table, keys, gv, gtype)}
}}
'''
    return g, h, m

def create_go_filter_w_str(table, keys, model):
    right = model['models'][table]['rights'] + '_READ'
    gtype = to_go(table)
    gv = table[0]

    m = f'''
        r.HandleFunc("/w_{table}_filter_str/{{fs}}/{{fs2}}",
            WrapAuth(GetW{gtype}ByFilterStr, {right})).Methods("GET")
    '''

    h = f'''
func GetW{gtype}ByFilterStr(req Req) {{
    req.Respond(W{gtype}GetByFilterStr(req.StrParam, req.Str2Param, req.WithDeleted, req.DeletedOnly))
}}
    '''
    add_sel, add_join = create_add_joins(table, keys)
    g = f'''

func W{gtype}GetByFilterStr(field string, param string, withDeleted bool, deletedOnly bool) ([]W{gtype}, error) {{
    {create_go_filter_w(table, keys, gv, gtype)}
}}
'''
    return g, h, m

def create_go_filter_w(table, keys, gv, gtype):
    add_sel, add_join = create_add_joins(table, keys)
    f = f'''
    if !{gtype}TestForExistingField(field) {{
        return nil, errors.New("field not exist")
    }}
    query := fmt.Sprintf(`SELECT {table}.*{add_sel} FROM {table}{add_join} WHERE {table}.%s=?`, field)
    if deletedOnly {{
        query += "  AND {table}.is_active = 0"
    }} else if !withDeleted {{
        query += "  AND {table}.is_active = 1"
    }}
    rows, err := db.Query(query, param)
    if err != nil {{
        return nil, err
    }}
    defer rows.Close()
    res := []W{gtype}{{}}
    for rows.Next() {{
        var {gv} W{gtype}
        if err := rows.Scan(
    {list_of_pointers_w(keys, gv)}
        ); err != nil {{
            return nil, err
        }}
        res = append(res, {gv})
    }}
    return res, nil
    '''
    return f

def create_add_joins_for_find(type, keys, find_keys):
    add_sel = ''
    add_join = ''
    for k in keys:
        if k.endswith('_id'):
            table_name = '_'.join(k.split('_')[:-1])
            add_sel += f', {table_name}.name'
            if table_name in find_keys:
                continue
            add_join += f'\n\tJOIN {table_name} ON {type}.{k} = {table_name}.id'

    return add_sel, add_join

def create_go_find_w(table, keys, find, model):
    right = model['models'][table]['rights'] + '_READ'
    gtype = to_go(table)
    gfind = ''
    fnd = ''
    for k, v in find.items():
        gfind += to_go(k)
        fnd += f"_{k}"
        for i in v:
            gfind += ('No' + to_go(i[1:]) if i.startswith('-') else to_go(i))
            fnd += f"_{'no_' + i[1:] if i.startswith('-') else i}"

    gv = table[0]

    m = f'''
        r.HandleFunc("/w_find_{table}{fnd}/{{fs}}",
            WrapAuth(GetW{gtype}FindBy{gfind}, {right})).Methods("GET")
    '''

    h = f'''
func GetW{gtype}FindBy{gfind}(req Req) {{
    req.Respond(W{gtype}FindBy{gfind}(req.StrParam))
}}
    '''
    add_sel, add_join = create_add_joins_for_find(table, keys, [find.keys()])
    func_name = f'W{gtype}FindBy{gfind}'
    find_query = finds[func_name]
    g = f'''

func {func_name}(fs string) ([]W{gtype}, error) {{
    fs = "%" + fs + "%"

    query := `{find_query}`

    rows, err := db.Query(query{(', fs'*len(find))})

    if err != nil {{
        return nil, err
    }}
    defer rows.Close()

    res := []W{gtype}{{}}
    for rows.Next() {{
        var {gv} W{gtype}
        if err := rows.Scan(
    {list_of_pointers_w(keys, gv)}
        ); err != nil {{
            return nil, err
        }}
        res = append(res, {gv})
    }}
    return res, nil
}}
'''
    return g, h, m

if __name__ == '__main__':
    with open ('models.json', "r") as f:
        model = json.loads(f.read())

    tables = model['models'].keys()

    create_go_models(model, tables)
