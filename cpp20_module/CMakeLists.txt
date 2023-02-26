cmake_minimum_required(VERSION 3.20)

set(CMAKE_CXX_STANDARD 20)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_CXX_EXTENSIONS OFF)


add_compile_options(-fmodules-ts)

project(cpp20_module)
set(CMAKE_INSTALL_PREFIX ${CMAKE_SOURCE_DIR}/bin)

file(GLOB_RECURSE SRC_FILES src/*.cpp)
add_executable(${PROJECT_NAME} ${SRC_FILES})
install(TARGETS ${PROJECT_NAME} DESTINATION .)
