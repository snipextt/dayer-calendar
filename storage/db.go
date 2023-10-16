package storage

func Init() {
	InitRedisConnection()
	InitMongodbConnection()
	CreateIndexes()
}
