package bootstrap

import (
	"github.com/web3sphere/backend/configs"
	"github.com/web3sphere/backend/pkg/cache"
	"github.com/web3sphere/backend/pkg/database"
	"github.com/web3sphere/backend/pkg/jwt"
	"github.com/web3sphere/backend/pkg/logger"
	"github.com/web3sphere/backend/pkg/mailer"
	"github.com/web3sphere/backend/pkg/queue/kafka"
	"github.com/web3sphere/backend/pkg/queue/rabbitmq"
	"github.com/web3sphere/backend/pkg/storage"
	"gorm.io/gorm"
)

// Container holds all the global dependencies for the application.
type Container struct {
	Config     *configs.Config
	Logger     *logger.Logger
	DB         *gorm.DB
	Redis      *cache.RedisClient
	RabbitMQ   *rabbitmq.Client
	Kafka      *kafka.Producer
	Mailer     mailer.Mailer
	Storage    storage.Storage
	JWTManager *jwt.Manager
}

// InitContainer initializes all the infrastructure tools and returns a Container
// along with a cleanup function to gracefully close connections on shutdown.
func InitContainer() (*Container, func(), error) {
	// 1. Load Configuration
	cfg := configs.Load()

	// 2. Initialize Logger
	log := logger.New(cfg.App.Env, cfg.App.Debug)
	log.Infof("Starting %s v%s in %s mode", cfg.App.Name, cfg.App.Version, cfg.App.Env)

	// 3. Initialize Database
	db, err := database.New(&cfg.Database, log)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
		return nil, nil, err
	}

	// 4. Initialize Redis
	redisClient, err := cache.New(&cfg.Redis, log)
	if err != nil {
		log.Fatalf("Failed to initialize Redis: %v", err)
		return nil, nil, err
	}

	// 5. Initialize RabbitMQ
	rmqClient, err := rabbitmq.New(&cfg.RabbitMQ, log)
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ: %v", err)
		return nil, nil, err
	}

	// 6. Initialize Kafka (Infrastructure only, continue if it fails)
	kafkaProducer, err := kafka.NewProducer(&cfg.Kafka, log)
	if err != nil {
		log.Warnf("Failed to initialize Kafka Producer: %v (continuing without Kafka)", err)
	}

	// 7. Initialize Mailer
	mailSvc, err := mailer.New(cfg, log)
	if err != nil {
		log.Fatalf("Failed to initialize Mailer: %v", err)
		return nil, nil, err
	}

	// 8. Initialize Storage
	storageSvc, err := storage.New(cfg, log)
	if err != nil {
		log.Fatalf("Failed to initialize Storage: %v", err)
		return nil, nil, err
	}

	// 9. Initialize JWT Manager
	jwtManager := jwt.NewManager(&cfg.JWT)

	// Create cleanup function
	cleanup := func() {
		log.Info("Cleaning up infrastructure connections...")
		database.Close(db, log)
		redisClient.Close()
		rmqClient.Close()
		if kafkaProducer != nil {
			kafkaProducer.Close()
		}
		_ = log.Sync() // Sync logger last
	}

	return &Container{
		Config:     cfg,
		Logger:     log,
		DB:         db,
		Redis:      redisClient,
		RabbitMQ:   rmqClient,
		Kafka:      kafkaProducer,
		Mailer:     mailSvc,
		Storage:    storageSvc,
		JWTManager: jwtManager,
	}, cleanup, nil
}
