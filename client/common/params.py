MIN_SEARCH_STR_LEN = 3
ID_ROLE = 110
SORT_ROLE = 111
FULL_VALUE_ROLE = 112
IS_EXT_ROLE = 113
TABLE_BUTTONS = {'reload':'Оновити', 'create':'Створити', 'edit':'Редагувати', 'copy':'Копіювати', 'delete':'Видалити'}
VIRTUAL = 1
REORDERING = 2
# access constant
ACCESS = [
    "Авторизуватись",
    "Вийти з програми",
    "Читати інформацію про користувачів",
    "Створити користувача",
    "Редагувати користувача",
    "Видалити користувача",
    "Читати інформацію про контрагентів",
    "Створити контрагента",
    "Редагувати контрагента",
    "Видалити контрагента",
    "Читати документи",
    "Читати власні документи",
    "Створити документ",
    "Редагувати документи",
    "Редагувати власні документи",
    "Видалити документ",
    "Читати довідники",
    "Створити запис в довіднику",
    "Редагувати запис в довіднику",
    "Видалити запис в довіднику",
    "Читати власні рахунки",
    "Додати власний рахунок",
    "Редагувати власний рахунок",
    "Видалити власний рахунок",
    "Користуватися чатом",
    "Керувати налаштуваннями",
]
CONF_TABS = {
    'Конфігурація': ("host", 
                    "port", 
                    "makets_path", 
                    "bs_makets_path", 
                    "new_makets_path", 
                    "program",
                    ),
    "Виглад": (
                "theme",
                "font_size",
                "color",
                "form_names_color",
                "form_values_color",
                "form_bg_color",
                "info_names_color",
                "info_values_color",
                "info_bg_color",
                "unreleazed_color",
                "product_groups",
                ),
    "Checkbox": (
                "license_key",
                "cashier_pin",
                "checkbox_url",
                ),
    "Призначення": (
                "copycenter warehouse id",
                "whs_in cash id",
                "contragent to production",
                "contact to production",
                "contragent copycenter default",
                "contact copycenter default",
                "contragent for delivery",
                "contact for delivery",
                ),
    "Константи": (
                "warehouse persent",
                "measure pieces",
                "measure linear",
                "measure square",
                "ordering state in work",
                "ordering state ready",
                "ordering state taken",
                "ordering state canceled",
                "product_to_ordering state ready",
                "clients_group",
                ),
}