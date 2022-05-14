cmake_minimum_required(VERSION 3.16)

project(MoodyPiScreen)

set(CMAKE_CXX_STANDARD 20)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_CXX_EXTENSIONS OFF)

set(CMAKE_INCLUDE_CURRENT_DIR ON)

find_package(PkgConfig REQUIRED)
pkg_check_modules(CairoPackages REQUIRED IMPORTED_TARGET fontconfig cairo cairo-fc)

add_executable(PiScreen
    I2CDevice.cpp
    I2CDevice.hpp
    main.cpp
    SH1106dev.cpp
    SH1106dev.hpp
)

target_link_libraries(PiScreen PRIVATE PkgConfig::CairoPackages)