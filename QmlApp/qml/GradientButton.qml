import QtQuick
import QtQuick.Controls

import QtGraphicalEffects

import client.api.mooody.me

Control {
    signal clicked

    property string text: "Button"
    property color textColor: "#2d2d2d"
    property color textColorDark: "#cccccc"

    hoverEnabled: PlatformHoverEnabled

    property color color1: "black"
    property color colorh1: Qt.lighter(color1, 1.3)
    property color colorc1: Qt.darker(color1, 1.3)

    property color borderColor: "grey"
    property color borderColorH: Qt.lighter(borderColor, 1.3)
    property color borderColorC: Qt.darker(borderColor, 1.3)

    implicitWidth: rootWindow.buttonSize
    implicitHeight: rootWindow.buttonSize / 2
    id: root

    Rectangle {
        anchors.fill: parent
        id: rectangle

        border.color: borderColor
        border.width: 5

        radius: rootWindow.buttonSize / 3
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
        font.pixelSize: rootWindow.buttonSize / 8
        color: AppCore.DarkMode ? root.textColorDark : root.textColor
        text: root.text
    }
}
