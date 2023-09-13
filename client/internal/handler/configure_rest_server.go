// This file is safe to edit. Once it exists it will not be overwritten

package handler

import (
	"context"
	"crypto/tls"
	"database/sql"
	"embed"
	"fmt"
	"net/http"
	"os"

	formatters "github.com/fabienm/go-logrus-formatters"
	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	er "github.com/pkg/errors"
	"github.com/pressly/goose/v3"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"client/internal/handler/operations"
	"client/internal/repository/hash_repository"
	hashservice "client/internal/service/hash"
	contact "proto/hash_service"
)

//go:generate swagger generate server --target ../../../client --name RestServer --spec ../../api/api-client.yml --server-package internal/handler --principal interface{}

//go:embed migrations/*.sql
var embedMigrations embed.FS

var gelFmt = formatters.NewGelf("client")
var Logger = logrus.Logger{
	Out:       os.Stdout,
	Formatter: gelFmt,
	Level:     logrus.DebugLevel,
	Hooks:     make(logrus.LevelHooks),
}

func configureFlags(api *operations.RestServerAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.RestServerAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})

	// Example:
	api.Logger = Logger.Printf
	api.UseSwaggerUI()
	// To continue using redoc as your UI, uncomment the following line
	// api.UseRedoc()

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()
	connStr := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable&connect_timeout=%d",
		os.Getenv("DB_TYPE"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		5)
	db, err := sql.Open(os.Getenv("DB_TYPE"), connStr)
	if err != nil {
		Logger.WithFields(logrus.Fields{"StackTrace": er.WithStack(err)}).Fatalf("sql.Open error: %s", err)
	}
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect(os.Getenv("DB_TYPE")); err != nil {
		panic(err)
	}

	if err := goose.Up(db, "migrations"); err != nil {
		panic(err)
	}
	grpcConn, err := grpc.Dial(fmt.Sprintf("%s:%s", os.Getenv("HASH_SERVICE_HOST"), os.Getenv("HASH_SERVICE_PORT")), grpc.WithInsecure())
	if err != nil {
		Logger.WithFields(logrus.Fields{"StackTrace": er.WithStack(err)}).Errorf("Dial error: %s", err)
	}
	grpcClient := contact.NewHashServiceClient(grpcConn)
	repo := hash_repository.New(db)
	service := hashservice.New(repo, grpcClient, &Logger)
	api.GetCheckHandler = operations.GetCheckHandlerFunc(func(params operations.GetCheckParams) middleware.Responder {
		res := service.GetHashedString(&params)
		return res
	})

	api.PostSendHandler = operations.PostSendHandlerFunc(func(params operations.PostSendParams) middleware.Responder {
		res := service.WriteHashedString(&params)
		return res
	})

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix".
func configureServer(s *http.Server, scheme, addr string) {}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation.
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics.
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		reqID := r.Header.Get("X-Request-Id")
		if reqID == "" {
			reqID = uuid.New().String()
		}
		ctx = context.WithValue(ctx, "X-Request-Id", reqID)
		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}
