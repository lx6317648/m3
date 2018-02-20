	"github.com/m3db/m3db/environment"
	"github.com/m3db/m3db/x/mmap"
	"github.com/m3db/m3db/x/xio"
	"github.com/m3db/m3x/context"
	"github.com/m3db/m3x/ident"
	mmapCfg := cfg.Filesystem.MmapConfiguration()
	shouldUseHugeTLB := mmapCfg.HugeTLB.Enabled
	if shouldUseHugeTLB {
		// Make sure the host supports HugeTLB before proceeding with it to prevent
		// excessive log spam.
		shouldUseHugeTLB, err = hostSupportsHugeTLB()
		if err != nil {
			logger.Fatalf("could not determine if host supports HugeTLB: %v", err)
		}
		if !shouldUseHugeTLB {
			logger.Warnf("host doesn't support HugeTLB, proceeding without it")
		}
	}
		SetMmapEnableHugeTLB(shouldUseHugeTLB).
		SetMmapHugeTLBThreshold(mmapCfg.HugeTLB.Threshold).
	// Apply pooling options
	opts = withEncodingAndPoolingOptions(logger, opts, cfg.PoolingPolicy)

	// Setup the block retriever
			SetSegmentReaderPool(opts.SegmentReaderPool()).
			SetIdentifierPool(opts.IdentifierPool())
	hostID, err := cfg.HostID.Resolve()
	if err != nil {
		logger.Fatalf("could not resolve local host ID: %v", err)
	}

		envCfg environment.ConfigureResults

	case cfg.EnvironmentConfig.Service != nil:
		envCfg, err = cfg.EnvironmentConfig.Configure(environment.ConfigurationParameters{
			InstrumentOpts: iopts,
			HashingSeed:    cfg.Hashing.Seed,
		})
			logger.Fatalf("could not initialize dynamic config: %v", err)
	case cfg.EnvironmentConfig.Static != nil:
		envCfg, err = cfg.EnvironmentConfig.Configure(environment.ConfigurationParameters{
			HostID: hostID,
		})
			logger.Fatalf("could not initialize static config: %v", err)
	opts = opts.SetNamespaceInitializer(envCfg.NamespaceInitializer)
	topo, err := envCfg.TopologyInitializer.Init()
		logger.Fatalf("could not initialize m3db topology: %v", err)
			TopologyInitializer: envCfg.TopologyInitializer,
	kvWatchBootstrappers(envCfg.KVStore, logger, timeout, cfg.Bootstrap.Bootstrappers,
	db, err := cluster.NewDatabase(hostID, envCfg.TopologyInitializer, opts)
		kvWatchNewSeriesLimitPerShard(envCfg.KVStore, logger, topo,
	segmentReaderPool := xio.NewSegmentReaderPool(
	var identifierPool ident.Pool
		identifierPool = ident.NewPool(
		identifierPool = ident.NewNativePool(
func hostSupportsHugeTLB() (bool, error) {
	// Try and determine if the host supports HugeTLB in the first place
	withHugeTLB, err := mmap.Bytes(10, mmap.Options{
		HugeTLB: mmap.HugeTLBOptions{
			Enabled:   true,
			Threshold: 0,
		},
	})
		return false, fmt.Errorf("could not mmap anonymous region: %v", err)
	defer mmap.Munmap(withHugeTLB.Result)
	if withHugeTLB.Warning == nil {
		// If there was no warning, then the host didn't complain about
		// usa of huge TLB
		return true, nil

	// If we got a warning, try mmap'ing without HugeTLB
	withoutHugeTLB, err := mmap.Bytes(10, mmap.Options{})
		return false, fmt.Errorf("could not mmap anonymous region: %v", err)
	defer mmap.Munmap(withoutHugeTLB.Result)
	if withoutHugeTLB.Warning == nil {
		// The machine doesn't support HugeTLB, proceed without it
		return false, nil
	}
	// The warning was probably caused by something else, proceed using HugeTLB
	return true, nil