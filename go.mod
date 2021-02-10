module github.com/public-awesome/stargaze

go 1.15

require (
	github.com/CosmWasm/wasmd v0.15.1-0.20210205121914-d0befd973db1
	github.com/bwmarrin/snowflake v0.3.0
	github.com/cosmos/cosmos-sdk v0.41.0
	github.com/gogo/protobuf v1.3.3
	github.com/golang/protobuf v1.4.3
	github.com/gorilla/mux v1.8.0
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/rakyll/statik v0.1.7
	github.com/spf13/cast v1.3.1
	github.com/spf13/cobra v1.1.1
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.7.0
	github.com/tendermint/liquidity v0.1.0-rc0.0.20210206010652-2ca7f61d0f15
	github.com/tendermint/tendermint v0.34.3
	github.com/tendermint/tm-db v0.6.3
	google.golang.org/genproto v0.0.0-20210114201628-6edceaf6022f
	google.golang.org/grpc v1.35.0
	gopkg.in/yaml.v2 v2.4.0
)

replace google.golang.org/grpc => google.golang.org/grpc v1.33.2

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
