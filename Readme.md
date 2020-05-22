# Kaco Powador Exporter
A Prometheus and MQTT Exporter for Kaco Powador Inverters.

Only tested with KACO Powador 10.0 TL3 and might requires some modification for other models.

Some Inverter models by `SUNPOWER`, `BLUEPLANET`, `SCHÜCO`, `INVENTUX` and `WÜRTH` might work too, as they seam to use the same controller software.

## Prometheus metrics
Prometheus metrics are available on `/metrics`

## MQTT Home-Assistant integration
Additionally, the exporter is able to publish metrics to Home-Assistant via MQTT.

After configuring the required environment variables (see Docker-Compose), HA should auto-discover the inverter.

## Docker Compose
```
  kacopowador_exporter:
    image: kacopowador_exporter
    ports:
      - 8080:8080
    restart: always
    environment:
      - KACOPOWADOR_HOST=KACOIP
      - KACOPOWADOR_METRICPORT=8080
#     - KACOPOWADOR_MQTTHOST=MQTTHOSTIP:1883
#     - KACOPOWADOR_MQTTUSER=hassio
#     - KACOPOWADOR_MQTTPASSWORD=PASSWORD
#     - KACOPOWADOR_MQTTNAME=kaco
#     - KACOPOWADOR_MQTTINTERVAL=30
```

