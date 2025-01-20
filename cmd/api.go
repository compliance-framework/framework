package cmd

import (
	"context"
	"github.com/compliance-framework/framework/api"
	"github.com/compliance-framework/framework/api/handler"
	"github.com/compliance-framework/framework/event/bus"
	apiRuntime "github.com/compliance-framework/framework/runtime"
	"github.com/compliance-framework/framework/service"
	mongoStore "github.com/compliance-framework/framework/store/mongo"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"log"
	"os"
)

func ApiCmd() *cobra.Command {
	var apiCmd = &cobra.Command{
		Use: "api",
		//		Short: "long-running agent for continuously checking policies against plugin data",
		//		Long: `The Continuous Compliance Agent is a long running process that continuously checks policy controls
		//with plugins to ensure continuous compliance.`,
		RunE: apiRunner,
	}

	return apiCmd
}

const (
	DefaultMongoURI = "mongodb://localhost:27017"
	DefaultNATSURI  = "nats://localhost:4222"
	DefaultPort     = ":8080"
)

type Config struct {
	MongoURI string
	NatsURI  string
}

//	@title			Compliance Framework Configuration Service API
//	@version		1.0
//	@description	This is the API for the Compliance Framework Configuration Service.

//	@host		localhost:8080
//	@BasePath	/api

// @externalDocs.description	OpenAPI
// @externalDocs.url			https://swagger.io/resources/open-api/
func apiRunner(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Can't initialize zap logger: %v", err)
	}
	sugar := logger.Sugar()

	config := loadConfig()

	mongoDatabase, err := connectMongo(ctx, options.Client().ApplyURI(config.MongoURI), "cf")
	if err != nil {
		sugar.Fatal(err)
	}
	defer mongoDatabase.Client().Disconnect(ctx)

	err = bus.Listen(config.NatsURI, sugar)
	if err != nil {
		sugar.Fatal(err)
	}

	server := api.NewServer(ctx, sugar)

	catalogStore := mongoStore.NewCatalogStore(mongoDatabase)
	catalogHandler := handler.NewCatalogHandler(catalogStore)
	catalogHandler.Register(server.API().Group("/catalog"))

	planService := service.NewPlanService(mongoDatabase, bus.Publish)
	planHandler := handler.NewPlanHandler(sugar, planService)
	planHandler.Register(server.API().Group("/plan"))

	resultService := service.NewResultsService(mongoDatabase)
	resultHandler := handler.NewResultsHandler(sugar, resultService, planService)
	resultHandler.Register(server.API().Group("/results"))

	resultProcessor := apiRuntime.NewProcessor(bus.Subscribe[apiRuntime.ExecutionResult], planService, resultService)
	resultProcessor.Listen()

	plansService := service.NewPlansService(mongoDatabase, bus.Publish)
	plansHandler := handler.NewPlansHandler(sugar, plansService)
	plansHandler.Register(server.API().Group("/plans"))

	systemPlanService := service.NewSSPService(mongoDatabase)
	systemPlanHandler := handler.NewSSPHandler(systemPlanService)
	systemPlanHandler.Register(server.API())

	metadataService := service.NewMetadataService(mongoDatabase)
	metadataHandler := handler.NewMetadataHandler(metadataService)
	metadataHandler.Register(server.API().Group("/metadata"))

	server.PrintRoutes()

	return server.Start(DefaultPort)
}

func connectMongo(ctx context.Context, clientOptions *options.ClientOptions, databaseName string) (*mongo.Database, error) {
	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return client.Database(databaseName), nil
}

func loadConfig() (config Config) {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = DefaultMongoURI
	}

	natsURI := os.Getenv("NATS_URI")
	if natsURI == "" {
		natsURI = DefaultNATSURI
	}

	config = Config{
		MongoURI: mongoURI,
		NatsURI:  natsURI,
	}
	return config
}

func checkErr(err error) {
	if err != nil {
		log.Fatalf("An error occurred: %v", err)
	}
}
