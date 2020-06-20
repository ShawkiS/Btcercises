from binascii import hexlify, unhexlify
from io import BytesIO
from unittest import TestCase

from helper import (
    double_sha256,
    int_to_little_endian,
    little_endian_to_int,
)


NETWORK_MAGIC = b'\xf9\xbe\xb4\xd9'

class NetworkEnvelope:

    def __init__(self, command, payload):
        self.command = command
        self.payload = payload

    def __repr__(self):
        return 'Network: Bitcoin mainnet\nCommand: {}\nPayload: {}'.format(
            self.command.decode('ascii'),
            hexlify(self.payload).decode('ascii'),
        )

    @classmethod
    def parse(cls, s):
        '''Takes a stream and creates a NetworkEnvelope'''
        # check the network magic b'\xf9\xbe\xb4\xd9'
        magic = s.read(4)
        if magic != NETWORK_MAGIC:
            raise RuntimeError('Network Magic not at beginning of stream')
        # command 12 bytes
        command = s.read(12)
        # payload length 4 bytes, little endian
        payload_length = little_endian_to_int(s.read(4))
        # checksum 4 bytes
        checksum = s.read(4)
        # payload
        payload = s.read(payload_length)
        # check the checksum
        if double_sha256(payload)[:4] != checksum:
            raise RuntimeError('Payload and Checksum do not match')
        return cls(command, payload)

    def serialize(self):
        '''Returns the byte serialization of the entire network message'''
        result = NETWORK_MAGIC + self.command
        payload_length = int_to_little_endian(len(self.payload), 4)
        checksum = double_sha256(self.payload)[:4]
        result += payload_length + checksum + self.payload
        return result


class NetworkEnvelopeTest(TestCase):

    def test_parse(self):
        msg = unhexlify('f9beb4d976657261636b000000000000000000005df6e0e2')
        stream = BytesIO(msg)
        envelope = NetworkEnvelope.parse(stream)
        self.assertEqual(envelope.command[:6], b'verack')
        self.assertEqual(envelope.payload, b'')
        msg = unhexlify('f9beb4d976657273696f6e0000000000650000005f1a69d2721101000100000000000000bc8f5e5400000000010000000000000000000000000000000000ffffc61b6409208d010000000000000000000000000000000000ffffcb0071c0208d128035cbc97953f80f2f5361746f7368693a302e392e332fcf05050001')
        stream = BytesIO(msg)
        envelope = NetworkEnvelope.parse(stream)
        self.assertEqual(envelope.command[:7], b'version')
        self.assertEqual(envelope.payload, msg[24:])

    def test_serialize(self):
        msg = unhexlify('f9beb4d976657261636b000000000000000000005df6e0e2')
        stream = BytesIO(msg)
        envelope = NetworkEnvelope.parse(stream)
        self.assertEqual(envelope.serialize(), msg)
        msg = unhexlify('f9beb4d976657273696f6e0000000000650000005f1a69d2721101000100000000000000bc8f5e5400000000010000000000000000000000000000000000ffffc61b6409208d010000000000000000000000000000000000ffffcb0071c0208d128035cbc97953f80f2f5361746f7368693a302e392e332fcf05050001')
        stream = BytesIO(msg)
        envelope = NetworkEnvelope.parse(stream)
        self.assertEqual(envelope.serialize(), msg)

