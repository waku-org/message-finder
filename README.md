# Query tool

Use this to query a storenode
```
make
```

To retrieve data using StoreV2
```
./build/query \
  --cluster-id=16 \
  --storenode=/dns4/store-01.do-ams3.shards.test.status.im/tcp/30303/p2p/16Uiu2HAmAUdrQ3uwzuE4Gy4D56hX6uLKEeerJAnhKEHZ3DxF1EfT \
  --pubsub-topic=/waku/2/rs/16/32 \
  --content-topic=/waku/1/0x242ed557/rfc26 \
  --content-topic=/waku/1/0xd811cd50/rfc26 \
  --content-topic=/waku/1/0x89bed93d/rfc26 \
  --content-topic=/waku/1/0xc95d2429/rfc26 \
  --content-topic=/waku/1/0xa0a6b41b/rfc26 \
  --start-time=1717507412000000000 \
  --end-time=1717593812000000000 \
  --pagesize=20 \
  --use-legacy=true


Page: 1, Record from 1 to 20
MessageHash                                                         Content Topic             Timestamp                              Unix Nanoseconds  
0x5142e23971e369c61c1fce403e28728e1e2bf950ea4988032aeab527f785a7bf  /waku/2/rs/16/32  /waku/1/0xd811cd50/rfc26  1717587985332236000  2024-06-05 11:46:25.332236 +0000 UTC
0x1f09e1534ff19ab2c3f4e9cc751d621ce66f86df2b6eb8efa2712edfa0162fde  /waku/2/rs/16/32  /waku/1/0x242ed557/rfc26  1717588588587729000  2024-06-05 11:56:28.587729 +0000 UTC     
...
...

Cursor: Digest(f4e9ceb8cb71b5f59e788858c1133d6853fbcf11094d80da53db235843351509); ReceiverTime:1717586544149237000, SenderTime: 1717586544149237000, PubsubTopic: /waku/2/rs/16/32

Page: 2, Record from 21 to 30
MessageHash                                                         Content Topic             Timestamp                              Unix Nanoseconds  
0x2dfc7a3bcfe70590b62c18f551db94a810561a14890d950adcb62f3a1b62868d  /waku/2/rs/16/32  /waku/1/0x242ed557/rfc26  1717588617605728000  2024-06-05 11:56:57.605728 +0000 UTC     
0x688bc1fcf806baa978fe6f59483b8acdcb1b7a19f920c98002d422b44fd6301f  /waku/2/rs/16/32  /waku/1/0x242ed557/rfc26  1717588680307133000  2024-06-05 11:58:00.307133 +0000 UTC
...
...
```

To retrieve data using StoreV3
```
# Using filter criteria
./build/query \
  --cluster-id=16 \
  --storenode=/dns4/store-01.do-ams3.shards.test.status.im/tcp/30303/p2p/16Uiu2HAmAUdrQ3uwzuE4Gy4D56hX6uLKEeerJAnhKEHZ3DxF1EfT \
  --pubsub-topic=/waku/2/rs/16/32 \
  --content-topic=/waku/1/0x242ed557/rfc26 \
  --content-topic=/waku/1/0xd811cd50/rfc26 \
  --content-topic=/waku/1/0x89bed93d/rfc26 \
  --content-topic=/waku/1/0xc95d2429/rfc26 \
  --content-topic=/waku/1/0xa0a6b41b/rfc26 \
  --start-time=1717507412000000000 \
  --end-time=1717593812000000000 \
  --pagesize=20

Page: 1, Record from 1 to 20
MessageHash                                                         Content Topic             Timestamp                                                     
0xf9becf36390b0790f4556829bab9fe177494fc3fd9309ceb6bc0e7d869e96485  /waku/1/0x242ed557/rfc26  1717710661383564000  2024-06-06 21:51:01.383564 +0000 UTC     
0x5b555b7986f6f76dfc199934b159327ed7535232012c0849bcc67ec1d2dacedf  /waku/1/0x242ed557/rfc26  1717710663235488000  2024-06-06 21:51:03.235488 +0000 UTC  
...
...

Cursor: acb469e0464aa6ebe7847807bf856c05a7ed3d26c5813f76c206ff9706eb686b


# Using message hashes
./build/query \
  --cluster-id=16 \
  --storenode=/dns4/store-01.do-ams3.shards.test.status.im/tcp/30303/p2p/16Uiu2HAmAUdrQ3uwzuE4Gy4D56hX6uLKEeerJAnhKEHZ3DxF1EfT \
  --hash=0xacb469e0464aa6ebe7847807bf856c05a7ed3d26c5813f76c206ff9706eb686b \
  --hash=0x518dd28a733ce491cedf698e4ed375e362fd435603f7d8807dc341303e838f3f

Page: 1, Record from 1 to 2
MessageHash                                                         PubsubTopic       Content Topic             Timestamp                                                  
0xacb469e0464aa6ebe7847807bf856c05a7ed3d26c5813f76c206ff9706eb686b  /waku/2/rs/16/32  /waku/1/0xd811cd50/rfc26  1717586544149237000  2024-06-05 11:22:24.149237 +0000 UTC  
0x518dd28a733ce491cedf698e4ed375e362fd435603f7d8807dc341303e838f3f  /waku/2/rs/16/32  /waku/1/0x242ed557/rfc26  1717586545164878000  2024-06-05 11:22:25.164878 +0000 UTC 

```

To see the content of a message
```
 ./build/query \
   --cluster-id=16 \
   --storenode=/dns4/store-01.do-ams3.shards.test.status.im/tcp/30303/p2p/16Uiu2HAmAUdrQ3uwzuE4Gy4D56hX6uLKEeerJAnhKEHZ3DxF1EfT \
   --hash 0xacb469e0464aa6ebe7847807bf856c05a7ed3d26c5813f76c206ff9706eb686b

PubsubTopic: /waku/2/rs/16/32
MessageHash: 0x518dd28a733ce491cedf698e4ed375e362fd435603f7d8807dc341303e838f3f
ContentTopic: /waku/1/0x242ed557/rfc26
Timestamp: 2024-06-05 11:22:25.164878 +0000 UTC
Timestamp (unix nano): 1717586545164878000
Version: 1
Payload:
00000000  c4 33 ca 2a 86 20 35 e8  5c 9b 3b aa 1d d4 e1 25  |.3.*. 5.\.;....%|
00000010  52 e4 62 8b 84 94 ca 1a  4a 02 4a e6 11 39 0c 99  |R.b.....J.J..9..|
00000020  12 aa 28 ae 70 0c b1 f2  31 e4 a1 10 ee 0f c8 6d  |..(.p...1......m|
00000030  c7 28 2e 75 c8 a0 a4 21  19 9a ee e1 07 5b 41 7f  |.(.u...!.....[A.|
00000040  bf f7 19 cd a9 f4 54 08  39 3f 55 2e ed 79 55 df  |......T.9?U..yU.|
00000050  77 cf fc f8 49 c3 04 c2  bf 77 b4 ce b9 95 a8 56  |w...I....w.....V|
00000060  46 06 c4 89 be 15 09 89  8b a5 06 30 90 96 36 8b  |F..........0..6.|
00000070  96 f8 df c3 84 c4 58 93  8f 76 58 d0 33 bb 14 cd  |......X..vX.3...|
00000080  8e 59 9b 51 1f 97 ab 8c  07 fe 3a ff f7 e7 6f fb  |.Y.Q......:...o.|
00000090  7f d0 aa 56 03 ca 49 eb  ef 08 9a ef ca 12 85 72  |...V..I........r|
000000a0  b0 47 79 2b 28 50 06 7d  89 69 d8 85 8e 3c 7f cb  |.Gy+(P.}.i...<..|
Meta: <nil>
RateLimitProof: <nil>
```

### Docker
```
# Build
docker build -t storequery:latest .

# Execute
docker run storequery:latest \
  --cluster-id=16 \
  --storenode=/dns4/store-01.do-ams3.shards.test.status.im/tcp/30303/p2p/16Uiu2HAmAUdrQ3uwzuE4Gy4D56hX6uLKEeerJAnhKEHZ3DxF1EfT \
  --pubsub-topic=/waku/2/rs/16/32 \
  --content-topic=/waku/1/0x242ed557/rfc26 \
  --content-topic=/waku/1/0xd811cd50/rfc26 \
  --content-topic=/waku/1/0x89bed93d/rfc26 \
  --content-topic=/waku/1/0xc95d2429/rfc26 \
  --content-topic=/waku/1/0xa0a6b41b/rfc26 \
  --start-time=1717507412000000000 \
  --end-time=1717593812000000000 \
  --pagesize=20 \
  --use-legacy=true
```
