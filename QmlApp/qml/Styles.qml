import QtQuick
pragma Singleton

import client.api.mooody.me

QtObject {
    property color background_1: AppSettings.darkMode ? "#1f4042" : "#00d8e9"
    property color background_2: AppSettings.darkMode ? "#003b1b" : "#00c159"
    property color text: AppSettings.darkMode ? "#a5a5a5" : "#2d2d2d"

    property color button_on: AppSettings.darkMode ? "#33712d" : "#1adf00"
    property color button_on_border: AppSettings.darkMode ? "#063026" : "#188300"

    property color button_off: AppSettings.darkMode ? "#84501d" : "#eb7500"
    property color button_off_border: AppSettings.darkMode ? "#452800" : "#834900"
}
