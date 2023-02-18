cmake_minimum_required(VERSION 3.13)
set(CMAKE_CXX_STANDARD 14)
project(cpptest_main)
set(CMAKE_INSTALL_PREFIX ${CMAKE_SOURCE_DIR}/bin)

enable_testing()

include(FetchContent)
FetchContent_Declare(
  googletest
  URL https://github.com/google/googletest/archive/03597a01ee50ed33e9dfd640b249b4be3799d395.zip
)
# For Windows: Prevent overriding the parent project's compiler/linker settings
set(gtest_force_shared_crt ON CACHE BOOL "" FORCE)
FetchContent_MakeAvailable(googletest)


file(GLOB_RECURSE SRC_FILES src/add.cpp src/add_test.cc)
add_executable(cpptest_main ${SRC_FILES})


target_link_libraries(cpptest_main GTest::gtest_main)

include(GoogleTest)
gtest_discover_tests(cpptest_main)

install(TARGETS cpptest_main DESTINATION .)
