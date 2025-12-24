package main

import (
	"fmt"
	"log"
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
	stChecker, _ := database_rules.NewDataBaseStatusChecker(true, true, true, true)
	errChecker, _ := database_rules.NewErrorDBChecker(true, true, true, nil, nil, rules.TypeAlertWarning)
	checkDb, _ := database_rules.NewDataBaseChecker(stChecker, errChecker)
	checkerBoard := board_info_rules.NewBoardInfoChecker(board_info_rules.NewRSsiChecker(0, 0), board_info_rules.NewTemperatureCpuCheck(), board_info_rules.NewVoltageChecker(0, 0, 0, 0))
	checkerStatus, _ := board_status_rules.NewBoardStatusChecker(true, true, true)
	bus := events.NewEventBus(10, 5*time.Second)
	loger := notifiers.LogNotifier{}
	am, err := alerts.NewAlertsManager(bus, []rules.Rule{checkDb, checkerBoard, checkerStatus}, 0, 10, &loger)
	if err != nil {
		log.Fatal(err)
	}
	am.Start()

	time.Sleep(2 * time.Second)

	sub := bus.Subscribe(events.TopicBoardInfoUpdateDTO, "main")
	bus.Publish(sub.Topic, events.Event{
		Payload: &dto.UpdateBoardInfo{
			CpuTemp:  getPtrFloat(100),
			RssiWifi: getPtrInt(-100),
			Voltage:  getPtrFloat(4),
		},
	}, sub.ID)

	sub1 := bus.Subscribe(events.TopicBoardsStatusUpdate, "main")
	log.Println(sub1.ID)
	err = bus.Publish(sub1.Topic, events.Event{
		BoardID: "esp32_2",
		Payload: "update status to offline",
	}, sub1.ID)
	if err != nil {
		log.Fatal(err.Error())
	}

	sub2 := bus.Subscribe(events.TopicErrorsDB, "main")
	bus.Publish(sub2.Topic, events.Event{
		Payload: fmt.Errorf("conection fail"),
	}, sub2.ID)
	sub3 := bus.Subscribe(events.TopicDataBaseStatus, "main")
	bus.Publish(events.TopicDataBaseStatus, events.Event{
		Payload: database_rules.DataBaseFail,
	}, sub3.ID)

	time.Sleep(10 * time.Second)
	// serv, err := server.NewServer(false, 30*time.Second, 150*time.Second, true, 10, 10*time.Second)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// serv.Start("localhost:8080")
	// err = <-serv.ErrorServerChan
	// log.Fatal(err)
	// check := database_rules.NewDataBaseChecker(true, database_rules.NewDataBaseStatusChecker(true, true, true, true))
	// rssi := board_info_rules.NewRSsiChecker(0, 0)
	// voltage := board_info_rules.NewVoltageChecker(0, 0, 0, 0)
	// temp := board_info_rules.NewTemperatureCpuCheck()
	// dto := &dto.UpdateBoardInfo{
	// 	RssiWifi: getPtrInt(60),
	// 	CpuTemp:  getPtrFloat(300),
	// 	Voltage:  getPtrFloat(3.1),
	// }
	// bch := board_status_rules.NewBoardStatusChecker(true, true, true)
	// ch := board_info_rules.NewBoardInfoChecker(rssi, temp, voltage)
	// eb := events.NewEventBus(50, 10*time.Second)
	// al := alerts.NewAlertsManager(eb, []rules.Rule{bch, ch, check}, 50, &notifiers.LogNotifier{})
	// al.Start()

	// sub := eb.Subscribe(alerts.TopicForBoardStatusChecker, "main")
	// sub2 := eb.Subscribe(alerts.TopicforBoardInfoChecker, "main")
	// sub3 := eb.Subscribe(events.TopicErrorsDB, "main")
	// sub4 := eb.Subscribe(events.TopicDataBaseStatus, "main")
	// eb.Publish(sub.Topic, events.Event{
	// 	BoardID: "esp32_1",
	// 	Payload: "offline",
	// }, sub.ID)
	// time.Sleep(2 * time.Second)
	// eb.Publish(sub2.Topic, events.Event{
	// 	BoardID: "esp32_2",
	// 	Payload: dto,
	// }, sub2.ID)
	// time.Sleep(1 * time.Second)
	// eb.Publish(events.TopicErrorsDB, events.Event{
	// 	Payload: fmt.Errorf("errordb"),
	// }, sub3.ID)

	// time.Sleep(2 * time.Second)
	// eb.Publish(sub4.Topic, events.Event{
	// 	Payload: fmt.Sprintf("DataBase fail, start Recover"),
	// }, sub4.ID)
	// fmt.Println(runtime.NumGoroutine())
	// time.Sleep(10 * time.Second)

}

func getPtrFloat(f float64) *float64 {
	return &f
}

func getPtrInt(i int) *int {
	return &i
}
