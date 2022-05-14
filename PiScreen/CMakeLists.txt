cmake_minimum_required(VERSION 3.16)

project(MoodyPiScreen)

set(CMAKE_CXX_STANDARD 20)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_CXX_EXTENSIONS OFF)

set(CMAKE_INCLUDE_CURRENT_DIR ON)

find_package(PkgConfig REQUIRED)
pkg_check_modules(CairoPackages REQUIRED IMPORTED_TARGET fontconfig cairomm-1.0 cairomm-ft-1.0)

option(PISCREEN_DUMMY_DEVICE "Enable the PiScreen dummy device" OFF)

add_executable(PiScreen
    ${CMAKE_SOURCE_DIR}/main.cpp
    ${CMAKE_SOURCE_DIR}/common/IPiScreenDevice.hpp
)

if(PISCREEN_DUMMY_DEVICE)
    target_compile_definitions(PiScreen PRIVATE "PISCREEN_DUMMY_DEVICE=1")
    target_sources(PiScreen PRIVATE
        ${CMAKE_SOURCE_DIR}/device/Dummy.cpp
        ${CMAKE_SOURCE_DIR}/device/Dummy.hpp
    )
else()
    target_sources(PiScreen PRIVATE
        ${CMAKE_SOURCE_DIR}/device/SH1106.cpp
        ${CMAKE_SOURCE_DIR}/device/SH1106.hpp

        ${CMAKE_SOURCE_DIR}/common/I2CDevice.hpp
        ${CMAKE_SOURCE_DIR}/common/I2CDevice.cpp
    )
endif()

target_link_libraries(PiScreen PRIVATE PkgConfig::CairoPackages)