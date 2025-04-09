cmake_minimum_required(VERSION 3.30.0)
set(EXECUTABLE project_name)
project(${EXECUTABLE})

# ====================================
# Default build settings
# ====================================
if(NOT CMAKE_BUILD_TYPE)
    set(CMAKE_BUILD_TYPE Debug)
endif()

set(CMAKE_CXX_STANDARD 23)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_EXPORT_COMPILE_COMMANDS ON)

# ====================================
# Sources and include directories
# ====================================
set(sources
${CMAKE_CURRENT_LIST_DIR}/src/main.cpp
${CMAKE_CURRENT_LIST_DIR}/src/command_runner.cpp

${CMAKE_CURRENT_LIST_DIR}/src/cli/cli.cpp
)

set(include_dirs
${CMAKE_CURRENT_LIST_DIR}/inc
${CMAKE_CURRENT_LIST_DIR}/inc/cli
)

# ====================================
# Creating executable target
# ====================================
add_executable(${EXECUTABLE} ${sources})
include_directories(${include_dirs})

# ====================================
# Add FetchContent to install libs
# ====================================
include(FetchContent)

# ====================================
# Boost
# ====================================
set(BOOST_ENABLE_CMAKE ON)
set(BOOST_INCLUDE_LIBRARIES process)

FetchContent_Declare(
  Boost
  URL https://github.com/boostorg/boost/releases/download/boost-1.84.0/boost-1.84.0.tar.xz
  URL_MD5 893b5203b862eb9bbd08553e24ff146a
  DOWNLOAD_NO_EXTRACT FALSE
)

FetchContent_MakeAvailable(Boost)

target_link_libraries(${EXECUTABLE} PRIVATE Boost::process)

# ====================================
# Google Test
# ====================================
FetchContent_Declare(
    googletest
    GIT_REPOSITORY https://github.com/google/googletest.git
    GIT_TAG release-1.12.1
)
FetchContent_MakeAvailable(googletest)

# ====================================
# Sources and include directories
# ====================================
set(test_sources
${CMAKE_CURRENT_LIST_DIR}/test/test_main.cpp
)
set(include_dirs
${CMAKE_CURRENT_LIST_DIR}/inc
${CMAKE_CURRENT_LIST_DIR}/googletest/include
)

# ====================================
# Creating unit test target
# ====================================
add_executable(unit_tests ${test_sources})
target_include_directories(unit_tests PRIVATE ${include_dirs})
target_link_libraries(unit_tests gtest_main pthread)

# ====================================
# Enable testing and add unit tests to ctest
# ====================================
enable_testing()
add_test(NAME all_tests COMMAND unit_tests)
