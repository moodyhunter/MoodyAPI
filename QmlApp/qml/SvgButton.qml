import QtQuick
import QtQuick.Controls
import QtQuick.Layouts

import QtGraphicalEffects
import client.api.mooody.me

Control {
    id: root

    signal clicked

    property alias source: image.source
    property int imageSize: 50

    hoverEnabled: PlatformHoverEnabled

    Image {
        id: image
        anchors.fill: parent
        sourceSize.height: imageSize
        sourceSize.width: imageSize

        states: [
            State {
                name: "hovered"
                when: root.hovered
                PropertyChanges {
                    target: image
                    rotation: 100
                }
            },
            State {
                name: "hovered"
                when: mouseArea.pressed
                PropertyChanges {
                    target: image
                    rotation: 100
                }
            }
        ]

        transitions: Transition {
            NumberAnimation {
                duration: 500
                properties: "rotation"
                easing.type: Easing.OutCubic
            }
        }

        ColorOverlay {
            anchors.fill: image
            source: image
            color: AppCore.DarkMode ? "#a5a5a5" : "#2d2d2d"
        }
    }

    MouseArea {
        id: mouseArea
        anchors.fill: parent
        onClicked: root.clicked()
    }
}
