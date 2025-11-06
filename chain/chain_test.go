package chain

import (
	"context"
	"testing"
	"time"

	"github.com/lance4117/blogd/api/blog/blog"
	blogmodules "github.com/lance4117/blogd/x/blog/module"
	"github.com/lance4117/gofuse/gen"
)

const HomeDir = "D:\\000_MyLibrary\\100_code\\142_blogd\\blogdata"

func TestNewClient(t *testing.T) {
	// blogmodules作为演示的区块链程序
	config := DefaultConfig("blog", HomeDir, blogmodules.AppModule{})

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
	config := DefaultConfig("blog", HomeDir, blogmodules.AppModule{})

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
		Content: gen.Sentence(100),
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

func TestBank(t *testing.T) {
	config := DefaultConfig("blog", HomeDir, blogmodules.AppModule{})

	client, err := New(config)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	aliceAddress, err := client.Address("alice")
	if err != nil {
		t.Fatal(err)
	}

	acc1Address, err := client.Address("acc1")
	if err != nil {
		t.Fatal(err)
	}

	coins, err := client.SendCoins(context.Background(), "alice", aliceAddress, acc1Address, "10stake")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(coins.String())

	// 等待交易上链
	tx, err := client.WaitForTx(context.Background(), coins.TxResponse.TxHash, time.Second*15, 1*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(tx.String())

	balance, err := client.Balance(context.Background(), aliceAddress, "stake")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(balance.String())
}
