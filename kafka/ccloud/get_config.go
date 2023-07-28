package ccloud

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func ReadConfig(configFile string) kafka.ConfigMap {
	m := make(map[string]kafka.ConfigValue)

	file, err := os.Open(configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open file: %s", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if !strings.HasPrefix(line, "#") && len(line) != 0 {
			kv := strings.Split(line, "=")
			parameter := strings.TrimSpace(kv[0])
			value := strings.TrimSpace(kv[1])
			m[parameter] = value
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Failed to read file: %s", err)
		os.Exit(1)
	}

	return m

}

func ProducerConfig() kafka.ConfigMap {
	var conf kafka.ConfigMap

	env := os.Getenv("ENV")
	if env == "" {
		// For local development purposes using locally held "properties" file.
		log.Println("defaulting Producer to Rando Caltestian Kafka")

		p, _ := os.Getwd()
		conf = ReadConfig(filepath.Join(p + "/kafka/producer/properties"))
	} else {
		conf = kafka.ConfigMap{
			"bootstrap.servers": os.Getenv("BOOT_STRAP_SERVER"),
			"security.protocol": os.Getenv("SECURITY_PROTOCOL"),
			"sasl.mechanisms":   os.Getenv("SASL_MECHANISMS"),
			"sasl.username":     os.Getenv("SASL_USERNAME"),
			"sasl.password":     os.Getenv("SASL_PASSWORD"),
			"acks":              os.Getenv("ACKS"),
		}
	}

	return conf
}

func ConsumerConfig() kafka.ConfigMap {
	var conf kafka.ConfigMap

	env := os.Getenv("ENV")
	if env == "" {
		// For local development purposes using locally held "properties" file.
		log.Println("defaulting Consumer to Rando Caltestian Kafka")

		p, _ := os.Getwd()
		conf = ReadConfig(filepath.Join(p + "/kafka/consumer/properties"))
	} else {
		conf = kafka.ConfigMap{
			"bootstrap.servers":  os.Getenv("BOOT_STRAP_SERVER"),
			"security.protocol":  os.Getenv("SECURITY_PROTOCOL"),
			"sasl.mechanisms":    os.Getenv("SASL_MECHANISMS"),
			"sasl.username":      os.Getenv("SASL_USERNAME"),
			"sasl.password":      os.Getenv("SASL_PASSWORD"),
			"enable.auto.commit": true,
			"auto.offset.reset":  "earliest",
		}
	}

	return conf
}
