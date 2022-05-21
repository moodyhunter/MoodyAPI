import colored
from colored import stylize

COLOR_YELLOW = colored.fg("light_yellow")
COLOR_BOLD_TEXT = colored.fg("white")
STYLE_BOLD_TEXT = colored.attr("bold")
STYLE_RESET = colored.attr("reset")

ITEM_FORMAT_BOLD = COLOR_YELLOW + STYLE_BOLD_TEXT + "{}" + COLOR_BOLD_TEXT
ITEM_FORMAT_NORMAL = COLOR_YELLOW + "{}" + STYLE_RESET


def print_status(arrow: str, color: str, msg: str, item: list[str], bold: bool = False, prefix: str = "", arrow_blink: bool = False):
    formatted_items = []
    for i in item:
        formatted_items.append((ITEM_FORMAT_BOLD if bold else ITEM_FORMAT_NORMAL).format(i))

    if arrow_blink:
        arrow = stylize(arrow, colored.attr("blink"))

    if bold:
        style = stylize(arrow, colored.fg(color) + STYLE_BOLD_TEXT)
        message = stylize(msg, COLOR_BOLD_TEXT + STYLE_BOLD_TEXT)
    else:
        style = stylize(arrow, colored.fg(color))
        message = msg

    if len(formatted_items) > 0:
        message = message.format(*formatted_items)

    print(prefix + style + " " + message, flush=True)


def progress(msg: str, item: list[str] = []):
    print_status("=>", "green", msg, item, bold=True)


def sub_progress(msg: str, item: list[str] = []):
    print_status("=>", "blue", msg, item,  prefix="  ")


def sub_sub_progress(msg: str, item: list[str] = []):
    print_status("->", "blue", msg, item,  prefix="  "*2)


def warning(msg: str, item: list[str] = []):
    print_status("=>", "yellow", msg, item, bold=True, arrow_blink=True)


def sub_warning(msg: str, item: list[str] = []):
    print_status("=>", "yellow", msg, item, prefix="  ", arrow_blink=True)


def error(msg: str, item: list[str] = []):
    print_status("=>", "red", msg, item, bold=True)
    print_status("=>", "red", "Program cannot continue, exiting...", [], bold=True)
    exit(1)
