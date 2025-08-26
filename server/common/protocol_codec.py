import struct
from common.utils import *
from common.socket_utils import *

AGENCY_LENGHT = 1
FRAME_LENGHT = 2
DOCUMENT_LENGHT = 8
BIRTHDATE_LENGHT = 10
NUMBER_LENGHT = 2
LEN_NAME_LENGHT = 2
LEN_LASTNAME_LENGHT = 2

def decode_bet(socket):

    total_lenght_bytes = avoid_short_reads(socket, FRAME_LENGHT)
    if total_lenght_bytes is None:
        return None
    
    (payload_len,) = struct.unpack("!H", total_lenght_bytes)
    
    payload = avoid_short_reads(socket, payload_len)
    if payload is None:
        return None

    offset = 0

    agency = payload[offset]
    offset += AGENCY_LENGHT

    (name_len,) = struct.unpack("!H", payload[offset:offset + LEN_NAME_LENGHT])
    offset += LEN_NAME_LENGHT
    name = payload[offset:offset+name_len].decode("utf-8")
    offset += name_len

    (lastname_len,) = struct.unpack("!H", payload[offset:offset + LEN_LASTNAME_LENGHT])
    offset += LEN_LASTNAME_LENGHT
    lastname = payload[offset:offset+lastname_len].decode("utf-8")
    offset += lastname_len    

    (document,) = struct.unpack('!Q', payload[offset:offset+DOCUMENT_LENGHT])
    offset += DOCUMENT_LENGHT

    birthdate_bytes = payload[offset:offset+BIRTHDATE_LENGHT]
    offset += BIRTHDATE_LENGHT
    birthdate = birthdate_bytes.decode('ascii')
    
    (number,) = struct.unpack('!H', payload[offset:offset+NUMBER_LENGHT])
    offset += NUMBER_LENGHT

    return Bet(agency, name, lastname, document, birthdate, number)

def decode_bet_batch(socket):
    len_batch = avoid_short_reads(socket, 1)[0]
    bets = [decode_bet(socket) for _ in range(len_batch)]
    return bets    
    

