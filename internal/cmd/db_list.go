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
	dbs, err := df.client.DatabasesV3.List(orgID, options)
	if err != nil {
		return turso.ListResponse{}, err
	}
	return turso.ListResponse{Databases: dbs}, nil
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
	}

	return df.client.Databases.List(options)
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
