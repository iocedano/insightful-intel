package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"insightful-intel/internal/database"
	"insightful-intel/internal/interactor"
	"insightful-intel/internal/repositories"

	"github.com/spf13/cobra"
)

var (
	query          string
	maxDepth       int
	skipDuplicates bool
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "cli",
	Short: "Insightful Intel CLI",
	Long:  "A CLI tool for running dynamic pipeline searches across multiple domains",
}

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run [query]",
	Short: "Run dynamic pipeline search",
	Long: `Run a dynamic pipeline search with the specified query across multiple domains.
The search will explore related entities across ONAPI, SCJ, DGII, PGR, and Google Docking.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		query := args[0]

		log.Println("Running CLI")
		db := database.New()
		pipelineResultRepo := repositories.NewPipelineRepository(db)
		scjRepo := repositories.NewScjRepository(db)
		dgiiRepo := repositories.NewDgiiRepository(db)
		pgrRepo := repositories.NewPgrRepository(db)
		googleDockingRepo := repositories.NewDockingRepository(db)
		onapiRepo := repositories.NewOnapiRepository(db)

		dynamicPipelineInteractor := interactor.NewDynamicPipelineInteractor(pipelineResultRepo, scjRepo, dgiiRepo, pgrRepo, googleDockingRepo, onapiRepo)

		log.Printf("Executing dynamic pipeline with query: %s, max depth: %d, skip duplicates: %v",
			query, maxDepth, skipDuplicates)

		dynamicPipelineInteractor.ExecuteDynamicPipeline(context.Background(), query, maxDepth, skipDuplicates)

		log.Println("Dynamic pipeline execution completed")
	},
}

func main() {
	// Initialize database
	db := database.New()

	log.Println("Running migrations")

	if err := runMigrations(db); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	log.Println("Migrations completed")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Add run command to root
	rootCmd.AddCommand(runCmd)

	// Add flags for run command
	runCmd.Flags().IntVarP(&maxDepth, "max-depth", "d", 5, "Maximum depth for pipeline execution")
	runCmd.Flags().BoolVarP(&skipDuplicates, "skip-duplicates", "s", true, "Skip duplicate searches")

	// Set description for flags
	runCmd.Flags().Lookup("max-depth").Usage = "Maximum depth to traverse in the pipeline (default: 5)"
	runCmd.Flags().Lookup("skip-duplicates").Usage = "Skip searching duplicate keywords across domains (default: true)"
}

// runMigrations executes database migrations
func runMigrations(db database.Service) error {
	migrationService := database.NewMigrationService(db.GetDB())
	migrations := database.GetInitialMigrations()
	return migrationService.RunMigrations(migrations)
}
