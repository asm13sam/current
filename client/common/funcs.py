from datetime import datetime
from data.model import Item

def update_cash_register(value: dict, field_name: str, is_add: bool = True, money_field='cash_sum'):
    reg = Item(field_name)
    reg.get(value[f'{field_name}_id'])
    if is_add:
        reg.value['total'] += value[money_field]
    else:
        reg.value['total'] -= value[money_field]
    reg.save()

def format_phone(p: str):
    if len(p) > 10:
        return p[:-10] + '-' + p[-10:-7] + '-' + p[-7:-4] + '-' + p[-4:-2] + '-' + p[-2:]
    if len(p) > 7:
        return p[:-7] + '-' + p[-7:-4] + '-' + p[-4:-2] + '-' + p[-2:]
    if len(p) > 4:
        return p[:-4] + '-' + p[-4:-2] + '-' + p[-2:]
    return p

def phonef(phone_number: str):
    res = []
    for phone in phone_number.split(','):
        res.append(format_phone(phone.strip()))
    return ', '.join(res)

def prepare_search_string(s: str):
    # Viber +380681317555
    # Telegram +380 98 742 78 18 or @Bondar_Yana
    # Email komar160386@ukr.net
    if s.startswith('+380 '):
        s = ''.join(s.split(' '))[3:]
    elif s.startswith('+38'):
        s = s[3:]
    elif s.startswith('@'):
        s = s[1:].lower()
    elif '@' in s:
        s = s[:s.index('@')]
    return s

def number_to_words(number):
    # Слова для одиниць
    units = ["", "одна", "дві", "три", "чотири", "п'ять", "шість", "сім", "вісім", "дев'ять"]

    # Слова для чисел від 10 до 19
    teens = ["десять", "одинадцять", "дванадцять", "тринадцять", "чотирнадцять",
             "п'ятнадцять", "шістнадцять", "сімнадцять", "вісімнадцять", "дев'ятнадцять"]

    # Слова для десятків
    tens = ["", "", "двадцять", "тридцять", "сорок", "п'ятдесят",
            "шістдесят", "сімдесят", "вісімдесят", "дев'яносто"]

    # Слова для сотень
    hundreds = ["", "сто", "двісті", "триста", "чотириста", "п'ятсот",
                "шістсот", "сімсот", "вісімсот", "дев'ятсот"]

    # Перевірка, чи число є нульовим
    if number == 0:
        return "нуль"

    # Витягуємо окремі частини числа
    hundreds_digit = number // 100
    tens_digit = number % 100 // 10
    units_digit = number % 10

    # Генеруємо текстове представлення числа
    text = hundreds[hundreds_digit] + " " if hundreds_digit != 0 else ""

    # Перевірка, чи є десятки від 10 до 19
    if tens_digit == 1:
        text += teens[units_digit]
    else:
        text += tens[tens_digit] + " " if tens_digit != 0 else ""
        text += units[units_digit] if units_digit != 0 else ""

    return text.strip()

def thousands_to_words(number):
    thousands = number // 1000
    if thousands:
        result = number_to_words(thousands)
        if thousands % 10 == 1 and thousands % 100 // 10 != 1:
            result += ' тисяча '
        elif thousands % 10 in (2, 3, 4) and thousands % 100 // 10 != 1:
            result += ' тисячі '
        else:
            result += ' тисяч '
    else:
        result = ''

    units = number % 1000
    result += number_to_words(units)
    if units % 10 == 1 and units % 100 // 10 != 1:
        result += ' гривня'
    elif units % 10 in (2, 3, 4) and units % 100 // 10 != 1:
        result += ' гривні'
    else:
        result += ' гривень'

    result += ' 00 копійок'
    return result

def dataiso_to_words(data_iso_str, only_year=False):
    months = ['', "січня", "лютого", "березня", "квітня", "травня",
             "червня", "липня", "серпня", "вересня", "жовтня", "листопада", "грудня"]
    dt = datetime.fromisoformat(data_iso_str)
    if only_year:
        date = f'{dt.year}р.'
    else:
        date = f'{dt.day} {months[dt.month]} {dt.year}р.'
    return date

def round_to(value: float, to: float):
    res = int(value/to)*to
    return res if to >= 1 else round(res, 2)

def id_generator():
    counter = 1
    while True:
        yield counter
        counter += 1

fake_id = id_generator()
