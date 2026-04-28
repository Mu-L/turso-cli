package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tursodatabase/turso-cli/internal/flags"
	"github.com/tursodatabase/turso-cli/internal/turso"
)

var groupFilter string
var schemaFilter string

func init() {
	dbCmd.AddCommand(listCmd)
	listCmd.Flags().StringVarP(&groupFilter, "group", "g", "", "Filter databases by group")
	listCmd.Flags().StringVarP(&schemaFilter, "schema", "s", "", "Filter databases by schema")
}

type DatabaseFetcher struct {
	client       *turso.Client
	SchemaFilter string
	GroupFilter  string
	ParentDbId   string
	LoadFullInfo bool
}

func (df *DatabaseFetcher) FetchPage(pageSize int, cursor *string) (turso.ListResponse, error) {
	if !flags.V3Api() {
		return df.fetchPageV2(pageSize, cursor)
	}
	if df.SchemaFilter != "" {
		return df.fetchPageV2(pageSize, cursor)
	}
	orgID, err := tryResolveOrgID(df.client)
	if err != nil {
		return turso.ListResponse{}, err
	}
	if orgID == "" {
		return df.fetchPageV2(pageSize, cursor)
	}
	groupID := ""
	if df.GroupFilter != "" {
		id, err := tryResolveGroupID(df.client, df.GroupFilter)
		if err != nil {
			return turso.ListResponse{}, err
		}
		if id == "" {
			return df.fetchPageV2(pageSize, cursor)
		}
		groupID = id
	}
	cursorStr := ""
	if cursor != nil {
		cursorStr = *cursor
	}
	options := turso.DatabaseV3ListOptions{
		GroupId:    groupID,
		Limit:      pageSize,
		Cursor:     cursorStr,
		ParentDbId: df.ParentDbId,
	}
	dbs, next, err := df.client.DatabasesV3.List(orgID, options)
	if err != nil {
		return turso.ListResponse{}, err
	}
	response := turso.ListResponse{Databases: dbs}
	if next != "" {
		response.Pagination = &turso.Pagination{Next: &next}
	}
	if df.LoadFullInfo {
		for i := range response.Databases {
			response.Databases[i], err = df.client.DatabasesV3.Get(orgID, response.Databases[i].ID)
			if err != nil {
				return turso.ListResponse{}, err
			}
		}
	}
	return response, nil
}

func (df *DatabaseFetcher) fetchPageV2(pageSize int, cursor *string) (turso.ListResponse, error) {
	cursorStr := ""
	if cursor != nil {
		cursorStr = *cursor
	}

	options := turso.DatabaseListOptions{
		Group:  groupFilter,
		Schema: schemaFilter,
		Limit:  pageSize,
		Cursor: cursorStr,
		Parent: df.ParentDbId,
	}

	response, err := df.client.Databases.List(options)
	if err != nil {
		return turso.ListResponse{}, err
	}
	if df.LoadFullInfo {
		for i := range response.Databases {
			response.Databases[i], err = df.client.Databases.Get(response.Databases[i].Name)
			if err != nil {
				return turso.ListResponse{}, err
			}
		}
	}
	return response, nil
}

var listCmd = &cobra.Command{
	Use:               "list",
	Short:             "List databases.",
	Args:              cobra.NoArgs,
	ValidArgsFunction: noFilesArg,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		client, err := authedTursoClient()
		if err != nil {
			return err
		}

		fetcher := &DatabaseFetcher{
			client:       client,
			SchemaFilter: schemaFilter,
			GroupFilter:  groupFilter,
		}
		return printDatabaseList(fetcher)
	},
}
