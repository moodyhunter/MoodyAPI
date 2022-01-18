cmake_minimum_required(VERSION 3.16)

project(MoodyPiScreen)

set(CMAKE_INCLUDE_CURRENT_DIR ON)

add_executable(PiScreen
    BitBang_I2C.cpp
    BitBang_I2C.h
    # multi_oled.cpp
    # sample.cpp
    speed_test.cpp
    ss_oled.cpp
    ss_oled.h
)

# target_link_libraries(PiScreen PRIVATE pigpio)
