from common.utils import *
from common.socket_utils import *

AGENCY_LENGHT = 1
FRAME_LENGHT = 2
DOCUMENT_LENGHT = 8
NUMBER_LENGHT = 2
LEN_NAME_LENGHT = 2
LEN_LASTNAME_LENGHT = 2
LEN_BIRTHDATE_LENGHT = 2

def decode_bet(socket):

    total_lenght_bytes = mustReadAll(socket, FRAME_LENGHT)
    if total_lenght_bytes is None:
        return None
    
    payload_len = int.from_bytes(total_lenght_bytes, "big")
    
    payload = mustReadAll(socket, payload_len)
    if payload is None:
        return None

    offset = 0

    agency = payload[offset]
    offset += AGENCY_LENGHT

    name_len = int.from_bytes(payload[offset:offset + LEN_NAME_LENGHT], "big")
    offset += LEN_NAME_LENGHT
    name = payload[offset:offset+name_len].decode("utf-8")
    offset += name_len

    lastname_len = int.from_bytes(payload[offset:offset + LEN_LASTNAME_LENGHT], "big")
    offset += LEN_LASTNAME_LENGHT
    lastname = payload[offset:offset+lastname_len].decode("utf-8")
    offset += lastname_len    

    document = int.from_bytes(payload[offset:offset+DOCUMENT_LENGHT], "big")
    offset += DOCUMENT_LENGHT

    birthdate_len = int.from_bytes(payload[offset:offset + LEN_BIRTHDATE_LENGHT], "big")
    offset += LEN_BIRTHDATE_LENGHT
    birthdate = payload[offset:offset+birthdate_len].decode("utf-8")
    offset += birthdate_len
    
    number = int.from_bytes(payload[offset:offset+NUMBER_LENGHT], "big")
    offset += NUMBER_LENGHT

    return Bet(agency, name, lastname, document, birthdate, number)

def decode_bet_batch(socket):
    len_batch = int.from_bytes(mustReadAll(socket, 1), "big")
    bets = [decode_bet(socket) for _ in range(len_batch)]
    return bets    

def packet_type(sock):
    b = mustReadAll(sock, 1)
    if not b or len(b) != 1:
        return None
    return int.from_bytes(b, "big")





    

