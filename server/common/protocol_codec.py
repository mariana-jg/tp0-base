import struct
from common.utils import *
from common.socket_utils import *

def decode_bet(socket):
    # Obtengo los bytes que hay que leer a continuacion, el payload
    # ! -> Big endian, H -> UINT16 - 2bytes
    total_lenght_bytes = avoid_short_reads(socket, 2)
    if total_lenght_bytes is None:
        return None
    (payload_len,) = struct.unpack("!H", total_lenght_bytes)
    # Traigo el payload completo
    payload = avoid_short_reads(socket, payload_len)
    if payload is None:
        return None

    offset = 0

    #Agencia con lenght 1
    agency = payload[offset]

    offset += 1

    #Nombre con lenght varaible

    (name_len,) = struct.unpack("!H", payload[offset:offset + 2])
    offset += 2
    name = payload[offset:offset+name_len].decode("utf-8")
    offset += name_len

    #Apellido con lenght varaible

    (lastname_len,) = struct.unpack("!H", payload[offset:offset + 2])
    offset += 2
    lastname = payload[offset:offset+lastname_len].decode("utf-8")
    offset += lastname_len    

    #documento con lenght 8b Q-> UINT64
    (document,) = struct.unpack('!Q', payload[offset:offset+8])
    offset += 8

    birth_bytes = payload[offset:offset+10]
    offset += 10
    birthdate = birth_bytes.decode('ascii')
    
    # numero con 2b
    (number,) = struct.unpack('!H', payload[offset:offset+2])
    offset += 2

    return Bet(agency, name, lastname, document, birthdate, number)



