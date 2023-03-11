cmake_minimum_required(VERSION 3.20.2)
set(CMAKE_CXX_STANDARD 20)
project(cpp20_coroutine)
set(CMAKE_INSTALL_PREFIX ${CMAKE_SOURCE_DIR}/bin)

file(GLOB_RECURSE SRC_FILES src/*.cpp)
add_executable(cpp20_coroutine ${SRC_FILES})
install(TARGETS cpp20_coroutine DESTINATION .)
