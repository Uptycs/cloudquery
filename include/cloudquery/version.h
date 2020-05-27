#pragma once

#include <string>

// clang-format off
#ifndef STR
#define STR_OF(x) #x
#define STR(x) STR_OF(x)
#endif

#ifdef WIN32
#define STR_EX(...) __VA_ARGS__
#else
#define STR_EX(x) x
#endif
#define CONCAT(x, y) STR(STR_EX(x)STR_EX(y))

#ifndef FRIEND_TEST
#define FRIEND_TEST(test_case_name, test_name) \
  friend class test_case_name##_##test_name##_Test
#endif
// clang-format on

#ifdef WIN32
#define USED_SYMBOL
#define EXPORT_FUNCTION __declspec(dllexport)
#else
#define USED_SYMBOL __attribute__((used))
#define EXPORT_FUNCTION
#endif


#if !defined(CLOUDQUERY_BUILD_SDK_VERSION)
#error The build must define CLOUDQUERY_BUILD_SDK_VERSION.
#elif !defined(CLOUDQUERY_BUILD_PLATFORM)
#error The build must define CLOUDQUERY_BUILD_PLATFORM.
#elif !defined(CLOUDQUERY_BUILD_DISTRO)
#error The build must define CLOUDQUERY_BUILD_DISTRO.
#endif

#define CLOUDQUERY_SDK_VERSION STR(CLOUDQUERY_BUILD_SDK_VERSION)
#define CLOUDQUERY_PLATFORM STR(CLOUDQUERY_BUILD_PLATFORM)

/**
 * @brief A series of platform-specific home folders.
 *
 * There are several platform-specific folders where osquery reads and writes
 * content. Most of the variance is due to legacy support.
 *
 * CLOUDQUERY_HOME: Configuration, flagfile, extensions and module autoload.
 * CLOUDQUERY_DB_HOME: Location of RocksDB persistent storage.
 * CLOUDQUERY_LOG_HOME: Location of log data when the filesystem plugin is used.
 */
#if defined(__linux__)
#define CLOUDQUERY_HOME "/etc/cloudquery/"
#define CLOUDQUERY_DB_HOME "/var/cloudquery/"
#define CLOUDQUERY_SOCKET CLOUDQUERY_DB_HOME
#define CLOUDQUERY_PIDFILE "/var/run/"
#define CLOUDQUERY_LOG_HOME "/var/log/cloudquery/"
#define CLOUDQUERY_CERTS_HOME "/usr/share/cloudquery/certs/"
#elif defined(WIN32)
#define CLOUDQUERY_HOME "\\Program Files\\uptycs\\cloudquery\\"
#define CLOUDQUERY_DB_HOME CLOUDQUERY_HOME
#define CLOUDQUERY_SOCKET "\\\\.\\pipe\\"
#define CLOUDQUERY_PIDFILE CLOUDQUERY_DB_HOME
#define CLOUDQUERY_LOG_HOME CLOUDQUERY_HOME "log\\"
#define CLOUDQUERY_CERTS_HOME CLOUDQUERY_HOME "certs\\"
#elif defined(FREEBSD)
#define CLOUDQUERY_HOME "/var/db/cloudquery/"
#define CLOUDQUERY_DB_HOME CLOUDQUERY_HOME
#define CLOUDQUERY_SOCKET "/var/run/"
#define CLOUDQUERY_PIDFILE "/var/run/"
#define CLOUDQUERY_LOG_HOME "/var/log/cloudquery/"
#define CLOUDQUERY_CERTS_HOME "/etc/ssl/"
#else
#define CLOUDQUERY_HOME "/var/cloudquery/"
#define CLOUDQUERY_DB_HOME CLOUDQUERY_HOME
#define CLOUDQUERY_SOCKET CLOUDQUERY_DB_HOME
#define CLOUDQUERY_PIDFILE CLOUDQUERY_DB_HOME
#define CLOUDQUERY_LOG_HOME "/var/log/cloudquery/"
#define CLOUDQUERY_CERTS_HOME CLOUDQUERY_HOME "certs/"
#endif

#define CLOUDQUERY_SHARED_SECTION "CLOUDQueryMemory"
#define CLOUDQUERY_SHARED_SECTION_SIZE 1024
#define CLOUDQUERY_CONFIG_NAME "Config Time"
#define CLOUDQUERY_DIST_NAME "Distributed Time"
#define CLOUDQUERY_SCHED_NAME "Scheduler Time"

/// A configuration error is catastrophic and should exit the watcher.
#define EXIT_CATASTROPHIC 78

namespace osquery
{

    namespace cloudquery
    {
        extern const std::string kCQVersion;
        extern const std::string kCQSDKVersion;

    } // namespace cloudquery
    /// Custom literal for size_t.
    size_t operator"" _sz(unsigned long long int x);
} // namespace osquery
