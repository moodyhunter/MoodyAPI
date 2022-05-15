cmake_minimum_required(VERSION 3.16)
project(MoodyPiScreen)

set(CMAKE_CXX_STANDARD 20)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_CXX_EXTENSIONS OFF)
set(CMAKE_INCLUDE_CURRENT_DIR ON)

find_package(PkgConfig REQUIRED)
pkg_check_modules(CairoPackages REQUIRED IMPORTED_TARGET fontconfig cairomm-1.0 cairomm-ft-1.0)

option(PISCREEN_DUMMY_DEVICE "Use the Dummy Device" OFF)
set(PISCREEN_OUTPUT_LIMIT 0 CACHE STRING "Limit the number of outputs, 0 means no limit")

macro(piscreen_add_datasource SOURCE)
    target_sources(PiScreen PRIVATE
        ${CMAKE_SOURCE_DIR}/datasource/${SOURCE}/source.cpp
        ${CMAKE_SOURCE_DIR}/datasource/${SOURCE}/source.hpp
    )
endmacro()

macro(piscreen_add_device DEVICE)
    target_sources(PiScreen PRIVATE
        ${CMAKE_SOURCE_DIR}/device/${DEVICE}.cpp
        ${CMAKE_SOURCE_DIR}/device/${DEVICE}.cpp
    )
endmacro()

add_executable(PiScreen
    ${CMAKE_SOURCE_DIR}/main.cpp
    ${CMAKE_SOURCE_DIR}/config/Config.hpp
    ${CMAKE_SOURCE_DIR}/config/Config.cpp

    # IPiScreen Interface
    ${CMAKE_SOURCE_DIR}/device/IPiScreenDevice.hpp

    # IDataSource Interface
    ${CMAKE_SOURCE_DIR}/datasource/IDataSource.hpp

    # Renderer
    ${CMAKE_SOURCE_DIR}/renderer/Renderer.hpp
    ${CMAKE_SOURCE_DIR}/renderer/Renderer.cpp

    # I2C Device Helper
    ${CMAKE_SOURCE_DIR}/common/I2CDevice.hpp
    ${CMAKE_SOURCE_DIR}/common/I2CDevice.cpp
)

piscreen_add_datasource(ip)
piscreen_add_datasource(cpu)
piscreen_add_datasource(mem)
piscreen_add_datasource(datetime)

# Systemd Service Inspection Support
pkg_check_modules(DBus REQUIRED IMPORTED_TARGET libsystemd)
piscreen_add_datasource(service)

piscreen_add_device(SH1106)
piscreen_add_device(Dummy)

target_compile_definitions(PiScreen PUBLIC "PISCREEN_OUTPUT_LIMIT=${PISCREEN_OUTPUT_LIMIT}")

if(PISCREEN_DUMMY_DEVICE)
    target_compile_definitions(PiScreen PRIVATE "PISCREEN_DUMMY_DEVICE=1")
endif()

target_link_libraries(PiScreen PRIVATE PkgConfig::CairoPackages PkgConfig::DBus)