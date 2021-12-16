cmake_minimum_required(VERSION 3.16)

project(MoodyApp VERSION 0.1 LANGUAGES CXX)

set(CMAKE_AUTOMOC ON)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

find_package(Qt6 6.2 COMPONENTS Quick REQUIRED)

add_subdirectory(QtGraphicalEffects)

list(APPEND QML_DIRS "${CMAKE_CURRENT_BINARY_DIR}")
set(QML_IMPORT_PATH "${QML_DIRS}" CACHE STRING "Qt Creator 4.1 extra qml import paths")

if(ANDROID)
    add_subdirectory(3rdparty/android_openssl)
endif()

qt_add_executable(MoodyAppMain
    src/main.cpp
)

set_target_properties(MoodyAppMain
    PROPERTIES
        QT_ANDROID_EXTRA_LIBS "${ANDROID_EXTRA_LIBS}"
        QT_ANDROID_VERSION_CODE "${CMAKE_PROJECT_VERSION_MAJOR}"
        QT_ANDROID_PACKAGE_SOURCE_DIR "${CMAKE_CURRENT_LIST_DIR}/android"
        QT_ANDROID_VERSION_NAME "${CMAKE_PROJECT_VERSION}"
    )

qt_add_qml_module(MoodyAppMain
    URI client.api.mooody.me
    VERSION 1.0
    QML_FILES qml/main.qml
    DEPENDENCIES qtgraphicaleffects
)

target_compile_definitions(MoodyAppMain
    PRIVATE
        $<$<OR:$<CONFIG:Debug>,$<CONFIG:RelWithDebInfo>>:QT_QML_DEBUG>)

target_link_libraries(MoodyAppMain
    PRIVATE
        Qt6::Quick
    )
