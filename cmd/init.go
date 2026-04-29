package cmd

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(doctorCmd)
	rootCmd.AddCommand(initCmd)
}