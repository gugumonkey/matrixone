// Copyright 2021 Matrix Origin
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mokafka

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/matrixorigin/matrixone/pkg/logutil"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry"
	"github.com/gogo/protobuf/proto"
	"github.com/google/uuid"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/jhump/protoreflect/dynamic"
	"github.com/matrixorigin/matrixone/pkg/common/moerr"
	"github.com/matrixorigin/matrixone/pkg/common/mpool"
	"github.com/matrixorigin/matrixone/pkg/container/batch"
	"github.com/matrixorigin/matrixone/pkg/container/types"
	"github.com/matrixorigin/matrixone/pkg/container/vector"
)

type ValueType string

const (
	TypeKey             = "type"
	TopicKey            = "topic"
	ValueKey            = "value"
	BootstrapServersKey = "bootstrap.servers"
	ProtobufSchemaKey   = "protobuf.schema"
	ProtobufMessagekey  = "protobuf.message"

	SchemaRegistryKey = "schema.registry"

	JSON       ValueType = "json"
	AVRO       ValueType = "avro"
	PROTOBUF   ValueType = "protobuf"
	PROTOBUFSR ValueType = "protobuf_sr"
)

type DataGetter interface {
	GetFieldValue(name string) (interface{}, bool)
}

type JsonDataGetter struct {
	Key   []byte
	Value []byte
	Data  map[string]interface{} // Cache the parsed JSON for efficiency
}

func (j *JsonDataGetter) GetFieldValue(name string) (interface{}, bool) {
	// If the JSON data hasn't been parsed, do it now
	if j.Data == nil {
		err := json.Unmarshal(j.Value, &j.Data)
		if err != nil {
			return nil, false
		}
	}

	val, ok := j.Data[name]
	return val, ok
}

type ProtoDataGetter struct {
	Value *dynamic.Message
	Key   any
}

func (p *ProtoDataGetter) GetFieldValue(name string) (interface{}, bool) {
	val := p.Value.GetFieldByName(name)
	return val, val != nil
}

type KafkaAdapterInterface interface {
	InitSchemaRegistry(url string) error
	Close()
	CreateTopic(ctx context.Context, topicName string, partitions int, replicationFactor int) error
	DescribeTopicDetails(ctx context.Context, topicName string) (*kafka.TopicMetadata, error)
	ReadMessagesFromPartition(topic string, partition int32, offset int64, limit int) ([]*kafka.Message, error)
	ReadMessagesFromTopic(topic string, offset int64, limit int64) ([]*kafka.Message, error)
	GetSchemaForTopic(topic string, isKey bool) (schemaregistry.SchemaMetadata, error)

	GetKafkaConsumer() (*kafka.Consumer, error)
	ProduceMessage(topic string, key, value []byte) (int64, error)
}
type KafkaAdapter struct {
	Producer       *kafka.Producer
	Consumer       *kafka.Consumer
	AdminClient    *kafka.AdminClient
	SchemaRegistry schemaregistry.Client
	ConfigMap      *kafka.ConfigMap
	Connected      bool
}

func (ka *KafkaAdapter) InitSchemaRegistry(url string) error {
	client, err := schemaregistry.NewClient(schemaregistry.NewConfig(url))
	if err != nil {
		return err
	}
	ka.SchemaRegistry = client
	return nil
}

func NewKafkaAdapter(configMap *kafka.ConfigMap) (KafkaAdapterInterface, error) {
	// Create a new admin client instance
	adminClient, err := kafka.NewAdminClient(configMap)
	if err != nil {
		return nil, err
	}

	// Create a new consumer client instance
	//todo : better handle the offset reset
	configMap.SetKey("auto.offset.reset", "earliest")
	consumer, err := kafka.NewConsumer(configMap)
	if err != nil {
		return nil, moerr.NewInternalError(context.Background(), fmt.Sprintf("unable to create confluent consumer client: %s", err))
	}

	// Create a new producer client instance
	producer, err := kafka.NewProducer(configMap)
	if err != nil {
		return nil, moerr.NewInternalError(context.Background(), fmt.Sprintf("unable to create confluent producer client: %s", err))
	}

	// Return a new KafkaAdapter instance
	return &KafkaAdapter{
		Producer:    producer,
		AdminClient: adminClient,
		Consumer:    consumer,
		ConfigMap:   configMap,
		Connected:   true,
	}, nil
}

func (ka *KafkaAdapter) GetKafkaConsumer() (*kafka.Consumer, error) {
	return ka.Consumer, nil
}
func (ka *KafkaAdapter) Close() {

	// Close the Producer if it's initialized
	if ka.Producer != nil {
		ka.Producer.Close()
	}

	// Close the Consumer if it's initialized
	if ka.Consumer != nil {
		ka.Consumer.Close()
	}

	// Close the AdminClient if it's initialized
	if ka.AdminClient != nil {
		ka.AdminClient.Close()
	}

	// Update the Connected status
	ka.Connected = false
}

func (ka *KafkaAdapter) CreateTopic(ctx context.Context, topicName string, partitions int, replicationFactor int) error {
	topicSpecification := kafka.TopicSpecification{
		Topic:             topicName,
		NumPartitions:     partitions,
		ReplicationFactor: replicationFactor,
		// can add more configs here
	}

	results, err := ka.AdminClient.CreateTopics(ctx, []kafka.TopicSpecification{topicSpecification})
	if err != nil {
		return err
	}

	// Check results for errors
	for _, result := range results {
		if result.Error.Code() != kafka.ErrNoError {
			return result.Error
		}
	}
	return nil
}

func (ka *KafkaAdapter) DescribeTopicDetails(ctx context.Context, topicName string) (*kafka.TopicMetadata, error) {

	// Fetch metadata
	meta, err := ka.AdminClient.GetMetadata(&topicName, false, int(10*time.Second.Milliseconds()))
	if err != nil {
		return nil, err
	}

	// Find and return the topic's metadata
	for _, topic := range meta.Topics {
		if topic.Topic == topicName {
			return &topic, nil
		}
	}

	return nil, moerr.NewInternalError(ctx, "topic not found")
}

func (ka *KafkaAdapter) ReadMessagesFromPartition(topic string, partition int32, offset int64, limit int) ([]*kafka.Message, error) {
	if ka.Consumer == nil {
		return nil, moerr.NewInternalError(context.Background(), "consumer not initialized")
	}

	// Assign the specific partition with the desired offset
	err := ka.Consumer.Assign([]kafka.TopicPartition{
		{Topic: &topic, Partition: partition, Offset: kafka.Offset(offset)},
	})
	if err != nil {
		return nil, moerr.NewInternalError(context.Background(), fmt.Sprintf("failed to assign partition: %s", err))
	}

	var messages []*kafka.Message
	for i := 0; i < limit; i++ {
		msg, err := ka.Consumer.ReadMessage(-1) // Wait indefinitely until a message is available
		if err != nil {
			// Check for timeout
			if kafkaErr, ok := err.(kafka.Error); ok && kafkaErr.Code() == kafka.ErrTimedOut {
				break // Exit the loop if a timeout occurs
			} else {
				return nil, moerr.NewInternalError(context.Background(), fmt.Sprintf("failed to read message: %s", err))
			}
		}
		messages = append(messages, msg)
	}

	return messages, nil
}

func (ka *KafkaAdapter) ReadMessagesFromTopic(topic string, offset int64, limit int64) ([]*kafka.Message, error) {
	if ka.Consumer == nil {
		return nil, moerr.NewInternalError(context.Background(), "consumer not initialized")
	}

	// Fetch metadata to get all partitions
	meta, err := ka.Consumer.GetMetadata(&topic, false, -1) // timeout in ms
	if err != nil {
		return nil, err
	}

	topicMetadata, ok := meta.Topics[topic]
	if !ok {
		return nil, moerr.NewInternalError(context.Background(), "topic not found in metadata")
	}

	var messages []*kafka.Message
	for _, p := range topicMetadata.Partitions {
		// Fetch the high watermark for the partition
		_, highwatermarkHigh, err := ka.Consumer.QueryWatermarkOffsets(topic, p.ID, -1)
		if err != nil {
			return nil, err
		}

		// Calculate the number of messages available to consume
		availableMessages := highwatermarkHigh - offset
		if availableMessages <= 0 {
			continue
		}

		// Determine the number of messages to consume from this partition
		partitionLimit := limit - int64(len(messages))
		if partitionLimit > availableMessages {
			partitionLimit = availableMessages
		}
		if limit == 0 {
			partitionLimit = availableMessages
		}

		// Assign the specific partition with the desired offset
		err = ka.Consumer.Assign([]kafka.TopicPartition{
			{Topic: &topic, Partition: p.ID, Offset: kafka.Offset(offset)},
		})
		if err != nil {
			return nil, err
		}

		for i := int64(0); i < partitionLimit; i++ {
			msg, err := ka.Consumer.ReadMessage(-1)
			if err != nil {
				// Check for timeout
				var kafkaErr kafka.Error
				if errors.As(err, &kafkaErr) && kafkaErr.Code() == kafka.ErrTimedOut {
					break // Exit the loop if a timeout occurs
				}
			}
			messages = append(messages, msg)
		}
	}

	return messages, nil
}

func (ka *KafkaAdapter) BatchRead(topic string, startOffset int64, limit int, batchSize int) ([]*kafka.Message, error) {
	// Fetch metadata to get all partitions
	meta, err := ka.Consumer.GetMetadata(&topic, false, 5000) // timeout in ms
	if err != nil {
		return nil, err
	}

	topicMetadata, ok := meta.Topics[topic]
	if !ok {
		return nil, moerr.NewInternalError(context.Background(), "topic not found in metadata")
	}

	numGoroutines := (limit + batchSize - 1) / batchSize

	messagesCh := make(chan []*kafka.Message, numGoroutines)
	errCh := make(chan error, numGoroutines)
	var wg sync.WaitGroup

	// Loop over each partition and start goroutines for reading
	for _, p := range topicMetadata.Partitions {
		wg.Add(1)
		go func(partition int32) {
			defer wg.Done()

			// Read a batch of messages
			messages, err := ka.ReadMessagesFromPartition(topic, partition, startOffset, batchSize)
			if err != nil {
				errCh <- err
				return
			}
			messagesCh <- messages
		}(p.ID)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	close(messagesCh)
	close(errCh)

	// Collect all messages
	var allMessages []*kafka.Message
	for batch := range messagesCh {
		allMessages = append(allMessages, batch...)
	}

	// Return the first error encountered, if any
	for err := range errCh {
		return nil, err
	}

	return allMessages, nil
}

func (ka *KafkaAdapter) GetSchemaForTopic(topic string, isKey bool) (schemaregistry.SchemaMetadata, error) {
	if ka.SchemaRegistry == nil {
		return schemaregistry.SchemaMetadata{}, moerr.NewInternalError(context.Background(), "schema registry not initialized")
	}

	subjectSuffix := "value"
	if isKey {
		subjectSuffix = "key"
	}
	subject := fmt.Sprintf("%s-%s", topic, subjectSuffix)

	// Fetch the schema for the subject
	return ka.SchemaRegistry.GetLatestSchemaMetadata(subject)
}

func (ka *KafkaAdapter) ProduceMessage(topic string, key, value []byte) (int64, error) {

	deliveryChan := make(chan kafka.Event)
	defer close(deliveryChan)

	message := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic},
		Key:            key,
		Value:          value,
	}

	err := ka.Producer.Produce(message, deliveryChan)
	if err != nil {
		return -1, moerr.NewInternalError(context.Background(), fmt.Sprintf("failed to produce message: %s", err))
	}

	e := <-deliveryChan
	m := e.(*kafka.Message)
	if m.TopicPartition.Error != nil {
		return -1, m.TopicPartition.Error
	}

	return int64(m.TopicPartition.Offset), nil
}

func newBatch(batchSize int, typs []types.Type, pool *mpool.MPool) (*batch.Batch, error) {

	//alloc space for vector
	batch := batch.NewWithSize(len(typs))
	for i, typ := range typs {
		typ.Size = int32(typ.Oid.TypeLen())
		switch typ.Oid {
		case types.T_datetime:
			typ.Scale = 6
		}
		vec := vector.NewVec(typ)
		err := vec.PreExtend(batchSize, pool)
		if err != nil {
			return nil, err
		}
		vec.SetLength(batchSize)
		batch.Vecs[i] = vec
	}
	return batch, nil
}

func populateBatchFromMSG(ctx context.Context, ka KafkaAdapterInterface, typs []types.Type, attrKeys []string, msgs []*kafka.Message, configs map[string]interface{}, mp *mpool.MPool) (*batch.Batch, error) {
	b, err := newBatch(len(msgs), typs, mp)
	if err != nil {
		return nil, err
	}
	unexpectEOF := false
	value, ok := configs[ValueKey].(string)
	if !ok {
		return nil, moerr.NewInternalError(ctx, "expected string value for key: %s", ValueKey)
	}
	switch ValueType(value) {
	case JSON:
		for i, msg := range msgs {
			err := populateOneRowData(ctx, b, attrKeys, &JsonDataGetter{Key: msg.Key, Value: msg.Value}, i, typs, mp)
			if err != nil {
				logutil.Error("populate row failed")
			}
		}
	case PROTOBUF:
		md, err := convertProtobufSchemaToMD(configs["protobuf.schema"].(string), configs["protobuf.message"].(string))
		if err != nil {
			return nil, err
		}
		for i, msg := range msgs {
			msgValue, err := deserializeProtobuf(md, msg.Value, false)
			if err != nil {
				return nil, err
			}
			err = populateOneRowData(ctx, b, attrKeys, &ProtoDataGetter{Value: msgValue, Key: msg.Key}, i, typs, mp)
			if err != nil {
				return nil, err
			}
		}

	case PROTOBUFSR:
		schema, err := ka.GetSchemaForTopic(configs["topic"].(string), false)
		if err != nil {
			return nil, err
		}
		md, err := convertProtobufSchemaToMD(schema.Schema, schema.SchemaInfo.Schema)
		if err != nil {
			return nil, err
		}
		for i, msg := range msgs {
			msgValue, _ := deserializeProtobuf(md, msg.Value, true)
			err := populateOneRowData(ctx, b, attrKeys, &ProtoDataGetter{Value: msgValue, Key: msg.Key}, i, typs, mp)
			if err != nil {
				return nil, err
			}
		}
	default:
		return nil, moerr.NewInternalError(ctx, "Unsupported value for key: %s", ValueKey)
	}

	n := b.Vecs[0].Length()
	if unexpectEOF && n > 0 {
		n--
		for i := 0; i < b.VectorCount(); i++ {
			vec := b.GetVector(int32(i))
			vec.SetLength(n)
		}
	}
	b.SetRowCount(n)
	return b, nil
}
func populateOneRowData(ctx context.Context, bat *batch.Batch, attrKeys []string, getter DataGetter, rowIdx int, typs []types.Type, mp *mpool.MPool) error {

	for colIdx, typ := range typs {
		fieldValue, ok := getter.GetFieldValue(attrKeys[colIdx])
		if !ok {
			return moerr.NewInternalError(ctx, "field not found: %s", attrKeys[colIdx])
		}

		id := typ.Oid
		vec := bat.Vecs[colIdx]
		switch id {
		case types.T_int32:
			var val int32
			switch v := fieldValue.(type) {
			case int32:
				val = v
			case float64:
				val = int32(v)
			// You can add more type cases if necessary
			default:
				return moerr.NewInternalError(ctx, "expected int32 compatible type for column %d but got %T", colIdx, fieldValue)
			}
			cols := vector.MustFixedCol[int32](vec)
			cols[rowIdx] = val
		case types.T_int64:
			var val int64
			switch v := fieldValue.(type) {
			case int64:
				val = v
			case float64:
				val = int64(v)
			// Add more type cases if necessary
			default:
				return moerr.NewInternalError(ctx, "expected int64 compatible type for column %d but got %T", colIdx, fieldValue)
			}
			cols := vector.MustFixedCol[int64](vec)
			cols[rowIdx] = val

		case types.T_uint64:
			val, ok := fieldValue.(uint64)
			if !ok {
				return moerr.NewInternalError(ctx, "expected uint64 type for column %d but got %T", colIdx, fieldValue)
			}
			cols := vector.MustFixedCol[uint64](vec)
			cols[rowIdx] = val

		case types.T_float64:
			val, ok := fieldValue.(float64)
			if !ok {
				return moerr.NewInternalError(ctx, "expected float64 type for column %d but got %T", colIdx, fieldValue)
			}
			cols := vector.MustFixedCol[float64](vec)
			cols[rowIdx] = val

		case types.T_char, types.T_varchar, types.T_binary, types.T_varbinary, types.T_blob, types.T_text:
			val, ok := fieldValue.(string)
			if !ok {
				return moerr.NewInternalError(ctx, "expected string type for column %d but got %T", colIdx, fieldValue)
			}
			err := vector.SetStringAt(vec, rowIdx, val, mp)
			if err != nil {
				return err
			}

		case types.T_bool:
			val, ok := fieldValue.(bool)
			if !ok {
				return moerr.NewInternalError(ctx, "expected bool type for column %d but got %T", colIdx, fieldValue)
			}
			cols := vector.MustFixedCol[bool](vec)
			cols[rowIdx] = val

		case types.T_json:
			val, ok := fieldValue.([]byte)
			if !ok || len(val) == 0 {
				strVal, strOk := fieldValue.(string)
				if !strOk {
					return moerr.NewInternalError(ctx, "expected bytes or string type for JSON column %d but got %T", colIdx, fieldValue)
				}
				val = []byte(strVal)
			}
			err := vector.SetBytesAt(vec, rowIdx, val, mp)
			if err != nil {
				return err
			}

		case types.T_datetime:
			val, ok := fieldValue.(string)
			if !ok {
				return moerr.NewInternalError(ctx, "expected string type for Datetime column %d but got %T", colIdx, fieldValue)
			}
			cols := vector.MustFixedCol[types.Datetime](vec)
			if len(val) == 0 {
				cols[rowIdx] = types.Datetime(0)
			} else {
				d, err := types.ParseDatetime(val, vec.GetType().Scale)
				if err != nil {
					return moerr.NewInternalError(ctx, "the input value is not Datetime type for column %d: %v", colIdx, fieldValue)
				}
				cols[rowIdx] = d
			}

		default:
			return moerr.NewInternalError(ctx, "the value type %s is not supported now", *vec.GetType())
		}
	}
	return nil
}

func convertProtobufSchemaToMD(schema string, msgTypeName string) (*desc.MessageDescriptor, error) {
	files := map[string]string{
		"test.proto": schema,
	}

	parser := protoparse.Parser{
		Accessor: protoparse.FileContentsFromMap(files),
	}
	fds, err := parser.ParseFiles("test.proto")

	if err != nil {
		return nil, err
	}
	fd := fds[0]
	md := fd.FindMessage(msgTypeName)
	return md, nil
}

func deserializeProtobuf(md *desc.MessageDescriptor, in []byte, isKafkSR bool) (*dynamic.Message, error) {
	dm := dynamic.NewMessage(md)
	var err error
	if isKafkSR {
		bytesRead, _, err := readMessageIndexes(in[5:])
		if err != nil {
			return nil, err
		}
		proto.Unmarshal(in[5+bytesRead:], dm)
	} else {
		err = dm.Unmarshal(in)
	}
	return dm, err
}

func readMessageIndexes(payload []byte) (int, []int, error) {
	arrayLen, bytesRead := binary.Varint(payload)
	if bytesRead <= 0 {
		return bytesRead, nil, moerr.NewInternalError(context.Background(), "unable to read message indexes")
	}
	if arrayLen == 0 {
		// Handle the optimization for the first message in the schema
		return bytesRead, []int{0}, nil
	}
	msgIndexes := make([]int, arrayLen)
	for i := 0; i < int(arrayLen); i++ {
		idx, read := binary.Varint(payload[bytesRead:])
		if read <= 0 {
			return bytesRead, nil, moerr.NewInternalError(context.Background(), "unable to read message indexes")
		}
		bytesRead += read
		msgIndexes[i] = int(idx)
	}
	return bytesRead, msgIndexes, nil
}

func convertToKafkaConfig(configs map[string]interface{}) *kafka.ConfigMap {
	kafkaConfigs := &kafka.ConfigMap{}
	allowedKeys := map[string]struct{}{
		"bootstrap.servers": {},
		"security.protocol": {},
		"sasl.mechanisms":   {},
		"sasl.username":     {},
		"sasl.password":     {},
		// Add other Kafka-specific properties here...
	}

	for key, value := range configs {
		if _, ok := allowedKeys[key]; ok {
			kafkaConfigs.SetKey(key, value)
		}
	}
	groupId := uuid.New().String()
	kafkaConfigs.SetKey("group.id", groupId)

	return kafkaConfigs
}

func ValidateConfig(ctx context.Context, configs map[string]interface{}, factory func(configMap *kafka.ConfigMap) (KafkaAdapterInterface, error)) error {
	var requiredKeys = []string{
		TypeKey,
		TopicKey,
		ValueKey,
		BootstrapServersKey,
	}

	for _, key := range requiredKeys {
		if _, exists := configs[key]; !exists {
			return moerr.NewInternalError(ctx, "missing required key: %s", key)
		}
	}

	value, ok := configs[ValueKey].(string)
	if !ok {
		return moerr.NewInternalError(ctx, "expected string value for key: %s", ValueKey)
	}

	switch ValueType(value) {
	case JSON:
		// no additional checks required
	case PROTOBUF:
		// check the schema and message name has been set or not
		if _, ok := configs[ProtobufSchemaKey]; !ok {
			return moerr.NewInternalError(ctx, "missing required key: %s", ProtobufSchemaKey)
		}
		if _, ok := configs[ProtobufMessagekey]; !ok {
			return moerr.NewInternalError(ctx, "missing required key: %s", ProtobufMessagekey)
		}
	case PROTOBUFSR:
		if _, ok := configs[ProtobufMessagekey]; !ok {
			return moerr.NewInternalError(ctx, "missing required key: %s", ProtobufMessagekey)
		}
		if _, ok := configs[SchemaRegistryKey]; !ok {
			return moerr.NewInternalError(ctx, "missing required key: %s", SchemaRegistryKey)
		}
	default:
		return moerr.NewInternalError(ctx, "Unsupported value for key: %s", ValueKey)
	}
	// Convert the configuration to map[string]string for Kafka
	kafkaConfigs := convertToKafkaConfig(configs)

	// Create the Kafka adapter
	ka, err := factory(kafkaConfigs)
	if err != nil {
		return err
	}
	defer ka.Close()

	// Check if Topic exists
	_, err = ka.DescribeTopicDetails(ctx, configs[TopicKey].(string))
	if err != nil {
		return err
	}
	return nil
}

type KafkaAdapterFactory func(configMap *kafka.ConfigMap) (KafkaAdapterInterface, error)

func GetStreamCurrentSize(ctx context.Context, configs map[string]interface{}, factory KafkaAdapterFactory) (int64, error) {
	err := ValidateConfig(ctx, configs, NewKafkaAdapter)
	if err != nil {
		return 0, err
	}

	configMap := convertToKafkaConfig(configs)

	ka, err := factory(configMap)
	if err != nil {
		return 0, err
	}
	defer ka.Close()

	meta, err := ka.DescribeTopicDetails(ctx, configs["topic"].(string))
	if err != nil {
		return 0, err
	}

	var totalSize int64
	kaConsumer, _ := ka.GetKafkaConsumer()
	for _, p := range meta.Partitions {
		// Fetch the high watermark for the partition
		_, highwatermarkHigh, err := kaConsumer.QueryWatermarkOffsets(configs["topic"].(string), p.ID, -1)
		if err != nil {
			return 0, err
		}
		totalSize += int64(highwatermarkHigh)
	}
	return totalSize, nil
}

func RetrieveData(ctx context.Context, configs map[string]interface{}, attrs []string, types []types.Type, offset int64, limit int64, mp *mpool.MPool, factory KafkaAdapterFactory) (*batch.Batch, error) {
	err := ValidateConfig(ctx, configs, NewKafkaAdapter)
	if err != nil {
		return nil, err
	}

	configMap := convertToKafkaConfig(configs)

	ka, err := factory(configMap)
	if err != nil {
		return nil, err
	}
	defer ka.Close()

	// init schema registry client if schema registry url is set
	if sr, ok := configs[SchemaRegistryKey]; ok {
		err = ka.InitSchemaRegistry(sr.(string))
		if err != nil {
			return nil, err
		}
	}

	messages, err := ka.ReadMessagesFromTopic(configs["topic"].(string), offset, limit)
	if err != nil {
		return nil, err
	}

	b, err := populateBatchFromMSG(ctx, ka, types, attrs, messages, configs, mp)
	if err != nil {
		return nil, err
	}

	return b, nil
}
