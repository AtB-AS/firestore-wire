package main

import (
	"encoding/json"
	"fmt"
	"os"

	pb "google.golang.org/genproto/googleapis/firestore/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

var (
	nullValue = &pb.Value{ValueType: &pb.Value_NullValue{}}
)

func main() {
	var m map[string]interface{}

	d := json.NewDecoder(os.Stdin)
	d.UseNumber()
	if err := d.Decode(&m); err != nil {
		panic(err)
	}

	fmt.Fprintln(os.Stdout, mustMarshal(mustFromJSON(m)))
}

func mustMarshal(d *pb.Document) string {
	b, err := protojson.Marshal(d)
	if err != nil {
		panic(err)
	}

	return string(b)
}

func mustFromJSON(m map[string]interface{}) *pb.Document {
	j, err := fromJSON(m)
	if err != nil {
		panic(err)
	}

	return j
}

func fromJSON(m map[string]interface{}) (*pb.Document, error) {
	d := pb.Document{Fields: make(map[string]*pb.Value, len(m))}
	for k, v := range m {
		pv, err := jsonToProtoValue(v)
		if err != nil {
			return nil, err
		}

		d.Fields[k] = pv
	}

	return &d, nil
}

func jsonToProtoValue(v interface{}) (pbv *pb.Value, err error) {
	switch vt := v.(type) {
	case nil:
		return nullValue, nil
	case string:
		return &pb.Value{ValueType: &pb.Value_StringValue{vt}}, nil
	case json.Number:
		if n, err := vt.Int64(); err == nil {
			return &pb.Value{ValueType: &pb.Value_IntegerValue{n}}, nil
		} else if f, err := vt.Float64(); err == nil {
			return &pb.Value{ValueType: &pb.Value_DoubleValue{f}}, nil
		} else {
			return &pb.Value{ValueType: &pb.Value_StringValue{vt.String()}}, nil
		}
	case bool:
		return &pb.Value{ValueType: &pb.Value_BooleanValue{vt}}, nil
	case []interface{}:
		return sliceToProtoValue(vt)
	case map[string]interface{}:
		return mapToProtoValue(vt)
	default:
		panic(fmt.Sprintf("cannot convert type: %T", v))
	}
}

func sliceToProtoValue(s []interface{}) (*pb.Value, error) {
	vals := make([]*pb.Value, len(s))
	for i, v := range s {
		val, err := jsonToProtoValue(v)
		if err != nil {
			return nil, err
		}
		vals[i] = val
	}
	return &pb.Value{ValueType: &pb.Value_ArrayValue{&pb.ArrayValue{Values: vals}}}, nil
}

func mapToProtoValue(m map[string]interface{}) (*pb.Value, error) {
	pm := map[string]*pb.Value{}
	for k, v := range m {
		val, err := jsonToProtoValue(v)
		if err != nil {
			return nil, err
		}
		pm[k] = val
	}

	return &pb.Value{ValueType: &pb.Value_MapValue{&pb.MapValue{Fields: pm}}}, nil
}
