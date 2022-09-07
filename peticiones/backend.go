package peticiones

import (
	"context"
	"fmt"
	"strings"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type Toner struct {
	Hostname      string
	Capacidad     int64
	Uso           int64
	PorcentajeUso float64
	Modelo        string
	Color         string
}

func NewToner(Hostname string, Capacidad int64, Uso int64, Modelo string, Color string) *Toner {
	if Capacidad < 0 {
		Capacidad = 0
	}
	if Uso < 0 {
		Uso = 0
	}

	PorcentajeUso := 0.0
	if Capacidad > 0 && Uso > 0 {
		PorcentajeUso = (float64(Uso) / float64(Capacidad)) * 100
	}

	return &Toner{Hostname, Capacidad, Uso, PorcentajeUso, Modelo, Color}

}

func (d *Toner) String() string {
	return fmt.Sprintf("Modelo: %s - Capacidad: %d - Uso: %d - Porcentaje: %.2f - Color: %s", d.Modelo, d.Capacidad, d.Uso, d.PorcentajeUso, d.Color)
}

func (d *Toner) Influenciar(estampa time.Time) string {
	contenido := fmt.Sprintf(
		"toner,hostname=%s,modelo=\"%s\",color=%s uso=%d,capacidad=%d,porcentaje_uso=%f %d",
		d.Hostname, d.Modelo, d.Color, d.Uso, d.Capacidad, d.PorcentajeUso, estampa.UnixNano())
	return contenido
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

func (b *Backend) Enviar(toners []Toner) error {
	client := influxdb2.NewClient(b.Endpoint, b.Token)
	defer client.Close()

	writeAPI := client.WriteAPIBlocking(b.Organizacion, b.Bucket)

	marcaTiempo := time.Now().Round(time.Second * 60)

	var sentencia strings.Builder
	for _, toner := range toners {
		sentencia.WriteString(toner.Influenciar(marcaTiempo))
		sentencia.WriteString("\n")
	}

	err := writeAPI.WriteRecord(context.Background(), sentencia.String())

	if err != nil {
		return fmt.Errorf("backend:> %v", err)
	}

	return nil
}
