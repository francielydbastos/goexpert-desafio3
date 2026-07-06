package main

import (
	"database/sql"
	"log"
	"net"
	"net/http"
	"time"

	"clean-architecture/configs"
	"clean-architecture/internal/infra/database"
	graphqlserver "clean-architecture/internal/infra/graph"
	"clean-architecture/internal/infra/grpc/pb"
	grpcservice "clean-architecture/internal/infra/grpc/service"
	"clean-architecture/internal/infra/web"
	"clean-architecture/internal/usecase"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg := configs.LoadConfig()

	db := waitForDB(cfg)
	defer db.Close()

	runMigrations(db)

	orderRepository := database.NewOrderRepository(db)
	createOrderUseCase := usecase.NewCreateOrderUseCase(orderRepository)
	listOrdersUseCase := usecase.NewListOrdersUseCase(orderRepository)

	// REST
	go func() {
		handler := web.NewOrderHandler(createOrderUseCase, listOrdersUseCase)
		server := web.NewWebServer(handler)
		addr := ":" + cfg.WebServerPort
		log.Printf("Starting REST server on %s", addr)
		if err := server.Start(addr); err != nil {
			log.Fatalf("REST server error: %v", err)
		}
	}()

	// gRPC
	go func() {
		grpcServer := grpc.NewServer()
		orderService := grpcservice.NewOrderService(createOrderUseCase, listOrdersUseCase)
		pb.RegisterOrderServiceServer(grpcServer, orderService)
		reflection.Register(grpcServer)
		addr := ":" + cfg.GRPCServerPort
		lis, err := net.Listen("tcp", addr)
		if err != nil {
			log.Fatalf("gRPC listen error: %v", err)
		}
		log.Printf("Starting gRPC server on %s", addr)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("gRPC server error: %v", err)
		}
	}()

	// GraphQL
	resolver := &graphqlserver.Resolver{
		CreateOrderUseCase: createOrderUseCase,
		ListOrdersUseCase:  listOrdersUseCase,
	}
	srv := handler.NewDefaultServer(graphqlserver.NewExecutableSchema(graphqlserver.Config{Resolvers: resolver}))
	mux := http.NewServeMux()
	mux.Handle("/", playground.Handler("GraphQL playground", "/query"))
	mux.Handle("/query", srv)
	addr := ":" + cfg.GraphQLServerPort
	log.Printf("Starting GraphQL server on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("GraphQL server error: %v", err)
	}
}

func waitForDB(cfg *configs.Config) *sql.DB {
	dsn := cfg.DSN()
	var db *sql.DB
	var err error
	for i := 0; i < 30; i++ {
		db, err = sql.Open(cfg.DBDriver, dsn)
		if err == nil {
			if pingErr := db.Ping(); pingErr == nil {
				log.Println("Database is ready")
				return db
			} else {
				err = pingErr
			}
		}
		log.Printf("Waiting for database... (%d/30): %v", i+1, err)
		time.Sleep(2 * time.Second)
	}
	log.Fatalf("Could not connect to database: %v", err)
	return nil
}

func runMigrations(db *sql.DB) {
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		log.Fatalf("migration driver error: %v", err)
	}
	m, err := migrate.NewWithDatabaseInstance("file://migrations", "mysql", driver)
	if err != nil {
		log.Fatalf("migration init error: %v", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("migration up error: %v", err)
	}
	log.Println("Migrations applied successfully")
}
