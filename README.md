# onamae-ddns

## Overview

When it was executed, it will update the A record for the specified domain on お名前.com to the current global IP address of the executing machine.

By running this periodically via cron, etc., DDNS can be made available by automatically updating the A record when the IP address changes.

## Usage

1. Update `config.json` with reference to `config.json.sample`

1. Run `go build` in `./onamae-ddns` and generate executable binary `onamae-ddns`

1. Set to run `onamae-ddns` by cron

## Notes

If you build with `go build`, the contents of `config.json` will be included in the executable binary. You would better remove `config.json` for security reasons.