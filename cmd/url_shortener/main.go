package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"

	grpcValidator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	urlShortenerV1 "github.com/iamtonydev/url-shortener/internal/app/api/url_shortener_v1"
	"github.com/iamtonydev/url-shortener/internal/app/repository"
	"github.com/iamtonydev/url-shortener/internal/app/repository/custom_storage"
	urlShortener "github.com/iamtonydev/url-shortener/internal/app/service/url_shortener"
	"github.com/iamtonydev/url-shortener/internal/config"
	desc "github.com/iamtonydev/url-shortener/pkg/url_shortener_v1"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
)

func main() {
	var useCustomStorage = flag.Bool("custom_storage", false, "use custom storage")
	var wg sync.WaitGroup
	var cfg *config.Config
	var err error

	flag.Parse()

	cfg, err = config.Read("config.yml")

	if err != nil {
		log.Fatal("failed to open configuration file")
		return
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Fatal(startGRPC(cfg, *useCustomStorage))
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Fatal(startHTTP(cfg))
	}()

	wg.Wait()
}

func startGRPC(cfg *config.Config, useCustomStorage bool) error {
	list, err := net.Listen("tcp", cfg.Grpc.Port)
	if err != nil {
		return fmt.Errorf("failed to mapping port: %s", err.Error())
	}
	defer list.Close()

	var urlsRepository repository.IUrlsRepository

	if useCustomStorage {
		urlsRepository = custom_storage.NewCustomStorage()
	} else {
		dbDsn := fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.Name, cfg.Database.SslMode,
		)

		db, err := sqlx.Open(cfg.Database.Driver, dbDsn)
		if err != nil {
			return fmt.Errorf("failed to open connection with db")
		}
		defer db.Close()

		urlsRepository = repository.NewUrlsRepository(db)
	}

	urlShortenerService := urlShortener.NewUrlShortenerService(urlsRepository)

	s := grpc.NewServer(
		grpc.StreamInterceptor(grpcValidator.StreamServerInterceptor()),
		grpc.UnaryInterceptor(grpcValidator.UnaryServerInterceptor()),
	)
	desc.RegisterUrlShortenerV1Server(s, urlShortenerV1.NewUrlShortenerV1(urlShortenerService))

	if err = s.Serve(list); err != nil {
		return fmt.Errorf("failed to url_shortener: %s", err.Error())
	}

	return nil
}

func startHTTP(cfg *config.Config) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	err := desc.RegisterUrlShortenerV1HandlerFromEndpoint(ctx, mux, cfg.Grpc.Port, opts)
	if err != nil {
		return err
	}

	return http.ListenAndServe(cfg.Http.Port, mux)
}
