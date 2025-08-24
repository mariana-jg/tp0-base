import socket
import logging
import signal
from common.utils import *
from common.socket_utils import *
from common.protocol_codec import *


class Server:
    def __init__(self, port, listen_backlog):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        self._server_socket.settimeout(1)
        self._running = True
        signal.signal(signal.SIGTERM, self.shutdown)
        self._client_sockets = []

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
        try:
            bet = decode_bet(client_sock)
            addr = client_sock.getpeername()
            logging.info(f'action: receive_message | result: success | ip: {addr[0]}')
            store_bets([bet])
            logging.info(f'action: apuesta_almacenada | result: success | dni: {bet.document} | numero: {bet.number}')
            avoid_short_writes(client_sock, struct.pack('>B', 1))
        except OSError as e:
            logging.error("action: receive_message | result: fail | error: {e}")
        finally:
            client_sock.close()
            self._client_sockets.remove(client_sock)

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
