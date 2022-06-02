package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"

	grpcvalidator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/iamtonydev/url-shortener/internal/app/api/url_shortener_v1"
	"github.com/iamtonydev/url-shortener/internal/app/repository"
	"github.com/iamtonydev/url-shortener/internal/app/service/url_shortener"
	desc "github.com/iamtonydev/url-shortener/pkg/url_shortener_v1"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
)

const (
	host       = "localhost"
	dbPort     = "54321"
	dbUser     = "postgres"
	dbPassword = "postgres"
	dbName     = "url_shortener"
	sslMode    = "disable"
	grpcPort   = ":50051"
	httpPort   = ":8000"
)

func main() {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Fatal(startGRPC())
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Fatal(startHTTP())
	}()

	wg.Wait()
}

func startGRPC() error {
	list, err := net.Listen("tcp", grpcPort)
	if err != nil {
		return fmt.Errorf("failed to mapping port: %s", err.Error())
	}
	defer list.Close()

	dbDsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, dbPort, dbUser, dbPassword, dbName, sslMode,
	)

	db, err := sqlx.Open("pgx", dbDsn)
	if err != nil {
		return fmt.Errorf("failed to open connection with db")
	}
	defer db.Close()

	urlShortenerRepository := repository.NewUrlShortenerRepository(db)
	urlShortenerService := url_shortener.NewUrlShortenerService(urlShortenerRepository)

	s := grpc.NewServer(
		grpc.StreamInterceptor(grpcvalidator.StreamServerInterceptor()),
		grpc.UnaryInterceptor(grpcvalidator.UnaryServerInterceptor()),
	)
	desc.RegisterUrlShortenerV1Server(s, url_shortener_v1.NewUrlShortenerV1(urlShortenerService))

	if err = s.Serve(list); err != nil {
		return fmt.Errorf("failed to server: %s", err.Error())
	}

	return nil
}

func startHTTP() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	err := desc.RegisterUrlShortenerV1HandlerFromEndpoint(ctx, mux, grpcPort, opts)
	if err != nil {
		return err
	}

	return http.ListenAndServe(httpPort, mux)
}
