cmake_policy(SET CMP0068 NEW)

set(CLOUDQUERY_SOURCES "")

set(CLOUDQUERY_LINKS "")

set(CLOUDQUERY_TESTS "")
set(CLOUDQUERY_BENCHMARKS "")

# This will contain the native cloudquery tables.
set(CLOUDQUERY_ADDITIONAL_SOURCES "")

set(CLOUDQUERY_ADDITIONAL_LINKS "")
set(CLOUDQUERY_ADDITIONAL_TESTS "")
set(CLOUDQUERY_TABLES_TESTS "")

# Add all and extra for cloudquery code.
if(CLANG AND POSIX)
  add_compile_options(
    -Wall
    -Wextra
    -pedantic
    -Wuseless-cast
    -Wno-c99-extensions
    -Wno-zero-length-array
    -Wno-unused-parameter
    -Wno-gnu-case-range
    -Weffc++
  )
  if(NOT FREEBSD)
    add_compile_options(
      -Wshadow-all
      -Wno-shadow-field
    )
  endif()
endif()

# Construct a set of all object files, starting with third-party and all
# of the cloudquery objects (sources from ADD_CORE_LIBRARY macros).
#if(FREEBSD)
 # set(CLOUDQUERY_OBJECTS $<TARGET_OBJECTS:cloudquery_sqlite>)
## set(CLOUDQUERY_OBJECTS
#    $<TARGET_OBJECTS:cloudquery_sqlite>
 #   $<TARGET_OBJECTS:linenoise-ng>
 #

add_library(libcloudquery STATIC main/lib.cpp)

#include(examples/CMakeLists.txt)

if(NOT SKIP_TABLES)
  add_subdirectory(tables)
  # Amalgamate the utility tables needed to compile.
  GENERATE_UTILITIES("${CMAKE_SOURCE_DIR}")
  AMALGAMATE("${CMAKE_SOURCE_DIR}" "utils" AMALGAMATION_UTILS)
  ADD_CLOUDQUERY_LIBRARY_CORE(cloudquery_amalgamation ${AMALGAMATION_UTILS})
endif()

if(WINDOWS)
  if(CMAKE_BUILD_TYPE EQUAL "DEBUG")
    set(WB_KEY "sgd")
    set(WO_KEY "_dbg")
    set(WT_KEY "d_dbg")
  else()
    set(WB_KEY "s")
    set(WO_KEY "")
  endif()
  ADD_CLOUDQUERY_LINK_CORE("ntdll.lib")
  ADD_CLOUDQUERY_LINK_CORE("ws2_32.lib")
  ADD_CLOUDQUERY_LINK_CORE("iphlpapi.lib")
  ADD_CLOUDQUERY_LINK_CORE("netapi32.lib")
  ADD_CLOUDQUERY_LINK_CORE("rpcrt4.lib")
  ADD_CLOUDQUERY_LINK_CORE("shlwapi.lib")
  ADD_CLOUDQUERY_LINK_CORE("version.lib")
  ADD_CLOUDQUERY_LINK_CORE("Wtsapi32.lib")
  ADD_CLOUDQUERY_LINK_CORE("wbemuuid.lib")
  ADD_CLOUDQUERY_LINK_CORE("Secur32.lib")
  ADD_CLOUDQUERY_LINK_CORE("taskschd.lib")
  ADD_CLOUDQUERY_LINK_CORE("dbghelp.lib")
  ADD_CLOUDQUERY_LINK_CORE("dbgeng.lib")
  ADD_CLOUDQUERY_LINK_CORE("libboost_chrono-mt-${WB_KEY}")
  ADD_CLOUDQUERY_LINK_CORE("libboost_system-mt-${WB_KEY}")
  ADD_CLOUDQUERY_LINK_CORE("libboost_regex-mt-${WB_KEY}")
  ADD_CLOUDQUERY_LINK_CORE("libboost_filesystem-mt-${WB_KEY}")
  ADD_CLOUDQUERY_LINK_CORE("libboost_context-mt-${WB_KEY}")
  ADD_CLOUDQUERY_LINK_CORE("libboost_thread-mt-${WB_KEY}")
  ADD_CLOUDQUERY_LINK_CORE("rocksdb${WO_KEY}")
  ADD_CLOUDQUERY_LINK_CORE("thriftmt${WT_KEY}")
  ADD_CLOUDQUERY_LINK_CORE("gflags_static${WO_KEY}")
  ADD_CLOUDQUERY_LINK_CORE("ssleay32")
  ADD_CLOUDQUERY_LINK_CORE("eay32")
  ADD_CLOUDQUERY_LINK_CORE("zlibstatic")
  ADD_CLOUDQUERY_LINK_CORE("libyara64${WO_KEY}")

  # Enable control flow guard
  ADD_CLOUDQUERY_LINK_CORE("-guard:cf")
else()
  ADD_CLOUDQUERY_LINK_CORE("libpthread")
  ADD_CLOUDQUERY_LINK_CORE("libz")
  ADD_CLOUDQUERY_LINK_CORE("gflags")
  ADD_CLOUDQUERY_LINK_CORE("thrift")
endif()

if(LINUX)
  ADD_CLOUDQUERY_LINK_CORE("tcmalloc_minimal")
endif()

if(APPLE OR LINUX)
  ADD_CLOUDQUERY_LINK_CORE("libdl")
  ADD_CLOUDQUERY_LINK_CORE("rocksdb_lite")
  ADD_CLOUDQUERY_LINK_CORE("rocksdb_lite")
elseif(FREEBSD)
  ADD_CLOUDQUERY_LINK_CORE("icuuc")
  ADD_CLOUDQUERY_LINK_CORE("linenoise")
  ADD_CLOUDQUERY_LINK_CORE("rocksdb-lite")
  ADD_CLOUDQUERY_LINK_CORE("rocksdb-lite")
endif()

if(POSIX)
  ADD_CLOUDQUERY_LINK_CORE("boost_system")
  ADD_CLOUDQUERY_LINK_CORE("boost_filesystem")
  ADD_CLOUDQUERY_LINK_CORE("boost_thread")
  ADD_CLOUDQUERY_LINK_CORE("boost_context")
  ADD_CLOUDQUERY_LINK_CORE("boost_regex")
endif()

if(LINUX OR FREEBSD)
  ADD_CLOUDQUERY_LINK_CORE("librt")
  ADD_CLOUDQUERY_LINK_CORE("libc")
endif()

# Remaining additional development libraries.
ADD_CLOUDQUERY_LINK_CORE("glog${WO_KEY}")
if (WINDOWS)
  ADD_CLOUDQUERY_LINK_CORE("grpc")
	ADD_CLOUDQUERY_LINK_CORE("grpc++")
	ADD_CLOUDQUERY_LINK_CORE("gpr")
	ADD_CLOUDQUERY_LINK_CORE("libprotobuf")
	ADD_CLOUDQUERY_LINK_CORE("cares_static")
	ADD_CLOUDQUERY_LINK_CORE("upb")
	ADD_CLOUDQUERY_LINK_CORE("address_sorting")
else()
	ADD_CLOUDQUERY_LINK_CORE("grpc++")
  ADD_CLOUDQUERY_LINK_CORE("grpc")
  ADD_CLOUDQUERY_LINK_CORE("protobuf")
  ADD_CLOUDQUERY_LINK_CORE("udev")
  ADD_CLOUDQUERY_LINK_CORE("audit")
  ADD_CLOUDQUERY_LINK_CORE("pcap")
  ADD_CLOUDQUERY_LINK_CORE("libresolv.so")
  ADD_CLOUDQUERY_LINK_CORE("tsk")
  ADD_CLOUDQUERY_LINK_CORE("aws-cpp-sdk-core")
  ADD_CLOUDQUERY_LINK_CORE("aws-cpp-sdk-sts")
  ADD_CLOUDQUERY_LINK_CORE("aws-cpp-sdk-ec2")
  ADD_CLOUDQUERY_LINK_CORE("yara")
  ADD_CLOUDQUERY_LINK_CORE("lldpctl")
  ADD_CLOUDQUERY_LINK_CORE("smartctl")
  ADD_CLOUDQUERY_LINK_CORE("elf++")
  ADD_CLOUDQUERY_LINK_CORE("cryptsetup devmapper lvm2app lvm-internal daemonclient")
  ADD_CLOUDQUERY_LINK_CORE("gcrypt gpg-error")
  ADD_CLOUDQUERY_LINK_CORE("blkid")
  ADD_CLOUDQUERY_LINK_CORE("ip4tc")
  ADD_CLOUDQUERY_LINK_CORE("dpkg lzma bz2")
  ADD_CLOUDQUERY_LINK_CORE("rpm rpmio popt db")
  ADD_CLOUDQUERY_LINK_CORE("elf++")
  ADD_CLOUDQUERY_LINK_CORE("augeas fa xml2")
  ADD_CLOUDQUERY_LINK_CORE("fuzzy")
  
endif()

if(POSIX)
  # Hashing methods in core use libcrypto.
  ADD_CLOUDQUERY_LINK_CORE("crypto")

  ADD_CLOUDQUERY_LINK_CORE("ssl")
  ADD_CLOUDQUERY_LINK_CORE("libpthread")
  ADD_CLOUDQUERY_LINK_CORE("magic")
endif()

if(APPLE)
  ADD_CLOUDQUERY_LINK_CORE("liblzma libbz2")
else()
  if(POSIX)
    ADD_CLOUDQUERY_LINK_CORE("lzma")
  endif()
  ADD_CLOUDQUERY_LINK_CORE("bz2")
endif()

# The platform-specific SDK + core linker flags.
if(POSIX)
  ADD_CLOUDQUERY_LINK_CORE("-rdynamic")
endif()

if(APPLE)
  ADD_CLOUDQUERY_LINK_CORE("-Wl,-dead_strip")
  ADD_CLOUDQUERY_LINK_CORE("-mmacosx-version-min=${OSX_VERSION_MIN}")
  ADD_CLOUDQUERY_LINK_CORE("-Wl,-cache_path_lto,${CMAKE_BINARY_DIR}/ltocache")
  ADD_CLOUDQUERY_LINK_CORE("-Wl,-no_weak_imports")
elseif(LINUX OR FREEBSD)
  ADD_CLOUDQUERY_LINK_CORE("-Wl,-zrelro -Wl,-znow")
  if(NOT DEFINED ENV{SANITIZE} AND (CMAKE_BUILD_TYPE STREQUAL "RELEASE"))
    ADD_CLOUDQUERY_LINK_CORE("-pie")
  endif()
endif()

if(LINUX)
  ADD_CLOUDQUERY_LINK_CORE("uuid")
  if(NOT DEFINED ENV{CLOUDQUERY_BUILD_LINK_SHARED})
    ADD_CLOUDQUERY_LINK_CORE("-static-libstdc++")
  endif()
  # For Ubuntu/CentOS packages add the build SHA1.
  ADD_CLOUDQUERY_LINK_CORE("-Wl,--build-id")
  if (CLANG AND DEPS)
    # If using GCC, libgcc_s may be needed.
    ADD_CLOUDQUERY_LINK_CORE("-fuse-ld=lld")
    ADD_CLOUDQUERY_LINK_CORE("c++")
    ADD_CLOUDQUERY_LINK_CORE("c++abi")
    ADD_CLOUDQUERY_LINK_CORE("unwind")
    if(NOT CMAKE_CXX_COMPILER_VERSION VERSION_LESS 5.0.0)
      ADD_CLOUDQUERY_LINK_CORE("-Wl,--thinlto-cache-dir=${CMAKE_BINARY_DIR}/cache")
      ADD_CLOUDQUERY_LINK_CORE("-Wl,--thinlto-cache-policy,cache_size_bytes=2g")
    endif()
    ADD_CLOUDQUERY_LINK_CORE("-B${BUILD_DEPS}/legacy/lib")
  endif()
endif()

if(NOT WINDOWS)
  ADD_CLOUDQUERY_LINK_CORE("archive")
  ADD_CLOUDQUERY_LINK_CORE("zstd")
else()
  ADD_CLOUDQUERY_LINK_CORE("zstd_static")
  ADD_CLOUDQUERY_LINK_CORE("archive_static")
endif()

list(APPEND CLOUDQUERY_OBJECTS ${CLOUDQUERY_SOURCES})
list(APPEND CLOUDQUERY_LIBS ${CLOUDQUERY_LINKS})

set(CMAKE_MACOSX_RPATH 0)
set(CMAKE_SKIP_RPATH TRUE)

# A final map from CMake build platform to a enum using for runtime detection.
# The goal is to provide a minimum set of compile code paths.
# See ./include/core.h for the enum class.
# POSIX   = 0x01
# WINDOWS = 0x02
# BSD     = 0x04
# LINUX   = 0x08 && POSIX
# OS X    = 0x10 && BSD && POSIX
# FREEBSD = 0x20 && BSD && POSIX
if(WINDOWS)
  math(EXPR PLATFORM_MASK "2")
elseif(LINUX)
  math(EXPR PLATFORM_MASK "1 + 8")
elseif(APPLE)
  math(EXPR PLATFORM_MASK "1 + 4 + 16")
elseif(FREEBSD)
  math(EXPR PLATFORM_MASK "1 + 4 + 32")
endif()

# Create the static libcloudquery (everything but non-utility tables).
set(CLOUDQUERY_LIBRARY_FLAGS
  -DCLOUDQUERY_BUILD_VERSION=${CLOUDQUERY_BUILD_VERSION}
  -DCLOUDQUERY_PLATFORM_MASK=${PLATFORM_MASK}
)
JOIN("${CLOUDQUERY_LIBRARY_FLAGS}" " " CLOUDQUERY_LIBRARY_FLAGS)

target_sources(libcloudquery PRIVATE ${CLOUDQUERY_OBJECTS})
target_link_libraries(libcloudquery ${CLOUDQUERY_LIBS})
set_target_properties(libcloudquery PROPERTIES OUTPUT_NAME cloudquery)
set_target_properties(libcloudquery PROPERTIES COMPILE_FLAGS "${CLOUDQUERY_LIBRARY_FLAGS}")

# Create the dynamic libcloudquery.
if(DEFINED ENV{CLOUDQUERY_BUILD_SHARED})
  add_library(a_shared SHARED main/lib.cpp ${CLOUDQUERY_OBJECTS})
  target_link_libraries(libcloudquery_shared PRIVATE ${CLOUDQUERY_LIBS})
  set_target_properties(libcloudquery_shared PROPERTIES COMPILE_FLAGS "${CLOUDQUERY_LIBRARY_FLAGS}")
  set_target_properties(libcloudquery_shared PROPERTIES OUTPUT_NAME cloudquery)
endif()



if(NOT SKIP_TABLES)
  # Generate the cloudquery additional tables (the non-util).
  GENERATE_TABLES("${CMAKE_SOURCE_DIR}")
  AMALGAMATE("${CMAKE_SOURCE_DIR}" "additional" AMALGAMATION)
  ADD_CLOUDQUERY_LIBRARY_ADDITIONAL(cloudquery_additional_amalgamation ${AMALGAMATION})

  #AMALGAMATE("${CMAKE_SOURCE_DIR}" "foreign" AMALGAMATION_FOREIGN)
 # ADD_CLOUDQUERY_LIBRARY_ADDITIONAL(cloudquery_foreign_amalgamation ${AMALGAMATION_FOREIGN})
endif()

#Create the static libcloudquery_additional.
add_library(libcloudquery_additional STATIC ${CLOUDQUERY_ADDITIONAL_SOURCES})
target_link_libraries(libcloudquery_additional ${CLOUDQUERY_ADDITIONAL_LINKS})
set_target_properties(libcloudquery_additional PROPERTIES OUTPUT_NAME cloudquery_additional)

if(DEFINED ENV{CLOUDQUERY_BUILD_SHARED})
  add_library(libcloudquery_additional_shared SHARED ${CLOUDQUERY_ADDITIONAL_SOURCES})
  target_link_libraries(libcloudquery_additional_shared PRIVATE ${CLOUDQUERY_LINKS} ${CLOUDQUERY_ADDITIONAL_LINKS})
  set_target_properties(libcloudquery_additional_shared PROPERTIES OUTPUT_NAME cloudquery_additional)
endif()

set(DAEMON_SOURCE_FILES
  devtools/shell.cpp
  main/main.cpp
  main/${PROCESS_FAMILY}/main.cpp)


# Add the versioning information for Windows binaries
if(WINDOWS)
  # The windows version resource requires a version string, as well
  # as version components separated by a comma. First we replace the `-`
  # values that may exist for commit hash versions, then we grab the
  # major version components
  string(REPLACE "-" "." CLOUDQUERY_VER_TEMP "${CLOUDQUERY_BUILD_VERSION}")
  string(REPLACE "." ";" CLOUDQUERY_VERSION_LIST "${CLOUDQUERY_VER_TEMP}")
  list(GET CLOUDQUERY_VERSION_LIST 0 CLOUDQUERY_PRODUCT_NUMBER)
  list(GET CLOUDQUERY_VERSION_LIST 1 CLOUDQUERY_PRODUCT_VERSION)
  list(GET CLOUDQUERY_VERSION_LIST 2 CLOUDQUERY_BUILD_NUMBER)

  configure_file(
    ${CMAKE_SOURCE_DIR}/tools/windows_resources.rc.in
    ${CMAKE_BINARY_DIR}/cloudquery/windows_resources.rc
    @ONLY)
  list(APPEND DAEMON_SOURCE_FILES ${CMAKE_BINARY_DIR}/cloudquery/windows_resources.rc)
  endif()
  add_executable(daemon
  ${DAEMON_SOURCE_FILES}
  )  
  add_dependencies(daemon libcloudquery)
  add_dependencies(daemon libcloudquery_additional)

  ADD_DEFAULT_LINKS(daemon TRUE)
  SET_CLOUDQUERY_COMPILE(daemon)
  set_target_properties(daemon PROPERTIES OUTPUT_NAME cloudqueryd)

  # A friendly echo printed before building the daemon
  add_custom_command(TARGET daemon PRE_BUILD
    COMMAND ${CMAKE_COMMAND} -E cmake_echo_color --cyan --bold
      "-- Building cloudqueryd: $<TARGET_FILE:daemon>")

  set(SHELL_EXECUTABLE "cloudqueryi")
  if(WINDOWS)

  set(SHELL_EXECUTABLE "${SHELL_EXECUTABLE}.exe")
endif()

add_custom_target(shell ALL
  COMMAND ${CMAKE_COMMAND} -E copy
    $<TARGET_FILE:daemon> "$<TARGET_FILE_DIR:daemon>/${SHELL_EXECUTABLE}")

add_custom_command(
  TARGET shell POST_BUILD
  COMMAND ${CMAKE_COMMAND} -E cmake_echo_color --cyan --bold
    "-- Creating cloudqueryi: $<TARGET_FILE_DIR:daemon>/${SHELL_EXECUTABLE}")


# make install (executables)
install(TARGETS daemon RUNTIME DESTINATION bin COMPONENT main)
install(PROGRAMS ${CMAKE_BINARY_DIR}/cloudquery/cloudqueryi DESTINATION bin COMPONENT main)
install(FILES "${CMAKE_SOURCE_DIR}/tools/deployment/cloudqueryctl"
  DESTINATION bin COMPONENT main)

# make install (config files)
install(FILES "${CMAKE_SOURCE_DIR}/tools/deployment/cloudquery.example.conf"
  DESTINATION "${CMAKE_INSTALL_PREFIX}/share/cloudquery/" COMPONENT main)
install(DIRECTORY DESTINATION "${CMAKE_INSTALL_PREFIX}/var/cloudquery")

install(DIRECTORY "${CMAKE_SOURCE_DIR}/packs/"
  DESTINATION "${CMAKE_INSTALL_PREFIX}/share/cloudquery/packs" COMPONENT main)
if(APPLE)
  install(FILES "${CMAKE_SOURCE_DIR}/tools/deployment/com.facebook.cloudqueryd.plist"
    DESTINATION "${CMAKE_INSTALL_PREFIX}/share/cloudquery/" COMPONENT main)
else()
  install(PROGRAMS "${CMAKE_SOURCE_DIR}/tools/deployment/cloudqueryd.initd"
    DESTINATION "${CMAKE_INSTALL_PREFIX}/etc/init.d/" RENAME "cloudqueryd" COMPONENT main)
endif()

if(NOT SKIP_TESTS)
  #add tests here
endif()
