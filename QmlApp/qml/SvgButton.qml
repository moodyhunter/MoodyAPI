import QtQuick
import QtQuick.Controls
import QtQuick.Layouts

import QtGraphicalEffects

Control {
    id: root

    signal clicked

    property alias source: image.source
    property int imageSize: 50

    Image {
        id: image
        anchors.fill: parent
        sourceSize.height: imageSize
        sourceSize.width: imageSize

        states: State {
            name: "hovered"
            when: root.hovered
            PropertyChanges {
                target: image
                rotation: 100
            }
        }

        transitions: Transition {
            NumberAnimation {
                duration: 500
                properties: "rotation"
                easing.type: Easing.OutCubic
            }
        }
    }

    MouseArea {
        anchors.fill: parent
        onClicked: root.clicked()
    }
}
