/*
 *  Copyright 2019 The searKing authors. All Rights Reserved.
 *
 *  Use of this source code is governed by a MIT-style license
 *  that can be found in the LICENSE file in the root of the source
 *  tree. An additional intellectual property rights grant can be found
 *  in the file PATENTS.  All contributing project authors may
 *  be found in the AUTHORS file in the root of the source tree.
 */
#include "stacktrace.h"

#include <cxxabi.h>    // for __cxa_demangle
#include <dlfcn.h>     // for dladdr
#include <execinfo.h>  // for backtrace

#include <sstream>
#include <string>
namespace searking {
namespace stacktrace {

std::string Stacktrace(int skip) {
  void *callstack[128];
  const int nMaxFrames = sizeof(callstack) / sizeof(callstack[0]);
  char buf[1024];
  int nFrames = backtrace(callstack, nMaxFrames);
  char **symbols = backtrace_symbols(callstack, nFrames);

  std::ostringstream trace_buf;
  for (int i = skip; i < nFrames; i++) {
    Dl_info info;
    if (dladdr(callstack[i], &info)) {
      char *demangled = NULL;
      int status;
      demangled = abi::__cxa_demangle(info.dli_sname, NULL, 0, &status);
      snprintf(buf, sizeof(buf), "%-3d %*p %s + %zd\n", i,
               static_cast<int>(2 + sizeof(void *) * 2), callstack[i],
               status == 0 ? demangled : info.dli_sname,
               reinterpret_cast<char *>(callstack[i]) -
                   reinterpret_cast<char *>(info.dli_saddr));
      free(demangled);
    } else {
      snprintf(buf, sizeof(buf), "%-3d %*p\n", i,
               static_cast<int>(2 + sizeof(void *) * 2), callstack[i]);
    }
    trace_buf << buf;

    snprintf(buf, sizeof(buf), "%s\n", symbols[i]);
    trace_buf << buf;
  }
  free(symbols);
  if (nFrames == nMaxFrames) trace_buf << "[truncated]\n";
  return trace_buf.str();
}

void SafeDumpToFd(int fd) {
  void *callstack[128];
  const int nMaxFrames = sizeof(callstack) / sizeof(callstack[0]);
  char buf[1024];
  int nFrames = backtrace(callstack, nMaxFrames);
  backtrace_symbols_fd(callstack, nFrames, fd);
}
}  // namespace stacktrace
}  // namespace searking