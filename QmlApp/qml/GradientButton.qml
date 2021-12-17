import QtQuick
import QtQuick.Controls

import QtGraphicalEffects

Control {
    signal clicked

    property string text: "Button"
    property color textColor: "#2d2d2d"

    hoverEnabled: true
    property color color1: "black"
    property color colorh1: Qt.lighter(color1)
    property color colorc1: Qt.darker(color1)

    property color borderColor: "grey"
    property color borderColorH: Qt.lighter(borderColor)
    property color borderColorC: Qt.darker(borderColor)

    implicitWidth: rootWindow.buttonSize
    implicitHeight: rootWindow.buttonSize / 2
    id: root

    Rectangle {
        anchors.fill: parent
        id: rectangle

        border.color: borderColor
        border.width: 5

        radius: rootWindow.buttonSize / 3
        color: color1 // mouse.pressed ? colorc1 : (root.hovered ? colorh1 : color1)

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

    //    LinearGradient {
    //        source: rectangle
    //        anchors.fill: rectangle
    //        anchors.margins: rectangle.border.width - 2

    //        start: Qt.point(0, 0)
    //        end: Qt.point(rectangle.width, rectangle.height)
    //        gradient: Gradient {
    //            GradientStop {
    //                position: 0.0
    //                color:
    //            }
    //            GradientStop {
    //                position: 1.0
    //                color: mouse.pressed ? colorc2 : (root.hovered ? colorh2 : color2)
    //            }
    //        }
    //    }
    MouseArea {
        anchors.fill: parent
        onClicked: root.clicked()
        id: mouse
    }

    Label {
        anchors.centerIn: parent
        font.pixelSize: rootWindow.buttonSize / 8
        color: root.textColor
        text: root.text
    }
}
