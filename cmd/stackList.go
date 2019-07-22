package cmd

import (
	"fmt"
	"github.com/greenled/portainer-stack-utils/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"text/template"
)

// stackListCmd represents the remove command
var stackListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List stacks",
	Aliases: []string{"ls"},
	Example: "psu stack list --endpoint 1",
	Run: func(cmd *cobra.Command, args []string) {
		stacks, err := common.GetAllStacksFiltered(common.StackListFilter{
			SwarmId:    viper.GetString("stack.list.swarm"),
			EndpointId: viper.GetUint32("stack.list.endpoint"),
		})
		common.CheckError(err)

		if viper.GetBool("stack.list.quiet") {
			// Print only stack names
			for _, s := range stacks {
				_, err := fmt.Println(s.Name)
				common.CheckError(err)
			}
		} else if viper.GetString("stack.list.format") != "" {
			// Print stack fields formatted
			template, templateParsingErr := template.New("stackTpl").Parse(viper.GetString("stack.list.format"))
			common.CheckError(templateParsingErr)
			for _, s := range stacks {
				templateExecutionErr := template.Execute(os.Stdout, s)
				common.CheckError(templateExecutionErr)
				fmt.Println()
			}
		} else {
			// Print all stack fields as a table
			writer, err := common.NewTabWriter([]string{
				"STACK ID",
				"NAME",
				"TYPE",
				"ENTRY POINT",
				"PROJECT PATH",
				"ENDPOINT ID",
				"SWARM ID",
			})
			common.CheckError(err)
			for _, s := range stacks {
				_, err := fmt.Fprintln(writer, fmt.Sprintf(
					"%v\t%s\t%v\t%s\t%s\t%v\t%s",
					s.Id,
					s.Name,
					s.GetTranslatedStackType(),
					s.EntryPoint,
					s.ProjectPath,
					s.EndpointID,
					s.SwarmID,
				))
				common.CheckError(err)
			}
			flushErr := writer.Flush()
			common.CheckError(flushErr)
		}
	},
}

func init() {
	stackCmd.AddCommand(stackListCmd)

	stackListCmd.Flags().String("swarm", "", "filter by swarm ID")
	stackListCmd.Flags().String("endpoint", "", "filter by endpoint ID")
	stackListCmd.Flags().BoolP("quiet", "q", false, "only display stack names")
	stackListCmd.Flags().String("format", "", "format output using a Go template")
	viper.BindPFlag("stack.list.swarm", stackListCmd.Flags().Lookup("swarm"))
	viper.BindPFlag("stack.list.endpoint", stackListCmd.Flags().Lookup("endpoint"))
	viper.BindPFlag("stack.list.quiet", stackListCmd.Flags().Lookup("quiet"))
	viper.BindPFlag("stack.list.format", stackListCmd.Flags().Lookup("format"))
}