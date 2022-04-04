cmake_minimum_required(VERSION 3.16)

project(MoodyApp VERSION 0.1 LANGUAGES CXX)

set(CMAKE_AUTOMOC ON)
set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_INCLUDE_CURRENT_DIR ON)

find_package(Qt6 6.2 COMPONENTS Core Gui Quick Svg RemoteObjects REQUIRED)

add_subdirectory(QtGraphicalEffects)

list(APPEND QML_DIRS "${CMAKE_CURRENT_BINARY_DIR}")
set(QML_IMPORT_PATH "${QML_DIRS}" CACHE STRING "Qt Creator 4.1 extra qml import paths")

if(ANDROID)
    add_subdirectory(3rdparty/android_openssl)
endif()

find_package(Protobuf REQUIRED)
find_package(gRPC REQUIRED)
find_program(GRPC_CPP_PLUGIN_EXECUTABLE grpc_cpp_plugin REQUIRED)



macro(generate_v2ray_grpc_source proto)
    get_filename_component(PROTO_ABS_DIR_PATH ${proto} DIRECTORY)
    file(RELATIVE_PATH _r ${PROTO_SOURCE_DIR} ${proto})
    get_filename_component(PROTO_FILE_DIR ${_r} DIRECTORY)
    get_filename_component(PROTO_FILE_NAME ${_r} NAME_WE)

    set(PROTO_HEADER_FILE "${PROTO_GENERATED_DIR}/${PROTO_FILE_DIR}/${PROTO_FILE_NAME}.pb.h")
    set(PROTO_SOURCE_FILE "${PROTO_GENERATED_DIR}/${PROTO_FILE_DIR}/${PROTO_FILE_NAME}.pb.cc")

    set(PROTO_GRPC_HEADER_FILE "${PROTO_GENERATED_DIR}/${PROTO_FILE_DIR}/${PROTO_FILE_NAME}.grpc.pb.h")
    set(PROTO_GRPC_SOURCE_FILE "${PROTO_GENERATED_DIR}/${PROTO_FILE_DIR}/${PROTO_FILE_NAME}.grpc.pb.cc")

    add_custom_command(
        OUTPUT "${PROTO_HEADER_FILE}" "${PROTO_SOURCE_FILE}" "${PROTO_GRPC_HEADER_FILE}" "${PROTO_GRPC_SOURCE_FILE}"
        COMMAND protobuf::protoc
        ARGS
            --grpc_out "${PROTO_GENERATED_DIR}"
            --cpp_out "${PROTO_GENERATED_DIR}"
            -I "${PROTO_SOURCE_DIR}"
            --plugin=protoc-gen-grpc="${GRPC_CPP_PLUGIN_EXECUTABLE}"
            "${proto}"
        DEPENDS "${proto}"
        )

    list(APPEND PROTO_HEADERS ${PROTO_HEADER_FILE} ${PROTO_GRPC_HEADER_FILE})
    list(APPEND PROTO_SOURCES ${PROTO_SOURCE_FILE} ${PROTO_GRPC_SOURCE_FILE})
endmacro()

set(PROTO_SOURCE_DIR "${CMAKE_SOURCE_DIR}/../proto")
set(PROTO_GENERATED_DIR "${CMAKE_CURRENT_BINARY_DIR}")

generate_v2ray_grpc_source("${PROTO_SOURCE_DIR}/MoodyAPI.proto")
generate_v2ray_grpc_source("${PROTO_SOURCE_DIR}/common/common.proto")
generate_v2ray_grpc_source("${PROTO_SOURCE_DIR}/privileged/privileged.proto")
generate_v2ray_grpc_source("${PROTO_SOURCE_DIR}/wg/wg.proto")

# # Adapted from https://github.com/ArkToria/ACross/blob/master/cmake/grpc.cmake
# # Generate gRPC and protobuf Sources
# set(GRPC_SOURCE "${CMAKE_CURRENT_BINARY_DIR}/MoodyAPI.grpc.pb.cc")
# set(GRPC_HEADER "${CMAKE_CURRENT_BINARY_DIR}/MoodyAPI.grpc.pb.h")
# set(PB_SOURCE "${CMAKE_CURRENT_BINARY_DIR}/MoodyAPI.pb.cc")
# set(PB_HEADER "${CMAKE_CURRENT_BINARY_DIR}/MoodyAPI.pb.h")
# set(PROTO_DIR "${CMAKE_SOURCE_DIR}/../")
# set(PROTO_FILE "${PROTO_DIR}/MoodyAPI.proto")

# add_custom_command(
#     OUTPUT "${GRPC_SOURCE}" "${GRPC_HEADER}" "${PB_SOURCE}" "${PB_HEADER}"
#     COMMAND protobuf::protoc
#     ARGS
#         --grpc_out="${CMAKE_CURRENT_BINARY_DIR}"
#         --cpp_out="${CMAKE_CURRENT_BINARY_DIR}"
#         -I="${PROTO_DIR}"
#         --plugin=protoc-gen-grpc="${GRPC_CPP_PLUGIN_EXECUTABLE}"
#         "${PROTO_FILE}"
#     DEPENDS "${PROTO_FILE}")

qt_add_executable(MoodyAppMain
    src/main.cpp
    src/AppCore.hpp
    src/AppCore.cpp
    src/AppSettings.hpp
    src/AppSettings.cpp
    src/ServerConnection.hpp
    src/ServerConnection.cpp
    src/JNICall.hpp
    src/JNICall.cpp

    "${PROTO_HEADERS}" "${PROTO_SOURCES}"
)

qt_add_qml_module(MoodyAppMain
    URI "client.api.mooody.me"
    VERSION 1.0
    QML_FILES
        qml/main.qml
        qml/GradientButton.qml
        qml/SvgButton.qml
        qml/SettingsPanel.qml
    DEPENDENCIES
        qtgraphicaleffects
)

qt6_add_repc_merged(MoodyAppMain
    IPC.rep)

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
        QT_ANDROID_VERSION_NAME "${CMAKE_PROJECT_VERSION}"
        QT_ANDROID_PACKAGE_SOURCE_DIR "${CMAKE_CURRENT_LIST_DIR}/android"
)

target_compile_definitions(MoodyAppMain
    PRIVATE
        $<$<OR:$<CONFIG:Debug>,$<CONFIG:RelWithDebInfo>>:QT_QML_DEBUG>)

target_link_libraries(MoodyAppMain
    PRIVATE
        Qt::Core
        Qt::Quick
        Qt::RemoteObjects
        Qt::Svg
        gRPC::grpc++
)

if (ANDROID)
    target_link_libraries(MoodyAppMain
        PRIVATE
            Qt::CorePrivate)
endif()
