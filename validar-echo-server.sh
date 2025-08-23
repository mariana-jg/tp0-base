#!/bin/bash

mensaje="Este es un mensaje de prueba"

PUERTO_SERVIDOR=$(grep SERVER_PORT server/config.ini | cut -d ' ' -f 3) 
IP_SERVIDOR=$(grep SERVER_IP server/config.ini | cut -d ' ' -f 3)

respuesta=$(docker run --rm --network tp0_testing_net busybox:latest sh -c "echo '$mensaje' | nc $IP_SERVIDOR $PUERTO_SERVIDOR")

if [ "$respuesta" = "$mensaje" ]; then
    echo "action: test_echo_server | result: success"
else
    echo "action: test_echo_server | result: fail"
fi