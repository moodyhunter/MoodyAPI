cmake_minimum_required(VERSION 3.16)

project(MoodyPiScreen)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_CXX_EXTENSIONS OFF)

set(CMAKE_INCLUDE_CURRENT_DIR ON)

add_executable(PiScreen
    I2CDevice.cpp
    I2CDevice.h
    # multi_oled.cpp
    # sample.cpp
    speed_test.cpp
    ss_oled.cpp
    ss_oled.h
)

# target_link_libraries(PiScreen PRIVATE pigpio)
