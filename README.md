# roboprof

Profile collection tool for Go programs. (WIP)

## Usage

Initialize and start the collector with a FS storage:

```go
c, err := collector.Start(
    collector.LogMode(collector.StdoutLog),
    collector.CollectionMode(collector.CollectionSerial),
    collector.WithCPUProfile(10*time.Second),
    collector.WithMemProfile(),
    collector.WithBlockProfile(10*time.Second),
    collector.WithGoroutineProfile(),
    collector.WithStorageConf(collector.StorageConfig{
        FSStorageConfig: collector.FSStorageConfig{
            Dir: "./profiles",
        },
    }),
    collector.WithTickInterval(time.Second*30),
)
if err != nil {
    panic(err)
}
defer c.Stop()
```

The collector will write these profiles to the `./profiles` directory.

## Status

In development.

## License

[Apache 2.0](LICENSE)