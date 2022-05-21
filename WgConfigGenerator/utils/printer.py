import colored
from colored import stylize

COLOR_YELLOW = colored.fg("light_yellow")
COLOR_BOLD_TEXT = colored.fg("white")
STYLE_BOLD_TEXT = colored.attr("bold")
STYLE_RESET = colored.attr("reset")

ITEM_FORMAT_BOLD = COLOR_YELLOW + STYLE_BOLD_TEXT + "{}" + COLOR_BOLD_TEXT
ITEM_FORMAT_NORMAL = COLOR_YELLOW + "{}" + STYLE_RESET


def print_status(arrow: str, color: str, msg: str, item: list[str], bold: bool = True, prefix: str = ""):
    formatted_items = []
    for i in item:
        formatted_items.append((ITEM_FORMAT_BOLD if bold else ITEM_FORMAT_NORMAL).format(i))

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
    print_status("=>", "green", msg, item)


def sub_progress(msg: str, item: list[str] = []):
    print_status("=>", "blue", msg, item, bold=False, prefix="  ")


def sub_sub_progress(msg: str, item: list[str] = []):
    print_status("->", "blue", msg, item, bold=False, prefix="  "*2)


def warning(msg: str, item: list[str] = []):
    print_status("=>", "yellow", msg, item)


def error(msg: str, item: list[str] = []):
    print_status("=>", "red", msg, item)
    print_status("=>", "red", "Program cannot continue, exiting...", [])
    exit(1)
