cmake_minimum_required (VERSION 3.5)

set(CMAKE_CXX_STANDARD 11)

project (double_free_test)

add_compile_options(-O3)

include_directories(".")
FILE(GLOB SRC_FILES *.cpp)
add_executable(test ${SRC_FILES})
#target_link_libraries(test -Wl,--start-group grpc++_reflection grpc++ grpc gpr upb absl cares crypto protobuf ssl z re2 address_sorting pthread -Wl,--end-group)
