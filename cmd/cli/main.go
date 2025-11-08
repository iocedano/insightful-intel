package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"insightful-intel/internal/database"
	"insightful-intel/internal/interactor"
	"insightful-intel/internal/repositories"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

// contextKey is a custom type for context keys
type contextKey string

const executionIDKey contextKey = "executionID"

var (
	query          string
	maxDepth       int
	skipDuplicates bool
)

// GetExecutionID retrieves the execution ID from the context
func GetExecutionID(ctx context.Context) (string, bool) {
	executionID, ok := ctx.Value(executionIDKey).(string)
	return executionID, ok
}

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
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		query := args[0]
		skipDuplicates := args[1] == "true"

		// Generate a unique execution ID
		executionID := uuid.New()

		log.Printf("Starting CLI execution with ID: %s", executionID.String())

		db := database.New()

		repositoryFactory := repositories.NewRepositoryFactory(db)
		dynamicPipelineInteractor := interactor.NewDynamicPipelineInteractor(repositoryFactory)

		// Create context with execution ID
		ctx := context.WithValue(context.Background(), executionIDKey, executionID.String())

		log.Printf("Executing dynamic pipeline [%s] with query: %s, max depth: %d, skip duplicates: %v",
			executionID.String(), query, maxDepth, skipDuplicates)

		err := dynamicPipelineInteractor.ExecuteDynamicPipeline(ctx, query, maxDepth, skipDuplicates)
		if err != nil {
			log.Fatalf("[%s] failed to execute dynamic pipeline: %v", executionID.String(), err)
		}

		log.Printf("[%s] Dynamic pipeline execution completed", executionID.String())
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
