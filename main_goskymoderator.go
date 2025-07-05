package main

import (
	"bskymoderator/config"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/api/bsky"
	lexutil "github.com/bluesky-social/indigo/lex/util"
	"github.com/bluesky-social/indigo/xrpc"
)

var _client *xrpc.Client

type BskyParam struct {
	UserId    string
	UserDid   string
	Password  string
	ListId    string
	ListAtUri string
	Query     string
}

func main() {

	conf, err := initializeConfig()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Did :", conf.UserDid)
	fmt.Println("AtUri :", conf.ListAtUri)

	if err := doMain(); err != nil {
		fmt.Println(err)
	}
}

func doMain() error {

	ctx := context.Background()
	registered := fetchExistingHandles(ctx)

	time.Sleep(2 * time.Second)

	cursor := ""
	count := 0
	total := 0
	for {
		resp, err := searchActors(ctx, cursor)
		if err != nil {
			return fmt.Errorf("検索失敗: %v", err)
		}

		for _, user := range resp.Actors {
			total++
			if registered[user.Did] {
				log.Printf("⚠️   %d スキップ（既登録）: %s", total, user.Handle)
				continue
			}

			err = register(ctx, user)
			if err != nil {
				log.Printf("❌  %d 登録失敗: %s (%v)", total, user.Handle, err)
			} else {
				log.Printf("✅  %d 登録成功: %s", total, user.Handle)
				count++
			}
		}

		if resp.Cursor == nil || *resp.Cursor == "" {
			break
		}
		cursor = *resp.Cursor
	}

	fmt.Printf("🎉 新規登録数: %d 件\n", count)

	return nil
}

func fetchExistingHandles(ctx context.Context) map[string]bool {
	conf := config.Instance()
	client := getClient(ctx)
	existing := make(map[string]bool)
	cursor := ""
	limit := int64(100)
	total := 1
	for {
		resp, err := atproto.RepoListRecords(ctx, client, "app.bsky.graph.listitem", cursor, limit, conf.UserDid, false)
		fmt.Print("\rfetch ExistingHandles ... ", total*100, "                     ")
		total++
		if err != nil {
			log.Fatalf("リスト項目取得失敗: %v", err)
		}

		for _, rec := range resp.Records {
			item := new(bsky.GraphListitem)

			raw, err := rec.Value.MarshalJSON()
			if err != nil {
				log.Printf("⚠️ MarshalJSON 失敗: %v", err)
				continue
			}
			if err := json.Unmarshal(raw, item); err != nil {
				log.Printf("⚠️ json.Unmarshal 失敗: %v", err)
				continue
			}

			if item.List != conf.ListAtUri {
				continue
			}
			existing[item.Subject] = true
		}

		if resp.Cursor == nil || *resp.Cursor == "" {
			break
		}
		cursor = *resp.Cursor
	}

	fmt.Println("")

	fmt.Printf("✅ 既存登録: %d 件取得\n", len(existing))

	return existing
}

func getClient(ctx context.Context) *xrpc.Client {

	if isValidSession(ctx) {
		return _client
	}

	conf := config.Instance()

	handle := conf.UserId
	appPassword := conf.Password

	xrpClient := &xrpc.Client{Host: "https://bsky.social"}
	sess, err := atproto.ServerCreateSession(ctx, xrpClient, &atproto.ServerCreateSession_Input{
		Identifier: handle,
		Password:   appPassword,
	})
	if err != nil {
		log.Fatal("ログイン失敗:", err)
	}
	xrpClient.Auth = &xrpc.AuthInfo{
		AccessJwt:  sess.AccessJwt,
		RefreshJwt: sess.RefreshJwt,
		Did:        sess.Did,
		Handle:     sess.Handle,
	}

	_client = xrpClient

	return _client
}

func searchActors(ctx context.Context, cursor string) (*bsky.ActorSearchActors_Output, error) {
	conf := config.Instance()
	xrpClient := getClient(ctx)
	resp, err := bsky.ActorSearchActors(ctx, xrpClient, cursor, 100, conf.Query, "")
	return resp, err
}

func register(ctx context.Context, user *bsky.ActorDefs_ProfileView) error {
	conf := config.Instance()
	client := getClient(ctx)
	_, err := atproto.RepoCreateRecord(ctx, client, &atproto.RepoCreateRecord_Input{
		Repo:       client.Auth.Did,
		Collection: "app.bsky.graph.listitem",
		Record: &lexutil.LexiconTypeDecoder{
			Val: &bsky.GraphListitem{
				Subject:   user.Did,
				List:      conf.ListAtUri,
				CreatedAt: time.Now().Format(time.RFC3339),
			},
		},
	})

	return err
}

func isValidSession(ctx context.Context) bool {
	if _client == nil {
		return false
	}

	return true

	// session, err := atproto.ServerGetSession(ctx, _client)
	// if err != nil {
	// 	fmt.Println("セション切れ", err)
	// }
	// return err == nil
}

func getDid(ctx context.Context) (string, error) {
	client := getClient(ctx)
	session, err := atproto.ServerGetSession(ctx, client)
	if err != nil {
		return "", err
	}

	return session.Did, nil
}

func initializeConfig() (*config.Config, error) {
	conf, err := config.InitializeConfig()
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	did, err := getDid(ctx)
	if err != nil {
		return nil, err
	}
	atUri := "at://" + did + "/app.bsky.graph.list/" + conf.ListId

	conf.UserDid = did
	conf.ListAtUri = atUri

	return conf, err
}
