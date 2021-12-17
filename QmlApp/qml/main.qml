import QtQuick
import QtQuick.Controls
import QtQuick.Layouts

import QtGraphicalEffects

import client.api.mooody.me

ApplicationWindow {
    width: 392
    height: 815
    visible: true
    title: qsTr("Moody Camera App")
    id: rootWindow
    readonly property double buttonSize: Math.min(width / 2.5, height / 4)

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

    SvgButton {
        anchors.right: parent.right
        anchors.top: parent.top
        anchors.margins: 10
        width: 36
        height: 36
        source: "/assets/settings.svg"

        onClicked: {
            console.log("Settings clicked")
        }
    }

    component VerticalSpacer: Item {
        Layout.fillHeight: true
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
            text: MoodyApi.CameraStatus ? "ON" : "OFF"
            color: "#2d2d2d"
            horizontalAlignment: Qt.AlignHCenter
            Layout.alignment: Qt.AlignHCenter
        }

        VerticalSpacer {}

        GradientButton {
            color1: "#1adf00"
            //            color2: "#00c902"
            borderColor: "#188300"
            Layout.alignment: Qt.AlignHCenter
            text: qsTr("Power On")

            onClicked: {
                MoodyApi.CameraStatus = true
            }
        }

        GradientButton {
            color1: "#eb7500"
            //            color2: "#e24e00"
            borderColor: "#834900"
            Layout.alignment: Qt.AlignHCenter
            text: qsTr("Power Off")

            onClicked: {
                MoodyApi.CameraStatus = false
            }
        }

        VerticalSpacer {}
        VerticalSpacer {}
        VerticalSpacer {}
    }
}
