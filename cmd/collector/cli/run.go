package cli

import (
	"fmt"
	"io/ioutil"

	"github.com/replicatedhq/troubleshoot/pkg/collect"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Run() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "run a single collector",
		Long:  `...`,
		PreRun: func(cmd *cobra.Command, args []string) {
			viper.BindPFlag("collector", cmd.Flags().Lookup("collector"))
			viper.BindPFlag("redact", cmd.Flags().Lookup("redact"))
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			v := viper.GetViper()

			specContents, err := ioutil.ReadFile(v.GetString("collector"))
			if err != nil {
				return err
			}

			c, err := collect.ParseSpec(string(specContents))
			if err != nil {
				return err
			}

			collector := collect.Collector{
				Collect: c,
				Redact:  v.GetBool("redact"),
			}
			b, err := collector.RunCollectorSync()
			if err != nil {
				return err
			}

			fmt.Printf("%s", b)

			return nil
		},
	}

	cmd.Flags().String("collector", "", "path to a single collector spec to collect")
	cmd.Flags().Bool("redact", true, "enable/disable default redactions")

	cmd.MarkFlagRequired("collector")

	viper.BindPFlags(cmd.Flags())

	return cmd
}
