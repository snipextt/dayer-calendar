package storage

func Init() {
	connectToRedis()
	connectToMongoDb()
	createIndexes()
  connectToKafka()
}
