package cmd

import (
	"log"
	"net"
	"sync"

	"github.com/cvetkovski98/zvax-common/gen/pbslot"
	"github.com/cvetkovski98/zvax-common/pkg/healthz"
	"github.com/cvetkovski98/zvax-common/pkg/redis"
	"github.com/cvetkovski98/zvax-slots/internal/config"
	"github.com/cvetkovski98/zvax-slots/internal/delivery"
	"github.com/cvetkovski98/zvax-slots/internal/repository"
	"github.com/cvetkovski98/zvax-slots/internal/service"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var (
	runCommand = &cobra.Command{
		Use:   "run",
		Short: "Run slots microservice",
		Long:  `Run slots microservice`,
		Run:   run,
	}
	network string
	address string
)

func init() {
	runCommand.Flags().StringVarP(&network, "network", "n", "tcp", "network to listen on")
	runCommand.Flags().StringVarP(&address, "address", "a", ":50052", "address to listen on")
}

func run(cmd *cobra.Command, args []string) {
	cfg := config.GetConfig()
	rdb, err := redis.NewRedisConn(cfg.Redis)
	if err != nil {
		log.Fatalf("failed to connect to Redis: %v", err)
	}
	slotRepository := repository.NewRedisSlotRepository(rdb)
	slotService := service.NewSlotServiceImpl(slotRepository)
	slotGrpc := delivery.NewSlotGrpcServerImpl(slotService)
	server := grpc.NewServer()
	pbslot.RegisterSlotGrpcServer(server, slotGrpc)

	wg := new(sync.WaitGroup)
	wg.Add(2)
	go func() {
		healtzSrv := healthz.CreateServer(80)
		if err := healtzSrv.ListenAndServe(); err != nil {
			log.Printf("error running healthz server: %v", err)
		}
		log.Println("Running healthz...")
		wg.Done()
	}()

	go func() {
		lis, err := net.Listen(network, address)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		if err := server.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
		log.Printf("Listening on %s://%s...", network, address)
		wg.Done()
	}()
	wg.Wait()
}
