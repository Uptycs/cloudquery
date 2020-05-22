#include <syscall.h>
#include <unistd.h>
#include <sys/eventfd.h>
#include <sys/types.h>
#include <sys/socket.h>

#include <osquery/logger.h>

#if defined(__GLIBC__)
#undef __GLIBC__
#endif
#include <boost/process.hpp>

extern "C" int execvpe(const char* file,
                       char* const argv[],
                       char* const envp[]) {
  char** env = const_cast<char**>(envp);
  return ::boost::process::detail::posix::execvpe(file, argv, env);
}

extern "C" int fallocate(int fd, int mode, off_t offset, off_t len) {
  return syscall(SYS_fallocate, fd, mode, offset, len);
}

extern "C" int epoll_create1(int flags) {
  int fd = syscall(SYS_epoll_create1, flags);
  return fd;
}

extern "C" int timerfd_settime(int fd,
                               int flags,
                               void* new_value,
                               void* old_value) {
  return syscall(SYS_timerfd_settime, fd, flags, new_value, old_value);
}

extern "C" int timerfd_create(int clockid, int flags) {
  return syscall(SYS_timerfd_create, clockid, flags);
}

extern "C" int eventfd(int initval, int flags) {
  return syscall(SYS_eventfd, initval, flags);
}

extern "C" int sync_file_range(int fd,
                               off64_t offset,
                               off64_t count,
                               unsigned int flags) {
  return syscall(SYS_sync_file_range, fd, offset, count, flags);
}

extern "C" int sched_getcpu(void) {
  unsigned cpu;
  if (syscall(SYS_getcpu, &cpu, nullptr, nullptr) == -1) {
    return -1;
  } else {
    return cpu;
  }
}

extern "C" int pthread_setname_np(pthread_t target_thread, const char* name) {
  return 0;
}

extern "C" int mkostemp(char* templt, int flags) {
  LOG(ERROR) << __FUNCTION__ <<": Not implemented";
  return 0;
}

extern "C" void __isoc99_sscanf() {
  LOG(ERROR) << __FUNCTION__ <<": Not implemented";
}

extern "C" int __vasprintf_chk(char** result_ptr,
                               int /* flags */,
                               const char* fmt,
                               va_list args) {
  return vasprintf(result_ptr, fmt, args);
}

extern "C" int __asprintf_chk(char** result_ptr,
                              int /*flags*/,
                              const char* fmt,
                              ...) {
  va_list args;
  int rc;

  va_start(args, fmt);
  rc = vasprintf(result_ptr, fmt, args);
  va_end(args);

  return rc;
}

extern "C" int eventfd_write(int fd, eventfd_t value) {
  int rc = write(fd, &value, sizeof(eventfd_t));
  if (rc != sizeof(eventfd_t)) {
    return -1;
  } else {
    return 0;
  }
}

extern "C" int eventfd_read(int fd, eventfd_t *value) {
  int rc = read(fd, value, sizeof(eventfd_t));
  if (rc != sizeof(eventfd_t)) {
    return -1;
  } else {
    return 0;
  }
}

extern "C" int accept4(int sockfd, struct sockaddr *addr,
                       socklen_t *addrlen, int flags) {
  LOG(ERROR) << __FUNCTION__ <<": Not implemented";
  return -1;
}
