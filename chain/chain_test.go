package chain

import (
	"context"
	"testing"
	"time"

	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/lance4117/blogd/api/blog/blog"
	blogmodules "github.com/lance4117/blogd/x/blog/module"
	"github.com/lance4117/gofuse/gen"
)

func TestNewClient(t *testing.T) {
	// blogmodules作为演示的区块链程序
	config := DefaultConfig("blog", "D:\\code\\blogd\\blogdata", blogmodules.AppModule{})

	client, err := New(config)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	status, err := client.Status()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(config)
	t.Log(status)
	height, err := client.LatestHeight(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	t.Log(height)
}

func TestBroadcastTx(t *testing.T) {
	config := DefaultConfig("blog", "D:\\code\\blogd\\blogdata", blogmodules.AppModule{})
	config.BroadcastMode = txtypes.BroadcastMode_BROADCAST_MODE_SYNC
	config.GasLimit = 200000000

	client, err := New(config)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	senderName := "alice"

	address, err := client.Address(senderName)
	if err != nil {
		t.Fatal(err)
	}

	// Creator 必填
	msg := blog.MsgCreateBlog{
		Creator: address.String(),
		Title:   "title",
		Content: gen.NewArticle(1, 1, 100),
	}

	response, err := client.BroadcastTx(context.Background(), senderName, &msg)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(response.TxResponse.String())

	tx, err := client.WaitForTx(context.Background(), response.TxResponse.TxHash, time.Second*15, 1*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(tx.String())
}
