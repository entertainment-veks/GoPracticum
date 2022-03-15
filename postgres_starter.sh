docker run -d --name postgres -e POSTGRES_DB=shortener -e POSTGRES_PASSWORD=postgres -p 5432:5432 postgres
docker start postgres