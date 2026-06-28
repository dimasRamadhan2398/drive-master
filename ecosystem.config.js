module.exports = {
  apps: [
    {
      name: "user-service",
      script: "./backend/user-service/bin/user-service",
      cwd: "/var/www/drive",
      args: "serve",
      instances: 1,
      autorestart: true,
      watch: false,
      env: {
        PORT: 8001,
        POSTGRES_HOST: "127.0.0.1",
        POSTGRES_PORT: 5432,
        POSTGRES_USER: "admin_drive",
        POSTGRES_PASSWORD: "drivemaster123",
        POSTGRES_DB: "drivemaster_user_service",
        REDIS_HOST: "127.0.0.1",
        REDIS_PORT: 6379,
        REDIS_PASSWORD: "",
        KAFKA_BROKERS: "127.0.0.1:9092",
        KAFKA_ENABLED: "true",
        KAFKA_TOPIC: "user.events",
        KAFKA_GROUP_ID: "user-service-consumer",
        CONFIG_PATH: "./backend/user-service/pkg/config/config.yaml",
        RUN_MIGRATIONS: "true",
        RUN_SEEDERS: "true"
      }
    },
    {
      name: "core-service",
      script: "./backend/core-service/bin/core-service",
      cwd: "/var/www/drive",
      args: "serve",
      instances: 1,
      autorestart: true,
      watch: false,
      env: {
        PORT: 8002,
        POSTGRES_HOST: "127.0.0.1",
        POSTGRES_PORT: 5432,
        POSTGRES_USER: "admin_drive",
        POSTGRES_PASSWORD: "drivemaster123",
        POSTGRES_DB: "drivemaster_core_service",
        REDIS_HOST: "127.0.0.1",
        REDIS_PORT: 6379,
        REDIS_PASSWORD: "",
        KAFKA_BROKERS: "127.0.0.1:9092",
        KAFKA_ENABLED: "true",
        KAFKA_TOPIC: "core.events",
        KAFKA_GROUP_ID: "core-service-consumer",
        CONFIG_PATH: "./backend/core-service/pkg/config/config.yaml",
        RUN_MIGRATIONS: "true",
        RUN_SEEDERS: "true",
        GIN_MODE: "release"
      }
    },
    {
      name: "booking-service",
      script: "./backend/booking-service/bin/booking-service",
      cwd: "/var/www/drive",
      args: "serve",
      instances: 1,
      autorestart: true,
      watch: false,
      env: {
        PORT: 8003,
        POSTGRES_HOST: "127.0.0.1",
        POSTGRES_PORT: 5432,
        POSTGRES_USER: "admin_drive",
        POSTGRES_PASSWORD: "drivemaster123",
        POSTGRES_DB: "drivemaster_booking_service",
        REDIS_HOST: "127.0.0.1",
        REDIS_PORT: 6379,
        REDIS_PASSWORD: "",
        KAFKA_BROKERS: "127.0.0.1:9092",
        KAFKA_ENABLED: "true",
        KAFKA_TOPIC: "booking.events",
        KAFKA_GROUP_ID: "booking-service-consumer",
        CONFIG_PATH: "./backend/booking-service/pkg/config/config.yaml",
        RUN_MIGRATIONS: "true",
        RUN_SEEDERS: "true",
        USER_SERVICE_URL: "http://127.0.0.1:8001",
        CORE_SERVICE_URL: "http://127.0.0.1:8002"
      }
    },
    {
      name: "payment-service",
      script: "./backend/payment-service/bin/payment-service",
      cwd: "/var/www/drive",
      args: "serve",
      instances: 1,
      autorestart: true,
      watch: false,
      env: {
        PORT: 8004,
        POSTGRES_HOST: "127.0.0.1",
        POSTGRES_PORT: 5432,
        POSTGRES_USER: "admin_drive",
        POSTGRES_PASSWORD: "drivemaster123",
        POSTGRES_DB: "drivemaster_payment_service",
        REDIS_HOST: "127.0.0.1",
        REDIS_PORT: 6379,
        REDIS_PASSWORD: "",
        KAFKA_BROKERS: "127.0.0.1:9092",
        KAFKA_ENABLED: "true",
        KAFKA_TOPIC: "payment.events",
        CONFIG_PATH: "./backend/payment-service/pkg/config/config.yaml",
        RUN_MIGRATIONS: "true",
        RUN_SEEDERS: "true"
      }
    },
    {
      name: "drive-frontend",
      script: "./frontend/.output/server/index.mjs",
      cwd: "/var/www/drive",
      instances: 1,
      autorestart: true,
      watch: false,
      env: {
        PORT: 3000,
        NUXT_PUBLIC_API_BASE: "https://drivemaster.id/api/v1",
        NUXT_PUBLIC_MODE: "prod"
      }
    }
  ]
};
