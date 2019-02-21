package message

const (
	MSG_PUB_BATCH     = 'a'
	MSG_PUB_ONE       = 'b'
	MSG_PUB_TIMER     = 'c'
	MSG_PUB_TIMER_ACK = 'd'
	MSG_PUB_RESTORE   = 'e'

	MSG_SUB    = 'f'
	MSG_SUBACK = 'g'
	MSG_UNSUB  = 'h'

	MSG_PING = 'i'
	MSG_PONG = 'j'

	MSG_COUNT = 'k'
	MSG_PULL  = 'l'

	MSG_CONNECT    = 'm'
	MSG_CONNECT_OK = 'n'

	MSG_BROADCAST = 'o'

	MSG_REDUCE_COUNT = 'p'
	MSG_MARK_READ    = 'r'

	MSG_JOIN_CHAT        = 's'
	MSG_LEAVE_CHAT       = 't'
	MSG_PRESENCE_ONLINE  = 'u'
	MSG_PRESENCE_OFFLINE = 'v'
	MSG_PRESENCE_ALL     = 'q'
	MSG_ALL_CHAT_USERS   = 'w'

	MSG_RETRIEVE = 'x'
)

const (
	NORMAL_MSG = 0
	TIMER_MSG  = 1
)

const (
	QOS0 = 0
	QOS1 = 1
)

var (
	DEFAULT_QUEUE     = []byte("meq.io")
	MSG_NEWEST_OFFSET = []byte("0")
)

const (
	MAX_PULL_COUNT = 100

	CacheFlushLen = 200

	REDUCE_ALL_COUNT = 0

	MAX_IDLE_TIME = 120

	NeverExpires = 0

	MSG_ID_LENGTH = 19
)
