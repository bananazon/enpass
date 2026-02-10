package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/gdanko/enpass/util"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var UsageTemplate = `Usage:{{if .Runnable}}
{{ PosixUsage . }}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}{{$cmds := .Commands}}{{if eq (len .Groups) 0}}

Available Commands:{{range $cmds}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{else}}{{range $group := .Groups}}

{{.Title}}{{range $cmds}}{{if (and (eq .GroupID $group.ID) (or .IsAvailableCommand (eq .Name "help")))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if not .AllChildCommandsHaveGroup}}

Additional Commands:{{range $cmds}}{{if (and (eq .GroupID "") (or .IsAvailableCommand (eq .Name "help")))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`

func GetflagLists(cmd *cobra.Command) {
	getListShowFlags(cmd)
}

func GetShowFlags(cmd *cobra.Command) {
	getListShowFlags(cmd)
}

func getListShowFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&flagTrashed, "trashed", false, "Show trashed items.")
	cmd.Flags().StringArrayVarP(&flagOrderBy, "orderby", "o", []string{}, fmt.Sprintf("Specify fields to sort by. Can be used multiple times. Valid: %s", strings.Join(sort.StringSlice(validOrderBy), ", ")))
	cmd.Flags().BoolVar(&flagList, "list", false, "Output the data as list, similar to SQLite line mode.")
	cmd.Flags().BoolVar(&flagYaml, "yaml", false, "Output the data as YAML.")
	cmd.Flags().BoolVar(&flagTable, "table", false, "Output the data as a table.")
}

func GetPersistentFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&flagVaultPath, "vault", "v", "", "Path to your Enpass vault.")
	cmd.PersistentFlags().StringVar(&flagCardType, "type", "password", "The type of your card. (password, ...)")
	cmd.PersistentFlags().StringArrayVarP(&flagRecordTitle, "title", "t", []string{}, "Filter based on record title. Wildcards (%) are allowed. Can be used multiple times.")
	cmd.PersistentFlags().StringArrayVarP(&flagRecordCategory, "category", "c", []string{}, "Filter based on record category. Wildcards (%) are allowed. Can be used multiple times.")
	cmd.PersistentFlags().StringArrayVarP(&flagRecordLogin, "login", "l", []string{}, "Filter based on record login. Wildcards (%) are allowed. Can be used multiple times.")
	cmd.PersistentFlags().StringArrayVarP(&flagLabel, "label", "y", []string{}, "Filter based on record field label. Can be used multiple times")
	cmd.PersistentFlags().StringArrayVarP(&flagRecordUuid, "uuid", "u", []string{}, "Filter based on record uuid. Can be used multiple times.")
	cmd.PersistentFlags().StringVarP(&flagKeyFilePath, "keyfile", "k", "", "Path to your Enpass vault keyfile.")
	cmd.PersistentFlags().StringVar(&logLevelStr, "log", defaultLogLevel, fmt.Sprintf("The log level, one of: %s", util.ReturnLogLevels(logLevelMap)))
	cmd.PersistentFlags().BoolVarP(&flagNonInteractive, "non-interactive", "n", false, "Disable prompts and fail instead.")
	cmd.PersistentFlags().BoolVar(&flagCaseSensitive, "sensitive", false, "Force category and title searches to be case-sensitive.")
	cmd.PersistentFlags().BoolVar(&flagNoColor, "nocolor", false, "Disable colorized output and logging.")
	cmd.PersistentFlags().BoolVarP(&flagEnablePin, "pin", "p", false, "Enable PIN.")
}

func GetCopyFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&flagClipboardPrimary, "primary", false, "Use primary X selection instead of clipboard.")
	cmd.Flags().StringArrayVarP(&flagOrderBy, "orderby", "o", []string{"title"}, "Specify fields to sort by. Can be used multiple times.")
}

func GetPassFlags(cmd *cobra.Command) {
	cmd.Flags().StringArrayVarP(&flagOrderBy, "orderby", "o", []string{"title"}, "Specify fields to sort by. Can be used multiple times.")
}

func addFlags(flagSlice *[]string, f *pflag.Flag) {
	flagType := strings.ToLower(f.Value.Type())

	if strings.Contains(flagType, "array") || strings.Contains(flagType, "slice") {
		*flagSlice = append(*flagSlice,
			fmt.Sprintf("[--%s <opt> ...]", f.Name),
		)
	} else if flagType == "count" {
		*flagSlice = append(*flagSlice,
			fmt.Sprintf("[--%s ...]", f.Name),
		)
	} else if strings.Contains(flagType, "int") || strings.Contains(flagType, "float") {
		*flagSlice = append(*flagSlice,
			fmt.Sprintf("[--%s <%s>]", f.Name, flagType),
		)
	} else if flagType == "string" {
		*flagSlice = append(*flagSlice,
			fmt.Sprintf("[--%s <%s>]", f.Name, flagType),
		)
	} else if flagType == "bool" {
		*flagSlice = append(*flagSlice,
			fmt.Sprintf("[--%s]", f.Name),
		)
	} else {
		*flagSlice = append(*flagSlice,
			fmt.Sprintf("[--%s <%s>]", f.Name, flagType),
		)
	}
}

func GenerateUsageFromFlags(cmd *cobra.Command, binaryName string, subcommand string, epilog string, examples []string, positionals []string) string {
	var (
		flagSlice      []string
		joined         [][]string
		maxWidth       int    = 70
		usageIndent    int    = 8
		usageIndentStr string = strings.Repeat(" ", usageIndent)
		usageSlice     []string
	)

	cmd.PersistentFlags().VisitAll(func(f *pflag.Flag) {
		addFlags(&flagSlice, f)
	})

	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		addFlags(&flagSlice, f)
	})

	sort.Strings(flagSlice)

	joined = ChunkWords(flagSlice, maxWidth)
	if len(positionals) > 0 {
		joined = append(joined, positionals)
	}

	usageSlice = []string{
		fmt.Sprintf("  %s %s %s", binaryName, subcommand, strings.Join(joined[0], " ")),
	}

	for _, line := range joined[1:] {
		usageSlice = append(usageSlice, usageIndentStr+strings.Join(line, " "))
	}

	return strings.Join(usageSlice, "\n")
}

func ChunkWords(words []string, maxWidth int) [][]string {
	var result [][]string
	var current []string
	currentLen := 0

	for _, w := range words {
		wordLen := len(w)

		// +1 accounts for the space between words
		if currentLen > 0 && currentLen+1+wordLen > maxWidth {
			result = append(result, current)
			current = nil
			currentLen = 0
		}

		if currentLen > 0 {
			currentLen++ // space
		}
		current = append(current, w)
		currentLen += wordLen
	}

	if len(current) > 0 {
		result = append(result, current)
	}

	return result
}
