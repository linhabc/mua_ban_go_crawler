# Run

go run crawler_phone_using_worker.go data_type.go util.go db.go

# using docker

sudo docker run --name=crawler_go --mount source=output,destination=/app/output --mount source=db,destination=/app/db linhabc/mua_ban_go_crawler

# Using docker-compose

docker-compose up

# output folder

output: store generated json file
ddb: store generated leveldb folder
