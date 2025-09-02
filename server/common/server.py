import socket
import logging
import signal
from common.utils import *
from common.socket_utils import *
from common.protocol_codec import *
from multiprocessing import Process, Manager, Barrier, Lock, BrokenBarrierError

TYPE_BET = 1
TYPE_DONE = 2
SERVER_SHUTDOWN = 255

class Server:
    def __init__(self, port, listen_backlog, expected_clients):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        self._server_socket.settimeout(1)
        self._running = True
        signal.signal(signal.SIGTERM, self.shutdown)
        self._expected_clients = expected_clients

        self.manager = Manager()
        self._winners_shared = self.manager.dict()
        self._barrier = Barrier(expected_clients) 
        self._io_lock = Lock()

    """
    Closing of file descriptors is contemplated before the main application thread dies
    """

    def shutdown(self, signum, frame):
        self._running = False
        try:
            self._barrier.abort()
        except Exception:
            pass
        try:
            self._server_socket.close()
        except OSError:
            pass
        logging.info("action: exit | result: success")

    def __send_shutdown(self, sock):
        try:
            mustWriteAll(sock, (SERVER_SHUTDOWN).to_bytes(1, "big"))
        except Exception:
            pass

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
                try: 
                    client_sock.close()
                except OSError:
                    pass
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
        agency = None
        try:
            while True:
                try:
                    type = packet_type(client_sock)
                    if type == TYPE_BET:
                        self.__process_bet(client_sock)
                    elif type == TYPE_DONE:
                        agency = self.__process_done(client_sock)
                        break
                    else:
                        break
                except OSError as e:
                    logging.error(f"action: receive_message | result: fail | error: {e}")
                    return
            try:
                index = self._barrier.wait()  
            except BrokenBarrierError:
                self.__send_shutdown(client_sock)
                return
            
            if index == 0:
                self.__draw_winners()
            try:
                self._barrier.wait()
            except BrokenBarrierError:
                self.__send_shutdown(client_sock)
                return
            if agency is not None:
                self.__send_winners(client_sock, agency)
        finally:
            try:
                client_sock.close()
            except:
                pass     

    def __process_bet(self, client_sock):
        bets = decode_bet_batch(client_sock)
        addr = client_sock.getpeername()
        logging.info(f'action: receive_message | result: success | ip: {addr[0]}')
        with self._io_lock:
            store_bets(bets)
        logging.info(f'action: apuesta_recibida | result: success | cantidad: {len(bets)}')
        mustWriteAll(client_sock, (1).to_bytes(1, "big"))

    def __process_done(self, client_sock):
        agency_bytes = mustReadAll(client_sock, 1)
        agency = int.from_bytes(agency_bytes, "big")
        logging.info(f"action: done | result: success | agency: {agency}")
        mustWriteAll(client_sock, (1).to_bytes(1, "big"))
        return agency

    def __draw_winners(self):
        logging.info("action: sorteo | result: in_progress")
        winners_local = {i: [] for i in range(1, self._expected_clients + 1)}
        for bet in load_bets():
            if has_won(bet):
                winners_local[int(bet.agency)].append(int(bet.document))
        for ag, docs in winners_local.items():
            self._winners_shared[ag] = docs
        logging.info("action: sorteo | result: success")

    def __send_winners(self, client_sock, agency):
        docs = list(self._winners_shared.get(int(agency), []))
        mustWriteAll(client_sock, len(docs).to_bytes(2, "big"))
        for doc in docs:
            mustWriteAll(client_sock, int(doc).to_bytes(8, "big"))

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
