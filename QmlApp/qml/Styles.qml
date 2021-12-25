import QtQuick
pragma Singleton

import client.api.mooody.me

QtObject {
    property color background_1: AppSettings.darkMode ? "#1f4042" : "#00b4c3"
    property color background_2: AppSettings.darkMode ? "#003b1b" : "#008b40"

    property color background_pure: AppSettings.darkMode ? "#34593d" : "#e3e3e3"
    property color background_pure_border: AppSettings.darkMode ? "#212121" : "#b6b6b6"

    property color text: AppSettings.darkMode ? "#c0c0c0" : "#2d2d2d"

    property color button_on: AppSettings.darkMode ? "#33712d" : "#1adf00"
    property color button_on_border: AppSettings.darkMode ? "#00241b" : "#188300"

    property color button_off: AppSettings.darkMode ? "#84501d" : "#eb7500"
    property color button_off_border: AppSettings.darkMode ? "#392100" : "#834900"
}
