Run a container

```bash
docker run --name pg-container \
  -e POSTGRES_PASSWORD=1234 \
  -e POSTGRES_DB=pgtictac \
  -p 5432:5432 \
  -d postgres
```

Create an .env file with the following names:

```env
DB_USER=postgres
DB_PASSWORD=1234
DB_HOST=localhost
DB_PORT=5432
DB_NAME=pgtictac
DB_SSLMODE=disable
```