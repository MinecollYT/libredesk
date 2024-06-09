package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/abhinavxd/artemis/internal/attachment"
	"github.com/abhinavxd/artemis/internal/cannedresp"
	"github.com/abhinavxd/artemis/internal/contact"
	"github.com/abhinavxd/artemis/internal/conversation"
	convtag "github.com/abhinavxd/artemis/internal/conversation/tag"
	"github.com/abhinavxd/artemis/internal/inbox"
	"github.com/abhinavxd/artemis/internal/initz"
	"github.com/abhinavxd/artemis/internal/message"
	"github.com/abhinavxd/artemis/internal/message/models"
	"github.com/abhinavxd/artemis/internal/tag"
	"github.com/abhinavxd/artemis/internal/team"
	"github.com/abhinavxd/artemis/internal/user"
	"github.com/abhinavxd/artemis/internal/ws"
	"github.com/knadh/koanf/v2"
	"github.com/valyala/fasthttp"
	"github.com/vividvilla/simplesessions"
	"github.com/zerodha/fastglue"
	"github.com/zerodha/logf"
)

var ko = koanf.New(".")

// App is the global app context which is passed and injected in the http handlers.
type App struct {
	constants           consts
	lo                  *logf.Logger
	cntctMgr            *contact.Manager
	userMgr             *user.Manager
	teamMgr             *team.Manager
	sessMgr             *simplesessions.Manager
	tagMgr              *tag.Manager
	msgMgr              *message.Manager
	inboxMgr            *inbox.Manager
	attachmentMgr       *attachment.Manager
	cannedRespMgr       *cannedresp.Manager
	conversationMgr     *conversation.Manager
	conversationTagsMgr *convtag.Manager
}

func main() {
	// Load command line flags into Koanf.
	initFlags()

	// Load the config file into Koanf.
	initz.Config(ko)

	var (
		shutdownCh = make(chan struct{})
		ctx, stop  = signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

		// Incoming messages from all inboxes are pushed to this queue.
		incomingMsgQ = make(chan models.IncomingMessage, ko.MustInt("message.incoming_queue_size"))

		lo = initz.Logger(ko.MustString("app.log_level"), ko.MustString("app.env"), "artemis")
		rd = initz.Redis(ko)
		db = initz.DB(ko)

		attachmentMgr   = initAttachmentsManager(db, &lo)
		cntctMgr        = initContactManager(db, &lo)
		conversationMgr = initConversations(db, &lo)
		inboxMgr        = initInboxManager(db, &lo, incomingMsgQ)

		// Websocket hub.
		wsHub = ws.NewHub()
	)

	// Init the app.
	var app = &App{
		lo:                  &lo,
		cntctMgr:            cntctMgr,
		inboxMgr:            inboxMgr,
		attachmentMgr:       attachmentMgr,
		conversationMgr:     conversationMgr,
		constants:           initConstants(),
		msgMgr:              initMessages(db, &lo, incomingMsgQ, wsHub, cntctMgr, attachmentMgr, conversationMgr, inboxMgr),
		sessMgr:             initSessionManager(rd),
		userMgr:             initUserDB(db, &lo),
		teamMgr:             initTeamMgr(db, &lo),
		tagMgr:              initTags(db, &lo),
		cannedRespMgr:       initCannedResponse(db, &lo),
		conversationTagsMgr: initConversationTags(db, &lo),
	}

	// Start receivers for all active inboxes.
	inboxMgr.Receive()

	// Start message inserter and dispatchers.
	go app.msgMgr.StartDBInserts(ctx, ko.MustInt("message.reader_concurrency"))
	go app.msgMgr.StartDispatcher(ctx, ko.MustInt("message.dispatch_concurrency"), ko.MustDuration("message.dispatch_read_interval"))

	// Init fastglue http server.
	g := fastglue.NewGlue()

	// Add app the request context.
	g.SetContext(app)

	// Init the handlers.
	initHandlers(g, wsHub)

	s := &fasthttp.Server{
		Name:                 ko.MustString("app.server.name"),
		ReadTimeout:          ko.MustDuration("app.server.read_timeout"),
		WriteTimeout:         ko.MustDuration("app.server.write_timeout"),
		MaxRequestBodySize:   ko.MustInt("app.server.max_body_size"),
		MaxKeepaliveDuration: ko.MustDuration("app.server.keepalive_timeout"),
		ReadBufferSize:       ko.MustInt("app.server.max_body_size"),
	}

	// Goroutine for handling interrupt signals & gracefully shutting down the server.
	go func() {
		<-ctx.Done()
		shutdownCh <- struct{}{}
		stop()
	}()

	// Start the HTTP server.
	log.Printf("server listening on %s %s", ko.String("app.server.address"), ko.String("app.server.socket"))
	if err := g.ListenServeAndWaitGracefully(ko.String("app.server.address"), ko.String("server.socket"), s, shutdownCh); err != nil {
		log.Fatalf("error starting frontend server: %v", err)
	}
	log.Println("bye")
}
