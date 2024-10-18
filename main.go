// This example simply captures data from your default microphone until you press Enter, after which it plays back the captured audio.
package main

import (
	"fmt"
	"github.com/jacobsa/go-serial/serial"
	"io"
	"os"

	"github.com/gen2brain/malgo"
)

func connectToDevice() (io.ReadWriteCloser, error) {
	options := serial.OpenOptions{
		PortName:        "COM4", // Change this to your serial port
		BaudRate:        115200,
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 1,
	}

	return serial.Open(options)
}

func main() {
	arduino, err := connectToDevice()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ctx, err := malgo.InitContext(nil, malgo.ContextConfig{}, func(message string) {
		fmt.Printf("LOG <%v>\n", message)
	})

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer func() {
		_ = ctx.Uninit()
		ctx.Free()
	}()

	deviceConfig := malgo.DefaultDeviceConfig(malgo.Loopback)
	deviceConfig.Capture.Format = malgo.FormatS16
	deviceConfig.Capture.Channels = 1
	deviceConfig.Playback.Format = malgo.FormatS16
	deviceConfig.Playback.Channels = 1
	//deviceConfig.SampleRate = 44100
	deviceConfig.SampleRate = 44100 / 1
	deviceConfig.Alsa.NoMMap = 1

	onRecvFrames := func(pSample2, pSample []byte, framecount uint32) {
		sum := 0
		max := 0
		min := 0

		fmt.Println(framecount)

		for i := 0; i < len(pSample); i += 2 {
			sample := int(int16(pSample[i]) | int16(pSample[i+1])<<8)
			sum += sample
			if sample > max {
				max = sample
			}
			if sample < min {
				min = sample
			}
		}

		max = max/400 - 20

		if max > 60 {
			max = 60
		}

		if max < 0 {
			max = 0
		}

		message := fmt.Sprintf("%d\n", max)

		_, err := arduino.Write([]byte(message))

		if err != nil {
			fmt.Println(err)
		}

		//output := ""

		//for i := 0; i < max/50; i++ {
		//	output += "#"
		//}
		//
		//fmt.Println(output)
	}

	fmt.Println("Recording...")
	captureCallbacks := malgo.DeviceCallbacks{
		Data: onRecvFrames,
	}
	device, err := malgo.InitDevice(ctx.Context, deviceConfig, captureCallbacks)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(device.Type())

	err = device.Start()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Press Enter to stop recording...")
	fmt.Scanln()

	device.Uninit()
}
