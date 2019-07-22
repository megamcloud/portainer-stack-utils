package cmd

import (
	"fmt"
	"github.com/greenled/portainer-stack-utils/common"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"net/url"

	"github.com/spf13/cobra"
)

// stackRemoveCmd represents the remove command
var stackRemoveCmd = &cobra.Command{
	Use:     "remove STACK_NAME",
	Short:   "Remove a stack",
	Aliases: []string{"rm", "down"},
	Example: "psu stack rm mystack",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		stackName := args[0]
		stack, err := common.GetStackByName(stackName)

		switch err.(type) {
		case nil:
			// The stack exists
			common.PrintVerbose(fmt.Sprintf("Stack %s exists.", stackName))

			stackId := stack.Id

			common.PrintVerbose(fmt.Sprintf("Removing stack %s...", stackName))
			reqUrl, err := url.Parse(fmt.Sprintf("%s/api/stacks/%d", viper.GetString("url"), stackId))
			common.CheckError(err)

			req, err := http.NewRequest(http.MethodDelete, reqUrl.String(), nil)
			common.CheckError(err)
			headerErr := common.AddAuthorizationHeader(req)
			common.CheckError(headerErr)
			common.PrintDebugRequest("Remove stack request", req)

			client := common.NewHttpClient()

			resp, err := client.Do(req)
			common.PrintDebugResponse("Remove stack response", resp)
			common.CheckError(err)

			common.CheckError(common.CheckResponseForErrors(resp))
		case *common.StackNotFoundError:
			// The stack does not exist
			common.PrintVerbose(fmt.Sprintf("Stack %s does not exist.", stackName))
			if viper.GetBool("stack.remove.strict") {
				log.Fatalln(fmt.Sprintf("Stack %s does not exist.", stackName))
			}
		default:
			// Something else happened
			common.CheckError(err)
		}
	},
}

func init() {
	stackCmd.AddCommand(stackRemoveCmd)

	stackRemoveCmd.Flags().Bool("strict", false, "fail if stack does not exist")
	viper.BindPFlag("stack.remove.strict", stackRemoveCmd.Flags().Lookup("strict"))
}