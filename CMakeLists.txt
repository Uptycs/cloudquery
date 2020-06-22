
cmake_minimum_required(VERSION 3.10)
string(ASCII 27 ESC)
macro(WARNING_LOG MESSAGE)
  message("-- ${ESC}[31m${MESSAGE}${ESC}[m")
endmacro(WARNING_LOG)

macro(LOG MESSAGE)
  message("-- ${MESSAGE}")
endmacro(LOG)

# Set expected directory for within-build dependencies
if(DEFINED ENV{CLOUDQUERY_DEPS})
  set(BUILD_DEPS "$ENV{CLOUDQUERY_DEPS}")
else()
  set(BUILD_DEPS "/usr/local/cloudquery")
endif()

if(WIN32)
  set(WINDOWS_DEP_DIR "$ENV{ChocolateyInstall}/lib")
  string(REGEX REPLACE "\\\\" "/" WINDOWS_DEP_DIR ${WINDOWS_DEP_DIR})
  list(APPEND WINDOWS_DEP_LIST ${WINDOWS_DEP_DIR})

  list(APPEND CMAKE_PROGRAM_PATH "${WINDOWS_DEP_DIR}/*/tools")
endif()

if(WIN32)
  add_definitions(
    -DCLOUDQUERY_WINDOWS
  )
else()
  add_definitions(
    -DCLOUDQUERY_POSIX
  )
endif()

# If no explicit compiler override and clang++ exists, prefer clang.
find_program(CLANGXX_EXECUTABLE "clang++" ENV PATH)
if(DEFINED ENV{CC})
  set(CMAKE_C_COMPILER "$ENV{CC}")
  LOG("Overriding C compiler from clang to $ENV{CC}")
elseif(CLANGXX_EXECUTABLE)
  set(CMAKE_C_COMPILER "clang")
endif()
if(DEFINED ENV{CXX})
  set(CMAKE_CXX_COMPILER "$ENV{CXX}")
  LOG("Overriding CXX compiler from clang++ to $ENV{CXX}")
elseif(CLANGXX_EXECUTABLE)
  set(CMAKE_CXX_COMPILER "clang++")
endif()

if(CMAKE_CXX_COMPILER MATCHES "clang")
  set(CLANG TRUE)
else()
  set(CLANG FALSE)
endif()

# We will add a lot of extra flags and options when using the cloudquery runtime.
# The "DEPS" environment is the cloudquery runtime compiler and 3rd-party deps.
if(NOT DEFINED ENV{SKIP_DEPS})
  set(DEPS TRUE)
else()
  set(DEPS FALSE)
endif()

if(NOT DEFINED ENV{SKIP_CCACHE})
  find_program(CCACHE_EXECUTABLE ccache ENV PATH)
  if(CCACHE_EXECUTABLE)
    if(NOT DEFINED CCACHE_FOUND_LOG)
      set(CCACHE_FOUND_LOG TRUE CACHE BOOL "Write log line for ccache.")
      LOG("Found ccache ${CCACHE_EXECUTABLE}")
      LOG("Using ccache to speed up compilation")
    endif()
    set(ENV{CCACHE_CPP2} "yes")
    set_property(GLOBAL PROPERTY RULE_LAUNCH_COMPILE "${CCACHE_EXECUTABLE}")
    set_property(GLOBAL PROPERTY RULE_LAUNCH_LINK "${CCACHE_EXECUTABLE}")
  endif()
endif()

# cloudquery additional compiler flags added by CMake.
set(CXX_STD "c++14")

find_program(PYTHON_EXECUTABLE "python" ${BUILD_DEPS} ENV PATH)

# Use cloudquery language to set platform/os
if(DEFINED ENV{CLOUDQUERY_PLATFORM})
  set(PLATFORM "$ENV{CLOUDQUERY_PLATFORM}")
else()
  execute_process(
    COMMAND "${PYTHON_EXECUTABLE}" "${CMAKE_SOURCE_DIR}/tools/get_platform.py"
    WORKING_DIRECTORY "${CMAKE_SOURCE_DIR}"
    OUTPUT_VARIABLE PLATFORM
    OUTPUT_STRIP_TRAILING_WHITESPACE
  )
endif()

if(PLATFORM STREQUAL "")
  WARNING_LOG("Unable to detect cloudquery platform: ./tools/get_platform.py failed")
  message(FATAL_ERROR "Cannot proceed without knowing the build platform")
endif()
list(GET PLATFORM 0 CLOUDQUERY_BUILD_PLATFORM)
list(GET PLATFORM 1 CLOUDQUERY_BUILD_DISTRO)
string(REPLACE "." "_" PLATFORM "${PLATFORM}")
string(TOUPPER "${PLATFORM}" PLATFORM)
list(GET PLATFORM 0 CLOUDQUERY_BUILD_PLATFORM_DEFINE)
list(GET PLATFORM 1 CLOUDQUERY_BUILD_DISTRO_DEFINE)

if(DEFINED ENV{FAST})
  set(ENV{CLOUDQUERY_BUILD_SHARED} TRUE)
endif()

# Set non-C compile flags and whole-loading linker flags.
# cloudquery needs ALL symbols in the libraries it includes for relaxed ctors
# late-loading modules and SQLite introspection utilities.
if(APPLE)
  if(DEFINED ENV{OSX_VERSION_MIN})
    set(OSX_VERSION_MIN "$ENV{OSX_VERSION_MIN}")
  elseif(DEFINED ENV{OSX_VERSION_NATIVE})
    set(OSX_VERSION_MIN "${CLOUDQUERY_BUILD_DISTRO}")
  else()
    set(OSX_VERSION_MIN "10.11")
  endif()
  add_compile_options(
    -mmacosx-version-min=${OSX_VERSION_MIN}
  )

  set(OS_WHOLELINK_PRE "-Wl,-force_load")
  set(OS_WHOLELINK_POST "")
  # Special compile flags for Objective-C++
  set(OBJCXX_COMPILE_FLAGS
    "-x objective-c++ -fobjc-arc -Wno-c++11-extensions")
  set(POSIX TRUE)
elseif(WIN32)
  # Windows
  # C++11 features should be on by default in 2012 and beyond.
  set(OS_WHOLELINK_PRE "-WHOLEARCHIVE:")
  set(OS_WHOLELINK_POST "")
  set(WINDOWS TRUE)
else()
  set(OS_WHOLELINK_PRE "-Wl,-whole-archive")
  set(OS_WHOLELINK_POST "-Wl,-no-whole-archive")
  # Set CMAKE variables depending on platform, to know which tables and what
  # component-specific globbing is needed.
  if(${CLOUDQUERY_BUILD_PLATFORM} STREQUAL "freebsd")
    set(FREEBSD TRUE)
    set(DEPS FALSE)
    include_directories("/usr/local/include")
    include_directories("/usr/include")
  else()
    set(LINUX TRUE)
  endif()
  set(POSIX TRUE)
endif()

# Configure the CXX standard and optionally set the stdlib to LLVM's.
if(NOT WIN32)
  if(CLANG)
    set(CXX_STDLIB "-stdlib=libc++")
  else()
    set(CXX_STDLIB "")
  endif()

  # Only allow adding CXX flags if using SKIP_DEPS.
  if(DEFINED ENV{SKIP_DEPS})
    set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -std=${CXX_STD} ${CXX_STDLIB}" CACHE STRING "" FORCE)
  else()
    set(CMAKE_CXX_FLAGS "-std=${CXX_STD} ${CXX_STDLIB}" CACHE STRING "" FORCE)
  endif()
endif()

if(POSIX)
  if(CLANG)
    add_compile_options(
      -Qunused-arguments
      -Wabi-tag
      -Wno-unused-local-typedef
      -Wno-deprecated-register
      -Wno-unknown-warning-option
      -Wstrict-aliasing
      -Wno-missing-field-initializers
      -Wnon-virtual-dtor
      -Wchar-subscripts
      -Wpointer-arith
      -Woverloaded-virtual
    )
  else()
    add_compile_options(
      -Wno-unknown-pragmas
    )
  endif()
  add_compile_options(
    -Wformat
    -Wformat-security
    -Werror=format-security
    -fpermissive
    -fstack-protector-all
    -pipe
    -fdata-sections
    -ffunction-sections
  )
  if(NOT FREEBSD AND NOT DEFINED ENV{CLOUDQUERY_BUILD_SHARED})
    add_compile_options(
      -fvisibility=hidden
      -fvisibility-inlines-hidden
    )
  endif()
elseif(WINDOWS)
  add_compile_options(
    /EHa
    /W3
    /guard:cf
  )
  add_definitions(
    -DWIN32_LEAN_AND_MEAN
    -DBOOST_ALL_NO_LIB
    -DAWS_DISABLE_DEPRECATION
  )
endif()

set(TEST_ARGS "")
if(LINUX AND CLANG AND DEPS)
  # Use the LLVM versions of ar and ranlib to support LTO builds
  find_program(LLVMRANLIB_EXECUTABLE "llvm-ranlib" ENV PATH)
  find_program(LLVMAR_EXECUTABLE "llvm-ar" ENV PATH)
  set(CMAKE_AR "${LLVMAR_EXECUTABLE}")
  set(CMAKE_RANLIB "${LLVMRANLIB_EXECUTABLE}")
  set(CMAKE_C_ARCHIVE_CREATE "${LLVMAR_EXECUTABLE} qc <TARGET> <LINK_FLAGS> <OBJECTS>")
  set(CMAKE_CXX_ARCHIVE_CREATE "${LLVMAR_EXECUTABLE} qc <TARGET> <LINK_FLAGS> <OBJECTS>")
  set(CMAKE_C_ARCHIVE_FINISH "${LLVMRANLIB_EXECUTABLE} <TARGET>")
  set(CMAKE_CXX_ARCHIVE_FINISH "${CMAKE_C_ARCHIVE_FINISH}")
endif()

# make debug (environment variable from Makefile)
if(DEFINED ENV{DEBUG})
  set(CMAKE_BUILD_TYPE "DEBUG")
  WARNING_LOG("Setting DEBUG build")
else()
  set(CMAKE_BUILD_TYPE "RELEASE")
endif()

if(CMAKE_BUILD_TYPE STREQUAL "DEBUG")
  add_definitions(-DDEBUG)
else()
  add_definitions(-DNDEBUG)
endif()

# cloudquery cmake is using CMake generator expressions to support multiconfig IDEs
# https://cmake.org/cmake/help/v3.10/manual/cmake-generator-expressions.7.html

if(WINDOWS)
  # TODO(#1985): Is there a MSVC flag congruent to -fstandalone-debug?
  add_compile_options(
    "$<$<CONFIG:DEBUG>:/Zi>"
    "$<$<CONFIG:DEBUG>:/Od>"
    "$<$<CONFIG:RELEASE>:/Ot>"
    "$<$<CONFIG:RELEASE>:/WX>"
  )
else()
  add_compile_options(
    "$<$<AND:$<OR:$<CXX_COMPILER_ID:Clang>,$<CXX_COMPILER_ID:AppleClang>>,$<CONFIG:RELEASE>>:-Oz>"
    "$<$<AND:$<NOT:$<OR:$<CXX_COMPILER_ID:Clang>,$<CXX_COMPILER_ID:AppleClang>>>,$<CONFIG:RELEASE>>:-Os>"
  )
endif()

add_compile_options(
  "$<$<CONFIG:DEBUG>:-g>"
  "$<$<CONFIG:DEBUG>:-O0>"
  "$<$<CONFIG:DEBUG>:-fstandalone-debug>"
)

# Do not enable fortify with clang: http://llvm.org/bugs/show_bug.cgi?id=16821
# Visual Studio does not support generation of PIC code
if (POSIX)
  # TSAN and MSAN do not like PIE/PIC on Linux 4.1.
  # To be safe, only include them when building 'release' outputs.
  add_compile_options(
    "$<$<CONFIG:RELEASE>:-g>"
    "$<$<CONFIG:RELEASE>:-fno-limit-debug-info>"
    "$<$<CONFIG:RELEASE>:-fPIE>"
    "$<$<CONFIG:RELEASE>:-fpie>"
    "$<$<CONFIG:RELEASE>:-fPIC>"
    "$<$<CONFIG:RELEASE>:-fpic>"
  )
  if (CLANG AND DEPS AND NOT FREEBSD)
    # Only enable LTO builds when using clang on Unix, for now
    add_compile_options("$<$<CONFIG:RELEASE>:-flto=thin>")
  endif()
endif()

# Generate a compile_commands.json for static analyzer input.
set(CMAKE_EXPORT_COMPILE_COMMANDS ON)

if(NOT DEFINED ENV{OPTIMIZED})
  if(WINDOWS)
    # It is unclear if AVX must be managed on Windows.
    add_compile_options(/favor:AMD64)
  else()
    add_compile_options(-march=x86-64 -mno-avx)
  endif()
endif()

# make analyze (environment variable from Makefile)
if(DEFINED ENV{ANALYZE})
  # TODO(#1985): clang-analyze on Windows
  if(WINDOWS)
    message(FATAL_ERROR "Cannot analyze on Windows right now." )
  else()
    set(CMAKE_CXX_COMPILER "${CMAKE_SOURCE_DIR}/tools/analysis/clang-analyze.sh")
  endif()
endif()

# make packages will set release to true and blacklist development features,
# development plugins, etc.
if(DEFINED ENV{PACKAGE})
  set(CLOUDQUERY_BUILD_RELEASE TRUE)
else()
  set(CLOUDQUERY_BUILD_RELEASE FALSE)
endif()

# Similar to the feature skips, there are default-build skips.
if(DEFINED ENV{SKIP_TESTS})
  set(SKIP_TESTS TRUE)
endif()
if(DEFINED ENV{SKIP_BENCHMARKS})
  set(SKIP_BENCHMARKS TRUE)
endif()
if(DEFINED ENV{SKIP_TABLES})
  set(SKIP_TABLES TRUE)
endif()

# Finished setting compiler/compiler flags.
project(CLOUDQUERY)
set(CMAKE_SKIP_INSTALL_ALL_DEPENDENCY TRUE)
set(CMAKE_MODULE_PATH "${CMAKE_CURRENT_SOURCE_DIR}/CMake" "${CMAKE_MODULE_PATH}")
include(CMakeLibs)

find_program(CLANG_TIDY NAMES "clang-tidy" ENV PATH)
if(DEFINED ENV{TIDY} AND CLANG_TIDY)
  message(STATUS "clang-tidy found: ${CLANG_TIDY}")
  set(TIDY_CHECKS
    "*"
    "-cert-err58-cpp"
    "-cppcoreguidelines-pro-bounds-array-to-pointer-decay"
    "-readability-else-after-return"
    "-cppcoreguidelines-pro-type-union-access"
    "-modernize-use-bool-literals"
    "-modernize-use-emplace"
    "-google-runtime-references"
    # This adds header guards when we already use pragmas
    "-llvm-header-guard"
    # This adds an explicit for operator TYPE() in flags.h
    "-google-explicit-constructor"
    # Too many notes about not defining all specials
    "-cppcoreguidelines-special-member-functions"
    # These default arguments are very valuable
    "-google-default-arguments"
  )
  JOIN("${TIDY_CHECKS}" "," TIDY_CHECKS)
  set(DO_CLANG_TIDY "${CLANG_TIDY}"
    "-checks=${TIDY_CHECKS}"
    "-p=${CMAKE_BINARY_DIR}"
    "-fix"
  )
endif()

if(WINDOWS)
  foreach(flag_var
          CMAKE_CXX_FLAGS CMAKE_CXX_FLAGS_DEBUG CMAKE_CXX_FLAGS_RELEASE
          CMAKE_C_FLAGS CMAKE_C_FLAGS_DEBUG CMAKE_C_FLAGS_RELEASE
          CMAKE_CXX_FLAGS_MINSIZEREL CMAKE_CXX_FLAGS_RELWITHDEBINFO
          CMAKE_C_FLAGS_MINSIZEREL CMAKE_C_FLAGS_RELWITHDEBINFO)
    if(${flag_var} MATCHES "/MD")
      string(REGEX REPLACE "/MD" "/MT" ${flag_var} "${${flag_var}}")
    endif(${flag_var} MATCHES "/MD")

    # Ignore unused parameter messages for C-based projects (e.g. sqlite3)
    if (flag_var MATCHES "^CMAKE_C_FLAGS*")
      set(${flag_var} "${${flag_var}} /wd4101")
    endif(flag_var MATCHES "^CMAKE_C_FLAGS*")
  endforeach(flag_var)

  set(CMAKE_EXE_LINKER_FLAGS "${CMAKE_EXE_LINKER_FLAGS} /LTCG")
elseif(CLANG AND DEPS AND NOT FREEBSD)
  # Clang on Unix uses LTO; we also need to pass -flto when linking
  set(CMAKE_EXE_LINKER_FLAGS "${CMAKE_EXE_LINKER_FLAGS} -flto=thin")
  set(CMAKE_SHARED_LINKER_FLAGS "${CMAKE_SHARED_LINKER_FLAGS} -flto=thin")
endif()

# Make sure deps were built before compiling (else show warning).
execute_process(
  COMMAND "${CMAKE_SOURCE_DIR}/tools/provision.sh" check "${CMAKE_BINARY_DIR}"
  WORKING_DIRECTORY "${CMAKE_SOURCE_DIR}"
  OUTPUT_VARIABLE CLOUDQUERY_DEPS_MESSAGE
  RESULT_VARIABLE CLOUDQUERY_DEPS_CHECK
  OUTPUT_STRIP_TRAILING_WHITESPACE
)
if(CLOUDQUERY_DEPS_CHECK)
  WARNING_LOG("${CLOUDQUERY_DEPS_MESSAGE}")
endif()

# Discover build version from an environment variable or from the git checkout.
if(DEFINED ENV{CLOUDQUERY_BUILD_VERSION})
  set(CLOUDQUERY_BUILD_VERSION "$ENV{CLOUDQUERY_BUILD_VERSION}")
else()
  # Generate version from git
  execute_process(
    COMMAND git describe --tags HEAD --abbrev=0 --always
    WORKING_DIRECTORY "${CMAKE_SOURCE_DIR}"
    OUTPUT_VARIABLE CLOUDQUERY_BUILD_VERSION
    OUTPUT_STRIP_TRAILING_WHITESPACE
  )
endif()

# Discover the SDK version from an environment variable or the build version.
if(CLOUDQUERY_BUILD_VERSION STREQUAL "")
  WARNING_LOG("Unable to detect cloudquery version: 'git describe' failed")
  WARNING_LOG("Setting the version to: 1.0.0-unknown (beware!)")
  set(CLOUDQUERY_BUILD_VERSION "1.0.0-unknown")
endif()
set(OSQUERY_BUILD_SDK_VERSION "4.2.0.4")

if(DEFINED ENV{CQ_SDK_VERSION})
  set(CLOUDQUERY_BUILD_SDK_VERSION "$ENV{CQ_SDK_VERSION}")
else()
  string(REPLACE "-" ";" CLOUDQUERY_BUILD_SDK_VERSION ${CLOUDQUERY_BUILD_VERSION})
  list(GET CLOUDQUERY_BUILD_SDK_VERSION 0 CLOUDQUERY_BUILD_SDK_VERSION)
endif()

# Set various platform/platform-version/build version/etc defines.
add_definitions(
  -DOSQUERY_BUILD_SDK_VERSION=${OSQUERY_BUILD_SDK_VERSION}
  -DCLOUDQUERY_BUILD_SDK_VERSION=${CLOUDQUERY_BUILD_SDK_VERSION}
  -DCLOUDQUERY_BUILD_VERSION=${CLOUDQUERY_BUILD_VERSION}
  -DCLOUDQUERY_BUILD_PLATFORM=${CLOUDQUERY_BUILD_PLATFORM}
  -DOSQUERY_BUILD_PLATFORM=${CLOUDQUERY_BUILD_PLATFORM}
  -DCLOUDQUERY_BUILD_DISTRO=${CLOUDQUERY_BUILD_DISTRO}
  -DOSQUERY_BUILD_DISTRO=${CLOUDQUERY_BUILD_DISTRO}
  -D${CLOUDQUERY_BUILD_PLATFORM_DEFINE}
  -D${CLOUDQUERY_BUILD_PLATFORM_DEFINE}_${CLOUDQUERY_BUILD_DISTRO_DEFINE}
  -DSTRIP_FLAG_HELP=1
  -DBOOST_NETWORK_ENABLE_HTTPS
  -DNOMINMAX
  -DBOOST_CONFIG_SUPPRESS_OUTDATED_MESSAGE
  -DRAPIDJSON_HAS_STDSTRING=1
  -DRAPIDJSON_NO_SIZETYPEDEFINE
  # See boost 1.65 breakage https://github.com/Homebrew/homebrew-core/pull/17150
  -D__ASSERT_MACROS_DEFINE_VERSIONS_WITHOUT_UNDERSCORES=0
)

if(APPLE)
  LOG_PLATFORM("OS X")
elseif(CLOUDQUERY_BUILD_PLATFORM STREQUAL "debian")
  LOG_PLATFORM("Debian")
elseif(CLOUDQUERY_BUILD_PLATFORM STREQUAL "ubuntu")
  LOG_PLATFORM("Ubuntu")
elseif(CLOUDQUERY_BUILD_PLATFORM STREQUAL "centos")
  LOG_PLATFORM("CentOS")
elseif(CLOUDQUERY_BUILD_PLATFORM STREQUAL "scientific")
  LOG_PLATFORM("Scientific Linux")
elseif(CLOUDQUERY_BUILD_PLATFORM STREQUAL "rhel")
  LOG_PLATFORM("RHEL")
elseif(CLOUDQUERY_BUILD_PLATFORM STREQUAL "amazon")
  LOG_PLATFORM("Amazon Linux")
elseif(CLOUDQUERY_BUILD_PLATFORM STREQUAL "oracle")
  LOG_PLATFORM("Oracle Linux")
elseif(CLOUDQUERY_BUILD_PLATFORM STREQUAL "gentoo")
  LOG_PLATFORM("Gentoo Linux")
elseif(CLOUDQUERY_BUILD_PLATFORM STREQUAL "funtoo")
  LOG_PLATFORM("Funtoo Linux")
elseif(CLOUDQUERY_BUILD_PLATFORM STREQUAL "freebsd")
  LOG_PLATFORM("FreeBSD")
elseif(CLOUDQUERY_BUILD_PLATFORM STREQUAL "windows")
  LOG_PLATFORM("Windows")
endif()

# Add sysroot overrides for each platform/distro.
# These should be prioritized over all other includes.
if(APPLE)
  include_directories(SYSTEM "${CMAKE_SOURCE_DIR}/third-party/osquery/third-party/sysroots/darwin")
elseif(LINUX)
  include_directories(SYSTEM "${CMAKE_SOURCE_DIR}/third-party/osquery/third-party/sysroots/linux")
endif()

if(POSIX AND DEPS)
  set(CMAKE_REQUIRED_INCLUDES
    "${BUILD_DEPS}/legacy/include"
    "${BUILD_DEPS}/include"
  )

  # The order here matters, please do not reorder.
  if(LINUX AND CLANG)
    include_directories(SYSTEM "${BUILD_DEPS}/include/c++/v1")
    include_directories(SYSTEM "${BUILD_DEPS}/legacy/include")
  endif()
  include_directories(SYSTEM "${BUILD_DEPS}/include")
  if(LINUX AND CLANG)
    include_directories(SYSTEM "${BUILD_DEPS}/lib/clang/4.0.0/include")
  endif()
endif()

include_directories("${CMAKE_SOURCE_DIR}/third-party/osquery/third-party/sqlite3")
include_directories("${CMAKE_SOURCE_DIR}/third-party/osquery/include")
include_directories("${CMAKE_SOURCE_DIR}/third-party/osquery/")
include_directories("${CMAKE_SOURCE_DIR}/third-party/curl/include")
include_directories("${CMAKE_SOURCE_DIR}/include")
include_directories("${CMAKE_SOURCE_DIR}")

if(NOT FREEBSD)
  include_directories("${CMAKE_SOURCE_DIR}/third-party/osquery/third-party/linenoise-ng")
endif()

set(MKDIR_OPTS "")
if(WINDOWS)
  set(PROCESS_FAMILY "windows")

  # Win64 specific gflags
  set(BUILD_SHARED_LIBS OFF CACHE BOOL "Build static gtest/gmock library" FORCE)

  file(GLOB INCLUDE_DIRS LIST_DIRECTORIES true "${WINDOWS_DEP_DIR}/*/*/include")
  include_directories(SYSTEM ${INCLUDE_DIRS})
  include_directories(SYSTEM "${WINDOWS_DEP_DIR}/boost-msvc14/local")
  include_directories(SYSTEM "${WINDOWS_DEP_DIR}/yara/local/include")

  ## TODO(FIXME): temporary for now
  include_directories("${WINDOWS_DEP_DIR}/linenoise-ng/local/include/linenoise")
else()
  set(PROCESS_FAMILY "posix")
  set(MKDIR_OPTS "-p")
endif()

if(WINDOWS)
  if(NOT OS_LIB_DIRS)
    set(OS_LIB_DIRS)
  endif()
  foreach(DEP_DIR ${WINDOWS_DEP_LIST})
    file(GLOB_RECURSE WINDOWS_LOCAL_LIB_DIRS LIST_DIRECTORIES false "${DEP_DIR}/*.lib")
    foreach(LIB_FILE ${WINDOWS_LOCAL_LIB_DIRS})
      get_filename_component(LIB_DIR ${LIB_FILE} DIRECTORY)
      list(APPEND OS_LIB_DIRS ${LIB_DIR})
    endforeach()
  endforeach()
  list(REMOVE_DUPLICATES OS_LIB_DIRS)
  link_directories(${OS_LIB_DIRS})
else()
  set(OS_LIB_DIRS
    "/lib"
    "/lib64"
    "/usr/lib"
    "/usr/lib64"
    "/usr/lib/x86_64-linux-gnu/"
    )
endif()

if(WINDOWS)
  list(APPEND OS_LIB_DIRS "${CMAKE_SOURCE_DIR}/third-party/osquery/build/windows10/osquery/RelWithDebInfo")
elseif(LINUX)
  list(APPEND OS_LIB_DIRS "${CMAKE_SOURCE_DIR}/third-party/osquery/build/linux/osquery")
elseif(APPLE)
  list(APPEND OS_LIB_DIRS "${CMAKE_SOURCE_DIR}/third-party/osquery/build/darwin/osquery")
endif()
list(APPEND OS_LIB_DIRS "${CMAKE_BINARY_DIR}/cloudquery")
#list(APPEND OS_LIB_DIRS "${CMAKE_BINARY_DIR}/third-party/curl/lib/.libs")

link_directories(${OS_LIB_DIRS})


# Make sure the generated paths exist
execute_process(COMMAND mkdir ${MKDIR_OPTS} "${CMAKE_BINARY_DIR}/generated")

# We need to link some packages as dynamic/dependent.
set(CMAKE_FIND_LIBRARY_SUFFIXES .lib .dylib .so .a)
set(CMAKE_SKIP_RPATH TRUE)

# Previously, this would find_package for Gtest.
# We can supply the known include paths within each test as compile flags.
enable_testing()
set(CTEST_TEST_TIMEOUT 600)
set(BUILD_GTEST TRUE)

# Detect TR1 usage.
include(CheckIncludeFileCXX)
if(WINDOWS)
  check_include_file_cxx("tuple" HAVE_TR1_TUPLE)
else()
  check_include_file_cxx("tr1/tuple" HAVE_TR1_TUPLE)
endif()

if(FREEBSD OR "${HAVE_TR1_TUPLE}" STREQUAL "")
  add_definitions(-DGTEST_USE_OWN_TR1_TUPLE=1)
 else()
  add_definitions(-DGTEST_USE_OWN_TR1_TUPLE=0)
endif()

set(GTEST_FLAGS
  "-I${CMAKE_SOURCE_DIR}/third-party/osquery/third-party/googletest/googletest/include"
  "-I${CMAKE_SOURCE_DIR}/third-party/osquery/third-party/googletest/googlemock/include"
)
join("${GTEST_FLAGS}" " " GTEST_FLAGS)

if(NOT SKIP_TESTS)
  add_subdirectory("${CMAKE_SOURCE_DIR}/third-party/osquery/third-party/googletest")
endif()

add_subdirectory(cloudquery)

# Packages must be included after all other CMake includes.
include(Packages)

# make format-all
find_program(CLANGFORMAT_EXECUTABLE "clang-format" ENV PATH)
add_custom_target(
  format_all
  find cloudquery include tools \( -name "*.h" -o -name "*.cpp" -o -name "*.mm" \)
    -exec "${CLANGFORMAT_EXECUTABLE}" -i {} +
  WORKING_DIRECTORY "${CMAKE_SOURCE_DIR}"
  COMMENT "Formatting all cloudquery code with clang-format" VERBATIM
)

# make format
add_custom_target(
  format
  "${PYTHON_EXECUTABLE}" "${CMAKE_SOURCE_DIR}/tools/formatting/git-clang-format.py" "-f" "--binary" "${BUILD_DEPS}/bin/clang-format"
  WORKING_DIRECTORY "${CMAKE_SOURCE_DIR}"
  COMMENT "Formatting code staged code changes with clang-format" VERBATIM
)

# make sync
add_custom_target(
  sync
  "${CMAKE_SOURCE_DIR}/tools/sync.sh" "${CMAKE_BINARY_DIR}" "${BUILD_DEPS}/lib:${BUILD_DEPS}/legacy/lib"
  WORKING_DIRECTORY "${CMAKE_SOURCE_DIR}"
  COMMENT "Generating sdk sync: ${CMAKE_BINARY_DIR}/sync"
  DEPENDS cloudquery_extensions cloudquery_amalgamation
)
