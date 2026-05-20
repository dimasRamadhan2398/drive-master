set -e

export CONFIG_PATH=/app/config.yaml

echo "Running database migrations..."
./core-service migrate

echo "Running database seeders..."
./core-service seed

echo "Starting User Service..."
./core-service serve