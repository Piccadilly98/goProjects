package main

import (
	"fmt"
	"runtime"
	"time"

	dto "github.com/Piccadilly98/goProjects/intellectHome2.0/src/DTO"
	"github.com/Piccadilly98/goProjects/intellectHome2.0/src/core/alerts"
	"github.com/Piccadilly98/goProjects/intellectHome2.0/src/core/alerts/notifiers"
	"github.com/Piccadilly98/goProjects/intellectHome2.0/src/core/alerts/rules"
	board_info_rules "github.com/Piccadilly98/goProjects/intellectHome2.0/src/core/alerts/rules/boardInfoRules"
	board_status_rules "github.com/Piccadilly98/goProjects/intellectHome2.0/src/core/alerts/rules/boardStatusRules"
	database_rules "github.com/Piccadilly98/goProjects/intellectHome2.0/src/core/alerts/rules/dataBaseRules"
	"github.com/Piccadilly98/goProjects/intellectHome2.0/src/core/events"
	_ "github.com/lib/pq"
)

func main() {
	// serv, err := server.NewServer(false, 30*time.Second, 150*time.Second, true, 10, 10*time.Second)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// serv.Start("localhost:8080")
	// err = <-serv.ErrorServerChan
	// log.Fatal(err)
	check := database_rules.NewDataBaseChecker(true, database_rules.NewDataBaseStatusChecker(true, true, true, true))
	rssi := board_info_rules.NewRSsiChecker(0, 0)
	voltage := board_info_rules.NewVoltageChecker(0, 0, 0, 0)
	temp := board_info_rules.NewTemperatureCpuCheck()
	dto := &dto.UpdateBoardInfo{
		RssiWifi: getPtrInt(60),
		CpuTemp:  getPtrFloat(300),
		Voltage:  getPtrFloat(3.1),
	}
	bch := board_status_rules.NewBoardStatusChecker(true, true, true)
	ch := board_info_rules.NewBoardInfoChecker(rssi, temp, voltage)
	eb := events.NewEventBus(50, 10*time.Second)
	al := alerts.NewAlertsManager(eb, []rules.Rule{bch, ch, check}, 50, &notifiers.LogNotifier{})
	al.Start()

	sub := eb.Subscribe(alerts.TopicForBoardStatusChecker, "main")
	sub2 := eb.Subscribe(alerts.TopicforBoardInfoChecker, "main")
	sub3 := eb.Subscribe(events.TopicErrorsDB, "main")
	sub4 := eb.Subscribe(events.TopicDataBaseStatus, "main")
	eb.Publish(sub.Topic, events.Event{
		BoardID: "esp32_1",
		Payload: "offline",
	}, sub.ID)
	time.Sleep(2 * time.Second)
	eb.Publish(sub2.Topic, events.Event{
		BoardID: "esp32_2",
		Payload: dto,
	}, sub2.ID)
	time.Sleep(1 * time.Second)
	eb.Publish(events.TopicErrorsDB, events.Event{
		Payload: fmt.Errorf("errordb"),
	}, sub3.ID)

	time.Sleep(2 * time.Second)
	eb.Publish(sub4.Topic, events.Event{
		Payload: fmt.Sprintf("DataBase fail, start Recover"),
	}, sub4.ID)
	fmt.Println(runtime.NumGoroutine())
	time.Sleep(10 * time.Second)

}

func getPtrFloat(f float64) *float64 {
	return &f
}

func getPtrInt(i int) *int {
	return &i
}
