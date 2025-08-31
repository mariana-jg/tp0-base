import socket
import logging
import signal
from common.utils import *
from common.socket_utils import *
from common.protocol_codec import *
from multiprocessing import Process, Manager, Barrier

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

        self.manager = Manager()
        self._waiting_clients = self.manager.dict()

        self._barrier = Barrier(expected_clients)


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
                    p = Process(target=self.__handle_client_connection, args=(client_sock,))
                    p.daemon = True
                    p.start()
            except socket.timeout:
                continue
            except OSError:
                break
                
    def __handle_client_connection(self, client_sock):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        self._client_sockets.append(client_sock)
        agency = None

        try:
            while True:
                try:
                    t = packet_type(client_sock)
                    if t == 1:
                        bets = decode_bet_batch(client_sock)
                        addr = client_sock.getpeername()
                        logging.info(f'action: receive_message | result: success | ip: {addr[0]}')
                        store_bets(bets)
                        logging.info(f'action: apuesta_recibida | result: success | cantidad: {len(bets)}')
                        mustWriteAll(client_sock, struct.pack('>B', 1))

                    elif t == 2:
                        # cliente terminó; recibo agencia
                        agency_bytes = mustReadAll(client_sock, 1)
                        agency = struct.unpack("!B", agency_bytes)[0]
                        logging.info(f"action: done | result: success | agency: {agency}")
                        mustWriteAll(client_sock, struct.pack('>B', 1))
                        break
                    else:
                        break
                except OSError as e:
                    logging.error(f"action: receive_message | result: fail | error: {e}")
                    return
            
            idx = self._barrier.wait()  

            if idx == 0:
                logging.info("action: sorteo | result: in_progress")
                winners_local = {i: [] for i in range(1, self._expected_clients + 1)}
                for bet in load_bets():
                    if has_won(bet):
                        winners_local[int(bet.agency)].append(int(bet.document))

                
                for ag, docs in winners_local.items():
                    self._winners_shared[ag] = docs
                logging.info("action: sorteo | result: success")

            self._barrier.wait()

            if agency is not None:
                docs = list(self._winners_shared.get(int(agency), []))
                mustWriteAll(client_sock, struct.pack('!H', len(docs)))
                for doc in docs:
                    mustWriteAll(client_sock, struct.pack('!Q', int(doc)))
        finally:
            try:
                client_sock.close()
            except:
                pass
            try:
                self._client_sockets.remove(client_sock)
            except:
                pass     

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
