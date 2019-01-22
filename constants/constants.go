package constants

const (
	Broker1 = "localhost:9092"
	Broker2 = "localhost:9093"
	Broker3 = "localhost:9094"
)

var (
	KafkaCluster = []string{Broker1, Broker2, Broker3}
)
