Why is TCP Reliable?
- It overcomes the effects of packet loss or receiving packets out of order.
- Packet loss usually occurs due to data transmission errors, such as wireless
network interference or network congestion.
- TCP adapts transfer rates to transmit data as quickly as possible while
minimizing packet loss. This process is called flow control.
- TCP keeps track of received packets and retransmits unacknowledged packets.
Packets are not guaranteed to take the same route and may not be in the correct
order; however, TCP will process them sequentially.

TCP Sessions
- Allows one node to send a stream of data to another while receiving feedback
in real time. This is so that errors can be corrected right away, rather than
after the entire stream is sent.

TCP Handshake
- The TCP handshake creates a TCP session over which two nodes can exchange data.
- Beforehand, the server must listen for incoming connections.
- Step 1: The client sends a packet to the server with the synchronize (SYN)
flag, which informs it of the client's capabilities and window settings for the
session.
- Step 2: The server responds with its own packet, with both the acknowledgement
(ACK) and the SYN flag set. ACK lets the client know that the server received
its SYN packet. The SYN packet that the server sends includes what settings
it's agreed to for the rest of the session.
- Step 3: The client sends and ACK packet to acknowledge the server's SYN
packet.
- After these three steps, the TCP session has been established, and the nodes
can now exchange data. The session will remain idle until either side has
data to transmit, which may result in wasted memory if the session is unmanaged.
- Things sent in the ACK packet:
  - Sequence number X that tells sender that all packets up to and including
  the packet with sequence number X have been received.
  - A receive buffer, which is a fixed-size block of memory that a node
  can accept without requiring another acknowledgement. The number of bytes
  remaining in the buffer is called the window size.

Terminating a TCP Session
- Either side of the connection can initiate the termination sequence by
sending a finish (FIN) packet. At this point, the sender's connection state
switches from ESTABLISHED to FIN_WAIT_1. The server sends back an ACK packet
to acknowledge this and switches its state from ESTABLISHED to CLOSE_WAIT. Then, the client's state switches to FIN_WAIT_2 and the server sends its own FIN
packet to the client and switches to LAST_ACK, indicating that it is waiting
for a final acknowledgement. At this point, the client switches to the TIME_WAIT state and sends back an ACK packet. Upon receiving this packet, the
server finally goes into the CLOSED state, and once the timer is done on the
client end, it also goes into the CLOSED state.

TCP in Go
- If the TCP handshake succeeds, a connection object will be returned.
- You can set read and write deadlines for a connection by using conn.SetReadDeadline, conn.SetWriteDeadline, or conn.SetDeadline for both. When a connection reaches its deadline, all currently blocked and future calls to the Read/Write
method immediately return a timeout error. This can make sure that a network connection is not idle for a long time. Exceeded deadlines do not mean that the session is terminated, you can always call the set methods again and reads/writes
will work again.
- For long-running network connections that may experience exteended idle periods at the application level, it is wise to implement a heartbeat between nodes to advance the deadlline. A heartbeat is a message sent to the remote side with the intention of eliciting a reply, which we can use to advance the deadlline like a heartbeat.