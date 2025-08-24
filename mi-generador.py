import sys

# Constants with the contents of the original .yaml file
HEADER = """
name: tp0
services:
"""

SERVER_BLOCK = """
  server:
    container_name: server
    image: server:latest
    entrypoint: python3 /main.py
    environment:
      - PYTHONUNBUFFERED=1
    networks:
      - testing_net
    volumes: 
      - ./server/config.ini:/config.ini
"""

NETWORK_BLOCK = """
networks:
  testing_net:
    ipam:
      driver: default
      config:
        - subnet: 172.25.125.0/24
"""

# Parameterizable client using format
def client_block(n):
    return f"""
  client{n}:
    container_name: client{n}
    image: client:latest
    entrypoint: /client
    environment:
      - CLI_ID={n}
      - NOMBRE=Santiago Lionel
      - APELLIDO=LORCA
      - DOCUMENTO=30904465
      - NACIMIENTO = 1999-03-17
      - NUMERO = 7574
    networks:
      - testing_net
    depends_on:
      - server
    volumes:
      - ./client/config.yaml:/config.yaml
"""

if __name__ == "__main__":
    arguments = sys.argv
    if len(arguments) > 2:
        output_file = arguments[1]
        clients = int(arguments[2])
        # with open to ensure the file is closed
        with open(output_file, 'w') as file:
            file.write(HEADER)
            file.write(SERVER_BLOCK)
            for n in range(1, clients + 1):
                file.write(client_block(n))
            file.write(NETWORK_BLOCK)    
    else:
        print("Insufficient number of arguments")    

        