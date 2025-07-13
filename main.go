package main

import (
	"datalogger/databases"
	"datalogger/models"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/goburrow/modbus"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
	logrus "github.com/sirupsen/logrus"
)

func initLogger() {
	// Create log file
	file, err := os.OpenFile("plc.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("❌ Failed to open log file:", err)
		os.Exit(1)
	}

	logrus.SetOutput(io.MultiWriter(os.Stdout, file))
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logrus.SetLevel(logrus.InfoLevel)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		logrus.Fatal("❌ Error loading .env file")
	}
	initLogger()
	logrus.Info("🚀 PLC Master Reader Started")
	err = databases.InitDb()
	if err != nil {
		logrus.Fatalf("❌ Error loading database : %v", err)

	}
	c := cron.New(cron.WithSeconds())
	_, err = c.AddFunc("*/30 * * * * *", func() {
		logrus.Println("⏱ Task triggered at:", time.Now().Format("15:04:05"))
		if err != nil {
			logrus.Errorf("Error Machine Data Logging : %v", err)
		}
		machinesDataLog()
	})
	if err != nil {
		logrus.Error(err)
		panic("❌ Failed to schedule job: " + err.Error())
	}
	c.Start()
	select {}
}

const (
	// MASTERS
	CountAddr  = 40361
	REC0Pulse  = 40362
	REC1Pulse  = 40363
	ResetPulse = 40364
	REC0_ADDR  = 40001
	REC1_ADDR  = 40015
	RECORD_LEN = 15
)

func readUint16(client modbus.Client, address uint16) (uint16, error) {
	results, err := client.ReadHoldingRegisters(address-40001, 1)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint16(results), nil
}

func readRecord(client modbus.Client, startAddress uint16, length uint16) ([]uint16, error) {
	results, err := client.ReadHoldingRegisters(startAddress-40001, length)
	if err != nil {
		return nil, err
	}
	data := make([]uint16, length)
	for i := 0; i < int(length); i++ {
		data[i] = binary.BigEndian.Uint16(results[i*2 : (i+1)*2])
	}
	return data, nil
}

func writePulse(client modbus.Client, address uint16, value uint16) error {
	_, err := client.WriteSingleCoil(address-40001, value)
	return err
}

func machinesDataLog() error {
	var machines []models.Machine
	err := databases.DiantaDB.Find(&machines).Error
	if err != nil {
		return err
	}

	for i, machine := range machines {
		logrus.Printf("Machine Processed: %d", i)
		handler := modbus.NewTCPClientHandler(fmt.Sprintf("%s:%s", machine.MachineIP, machine.MachinePort))
		handler.SlaveId = 1
		handler.Timeout = 5 * time.Second
		client := modbus.NewClient(handler)

		if err := handler.Connect(); err != nil {
			logrus.Errorf("🔌 Failed to connect to PLC: %v", err)
			continue
		}
		defer handler.Close()

		count, err := readUint16(client, CountAddr)
		if err != nil {
			logrus.Errorf("❌ Failed to read count: %v", err)
			continue
		}
		logrus.Infof("📦 Record Count: %d", count)

		if count == 0 {
			continue
		}

		// Process REC0 first
		rec0Read := 0
		rec0ReadInt, err := readUint16(client, REC0Pulse)
		if err != nil {
			logrus.Errorf("❌ REC0 pulse read error: %v", err)
			continue
		}
		rec0Read = int(rec0ReadInt + 1)
		for rec0Read <= 242 && rec0Read < int(count) {
			err = writePulse(client, REC0Pulse, uint16(rec0Read))
			if err != nil {
				logrus.Errorf("❌ REC0 pulse write error: %v", err)
				break
			}
			time.Sleep(300 * time.Millisecond)

			rec0, err := readRecord(client, REC0_ADDR, RECORD_LEN)
			if err != nil {
				logrus.Errorf("❌ Error reading REC0: %v", err)
				err = writePulse(client, REC0Pulse, uint16(rec0Read-1))
				if err != nil {
					logrus.Errorf("❌ REC0 pulse write error: %v", err)
					break
				}
				break
			}
			if !isRecordEmpty(rec0) {
				if _, err := storeMasterDB(rec0, int(machine.ID)); err != nil {
					logrus.Errorf("Error storing REC0 to DB: %v", err)
					err = writePulse(client, REC0Pulse, uint16(rec0Read-1))
					if err != nil {
						logrus.Errorf("❌ REC0 pulse write error: %v", err)
						break
					}
					break
				}
				logrus.Infof("✅ REC0 record %d stored successfully", rec0Read+1)
			} else {
				logrus.Warnf("⚠️ REC0 record %d skipped: Empty timestamp", rec0Read+1)
			}
			rec0Read++
		}

		// Process REC1
		rec1Read := 0
		rec1ReadInt, err := readUint16(client, REC1Pulse)
		if err != nil {
			logrus.Errorf("❌ REC1 pulse read error: %v", err)
			continue
		}
		rec1Read = int(rec1ReadInt + 1)
		remaining := int(count) - rec0Read
		for rec1Read <= 242 && rec1Read < remaining {
			err = writePulse(client, REC1Pulse, uint16(rec1Read))
			if err != nil {
				logrus.Errorf("❌ REC1 pulse write error: %v", err)
				break
			}
			time.Sleep(300 * time.Millisecond)
			rec1, err := readRecord(client, REC1_ADDR, RECORD_LEN)
			if err != nil {
				logrus.Errorf("❌ Error reading REC1: %v", err)
				err = writePulse(client, REC1Pulse, uint16(rec1Read-1))
				if err != nil {
					logrus.Errorf("❌ REC1 pulse write error: %v", err)
					break
				}
				break
			}
			if !isRecordEmpty(rec1) {
				if _, err := storeMasterDB(rec1, int(machine.ID)); err != nil {
					logrus.Errorf("Error storing REC1 to DB: %v", err)
					err = writePulse(client, REC0Pulse, uint16(rec1Read-1))
					if err != nil {
						logrus.Errorf("❌ REC0 pulse write error: %v", err)
						break
					}
					break
				}
				logrus.Infof("✅ REC1 record %d stored successfully", rec1Read+1)
			} else {
				logrus.Warnf("⚠️ REC1 record %d skipped: Empty timestamp", rec1Read+1)
			}
			rec1Read++
		}
		if err != nil {
			logrus.Errorf("❌ Failed to reset PLC data: %v", err)
		} else {
			logrus.Info("🔄 PLC data reset complete")
		}
	}
	return nil
}

func isRecordEmpty(data []uint16) bool {
	return data[0] == 0 && data[1] == 0 && data[2] == 0
}

func storeMasterDB(rec0 []uint16, machineMasterId int) (*models.MasterData, error) {
	logrus.Printf("✅ Stored Record: %04d-%02d-%02d %02d:%02d:%02d | Temp: %d | Humidity: %d\n",
		rec0[0], rec0[1], rec0[2], rec0[3], rec0[4], rec0[5],
		rec0[6], rec0[8],
	)

	tempSet := int(rec0[6])
	humSet := int(rec0[7])
	tempAct := int(rec0[8])
	humAct := int(rec0[9])

	mastersData := &models.MasterData{
		MachineMasterID: machineMasterId,
		MachineYear:     int(rec0[0]),
		MachineMonth:    int(rec0[1]),
		MachineDay:      int(rec0[2]),
		MachineHour:     int(rec0[3]),
		MachineMinute:   int(rec0[4]),
		MachineSecond:   int(rec0[5]),
		TempSet:         &tempSet,
		HumSet:          &humSet,
		TempAct:         &tempAct,
		HumAct:          &humAct,
	}
	err := databases.DiantaDB.Save(&mastersData).Error
	return mastersData, err
}
