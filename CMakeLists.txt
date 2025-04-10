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
${CMAKE_CURRENT_LIST_DIR}/src/http/http_client.cpp
)

set(include_dirs
${CMAKE_CURRENT_LIST_DIR}/inc
${CMAKE_CURRENT_LIST_DIR}/inc/cli
${CMAKE_CURRENT_LIST_DIR}/inc/http
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
set(Boost_USE_STATIC_LIBS ON)
set(Boost_USE_MULTITHREADED ON)
set(Boost_USE_STATIC_RUNTIME OFF)
set(BOOST_INCLUDE_LIBRARIES
    process
    system
    beast
    asio
)

FetchContent_Declare(
  Boost
  URL https://github.com/boostorg/boost/releases/download/boost-1.84.0/boost-1.84.0.tar.xz
  URL_MD5 893b5203b862eb9bbd08553e24ff146a
  DOWNLOAD_NO_EXTRACT FALSE
)

FetchContent_MakeAvailable(Boost)

# Fetch OpenSSL
FetchContent_Declare(
  OpenSSL
  GIT_REPOSITORY https://github.com/openssl/openssl.git
  GIT_TAG OpenSSL_1_1_1-stable
)
FetchContent_MakeAvailable(OpenSSL)
find_package(OpenSSL REQUIRED)

target_link_libraries(${EXECUTABLE}
                        PRIVATE
                        Boost::process
                        Boost::system
                        Boost::beast
                        Boost::asio
                        OpenSSL::SSL
                        OpenSSL::Crypto
)

target_link_libraries(${EXECUTABLE} PRIVATE ${OPENSSL_INCLUDE_DIR})

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
