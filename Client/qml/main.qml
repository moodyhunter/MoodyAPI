import QtQuick
import QtQuick.Controls
import QtQuick.Layouts

import QtGraphicalEffects

import client.api.mooody.me

ApplicationWindow {
    property color color_background_1: AppSettings.darkMode ? "#1f4042" : "#00b4c3"
    property color color_background_2: AppSettings.darkMode ? "#003b1b" : "#008b40"

    property color color_background_pure: AppSettings.darkMode ? "#34593d" : "#e3e3e3"
    property color color_background_pure_border: AppSettings.darkMode ? "#212121" : "#b6b6b6"

    property color color_text: AppSettings.darkMode ? "#c0c0c0" : "#2d2d2d"

    property color color_button_on: AppSettings.darkMode ? "#33712d" : "#1adf00"
    property color color_button_on_border: AppSettings.darkMode ? "#00241b" : "#188300"

    property color color_button_off: AppSettings.darkMode ? "#84501d" : "#eb7500"
    property color color_button_off_border: AppSettings.darkMode ? "#392100" : "#834900"

    width: 392
    height: 815
    visible: true
    title: qsTr("Moody Camera App")
    id: rootWindow
    readonly property double standardSize: Math.min(width / 2.5, height / 4)

    LinearGradient {
        anchors.fill: parent
        start: Qt.point(0, 0)
        end: Qt.point(rootWindow.width, rootWindow.height)
        gradient: Gradient {
            GradientStop {
                position: 0.0
                color: color_background_1
            }
            GradientStop {
                position: 1.0
                color: color_background_2
            }
        }
    }

    SvgButton {
        anchors.right: parent.right
        anchors.top: parent.top
        anchors.margins: 10
        width: 36
        height: 36
        source: "/assets/settings.svg"

        onClicked: {
            settingsPopup.open()
        }
    }

    component VerticalSpacer: Item {
        Layout.fillHeight: true
    }

    ColumnLayout {
        anchors.fill: parent

        ColumnLayout {
            Layout.fillWidth: true
            Layout.fillHeight: true
            Layout.alignment: Qt.AlignHCenter

            VerticalSpacer {}
            VerticalSpacer {}

            Label {
                font.pixelSize: standardSize / 5
                font.family: "System-ui"
                font.bold: true
                text: "Camera Status"
                color: color_text
                horizontalAlignment: Qt.AlignHCenter
                Layout.alignment: Qt.AlignHCenter
                Layout.bottomMargin: standardSize / 8
            }

            Label {
                font.pixelSize: standardSize / 3
                font.family: "System-ui"
                font.bold: true
                text: AppCore.ServerConnected ? (AppCore.IsRecording ? "ON" : "OFF") : "N/A"
                color: color_text
                horizontalAlignment: Qt.AlignHCenter
                Layout.alignment: Qt.AlignHCenter
            }

            VerticalSpacer {}

            GradientButton {
                color1: color_button_on
                borderColor: color_button_on_border
                Layout.alignment: Qt.AlignHCenter
                text: qsTr("Power On")
                buttonSize: rootWindow.standardSize

                onClicked: {
                    AppCore.startRecording()
                }

                Layout.bottomMargin: standardSize / 8
            }

            GradientButton {
                color1: color_button_off
                borderColor: color_button_off_border
                Layout.alignment: Qt.AlignHCenter
                text: qsTr("Power Off")
                buttonSize: rootWindow.standardSize

                onClicked: {
                    AppCore.stopRecording()
                }
            }
            VerticalSpacer {}
            VerticalSpacer {}
        }

        Label {
            font.pixelSize: 16
            font.family: "System-ui"
            font.bold: true
            text: AppCore.ServerConnected ? "Connected" : "Disconnected"
            color: color_text
            horizontalAlignment: Qt.AlignHCenter
            verticalAlignment: Qt.AlignBottom
            Layout.alignment: Qt.AlignHCenter
            Layout.fillWidth: true
            Layout.fillHeight: true
            Layout.bottomMargin: 10
        }
    }

    DropShadow {
        source: settingsPopup
        anchors.fill: settingsPopup
        transparentBorder: true
        radius: 10
        opacity: settingsPopup.opacity
        scale: settingsPopup.scale
    }
    SettingsPanel {
        id: settingsPopup
        y: 100
        anchors.horizontalCenter: parent.horizontalCenter
    }
}
