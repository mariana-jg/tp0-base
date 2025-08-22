import sys

# Constantes con el contenido del archivo .yaml original
HEADER = """
name: tp0
services:
"""

BLOQUE_SERVIDOR = """
  server:
    container_name: server
    image: server:latest
    entrypoint: python3 /main.py
    environment:
      - PYTHONUNBUFFERED=1
      - LOGGING_LEVEL=DEBUG
    networks:
      - testing_net
    volumes: 
      - ./server/config.ini:/config.ini
"""

BLOQUE_REDES = """
networks:
  testing_net:
    ipam:
      driver: default
      config:
        - subnet: 172.25.125.0/24
"""

# Cliente parametrizable utilizando format
def bloque_cliente(n):
    return f"""
  client{n}:
    container_name: client{n}
    image: client:latest
    entrypoint: /client
    environment:
      - CLI_ID={n}
      - CLI_LOG_LEVEL=DEBUG
    networks:
      - testing_net
    depends_on:
      - server
    volumes:
      - ./client/config.yaml:/config.yaml
"""

if __name__ == "__main__":
    argumentos = sys.argv
    if len(argumentos) > 2:
        salida = argumentos[1]
        cantidad_de_clientes = int(argumentos[2])
        # with open para asegurar el cierre del archivo
        with open(salida, 'w') as archivo:
            archivo.write(HEADER)
            archivo.write(BLOQUE_SERVIDOR)
            for n in range(1, cantidad_de_clientes + 1):
                archivo.write(bloque_cliente(n))
            archivo.write(BLOQUE_REDES)    
    else:
        print("Cantidad de argumentos insuficiente")    

        