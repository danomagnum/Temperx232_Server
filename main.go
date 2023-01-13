package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/tarm/serial"
)

const comport = "/dev/ttyUSB1"

func main() {

	c := &serial.Config{Name: comport, Baud: 9600}
	f, err := serial.OpenPort(c)
	if err != nil {
		log.Fatalf("error opening serial port: %v", err)
		os.Exit(1)
	}

	defer f.Close()

	trigger := []byte("ReadTemp\n")
	err_ctr := 0

	go class1serve()

	for {
		time.Sleep(time.Millisecond * 100)

		n, err := f.Write(trigger)
		if err != nil {
			log.Printf("wrote %d bytes (%v) with error %v", n, trigger, err)
			err_ctr++
			if err_ctr > 100 {
				log.Fatal("100 bad reads in a row. Bailing out.")
				panic(err)
			}
			continue
		}
		b := make([]byte, 256)
		_, err = f.Read(b)
		if err != nil {
			log.Printf("problem reading: %v", err)
			os.Exit(1)
		}
		//log.Printf("got %d bytes: %s, ", n, string(b))
		t, h, err := ParseTempHum(string(b))
		if err != nil {
			log.Printf("bad parse. %v", err)
			err_ctr++
			continue
		}
		//log.Printf("Temp: %2.2f     Hum:%2.2f", t, h)
		ioProvider.Mutex.Lock()
		inInstance.Temperature = t
		inInstance.Humidity = h
		ioProvider.Mutex.Unlock()
		err_ctr = 0

	}

}

func ParseTempHum(result string) (t, h float32, err error) {
	data_start_pos := strings.Index(result, "Temp-Inner:")
	data := result[data_start_pos+11:]
	split := strings.Split(data, ",")
	if len(split) != 2 {
		err = fmt.Errorf("expected two items but got %d from %s", len(split), result)
		return
	}
	temp_raw := split[0]
	hum_raw := split[1]

	temp_str := strings.Split(temp_raw, " ")[0]
	hum_str := strings.Split(hum_raw, " ")[0]

	t64, err := strconv.ParseFloat(temp_str, 32)
	if err != nil {
		return
	}
	h64, err := strconv.ParseFloat(hum_str, 32)

	t = float32(t64)
	h = float32(h64)

	return
}
