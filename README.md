# Decolecta Sunat Ruc

```sh
cd ~/deco/decolecta-ruc
export DATABASE_URI="postgres://juan:createdat2024@localhost:5432/decolecta_rucs?sslmode=disable"
export REDIS_URI="localhost:6379"
go run .
```

# To override the entrypoint and run a shell in docker-compose:
# In your docker-compose.yml, add:
#
#   entrypoint: ["sh"]
#
# Example service:
#   services:
#     app:
#       image: your-image
#       entrypoint: ["sh"]