package constants

const (
	Broker1 = "localhost:9092"
	Broker2 = "localhost:9093"
	Broker3 = "localhost:9094"

	Broker1SASLPlain = "localhost:9082"
	Broker2SASLPlain = "localhost:9083"
	Broker3SASLPlain = "localhost:9084"
)

var (
	KafkaCluster = []string{Broker1, Broker2, Broker3}

	KafkaClusterSecure = []string{Broker1SASLPlain, Broker2SASLPlain, Broker3SASLPlain}
)
