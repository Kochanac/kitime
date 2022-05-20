package main

import (
	"fmt"
	"github.com/bojand/ghz/printer"
	"github.com/bojand/ghz/runner"
	"github.com/jhump/protoreflect/desc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"math"
	"math/rand"
	"os"
	pb "perf_test/internal/api"
	"time"
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
	var rps uint = 8000
	report, err := runner.Run(
		"head.head.Set",
		"185.46.8.166:80",
		runner.WithProtoFile("../service/api/api.proto", []string{}),
		runner.WithInsecure(true),
		runner.WithBinaryDataFunc(dataFunc),
		runner.WithTotalRequests(rps*2000),
		runner.WithRPS(rps),
		runner.WithAsync(true),
		runner.WithConnections(8),
		runner.WithLoadStart(rps),
		runner.WithLoadEnd(rps*2),
		runner.WithLoadStepDuration(time.Second*30),
		runner.WithLoadStep(int(math.Floor(float64(rps)/5))),
		runner.WithLoadSchedule("step"),
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
