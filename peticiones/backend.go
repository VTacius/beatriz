package peticiones

import (
	"context"
	"fmt"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type Datos struct {
	Hostname      string
	UsoPorcentaje float64
}

type Backend struct {
	Endpoint     string
	Token        string
	Organizacion string
	Bucket       string
}

func NewBackend(Endpoint string, Token string, Organizacion string, Bucket string) *Backend {
	return &Backend{Endpoint, Token, Organizacion, Bucket}
}

func (b *Backend) Enviar(datos Datos) error {
	client := influxdb2.NewClient(b.Endpoint, b.Token)
	defer client.Close()

	writeAPI := client.WriteAPIBlocking(b.Organizacion, b.Bucket)

	marcaTiempo := time.Now().Round(time.Second * 60)
	// Create point using fluent style
	temperatura := influxdb2.NewPointWithMeasurement("toner").
		AddTag("host", datos.Hostname).
		AddField("usoPorcentaje", datos.UsoPorcentaje).
		SetTime(marcaTiempo)

	err := writeAPI.WritePoint(context.Background(), temperatura)

	return fmt.Errorf("backend:> %v", err)
}
