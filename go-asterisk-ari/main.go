package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/CyCoreSystems/ari/v6"
	"github.com/CyCoreSystems/ari/v6/client/native"
	"github.com/CyCoreSystems/ari/v6/ext/play"
)

func main() {
	log.Print("Conectando a ARI")

	// Crear un contexto con cancelación
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Configurar el cliente ARI
	cl, err := native.Connect(&native.Options{
		Application:  "my_app",
		Username:     "test1",
		Password:     "test1",
		URL:          "http://localhost:8088/ari",
		WebsocketURL: "ws://localhost:8088/ari/events",
	})
	if err != nil {
		log.Printf("Error al conectar con ARI: %v", err)
		return
	}
	defer cl.Close()

	log.Print("Conectado exitosamente a ARI")

	// Obtener información de Asterisk
	info, err := cl.Asterisk().Info(nil)
	if err != nil {
		log.Printf("Error al obtener información de Asterisk: %v", err)
		return
	}
	log.Printf("Información de Asterisk: %+v", info)

	// Suscribirse a todos los eventos
	log.Print("Intentando suscribirse a eventos...")
	sub := cl.Bus().Subscribe(nil, ari.Events.All)
	if sub == nil {
		log.Print("Error: No se pudo crear la suscripción")
		return
	}
	defer sub.Cancel()

	log.Print("Suscrito a eventos de Asterisk")

	// Crear un canal para manejar señales de terminación
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Bucle principal para escuchar eventos
	log.Print("Iniciando bucle de eventos...")
	for {
		select {
		case event := <-sub.Events():
			// Imprimir el tipo de evento
			log.Printf("Evento recibido: %T", event)

			// Imprimir el JSON del evento
			jsonData, err := json.MarshalIndent(event, "", "  ")
			if err != nil {
				log.Printf("Error al convertir evento a JSON: %v", err)
			} else {
				log.Printf("JSON del evento:\n%s", string(jsonData))
			}

			// Manejar diferentes tipos de eventos
			switch e := event.(type) {
			case *ari.StasisStart:
				log.Printf("Evento: StasisStart")
				log.Printf("Inicio de Stasis en canal: %s", e.Channel.ID)
				log.Printf("Detalles del canal: %+v", e.Channel)
				go handleChannel(ctx, cl, e.Channel.ID)

			case *ari.ChannelCreated:
				log.Printf("Evento: ChannelCreated")
				log.Printf("Canal creado: %s", e.Channel.ID)
				log.Printf("Detalles: %+v", e.Channel)

			case *ari.ChannelDestroyed:
				log.Printf("Evento: ChannelDestroyed")
				log.Printf("Canal destruido: %s, Causa: %s", e.Channel.ID, e.CauseTxt)

			case *ari.ChannelStateChange:
				log.Printf("Evento: ChannelStateChange")
				log.Printf("Cambio de estado en canal: %s", e.Channel.ID)
				log.Printf("Estado: %+v", e.Channel)

			case *ari.ChannelDtmfReceived:
				log.Printf("Evento: ChannelDtmfReceived")
				log.Printf("DTMF recibido en canal: %s, Dígito: %s, Duración: %dms",
					e.Channel.ID, e.Digit, e.DurationMs)

			case *ari.ChannelHangupRequest:
				log.Printf("Evento: ChannelHangupRequest")
				log.Printf("Solicitud de colgar en canal: %s, Causa: %d",
					e.Channel.ID, e.Cause)

			case *ari.BridgeCreated:
				log.Printf("Evento: BridgeCreated")
				log.Printf("Bridge creado: %s", e.Bridge.ID)

			case *ari.BridgeDestroyed:
				log.Printf("Evento: BridgeDestroyed")
				log.Printf("Bridge destruido: %s", e.Bridge.ID)

			case *ari.ChannelEnteredBridge:
				log.Printf("Evento: ChannelEnteredBridge")
				log.Printf("Canal entró al bridge: Canal=%s, Bridge=%s",
					e.Channel.ID, e.Bridge.ID)

			case *ari.ChannelLeftBridge:
				log.Printf("Evento: ChannelLeftBridge")
				log.Printf("Canal salió del bridge: Canal=%s, Bridge=%s",
					e.Channel.ID, e.Bridge.ID)

			case *ari.PlaybackStarted:
				log.Printf("Evento: PlaybackStarted")
				log.Printf("Inicio de reproducción: %s", e.Playback.ID)

			case *ari.PlaybackFinished:
				log.Printf("Evento: PlaybackFinished")
				log.Printf("Fin de reproducción: %s", e.Playback.ID)

			case *ari.RecordingStarted:
				log.Printf("Evento: RecordingStarted")
				log.Printf("Inicio de grabación: %s", e.Recording.ID())

			case *ari.RecordingFinished:
				log.Printf("Evento: RecordingFinished")
				log.Printf("Fin de grabación: %s", e.Recording.ID())

			case *ari.StasisEnd:
				log.Printf("Evento: StasisEnd")
				log.Printf("Fin de Stasis en canal: %s", e.Channel.ID)

			default:
				log.Printf("Evento no manejado: %T", event)
			}

		case <-ctx.Done():
			log.Print("Contexto cancelado, cerrando...")
			return
		case <-sigChan:
			log.Print("Señal de terminación recibida, cerrando...")
			cancel()
			return
		}
	}
}


func handleChannel(ctx context.Context, cl ari.Client, channelID string) {
	log.Printf("Manejando canal: %s", channelID)

	// Obtener el handle del canal
	channel := cl.Channel().Get(&ari.Key{ID: channelID, App: "iperfex", Node: ""})
	if channel == nil {
		log.Printf("No se pudo obtener el handle del canal: %s", channelID)
		return
	}
	defer channel.Hangup()

	// Crear un contexto para este canal
	channelCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Suscribirse a eventos de fin de Stasis
	endSub := channel.Subscribe(ari.Events.StasisEnd)
	defer endSub.Cancel()

	// Manejar el fin de Stasis
	go func() {
		<-endSub.Events()
		cancel()
	}()

	// Contestar la llamada
	if err := channel.Answer(); err != nil {
		log.Printf("Error al contestar la llamada: %v", err)
		return
	}
	log.Printf("Llamada contestada en canal: %s", channelID)

	// Reproducir audio
	if err := play.Play(channelCtx, channel, play.URI("sound:tt-monkeys")).Err(); err != nil {
		log.Printf("Error al reproducir audio: %v", err)
		return
	}
	log.Printf("Reproducción de audio completada en canal: %s", channelID)

	// Esperar a que termine el contexto
	<-channelCtx.Done()
	log.Printf("Canal terminado: %s", channelID)
}
