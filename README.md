# lightOtp

## What?
lightOtp is a command line TOTP code generation tool built with Go. It's cross platform and lightweight.

## Supported platforms
* OSX (theoretically, this is untested)
* Windows
* Linux and Unix 
  * Optional: `xclip` or `xsel` for clipboard support ([atotto/clipboard](https://github.com/atotto/clipboard))
  
## Building from source

```
go get github.com/codemicro/lightOtp
go build github.com/codemicro/lightOtp/cmd/lightOtp
```

## Exporting tokens from Authy

[https://gist.github.com/gboudreau/94bb0c11a6209c82418d01a59d958c93](https://gist.github.com/gboudreau/94bb0c11a6209c82418d01a59d958c93)
