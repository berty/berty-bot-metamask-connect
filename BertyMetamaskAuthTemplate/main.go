package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"syscall"

	"berty.tech/berty/v2/go/pkg/bertybot"
	"berty.tech/berty/v2/go/pkg/bertyversion"
	"berty.tech/berty/v2/go/pkg/messengertypes"

	qrterminal "github.com/mdp/qrterminal/v3"
	"github.com/oklog/run"
	"github.com/peterbourgon/ff/v3"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"

	"moul.io/climan"
	"moul.io/motd"
	"moul.io/srand"
	"moul.io/zapconfig"
)

func main() {
	if err := mainRun(os.Args[1:]); err != nil {
		if !errors.Is(err, flag.ErrHelp) {
			fmt.Fprintf(os.Stderr, "error: %v+\n", err)
		}
		os.Exit(1)
	}
}

var opts struct { // nolint:maligned
	Debug          bool
	BertyNodeAddr  string
	authAddr       string
	rootLogger     *zap.Logger
	publickeyPath  string
	privatekeyPath string
	generateKeys   bool
}

func mainRun(args []string) error {
	// parse CLI
	name := os.Args[0]
	root := &climan.Command{
		Name:       name,
		ShortUsage: name + " [global flags] <subcommand> [flags] [args]",
		ShortHelp:  "More info on https://github.com/pmg-tools/berty-bot-experiments.",
		FlagSetBuilder: func(fs *flag.FlagSet) {
			// opts.BertyNodeAddr = ""
			fs.BoolVar(&opts.Debug, "debug", false, "debug mode")
			fs.StringVar(&opts.BertyNodeAddr, "berty-node-addr", "127.0.0.1:9091", "Berty node address")
			fs.StringVar(&opts.authAddr, "authAddress", "http://localhost:3000/", "metamask oauth address")
		},
		Exec:      doRoot,
		FFOptions: []ff.Option{ff.WithEnvVarPrefix(name)},
	}
	if err := root.Parse(args); err != nil {
		return fmt.Errorf("parse error: %w", err)
	}

	// init runtime
	{
		// prng
		rand.Seed(srand.Fast())

		// concurrency
		runtime.GOMAXPROCS(1)

		// logger
		config := zapconfig.New().SetPreset("light-console")
		if opts.Debug {
			config = config.SetLevel(zapcore.DebugLevel)
		} else {
			config = config.SetLevel(zapcore.InfoLevel)
		}
		var err error
		opts.rootLogger, err = config.Build()
		if err != nil {
			return fmt.Errorf("logger init: %w", err)
		}
	}

	// run
	if err := root.Run(context.Background()); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

func doRoot(ctx context.Context, args []string) error { // nolint:gocognit
	logger := opts.rootLogger.Named("app")
	logger.Debug("init", zap.Strings("args", args), zap.Any("opts", opts))

	if len(args) > 0 {
		return flag.ErrHelp
	}

	if opts.BertyNodeAddr == "" {
		// FIXME: implement inmem bot.
		return fmt.Errorf("missing --berty-node-addr: %w", flag.ErrHelp)
	}

	fmt.Print(motd.Default())

	var g run.Group
	ctx, cancel := context.WithCancel(ctx)
	g.Add(func() error {
		<-ctx.Done()
		return nil
	}, func(err error) {
		logger.Info("Exiting...", zap.Error(err))
		cancel()
	})
	// signal handling
	g.Add(run.SignalHandler(ctx, syscall.SIGTERM, syscall.SIGINT, os.Interrupt, os.Kill))

	// berty bot
	g.Add(func() error {
		db, err := NewSqliteDB()
		if err != nil {
			return fmt.Errorf("db init: %w", err)
		}
		versionCommand := func(ctx bertybot.Context) {
			_ = ctx.ReplyString("version: " + bertyversion.Version)
		}

		cc, err := grpc.Dial(opts.BertyNodeAddr, grpc.WithInsecure())
		if err != nil {
			return fmt.Errorf("dial error: %w", err)
		}
		client := messengertypes.NewMessengerServiceClient(cc)

		newOpts := []bertybot.NewOption{}
		newOpts = append(newOpts,
			bertybot.WithLogger(logger.Named("berty")), // configure a logger
			bertybot.WithDisplayName("IOnlyTalkToETH"), // bot name
			bertybot.WithCommand("version", "show version", versionCommand),
			bertybot.WithRecipe(bertybot.AutoAcceptIncomingContactRequestRecipe()),
			bertybot.WithRecipe(bertybot.AutoAcceptIncomingGroupInviteRecipe()),
			bertybot.WithRecipe(bertybot.WelcomeMessageRecipe("Hello dear peroquet !")),
			bertybot.WithCommand("ping", "ping", Pingpong),

			bertybot.WithCommand("link-metamask", "link metamask", LinkMetamask(db)),
			bertybot.WithCommand("verify-metamask", "verify metamask", VerifyMetamask(db)),

			bertybot.WithCommand("clear", "clear", func(ctx bertybot.Context) {
				if ctx.IsReplay || !ctx.IsNew {
					return
				}

				for i := 0; i < 40; i++ {
					ctx.ReplyString("\n ")
				}
			}),

			bertybot.WithMessengerClient(client),
		)
		if opts.Debug {
			newOpts = append(newOpts, bertybot.WithRecipe(bertybot.DebugEventRecipe(logger.Named("debug"))))
		}

		bot, err := bertybot.New(newOpts...)
		if err != nil {
			return fmt.Errorf("bot initialization failed: %w", err)
		}
		logger.Info("retrieve instance Berty ID",
			zap.String("pk", bot.PublicKey()),
			zap.String("link", bot.BertyIDURL()),
		)
		if opts.Debug {
			qrterminal.GenerateHalfBlock(bot.BertyIDURL(), qrterminal.L, os.Stdout)
		}

		return bot.Start(ctx)
	}, func(error) {})

	logger.Info("Starting...")
	return g.Run()
}
