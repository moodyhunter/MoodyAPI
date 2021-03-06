import QtQuick
import QtQuick.Controls
import QtQuick.Layouts

import client.api.mooody.me

Rectangle {
    id: root

    MouseArea {
        anchors.fill: parent
        hoverEnabled: true
        preventStealing: true
        enabled: root.state == "opened"
    }

    color: color_background_pure
    radius: 20
    border.color: color_background_pure_border
    border.width: 5

    implicitHeight: layout.implicitHeight + 2 * layout.anchors.margins
    width: 325

    opacity: 0
    scale: 0

    states: [
        State {
            name: "closed"
            PropertyChanges {
                target: root
                opacity: 0
                scale: 0
            }
        },
        State {
            name: "opened"
            PropertyChanges {
                target: root
                opacity: 1
                scale: 1
            }
        }
    ]

    transitions: Transition {
        NumberAnimation {
            duration: 250
            properties: "opacity,scale"
            easing.type: Easing.OutQuart
        }
    }

    property int standardFontSize: 16

    function open() {
        root.enabled = true
        root.state = "opened"
    }

    function close(accept) {
        root.enabled = false
        root.state = "closed"
        if (accept)
            AppCore.connectToServer(AppSettings.apiHost, AppSettings.apiSecret,
                                    AppSettings.disableTLS)
    }

    component Spacer: Item {
        property int spacerHeight: 15
        implicitHeight: spacerHeight
        Layout.fillWidth: true
    }

    component BackgroundRectangle: Rectangle {
        color: AppSettings.darkMode ? Qt.darker(
                                          color_background_pure) : Qt.lighter(
                                          color_background_pure)
        border.color: {
            if (AppSettings.darkMode)
                Qt.lighter(color_background_pure_border, parent.focus ? 4 : 1)
            else
                Qt.darker(color_background_pure_border, parent.focus ? 4 : 1)
        }
        border.width: 2
        radius: 5
    }

    ColumnLayout {
        id: layout
        anchors.fill: parent
        anchors.margins: 20

        RowLayout {
            Layout.fillWidth: true

            Label {
                Layout.fillWidth: true
                font.pixelSize: standardFontSize
                font.family: "System-ui"
                font.bold: true
                text: "DarkMode"
                color: color_text
            }

            Switch {
                checked: AppSettings.darkMode
                onCheckedChanged: AppSettings.darkMode = checked
            }
        }

        RowLayout {
            Layout.fillWidth: true

            Label {
                Layout.fillWidth: true
                font.pixelSize: standardFontSize
                font.family: "System-ui"
                font.bold: true
                text: "Disable TLS"
                color: color_text
            }

            Switch {
                checked: AppSettings.disableTLS
                onCheckedChanged: AppSettings.disableTLS = checked
            }
        }

        Label {
            Layout.fillWidth: true
            font.pixelSize: standardFontSize
            font.family: "System-ui"
            font.bold: true
            text: "API Server"
            color: color_text
        }

        TextField {
            padding: 10
            background: BackgroundRectangle {}
            Layout.fillWidth: true
            Layout.leftMargin: 10
            Layout.rightMargin: 10
            selectByMouse: true
            wrapMode: TextEdit.WrapAnywhere
            color: color_text
            font.family: fixedFont
            text: AppSettings.apiHost
            onTextChanged: {
                AppSettings.apiHost = text
            }
        }

        Spacer {}

        Label {
            Layout.fillWidth: true
            font.pixelSize: standardFontSize
            font.family: "System-ui"
            font.bold: true
            color: color_text
            text: "API Secret"
        }
        TextField {
            padding: 10
            background: BackgroundRectangle {}
            Layout.fillWidth: true
            Layout.leftMargin: 10
            Layout.rightMargin: 10
            selectByMouse: true
            wrapMode: TextEdit.WrapAnywhere
            color: color_text
            font.family: fixedFont
            text: AppSettings.apiSecret
            onTextChanged: {
                AppSettings.apiSecret = text
            }
        }

        Spacer {
            spacerHeight: 30
        }

        RowLayout {
            Layout.fillWidth: true

            GradientButton {
                Layout.fillWidth: true
                buttonBorderWidth: 1
                color1: color_button_off
                borderColor: color_button_off_border
                buttonSize: 60
                fontSize: 18
                text: "Cancel"
                onClicked: {
                    root.close(false)
                }
            }

            GradientButton {
                Layout.fillWidth: true
                buttonBorderWidth: 1
                color1: color_button_on
                borderColor: color_button_on_border
                buttonSize: 60
                fontSize: 18
                text: "OK"
                onClicked: {
                    root.close(true)
                }
            }
        }
    }
}
