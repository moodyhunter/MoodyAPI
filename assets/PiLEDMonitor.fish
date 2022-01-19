#!/bin/fish

set oldState "off"
./PiScreen 75 2 "Camera:             Off"

while true
    if [ (systemctl show motion | grep 'SubState=') = "SubState=running" ]
        echo "1" | tee /sys/class/leds/led1/brightness
        if [ $oldState = "off" ]
            ./PiScreen 75 2 "Camera:             On"
            set oldState "on"
        end
    else
        if [ $oldState = "on" ]
            ./PiScreen 75 2 "Camera:             Off"
            set oldState "off"
    end
    end
    sleep 2
    echo "0" | tee /sys/class/leds/led1/brightness
    sleep 2
end
