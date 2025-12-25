package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/ipfs/go-datastore"
	"github.com/post-quantumqoin/address"
	"github.com/post-quantumqoin/qoin-shor/api/v1api"
	"github.com/post-quantumqoin/qoin-shor/core/gen/slashfilter"
	"github.com/post-quantumqoin/qoin-shor/core/types"
	"github.com/post-quantumqoin/qoin-shor/journal"
	"github.com/post-quantumqoin/qoin-shor/journal/fsjournal"
	"github.com/post-quantumqoin/qoin-shor/node/repo"
	"github.com/urfave/cli/v2"
	"golang.org/x/xerrors"

	lapi "github.com/post-quantumqoin/qoin-shor/api"
	pqcminer "github.com/post-quantumqoin/qoin-shor/miner"
)

const (
	FlagMinerRepo   = "miner-repo"
	FlagMarketsRepo = "markets-repo"
)
const FlagMinerRepoDeprecation = "minerepo"

var minerCmd = &cli.Command{
	Name:  "miner",
	Usage: "Qoin Decentralized Quantum-Resistant AI Supercomputing Network miner",
	Subcommands: []*cli.Command{
		startCmd,
	},
}

var startCmd = &cli.Command{
	Name:  "start",
	Usage: "Start a running shor miner",
	Flags: []cli.Flag{
		&cli.DurationFlag{
			Name:  "timeout",
			Usage: "duration to wait till fail",
			Value: time.Second * 30,
		},
		&cli.StringFlag{
			Name:  "gas-premium",
			Usage: "set gas premium for initialization messages in AttoFIL",
			Value: "0",
		},
		&cli.StringFlag{
			Name:    FlagMinerRepo,
			Aliases: []string{FlagMinerRepoDeprecation},
			EnvVars: []string{"SHOR_MINER_PATH", "SHOR_STORAGE_PATH"},
			Value:   "~/.shorminer", // TODO: Consider XDG_DATA_HOME
			Usage:   fmt.Sprintf("Specify miner repo path. flag(%s) and env(SHOR_STORAGE_PATH) are DEPRECATION, will REMOVE SOON", FlagMinerRepoDeprecation),
		},
		&cli.StringFlag{
			Name:    "miner-address",
			Aliases: []string{"o"},
			Usage:   "miner key to use",
		},
	},

	Action: func(cctx *cli.Context) error {
		log.Info("Initializing shor miner")

		// ssize, err := abi.RegisteredSealProof_StackedDrg32GiBV1.SectorSize()
		// if err != nil {
		// 	return xerrors.Errorf("failed to calculate default sector size: %w", err)
		// }

		// if cctx.IsSet("sector-size") {
		// 	sectorSizeInt, err := units.RAMInBytes(cctx.String("sector-size"))
		// 	if err != nil {
		// 		return err
		// 	}
		// 	ssize = abi.SectorSize(sectorSizeInt)
		// }

		gasPrice, err := types.BigFromString(cctx.String("gas-premium"))
		if err != nil {
			return xerrors.Errorf("failed to parse gas-price flag: %s", err)
		}

		// symlink := cctx.Bool("symlink-imported-sectors")
		// if symlink {
		// 	log.Info("will attempt to symlink to imported sectors")
		// }

		ctx := ReqContext(cctx)

		// log.Info("Checking proof parameters")

		// if err := paramfetch.GetParams(ctx, build.ParametersJSON(), build.SrsJSON(), uint64(ssize)); err != nil {
		// 	return xerrors.Errorf("fetching proof parameters: %w", err)
		// }

		log.Info("Trying to connect to full node RPC")

		// if err := checkV1ApiSupport(ctx, cctx); err != nil {
		// 	return err
		// }

		api, closer, err := GetFullNodeAPIV1(cctx) // TODO: consider storing full node address in config
		if err != nil {
			return err
		}
		defer closer()

		log.Info("Checking full node sync status")

		// if !cctx.Bool("genesis-miner") && !cctx.Bool("nosync") {
		// 	if err := lcli.SyncWait(ctx, &v0api.WrapperV1Full{FullNode: api}, false); err != nil {
		// 		return xerrors.Errorf("sync wait: %w", err)
		// 	}
		// }

		log.Info("Checking if repo exists")

		// repoPath := cctx.String("repo")
		// r, err := repo.NewFS(repoPath)
		// if err != nil {
		// 	return err
		// }

		// ok, err := r.Exists()
		// if err != nil {
		// 	return err
		// }
		// if ok {
		// 	return xerrors.Errorf("repo at '%s' is already initialized", cctx.String(FlagMinerRepo))
		// }

		r, err := repo.NewFS(cctx.String(FlagMinerRepo))
		if err != nil {
			return xerrors.Errorf("opening fs repo: %w", err)
		}

		if cctx.String("config") != "" {
			r.SetConfigPath(cctx.String("config"))
		}

		err = r.Init(repo.StorageMiner)
		if err != nil && err != repo.ErrRepoExists {
			return xerrors.Errorf("repo init error: %w", err)
		}

		log.Info("Checking full node version")

		v, err := api.Version(ctx)
		if err != nil {
			return err
		}

		if !v.APIVersion.EqMajorMinor(lapi.FullAPIVersion1) {
			return xerrors.Errorf("Remote API version didn't match (expected %s, remote %s)", lapi.FullAPIVersion1, v.APIVersion)
		}

		log.Info("Initializing repo")

		// if err := r.Init(repo.StorageMiner); err != nil {
		// 	return xerrors.Errorf("miner repo init error: %w", err)
		// }

		// {
		// 	lr, err := r.Lock(repo.StorageMiner)
		// 	if err != nil {
		// 		return err
		// 	}

		// 	var localPaths []storiface.LocalPath

		// 	if pssb := cctx.StringSlice("pre-sealed-sectors"); len(pssb) != 0 {
		// 		log.Infof("Setting up storage config with presealed sectors: %v", pssb)

		// 		for _, psp := range pssb {
		// 			psp, err := homedir.Expand(psp)
		// 			if err != nil {
		// 				return err
		// 			}
		// 			localPaths = append(localPaths, storiface.LocalPath{
		// 				Path: psp,
		// 			})
		// 		}
		// 	}

		// 	if !cctx.Bool("no-local-storage") {
		// 		b, err := json.MarshalIndent(&storiface.LocalStorageMeta{
		// 			ID:       storiface.ID(uuid.New().String()),
		// 			Weight:   10,
		// 			CanSeal:  true,
		// 			CanStore: true,
		// 		}, "", "  ")
		// 		if err != nil {
		// 			return xerrors.Errorf("marshaling storage config: %w", err)
		// 		}

		// 		if err := os.WriteFile(filepath.Join(lr.Path(), "sectorstore.json"), b, 0644); err != nil {
		// 			return xerrors.Errorf("persisting storage metadata (%s): %w", filepath.Join(lr.Path(), "sectorstore.json"), err)
		// 		}

		// 		localPaths = append(localPaths, storiface.LocalPath{
		// 			Path: lr.Path(),
		// 		})
		// 	}

		// 	if err := lr.SetStorage(func(sc *storiface.StorageConfig) {
		// 		sc.StoragePaths = append(sc.StoragePaths, localPaths...)
		// 	}); err != nil {
		// 		return xerrors.Errorf("set storage config: %w", err)
		// 	}

		// 	if err := lr.Close(); err != nil {
		// 		return err
		// 	}
		// }

		if err := minerInit(ctx, cctx, api, r, gasPrice); err != nil {
			log.Errorf("Failed to initialize shor-miner: %+v", err)
			// path, err := homedir.Expand(repoPath)
			// if err != nil {
			// 	return err
			// }
			// log.Infof("Cleaning up %s after attempt...", path)
			// if err := os.RemoveAll(path); err != nil {
			// 	log.Errorf("Failed to clean up failed storage repo: %s", err)
			// }
			return xerrors.Errorf("Shor miner init failed")
		}

		// TODO: Point to setting storage price, maybe do it interactively or something
		// log.Info("Miner successfully created, you can now start it with 'lotus-miner run'")

		return nil
	},
}

var stopCmd = &cli.Command{
	Name:  "stop",
	Usage: "Stop a running miner",
	Flags: []cli.Flag{},
	Action: func(cctx *cli.Context) error {
		api, closer, err := GetAPI(cctx)
		if err != nil {
			return err
		}
		defer closer()

		err = api.Shutdown(ReqContext(cctx))
		if err != nil {
			return err
		}

		return nil
	},
}
var getDeviceCountCmd = &cli.Command{
	Name:  "get-device",
	Usage: "Getting the number of devices that can work",
	Flags: []cli.Flag{},
	Action: func(cctx *cli.Context) error {
		api, closer, err := GetAPI(cctx)
		if err != nil {
			return err
		}
		defer closer()

		err = api.Shutdown(ReqContext(cctx))
		if err != nil {
			return err
		}

		return nil
	},
}

var configCountCmd = &cli.Command{
	Name:  "config",
	Usage: "Getting the number of devices that can work",
	Flags: []cli.Flag{},
	Action: func(cctx *cli.Context) error {
		api, closer, err := GetAPI(cctx)
		if err != nil {
			return err
		}
		defer closer()

		err = api.Shutdown(ReqContext(cctx))
		if err != nil {
			return err
		}

		return nil
	},
}

var setMinerCmd = &cli.Command{
	Name:  "set-miner",
	Usage: "Configuring the information a miner needs before proceeding with mining",
	Flags: []cli.Flag{},
	Action: func(cctx *cli.Context) error {
		api, closer, err := GetAPI(cctx)
		if err != nil {
			return err
		}
		defer closer()

		err = api.Shutdown(ReqContext(cctx))
		if err != nil {
			return err
		}

		return nil
	},
	//		MINER OPTIONS:
	//	  --mine            Enable mining
	//	  --minerthreads value      Number of CPU threads to use for mining (default: 8)
	//	  --minergpus value     List of GPUs to use for mining (e.g. '0,1' will use the first two GPUs found)
	//	  --autodag         Enable automatic DAG pregeneration
	//	  --etherbase value     Public address for block mining rewards (default = first account created) (default: "0")
	//	  --targetgaslimit value    Target gas limit sets the artificial target gas floor for the blocks to mine (default: "4712388")
	//	  --gasprice value      Minimal gas price to accept for mining a transactions (default: "20000000000")
	//	  --extradata value     Block extra data set by the miner (default = client version)
	//
	// 	GAS PRICE ORACLE OPTIONS:
	//   --gpomin value    Minimum suggested gas price (default: "20000000000")
	//   --gpomax value    Maximum suggested gas price (default: "500000000000")
	//   --gpofull value   Full block threshold for gas price calculation (%) (default: 80)
	//   --gpobasedown value   Suggested gas price base step down ratio (1/1000) (default: 10)
	//   --gpobaseup value Suggested gas price base step up ratio (1/1000) (default: 100)
	//   --gpobasecf value Suggested gas price base correction factor (%) (default: 110)
	// --gas-premium value
}

var getMinerCmd = &cli.Command{
	Name:  "get-miner",
	Usage: "Get the address to send rewards to miners",
	Flags: []cli.Flag{},
	Action: func(cctx *cli.Context) error {
		api, closer, err := GetAPI(cctx)
		if err != nil {
			return err
		}
		defer closer()

		err = api.Shutdown(ReqContext(cctx))
		if err != nil {
			return err
		}

		return nil
	},
}

func minerInit(ctx context.Context, cctx *cli.Context, api v1api.FullNode, r repo.Repo, gasPrice types.BigInt) error {
	lr, err := r.Lock(repo.StorageMiner)
	if err != nil {
		return err
	}
	defer lr.Close() //nolint:errcheck
	mds, err := lr.Datastore(ctx, "/metadata")
	if err != nil {
		return err

	}
	if ad := cctx.String("miner-address"); ad == "" {
		return xerrors.Errorf("address is null")
	}

	a, err := address.NewFromString(cctx.String("miner-address"))

	if err := mds.Put(ctx, datastore.NewKey("miner-address"), a.Bytes()); err != nil {
		return err
	}

	// mid, err := address.IDFromAddress(a)
	// if err != nil {
	// 	return xerrors.Errorf("getting id address: %w", err)
	// }

	// sa, err := modules.StorageAuth(ctx, api)
	// if err != nil {
	// 	return err
	// }

	// wsts := statestore.New(namespace.Wrap(mds, modules.WorkerCallsPrefix))
	// smsts := statestore.New(namespace.Wrap(mds, modules.ManagerWorkPrefix))

	// si := paths.NewMemIndex(nil)

	// lstor, err := paths.NewLocal(ctx, lr, si, nil)
	// if err != nil {
	// 	return err
	// }
	// stor := paths.NewRemote(lstor, si, http.Header(sa), 10, &paths.DefaultPartialFileHandler{})

	// smgr, err := sealer.New(ctx, lstor, stor, lr, si, config.SealerConfig{
	// 	ParallelFetchLimit:       10,
	// 	AllowAddPiece:            true,
	// 	AllowPreCommit1:          true,
	// 	AllowPreCommit2:          true,
	// 	AllowCommit:              true,
	// 	AllowUnseal:              true,
	// 	AllowReplicaUpdate:       true,
	// 	AllowProveReplicaUpdate2: true,
	// 	AllowRegenSectorKey:      true,
	// }, config.ProvingConfig{}, wsts, smsts)
	// if err != nil {
	// 	return err
	// }

	// epp, err := storage.NewWinningPoStProver(api, smgr, ffiwrapper.ProofVerifier, dtypes.MinerID(mid))
	// if err != nil {
	// 	return err
	// }

	j, err := fsjournal.OpenFSJournal(lr, journal.EnvDisabledEvents())
	if err != nil {
		return fmt.Errorf("failed to open filesystem journal: %w", err)
	}
	m := pqcminer.NewPqcMiner(api, a, slashfilter.New(mds), j)
	{
		if err := m.Start(ctx); err != nil {
			return xerrors.Errorf("failed to start up genesis miner: %w", err)
		}

		if err := m.Stop(ctx); err != nil {
			log.Error("failed to shut down miner: ", err)
		}

	}
	return nil
}

// func storageMinerInit(ctx context.Context, cctx *cli.Context, api v1api.FullNode, r repo.Repo, ssize abi.SectorSize, gasPrice types.BigInt) error {
// 	lr, err := r.Lock(repo.StorageMiner)
// 	if err != nil {
// 		return err
// 	}
// 	defer lr.Close() //nolint:errcheck

// 	log.Info("Initializing libp2p identity")

// 	// p2pSk, err := makeHostKey(lr)
// 	// if err != nil {
// 	// 	return xerrors.Errorf("make host key: %w", err)
// 	// }

// 	// peerid, err := peer.IDFromPrivateKey(p2pSk)
// 	// if err != nil {
// 	// 	return xerrors.Errorf("peer ID from private key: %w", err)
// 	// }

// 	mds, err := lr.Datastore(ctx, "/metadata")
// 	if err != nil {
// 		return err
// 	}

// 	var addr address.Address
// 	if act := cctx.String("actor"); act != "" {
// 		a, err := address.NewFromString(act)
// 		if err != nil {
// 			return xerrors.Errorf("failed parsing actor flag value (%q): %w", act, err)
// 		}

// 		if cctx.Bool("genesis-miner") {
// 			if err := mds.Put(ctx, datastore.NewKey("miner-address"), a.Bytes()); err != nil {
// 				return err
// 			}

// 			mid, err := address.IDFromAddress(a)
// 			if err != nil {
// 				return xerrors.Errorf("getting id address: %w", err)
// 			}

// 			sa, err := modules.StorageAuth(ctx, api)
// 			if err != nil {
// 				return err
// 			}

// 			wsts := statestore.New(namespace.Wrap(mds, modules.WorkerCallsPrefix))
// 			smsts := statestore.New(namespace.Wrap(mds, modules.ManagerWorkPrefix))

// 			si := paths.NewMemIndex(nil)

// 			lstor, err := paths.NewLocal(ctx, lr, si, nil)
// 			if err != nil {
// 				return err
// 			}
// 			stor := paths.NewRemote(lstor, si, http.Header(sa), 10, &paths.DefaultPartialFileHandler{})

// 			smgr, err := sealer.New(ctx, lstor, stor, lr, si, config.SealerConfig{
// 				ParallelFetchLimit:       10,
// 				AllowAddPiece:            true,
// 				AllowPreCommit1:          true,
// 				AllowPreCommit2:          true,
// 				AllowCommit:              true,
// 				AllowUnseal:              true,
// 				AllowReplicaUpdate:       true,
// 				AllowProveReplicaUpdate2: true,
// 				AllowRegenSectorKey:      true,
// 			}, config.ProvingConfig{}, wsts, smsts)
// 			if err != nil {
// 				return err
// 			}
// 			epp, err := storage.NewWinningPoStProver(api, smgr, ffiwrapper.ProofVerifier, dtypes.MinerID(mid))
// 			if err != nil {
// 				return err
// 			}

// 			j, err := fsjournal.OpenFSJournal(lr, journal.EnvDisabledEvents())
// 			if err != nil {
// 				return fmt.Errorf("failed to open filesystem journal: %w", err)
// 			}

// 			m := storageminer.NewMiner(api, epp, a, slashfilter.New(mds), j)
// 			{
// 				if err := m.Start(ctx); err != nil {
// 					return xerrors.Errorf("failed to start up genesis miner: %w", err)
// 				}

// 				cerr := configureStorageMiner(ctx, api, a, peerid, gasPrice)

// 				if err := m.Stop(ctx); err != nil {
// 					log.Error("failed to shut down miner: ", err)
// 				}

// 				if cerr != nil {
// 					return xerrors.Errorf("failed to configure miner: %w", cerr)
// 				}
// 			}

// 			if pssb := cctx.String("pre-sealed-metadata"); pssb != "" {
// 				pssb, err := homedir.Expand(pssb)
// 				if err != nil {
// 					return err
// 				}

// 				log.Infof("Importing pre-sealed sector metadata for %s", a)

// 				if err := migratePreSealMeta(ctx, api, pssb, a, mds); err != nil {
// 					return xerrors.Errorf("migrating presealed sector metadata: %w", err)
// 				}
// 			}

// 			return nil
// 		}

// 		if pssb := cctx.String("pre-sealed-metadata"); pssb != "" {
// 			pssb, err := homedir.Expand(pssb)
// 			if err != nil {
// 				return err
// 			}

// 			log.Infof("Importing pre-sealed sector metadata for %s", a)

// 			if err := migratePreSealMeta(ctx, api, pssb, a, mds); err != nil {
// 				return xerrors.Errorf("migrating presealed sector metadata: %w", err)
// 			}
// 		}

// 		if err := configureStorageMiner(ctx, api, a, peerid, gasPrice); err != nil {
// 			return xerrors.Errorf("failed to configure miner: %w", err)
// 		}

// 		addr = a
// 	} else {
// 		a, err := createStorageMiner(ctx, api, ssize, peerid, gasPrice, cctx)
// 		if err != nil {
// 			return xerrors.Errorf("creating miner failed: %w", err)
// 		}

// 		addr = a
// 	}

// 	log.Infof("Created new miner: %s", addr)
// 	if err := mds.Put(ctx, datastore.NewKey("miner-address"), addr.Bytes()); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func createStorageMiner(ctx context.Context, api v1api.FullNode, ssize abi.SectorSize, peerid peer.ID, gasPrice types.BigInt, cctx *cli.Context) (address.Address, error) {
// 	var err error
// 	var owner address.Address
// 	if cctx.String("owner") != "" {
// 		owner, err = address.NewFromString(cctx.String("owner"))
// 	} else {
// 		owner, err = api.WalletDefaultAddress(ctx)
// 	}
// 	if err != nil {
// 		return address.Undef, err
// 	}

// 	worker := owner
// 	if cctx.String("worker") != "" {
// 		worker, err = address.NewFromString(cctx.String("worker"))
// 	} else if cctx.Bool("create-worker-key") { // TODO: Do we need to force this if owner is Secpk?
// 		worker, err = api.WalletNew(ctx, types.KTBLS)
// 	}
// 	if err != nil {
// 		return address.Address{}, err
// 	}

// 	sender := owner
// 	if fromstr := cctx.String("from"); fromstr != "" {
// 		faddr, err := address.NewFromString(fromstr)
// 		if err != nil {
// 			return address.Undef, fmt.Errorf("could not parse from address: %w", err)
// 		}
// 		sender = faddr
// 	}

// 	// make sure the sender account exists on chain
// 	_, err = api.StateLookupID(ctx, owner, types.EmptyTSK)
// 	if err != nil {
// 		return address.Undef, xerrors.Errorf("sender must exist on chain: %w", err)
// 	}

// 	// make sure the worker account exists on chain
// 	_, err = api.StateLookupID(ctx, worker, types.EmptyTSK)
// 	if err != nil {
// 		signed, err := api.MpoolPushMessage(ctx, &types.Message{
// 			From:  sender,
// 			To:    worker,
// 			Value: types.NewInt(0),
// 		}, nil)
// 		if err != nil {
// 			return address.Undef, xerrors.Errorf("push worker init: %w", err)
// 		}

// 		log.Infof("Initializing worker account %s, message: %s", worker, signed.Cid())
// 		log.Infof("Waiting for confirmation")

// 		mw, err := api.StateWaitMsg(ctx, signed.Cid(), build.MessageConfidence, lapi.LookbackNoLimit, true)
// 		if err != nil {
// 			return address.Undef, xerrors.Errorf("waiting for worker init: %w", err)
// 		}
// 		if mw.Receipt.ExitCode != 0 {
// 			return address.Undef, xerrors.Errorf("initializing worker account failed: exit code %d", mw.Receipt.ExitCode)
// 		}
// 	}

// 	// make sure the owner account exists on chain
// 	_, err = api.StateLookupID(ctx, owner, types.EmptyTSK)
// 	if err != nil {
// 		signed, err := api.MpoolPushMessage(ctx, &types.Message{
// 			From:  sender,
// 			To:    owner,
// 			Value: types.NewInt(0),
// 		}, nil)
// 		if err != nil {
// 			return address.Undef, xerrors.Errorf("push owner init: %w", err)
// 		}

// 		log.Infof("Initializing owner account %s, message: %s", worker, signed.Cid())
// 		log.Infof("Waiting for confirmation")

// 		mw, err := api.StateWaitMsg(ctx, signed.Cid(), build.MessageConfidence, lapi.LookbackNoLimit, true)
// 		if err != nil {
// 			return address.Undef, xerrors.Errorf("waiting for owner init: %w", err)
// 		}
// 		if mw.Receipt.ExitCode != 0 {
// 			return address.Undef, xerrors.Errorf("initializing owner account failed: exit code %d", mw.Receipt.ExitCode)
// 		}
// 	}

// 	// Note: the correct thing to do would be to call SealProofTypeFromSectorSize if actors version is v3 or later, but this still works
// 	nv, err := api.StateNetworkVersion(ctx, types.EmptyTSK)
// 	if err != nil {
// 		return address.Undef, xerrors.Errorf("failed to get network version: %w", err)
// 	}
// 	spt, err := miner.WindowPoStProofTypeFromSectorSize(ssize, nv)
// 	if err != nil {
// 		return address.Undef, xerrors.Errorf("getting post proof type: %w", err)
// 	}

// 	params, err := actors.SerializeParams(&power6.CreateMinerParams{
// 		Owner:               owner,
// 		Worker:              worker,
// 		WindowPoStProofType: spt,
// 		Peer:                abi.PeerID(peerid),
// 	})
// 	if err != nil {
// 		return address.Undef, err
// 	}

// 	createStorageMinerMsg := &types.Message{
// 		To:    power.Address,
// 		From:  sender,
// 		Value: big.Zero(),

// 		Method: power.Methods.CreateMiner,
// 		Params: params,

// 		GasLimit:   0,
// 		GasPremium: gasPrice,
// 	}

// 	signed, err := api.MpoolPushMessage(ctx, createStorageMinerMsg, nil)
// 	if err != nil {
// 		return address.Undef, xerrors.Errorf("pushing createMiner message: %w", err)
// 	}

// 	log.Infof("Pushed CreateMiner message: %s", signed.Cid())
// 	log.Infof("Waiting for confirmation")

// 	mw, err := api.StateWaitMsg(ctx, signed.Cid(), build.MessageConfidence, lapi.LookbackNoLimit, true)
// 	if err != nil {
// 		return address.Undef, xerrors.Errorf("waiting for createMiner message: %w", err)
// 	}

// 	if mw.Receipt.ExitCode != 0 {
// 		return address.Undef, xerrors.Errorf("create miner failed: exit code %d", mw.Receipt.ExitCode)
// 	}

// 	var retval power2.CreateMinerReturn
// 	if err := retval.UnmarshalCBOR(bytes.NewReader(mw.Receipt.Return)); err != nil {
// 		return address.Undef, err
// 	}

// 	log.Infof("New miners address is: %s (%s)", retval.IDAddress, retval.RobustAddress)
// 	return retval.IDAddress, nil
// }
