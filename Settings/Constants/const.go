package Constants

const (
	TCPHost             = "127.0.0.1"
	TCPPort             = "6060"
	
	ClientServerPort    = "8080"

	ClientTokenSecret   = "Hoolalalaleleho"
	TokenExpiryDuration = "5m"

	TokenSecret         = "YayThisIsServerSecret"

	//Key used to decrypt excryption key sent by any server
	RSAPrivateKeyFile   = "rsa_private.pem"

	//Channel size for each channel in websocket messaging queue
	WSChannelSize = 1024
	ProcessesChannelSize = 4096
	FileListChannelSize = 4096
	RequestChannelSize = 4096
)
