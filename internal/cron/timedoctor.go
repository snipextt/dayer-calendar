package cron

import (
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/snipextt/dayer/internal/timedoctor"
	"github.com/snipextt/dayer/models"
	"github.com/snipextt/dayer/models/connection"
	"github.com/snipextt/dayer/models/workspace"
	"github.com/snipextt/dayer/storage"
	"github.com/snipextt/dayer/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func pushToKafka(report models.TimeDoctorReportForAnalysis) {
	ctx, cancel := utils.NewContext()
	defer cancel()
	defer catchError()
	reportBytes, err := report.ToBytes()
	utils.CheckError(err)
	err = storage.KafkaWriter().WriteMessages(ctx, kafka.Message{
		Value: reportBytes,
	})
	utils.CheckError(err)
}

func timedoctorCron() {
	defer catchError()
	connections, err := connection.GetTimedoctorConnections()
	utils.CheckError(err)
	var wg sync.WaitGroup
	for _, connection := range connections {
		wg.Add(1)
		go generateTimedoctorReportConnection(connection, &wg)
	}
}

func generateTimedoctorReportConnection(connection connection.Model, wg *sync.WaitGroup) {
	defer wg.Done()
	defer catchError()

	users, err := workspace.FindWorkspaceMembers(connection.Workspace.(primitive.ObjectID))
	utils.CheckError(err)

	var wg2 sync.WaitGroup

	for _, user := range users {
		if user.Meta.TimeDoctorId == "" {
			continue
		}
		wg2.Add(1)
		go generateTimedoctorReport(connection, user, &wg2)
	}
}

func generateTimedoctorReport(connection connection.Model, user workspace.Member, wg *sync.WaitGroup) {
	defer wg.Done()
	defer catchError()

	location, err := time.LoadLocation("Asia/Kolkata")
	utils.CheckError(err)

	date := time.Now().In(location).AddDate(0, 0, -1)

	report, err := timedoctor.GenerateReportFromTimedoctor(connection.Token, connection.Meta.TimeDoctorCompanyID, user.Meta.TimeDoctorId, connection.Meta.TimeDoctorParseScreencast, date)
	utils.CheckError(err)

	report.MemberId = user.Id
	report.WorkspaceId = connection.Workspace.(primitive.ObjectID)

	pushToKafka(report)
}
