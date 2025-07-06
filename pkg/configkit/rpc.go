package configkit

type RPC struct {
	ListenPort string `validate:"required"`
}

func InitRPC(c *C) RPC {
	return RPC{
		ListenPort: c.Viper.GetString("RPC_LISTEN_PORT"),
	}
}
