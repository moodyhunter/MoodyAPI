import QtQuick
import QtQuick.Controls
import QtQuick.Layouts

import QtGraphicalEffects

ApplicationWindow {
    width: 392
    height: 815
    visible: true
    title: qsTr("Moody Camera App")
    id: rootWindow

    LinearGradient {
        anchors.fill: parent
        start: Qt.point(0, 0)
        end: Qt.point(rootWindow.width, rootWindow.height)
        gradient: Gradient {
            GradientStop {
                position: 0.0
                color: "#00d8e9"
            }
            GradientStop {
                position: 1.0
                color: "#00c159"
            }
        }
    }

    readonly property double buttonSize: Math.min(width / 2.5, height / 4)

    component VerticalSpacer: Item {
        Layout.fillHeight: true
    }

    component GradientButton: Control {
        signal clicked
        hoverEnabled: true
        property color color1: "black"
        property color colorh1: Qt.lighter(color1)
        property color colorc1: Qt.darker(color1)

        property color color2: "white"
        property color colorh2: Qt.lighter(color2)
        property color colorc2: Qt.darker(color2)

        property color borderColor: "grey"
        property color borderColorH: Qt.lighter(borderColor)
        property color borderColorC: Qt.darker(borderColor)

        implicitWidth: buttonSize
        implicitHeight: buttonSize / 2
        id: root

        Rectangle {
            anchors.fill: parent
            id: rectangle

            border.color: mouse.pressed ? borderColorC : (root.hovered ? borderColorH : borderColor)
            border.width: 5

            radius: buttonSize / 3
        }

        LinearGradient {
            source: rectangle
            anchors.fill: rectangle
            anchors.margins: rectangle.border.width - 2

            start: Qt.point(0, 0)
            end: Qt.point(rectangle.width, rectangle.height)
            gradient: Gradient {
                GradientStop {
                    position: 0.0
                    color: mouse.pressed ? colorc1 : (root.hovered ? colorh1 : color1)
                }
                GradientStop {
                    position: 1.0
                    color: mouse.pressed ? colorc2 : (root.hovered ? colorh2 : color2)
                }
            }
        }

        MouseArea {
            anchors.fill: parent
            id: mouse
        }

        Label {
            anchors.centerIn: parent
            text: "test"
        }
    }

    ColumnLayout {
        anchors.fill: parent
        spacing: buttonSize / 4

        VerticalSpacer {}
        VerticalSpacer {}

        Label {
            font.pixelSize: buttonSize / 5
            font.family: "System-ui"
            font.bold: true
            text: "Camera Status"
            color: "#2d2d2d"
            horizontalAlignment: Qt.AlignHCenter
            Layout.alignment: Qt.AlignHCenter
        }

        Label {
            font.pixelSize: buttonSize / 3
            font.family: "System-ui"
            font.bold: true
            text: "ON"
            color: "#2d2d2d"
            horizontalAlignment: Qt.AlignHCenter
            Layout.alignment: Qt.AlignHCenter
        }

        VerticalSpacer {}

        GradientButton {
            color1: "#1adf00"
            color2: "#00c902"
            borderColor: "#188300"
            Layout.alignment: Qt.AlignHCenter

            onClicked: {

            }
        }

        GradientButton {
            color1: "#eb7500"
            color2: "#e24e00"
            borderColor: "#834900"
            Layout.alignment: Qt.AlignHCenter

            onClicked: {

            }
        }

        VerticalSpacer {}
        VerticalSpacer {}
        VerticalSpacer {}
    }
}
