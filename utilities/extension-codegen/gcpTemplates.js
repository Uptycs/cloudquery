/*
 * Copyright (c) 2020 Uptycs, Inc. All rights reserved
 */

'use strict';

module.exports = {
  list_call_template: `
	listCall := service.{{ paging_api }}({{ expanded_paging_api_args }})
`,

  list_call_service_api_template: `
	myApiService := {{ service_api }}(service)
	if myApiService == nil {
		fmt.Println("{{ service_api }}() returned nil")
		return resultMap, fmt.Errorf("{{ service_api }}() returned nil")
	}

	listCall := myApiService.{{ paging_api }}({{ expanded_paging_api_args }})`,

  page_invocation_template: `	if err := listCall.Pages(ctx, func(page *{{ page_type }}) error {
		for _, item := range page.{{ page_list_name }} {
			result := make(map[string]string)
{{ expanded_col_mapping }}
			resultMap = append(resultMap, result)
		}
		return nil
	}); err != nil {
		fmt.Println("listCall.Page: ", err)
		//log.Fatal(err)
		return resultMap, err
	}`,

  pageable_iteration_template: `	p := iterator.NewPager(listCall, 10, "")
	for {
		var items []*{{ page_type }}
		pageToken, err := p.NextPage(&items)
		if err != nil {
			fmt.Println("NextPage() error: ", err)
			return resultMap, err
		}
		for _, item := range items {
			result := make(map[string]string)
{{ expanded_col_mapping }}
			resultMap = append(resultMap, result)
		}
		if pageToken == "" {
			break
		}
	}`,

  base_template: `
package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/kolide/osquery-go/plugin/table"
	"google.golang.org/api/option"

{{ expanded_imports }}
)

func {{ name }}Columns() []table.ColumnDefinition {
	var _, _ = strconv.Atoi("123") // Disables warning when strcov is not used
	return []table.ColumnDefinition{
{{ expanded_col_list }}
	}
}

func {{ name }}Generate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	var _ = queryContext
	resultMap := make([]map[string]string, 0)
	ctx, cancel := context.WithCancel(osqCtx)
	defer cancel()
	service, err := {{ client_api }}(ctx, option.WithCredentialsFile(*keyFile))
	if err != nil {
		fmt.Println("NewService() error: ", err)
		return resultMap, err
	}{{ list_call_template }}
	if listCall == nil {
		fmt.Println("listCall is nil")
		return resultMap, nil
	}
{{ expanded_pagination }}
	return resultMap, nil
}
`,

  plugin_helper_template: `
package main

import (
	"github.com/kolide/osquery-go"
	"github.com/kolide/osquery-go/plugin/table"
)

func registerPlugins(server *osquery.ExtensionManagerServer) {
{{ plugin_list }}
}
`
};
