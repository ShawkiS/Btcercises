from binascii import hexlify, unhexlify
from io import BytesIO
from unittest import TestCase

from helper import (
    double_sha256,
    int_to_little_endian,
    little_endian_to_int,
    merkle_parent,
    merkle_parent_level,
    merkle_path,
    merkle_root,
)


class Proof:

    def __init__(self, merkle_root, tx_hash, index, merkle_proof):
        self.merkle_root = merkle_root
        self.tx_hash = tx_hash
        self.index = index
        self.merkle_proof = merkle_proof

    def __repr__(self):
        s = '{}:{}:{}:['.format(
            hexlify(self.merkle_root).decode('ascii'),
            hexlify(self.tx_hash).decode('ascii'),
            self.index,
        )
        for p in self.merkle_proof:
            s += '{},'.format(hexlify(p).decode('ascii'))
        s += ']'
        return s

    def verify(self):
        '''Returns whether this proof is valid'''
        current = self.tx_hash[::-1]
        path = merkle_path(self.index, 2**len(self.merkle_proof))
        for i, proof_hash in enumerate(self.merkle_proof):
            if path[i] % 2 == 1:
                current = merkle_parent(proof_hash, current)
            else:
                current = merkle_parent(current, proof_hash)
        return current[::-1] == self.merkle_root

    
class Block:

    def __init__(self, version, prev_block, merkle_root, timestamp, bits, nonce, tx_hashes=None):
        self.version = version
        self.prev_block = prev_block
        self.merkle_root = merkle_root
        self.timestamp = timestamp
        self.bits = bits
        self.nonce = nonce
        self.tx_hashes = tx_hashes
        self.merkle_tree = None

    @classmethod
    def parse(cls, s):
        '''Takes a byte stream and parses a block. Returns a Block object'''
        # version - 4 bytes, little endian
        version = little_endian_to_int(s.read(4))
        # prev_block - 32 bytes, little endian
        prev_block = s.read(32)[::-1]
        # merkle_root - 32 bytes, little endian
        merkle_root = s.read(32)[::-1]
        # timestamp - 4 bytes, little endian
        timestamp = little_endian_to_int(s.read(4))
        # bits - 4 bytes
        bits = s.read(4)
        # nonce - 4 bytes
        nonce = s.read(4)
        return cls(version, prev_block, merkle_root, timestamp, bits, nonce)

    def serialize(self):
        '''Returns the 80 byte block header'''
        result = int_to_little_endian(self.version, 4)
        result += self.prev_block[::-1]
        result += self.merkle_root[::-1]
        result += int_to_little_endian(self.timestamp, 4)
        result += self.bits
        result += self.nonce
        return result

    def hash(self):
        '''Returns the double-sha256 interpreted little endian of the block'''
        return double_sha256(self.serialize())[::-1]

    def bip9(self):
        '''Returns whether this block is signaling readiness for BIP9'''
        # BIP9 is signalled if the top 3 bits are 001
        return self.version >> 29 == 0b001

    def bip91(self):
        '''Returns whether this block is signaling readiness for BIP91'''
        # BIP91 is signalled if the top 5th bit from the right is 1
        return (self.version >> 4) & 1 == 1
    
    def bip141(self):
        '''Returns whether this block is signaling readiness for BIP141'''
        # BIP91 is signalled if the top 2nd bit from the right is 1
        return (self.version >> 1) & 1 == 1

    def target(self):
        '''Returns the proof-of-work target based on the bits'''
        # note bits is in little-endian and the first byte is the exponent
        # the other three bytes are the coefficient.
        # the formula is:
        # coefficient * 2**(8*(exponent-3))
        exponent = self.bits[-1]
        coefficient = little_endian_to_int(self.bits[:-1])
        return coefficient * 2**(8*(exponent - 3))

    def difficulty(self):
        '''Returns the block difficulty based on the bits'''
        # note difficulty is (target of lowest difficulty) / (self's target)
        # lowest difficulty has bits that equal 0xffff001d
        exponent = 0x1d
        minimum_target = 0xffff * 2**(8*(0x1d-3))
        return minimum_target / self.target()

    def check_pow(self):
        '''Returns whether this block satisfies proof of work'''
        # You will need to get the hash of this block and interpret it
        # as an integer. If the hash of the block is lower, pow is good.
        # hint: int.from_bytes('', 'big')
        s256 = self.hash()
        return int.from_bytes(s256, 'big') < self.target()

    def validate_merkle_root(self):
        '''Gets the merkle root of the tx_hashes and checks that it's
        the same as the merkle root of this block.
        '''
        # reverse all the transaction hashes
        current_level = [x[::-1] for x in self.tx_hashes]
        # get the Merkle Root
        root = merkle_root(current_level)
        # check that this block's merkle root is the same as reverse of root
        return root[::-1] == self.merkle_root

    def calculate_merkle_tree(self):
        '''Calculate and store the entire Merkle Tree'''
        # store the result in self.merkle_tree, an array, 0 representing
        # the bottom level and 1 the parent level of level 0 and so on.
        self.merkle_tree = []
        if self.tx_hashes is None:
            raise RuntimeError('Transaction Hashes needed to calculate Merkle Tree')
        # reverse all the transaction hashes
        current = [x[::-1] for x in self.tx_hashes]
        # if there is more than 1 hash:
        while len(current) > 1:
            #   store current level
            self.merkle_tree.append(current)
            #   Make current level Merkle Parent level
            current = merkle_parent_level(current)
        # store root as the final level
        self.merkle_tree.append(current)
    
    def create_merkle_proof(self, tx_hash):
        if self.tx_hashes is None:
            return None
        elif self.merkle_tree is None:
            self.calculate_merkle_tree()
        index = self.merkle_tree[0].index(tx_hash[::-1])
        current = self.merkle_tree[0][index]
        proof_hashes = []
        for level, level_index in enumerate(merkle_path(index, len(self.tx_hashes))):
            if level_index % 2 == 0:
                partner = self.merkle_tree[level][level_index + 1]
                current = merkle_parent(current, partner)
            else:
                partner = self.merkle_tree[level][level_index - 1]
                current = merkle_parent(partner, current)
            proof_hashes.append(partner)
            
        # sanity check
        if current != self.merkle_tree[-1][0]:
            raise RuntimeError('merkle tree looks invalid')
        return Proof(self.merkle_root, tx_hash, index, proof_hashes)


class BlockTest(TestCase):

    def test_parse(self):
        block_raw = unhexlify('020000208ec39428b17323fa0ddec8e887b4a7c53b8c0a0a220cfd0000000000000000005b0750fce0a889502d40508d39576821155e9c9e3f5c3157f961db38fd8b25be1e77a759e93c0118a4ffd71d')
        stream = BytesIO(block_raw)
        block = Block.parse(stream)
        self.assertEqual(block.version, 0x20000002)
        want = unhexlify('000000000000000000fd0c220a0a8c3bc5a7b487e8c8de0dfa2373b12894c38e')
        self.assertEqual(block.prev_block, want)
        want = unhexlify('be258bfd38db61f957315c3f9e9c5e15216857398d50402d5089a8e0fc50075b')
        self.assertEqual(block.merkle_root, want)
        self.assertEqual(block.timestamp, 0x59a7771e)
        self.assertEqual(block.bits, unhexlify('e93c0118'))
        self.assertEqual(block.nonce, unhexlify('a4ffd71d'))

    def test_serialize(self):
        block_raw = unhexlify('020000208ec39428b17323fa0ddec8e887b4a7c53b8c0a0a220cfd0000000000000000005b0750fce0a889502d40508d39576821155e9c9e3f5c3157f961db38fd8b25be1e77a759e93c0118a4ffd71d')
        stream = BytesIO(block_raw)
        block = Block.parse(stream)
        self.assertEqual(block.serialize(), block_raw)

    def test_hash(self):
        block_raw = unhexlify('020000208ec39428b17323fa0ddec8e887b4a7c53b8c0a0a220cfd0000000000000000005b0750fce0a889502d40508d39576821155e9c9e3f5c3157f961db38fd8b25be1e77a759e93c0118a4ffd71d')
        stream = BytesIO(block_raw)
        block = Block.parse(stream)
        self.assertEqual(block.hash(), unhexlify('0000000000000000007e9e4c586439b0cdbe13b1370bdd9435d76a644d047523'))


    def test_bip9(self):
        block_raw = unhexlify('020000208ec39428b17323fa0ddec8e887b4a7c53b8c0a0a220cfd0000000000000000005b0750fce0a889502d40508d39576821155e9c9e3f5c3157f961db38fd8b25be1e77a759e93c0118a4ffd71d')
        stream = BytesIO(block_raw)
        block = Block.parse(stream)
        self.assertTrue(block.bip9())
        block_raw = unhexlify('0400000039fa821848781f027a2e6dfabbf6bda920d9ae61b63400030000000000000000ecae536a304042e3154be0e3e9a8220e5568c3433a9ab49ac4cbb74f8df8e8b0cc2acf569fb9061806652c27')
        stream = BytesIO(block_raw)
        block = Block.parse(stream)
        self.assertFalse(block.bip9())

    def test_bip91(self):
        block_raw = unhexlify('1200002028856ec5bca29cf76980d368b0a163a0bb81fc192951270100000000000000003288f32a2831833c31a25401c52093eb545d28157e200a64b21b3ae8f21c507401877b5935470118144dbfd1')
        stream = BytesIO(block_raw)
        block = Block.parse(stream)
        self.assertTrue(block.bip91())
        block_raw = unhexlify('020000208ec39428b17323fa0ddec8e887b4a7c53b8c0a0a220cfd0000000000000000005b0750fce0a889502d40508d39576821155e9c9e3f5c3157f961db38fd8b25be1e77a759e93c0118a4ffd71d')
        stream = BytesIO(block_raw)
        block = Block.parse(stream)
        self.assertFalse(block.bip91())

    def test_bip141(self):
        block_raw = unhexlify('020000208ec39428b17323fa0ddec8e887b4a7c53b8c0a0a220cfd0000000000000000005b0750fce0a889502d40508d39576821155e9c9e3f5c3157f961db38fd8b25be1e77a759e93c0118a4ffd71d')
        stream = BytesIO(block_raw)
        block = Block.parse(stream)
        self.assertTrue(block.bip141())
        block_raw = unhexlify('0000002066f09203c1cf5ef1531f24ed21b1915ae9abeb691f0d2e0100000000000000003de0976428ce56125351bae62c5b8b8c79d8297c702ea05d60feabb4ed188b59c36fa759e93c0118b74b2618')
        stream = BytesIO(block_raw)
        block = Block.parse(stream)
        self.assertFalse(block.bip141())

    def test_target(self):
        block_raw = unhexlify('020000208ec39428b17323fa0ddec8e887b4a7c53b8c0a0a220cfd0000000000000000005b0750fce0a889502d40508d39576821155e9c9e3f5c3157f961db38fd8b25be1e77a759e93c0118a4ffd71d')
        stream = BytesIO(block_raw)
        block = Block.parse(stream)
        self.assertEqual(block.target(), 0x13ce9000000000000000000000000000000000000000000)
        self.assertEqual(int(block.difficulty()), 888171856257)

    def test_check_pow(self):
        block_raw = unhexlify('04000000fbedbbf0cfdaf278c094f187f2eb987c86a199da22bbb20400000000000000007b7697b29129648fa08b4bcd13c9d5e60abb973a1efac9c8d573c71c807c56c3d6213557faa80518c3737ec1')
        stream = BytesIO(block_raw)
        block = Block.parse(stream)
        self.assertTrue(block.check_pow())
        block_raw = unhexlify('04000000fbedbbf0cfdaf278c094f187f2eb987c86a199da22bbb20400000000000000007b7697b29129648fa08b4bcd13c9d5e60abb973a1efac9c8d573c71c807c56c3d6213557faa80518c3737ec0')
        stream = BytesIO(block_raw)
        block = Block.parse(stream)
        self.assertFalse(block.check_pow())

    def test_validate_merkle_root(self):
        hashes_hex = [
            'f54cb69e5dc1bd38ee6901e4ec2007a5030e14bdd60afb4d2f3428c88eea17c1',
            'c57c2d678da0a7ee8cfa058f1cf49bfcb00ae21eda966640e312b464414731c1',
            'b027077c94668a84a5d0e72ac0020bae3838cb7f9ee3fa4e81d1eecf6eda91f3',
            '8131a1b8ec3a815b4800b43dff6c6963c75193c4190ec946b93245a9928a233d',
            'ae7d63ffcb3ae2bc0681eca0df10dda3ca36dedb9dbf49e33c5fbe33262f0910',
            '61a14b1bbdcdda8a22e61036839e8b110913832efd4b086948a6a64fd5b3377d',
            'fc7051c8b536ac87344c5497595d5d2ffdaba471c73fae15fe9228547ea71881',
            '77386a46e26f69b3cd435aa4faac932027f58d0b7252e62fb6c9c2489887f6df',
            '59cbc055ccd26a2c4c4df2770382c7fea135c56d9e75d3f758ac465f74c025b8',
            '7c2bf5687f19785a61be9f46e031ba041c7f93e2b7e9212799d84ba052395195',
            '08598eebd94c18b0d59ac921e9ba99e2b8ab7d9fccde7d44f2bd4d5e2e726d2e',
            'f0bb99ef46b029dd6f714e4b12a7d796258c48fee57324ebdc0bbc4700753ab1',
        ]
        hashes = [unhexlify(x) for x in hashes_hex]
        stream = BytesIO(unhexlify('00000020fcb19f7895db08cadc9573e7915e3919fb76d59868a51d995201000000000000acbcab8bcc1af95d8d563b77d24c3d19b18f1486383d75a5085c4e86c86beed691cfa85916ca061a00000000'))
        block = Block.parse(stream)
        block.tx_hashes = hashes
        self.assertTrue(block.validate_merkle_root())

    def test_calculate_merkle_tree(self):
        hashes_hex = [
            'f54cb69e5dc1bd38ee6901e4ec2007a5030e14bdd60afb4d2f3428c88eea17c1',
            'c57c2d678da0a7ee8cfa058f1cf49bfcb00ae21eda966640e312b464414731c1',
            'b027077c94668a84a5d0e72ac0020bae3838cb7f9ee3fa4e81d1eecf6eda91f3',
            '8131a1b8ec3a815b4800b43dff6c6963c75193c4190ec946b93245a9928a233d',
            'ae7d63ffcb3ae2bc0681eca0df10dda3ca36dedb9dbf49e33c5fbe33262f0910',
            '61a14b1bbdcdda8a22e61036839e8b110913832efd4b086948a6a64fd5b3377d',
            'fc7051c8b536ac87344c5497595d5d2ffdaba471c73fae15fe9228547ea71881',
            '77386a46e26f69b3cd435aa4faac932027f58d0b7252e62fb6c9c2489887f6df',
            '59cbc055ccd26a2c4c4df2770382c7fea135c56d9e75d3f758ac465f74c025b8',
            '7c2bf5687f19785a61be9f46e031ba041c7f93e2b7e9212799d84ba052395195',
            '08598eebd94c18b0d59ac921e9ba99e2b8ab7d9fccde7d44f2bd4d5e2e726d2e',
            'f0bb99ef46b029dd6f714e4b12a7d796258c48fee57324ebdc0bbc4700753ab1',
        ]
        hashes = [unhexlify(x) for x in hashes_hex]
        stream = BytesIO(unhexlify('00000020fcb19f7895db08cadc9573e7915e3919fb76d59868a51d995201000000000000acbcab8bcc1af95d8d563b77d24c3d19b18f1486383d75a5085c4e86c86beed691cfa85916ca061a00000000'))
        block = Block.parse(stream)
        block.tx_hashes = hashes
        block.calculate_merkle_tree()
        want0 = [
            'c117ea8ec828342f4dfb0ad6bd140e03a50720ece40169ee38bdc15d9eb64cf5',
            'c131474164b412e3406696da1ee20ab0fc9bf41c8f05fa8ceea7a08d672d7cc5',
            'f391da6ecfeed1814efae39e7fcb3838ae0b02c02ae7d0a5848a66947c0727b0',
            '3d238a92a94532b946c90e19c49351c763696cff3db400485b813aecb8a13181',
            '10092f2633be5f3ce349bf9ddbde36caa3dd10dfa0ec8106bce23acbff637dae',
            '7d37b3d54fa6a64869084bfd2e831309118b9e833610e6228adacdbd1b4ba161',
            '8118a77e542892fe15ae3fc771a4abfd2f5d5d5997544c3487ac36b5c85170fc',
            'dff6879848c2c9b62fe652720b8df5272093acfaa45a43cdb3696fe2466a3877',
            'b825c0745f46ac58f7d3759e6dc535a1fec7820377f24d4c2c6ad2cc55c0cb59',
            '95513952a04bd8992721e9b7e2937f1c04ba31e0469fbe615a78197f68f52b7c',
            '2e6d722e5e4dbdf2447ddecc9f7dabb8e299bae921c99ad5b0184cd9eb8e5908',
            'b13a750047bc0bdceb2473e5fe488c2596d7a7124b4e716fdd29b046ef99bbf0',
        ]
        want1 = [
            '8b30c5ba100f6f2e5ad1e2a742e5020491240f8eb514fe97c713c31718ad7ecd',
            '7f4e6f9e224e20fda0ae4c44114237f97cd35aca38d83081c9bfd41feb907800',
            'ade48f2bbb57318cc79f3a8678febaa827599c509dce5940602e54c7733332e7',
            '68b3e2ab8182dfd646f13fdf01c335cf32476482d963f5cd94e934e6b3401069',
            '43e7274e77fbe8e5a42a8fb58f7decdb04d521f319f332d88e6b06f8e6c09e27',
            '4f492e893bf854111c36cb5eff4dccbdd51b576e1cfdc1b84b456cd1c0403ccb',
        ]
        want2 = [
            '26906cb2caeb03626102f7606ea332784281d5d20e2b4839fbb3dbb37262dbc1',
            '717a0d17538ff5ad2c020bab38bdcde66e63f3daef88f89095f344918d5d4f96',
            'd20629030c7e48e778c1c837d91ebadc2f2ee319a0a0a461f4a9538b5cae2a69',
            'd20629030c7e48e778c1c837d91ebadc2f2ee319a0a0a461f4a9538b5cae2a69',
        ]
        want3 = [
            'b9f5560ce9630ea4177a7ac56d18dea73c8f76b59e02ab4805eaeebd84a4c5b1',
            '00aa9ad6a7841ffbbf262eb775f8357674f1ea23af11c01cfb6d481fec879701',
        ]
        want4 = [
            'acbcab8bcc1af95d8d563b77d24c3d19b18f1486383d75a5085c4e86c86beed6',
        ]
        self.assertEqual(block.merkle_tree[0], [unhexlify(x) for x in want0])
        self.assertEqual(block.merkle_tree[1], [unhexlify(x) for x in want1])
        self.assertEqual(block.merkle_tree[2], [unhexlify(x) for x in want2])
        self.assertEqual(block.merkle_tree[3], [unhexlify(x) for x in want3])
        self.assertEqual(block.merkle_tree[4], [unhexlify(x) for x in want4])

    def test_create_merkle_proof(self):
        hashes_hex = [
            'f54cb69e5dc1bd38ee6901e4ec2007a5030e14bdd60afb4d2f3428c88eea17c1',
            'c57c2d678da0a7ee8cfa058f1cf49bfcb00ae21eda966640e312b464414731c1',
            'b027077c94668a84a5d0e72ac0020bae3838cb7f9ee3fa4e81d1eecf6eda91f3',
            '8131a1b8ec3a815b4800b43dff6c6963c75193c4190ec946b93245a9928a233d',
            'ae7d63ffcb3ae2bc0681eca0df10dda3ca36dedb9dbf49e33c5fbe33262f0910',
            '61a14b1bbdcdda8a22e61036839e8b110913832efd4b086948a6a64fd5b3377d',
            'fc7051c8b536ac87344c5497595d5d2ffdaba471c73fae15fe9228547ea71881',
            '77386a46e26f69b3cd435aa4faac932027f58d0b7252e62fb6c9c2489887f6df',
            '59cbc055ccd26a2c4c4df2770382c7fea135c56d9e75d3f758ac465f74c025b8',
            '7c2bf5687f19785a61be9f46e031ba041c7f93e2b7e9212799d84ba052395195',
            '08598eebd94c18b0d59ac921e9ba99e2b8ab7d9fccde7d44f2bd4d5e2e726d2e',
            'f0bb99ef46b029dd6f714e4b12a7d796258c48fee57324ebdc0bbc4700753ab1',
        ]
        hashes = [unhexlify(x) for x in hashes_hex]
        stream = BytesIO(unhexlify('00000020fcb19f7895db08cadc9573e7915e3919fb76d59868a51d995201000000000000acbcab8bcc1af95d8d563b77d24c3d19b18f1486383d75a5085c4e86c86beed691cfa85916ca061a00000000'))
        block = Block.parse(stream)
        block.tx_hashes = hashes
        h = hashes[7]
        proof = block.create_merkle_proof(h)
        self.assertEqual(proof.index, 7)
        want = [
            '8118a77e542892fe15ae3fc771a4abfd2f5d5d5997544c3487ac36b5c85170fc',
            'ade48f2bbb57318cc79f3a8678febaa827599c509dce5940602e54c7733332e7',
            '26906cb2caeb03626102f7606ea332784281d5d20e2b4839fbb3dbb37262dbc1',
            '00aa9ad6a7841ffbbf262eb775f8357674f1ea23af11c01cfb6d481fec879701',
        ]
        self.assertEqual(proof.merkle_proof, [unhexlify(x) for x in want])

    def test_verify_merkle_proof(self):
        merkle_root = unhexlify('d6ee6bc8864e5c08a5753d3886148fb1193d4cd2773b568d5df91acc8babbcac')
        tx_hash = unhexlify('77386a46e26f69b3cd435aa4faac932027f58d0b7252e62fb6c9c2489887f6df')
        index = 7
        proof_hex_hashes = [
            '8118a77e542892fe15ae3fc771a4abfd2f5d5d5997544c3487ac36b5c85170fc',
            'ade48f2bbb57318cc79f3a8678febaa827599c509dce5940602e54c7733332e7',
            '26906cb2caeb03626102f7606ea332784281d5d20e2b4839fbb3dbb37262dbc1',
            '00aa9ad6a7841ffbbf262eb775f8357674f1ea23af11c01cfb6d481fec879701',
        ]
        proof_hashes = [unhexlify(x) for x in proof_hex_hashes]
        proof = Proof(merkle_root=merkle_root, tx_hash=tx_hash, index=index, merkle_proof=proof_hashes)
        self.assertTrue(proof.verify())
