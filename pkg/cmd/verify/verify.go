package verify

import (
	"fmt"
	"os"
	"strings"

	"github.com/jenkins-x/jx-extsecret/pkg/extsecrets/secretfacade"
	"github.com/jenkins-x/jx-logging/pkg/log"
	"github.com/jenkins-x/jx-promote/pkg/common"
	"github.com/jenkins-x/jx/v2/pkg/cmd/helper"
	"github.com/jenkins-x/jx/v2/pkg/cmd/templates"
	"github.com/jenkins-x/jx/v2/pkg/table"
	"github.com/jenkins-x/jx/v2/pkg/util"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	verifyLong = templates.LongDesc(`
		Verifies that the ExternalSecret resources have the required properties populated in the underlying secret storage
`)

	verifyExample = templates.Examples(`
		%s verify
	`)
)

// Options the options for the command
type Options struct {
	secretfacade.Options

	Results []*secretfacade.SecretError
}

// NewCmdVerify creates a command object for the command
func NewCmdVerify() (*cobra.Command, *Options) {
	o := &Options{}

	cmd := &cobra.Command{
		Use:     "verify",
		Short:   "Verifies that the ExternalSecret resources have the required properties populated in the underlying secret storage",
		Long:    verifyLong,
		Example: fmt.Sprintf(verifyExample, common.BinaryName),
		Run: func(cmd *cobra.Command, args []string) {
			err := o.Run()
			helper.CheckErr(err)
		},
	}
	cmd.Flags().StringVarP(&o.Namespace, "ns", "n", "", "the namespace to filter the ExternalSecret resources")
	return cmd, o
}

// Run implements the command
func (o *Options) Run() error {
	results, err := o.Verify()
	if err != nil {
		return errors.Wrap(err, "failed to verify secrets")
	}
	o.Results = results

	if len(results) == 0 {
		log.Logger().Infof("the %d ExternalSecrets are %s", len(o.ExternalSecrets), util.ColorInfo("valid"))
		return nil
	}

	t := table.CreateTable(os.Stdout)
	t.AddRow("SECRET", "KEY", "MISSING PROPERTIES")
	for _, r := range results {
		for _, e := range r.EntryErrors {
			t.AddRow(r.ExternalSecret.Name, e.Key, strings.Join(e.Properties, ", "))
		}
	}
	t.Render()
	return nil
}