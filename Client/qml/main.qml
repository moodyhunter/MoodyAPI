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
    readonly property double standardSize: Math.min(width / 2.5, height / 4)

    LinearGradient {
        anchors.fill: parent
        start: Qt.point(0, 0)
        end: Qt.point(rootWindow.width, rootWindow.height)
        gradient: Gradient {
            GradientStop {
                position: 0.0
                color: Styles.background_1
            }
            GradientStop {
                position: 1.0
                color: Styles.background_2
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
                color: Styles.text
                horizontalAlignment: Qt.AlignHCenter
                Layout.alignment: Qt.AlignHCenter
                Layout.bottomMargin: standardSize / 8
            }

            Label {
                font.pixelSize: standardSize / 3
                font.family: "System-ui"
                font.bold: true
                text: AppCore.ServerConnected ? (AppCore.IsRecording ? "ON" : "OFF") : "N/A"
                color: Styles.text
                horizontalAlignment: Qt.AlignHCenter
                Layout.alignment: Qt.AlignHCenter
            }

            VerticalSpacer {}

            GradientButton {
                color1: Styles.button_on
                borderColor: Styles.button_on_border
                Layout.alignment: Qt.AlignHCenter
                text: qsTr("Power On")
                buttonSize: rootWindow.standardSize

                onClicked: {
                    AppCore.startRecording()
                }

                Layout.bottomMargin: standardSize / 8
            }

            GradientButton {
                color1: Styles.button_off
                borderColor: Styles.button_off_border
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
            color: Styles.text
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
