"""
Made to avoid short writes.
"""
def avoid_short_writes(socket, data):
    total = len(data)
    sent = 0
    while sent < total:
        pending = data[sent:]
        sent += socket.send(pending)

"""
Made to avoid short reads.
"""
def avoid_short_reads(socket, expected_length):
    data = bytearray()
    while len(data) < expected_length:
        packet = socket.recv(expected_length-len(data))
        if len(packet) == 0:
            return None
        else:
            data.extend(packet)
    return data  