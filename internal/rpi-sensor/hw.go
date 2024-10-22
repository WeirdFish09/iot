package rpi_sensor

import (
	"fmt"
	"log"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/conn/v3/spi"
	"periph.io/x/conn/v3/spi/spireg"
	"periph.io/x/host/v3"
	"sync"
)

const (
	channel = 7
)

var (
	tempretures []int32 = []int32{-20, -19, -18, -17, -16, -15, -14, -13, -12, -11,
		-10, -9, -8, -7, -6, -5, -4, -3, -2, -1,
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9,
		10, 11, 12, 13, 14, 15, 16, 17, 18, 19,
		20, 21, 22, 23, 24, 25, 26, 27, 28, 29,
		30, 31, 32, 33, 34, 35, 36, 37, 38, 39,
		40}
	resistances []float32 = []float32{89.776,
		85.137,
		80.750,
		76.600,
		72.676,
		68.963,
		65.451,
		62.129,
		58.986,
		56.012,
		53.198,
		50.534,
		48.013,
		45.627,
		43.368,
		41.229,
		39.204,
		37.285,
		35.468,
		33.747,
		32.116,
		30.570,
		29.105,
		27.716,
		26.399,
		25.150,
		23.965,
		22.842,
		21.776,
		20.764,
		19.783,
		18.892,
		18.026,
		17.204,
		16.423,
		15.681,
		14.976,
		14.306,
		13.669,
		13.063,
		12.487,
		11.939,
		11.418,
		10.921,
		10.449,
		10.000,
		9.571,
		9.164,
		8.775,
		8.405,
		8.052,
		7.716,
		7.396,
		7.090,
		6.798,
		6.520,
		6.255,
		6.002,
		5.760,
		5.529,
		5.309}
	ranges [][]float32 = [][]float32{}
	m      sync.Mutex
)

type SPIHandler struct {
	Conn spi.Conn
}

func HWInit() (SPIHandler, error) {
	ranges = make([][]float32, 60)
	for i := 0; i < len(resistances)-1; i++ {
		ranges[i] = make([]float32, 2)
		ranges[i][0] = resistances[i]
		ranges[i][1] = resistances[i+1]
	}
	if _, err := host.Init(); err != nil {
		log.Fatalf("failed to initialize periph: %v", err)
	}
	// Open the default SPI
	port, err := spireg.Open("")
	if err != nil {
		log.Fatalf("failed to open SPI: %v", err)
	}
	// defer port.Close()
	// Configure SPI connection parameters (MCP3008 works at up to 3.6 MHz)
	conn, err := port.Connect(physic.MegaHertz*1, spi.Mode0, 8)
	if err != nil {
		log.Fatalf("failed to connect to SPI device: %v", err)
	}
	var spiHandler SPIHandler
	spiHandler.Conn = conn
	return spiHandler, nil
}

func (spi *SPIHandler) ReadTemp() (int, error) {
	m.Lock()
	val, err := readMCP3008(spi.Conn, channel)
	if err != nil {
		return 0, err
	}
	m.Unlock()
	return parseResistence(float32(val)), nil
}

func getTempFromResistence(res float32) int {
	index := -1
	for i := 0; i < len(ranges); i++ {
		if res > ranges[i][1] && res < ranges[i][0] {
			index = i
			break
		}
	}
	if index == -1 {
		return -100
	} else {
		mid := (ranges[index][0] + ranges[index][1]) / 2
		if res < mid {
			return index - 20
		} else {
			return index - 20 + 1
		}
	}
}

func parseResistence(val float32) int {
	voltsOut := (val / 1024) * 3.3
	resistence := 10 * ((3.3 / voltsOut) - 1) // other resistor is 10KOhm
	fmt.Println("Read resistence: %v", resistence)
	return getTempFromResistence(resistence)
}

func readMCP3008(conn spi.Conn, channel int) (int, error) {
	if channel < 0 || channel > 7 {
		return 0, fmt.Errorf("invalid channel: %d", channel)
	}
	// Create the command to send (3 bytes):
	// Byte 1: Start bit (1), single-ended (1), channel bits (D2, D1, D0)
	command := []byte{
		0x01,                        // Start bit (1)
		byte(0x80 | (channel << 4)), // SGL/DIFF = 1, channel bits D2-D0
		0x00,                        // Empty byte
	}
	read := make([]byte, 3)
	if conn == nil {
		log.Fatal("conn is nil")
	}
	// Send the command and read the response
	if err := conn.Tx(command, read); err != nil {
		return 0, fmt.Errorf("failed to communicate with MCP3008: %v", err)
	}
	// Extract the 10-bit result from the received data
	//The result is in the last 10 bits of the second and third bytes
	result := ((int(read[1]) & 0x03) << 8) | int(read[2])
	fmt.Println("result: %v", result)
	// (byte2 & 11) << 8 | byte3
	return result, nil
}
