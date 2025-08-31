import socket
import logging
import signal
from common.utils import *
from common.socket_utils import *
from common.protocol_codec import *

TYPE_BET = 1
TYPE_DONE = 2

class Server:
    def __init__(self, port, listen_backlog, expected_clients):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        self._server_socket.settimeout(1)
        self._running = True
        signal.signal(signal.SIGTERM, self.shutdown)
        self._client_sockets = []
        self._done_clients = 0
        self._expected_clients = expected_clients
        self._waiting_winners = {}

    """
    Closing of file descriptors is contemplated before the main application thread dies
    """

    def shutdown(self, signum, frame):
        self._running = False
        self._server_socket.close()
        for socket in self._client_sockets:
            try:
                socket.close()
            except OSError:
                pass    
        logging.info("action: exit | result: success")

    def run(self):
        """
        Dummy Server loop

        Server that accept a new connections and establishes a
        communication with a client. After client with communucation
        finishes, servers starts to accept new connections again
        """

        while self._running:
            try:
                client_sock = self.__accept_new_connection()
                if client_sock: 
                    self.__handle_client_connection(client_sock)
            except socket.timeout:
                continue
            except OSError as error:
                break                

    def __handle_client_connection(self, client_sock):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        self._client_sockets.append(client_sock)
        while True:
            try:
                type = packet_type(client_sock)
                if type == TYPE_BET:
                    self.__process_bet(client_sock)
                elif type == TYPE_DONE:
                    self._done_clients += 1
                    agency_bytes = mustReadAll(client_sock, 1)
                    agency = struct.unpack("!B", agency_bytes)[0]
                    logging.info(f"action: done | result: success | agency: {agency}")
                    self._waiting_winners[agency] = client_sock
                    mustWriteAll(client_sock, struct.pack('>B', 1))
                    break
                else:
                    break
            except OSError as e:
                logging.error("action: receive_message | result: fail | error: {e}")

        if self._done_clients == self._expected_clients:
            logging.info("action: sorteo | result: success")
            
            winners = {i: [] for i in range(1, self._expected_clients+1)}
            
            for bet in load_bets():
                if has_won(bet):
                    winners[int(bet.agency)].append(int(bet.document))

            for agency, sock in list(self._waiting_winners.items()):
                documents = winners.get(int(agency), [])
                mustWriteAll(sock, struct.pack('!H', len(documents)))  

                for doc in documents:
                    mustWriteAll(sock, struct.pack('!Q', int(doc)))  
                try:
                    sock.close()
                except:
                    pass
                self._waiting_winners.pop(agency, None)        

    def __process_bet(self, client_sock):
        bets = decode_bet_batch(client_sock)
        len_bets = len(bets)
        addr = client_sock.getpeername()
        logging.info(f'action: receive_message | result: success | ip: {addr[0]}')
        store_bets(bets)
        logging.info(f'action: apuesta_recibida | result: success | cantidad: {len_bets}')
        mustWriteAll(client_sock, struct.pack('>B', 1))  
    
    def __accept_new_connection(self):
        """
        Accept new connections

        Function blocks until a connection to a client is made.
        Then connection created is printed and returned
        """

        # Connection arrived
        logging.info('action: accept_connections | result: in_progress')
        c, addr = self._server_socket.accept()
        logging.info(f'action: accept_connections | result: success | ip: {addr[0]}')
        return c
