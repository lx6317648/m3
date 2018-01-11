	"github.com/m3db/m3db/persist/fs/commitlog"
		logger.Fatalf("could not parse new file mode: %v", err)
		logger.Fatalf("could not parse new directory mode: %v", err)
		logger.Fatalf("unknown commit log queue size type: %v",
		SetStrategy(commitlog.StrategyWriteBehind).
	// Set the series cache policy
	seriesCachePolicy := cfg.Cache.SeriesConfiguration().Policy
	opts = opts.SetSeriesCachePolicy(seriesCachePolicy)

	switch seriesCachePolicy {
	case series.CacheAll:
		// No options needed to be set
	default:
		// All other caching strategies require retrieving series from disk
		// to service a cache miss
		retrieverOpts := fs.NewBlockRetrieverOptions().
			SetBytesPool(opts.BytesPool()).
			SetSegmentReaderPool(opts.SegmentReaderPool())
		if blockRetrieveCfg := cfg.BlockRetrieve; blockRetrieveCfg != nil {
			retrieverOpts = retrieverOpts.
				SetFetchConcurrency(blockRetrieveCfg.FetchConcurrency)
		}
		blockRetrieverMgr := block.NewDatabaseBlockRetrieverManager(
			func(md namespace.Metadata) (block.DatabaseBlockRetriever, error) {
				retriever := fs.NewBlockRetriever(retrieverOpts, fsopts)
				if err := retriever.Open(md); err != nil {
					return nil, err
				}
				return retriever, nil
			})
		opts = opts.SetDatabaseBlockRetrieverManager(blockRetrieverMgr)
	}

			SetHashGen(sharding.NewHashGenWithSeed(cfg.Hashing.Seed))
	bs, err := cfg.Bootstrap.New(opts, m3dbClient)
			updated, err := cfg.Bootstrap.New(opts, m3dbClient)
		logger.Infof("bytes pool registering bucket capacity=%d, size=%d, "+
	logger.Infof("bytes pool %s init", policy.Type)