set(TABLE_CATEGORIES
  "applications"
  "compliance"
  "events"
  "networking"
  "system"
)
# For each general category, add the C++ files in the base.
foreach(CATEGORY ${TABLE_CATEGORIES})
  file(GLOB CLOUDQUERY_${CATEGORY}_TABLES
    "${CATEGORY}/*.[cpm]*"
    "${CATEGORY}/${TABLE_PLATFORM}/*.[cpm]*"
  )
  #if(CLOUDQUERY_${CATEGORY}_TABLES)
  ADD_CLOUDQUERY_LIBRARY_ADDITIONAL(cloudquery_${CATEGORY}_tables 
                                    ${CLOUDQUERY_${CATEGORY}_TABLES} )
  #endif()

endforeach()

ADD_CLOUDQUERY_LIBRARY_CORE(cloudquery_tables_utility
  "${CMAKE_CURRENT_LIST_DIR}/utility/file.cpp"
  "${CMAKE_CURRENT_LIST_DIR}/utility/cloudquery.cpp"
  "${CMAKE_CURRENT_LIST_DIR}/utility/time.cpp"
)

