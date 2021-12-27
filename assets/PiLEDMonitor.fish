#!/bin/fish

while true
    if [ (systemctl show motion | grep 'SubState=') = "SubState=running" ]
        echo "1" > /sys/class/leds/led?/brightness
    end
    sleep 2
    echo "0" > /sys/class/leds/led?/brightness
    sleep 2
end
