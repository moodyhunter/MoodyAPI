import QtQuick
import QtQuick.Controls

import QtGraphicalEffects

import client.api.mooody.me

Control {
    signal clicked

    property string text: "Button"

    hoverEnabled: PlatformHoverEnabled

    property color color1: "black"
    property color colorh1: Qt.lighter(color1, 1.3)
    property color colorc1: Qt.darker(color1, 1.3)

    property color borderColor: "grey"
    property color borderColorH: Qt.lighter(borderColor, 1.3)
    property color borderColorC: Qt.darker(borderColor, 1.3)

    property double buttonSize: 50
    property int buttonBorderWidth: 5

    property int fontSize: buttonSize / 8

    implicitWidth: buttonSize
    implicitHeight: buttonSize / 2

    id: root

    Rectangle {
        anchors.fill: parent
        id: rectangle

        border.color: borderColor
        border.width: buttonBorderWidth

        radius: buttonSize / 3
        color: color1

        states: [
            State {
                name: "pressed"
                when: mouse.pressed
                PropertyChanges {
                    target: rectangle
                    color: colorc1
                    border.color: borderColorC
                }
            },
            State {
                name: "hovered"
                when: root.hovered
                PropertyChanges {
                    target: rectangle
                    color: colorh1
                    border.color: borderColorH
                }
            }
        ]

        transitions: Transition {
            ColorAnimation {
                easing.type: Easing.OutSine
                duration: 140
            }
        }
    }

    MouseArea {
        anchors.fill: parent
        onClicked: root.clicked()
        id: mouse
    }

    Label {
        anchors.centerIn: parent
        font.pixelSize: fontSize
        color: color_text
        text: root.text
    }
}
