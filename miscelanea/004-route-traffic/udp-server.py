#!/usr/bin/env python3
# -*- coding: utf-8 -*-

# You can test this by:
#   ECHO_ADDRESS="0.0.0.0" ECHO_PORT=31337 ./udp-server.py &
#   echo "hello" | nc -4u -w1 localhost 31337

from socket import socket, gethostname, AF_INET, SOCK_DGRAM
import logging
from os import environ


__ENV_ADDRESS = 'ECHO_ADDRESS'
__ENV_PORT = 'ECHO_PORT'


class echo_udp_server:
	def __init__(self, server_address = '0.0.0.0', server_port = 7):
		self.server_address = server_address
		self.server_port = server_port
		self.sock = socket(AF_INET, SOCK_DGRAM)
		server = (server_address, server_port)
		self.sock.bind(server)
		logging.info("Listening on " + server_address + ":" + str(server_port) + "/udp")
		logging.info("You can test the connection by: nc -4u localhost " + str(server_port))

	def run(self):
		while True:
			payload, client_address = self.sock.recvfrom(512)
			logging.info("<< '%s'" % payload.decode('utf-8'))
			sent = self.sock.sendto(payload, client_address)

if __name__ == '__main__':
	logging.basicConfig(level=environ.get("LOGLEVEL", "INFO"))
	__server_address = environ.get(__ENV_ADDRESS, '0.0.0.0')
	__server_port = int(environ.get(__ENV_PORT, 7))
	echo_udp_server(__server_address, __server_port).run()
