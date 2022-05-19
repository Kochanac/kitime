package main

import (
	"fmt"
	"github.com/bojand/ghz/printer"
	"github.com/bojand/ghz/runner"
	"github.com/jhump/protoreflect/desc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"math/rand"
	"os"
	pb "perf_test/internal/api"
)

func dataFunc(mtd *desc.MethodDescriptor, cd *runner.CallData) []byte {
	msg := &pb.SetRequest{
		UserId:    uint32(rand.Intn(1000)),
		EventTime: timestamppb.Now(),
		EventType: pb.SetRequest_EVENT_TYPE(rand.Intn(2)),
		VideoId:   uint32(rand.Intn(1000)),
		VideoTime: uint32(rand.Intn(1000)),
	}
	binData, _ := proto.Marshal(msg)
	return binData
}

func main() {
	var rps uint = 16000
	report, err := runner.Run(
		"head.head.Set",
		"151.248.117.121:80",
		runner.WithProtoFile("../service/api/api.proto", []string{}),
		runner.WithInsecure(true),
		runner.WithBinaryDataFunc(dataFunc),
		runner.WithTotalRequests(rps*200),
		runner.WithRPS(rps),
		runner.WithAsync(true),
	)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	reportPrinter := printer.ReportPrinter{
		Out:    os.Stdout,
		Report: report,
	}

	_ = reportPrinter.Print("summary")
}
