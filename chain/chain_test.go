package chain

import (
	"context"
	"testing"

	"github.com/lance4117/blogd/api/blog/blog"
	blogmodules "github.com/lance4117/blogd/x/blog/module"
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
}

func TestBroadcastTx(t *testing.T) {
	config := DefaultConfig("blog", "D:\\code\\blogd\\blogdata", blogmodules.AppModule{})

	client, err := New(config)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	senderName := "alice"

	msg := blog.MsgCreateBlog{
		Title:   "title",
		Content: "content",
	}

	response, err := client.BroadcastTx(context.Background(), senderName, &msg)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(response.TxResponse)

}
