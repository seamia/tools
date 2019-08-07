package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/hashicorp/hcl"
)

type msi = map[string]interface{}
type slice = []interface{}

const (
	defaultCapacityUnits      = 10
	defaultReadCapacityUnits  = defaultCapacityUnits
	defaultWriteCapacityUnits = defaultCapacityUnits

	tableMarker     = "aws_dynamodb_table"
	failureExitCode = 7

	usage = `This utility converts terraform definition of AWS DynamoDB table into JSON-ed representation of dynamodb.CreateTableInput, which can be used to (re-)create dynamodb tables...
the usage:
    dynamoSchema Name-of-the-terraform-file.tf > name-of-the-output-file.json

(disclaimer: there were a few corners cut during the development of this utility - there is a chance it may fail ingesting certain input files)
`
)

func println(a ...interface{}) {
	_, _ = fmt.Fprintln(os.Stderr, a...)
}

func main() {
	input := ""
	if len(os.Args) > 1 {
		input = os.ExpandEnv(os.Args[1])
	}

	if len(input) == 0 || input == "?" || input == "/?" || input == "-?" || input == "help" {
		println(usage)
		os.Exit(failureExitCode)
	}

	content, err := ioutil.ReadFile(input)
	if err != nil {
		println("Failed to load file [", input, "] due to error:", err)
		os.Exit(failureExitCode)
	}

	what := make(msi)
	err = hcl.Unmarshal(content, &what)
	if err != nil {
		println("Failed to process file [", input, "] due to error:", err)
		os.Exit(failureExitCode)
	}
	ingest(what)
}

// slice
func getSlice(raw interface{}, name string) ([]msi, bool) {
	if raw != nil {
		if data, converts := raw.(msi); converts {
			if branch, present := data[name]; present {
				if leaf, converts := branch.([]msi); converts {
					return leaf, true
				}
			}
		}
	}
	return nil, false
}

func ingest(data msi) {
	if resources, present := getSlice(data, "resource"); present {
		for _, entry := range resources {
			for key, value := range entry {
				println("=======================", key)
				_ = value

				if key == tableMarker {
					if tables, present := getSlice(entry, key); present {
						for _, table := range tables {
							create(table)
						}
					}
				}
			}
		}
	}
}

func add_attributes(from []msi, to *dynamodb.CreateTableInput) {
	if len(from) > 0 {
		to.AttributeDefinitions = []*dynamodb.AttributeDefinition{}

		for _, data := range from {
			attribute := dynamodb.AttributeDefinition{}
			for k, v := range data {
				txt := v.(string)
				switch k {
				case "name":
					attribute.SetAttributeName(txt)
				case "type":
					attribute.SetAttributeType(txt)
				default:
					println("----------- ignoring:", k)
				}
			}
			to.AttributeDefinitions = append(to.AttributeDefinitions, &attribute)
		}
	}
}

func add_global_secondary_index(from []msi, to *dynamodb.CreateTableInput) {

	if len(from) > 0 {
		to.GlobalSecondaryIndexes = []*dynamodb.GlobalSecondaryIndex{}

		for _, data := range from {
			index := dynamodb.GlobalSecondaryIndex{}
			keys := make([]*dynamodb.KeySchemaElement, 0)

			throughput := dynamodb.ProvisionedThroughput{}
			throughput.SetReadCapacityUnits(defaultReadCapacityUnits)
			throughput.SetWriteCapacityUnits(defaultWriteCapacityUnits)

			for k, v := range data {
				txt := v.(string)
				switch k {
				case "name":
					index.SetIndexName(txt)

				case "hash_key":
					element := dynamodb.KeySchemaElement{}
					element.SetAttributeName(txt)
					element.SetKeyType("HASH")
					keys = append(keys, &element)

				case "range_key":
					element := dynamodb.KeySchemaElement{}
					element.SetAttributeName(txt)
					element.SetKeyType("RANGE")
					keys = append(keys, &element)

				case "projection_type":
					projection := dynamodb.Projection{}
					projection.SetProjectionType(txt)
					index.SetProjection(&projection)

				case "read_capacity":
				case "write_capacity":
					// these are ignored

				default:
					println("----------- ignoring:", k)
				}
			}
			index.SetKeySchema(keys)

			// ProvisionedThroughput: &dynamodb.ProvisionedThroughput{ReadCapacityUnits: aws.Int64(50), WriteCapacityUnits: aws.Int64(50)},
			index.SetProvisionedThroughput(&throughput)

			to.GlobalSecondaryIndexes = append(to.GlobalSecondaryIndexes, &index)
		}
	}
}

func create(from msi) (*dynamodb.CreateTableInput, error) {
	if len(from) != 1 {
		println("wrong count", from)
	}
	// println(from)

	result := &dynamodb.CreateTableInput{
		AttributeDefinitions:   []*dynamodb.AttributeDefinition{},
		KeySchema:              []*dynamodb.KeySchemaElement{},
		GlobalSecondaryIndexes: []*dynamodb.GlobalSecondaryIndex{},
	}

	throughput := dynamodb.ProvisionedThroughput{}
	throughput.SetReadCapacityUnits(defaultReadCapacityUnits)
	throughput.SetWriteCapacityUnits(defaultWriteCapacityUnits)

	for name, value := range from {
		_ = name
		if definitions, converts := value.([]msi); converts {
			if len(definitions) != 1 {
				println("more than 1 definition ?????????")
			}
			for _, definition := range definitions {

				for property, value := range definition {
					switch property {
					case "name":
						result.TableName = aws.String(value.(string))

					case "hash_key":
						index := dynamodb.KeySchemaElement{
							AttributeName: aws.String(value.(string)),
							KeyType:       aws.String("HASH"),
						}
						result.KeySchema = append(result.KeySchema, &index)

					case "range_key":
						index := dynamodb.KeySchemaElement{
							AttributeName: aws.String(value.(string)),
							KeyType:       aws.String("RANGE"),
						}
						result.KeySchema = append(result.KeySchema, &index)

					case "attribute":
						if attributes, converts := value.([]msi); converts {
							add_attributes(attributes, result)
						} else {
							println("====== no conversion:", value)
						}

					case "global_secondary_index":
						if attributes, converts := value.([]msi); converts {
							add_global_secondary_index(attributes, result)
						} else {
							println("====== no conversion:", value)
						}

					case "write_capacity", "read_capacity", "point_in_time_recovery", "tags", "lifecycle":
						// ignoring these (for now)

					default:
						println("====== unhandled property:", property)
					}
				}
			}
		} else {
			println("no conversion: ", value)
		}
	}

	result.SetProvisionedThroughput(&throughput)

	reArrangeFields(result)

	if result != nil {
		if bytes, err := json.MarshalIndent(result, "", "\t"); err == nil {
			os.Stdout.Write(bytes)
		}
	}

	return result, nil
}

// Q: why do we need this?
// A: aws wants "hash" elements to come before "range" ones (in the list)... go figure...
func reArrangeFields(result *dynamodb.CreateTableInput) {
	if result == nil {
		return
	}

	// KeySchema: []*dynamodb.KeySchemaElement{},
	sort.Slice(result.KeySchema, func(i, j int) bool {
		return *result.KeySchema[i].KeyType < *result.KeySchema[j].KeyType
	})

	// GlobalSecondaryIndexes: []*dynamodb.GlobalSecondaryIndex{},
	if len(result.GlobalSecondaryIndexes) > 0 {
		for _, index := range result.GlobalSecondaryIndexes {
			sort.Slice(index.KeySchema, func(i, j int) bool {
				return *index.KeySchema[i].KeyType < *index.KeySchema[j].KeyType
			})
		}
	}
}
