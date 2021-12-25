cmake_minimum_required(VERSION 3.16)

project(MoodyApp VERSION 0.1 LANGUAGES CXX)

set(CMAKE_AUTOMOC ON)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_INCLUDE_CURRENT_DIR ON)

find_package(Qt6 6.2 COMPONENTS Quick REQUIRED)

add_subdirectory(QtGraphicalEffects)

list(APPEND QML_DIRS "${CMAKE_CURRENT_BINARY_DIR}")
set(QML_IMPORT_PATH "${QML_DIRS}" CACHE STRING "Qt Creator 4.1 extra qml import paths")

if(ANDROID)
    add_subdirectory(3rdparty/android_openssl)
endif()

find_package(Protobuf REQUIRED)
find_package(gRPC REQUIRED)
find_program(GRPC_CC_PLUGIN_EXECUTABLE grpc_cpp_plugin REQUIRED)

list(APPEND PROTO_NAME_LISTS
    "MoodyApi")

# Adapted from https://github.com/ArkToria/ACross/blob/master/cmake/grpc.cmake
# Generate gRPC and protobuf Sources
foreach(PROTO_NAME ${PROTO_NAME_LISTS})
    set(GRPC_SOURCE "${CMAKE_CURRENT_BINARY_DIR}/${PROTO_NAME}.grpc.pb.cc")
    set(GRPC_HEADER "${CMAKE_CURRENT_BINARY_DIR}/${PROTO_NAME}.grpc.pb.h")
    set(PB_SOURCE "${CMAKE_CURRENT_BINARY_DIR}/${PROTO_NAME}.pb.cc")
    set(PB_HEADER "${CMAKE_CURRENT_BINARY_DIR}/${PROTO_NAME}.pb.h")
    set(PROTO_DIR "${CMAKE_SOURCE_DIR}/assets/")
    set(PROTO_FILE "${PROTO_DIR}/${PROTO_NAME}.proto")

    add_custom_command(
        OUTPUT "${GRPC_SOURCE}" "${GRPC_HEADER}" "${PB_SOURCE}" "${PB_HEADER}"
        COMMAND protobuf::protoc
        ARGS
            --grpc_out="${CMAKE_CURRENT_BINARY_DIR}"
            --cpp_out="${CMAKE_CURRENT_BINARY_DIR}"
            -I="${PROTO_DIR}"
            --plugin=protoc-gen-grpc="${GRPC_CC_PLUGIN_EXECUTABLE}"
            "${PROTO_FILE}"
        DEPENDS "${PROTO_FILE}")

    list(APPEND PROTOBUF_HEADERS "${GRPC_HEADER}" "${PB_HEADER}")
    list(APPEND PROTOBUF_SOURCES "${GRPC_SOURCE}" "${PB_SOURCE}")
endforeach()

qt_add_executable(MoodyAppMain
    src/main.cpp
    src/AppCore.hpp
    src/AppCore.cpp
    src/AppSettings.hpp
    src/AppSettings.cpp
    ${PROTOBUF_HEADERS}
    ${PROTOBUF_SOURCES}
)

set_source_files_properties(qml/Styles.qml PROPERTIES
    QT_QML_SINGLETON_TYPE TRUE
)

qt_add_qml_module(MoodyAppMain
    URI "client.api.mooody.me"
    VERSION 1.0
    QML_FILES
        qml/main.qml
        qml/GradientButton.qml
        qml/SvgButton.qml
        qml/Styles.qml
        qml/SettingsPanel.qml
    DEPENDENCIES
        qtgraphicaleffects
)

qt_add_resources(MoodyAppMain "assets"
    PREFIX "/assets"
    BASE "${CMAKE_CURRENT_LIST_DIR}/assets"
    FILES
        ${CMAKE_CURRENT_LIST_DIR}/assets/settings.svg
        ${CMAKE_CURRENT_LIST_DIR}/assets/Logo.svg
)

set_target_properties(MoodyAppMain
    PROPERTIES
        QT_ANDROID_EXTRA_LIBS "${ANDROID_EXTRA_LIBS}"
        QT_ANDROID_VERSION_CODE "${CMAKE_PROJECT_VERSION_MAJOR}"
        QT_ANDROID_PACKAGE_SOURCE_DIR "${CMAKE_CURRENT_LIST_DIR}/android"
        QT_ANDROID_VERSION_NAME "${CMAKE_PROJECT_VERSION}"
)

target_compile_definitions(MoodyAppMain
    PRIVATE
        $<$<OR:$<CONFIG:Debug>,$<CONFIG:RelWithDebInfo>>:QT_QML_DEBUG>)

target_link_libraries(MoodyAppMain
    PRIVATE
        Qt6::Quick
        gRPC::grpc++
)
