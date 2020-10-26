#!/usr/bin/env python3
# -*- coding: utf-8 -*-

# You can test this by:
#   ECHO_ADDRESS="0.0.0.0" ECHO_PORT=31337 ./tcp-server.py &
#   echo "hello" | nc -4 -w1 localhost 31337

from socket import socket, gethostname, AF_INET, SOCK_STREAM, SOL_SOCKET, SO_REUSEADDR
import logging
from os import environ
from threading import Thread

__ENV_ADDRESS = 'ECHO_ADDRESS'
__ENV_PORT = 'ECHO_PORT'
__ENV_MAX_CLIENTS = 'ECHO_MAX_CLIENTS'

class DispatchThread(Thread):
	def __init__(self, connection):
		Thread.__init__(self)
		self.connection = connection
		self.client_address = connection.getpeername()
		logging.info('Connection from %s' % str(self.client_address))

	def run(self):
		try:
			# Receive the data in small chuncks and retransmit it
			while True:
				payload = self.connection.recv(512)
				logging.info("<< '%s'" % payload.decode('utf-8'))
				if payload:
					if len(payload)>0:
						self.connection.sendall(payload)
				else:
					logging.info('No more data from %s' % str(self.client_address))
					break
		finally:
			self.connection.close()

class echo_tcp_server:
	def __init__(self, server_address = '0.0.0.0', server_port = 7, server_max_clients = 16):
		self.server_address = server_address
		self.server_port = server_port
		self.server_max_clients = server_max_clients
		self.sock = socket(AF_INET, SOCK_STREAM)
		self.sock.setsockopt(SOL_SOCKET, SO_REUSEADDR, 1)
		server = (server_address, server_port)
		self.sock.bind(server)
		self.sock.listen(server_max_clients)
		self.threads = []
		logging.info("Listening on " + server_address + ":" + str(server_port) + "/tcp")
		logging.info("You can test the connection by: nc localhost " + str(server_port))

	def dispatch(self, connection):
		new_thread = DispatchThread(connection)
		new_thread.start()
		self.threads.append(new_thread)

	def run(self):
		while True:
			logging.info("Waiting for connection")
			connection, (client_ip, client_port) = self.sock.accept()
			logging.info("Dispatching client %s:%s" % (str(client_ip), str(client_port)))
			self.dispatch(connection)

		for item in self.threads:
			item.join()

if __name__ == '__main__':
	logging.basicConfig(level=environ.get("LOGLEVEL", "INFO"))
	__server_address = environ.get(__ENV_ADDRESS, '0.0.0.0')
	__server_port = int(environ.get(__ENV_PORT, 7))
	__server_max_clients = int(environ.get(__ENV_MAX_CLIENTS, 16))
	echo_tcp_server(__server_address, __server_port, __server_max_clients).run()
