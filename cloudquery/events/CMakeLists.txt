#set (SKIP_EVENTS TRUE)
if(NOT SKIP_EVENTS)
  if(APPLE)
    #no-op
  elseif(FREEBSD)
    # no-op
    # no files for FREEBSD
  elseif(WINDOWS)
  
  else()
    #no-op
  endif()

  if(POSIX)
    #ADD_CLOUDQUERY_LINK_CORE("curl")
    #ADD_CLOUDQUERY_LINK_CORE("idn")
    ADD_CLOUDQUERY_LIBRARY(FALSE osquery_events_linux
    "${CMAKE_CURRENT_LIST_DIR}/api_events_publisher.cpp"
    "${CMAKE_CURRENT_LIST_DIR}/api_events_publisher.h"
    )
  endif()
endif()

