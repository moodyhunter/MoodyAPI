import QtQuick
import QtQuick.Controls
import QtQuick.Layouts

import client.api.mooody.me

Rectangle {
    id: root

    color: Styles.background_pure
    radius: 20
    border.color: Styles.background_pure_border
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
        root.state = "opened"
    }

    function close() {
        root.state = "closed"
    }

    component Spacer: Item {
        property int spacerHeight: 20
        implicitHeight: spacerHeight
        Layout.fillWidth: true
        Layout.columnSpan: 2
    }

    GridLayout {
        id: layout
        anchors.fill: parent
        anchors.margins: 20
        columns: 2

        Label {
            Layout.fillWidth: true
            font.pixelSize: standardFontSize
            font.family: "System-ui"
            font.bold: true
            text: "DarkMode"
            color: Styles.text
        }

        Switch {
            Layout.fillWidth: true
            checked: AppSettings.darkMode
            onCheckedChanged: AppSettings.darkMode = checked
        }

        Label {
            Layout.columnSpan: 2
            Layout.fillWidth: true
            font.pixelSize: standardFontSize
            font.family: "System-ui"
            font.bold: true
            text: "API Server"
            color: Styles.text
        }

        TextEdit {
            Layout.columnSpan: 2
            Layout.fillWidth: true
            Layout.leftMargin: 10
            Layout.rightMargin: 10
            selectByMouse: true
            wrapMode: TextEdit.WrapAnywhere
            color: Styles.text
            font.family: fixedFont
            onTextChanged: {
                AppSettings.apiHost = text
            }
        }

        Spacer {}

        Label {
            Layout.columnSpan: 2
            Layout.fillWidth: true
            font.pixelSize: standardFontSize
            font.family: "System-ui"
            font.bold: true
            color: Styles.text
            text: "API Secret"
        }
        TextEdit {
            Layout.columnSpan: 2
            Layout.fillWidth: true
            Layout.leftMargin: 10
            Layout.rightMargin: 10
            selectByMouse: true
            wrapMode: TextEdit.WrapAnywhere
            color: Styles.text
            font.family: fixedFont
            text: AppSettings.apiSecret
            onTextChanged: {
                AppSettings.apiSecret = text
            }
        }

        Spacer {
            spacerHeight: 40
        }

        GradientButton {
            Layout.columnSpan: 2
            Layout.fillWidth: true
            buttonBorderWidth: 1
            color1: Styles.button_on
            borderColor: Styles.button_on_border
            buttonSize: 60
            fontSize: 18
            text: "OK"
            onClicked: {
                root.close()
            }
        }
    }
}
