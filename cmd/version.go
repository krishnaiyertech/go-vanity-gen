// SPDX-FileCopyrightText: Copyright 2025 Krishna Iyer <www.krishnaiyer.tech>
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var (
	version   = ""
	gitCommit = ""
	buildDate = ""
)

// VersionCommand prints version information to the output stream.
func VersionCommand(root *cobra.Command) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Display version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("%s\n", root.Name())
			fmt.Println("----------------")
			fmt.Printf("Version: %s\n", version)
			fmt.Printf("Git Commit: %s\n", gitCommit)
			fmt.Printf("Build Date: %s\n", buildDate)
			fmt.Printf("Go version: %s\n", runtime.Version())
			fmt.Printf("OS/Arch: %s\n", runtime.GOOS+"/"+runtime.GOARCH)
		},
	}
}
