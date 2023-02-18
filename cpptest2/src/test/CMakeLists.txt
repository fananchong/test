include(FetchContent)
FetchContent_Declare(
  googletest
  URL https://github.com/google/googletest/archive/03597a01ee50ed33e9dfd640b249b4be3799d395.zip
)
set(gtest_force_shared_crt ON CACHE BOOL "" FORCE)
FetchContent_MakeAvailable(googletest)

set(TEST_BINARY ${PROJECT_NAME}_test)
file(GLOB_RECURSE TEST_FILES ${CMAKE_SOURCE_DIR}/src/*.cc)
file(GLOB_RECURSE APP_FILES ${CMAKE_SOURCE_DIR}/src/*.cpp)
list (REMOVE_ITEM APP_FILES ${CMAKE_SOURCE_DIR}/src/main.cpp)
add_executable(${TEST_BINARY} ${TEST_FILES} ${APP_FILES})
target_link_libraries(${TEST_BINARY} GTest::gtest_main)
include(GoogleTest)
gtest_discover_tests(${TEST_BINARY})
